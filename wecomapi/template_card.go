package wecomapi

// TemplateCard 模板卡片消息。
// 字段含义以官方文档为准：
// https://developer.work.weixin.qq.com/document/path/101031
type TemplateCard struct {
	CardType              TemplateCardType    `json:"card_type"` // 模板卡片类型
	Source                *Source             `json:"source,omitempty"`
	ActionMenu            *ActionMenu         `json:"action_menu,omitempty"`
	MainTitle             *MainTitle          `json:"main_title,omitempty"`
	EmphasisContent       *EmphasisContent    `json:"emphasis_content,omitempty"`
	QuoteArea             *QuoteArea          `json:"quote_area,omitempty"`
	SubTitleText          string              `json:"sub_title_text,omitempty"`
	HorizontalContentList []HorizontalContent `json:"horizontal_content_list,omitempty"`
	JumpList              []JumpAction        `json:"jump_list,omitempty"`
	CardAction            *CardAction         `json:"card_action,omitempty"`
	CardImage             *CardImage          `json:"card_image,omitempty"`
	ImageTextArea         *ImageTextArea      `json:"image_text_area,omitempty"`
	VerticalContentList   []VerticalContent   `json:"vertical_content_list,omitempty"`
	ButtonSelection       *SelectionItem      `json:"button_selection,omitempty"`
	ButtonList            []Button            `json:"button_list,omitempty"`
	Checkbox              *Checkbox           `json:"checkbox,omitempty"`
	SelectList            []SelectionItem     `json:"select_list,omitempty"`
	SubmitButton          *SubmitButton       `json:"submit_button,omitempty"`
	TaskID                string              `json:"task_id,omitempty"`
	Feedback              *Feedback           `json:"feedback,omitempty"`
}

// SourceDescColor 来源描述文字颜色。
type SourceDescColor int

const (
	SourceDescColorGray  SourceDescColor = 0
	SourceDescColorBlack SourceDescColor = 1
	SourceDescColorRed   SourceDescColor = 2
	SourceDescColorGreen SourceDescColor = 3
)

// Source 卡片来源样式。
type Source struct {
	IconURL   string          `json:"icon_url,omitempty"`   // 来源图片URL
	Desc      string          `json:"desc,omitempty"`       // 来源描述，建议不超过13个字
	DescColor SourceDescColor `json:"desc_color,omitempty"` // 来源文字颜色：0-灰色 1-黑色 2-红色 3-绿色
}

// ActionMenu 卡片右上角更多操作按钮。
type ActionMenu struct {
	Desc       string       `json:"desc"`        // 更多操作界面的描述
	ActionList []ActionItem `json:"action_list"` // 操作列表，长度1-3
}

// ActionItem 操作菜单项。
type ActionItem struct {
	Text string `json:"text"` // 操作描述文案
	Key  string `json:"key"`  // 操作key值，最长1024字节，不可重复
}

// MainTitle 主要内容（标题/辅助信息）。
type MainTitle struct {
	Title string `json:"title,omitempty"` // 一级标题，建议不超过26个字
	Desc  string `json:"desc,omitempty"`  // 标题辅助信息，建议不超过30个字
}

// EmphasisContent 关键数据样式。
type EmphasisContent struct {
	Title string `json:"title,omitempty"` // 关键数据内容，建议不超过10个字
	Desc  string `json:"desc,omitempty"`  // 数据描述，建议不超过15个字
}

// QuoteAreaType 引用区域点击类型。
type QuoteAreaType int

const (
	QuoteAreaTypeNone   QuoteAreaType = 0
	QuoteAreaTypeURL    QuoteAreaType = 1
	QuoteAreaTypeApplet QuoteAreaType = 2
)

// QuoteArea 引用样式区域。
type QuoteArea struct {
	Type      QuoteAreaType `json:"type,omitempty"`       // 点击事件类型：0-无 1-跳转URL 2-跳转小程序
	URL       string        `json:"url,omitempty"`        // 跳转URL（type=1时必填）
	AppID     string        `json:"appid,omitempty"`      // 小程序AppID（type=2时必填）
	PagePath  string        `json:"pagepath,omitempty"`   // 小程序页面路径（type=2时选填）
	Title     string        `json:"title,omitempty"`      // 引用标题
	QuoteText string        `json:"quote_text,omitempty"` // 引用文案
}

// HorizontalContentType 二级内容链接类型。
type HorizontalContentType int

const (
	HorizontalContentTypeText   HorizontalContentType = 0
	HorizontalContentTypeURL    HorizontalContentType = 1
	HorizontalContentTypeUserID HorizontalContentType = 3
)

// HorizontalContent 二级标题+文本项。
type HorizontalContent struct {
	Type    HorizontalContentType `json:"type,omitempty"`   // 链接类型：0-文本 1-跳转URL 3-成员详情
	KeyName string                `json:"keyname"`          // 二级标题，建议不超过5个字
	Value   string                `json:"value,omitempty"`  // 二级文本，建议不超过26个字
	URL     string                `json:"url,omitempty"`    // 跳转URL（type=1时必填）
	UserID  string                `json:"userid,omitempty"` // 成员UserID（type=3时必填）
}

// JumpActionType 跳转指引类型。
type JumpActionType int

const (
	JumpActionTypeNone     JumpActionType = 0
	JumpActionTypeURL      JumpActionType = 1
	JumpActionTypeApplet   JumpActionType = 2
	JumpActionTypeQuestion JumpActionType = 3
)

// JumpAction 跳转指引项。
type JumpAction struct {
	Type     JumpActionType `json:"type,omitempty"`     // 跳转类型：0-无 1-URL 2-小程序 3-智能回复
	Question string         `json:"question,omitempty"` // 智能问答问题（type=3时必填，最长200字节）
	Title    string         `json:"title"`              // 文案内容，建议不超过13个字
	URL      string         `json:"url,omitempty"`      // 跳转URL（type=1时必填）
	AppID    string         `json:"appid,omitempty"`    // 小程序AppID（type=2时必填）
	PagePath string         `json:"pagepath,omitempty"` // 小程序页面路径（type=2时选填）
}

// CardActionType 整体跳转类型。
type CardActionType int

const (
	CardActionTypeNone   CardActionType = 0
	CardActionTypeURL    CardActionType = 1
	CardActionTypeApplet CardActionType = 2
)

// CardAction 整体卡片点击跳转。
type CardAction struct {
	Type     CardActionType `json:"type"`               // 跳转类型：0-无 1-URL 2-小程序（text_notice必须为1或2）
	URL      string         `json:"url,omitempty"`      // 跳转URL（type=1时必填）
	AppID    string         `json:"appid,omitempty"`    // 小程序AppID（type=2时必填）
	PagePath string         `json:"pagepath,omitempty"` // 小程序页面路径（type=2时选填）
}

// VerticalContent 二级垂直内容。
type VerticalContent struct {
	Title string `json:"title"`          // 二级标题，建议不超过26个字
	Desc  string `json:"desc,omitempty"` // 二级文本，建议不超过112个字
}

// CardImage 图片样式。
type CardImage struct {
	URL         string  `json:"url"`                    // 图片URL
	AspectRatio float64 `json:"aspect_ratio,omitempty"` // 图片宽高比：1.3-2.25，默认1.3
}

// ImageTextAreaType 左图右文区域点击类型。
type ImageTextAreaType int

const (
	ImageTextAreaTypeNone   ImageTextAreaType = 0
	ImageTextAreaTypeURL    ImageTextAreaType = 1
	ImageTextAreaTypeApplet ImageTextAreaType = 2
)

// ImageTextArea 左图右文样式。
type ImageTextArea struct {
	Type     ImageTextAreaType `json:"type,omitempty"`     // 点击事件类型：0-无 1-URL 2-小程序
	URL      string            `json:"url,omitempty"`      // 跳转URL（type=1时必填）
	AppID    string            `json:"appid,omitempty"`    // 小程序AppID（type=2时必填）
	PagePath string            `json:"pagepath,omitempty"` // 小程序页面路径（type=2时选填）
	Title    string            `json:"title,omitempty"`    // 标题
	Desc     string            `json:"desc,omitempty"`     // 描述
	ImageURL string            `json:"image_url"`          // 图片URL
}

// SelectionItem 下拉选择器。
type SelectionItem struct {
	QuestionKey string         `json:"question_key"`          // 选择器题目key，最长1024字节，不可重复
	Title       string         `json:"title,omitempty"`       // 选择器标题，建议不超过13个字
	Disable     bool           `json:"disable,omitempty"`     // 是否不可选（仅更新时有效）
	SelectedID  string         `json:"selected_id,omitempty"` // 默认选定的id
	OptionList  []SelectOption `json:"option_list"`           // 选项列表，1-10个
}

// SelectOption 下拉选项。
type SelectOption struct {
	ID   string `json:"id"`   // 选项id，最长128字节，不可重复
	Text string `json:"text"` // 选项文案，建议不超过10个字
}

// ButtonStyle 按钮样式。
type ButtonStyle int

const (
	ButtonStylePrimary ButtonStyle = 1
	ButtonStyleBlue    ButtonStyle = 2
	ButtonStyleRed     ButtonStyle = 3
	ButtonStyleGray    ButtonStyle = 4
)

// Button 按钮项。
type Button struct {
	Text  string      `json:"text"`            // 按钮文案，建议不超过10个字
	Style ButtonStyle `json:"style,omitempty"` // 按钮样式：1-4，默认1
	Key   string      `json:"key"`             // 按钮key值，最长1024字节，不可重复
}

// CheckboxMode 选择题模式。
type CheckboxMode int

const (
	CheckboxModeSingle CheckboxMode = 0
	CheckboxModeMulti  CheckboxMode = 1
)

// Checkbox 选择题样式。
type Checkbox struct {
	QuestionKey string           `json:"question_key"`      // 选择题key值，最长1024字节
	Disable     bool             `json:"disable,omitempty"` // 是否不可选（仅更新时有效）
	Mode        CheckboxMode     `json:"mode,omitempty"`    // 选择模式：0-单选 1-多选，默认0
	OptionList  []CheckboxOption `json:"option_list"`       // 选项列表，1-20个
}

// CheckboxOption 选择题选项。
type CheckboxOption struct {
	ID        string `json:"id"`                   // 选项id，最长128字节，不可重复
	Text      string `json:"text"`                 // 选项文案，建议不超过11个字
	IsChecked bool   `json:"is_checked,omitempty"` // 是否默认选中
}

// SubmitButton 提交按钮。
type SubmitButton struct {
	Text string `json:"text"` // 按钮文案，建议不超过10个字
	Key  string `json:"key"`  // 提交按钮key，最长1024字节
}

// Feedback 反馈信息配置。
type Feedback struct {
	ID string `json:"id"` // 反馈id，最长256字节
}
