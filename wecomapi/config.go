package wecomapi

// Config 企业微信智能机器人的配置
// Token和AESKey可以在企业微信管理后台的机器人配置页面获取
type Config struct {
	Token     string `json:"token" yaml:"token"`           // 用于签名验证的Token
	AESKey    string `json:"aes_key" yaml:"aes_key"`       // Base64编码的AES密钥（43个字符）
	ReceiveID string `json:"receive_id" yaml:"receive_id"` // 接收者ID（企业内部机器人使用空字符串）
}
