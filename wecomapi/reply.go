package wecomapi

// PassiveReply represents a passive reply to callback
type PassiveReply struct {
	MsgType      string        `json:"msgtype"`
	Text         *Text         `json:"text,omitempty"`
	TemplateCard *TemplateCard `json:"template_card,omitempty"`
	Stream       *StreamReply  `json:"stream,omitempty"`
	ResponseType string        `json:"response_type,omitempty"`
	UserIDs      []string      `json:"userids,omitempty"`
}

// StreamReply represents a stream message reply
type StreamReply struct {
	ID       string    `json:"id,omitempty"`
	Finish   bool      `json:"finish,omitempty"`
	Content  string    `json:"content,omitempty"`
	MsgItem  []MsgItem `json:"msg_item,omitempty"`
	Feedback *Feedback `json:"feedback,omitempty"`
}

// ImageBase64 表示Base64编码的图片
// 用于流式消息中的图片混排
type ImageBase64 struct {
	Base64 string `json:"base64"` // 图片内容的Base64编码（编码前最大10M，支持JPG/PNG）
	MD5    string `json:"md5"`    // 图片内容（Base64编码前）的MD5值
}

// NewTextReply 创建文本回复
// 注意：目前仅支持进入会话回调事件时被动回复文本消息
func NewTextReply(content string) *PassiveReply {
	return &PassiveReply{
		MsgType: "text",
		Text: &Text{
			Content: content,
		},
	}
}

// NewTemplateCardReply 创建模板卡片回复
// 支持进入会话回调事件或接收消息回调时被动回复
func NewTemplateCardReply(card *TemplateCard) *PassiveReply {
	return &PassiveReply{
		MsgType:      "template_card",
		TemplateCard: card,
	}
}

// NewStreamReply 创建流式消息回复
// id: 自定义的唯一ID，首次回复时必须设置
// content: 消息内容，支持Markdown格式和<think>标签
// finish: 是否结束流式消息
func NewStreamReply(id, content string, finish bool) *PassiveReply {
	return &PassiveReply{
		MsgType: "stream",
		Stream: &StreamReply{
			ID:      id,
			Content: content,
			Finish:  finish,
		},
	}
}

// NewStreamWithTemplateCardReply 创建流式消息+模板卡片回复
// 首次回复时必须返回stream的id
// template_card可首次回复，也可在收到流式消息刷新事件时回复
// 注意：同一个消息只能回复一次模板卡片
func NewStreamWithTemplateCardReply(stream *StreamReply, card *TemplateCard) *PassiveReply {
	return &PassiveReply{
		MsgType:      "stream_with_template_card",
		Stream:       stream,
		TemplateCard: card,
	}
}

// NewUpdateTemplateCardReply 创建更新模板卡片回复
// 当接收到模板卡片事件后，可以更新卡片内容
// userIDs: 要更新卡片的用户ID列表，为空则更新所有用户
// card: 新的模板卡片内容，task_id需与回调收到的task_id一致
func NewUpdateTemplateCardReply(userIDs []string, card *TemplateCard) *PassiveReply {
	return &PassiveReply{
		ResponseType: "update_template_card",
		UserIDs:      userIDs,
		TemplateCard: card,
	}
}

// NewEmptyReply 创建空回复（不回复任何消息）
// 用于仅接收消息但不需要回复的场景
func NewEmptyReply() *PassiveReply {
	return nil
}
