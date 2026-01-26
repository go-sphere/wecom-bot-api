package wecombot

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// Handler 处理接收到的回调消息
// ctx: 上下文，用于超时控制和取消
// callback: 解密后的回调消息
// 返回: 被动回复消息或错误
type Handler func(ctx context.Context, callback *Callback) (*PassiveReply, error)

// Server 处理来自企业微信的HTTP回调
// 实现http.Handler接口，可直接用于http.Handle
// 支持GET请求的URL验证和POST请求的消息回调
type Server struct {
	crypto   *Crypto // 加解密处理器
	handler  Handler // 业务处理函数
	deduper  Deduper // 消息去重器
	maxBytes int64   // 最大请求体大小
}

// ServerOption configures the server
type ServerOption func(*Server)

// WithDeduper sets a custom deduplicator
func WithDeduper(d Deduper) ServerOption {
	return func(s *Server) {
		s.deduper = d
	}
}

// WithMaxBytes sets the maximum request body size
func WithMaxBytes(n int64) ServerOption {
	return func(s *Server) {
		s.maxBytes = n
	}
}

// NewServer 创建一个新的回调服务器
// config: 包含Token、AESKey和ReceiveID的配置
// handler: 处理回调消息的业务函数
// opts: 可选配置项（去重器、最大请求体大小等）
// 返回: Server实例或错误
func NewServer(config Config, handler Handler, opts ...ServerOption) (*Server, error) {
	slog.Info("Creating new server", "receive_id", config.ReceiveID)

	crypto, err := NewCrypto(config)
	if err != nil {
		slog.Error("Failed to create crypto", "error", err)
		return nil, err
	}

	s := &Server{
		crypto:   crypto,
		handler:  handler,
		maxBytes: 10 * 1024 * 1024, // 10MB default
	}

	for _, opt := range opts {
		opt(s)
	}

	slog.Info("Server created successfully", "max_bytes", s.maxBytes)
	return s, nil
}

// ServeHTTP 实现http.Handler接口
// GET请求：处理URL验证（配置回调URL时）
// POST请求：处理消息回调
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	slog.Debug("Request received", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)

	if r.Method == http.MethodGet {
		s.handleURLVerification(w, r)
		return
	}

	if r.Method != http.MethodPost {
		slog.Warn("Method not allowed", "method", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s.handleCallback(ctx, w, r)
}

// handleURLVerification 处理URL验证的GET请求
// 企业微信在配置回调URL时会发送验证请求
// 需要在1秒内返回解密后的echostr明文
func (s *Server) handleURLVerification(w http.ResponseWriter, r *http.Request) {
	slog.Info("URL verification request received", "url", r.URL.String())
	query := r.URL.Query()
	msgSignature := query.Get("msg_signature")
	timestamp := query.Get("timestamp")
	nonce := query.Get("nonce")
	echoStr := query.Get("echostr")

	if msgSignature == "" || timestamp == "" || nonce == "" || echoStr == "" {
		slog.Warn("URL verification missing parameters",
			"has_msg_signature", msgSignature != "",
			"has_timestamp", timestamp != "",
			"has_nonce", nonce != "",
			"has_echostr", echoStr != "")
		http.Error(w, "missing parameters", http.StatusBadRequest)
		return
	}

	slog.Debug("Verifying URL", "echostr_length", len(echoStr))
	plaintext, err := s.crypto.VerifyURL(msgSignature, timestamp, nonce, echoStr)
	if err != nil {
		slog.Error("URL verification failed", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slog.Info("URL verification successful", "plaintext_length", len(plaintext))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(plaintext))
}

// handleCallback 处理消息回调的POST请求
// 流程：
// 1. 读取并限制请求体大小
// 2. 解析加密的JSON
// 3. 验证签名
// 4. 解密消息
// 5. 反序列化回调结构
// 6. 消息去重
// 7. 调用业务处理函数
// 8. 加密并返回回复消息
func (s *Server) handleCallback(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	msgSignature := query.Get("msg_signature")
	timestamp := query.Get("timestamp")
	nonce := query.Get("nonce")

	slog.Debug("Callback request received",
		"has_msg_signature", msgSignature != "",
		"has_timestamp", timestamp != "",
		"has_nonce", nonce != "")

	// Read and limit body size
	body, err := io.ReadAll(io.LimitReader(r.Body, s.maxBytes))
	if err != nil {
		slog.Error("Failed to read request body", "error", err)
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	slog.Debug("Request body read", "size", len(body))

	// Parse encrypted request
	var encReq struct {
		Encrypt string `json:"encrypt"`
	}
	if err := json.Unmarshal(body, &encReq); err != nil {
		slog.Error("Failed to parse encrypted request JSON", "error", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Verify signature
	if err := s.crypto.VerifySignature(msgSignature, timestamp, nonce, encReq.Encrypt); err != nil {
		slog.Error("Signature verification failed", "error", err)
		http.Error(w, "signature verification failed", http.StatusUnauthorized)
		return
	}

	slog.Debug("Signature verified successfully")

	// Decrypt message
	plaintext, err := s.crypto.Decrypt(encReq.Encrypt)
	if err != nil {
		slog.Error("Message decryption failed", "error", err)
		http.Error(w, "decryption failed", http.StatusBadRequest)
		return
	}

	slog.Debug("Message decrypted", "plaintext_length", len(plaintext))

	// Parse callback
	var callback Callback
	if err := json.Unmarshal(plaintext, &callback); err != nil {
		slog.Error("Failed to parse callback JSON", "error", err)
		http.Error(w, "invalid callback JSON", http.StatusBadRequest)
		return
	}

	slog.Info("Callback message received",
		"msg_id", callback.MsgID,
		"msg_type", callback.MsgType,
		"chat_type", callback.ChatType,
		"aibot_id", callback.AIBotID)

	// Deduplicate
	if s.deduper != nil {
		if s.deduper.IsDuplicate(callback.MsgID) {
			slog.Info("Duplicate message ignored", "msg_id", callback.MsgID)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	// Call handler
	reply, err := s.handler(ctx, &callback)
	if err != nil {
		slog.Error("Handler error", "error", err, "msg_id", callback.MsgID)
		http.Error(w, "handler error", http.StatusInternalServerError)
		return
	}

	// Send reply
	if reply == nil {
		slog.Debug("No reply to send", "msg_id", callback.MsgID)
		w.WriteHeader(http.StatusOK)
		return
	}

	slog.Debug("Sending reply", "msg_id", callback.MsgID, "reply_type", reply.MsgType)
	s.sendEncryptedReply(w, reply, nonce)
}

// sendEncryptedReply encrypts and sends the reply
func (s *Server) sendEncryptedReply(w http.ResponseWriter, reply *PassiveReply, nonce string) {
	// Marshal reply to JSON
	replyJSON, err := json.Marshal(reply)
	if err != nil {
		slog.Error("Failed to marshal reply", "error", err)
		http.Error(w, "failed to marshal reply", http.StatusInternalServerError)
		return
	}

	slog.Debug("Reply marshaled", "size", len(replyJSON))

	// Encrypt reply
	encrypted, err := s.crypto.Encrypt(replyJSON)
	if err != nil {
		slog.Error("Failed to encrypt reply", "error", err)
		http.Error(w, "failed to encrypt reply", http.StatusInternalServerError)
		return
	}

	// Calculate signature
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	msgSignature := s.crypto.calculateSignature(timestamp, nonce, encrypted)

	// Build response
	response := map[string]interface{}{
		"encrypt":      encrypted,
		"msgsignature": msgSignature,
		"timestamp":    timestamp,
		"nonce":        nonce,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("Failed to encode response", "error", err)
		return
	}

	slog.Info("Reply sent successfully", "reply_type", reply.MsgType)
}

// Deduper 消息去重器接口
// 用于防止因网络原因导致的重复回调
type Deduper interface {
	IsDuplicate(msgID string) bool // 检查消息ID是否重复
}

// MemoryDeduper 基于内存的消息去重器
// 使用TTL机制自动清理过期的消息ID
// 注意：在分布式环境中应使用Redis等分布式缓存
type MemoryDeduper struct {
	cache map[string]time.Time // 消息ID到过期时间的映射
	ttl   time.Duration        // 消息ID的生存时间
}

// NewMemoryDeduper 创建一个新的内存去重器
// ttl: 消息ID的生存时间，建议设置为2小时
func NewMemoryDeduper(ttl time.Duration) *MemoryDeduper {
	slog.Info("Creating memory deduper", "ttl", ttl)
	d := &MemoryDeduper{
		cache: make(map[string]time.Time),
		ttl:   ttl,
	}
	go d.cleanup()
	return d
}

// IsDuplicate checks if a message ID is duplicate
func (d *MemoryDeduper) IsDuplicate(msgID string) bool {
	now := time.Now()
	if expiry, exists := d.cache[msgID]; exists && now.Before(expiry) {
		return true
	}
	d.cache[msgID] = now.Add(d.ttl)
	return false
}

// cleanup periodically removes expired entries
func (d *MemoryDeduper) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		removed := 0
		for msgID, expiry := range d.cache {
			if now.After(expiry) {
				delete(d.cache, msgID)
				removed++
			}
		}
		if removed > 0 {
			slog.Debug("Deduper cleanup completed", "removed_count", removed, "remaining_count", len(d.cache))
		}
	}
}
