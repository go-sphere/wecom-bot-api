package wecomapi

import "encoding/json"

// ChatType 回调会话类型。
type ChatType string

const (
	ChatTypeSingle ChatType = "single"
	ChatTypeGroup  ChatType = "group"
)

// CallbackMsgType 回调消息类型。
type CallbackMsgType string

const (
	CallbackMsgTypeText   CallbackMsgType = "text"
	CallbackMsgTypeImage  CallbackMsgType = "image"
	CallbackMsgTypeMixed  CallbackMsgType = "mixed"
	CallbackMsgTypeVoice  CallbackMsgType = "voice"
	CallbackMsgTypeFile   CallbackMsgType = "file"
	CallbackMsgTypeStream CallbackMsgType = "stream"
	CallbackMsgTypeEvent  CallbackMsgType = "event"
)

// MsgItemType 混排消息项类型。
type MsgItemType string

const (
	MsgItemTypeText  MsgItemType = "text"
	MsgItemTypeImage MsgItemType = "image"
)

// QuoteMsgType 引用消息类型。
type QuoteMsgType string

const (
	QuoteMsgTypeText  QuoteMsgType = "text"
	QuoteMsgTypeImage QuoteMsgType = "image"
	QuoteMsgTypeMixed QuoteMsgType = "mixed"
	QuoteMsgTypeVoice QuoteMsgType = "voice"
	QuoteMsgTypeFile  QuoteMsgType = "file"
)

// EventType 回调事件类型。
type EventType string

const (
	EventTypeEnterChat    EventType = "enter_chat"
	EventTypeTemplateCard EventType = "template_card_event"
	EventTypeFeedback     EventType = "feedback_event"
)

// TemplateCardType 模板卡片类型。
type TemplateCardType string

const (
	TemplateCardTypeTextNotice          TemplateCardType = "text_notice"
	TemplateCardTypeNewsNotice          TemplateCardType = "news_notice"
	TemplateCardTypeButtonInteraction   TemplateCardType = "button_interaction"
	TemplateCardTypeVoteInteraction     TemplateCardType = "vote_interaction"
	TemplateCardTypeMultipleInteraction TemplateCardType = "multiple_interaction"
)

// FeedbackType 反馈类型。
type FeedbackType int

const (
	FeedbackTypeAccurate   FeedbackType = 1
	FeedbackTypeInaccurate FeedbackType = 2
	FeedbackTypeCancel     FeedbackType = 3
)

// FeedbackInaccurateReason 负反馈原因。
type FeedbackInaccurateReason int

const (
	FeedbackInaccurateReasonIrrelevant   FeedbackInaccurateReason = 1
	FeedbackInaccurateReasonIncomplete   FeedbackInaccurateReason = 2
	FeedbackInaccurateReasonIncorrect    FeedbackInaccurateReason = 3
	FeedbackInaccurateReasonDataAnalysis FeedbackInaccurateReason = 4
)

// Callback 回调消息的通用结构。
type Callback struct {
	MsgID       string          `json:"msgid"`
	CreateTime  int64           `json:"create_time,omitempty"`
	AIBotID     string          `json:"aibotid"`
	ChatID      string          `json:"chatid,omitempty"`
	ChatType    ChatType        `json:"chattype"` // single 或 group
	From        From            `json:"from"`
	ResponseURL string          `json:"response_url,omitempty"`
	MsgType     CallbackMsgType `json:"msgtype"` // text/image/mixed/voice/file/stream/event
	Text        *Text           `json:"text,omitempty"`
	Image       *Image          `json:"image,omitempty"`
	Mixed       *Mixed          `json:"mixed,omitempty"`
	Voice       *Voice          `json:"voice,omitempty"`
	File        *File           `json:"file,omitempty"`
	Stream      *Stream         `json:"stream,omitempty"`
	Quote       *Quote          `json:"quote,omitempty"`
	Event       *Event          `json:"event,omitempty"`
}

// From 消息发送者信息。
type From struct {
	CorpID string `json:"corpid,omitempty"` // 企业ID（内部机器人不返回）
	UserID string `json:"userid"`           // 操作者UserID
}

// Text 文本消息内容。
type Text struct {
	Content string `json:"content"` // 文本消息内容
}

// Image 图片消息内容。
type Image struct {
	URL string `json:"url"` // 图片下载URL（已加密）
}

// Mixed 图文混排消息内容。
type Mixed struct {
	MsgItem []MsgItem `json:"msg_item"` // 图文混排消息列表
}

// MsgItem 混排中的单个元素。
type MsgItem struct {
	MsgType MsgItemType `json:"msgtype"`         // 类型：text 或 image
	Text    *Text       `json:"text,omitempty"`  // 文本内容
	Image   *Image      `json:"image,omitempty"` // 图片内容
}

// Voice 语音消息内容。
type Voice struct {
	Content string `json:"content"` // 语音转换成的文本内容
}

// File 文件消息内容。
type File struct {
	URL string `json:"url"` // 文件下载URL（已加密）
}

// Stream 流式消息刷新事件。
type Stream struct {
	ID string `json:"id"` // 流式消息的ID
}

// Quote 引用的消息内容。
type Quote struct {
	MsgType QuoteMsgType `json:"msgtype"`         // 引用的消息类型
	Text    *Text        `json:"text,omitempty"`  // 引用的文本内容
	Image   *Image       `json:"image,omitempty"` // 引用的图片内容
	Mixed   *Mixed       `json:"mixed,omitempty"` // 引用的图文混排内容
	Voice   *Voice       `json:"voice,omitempty"` // 引用的语音内容
	File    *File        `json:"file,omitempty"`  // 引用的文件内容
}

// Event 事件回调。
type Event struct {
	EventType         EventType          `json:"eventtype"`                     // 事件类型
	EnterChat         *EnterChatEvent    `json:"enter_chat,omitempty"`          // 进入会话事件
	TemplateCardEvent *TemplateCardEvent `json:"template_card_event,omitempty"` // 模板卡片事件
	FeedbackEvent     *FeedbackEvent     `json:"feedback_event,omitempty"`      // 用户反馈事件
	RawData           map[string]any     `json:"-"`                             // 原始数据，用于自定义事件解析
}

// UnmarshalJSON 自定义反序列化，保存原始数据。
func (e *Event) UnmarshalJSON(data []byte) error {
	type Alias Event
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// 保存原始数据，便于自定义解析
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.RawData = raw

	return nil
}

// EnterChatEvent 进入会话事件。
type EnterChatEvent struct{}

// TemplateCardEvent 模板卡片交互事件。
type TemplateCardEvent struct {
	CardType      TemplateCardType `json:"card_type"`                // 模板卡片类型
	EventKey      string           `json:"event_key"`                // 用户点击的按钮key
	TaskID        string           `json:"task_id"`                  // 交互模板卡片的任务ID
	SelectedItems *SelectedItems   `json:"selected_items,omitempty"` // 用户提交的选择框数据
}

// SelectedItems 选中项列表。
type SelectedItems struct {
	SelectedItem []SelectedItem `json:"selected_item"` // 选中项列表
}

// SelectedItem 单个选中项。
type SelectedItem struct {
	QuestionKey string    `json:"question_key"` // 选择框的key值
	OptionIDs   OptionIDs `json:"option_ids"`   // 选中的选项ID列表
}

// OptionIDs 选中的选项ID列表。
type OptionIDs struct {
	OptionID []string `json:"option_id"` // 选项ID数组
}

// FeedbackEvent 用户反馈事件。
type FeedbackEvent struct {
	ID                   string                     `json:"id"`                               // 反馈ID
	Type                 FeedbackType               `json:"type"`                             // 反馈类型：1-准确 2-不准确 3-取消
	Content              string                     `json:"content,omitempty"`                // 用户输入的反馈内容（仅不准确时返回）
	InaccurateReasonList []FeedbackInaccurateReason `json:"inaccurate_reason_list,omitempty"` // 负反馈原因列表：1-与问题无关 2-内容不完整 3-内容有错误 4-数据分析错误
}
