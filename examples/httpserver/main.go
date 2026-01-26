package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	wecombot "github.com/TBXark/wecom-bot-api"
	"github.com/gin-gonic/gin"
)

var (
	activeClient *wecombot.Client
)

func main() {
	token := flag.String("token", "", "WeCom bot token")
	aesKey := flag.String("aes-key", "", "WeCom AES key")
	addr := flag.String("addr", "0.0.0.0:9090", "HTTP server address")
	debug := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	// Set Gin mode
	if !*debug {
		gin.SetMode(gin.ReleaseMode)
	}

	slog.Info("Starting WeCom bot server")

	// Load configuration
	config := wecombot.Config{
		Token:     *token,
		AESKey:    *aesKey,
		ReceiveID: "", // Empty for internal enterprise bots
	}

	// Create server with deduplication
	deduper := wecombot.NewMemoryDeduper(2 * time.Hour)
	server, err := wecombot.NewServer(config, handleCallback,
		wecombot.WithDeduper(deduper),
		wecombot.WithMaxBytes(10*1024*1024),
	)
	if err != nil {
		slog.Error("Failed to create server", "error", err)
		return
	}

	// Create active reply client
	activeClient = wecombot.NewClient()

	// Setup Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", healthHandler)

	// WeCom callback endpoint
	r.Any("/callback", gin.WrapH(server))

	// Active reply endpoint (for demonstration)
	r.POST("/active-reply", activeReplyHandler)

	slog.Info("Server listening", "address", *addr)
	if err := r.Run(*addr); err != nil {
		slog.Error("Server failed", "error", err)
	}
}

// healthHandler handles health check requests
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Unix(),
	})
}

// activeReplyHandler handles active reply requests
func activeReplyHandler(c *gin.Context) {
	responseURL := c.Query("response_url")
	if responseURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing response_url parameter",
		})
		return
	}

	err := activeClient.ReplyMarkdown(
		c.Request.Context(),
		responseURL,
		"# Hello from active reply!\nThis is a **markdown** message.",
		"",
	)
	if err != nil {
		slog.Error("Failed to send active reply", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sent",
	})
}

// handleCallback processes incoming callbacks
func handleCallback(ctx context.Context, callback *wecombot.Callback) (*wecombot.PassiveReply, error) {
	slog.Info("Callback received",
		"msg_id", callback.MsgID,
		"msg_type", callback.MsgType,
		"chat_type", callback.ChatType,
		"aibot_id", callback.AIBotID,
		"user_id", callback.From.UserID)

	// Pretty print callback for debugging
	if data, err := json.MarshalIndent(callback, "", "  "); err == nil {
		slog.Debug("Callback details", "data", string(data))
	}

	switch callback.MsgType {
	case "text":
		return handleTextMessage(ctx, callback)
	case "image":
		return handleImageMessage(ctx, callback)
	case "stream":
		return handleStreamRefresh(ctx, callback)
	case "event":
		return handleEvent(ctx, callback)
	default:
		slog.Warn("Unknown message type", "msg_type", callback.MsgType)
		return wecombot.NewEmptyReply(), nil
	}
}

// handleTextMessage handles text messages
func handleTextMessage(ctx context.Context, callback *wecombot.Callback) (*wecombot.PassiveReply, error) {
	if callback.Text == nil {
		return wecombot.NewEmptyReply(), nil
	}

	content := callback.Text.Content
	slog.Info("Text message received", "content", content, "from_user", callback.From.UserID)

	// Example: Echo with stream reply
	streamID := fmt.Sprintf("stream_%s", callback.MsgID)
	return wecombot.NewStreamReply(streamID, fmt.Sprintf("You said: %s", content), false), nil
}

// handleImageMessage handles image messages
func handleImageMessage(ctx context.Context, callback *wecombot.Callback) (*wecombot.PassiveReply, error) {
	if callback.Image == nil {
		return wecombot.NewEmptyReply(), nil
	}

	slog.Info("Image message received", "url", callback.Image.URL, "from_user", callback.From.UserID)

	// Reply with a template card
	card := &wecombot.TemplateCard{
		CardType: "text_notice",
		MainTitle: &wecombot.MainTitle{
			Title: "Image Received",
			Desc:  "We've received your image",
		},
		SubTitleText: "Processing your image...",
		CardAction: &wecombot.CardAction{
			Type: 0,
		},
	}

	return wecombot.NewTemplateCardReply(card), nil
}

// handleStreamRefresh handles stream message refresh
func handleStreamRefresh(ctx context.Context, callback *wecombot.Callback) (*wecombot.PassiveReply, error) {
	if callback.Stream == nil {
		return wecombot.NewEmptyReply(), nil
	}

	streamID := callback.Stream.ID
	slog.Info("Stream refresh received", "stream_id", streamID)

	// Continue streaming (in real app, fetch from AI model)
	content := "This is a streaming response... "
	return &wecombot.PassiveReply{
		MsgType: "stream",
		Stream: &wecombot.StreamReply{
			ID:      streamID,
			Content: content,
			Finish:  false,
		},
	}, nil
}

// handleEvent handles event callbacks
func handleEvent(ctx context.Context, callback *wecombot.Callback) (*wecombot.PassiveReply, error) {
	if callback.Event == nil {
		return wecombot.NewEmptyReply(), nil
	}

	slog.Info("Event received", "event_type", callback.Event.EventType)

	switch callback.Event.EventType {
	case "enter_chat":
		return handleEnterChat(ctx, callback)
	case "template_card_event":
		return handleTemplateCardEvent(ctx, callback)
	case "feedback_event":
		return handleFeedbackEvent(ctx, callback)
	default:
		slog.Warn("Unknown event type", "event_type", callback.Event.EventType)
		return wecombot.NewEmptyReply(), nil
	}
}

// handleEnterChat handles user entering chat
func handleEnterChat(ctx context.Context, callback *wecombot.Callback) (*wecombot.PassiveReply, error) {
	slog.Info("User entered chat", "user_id", callback.From.UserID, "chat_id", callback.ChatID)

	// Send welcome message
	card := &wecombot.TemplateCard{
		CardType: "text_notice",
		Source: &wecombot.Source{
			IconURL: "https://example.com/icon.png",
			Desc:    "AI Assistant",
		},
		MainTitle: &wecombot.MainTitle{
			Title: "Welcome!",
			Desc:  "I'm your AI assistant. How can I help you today?",
		},
		JumpList: []wecombot.JumpAction{
			{
				Type:     3,
				Title:    "How to use this bot?",
				Question: "How do I use this AI assistant?",
			},
		},
		CardAction: &wecombot.CardAction{
			Type: 0,
		},
	}

	return wecombot.NewTemplateCardReply(card), nil
}

// handleTemplateCardEvent handles template card interactions
func handleTemplateCardEvent(ctx context.Context, callback *wecombot.Callback) (*wecombot.PassiveReply, error) {
	event := callback.Event.TemplateCardEvent
	if event == nil {
		return wecombot.NewEmptyReply(), nil
	}

	slog.Info("Template card event received",
		"card_type", event.CardType,
		"event_key", event.EventKey,
		"task_id", event.TaskID,
		"user_id", callback.From.UserID)

	// Update the card
	card := &wecombot.TemplateCard{
		CardType: "text_notice",
		MainTitle: &wecombot.MainTitle{
			Title: "Action Received",
			Desc:  fmt.Sprintf("You clicked: %s", event.EventKey),
		},
		TaskID: event.TaskID,
		CardAction: &wecombot.CardAction{
			Type: 0,
		},
	}

	return wecombot.NewUpdateTemplateCardReply(nil, card), nil
}

// handleFeedbackEvent handles user feedback
func handleFeedbackEvent(ctx context.Context, callback *wecombot.Callback) (*wecombot.PassiveReply, error) {
	event := callback.Event.FeedbackEvent
	if event == nil {
		return wecombot.NewEmptyReply(), nil
	}

	feedbackType := "unknown"
	switch event.Type {
	case 1:
		feedbackType = "accurate"
	case 2:
		feedbackType = "inaccurate"
	case 3:
		feedbackType = "cancelled"
	}

	slog.Info("Feedback received",
		"feedback_id", event.ID,
		"feedback_type", feedbackType,
		"content", event.Content,
		"user_id", callback.From.UserID)

	// Just acknowledge, no reply needed
	return wecombot.NewEmptyReply(), nil
}
