// Package wecomapi provides a complete SDK for building WeChat Work (企业微信) AI Bot applications.
//
// This package implements the full WeChat Work AI Bot API with support for:
//   - Message callbacks (text, image, mixed, voice, file)
//   - Event callbacks (enter chat, template card interactions, feedback)
//   - Encryption/decryption (AES-256-CBC with PKCS#7 padding)
//   - Signature verification (SHA1 with constant-time comparison)
//   - Template cards (all 5 types with full field support)
//   - Stream messages (for AI responses)
//   - Active replies via response_url
//   - Message deduplication
//
// # Quick Start
//
// Create a callback server:
//
//	config := wecombot.Config{
//	    Token:     "your_token",
//	    AESKey:    "your_43_char_aes_key",
//	    ReceiveID: "", // Empty for internal enterprise bots
//	}
//
//	server, err := wecombot.NewServer(config, func(ctx context.Context, callback *wecombot.Callback) (*wecombot.PassiveReply, error) {
//	    if callback.MsgType == "text" {
//	        return wecombot.NewTextReply("Echo: " + callback.Text.Content), nil
//	    }
//	    return wecombot.NewEmptyReply(), nil
//	})
//
//	http.Handle("/callback", server)
//	http.ListenAndServe(":8080", nil)
//
// Send active replies:
//
//	client := wecombot.NewClient()
//	err := client.ReplyMarkdown(ctx, responseURL, "# Hello\nMarkdown content", "")
//
// # Message Types
//
// The SDK supports all message types defined in the WeChat Work API:
//   - Text messages
//   - Image messages
//   - Mixed text/image messages
//   - Voice messages (with transcription)
//   - File messages
//   - Stream messages (for AI responses)
//
// # Template Cards
//
// All 5 template card types are fully supported:
//   - Text Notice (text_notice)
//   - News Notice (news_notice)
//   - Button Interaction (button_interaction)
//   - Vote Interaction (vote_interaction)
//   - Multiple Interaction (multiple_interaction)
//
// # Security
//
// The package implements the WeChat Work encryption scheme:
//   - AES-256-CBC encryption with PKCS#7 padding
//   - SHA1 signature verification with constant-time comparison
//   - URL parameter validation and decoding
//
// # Architecture
//
// The package is organized into several components:
//   - Config: Configuration for token, AES key, and receive ID
//   - Crypto: Encryption, decryption, and signature verification
//   - Server: HTTP handler for callbacks with automatic encryption/decryption
//   - Client: HTTP client for active replies via response_url
//   - Callback: Strongly-typed structures for all callback types
//   - TemplateCard: Strongly-typed structures for all template card types
//   - PassiveReply: Structures for passive replies to callbacks
//
// For more information, see the official WeChat Work documentation:
// https://developer.work.weixin.qq.com/document/path/100719
package wecomapi
