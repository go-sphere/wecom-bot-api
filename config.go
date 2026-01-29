package wecombot

import "log/slog"

// Config 企业微信智能机器人的配置
// Token和AESKey可以在企业微信管理后台的机器人配置页面获取
type Config struct {
	Token     string `json:"token" yaml:"token"`           // 用于签名验证的Token
	AESKey    string `json:"aes_key" yaml:"aes_key"`       // Base64编码的AES密钥（43个字符）
	ReceiveID string `json:"receive_id" yaml:"receive_id"` // 接收者ID（企业内部机器人使用空字符串）
}

// Validate 检查配置是否有效
// 返回: 如果配置无效则返回错误
func (c *Config) Validate() error {
	if c.Token == "" {
		slog.Error("Config validation failed: token is empty")
		return ErrInvalidConfig
	}
	if len(c.AESKey) != 43 {
		slog.Error("Config validation failed: invalid AES key length", "expected", 43, "got", len(c.AESKey))
		return ErrInvalidAESKey
	}
	slog.Debug("Config validated successfully")
	return nil
}
