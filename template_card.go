package wecombot

// TemplateCard 表示模板卡片消息
// 企业微信智能机器人支持5种模板卡片类型：
//   - text_notice: 文本通知模板卡片
//   - news_notice: 图文展示模板卡片
//   - button_interaction: 按钮交互模板卡片
//   - vote_interaction: 投票选择模板卡片
//   - multiple_interaction: 多项选择模板卡片
//
// 参数说明：
//   - CardType: 模板卡片的类型（必填）
//     text_notice: 文本通知模板卡片
//     news_notice: 图文展示模板卡片
//     button_interaction: 按钮交互模板卡片（需要设置回调URL）
//     vote_interaction: 投票选择模板卡片（需要设置回调URL）
//     multiple_interaction: 多项选择模板卡片（需要设置回调URL）
//   - Source: 卡片来源样式信息，不需要来源样式可不填写
//   - ActionMenu: 卡片右上角更多操作按钮
//   - MainTitle: 模板卡片的主要内容，包括一级标题和标题辅助信息
//     注意：main_title.title和sub_title_text必须有一项填写
//   - EmphasisContent: 关键数据样式，建议不与引用样式共用
//   - QuoteArea: 引用文献样式，建议不与关键数据样式共用
//   - SubTitleText: 二级普通文本，建议不超过112个字
//     注意：main_title.title和sub_title_text必须有一项填写
//   - HorizontalContentList: 二级标题+文本列表，列表长度不超过6
//   - JumpList: 跳转指引样式的列表，列表长度不超过3
//     注意：点击文本通知卡片以及图文通知卡片的"跳转指引"区域支持消息智能回复
//   - CardAction: 整体卡片的点击跳转事件
//     text_notice和news_notice类型的卡片该字段为必填项
//   - CardImage: 图片样式
//     news_notice类型的卡片，card_image和image_text_area两者必填一个字段
//   - ImageTextArea: 左图右文样式
//     news_notice类型的卡片，card_image和image_text_area两者必填一个字段
//   - VerticalContentList: 卡片二级垂直内容，列表长度不超过4
//   - ButtonSelection: 下拉式的选择器（用于按钮交互模板卡片）
//   - ButtonList: 按钮列表，列表长度不超过6
//     button_interaction类型的卡片该字段为必填项
//   - Checkbox: 选择题样式（用于投票选择模板卡片）
//     vote_interaction类型的卡片该字段为必填项
//   - SelectList: 下拉式的选择器列表
//     multiple_interaction类型的卡片该字段不可为空，一个消息最多支持3个选择器
//   - SubmitButton: 提交按钮样式
//     vote_interaction和multiple_interaction类型的卡片该字段为必填项
//   - TaskID: 任务id，同一个机器人任务id不能重复
//     只能由数字、字母和"_-@"组成，最长128字节
//     任务id只在发消息时候有效，更新消息的时候无效
//     任务id将会在相应的回调事件中返回
//     当模板卡片有action_menu字段时，该字段必填
//     button_interaction、vote_interaction、multiple_interaction类型的卡片该字段为必填项
//   - Feedback: 反馈信息配置
//     若字段不为空值，回复的消息被用户反馈时会触发回调事件
//     参考：https://developer.work.weixin.qq.com/document/path/101027#59058/用户反馈事件
//
// 官方文档：https://developer.work.weixin.qq.com/document/path/101031#59098
type TemplateCard struct {
	CardType              string              `json:"card_type"` // 模板卡片类型
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

// Source 表示卡片来源样式信息
// 参数说明：
//   - IconURL:   来源图片的URL
//   - Desc:      来源图片的描述，建议不超过13个字
//   - DescColor: 来源文字的颜色
//     0（默认）: 灰色
//     1: 黑色
//     2: 红色
//     3: 绿色
type Source struct {
	IconURL   string `json:"icon_url,omitempty"`   // 来源图片URL
	Desc      string `json:"desc,omitempty"`       // 来源描述，建议不超过13个字
	DescColor int    `json:"desc_color,omitempty"` // 来源文字颜色：0-灰色 1-黑色 2-红色 3-绿色
}

// ActionMenu 表示卡片右上角更多操作按钮
// 参数说明：
//   - Desc:       更多操作界面的描述（必填）
//   - ActionList: 操作列表，列表长度取值范围为[1, 3]（必填）
type ActionMenu struct {
	Desc       string       `json:"desc"`        // 更多操作界面的描述
	ActionList []ActionItem `json:"action_list"` // 操作列表，长度1-3
}

// ActionItem 表示操作菜单中的操作项
// 参数说明：
//   - Text: 操作的描述文案（必填）
//   - Key:  操作key值，用户点击后会产生回调事件将本参数作为EventKey返回
//     回调事件会带上该key值，最长支持1024字节，不可重复（必填）
type ActionItem struct {
	Text string `json:"text"` // 操作描述文案
	Key  string `json:"key"`  // 操作key值，最长1024字节，不可重复
}

// MainTitle 表示模板卡片的主要内容，包括一级标题和标题辅助信息
// 参数说明：
//   - Title: 一级标题，建议不超过26个字
//     注意：模板卡片主要内容的一级标题main_title.title和二级普通文本sub_title_text必须有一项填写
//   - Desc:  标题辅助信息，建议不超过30个字
type MainTitle struct {
	Title string `json:"title,omitempty"` // 一级标题，建议不超过26个字
	Desc  string `json:"desc,omitempty"`  // 标题辅助信息，建议不超过30个字
}

// EmphasisContent 表示关键数据样式
// 建议不与引用样式共用
// 参数说明：
//   - Title: 关键数据样式的数据内容，建议不超过10个字
//   - Desc:  关键数据样式的数据描述内容，建议不超过15个字
type EmphasisContent struct {
	Title string `json:"title,omitempty"` // 关键数据内容，建议不超过10个字
	Desc  string `json:"desc,omitempty"`  // 数据描述，建议不超过15个字
}

// QuoteArea 表示引用文献样式
// 建议不与关键数据样式共用
// 参数说明：
//   - Type:      引用文献样式区域点击事件
//     0或不填: 没有点击事件
//     1: 跳转URL
//     2: 跳转小程序
//   - URL:       点击跳转的URL，type是1时必填
//   - AppID:     点击跳转的小程序的appid，必须是与当前应用关联的小程序，type是2时必填
//   - PagePath:  点击跳转的小程序的pagepath，type是2时选填
//   - Title:     引用文献样式的标题
//   - QuoteText: 引用文献样式的引用文案
type QuoteArea struct {
	Type      int    `json:"type,omitempty"`       // 点击事件类型：0-无 1-跳转URL 2-跳转小程序
	URL       string `json:"url,omitempty"`        // 跳转URL（type=1时必填）
	AppID     string `json:"appid,omitempty"`      // 小程序AppID（type=2时必填）
	PagePath  string `json:"pagepath,omitempty"`   // 小程序页面路径（type=2时选填）
	Title     string `json:"title,omitempty"`      // 引用标题
	QuoteText string `json:"quote_text,omitempty"` // 引用文案
}

// HorizontalContent 表示二级标题+文本列表
// 列表长度不超过6
// 参数说明：
//   - Type:    链接类型
//     0或不填: 普通文本
//     1: 跳转URL
//     3: 点击跳转成员详情
//   - KeyName: 二级标题，建议不超过5个字（必填）
//   - Value:   二级文本，建议不超过26个字
//   - URL:     链接跳转的URL，type是1时必填
//   - UserID:  成员详情的userid，type是3时必填
type HorizontalContent struct {
	Type    int    `json:"type,omitempty"`   // 链接类型：0-文本 1-跳转URL 3-成员详情
	KeyName string `json:"keyname"`          // 二级标题，建议不超过5个字
	Value   string `json:"value,omitempty"`  // 二级文本，建议不超过26个字
	URL     string `json:"url,omitempty"`    // 跳转URL（type=1时必填）
	UserID  string `json:"userid,omitempty"` // 成员UserID（type=3时必填）
}

// JumpAction 表示跳转指引样式的列表
// 列表长度不超过3
// 注意：点击文本通知卡片以及图文通知卡片的"跳转指引"区域支持消息智能回复
// 参数说明：
//   - Type:     跳转链接类型
//     0或不填: 不是链接
//     1: 跳转URL
//     2: 跳转小程序
//     3: 触发消息智能回复
//   - Question: 智能问答问题，最长不超过200个字节，若type为3则必填
//   - Title:    跳转链接样式的文案内容，建议不超过13个字（必填）
//   - URL:      跳转链接的URL，type是1时必填
//   - AppID:    跳转链接的小程序的appid，type是2时必填
//   - PagePath: 跳转链接的小程序的pagepath，type是2时选填
type JumpAction struct {
	Type     int    `json:"type,omitempty"`     // 跳转类型：0-无 1-URL 2-小程序 3-智能回复
	Question string `json:"question,omitempty"` // 智能问答问题（type=3时必填，最长200字节）
	Title    string `json:"title"`              // 文案内容，建议不超过13个字
	URL      string `json:"url,omitempty"`      // 跳转URL（type=1时必填）
	AppID    string `json:"appid,omitempty"`    // 小程序AppID（type=2时必填）
	PagePath string `json:"pagepath,omitempty"` // 小程序页面路径（type=2时选填）
}

// CardAction 表示整体卡片的点击跳转事件
// 参数说明：
//   - Type:     卡片跳转类型（必填）
//     0或不填: 不是链接
//     1: 跳转URL
//     2: 打开小程序
//     注意：text_notice模板卡片中该字段取值范围为[1,2]
//   - URL:      跳转事件的URL，type是1时必填
//   - AppID:    跳转事件的小程序的appid，type是2时必填
//   - PagePath: 跳转事件的小程序的pagepath，type是2时选填
type CardAction struct {
	Type     int    `json:"type"`               // 跳转类型：0-无 1-URL 2-小程序（text_notice必须为1或2）
	URL      string `json:"url,omitempty"`      // 跳转URL（type=1时必填）
	AppID    string `json:"appid,omitempty"`    // 小程序AppID（type=2时必填）
	PagePath string `json:"pagepath,omitempty"` // 小程序页面路径（type=2时选填）
}

// VerticalContent 表示卡片二级垂直内容
// 列表长度不超过4
// 参数说明：
//   - Title: 卡片二级标题，建议不超过26个字（必填）
//   - Desc:  二级普通文本，建议不超过112个字
type VerticalContent struct {
	Title string `json:"title"`          // 二级标题，建议不超过26个字
	Desc  string `json:"desc,omitempty"` // 二级文本，建议不超过112个字
}

// CardImage 表示图片样式
// news_notice类型的卡片，card_image和image_text_area两者必填一个字段，不可都不填
// 参数说明：
//   - URL:         图片的URL（必填）
//   - AspectRatio: 图片的宽高比，宽高比要小于2.25，大于1.3，不填该参数默认1.3
type CardImage struct {
	URL         string  `json:"url"`                    // 图片URL
	AspectRatio float64 `json:"aspect_ratio,omitempty"` // 图片宽高比：1.3-2.25，默认1.3
}

// ImageTextArea 表示左图右文样式
// news_notice类型的卡片，card_image和image_text_area两者必填一个字段，不可都不填
// 参数说明：
//   - Type:     左图右文样式区域点击事件
//     0或不填: 没有点击事件
//     1: 跳转URL
//     2: 跳转小程序
//   - URL:      点击跳转的URL，type是1时必填
//   - AppID:    点击跳转的小程序的appid，必须是与当前应用关联的小程序，type是2时必填
//   - PagePath: 点击跳转的小程序的pagepath，type是2时选填
//   - Title:    左图右文样式的标题
//   - Desc:     左图右文样式的描述
//   - ImageURL: 左图右文样式的图片URL（必填）
type ImageTextArea struct {
	Type     int    `json:"type,omitempty"`     // 点击事件类型：0-无 1-URL 2-小程序
	URL      string `json:"url,omitempty"`      // 跳转URL（type=1时必填）
	AppID    string `json:"appid,omitempty"`    // 小程序AppID（type=2时必填）
	PagePath string `json:"pagepath,omitempty"` // 小程序页面路径（type=2时选填）
	Title    string `json:"title,omitempty"`    // 标题
	Desc     string `json:"desc,omitempty"`     // 描述
	ImageURL string `json:"image_url"`          // 图片URL
}

// SelectionItem 表示下拉式的选择器列表
// multiple_interaction类型的卡片该字段不可为空，一个消息最多支持3个选择器
// 参数说明：
//   - QuestionKey: 下拉式的选择器题目的key，用户提交选项后会产生回调事件
//     回调事件会带上该key值表示该题，最长支持1024字节，不可重复（必填）
//   - Title:       选择器的标题，建议不超过13个字
//   - Disable:     下拉式的选择器是否不可选，false为可选，true为不可选
//     仅在更新模板卡片的时候该字段有效
//   - SelectedID:  默认选定的id，不填或错填默认第一个
//   - OptionList:  选项列表，下拉选项不超过10个，最少1个（必填）
type SelectionItem struct {
	QuestionKey string         `json:"question_key"`          // 选择器题目key，最长1024字节，不可重复
	Title       string         `json:"title,omitempty"`       // 选择器标题，建议不超过13个字
	Disable     bool           `json:"disable,omitempty"`     // 是否不可选（仅更新时有效）
	SelectedID  string         `json:"selected_id,omitempty"` // 默认选定的id
	OptionList  []SelectOption `json:"option_list"`           // 选项列表，1-10个
}

// SelectOption 表示下拉式的选择器选项
// 参数说明：
//   - ID:   下拉式的选择器选项的id，用户提交选项后会产生回调事件
//     回调事件会带上该id值表示该选项，最长支持128字节，不可重复（必填）
//   - Text: 下拉式的选择器选项的文案，建议不超过10个字（必填）
type SelectOption struct {
	ID   string `json:"id"`   // 选项id，最长128字节，不可重复
	Text string `json:"text"` // 选项文案，建议不超过10个字
}

// Button 表示按钮列表
// 列表长度不超过6
// 参数说明：
//   - Text:  按钮文案，建议不超过10个字（必填）
//   - Style: 按钮样式，目前可填1~4，不填或错填默认1
//     按钮样式参考：https://wework.qpic.cn/wwpic/805842_iKxTyYPiRBamTcX_1628665323/0
//   - Key:   按钮key值，用户点击后会产生回调事件将本参数作为event_key返回
//     最长支持1024字节，不可重复（必填）
type Button struct {
	Text  string `json:"text"`            // 按钮文案，建议不超过10个字
	Style int    `json:"style,omitempty"` // 按钮样式：1-4，默认1
	Key   string `json:"key"`             // 按钮key值，最长1024字节，不可重复
}

// Checkbox 表示选择题样式（用于投票选择模板卡片）
// 参数说明：
//   - QuestionKey: 选择题key值，用户提交选项后会产生回调事件
//     回调事件会带上该key值表示该题，最长支持1024字节（必填）
//   - Disable:     投票选择框是否不可选，false为可选，true为不可选
//     仅在更新模板卡片的时候该字段有效
//   - Mode:        选择题模式，单选：0，多选：1，不填默认0
//   - OptionList:  选项list，选项个数不超过20个，最少1个（必填）
type Checkbox struct {
	QuestionKey string           `json:"question_key"`      // 选择题key值，最长1024字节
	Disable     bool             `json:"disable,omitempty"` // 是否不可选（仅更新时有效）
	Mode        int              `json:"mode,omitempty"`    // 选择模式：0-单选 1-多选，默认0
	OptionList  []CheckboxOption `json:"option_list"`       // 选项列表，1-20个
}

// CheckboxOption 表示选择题的选项
// 参数说明：
//   - ID:        选项id，用户提交选项后会产生回调事件
//     回调事件会带上该id值表示该选项，最长支持128字节，不可重复（必填）
//   - Text:      选项文案描述，建议不超过11个字（必填）
//   - IsChecked: 该选项是否要默认选中
type CheckboxOption struct {
	ID        string `json:"id"`                   // 选项id，最长128字节，不可重复
	Text      string `json:"text"`                 // 选项文案，建议不超过11个字
	IsChecked bool   `json:"is_checked,omitempty"` // 是否默认选中
}

// SubmitButton 表示提交按钮样式
// 参数说明：
//   - Text: 按钮文案，建议不超过10个字（必填）
//   - Key:  提交按钮的key，会产生回调事件将本参数作为EventKey返回
//     最长支持1024字节（必填）
type SubmitButton struct {
	Text string `json:"text"` // 按钮文案，建议不超过10个字
	Key  string `json:"key"`  // 提交按钮key，最长1024字节
}

// Feedback 表示反馈信息配置
// 若字段不为空值，回复的消息被用户反馈时会触发回调事件
// 参考：https://developer.work.weixin.qq.com/document/path/101027#59058/用户反馈事件
// 参数说明：
//   - ID: 反馈id，有效长度为256字节以内，必须是utf-8编码
type Feedback struct {
	ID string `json:"id"` // 反馈id，最长256字节
}
