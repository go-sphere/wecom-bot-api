# WeChat Work AI Bot SDK for Go

A complete Go SDK for building WeChat Work (企业微信) AI Bot applications with callbacks, encryption, template cards, and streaming messages.

## Features

- Complete API Coverage: All message types, events, and template cards
- Strong Typing: Type-safe structures for all API objects
- Encryption/Decryption: AES-256-CBC encryption with PKCS#7 padding
- Signature Verification: Constant-time signature comparison
- Active Replies: Client for sending messages via `response_url`
- Stream Support: Handle streaming AI responses
- Deduplication: Built-in message deduplication with TTL
- Template Cards: All 5 card types with full field support
- Zero Dependencies: Uses only Go standard library

## Documentation

- [API Documentation](API.md)
- [WeChat Work Official Docs](https://developer.work.weixin.qq.com/document/path/100719)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.