package wecomapi

// ReplyMsgType 被动回复的消息类型。
type ReplyMsgType string

const (
	ReplyMsgTypeText                   ReplyMsgType = "text"
	ReplyMsgTypeMarkdown               ReplyMsgType = "markdown"
	ReplyMsgTypeTemplateCard           ReplyMsgType = "template_card"
	ReplyMsgTypeStream                 ReplyMsgType = "stream"
	ReplyMsgTypeStreamWithTemplateCard ReplyMsgType = "stream_with_template_card"
)

// ReplyResponseType 被动回复的响应类型。
type ReplyResponseType string

const (
	ReplyResponseTypeUpdateTemplateCard ReplyResponseType = "update_template_card"
)

// PassiveReply 回调的被动回复体。
type PassiveReply struct {
	MsgType      ReplyMsgType      `json:"msgtype,omitempty"`
	Text         *Text             `json:"text,omitempty"`
	Markdown     *Markdown         `json:"markdown,omitempty"`
	TemplateCard *TemplateCard     `json:"template_card,omitempty"`
	Stream       *StreamReply      `json:"stream,omitempty"`
	ResponseType ReplyResponseType `json:"response_type,omitempty"`
	UserIDs      []string          `json:"userids,omitempty"`
}

// StreamReply 流式消息的回复体。
type StreamReply struct {
	ID       string    `json:"id,omitempty"`
	Finish   bool      `json:"finish,omitempty"`
	Content  string    `json:"content,omitempty"`
	MsgItem  []MsgItem `json:"msg_item,omitempty"`
	Feedback *Feedback `json:"feedback,omitempty"`
}

// ImageBase64 表示Base64编码图片，用于流式混排。
type ImageBase64 struct {
	Base64 string `json:"base64"` // 图片内容的Base64编码（编码前最大10M，支持JPG/PNG）
	MD5    string `json:"md5"`    // 图片内容（Base64编码前）的MD5值
}

// Markdown 表示Markdown消息内容。
type Markdown struct {
	Content string `json:"content"` // Markdown消息内容
}

// NewTextReply 创建文本被动回复。
func NewTextReply(content string) *PassiveReply {
	return &PassiveReply{
		MsgType: ReplyMsgTypeText,
		Text: &Text{
			Content: content,
		},
	}
}

// NewMarkdownReply 创建Markdown被动回复。
func NewMarkdownReply(content string) *PassiveReply {
	return &PassiveReply{
		MsgType: ReplyMsgTypeMarkdown,
		Markdown: &Markdown{
			Content: content,
		},
	}
}

// NewTemplateCardReply 创建模板卡片被动回复。
func NewTemplateCardReply(card *TemplateCard) *PassiveReply {
	return &PassiveReply{
		MsgType:      ReplyMsgTypeTemplateCard,
		TemplateCard: card,
	}
}

// NewStreamReply 创建流式消息被动回复。
func NewStreamReply(id, content string, finish bool) *PassiveReply {
	return &PassiveReply{
		MsgType: ReplyMsgTypeStream,
		Stream: &StreamReply{
			ID:      id,
			Content: content,
			Finish:  finish,
		},
	}
}

// NewStreamWithTemplateCardReply 创建流式消息+模板卡片被动回复。
func NewStreamWithTemplateCardReply(stream *StreamReply, card *TemplateCard) *PassiveReply {
	return &PassiveReply{
		MsgType:      ReplyMsgTypeStreamWithTemplateCard,
		Stream:       stream,
		TemplateCard: card,
	}
}

// NewUpdateTemplateCardReply 创建更新模板卡片的被动回复。
func NewUpdateTemplateCardReply(userIDs []string, card *TemplateCard) *PassiveReply {
	return &PassiveReply{
		ResponseType: ReplyResponseTypeUpdateTemplateCard,
		UserIDs:      userIDs,
		TemplateCard: card,
	}
}

// NewEmptyReply 返回空回复（不回复任何消息）。
func NewEmptyReply() *PassiveReply {
	return nil
}
