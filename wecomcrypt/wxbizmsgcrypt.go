package wecomcrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	ErrValidateSignature int = -40001
	ErrParseJSON         int = -40002
	ErrComputeSignature  int = -40003
	ErrIllegalAESKey     int = -40004
	ErrValidateCorpID    int = -40005
	ErrEncryptAES        int = -40006
	ErrDecryptAES        int = -40007
	ErrIllegalBuffer     int = -40008
	ErrEncodeBase64      int = -40009
	ErrDecodeBase64      int = -40010
	ErrGenJSON           int = -40011
	ErrIllegalProtocol   int = -40012
)

type ProtocolType int

const (
	JSONProtocol ProtocolType = 1
)

type CryptError struct {
	ErrCode int
	ErrMsg  string
}

func (e *CryptError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("crypt error %d: %s", e.ErrCode, e.ErrMsg)
}

func NewCryptError(errCode int, errMsg string) error {
	return &CryptError{ErrCode: errCode, ErrMsg: errMsg}
}

type WXBizJSONMessageRecv struct {
	ToUsername string `json:"tousername"`
	Encrypt    string `json:"encrypt"`
	AgentID    string `json:"agentid"`
}

type WXBizJSONMessageSend struct {
	Encrypt      string `json:"encrypt"`
	MsgSignature string `json:"msgsignature"`
	Timestamp    int    `json:"timestamp"`
	Nonce        string `json:"nonce"`
}

func NewWXBizJSONMessageSend(encrypt, signature string, timestamp int, nonce string) *WXBizJSONMessageSend {
	return &WXBizJSONMessageSend{Encrypt: encrypt, MsgSignature: signature, Timestamp: timestamp, Nonce: nonce}
}

type ProtocolProcessor interface {
	Parse(srcData []byte) (*WXBizJSONMessageRecv, error)
	Serialize(msgSend *WXBizJSONMessageSend) ([]byte, error)
}

type WXBizMsgCrypt struct {
	token             string
	encodingAESKey    string
	receiverID        string
	protocolProcessor ProtocolProcessor
}

type JsonProcessor struct{}

func (p *JsonProcessor) Parse(srcData []byte) (*WXBizJSONMessageRecv, error) {
	var msgRecv WXBizJSONMessageRecv
	err := json.Unmarshal(srcData, &msgRecv)
	if nil != err {
		fmt.Println("Unmarshal fail", err)
		return nil, NewCryptError(ErrParseJSON, "json to msg fail")
	}
	return &msgRecv, nil
}

func (p *JsonProcessor) Serialize(msgSend *WXBizJSONMessageSend) ([]byte, error) {
	jsonMsg, err := json.Marshal(msgSend)
	if nil != err {
		return nil, NewCryptError(ErrGenJSON, err.Error())
	}

	return jsonMsg, nil
}

func NewWXBizMsgCrypt(token, encodingAESKey, receiverID string, protocolType ProtocolType) (*WXBizMsgCrypt, error) {
	var protocolProcessor ProtocolProcessor
	if protocolType != JSONProtocol {
		return nil, NewCryptError(ErrIllegalProtocol, "protocol type not support")
	}
	protocolProcessor = new(JsonProcessor)
	return &WXBizMsgCrypt{token: token, encodingAESKey: encodingAESKey + "=", receiverID: receiverID, protocolProcessor: protocolProcessor}, nil
}

func (c *WXBizMsgCrypt) randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func (c *WXBizMsgCrypt) pkcs7Padding(plaintext string, blockSize int) []byte {
	padding := blockSize - (len(plaintext) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	var buffer bytes.Buffer
	buffer.WriteString(plaintext)
	buffer.Write(padtext)
	return buffer.Bytes()
}

func (c *WXBizMsgCrypt) pkcs7Unpadding(plaintext []byte, blockSize int) ([]byte, error) {
	plaintextLen := len(plaintext)
	if nil == plaintext || plaintextLen == 0 {
		return nil, NewCryptError(ErrDecryptAES, "pkcs7Unpadding error nil or zero")
	}
	if plaintextLen%blockSize != 0 {
		return nil, NewCryptError(ErrDecryptAES, "pkcs7Unpadding text not a multiple of the block size")
	}
	paddingLen := int(plaintext[plaintextLen-1])
	return plaintext[:plaintextLen-paddingLen], nil
}

func (c *WXBizMsgCrypt) cbcEncrypt(plaintext string) ([]byte, error) {
	aesKey, err := base64.StdEncoding.DecodeString(c.encodingAESKey)
	if nil != err {
		return nil, NewCryptError(ErrDecodeBase64, err.Error())
	}
	const blockSize = 32
	padMsg := c.pkcs7Padding(plaintext, blockSize)

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, NewCryptError(ErrEncryptAES, err.Error())
	}

	ciphertext := make([]byte, len(padMsg))
	iv := aesKey[:aes.BlockSize]

	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(ciphertext, padMsg)
	base64Msg := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(base64Msg, ciphertext)

	return base64Msg, nil
}

func (c *WXBizMsgCrypt) cbcDecrypt(base64EncryptMsg string) ([]byte, error) {
	aesKey, err := base64.StdEncoding.DecodeString(c.encodingAESKey)
	if nil != err {
		return nil, NewCryptError(ErrDecodeBase64, err.Error())
	}

	encryptMsg, err := base64.StdEncoding.DecodeString(base64EncryptMsg)
	if nil != err {
		return nil, NewCryptError(ErrDecodeBase64, err.Error())
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, NewCryptError(ErrDecryptAES, err.Error())
	}

	if len(encryptMsg) < aes.BlockSize {
		return nil, NewCryptError(ErrDecryptAES, "encrypt_msg size is not valid")
	}

	iv := aesKey[:aes.BlockSize]

	if len(encryptMsg)%aes.BlockSize != 0 {
		return nil, NewCryptError(ErrDecryptAES, "encrypt_msg not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(encryptMsg, encryptMsg)

	return encryptMsg, nil
}

func (c *WXBizMsgCrypt) cbcDecryptRaw(encryptMsg []byte) ([]byte, error) {
	aesKey, err := base64.StdEncoding.DecodeString(c.encodingAESKey)
	if nil != err {
		return nil, NewCryptError(ErrDecodeBase64, err.Error())
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, NewCryptError(ErrDecryptAES, err.Error())
	}

	if len(encryptMsg) < aes.BlockSize {
		return nil, NewCryptError(ErrDecryptAES, "encrypt_msg size is not valid")
	}

	iv := aesKey[:aes.BlockSize]

	if len(encryptMsg)%aes.BlockSize != 0 {
		return nil, NewCryptError(ErrDecryptAES, "encrypt_msg not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(encryptMsg))
	mode.CryptBlocks(plaintext, encryptMsg)
	return plaintext, nil
}

func (c *WXBizMsgCrypt) calcSignature(timestamp, nonce, data string) string {
	sortArr := []string{c.token, timestamp, nonce, data}
	sort.Strings(sortArr)
	var buffer bytes.Buffer
	for _, value := range sortArr {
		buffer.WriteString(value)
	}

	sha := sha1.New()
	sha.Write(buffer.Bytes())
	signature := fmt.Sprintf("%x", sha.Sum(nil))
	return string(signature)
}

func (c *WXBizMsgCrypt) ParsePlaintext(plaintext []byte) ([]byte, uint32, []byte, []byte, error) {
	const blockSize = 32
	plaintext, err := c.pkcs7Unpadding(plaintext, blockSize)
	if nil != err {
		return nil, 0, nil, nil, err
	}

	textLen := uint32(len(plaintext))
	if textLen < 20 {
		return nil, 0, nil, nil, NewCryptError(ErrIllegalBuffer, "plain is to small 1")
	}
	random := plaintext[:16]
	msg_len := binary.BigEndian.Uint32(plaintext[16:20])
	if textLen < (20 + msg_len) {
		return nil, 0, nil, nil, NewCryptError(ErrIllegalBuffer, "plain is to small 2")
	}

	msg := plaintext[20 : 20+msg_len]
	receiverID := plaintext[20+msg_len:]

	return random, msg_len, msg, receiverID, nil
}

func (c *WXBizMsgCrypt) VerifyURL(msgSignature, timestamp, nonce, echoStr string) ([]byte, error) {
	signature := c.calcSignature(timestamp, nonce, echoStr)

	if strings.Compare(signature, msgSignature) != 0 {
		return nil, NewCryptError(ErrValidateSignature, "signature not equal")
	}

	plaintext, err := c.cbcDecrypt(echoStr)
	if nil != err {
		return nil, err
	}

	_, _, msg, receiverID, err := c.ParsePlaintext(plaintext)
	if nil != err {
		return nil, err
	}

	if len(c.receiverID) > 0 && strings.Compare(string(receiverID), c.receiverID) != 0 {
		fmt.Println(string(receiverID), c.receiverID, len(receiverID), len(c.receiverID))
		return nil, NewCryptError(ErrValidateCorpID, "receiver_id is not equil")
	}

	return msg, nil
}

func (c *WXBizMsgCrypt) EncryptMessage(replyMsg, timestamp, nonce string) ([]byte, error) {
	timeInt, tErr := strconv.Atoi(timestamp)
	if tErr != nil {
		return nil, NewCryptError(ErrComputeSignature, "timestamp atoi fail")
	}
	randStr := c.randString(16)
	var buffer bytes.Buffer
	buffer.WriteString(randStr)

	msgLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLenBuf, uint32(len(replyMsg)))
	buffer.Write(msgLenBuf)
	buffer.WriteString(replyMsg)
	buffer.WriteString(c.receiverID)

	tmpCiphertext, err := c.cbcEncrypt(buffer.String())
	if nil != err {
		return nil, err
	}
	ciphertext := string(tmpCiphertext)

	signature := c.calcSignature(timestamp, nonce, ciphertext)

	msgSend := NewWXBizJSONMessageSend(ciphertext, signature, timeInt, nonce)
	return c.protocolProcessor.Serialize(msgSend)
}

func (c *WXBizMsgCrypt) DecryptMessage(msgSignature, timestamp, nonce string, postData []byte) ([]byte, error) {
	msgRecv, cryptErr := c.protocolProcessor.Parse(postData)
	if nil != cryptErr {
		return nil, cryptErr
	}

	signature := c.calcSignature(timestamp, nonce, msgRecv.Encrypt)

	if strings.Compare(signature, msgSignature) != 0 {
		return nil, NewCryptError(ErrValidateSignature, "signature not equal")
	}

	plaintext, cryptErr := c.cbcDecrypt(msgRecv.Encrypt)
	if nil != cryptErr {
		return nil, cryptErr
	}

	_, _, msg, receiverID, cryptErr := c.ParsePlaintext(plaintext)
	if nil != cryptErr {
		return nil, cryptErr
	}

	if len(c.receiverID) > 0 && strings.Compare(string(receiverID), c.receiverID) != 0 {
		return nil, NewCryptError(ErrValidateCorpID, "receiver_id is not equil")
	}

	return msg, nil
}

// DecryptFile decrypts encrypted file content downloaded from WeCom.
// The input should be raw encrypted bytes (not base64-encoded).
// And fuck wecom and wechat !!!
func (c *WXBizMsgCrypt) DecryptFile(encryptData []byte) ([]byte, error) {
	plaintext, cryptErr := c.cbcDecryptRaw(encryptData)
	if nil != cryptErr {
		return nil, cryptErr
	}

	const blockSize = 32
	plaintext, cryptErr = c.pkcs7Unpadding(plaintext, blockSize)
	if nil != cryptErr {
		return nil, cryptErr
	}

	return plaintext, nil
}
