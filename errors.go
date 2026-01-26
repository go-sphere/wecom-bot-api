package wecombot

import "errors"

var (
	// 配置错误
	ErrInvalidConfig = errors.New("wecombot: 无效的配置")
	ErrInvalidAESKey = errors.New("wecombot: AES密钥必须是43个字符")

	// 加密错误
	ErrBadSignature = errors.New("wecombot: 签名验证失败")
	ErrDecrypt      = errors.New("wecombot: 解密失败")
	ErrEncrypt      = errors.New("wecombot: 加密失败")
	ErrInvalidPKCS7 = errors.New("wecombot: 无效的PKCS#7填充")

	// 请求错误
	ErrBadRequest    = errors.New("wecombot: 错误的请求")
	ErrDuplicate     = errors.New("wecombot: 重复的消息")
	ErrInvalidJSON   = errors.New("wecombot: 无效的JSON")
	ErrMissingParams = errors.New("wecombot: 缺少必需的参数")

	// 响应错误
	ErrResponseFailed = errors.New("wecombot: 响应失败")
	ErrTimeout        = errors.New("wecombot: 请求超时")
)
