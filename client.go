package wecombot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// Client 通过response_url处理主动回复
// 当收到回调消息后，可以使用response_url主动发送消息
// 注意：每个response_url只能使用一次，有效期为1小时
type Client struct {
	httpClient *http.Client // HTTP客户端，可自定义超时等配置
}

// ClientOption configures the client
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(c *http.Client) ClientOption {
	return func(client *Client) {
		client.httpClient = c
	}
}

// NewClient creates a new active reply client
func NewClient(opts ...ClientOption) *Client {
	slog.Info("Creating new client")

	c := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// MarkdownReply represents a markdown message
type MarkdownReply struct {
	MsgType  string          `json:"msgtype"`
	Markdown MarkdownContent `json:"markdown"`
}

// MarkdownContent represents markdown content
type MarkdownContent struct {
	Content  string    `json:"content"`
	Feedback *Feedback `json:"feedback,omitempty"`
}

// TemplateCardReply represents a template card message
type TemplateCardReply struct {
	MsgType      string        `json:"msgtype"`
	TemplateCard *TemplateCard `json:"template_card"`
}

// ReplyMarkdown 通过response_url发送Markdown消息
// 支持常见的Markdown格式：标题、列表、链接、图片、代码块、表格等
// 参数:
//
//	ctx: 上下文
//	responseURL: 回调消息中返回的response_url
//	content: Markdown内容，最长20480字节
//	feedbackID: 反馈ID，设置后用户反馈时会触发回调事件
func (c *Client) ReplyMarkdown(ctx context.Context, responseURL, content string, feedbackID string) error {
	slog.Info("Replying with markdown", "content_length", len(content), "has_feedback_id", feedbackID != "")

	reply := &MarkdownReply{
		MsgType: "markdown",
		Markdown: MarkdownContent{
			Content: content,
		},
	}

	if feedbackID != "" {
		reply.Markdown.Feedback = &Feedback{ID: feedbackID}
	}

	return c.sendRequest(ctx, responseURL, reply)
}

// ReplyTemplateCard 通过response_url发送模板卡片消息
// 注意：仅当回调的会话类型为单聊时支持
// 群聊中主动回复会引用触发回调的用户消息
func (c *Client) ReplyTemplateCard(ctx context.Context, responseURL string, card *TemplateCard) error {
	slog.Info("Replying with template card", "card_type", card.CardType)

	reply := &TemplateCardReply{
		MsgType:      "template_card",
		TemplateCard: card,
	}

	return c.sendRequest(ctx, responseURL, reply)
}

// sendRequest sends an HTTP POST request
func (c *Client) sendRequest(ctx context.Context, url string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("Failed to marshal payload", "error", err)
		return fmt.Errorf("%w: marshal failed: %v", ErrBadRequest, err)
	}

	slog.Debug("Sending request", "url", url, "payload_size", len(data))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return fmt.Errorf("%w: create request failed: %v", ErrBadRequest, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("Request failed", "error", err, "url", url)
		return fmt.Errorf("%w: %v", ErrResponseFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("Request returned non-OK status", "status", resp.StatusCode, "body", string(body))
		return fmt.Errorf("%w: status %d, body: %s", ErrResponseFailed, resp.StatusCode, string(body))
	}

	slog.Info("Request sent successfully", "url", url, "status", resp.StatusCode)
	return nil
}
