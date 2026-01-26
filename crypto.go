package wecombot

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"sort"
	"strings"
)

// Crypto 处理加密、解密和签名验证
// 实现企业微信智能机器人的加解密方案
// 加密方式：AES-256-CBC，数据采用PKCS#7填充至32字节的倍数
// IV初始向量：取AESKey前16字节
type Crypto struct {
	token     string // 用于签名验证的Token
	aesKey    []byte // AES加密密钥（32字节）
	receiveID string // 接收者ID（企业内部机器人为空字符串）
}

// NewCrypto 创建一个新的Crypto实例
// config: 包含Token、AESKey和ReceiveID的配置
// 返回: Crypto实例或错误
func NewCrypto(config Config) (*Crypto, error) {
	if err := config.Validate(); err != nil {
		slog.Error("Config validation failed", "error", err)
		return nil, err
	}

	// EncodingAESKey is 43 characters (Base64 without padding)
	// Decode it to get 32 bytes key
	aesKey, err := base64.StdEncoding.DecodeString(config.AESKey + "=")
	if err != nil {
		slog.Error("Failed to decode AES key", "error", err)
		return nil, fmt.Errorf("%w: %v", ErrInvalidAESKey, err)
	}

	if len(aesKey) != 32 {
		slog.Error("Invalid AES key length", "expected", 32, "got", len(aesKey))
		return nil, fmt.Errorf("%w: decoded key must be 32 bytes, got %d", ErrInvalidAESKey, len(aesKey))
	}

	slog.Debug("Crypto instance created successfully")
	return &Crypto{
		token:     config.Token,
		aesKey:    aesKey,
		receiveID: config.ReceiveID,
	}, nil
}

// VerifyURL 验证URL回调的有效性
// 当配置智能机器人回调URL时，企业微信会发送验证请求
// 需要在1秒内响应，返回解密后的echostr明文
// 参数:
//
//	msgSignature: 企业微信加密签名
//	timestamp: 时间戳
//	nonce: 随机数（两小时内不重复）
//	echoStr: 加密的字符串（已被HTTP框架自动URL decode）
//
// 返回: 解密后的消息明文或错误
func (c *Crypto) VerifyURL(msgSignature, timestamp, nonce, echoStr string) (string, error) {
	// Verify signature
	if err := c.VerifySignature(msgSignature, timestamp, nonce, echoStr); err != nil {
		return "", err
	}

	// Decrypt echostr
	msg, err := c.Decrypt(echoStr)
	if err != nil {
		return "", err
	}

	return string(msg), nil
}

// VerifySignature 验证消息签名
// 使用SHA1算法对token、timestamp、nonce、encrypt进行排序后计算签名
// 使用常量时间比较防止时序攻击
func (c *Crypto) VerifySignature(msgSignature, timestamp, nonce, encrypt string) error {
	expected := c.calculateSignature(timestamp, nonce, encrypt)
	if !constantTimeCompare(expected, msgSignature) {
		slog.Warn("Signature verification failed", "expected", expected, "got", msgSignature)
		return ErrBadSignature
	}
	slog.Debug("Signature verified successfully")
	return nil
}

// calculateSignature computes SHA1 signature
func (c *Crypto) calculateSignature(timestamp, nonce, encrypt string) string {
	strs := []string{c.token, timestamp, nonce, encrypt}
	sort.Strings(strs)
	h := sha1.New()
	h.Write([]byte(strings.Join(strs, "")))
	return hex.EncodeToString(h.Sum(nil))
}

// Decrypt 解密加密的消息
// 消息格式：random(16字节) + msgLen(4字节) + msg + receiveID
// 加密方式：AES-256-CBC，PKCS#7填充至32字节的倍数（企业微信特殊实现）
// IV：AESKey的前16字节
func (c *Crypto) Decrypt(encryptedMsg string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedMsg)
	if err != nil {
		slog.Error("Failed to decode base64", "error", err)
		return nil, fmt.Errorf("%w: base64 decode failed: %v", ErrDecrypt, err)
	}

	slog.Debug("Decrypting message", "ciphertext_length", len(ciphertext))

	block, err := aes.NewCipher(c.aesKey)
	if err != nil {
		slog.Error("Failed to create cipher", "error", err)
		return nil, fmt.Errorf("%w: create cipher failed: %v", ErrDecrypt, err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		slog.Error("Invalid ciphertext length", "length", len(ciphertext), "block_size", aes.BlockSize)
		return nil, fmt.Errorf("%w: ciphertext not multiple of block size (len=%d)", ErrDecrypt, len(ciphertext))
	}

	// IV is first 16 bytes of key
	iv := c.aesKey[:aes.BlockSize]
	mode := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove PKCS#7 padding (32-byte block size for WeCom)
	plaintext, err = pkcs7Unpad(plaintext, 32)
	if err != nil {
		slog.Error("Failed to unpad plaintext", "error", err)
		return nil, fmt.Errorf("%w (ciphertext_len=%d, plaintext_len=%d)", err, len(ciphertext), len(plaintext))
	}

	// Parse: random(16) + msgLen(4) + msg + receiveID
	if len(plaintext) < 20 {
		slog.Error("Plaintext too short", "length", len(plaintext))
		return nil, fmt.Errorf("%w: plaintext too short", ErrDecrypt)
	}

	msgLen := binary.BigEndian.Uint32(plaintext[16:20])
	if len(plaintext) < int(20+msgLen) {
		slog.Error("Invalid message length", "msg_len", msgLen, "plaintext_len", len(plaintext))
		return nil, fmt.Errorf("%w: invalid message length", ErrDecrypt)
	}

	msg := plaintext[20 : 20+msgLen]
	slog.Debug("Message decrypted successfully", "msg_length", len(msg))
	return msg, nil
}

// Encrypt 加密消息用于回复
// 消息格式：random(16字节) + msgLen(4字节) + msg + receiveID
// 加密方式：AES-256-CBC，PKCS#7填充至32字节的倍数（企业微信特殊实现）
// IV：AESKey的前16字节
func (c *Crypto) Encrypt(msg []byte) (string, error) {
	slog.Debug("Encrypting message", "msg_length", len(msg))

	// Build: random(16) + msgLen(4) + msg + receiveID
	random := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, random); err != nil {
		slog.Error("Failed to generate random bytes", "error", err)
		return "", fmt.Errorf("%w: generate random failed: %v", ErrEncrypt, err)
	}

	buf := new(bytes.Buffer)
	buf.Write(random)

	msgLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLenBuf, uint32(len(msg)))
	buf.Write(msgLenBuf)
	buf.Write(msg)
	buf.WriteString(c.receiveID)

	// Apply PKCS#7 padding (32-byte block size for WeCom)
	plaintext := pkcs7Pad(buf.Bytes(), 32)

	block, err := aes.NewCipher(c.aesKey)
	if err != nil {
		slog.Error("Failed to create cipher", "error", err)
		return "", fmt.Errorf("%w: create cipher failed: %v", ErrEncrypt, err)
	}

	ciphertext := make([]byte, len(plaintext))
	iv := c.aesKey[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	encrypted := base64.StdEncoding.EncodeToString(ciphertext)
	slog.Debug("Message encrypted successfully", "encrypted_length", len(encrypted))
	return encrypted, nil
}

// pkcs7Pad applies PKCS#7 padding
// WeCom uses 32-byte block size instead of standard 16-byte AES block size
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7Unpad removes PKCS#7 padding
// blockSize parameter allows for WeCom's non-standard 32-byte block size
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, ErrInvalidPKCS7
	}

	padding := int(data[length-1])
	if padding > length || padding > blockSize {
		return nil, ErrInvalidPKCS7
	}

	for i := 0; i < padding; i++ {
		if data[length-1-i] != byte(padding) {
			return nil, ErrInvalidPKCS7
		}
	}

	return data[:length-padding], nil
}

// constantTimeCompare performs constant-time string comparison
func constantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
