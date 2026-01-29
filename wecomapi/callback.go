package wecomapi

import "encoding/json"

// ChatType represents the chat scope for callbacks.
type ChatType string

const (
	ChatTypeSingle ChatType = "single"
	ChatTypeGroup  ChatType = "group"
)

// CallbackMsgType represents msgtype values for callback messages.
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

// MsgItemType represents msgtype values within mixed message items.
type MsgItemType string

const (
	MsgItemTypeText  MsgItemType = "text"
	MsgItemTypeImage MsgItemType = "image"
)

// QuoteMsgType represents msgtype values for quoted messages.
type QuoteMsgType string

const (
	QuoteMsgTypeText  QuoteMsgType = "text"
	QuoteMsgTypeImage QuoteMsgType = "image"
	QuoteMsgTypeMixed QuoteMsgType = "mixed"
	QuoteMsgTypeVoice QuoteMsgType = "voice"
	QuoteMsgTypeFile  QuoteMsgType = "file"
)

// EventType represents eventtype values for callback events.
type EventType string

const (
	EventTypeEnterChat    EventType = "enter_chat"
	EventTypeTemplateCard EventType = "template_card_event"
	EventTypeFeedback     EventType = "feedback_event"
)

// TemplateCardType represents card_type values for template card events.
type TemplateCardType string

const (
	TemplateCardTypeTextNotice          TemplateCardType = "text_notice"
	TemplateCardTypeNewsNotice          TemplateCardType = "news_notice"
	TemplateCardTypeButtonInteraction   TemplateCardType = "button_interaction"
	TemplateCardTypeVoteInteraction     TemplateCardType = "vote_interaction"
	TemplateCardTypeMultipleInteraction TemplateCardType = "multiple_interaction"
)

// FeedbackType represents feedback type values.
type FeedbackType int

const (
	FeedbackTypeAccurate   FeedbackType = 1
	FeedbackTypeInaccurate FeedbackType = 2
	FeedbackTypeCancel     FeedbackType = 3
)

// FeedbackInaccurateReason represents reasons for inaccurate feedback.
type FeedbackInaccurateReason int

const (
	FeedbackInaccurateReasonIrrelevant   FeedbackInaccurateReason = 1
	FeedbackInaccurateReasonIncomplete   FeedbackInaccurateReason = 2
	FeedbackInaccurateReasonIncorrect    FeedbackInaccurateReason = 3
	FeedbackInaccurateReasonDataAnalysis FeedbackInaccurateReason = 4
)

// Callback 表示所有回调消息的通用结构
// 当用户与智能机器人交互时，企业微信会将加密的回调消息推送到开发者设置的回调URL
type Callback struct {
	MsgID       string          `json:"msgid"`
	CreateTime  int64           `json:"create_time,omitempty"`
	AIBotID     string          `json:"aibotid"`
	ChatID      string          `json:"chatid,omitempty"`
	ChatType    ChatType        `json:"chattype"` // "single" or "group"
	From        From            `json:"from"`
	ResponseURL string          `json:"response_url,omitempty"`
	MsgType     CallbackMsgType `json:"msgtype"` // "text", "image", "mixed", "voice", "file", "stream", "event"
	Text        *Text           `json:"text,omitempty"`
	Image       *Image          `json:"image,omitempty"`
	Mixed       *Mixed          `json:"mixed,omitempty"`
	Voice       *Voice          `json:"voice,omitempty"`
	File        *File           `json:"file,omitempty"`
	Stream      *Stream         `json:"stream,omitempty"`
	Quote       *Quote          `json:"quote,omitempty"`
	Event       *Event          `json:"event,omitempty"`
}

// From 表示消息发送者信息
type From struct {
	CorpID string `json:"corpid,omitempty"` // 企业ID，企业内部机器人不返回
	UserID string `json:"userid"`           // 操作者的UserID
}

// Text 表示文本消息内容
type Text struct {
	Content string `json:"content"` // 文本消息内容
}

// Image 表示图片消息内容
// 注意：图片URL五分钟内有效，获取到的文件是已加密的
// 加密方式：AES-256-CBC，IV为AESKey前16字节
type Image struct {
	URL string `json:"url"` // 图片下载URL（已加密）
}

// Mixed 表示图文混排消息内容
type Mixed struct {
	MsgItem []MsgItem `json:"msg_item"` // 图文混排消息列表
}

// MsgItem 表示图文混排中的单个元素
type MsgItem struct {
	MsgType MsgItemType `json:"msgtype"`         // 类型：text 或 image
	Text    *Text       `json:"text,omitempty"`  // 文本内容
	Image   *Image      `json:"image,omitempty"` // 图片内容
}

// Voice 表示语音消息内容
type Voice struct {
	Content string `json:"content"` // 语音转换成的文本内容
}

// File 表示文件消息内容
// 注意：智能机器人目前仅支持100M大小以内的文件回调
// 文件URL五分钟内有效，获取到的文件是已加密的
type File struct {
	URL string `json:"url"` // 文件下载URL（已加密）
}

// Stream 表示流式消息刷新
// 企业微信会不断推送流式消息刷新事件（最多等待6分钟）
type Stream struct {
	ID string `json:"id"` // 流式消息的ID
}

// Quote 表示引用的消息内容
// 用户可以引用其他消息进行回复
type Quote struct {
	MsgType QuoteMsgType `json:"msgtype"`         // 引用的消息类型
	Text    *Text        `json:"text,omitempty"`  // 引用的文本内容
	Image   *Image       `json:"image,omitempty"` // 引用的图片内容
	Mixed   *Mixed       `json:"mixed,omitempty"` // 引用的图文混排内容
	Voice   *Voice       `json:"voice,omitempty"` // 引用的语音内容
	File    *File        `json:"file,omitempty"`  // 引用的文件内容
}

// Event 表示事件回调
// 当用户与智能机器人发生交互时，会触发各种事件
type Event struct {
	EventType         EventType              `json:"eventtype"`                     // 事件类型
	EnterChat         *EnterChatEvent        `json:"enter_chat,omitempty"`          // 进入会话事件
	TemplateCardEvent *TemplateCardEvent     `json:"template_card_event,omitempty"` // 模板卡片事件
	FeedbackEvent     *FeedbackEvent         `json:"feedback_event,omitempty"`      // 用户反馈事件
	RawData           map[string]interface{} `json:"-"`                             // 原始数据，用于自定义事件解析
}

// UnmarshalJSON 自定义JSON反序列化，用于处理动态事件类型
// 将原始JSON数据保存到RawData字段，方便自定义事件解析
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

	// Store raw data for custom parsing
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.RawData = raw

	return nil
}

// EnterChatEvent 表示进入会话事件
// 当用户当天首次进入智能机器人单聊会话时触发
// 注意：若未回复消息，用户当天再次进入也不再推送此事件
type EnterChatEvent struct{}

// TemplateCardEvent 表示模板卡片交互事件
// 当用户点击模板卡片中的按钮、选择器等交互元素时触发
// 注意：企业微信服务器在5秒内收不到响应会断开连接，丢弃该回调事件
type TemplateCardEvent struct {
	CardType      TemplateCardType `json:"card_type"`                // 模板卡片类型
	EventKey      string           `json:"event_key"`                // 用户点击的按钮key
	TaskID        string           `json:"task_id"`                  // 交互模板卡片的任务ID
	SelectedItems *SelectedItems   `json:"selected_items,omitempty"` // 用户提交的选择框数据
}

// SelectedItems 表示模板卡片中选中的项目列表
type SelectedItems struct {
	SelectedItem []SelectedItem `json:"selected_item"` // 选中项列表
}

// SelectedItem 表示单个选中项
type SelectedItem struct {
	QuestionKey string    `json:"question_key"` // 选择框的key值
	OptionIDs   OptionIDs `json:"option_ids"`   // 选中的选项ID列表
}

// OptionIDs 表示选中的选项ID列表
// 单选时只有一个ID，多选时可能有多个ID
type OptionIDs struct {
	OptionID []string `json:"option_id"` // 选项ID数组
}

// FeedbackEvent 表示用户反馈事件
// 当用户对机器人回复进行反馈时触发
// 注意：该事件目前仅支持回复空包，不支持回复新消息或更新卡片
type FeedbackEvent struct {
	ID                   string                     `json:"id"`                               // 反馈ID
	Type                 FeedbackType               `json:"type"`                             // 反馈类型：1-准确 2-不准确 3-取消
	Content              string                     `json:"content,omitempty"`                // 用户输入的反馈内容（仅不准确时返回）
	InaccurateReasonList []FeedbackInaccurateReason `json:"inaccurate_reason_list,omitempty"` // 负反馈原因列表：1-与问题无关 2-内容不完整 3-内容有错误 4-数据分析错误
}
