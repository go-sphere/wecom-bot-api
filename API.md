# 企业内部开发/服务端API/消息接收与发送/智能机器人

## API设置

若智能机器人开启API模式，当用户跟智能机器人交互时，企业微信会向智能机器人API设置中的URL的回调地址推送相关消息跟事件。
![img](https://wework.qpic.cn/wwpic3az/740428_c6nXsGv7RAaRVtH_1750834300/0)

> 配置URL时，需要[验证url有效性](https://developer.work.weixin.qq.com/document/path/101039#59137/验证url有效性)。

## 接收回调与被动回复

当用户跟智能机器人交互时，企业微信会向智能机器人的回调URL上推送相关消息或者事件，开发者可根据接收的[消息](https://developer.work.weixin.qq.com/document/path/101039#57141)或者[事件](https://developer.work.weixin.qq.com/document/path/101039#59058)，[被动回复消息](https://developer.work.weixin.qq.com/document/path/101039#59068)。具体流程如下：

![企业微信截图](https://wework.qpic.cn/wwpic3az/919943_FZ0ZWbxJQrqKILh_1751957879/0)

> 接收消息与被动回复消息都需加密，加密方式参考[回调和回复的加解密方案](https://developer.work.weixin.qq.com/document/path/101039#59137)。

## 接收消息

最后更新：2025/12/11

### 概述

当用户与智能机器人发生交互，向机器人发送消息时，交互事件将加密回调给机器人接受消息url，智能机器人服务通过接收并处理回调消息，实现更加丰富的自定义功能。

目前支持触发消息回调交互场景：

1. 用户群里@智能机器人或者单聊中向智能机器人发送文本消息

2. 用户群里@智能机器人或者单聊中向智能机器人发送图文混排消息

3. 用户单聊中向智能机器人发送图片消息

4. 用户单聊中向智能机器人发送语音消息

5. 用户单聊中向智能机器人发送本地文件消息

6. 用户群里@智能机器人或者单聊中向智能机器人发送引用消息

   ![3 (1).png](https://wework.qpic.cn/wwpic3az/697572_q7ZifS-8RiSJyHw_1750159499/0) ![4.png](https://wework.qpic.cn/wwpic3az/108624_G2fA8Qr6TLSgnKp_1750159518/0)

   交互流程如下图所示：
   ![企业微信截图_17509401646139.png](https://wework.qpic.cn/wwpic3az/452410_NWYPvpVkTneezum_1750940240/0)

   流程说明：
   1.当用户跟智能机器人交互发送支持的消息类型时，企业微信后台会向开发者后台推送[消息推送](https://developer.work.weixin.qq.com/document/path/100719#消息推送)。用户跟同一个智能机器人最多同时有三条消息交互中。
   2.开发者回调url接收到新消息推送后：可选择生成[流式消息回复](https://developer.work.weixin.qq.com/document/path/100719#59068/流式消息回复)，并使用用户消息内容调用大模型/AIAgent；也可直接回复[模板卡片消息](https://developer.work.weixin.qq.com/document/path/100719#59098)。
   3.若开发者回复消息类型包含流式消息，企业微信在未收到流式消息回复结束前，会不断向开发者回调url推送[流式消息刷新](https://developer.work.weixin.qq.com/document/path/100719#流式消息刷新)（从用户发消息开始最多等待6min，超过6min结束推送）。开发者接收到流式消息刷新后，生成[流式消息回复](https://developer.work.weixin.qq.com/document/path/100719#59068/流式消息回复)。

   

   > 接收消息与被动回复消息都是加密的，加密方式参考[回调和回复的加解密方案](https://developer.work.weixin.qq.com/document/path/100719#59137)。

### 消息推送

#### 文本消息

**协议格式如下：**

```javascript
{
    "msgid": "CAIQ16HMjQYY\/NGagIOAgAMgq4KM0AI=",
    "aibotid": "AIBOTID",
    "chatid": "CHATID",
    "chattype": "group",
    "from": {
        "userid": "USERID"
    },
    "response_url": "RESPONSEURL",
    "msgtype": "text",
    "text": {
        "content": "@RobotA hello robot"
    },
    "quote": {
        "msgtype": "text",
        "text": {
            "content": "这是今日的测试情况"
        }
    }
}
```

**参数说明：**

| 参数         | 说明                                                         |
| ------------ | ------------------------------------------------------------ |
| msgid        | 本次回调的唯一性标志，开发者需据此进行事件排重（可能因为网络等原因重复回调） |
| aibotid      | 智能机器人id                                                 |
| chatid       | 会话id，仅群聊类型时候返回                                   |
| chattype     | 会话类型，single\group，分别表示：单聊\群聊                  |
| from         | 该事件触发者的信息                                           |
| from.userid  | 操作者的userid                                               |
| response_url | 支持主动回复消息的临时[url](https://developer.work.weixin.qq.com/document/path/100719#59947) |
| msgtype      | 消息类型，此时固定是text                                     |
| text         | 文本消息内容，可参考 [文本](https://developer.work.weixin.qq.com/document/path/100719#57141/文本) 结构体说明 |
| quote        | 引用内容，若用户引用了其他消息则有该字段，可参考 [引用](https://developer.work.weixin.qq.com/document/path/100719#57141/引用) 结构体说明 |

#### 图片消息

```javascript
{
    "msgid": "CAIQz7/MjQYY/NGagIOAgAMgl8jK/gI=",
    "aibotid": "AIBOTID",
    "chattype": "single",
    "from": {
        "userid": "USERID"
    },
    "response_url": "RESPONSEURL",
    "msgtype": "image",
    "image": {
        "url": "https://ww-aibot-img-1258476243.cos.ap-guangzhou.myqcloud.com/BHoPdA3/7571665296904772241?sign=q-sign-algorithm%3Dsha1%26q-ak%3DAKIDbBpaJvdLBvpnibcYcfyPuaO5f9U1UoGo%26q-sign-time%3D1733467811%3B1733468111%26q-key-time%3D1733467811%3B1733468111%26q-header-list%3Dhost%26q-url-param-list%3D%26q-signature%3D0f7b6576943685f82870bc8db306dbf09dfe0fd6 "
    }
}
```

**参数说明：**

| 参数         | 说明                                                         |
| ------------ | ------------------------------------------------------------ |
| msgid        | 本次回调的唯一性标志，开发者需据此进行事件排重（可能因为网络等原因重复回调） |
| aibotid      | 智能机器人id                                                 |
| chattype     | 会话类型，single，表示单聊。该回调消息类型仅支持单聊中发送   |
| from         | 该事件触发者的信息                                           |
| from.userid  | 操作者的userid                                               |
| response_url | 支持主动回复消息的临时[url](https://developer.work.weixin.qq.com/document/path/100719#59947) |
| msgtype      | 消息类型，此时固定是image                                    |
| image        | 图片的内容，可参考 [图片](https://developer.work.weixin.qq.com/document/path/100719#57141/图片) 结构体说明 |

#### 图文混排消息

```javascript
{
    "msgid": "CAIQrcjMjQYY/NGagIOAgAMg6PDc/w0=",
    "aibotid": "AIBOTID",
    "chatid": "CHATID",
    "chattype": "group",
    "from": {
        "userid": "USERID"
    },
    "response_url": "RESPONSEURL",
    "msgtype": "mixed",
    "mixed": {
        "msg_item": [
            {
                "msgtype": "text",
                "text": {
                    "content": "@机器人 这是今日的测试情况"
                }
            },
            {
                "msgtype": "image",
                "image": {
                    "url": "https://ww-aibot-img-1258476243.cos.ap-guangzhou.myqcloud.com/BHoPdA3/7571665296904772241?sign=q-sign-algorithm%3Dsha1%26q-ak%3DAKIDbBpaJvdLBvpnibcYcfyPuaO5f9U1UoGo%26q-sign-time%3D1733467811%3B1733468111%26q-key-time%3D1733467811%3B1733468111%26q-header-list%3Dhost%26q-url-param-list%3D%26q-signature%3D0f7b6576943685f82870bc8db306dbf09dfe0fd6 "
                }
            }
        ]
    },
    "quote": {
        "msgtype": "text",
        "text": {
            "content": "这是今日的测试情况"
        }
    }
}
```

**参数说明：**

| 参数         | 说明                                                         |
| ------------ | ------------------------------------------------------------ |
| msgid        | 本次回调的唯一性标志，开发者需据此进行事件排重（可能因为网络等原因重复回调） |
| aibotid      | 智能机器人id                                                 |
| chatid       | 会话id，仅群聊类型时候返回                                   |
| chattype     | 会话类型，single\group，分别表示：单聊\群聊                  |
| from         | 该事件触发者的信息                                           |
| from.userid  | 操作者的userid                                               |
| response_url | 支持主动回复消息的临时[url](https://developer.work.weixin.qq.com/document/path/100719#59947) |
| msgtype      | 消息类型，此时固定是mixed                                    |
| mixed        | 图文混排内容，可参考 [图文混排](https://developer.work.weixin.qq.com/document/path/100719#57141/图文混排) 结构体说明 |
| quote        | 引用内容，若用户引用了其他消息则有该字段，可参考 [引用](https://developer.work.weixin.qq.com/document/path/100719#57141/引用) 结构体说明 |

#### 语音消息

```javascript
{
    "msgid": "CAIQrcjMjQYY/NGagIOAgAMg6PDc/w0=",
    "aibotid": "AIBOTID",
    "chattype": "single",
    "from": {
        "userid": "USERID"
    },
    "response_url": "RESPONSEURL",
    "msgtype": "voice",
    "voice": {
        "content": "这是语音转成文本的内容"
    }
}
```

**参数说明：**

| 参数         | 说明                                                         |
| ------------ | ------------------------------------------------------------ |
| msgid        | 本次回调的唯一性标志，开发者需据此进行事件排重（可能因为网络等原因重复回调） |
| aibotid      | 智能机器人id                                                 |
| chattype     | 会话类型，single，表示单聊。该回调消息类型仅支持单聊中发送   |
| from         | 该事件触发者的信息                                           |
| from.userid  | 操作者的userid                                               |
| response_url | 支持主动回复消息的临时[url](https://developer.work.weixin.qq.com/document/path/100719#59947) |
| msgtype      | 消息类型，此时固定是voice                                    |
| voice        | 语音内容，可参考 [语音](https://developer.work.weixin.qq.com/document/path/100719#57141/语音) 结构体说明 |

 

#### 文件消息

```javascript
{
    "msgid": "CAIQrcjMjQYY/NGagIOAgAMg6PDc/w0=",
    "aibotid": "AIBOTID",
    "chattype": "single",
    "from": {
        "userid": "USERID"
    },
    "response_url": "RESPONSEURL",
    "msgtype": "file",
    "file": {
        "url": "https://ww-aibot-img-1258476243.cos.ap-guangzhou.myqcloud.com/BHoPdA3/7571665296904772241?sign=q-sign-algorithm%3Dsha1%26q-ak%3DAKIDbBpaJvdLBvpnibcYcfyPuaO5f9U1UoGo%26q-sign-time%3D1733467811%3B1733468111%26q-key-time%3D1733467811%3B1733468111%26q-header-list%3Dhost%26q-url-param-list%3D%26q-signature%3D0f7b6576943685f82870bc8db306dbf09df00000 "
    }
}
```

**参数说明：**

| 参数         | 说明                                                         |
| ------------ | ------------------------------------------------------------ |
| msgid        | 本次回调的唯一性标志，开发者需据此进行事件排重（可能因为网络等原因重复回调） |
| aibotid      | 智能机器人id                                                 |
| chatid       | 会话id，仅群聊类型时候返回                                   |
| chattype     | 会话类型，single，表示单聊。该回调消息类型仅支持单聊中发送   |
| from         | 该事件触发者的信息                                           |
| from.userid  | 操作者的userid                                               |
| response_url | 支持主动回复消息的临时[url](https://developer.work.weixin.qq.com/document/path/100719#59947) |
| msgtype      | 消息类型，此时固定是file。特殊的，**智能机器人目前仅支持100M大小以内的文件回调** |
| file         | 可参考 [文件](https://developer.work.weixin.qq.com/document/path/100719#57141/文件) 结构体说明 |

 

#### 结构体说明

##### 文本

```javascript
{
    "content": "@RobotA hello robot"
}
```

**参数说明：**

| 参数    | 说明         |
| ------- | ------------ |
| content | 文本消息内容 |

##### 图片

```javascript
{
    "url": "https://ww-aibot-img-1258476243.cos.ap-guangzhou.myqcloud.com/BHoPdA3/7571665296904772241?sign=q-sign-algorithm%3Dsha1%26q-ak%3DAKIDbBpaJvdLBvpnibcYcfyPuaO5f9U1UoGo%26q-sign-time%3D1733467811%3B1733468111%26q-key-time%3D1733467811%3B1733468111%26q-header-list%3Dhost%26q-url-param-list%3D%26q-signature%3D0f7b6576943685f82870bc8db306dbf09dfe0fd6 "
}
```

**参数说明：**

| 参数 | 说明                                                         |
| ---- | ------------------------------------------------------------ |
| url  | 图片的下载url, 该url五分钟内有效。注意获取到的文件是已加密的，不能直接打开。加密AESKey与[回调加解密的AESKey](https://developer.work.weixin.qq.com/document/path/100719#12976)相同。加密方式：AES-256-CBC，数据采用**PKCS#7**填充至32字节的倍数；IV初始向量大小为16字节，取AESKey前16字节，详见：https://datatracker.ietf.org/doc/html/rfc2315 |

##### 图文混排

```javascript
{
    "msg_item": [
        {
            "msgtype": "text",
            "text": {
                "content": "@机器人 这是今日的测试情况"
            }
        },
        {
            "msgtype": "image",
            "image": {
                "url": "URL"
            }
        }
    ]
}
```

**参数说明：**

| 参数             | 说明                                                         |
| ---------------- | ------------------------------------------------------------ |
| msg_item.msgtype | 图文混排中的类型，text/image，分别表示：文本和图片           |
| msg_item.text    | 图文混排中的文本内容，可参考 [文本](https://developer.work.weixin.qq.com/document/path/100719#57141/文本) 结构体说明 |
| msg_item.image   | 图文混排中的图片内容，可参考 [图片](https://developer.work.weixin.qq.com/document/path/100719#57141/图片) 结构体说明 |

##### 语音

```javascript
{
    "content": "这是语音转成文本的内容"
}
```

**参数说明：**

| 参数    | 说明                 |
| ------- | -------------------- |
| content | 语音转换成文本的内容 |

 

##### 文件

```javascript
{
    "url": "https://ww-aibot-img-1258476243.cos.ap-guangzhou.myqcloud.com/BHoPdA3/7571665296904772241?sign=q-sign-algorithm%3Dsha1%26q-ak%3DAKIDbBpaJvdLBvpnibcYcfyPuaO5f9U1UoGo%26q-sign-time%3D1733467811%3B1733468111%26q-key-time%3D1733467811%3B1733468111%26q-header-list%3Dhost%26q-url-param-list%3D%26q-signature%3D0f7b6576943685f82870bc8db306dbf09dfe0fd6 "
}
```

**参数说明：**

| 参数 | 说明                                                         |
| ---- | ------------------------------------------------------------ |
| url  | 文件的下载url, 该url五分钟内有效。注意获取到的文件是已加密的，不能直接打开。加密AESKey与[回调加解密的AESKey](https://developer.work.weixin.qq.com/document/path/100719#12976)相同。加密方式：AES-256-CBC，数据采用**PKCS#7**填充至32字节的倍数；IV初始向量大小为16字节，取AESKey前16字节，详见：https://datatracker.ietf.org/doc/html/rfc2315 |

 

##### 引用

```javascript
{
    "msgtype": "text",
    "text": {
        "content": "这是今日的测试情况"
    },
    "image": {
        "url": "URL"
    },
    "mixed": {
        "msg_item": [
            {
                "msgtype": "text",
                "text": {
                    "content": "@机器人 这是今日的测试情况"
                }
            },
            {
                "msgtype": "image",
                "image": {
                    "url": "URL"
                }
            }
        ]
    },
    "voice": {
        "content": "这是语音转成文本的内容"
    },
    "file": {
        "url": "URL"
    }
}
```

**参数说明：**

| 参数    | 说明                                                         |
| ------- | ------------------------------------------------------------ |
| msgtype | 引用的类型，text/image/mixed/voice/file，分别表示：文本/图片/图文混排/语音/文件消息类型 |
| text    | 引用的文本内容，可参考 [文本](https://developer.work.weixin.qq.com/document/path/100719#57141/文本) 结构体说明 |
| image   | 引用的图片内容，可参考 [图片](https://developer.work.weixin.qq.com/document/path/100719#57141/图片) 结构体说明 |
| mixed   | 引用的图文混排内容，可参考 [图文混排](https://developer.work.weixin.qq.com/document/path/100719#57141/图文混排) 结构体说明 |
| voice   | 引用的语音内容，可参考 [语音](https://developer.work.weixin.qq.com/document/path/100719#57141/语音) 结构体说明 |
| file    | 引用的文件内容，可参考 [文件](https://developer.work.weixin.qq.com/document/path/100719#57141/文件) 结构体说明 |

 

### 流式消息刷新

```javascript
{
    "msgid": "CAIQz7/MjQYY/NGagIOAgAMgl8jK/gI=",
    "aibotid": "AIBOTID",
    "chatid": "CHATID",
    "chattype": "group",
    "from": {
        "userid": "USERID"
    },
    "msgtype": "stream",
    "stream": {
        "id": "STREAMID"
    }
}
```

**参数说明：**

| 参数        | 说明                                                         |
| ----------- | ------------------------------------------------------------ |
| msgid       | 本次回调的唯一性标志，开发者需据此进行事件排重（可能因为网络等原因重复回调） |
| aibotid     | 智能机器人id                                                 |
| chatid      | 会话id，仅群聊类型时候返回                                   |
| chattype    | 会话类型，single\group，分别表示：单聊\群聊                  |
| from        | 该事件触发者的信息                                           |
| from.userid | 操作者的userid                                               |
| msgtype     | 消息类型，此时固定是stream。特殊的，**该消息事件仅支持[流式消息的回复](https://developer.work.weixin.qq.com/document/path/100719#57141/被动回复消息格式/流式消息回复)** |
| stream.id   | 流式消息的id，智能机器人根据该id返回对应的流式消息           |

## 接收事件

最后更新：2025/11/25

### 概述

智能机器人回调事件通用协议示例：

```javascript
{
   
    "msgid": "CAIQ16HMjQYY\/NGagIOAgAMgq4KM0AI=",
	"create_time":1700000000,
    "aibotid": "AIBOTID",
	"chatid":"CHATID",
	"chattype":"single",
    "from": {
	 	"corpid": "wpxxxx",
        "userid": "USERID"
    },
    "msgtype": "event",
    "event": {
        "eventtype": "eventtype_name",
		      "eventtype_name":{
			  }
    }
}
```

**参数说明：**

| 参数        | 是否必填 | 说明                                                         |
| ----------- | -------- | ------------------------------------------------------------ |
| msgid       | 是       | 本次回调的唯一性标志，开发者需据此进行事件排重（可能因为网络等原因重复回调） |
| create_time | 是       | 本次回调事件产生的时间                                       |
| aibotid     | 是       | 智能机器人id                                                 |
| chatid      | 否       | 群聊id                                                       |
| chattype    | 是       | 会话类型，single\group，分别表示：单聊\群聊                  |
| from        | 是       | 该事件触发者的信息，详见From结构体                           |
| msgtype     | 是       | 消息类型，若为事件通知，固定为event                          |
| event       | 是       | 若为事件通知，事件结构体，参考Event结构体说明                |

 

**From结构体说明：**

| 参数   | 是否必填 | 说明                                         |
| ------ | -------- | -------------------------------------------- |
| corpid | 否       | 操作者的corpid，若为企业内部智能机器人不返回 |
| userid | 是       | 操作者的userid                               |

**Event结构体说明：**

| 参数           | 是否必填 | 说明                                                         |
| -------------- | -------- | ------------------------------------------------------------ |
| eventtype      | 是       | 事件类型，例如template_card_event                            |
| eventtype_name | 否       | 具体的事件结构体。例如当eventtype为template_card_event时，eventtype_name字段名为template_card_event。具体可参考[模板卡片事件](https://developer.work.weixin.qq.com/document/path/101027#模板卡片事件) |

 

### 事件格式

所有的回调事件都遵循通用协议格式。

#### 进入会话事件

当用户当天首次进入智能机器人单聊会话时，触发该事件。开发者可回复一条文本消息或者模板卡片消息。

> 若未回复消息，用户当天再次进入也不再推送进入会话事件。

**协议格式如下：**

```javascript
{
    "msgid": "CAIQ16HMjQYY\/NGagIOAgAMgq4KM0AI=",
	"create_time":1700000000,
    "aibotid": "AIBOTID",
    "from": {
		"corpid": "wpxxxx",
        "userid": "USERID"
    },
    "msgtype": "event",
    "event": {
        "eventtype": "enter_chat"
    }
}
```

**参数说明：**

| 参数      | 说明                           |
| --------- | ------------------------------ |
| eventtype | 事件类型，此处固定为enter_chat |

 

#### 模板卡片事件

按钮交互、投票选择和多项选择模版卡片中的按钮点击后，企业微信会将相应事件发送给机器人

> 当有模版卡片回调事件的时候，企业微信服务只会发起一次请求，企业微信服务器在五秒内收不到响应会断掉连接，丢弃该回调事件。

模板卡片事件通用协议示例：

```json
{
    "msgid": "CAIQ16HMjQYY\/NGagIOAgAMgq4KM0AI=",
    "create_time": 1700000000,
    "aibotid": "AIBOTID",
    "from": {
        "corpid": "CORPID",
        "userid": "USERID"
    },
    "chatid": "CHATID",
    "chattype": "group",
    "response_url": "RESPONSEURL",
    "msgtype": "event",
    "event": {
        "eventtype": "template_card_event",
        "template_card_event": {
            "card_type": "vote_interaction",
            "event_key": "button_replace_text",
            "task_id": "fBmjTL7ErRCQSNA6GZKMlcFiWX1shOvg",
            "selected_items": {
                "selected_item": [
                    {
                        "question_key": "button_selection_key1",
                        "option_ids": {
                            "option_id": [
                                "button_selection_id1"
                            ]
                        }
                    }
                ]
            }
        }
    }
}
```

其中，eventtype固定为template_card_event。对应结构体TemplateCardEvent。

**参数说明：**

| 参数         | 说明                                                         |
| ------------ | ------------------------------------------------------------ |
| response_url | 支持主动回复消息的[url](https://developer.work.weixin.qq.com/document/path/101027#59947) |

**TemplateCardEvent结构说明：**

| 参数           | 是否必填 | 说明                                              |
| -------------- | -------- | ------------------------------------------------- |
| cardtype       | 是       | 模版卡片的模版类型,此处固定为`button_interaction` |
| eventkey       | 是       | 用户点击的按钮交互模版卡片的按钮key               |
| task_id        | 是       | 用户点击的交互模版卡片的task_id                   |
| selected_items | 否       | 用户点击提交的选择框数据，参考SeletedItem结构说明 |

**参考SeletedItem结构说明：**

| 参数               | 是否必填                                                     | 说明 |
| ------------------ | ------------------------------------------------------------ | ---- |
| question_key       | 用户点击提交的选择框的key值                                  |      |
| optionids.optionid | 用户在选择框选择的数据。当选择框为单选的时候，optionids数组只有一个选项key值; 当选择框为多选的时候，optionids数组可能有多个选项key值 |      |

 

##### 按钮交互模版卡片的事件

当用户点击机器人的按钮交互卡片模块消息的按钮时候，触发相应的回调事件
回调示例

```json
{
    "msgid": "CAIQ16HMjQYY\/NGagIOAgAMgq4KM0AI=",
    "create_time": 1700000000,
    "aibotid": "AIBOTID",
    "from": {
        "corpid": "CORPID",
        "userid": "USERID"
    },
    "chatid": "CHATID",
    "chattype": "group",
    "response_url": "RESPONSEURL",
    "msgtype": "event",
    "event": {
        "eventtype": "template_card_event",
        "template_card_event": {
            "card_type": "button_interaction",
            "event_key": "button_replace_text",
            "task_id": "fBmjTL7ErRCQSNA6GZKMlcFiWX1shOvg",
            "selected_items": {
                "selected_item": [
                    {
                        "question_key": "button_selection_key1",
                        "option_ids": {
                            "option_id": [
                                "button_selection_id1"
                            ]
                        }
                    }
                ]
            }
        }
    }
}
```

**参数说明：**

| 参数           | 说明                                   |
| -------------- | -------------------------------------- |
| cardtype       | 事件类型，此处固定为button_interaction |
| selected_items | 下拉式的选择器选择的数据               |

 

##### 投票选择模版卡片的事件

当用户在选择选项后，点击按钮触发相应的回调事件
回调示例

```json
{
    "msgid": "CAIQ16HMjQYY\/NGagIOAgAMgq4KM0AI=",
    "create_time": 1700000000,
    "aibotid": "AIBOTID",
    "from": {
        "corpid": "CORPID",
        "userid": "USERID"
    },
    "chatid": "CHATID",
    "chattype": "group",
    "response_url": "RESPONSEURL",
    "msgtype": "event",
    "event": {
        "eventtype": "template_card_event",
        "template_card_event": {
            "card_type": "vote_interaction",
            "event_key": "button_replace_text",
            "task_id": "fBmjTL7ErRCQSNA6GZKMlcFiWX1shOvg",
            "selected_items": {
                "selected_item": [
                    {
                        "question_key": "button_selection_key1",
                        "option_ids": {
                            "option_id": [
                                "one",
                                "two"
                            ]
                        }
                    }
                ]
            }
        }
    }
}
```

| 参数           | 说明                                             |
| -------------- | ------------------------------------------------ |
| cardtype       | 模版卡片的模版类型,此处固定为 `vote_interaction` |
| selected_items | 用户点击提交的投票选择框数据                     |

 

##### 多项选择模版卡片的事件

当用户在下拉框选择选项后，点击按钮触发相应的回调事件
回调示例

```json
{
    "msgid": "CAIQ16HMjQYY\/NGagIOAgAMgq4KM0AI=",
    "create_time": 1700000000,
    "aibotid": "AIBOTID",
    "from": {
        "userid": "USERID"
    },
    "chatid": "CHATID",
    "chattype": "group",
    "response_url": "RESPONSEURL",
    "msgtype": "event",
    "event": {
        "eventtype": "template_card_event",
        "template_card_event": {
            "card_type": "multiple_interaction",
            "event_key": "button_replace_text",
            "task_id": "fBmjTL7ErRCQSNA6GZKMlcFiWX1shOvg",
            "selected_items": {
                "selected_item": [
                    {
                        "question_key": "button_selection_key1",
                        "option_ids": {
                            "option_id": [
                                "button_selection_id1"
                            ]
                        }
                    },
                    {
                        "question_key": "button_selection_key2",
                        "option_ids": {
                            "option_id": [
                                "button_selection_id2"
                            ]
                        }
                    }
                ]
            }
        }
    }
}
```

| 参数           | 说明                                                 |
| -------------- | ---------------------------------------------------- |
| cardtype       | 模版卡片的模版类型,此处固定为 `multiple_interaction` |
| selected_items | 用户点击提交的下拉菜单选择框列表数据                 |

##### 模版卡片右上角菜单事件

用户点击文本通知，图文展示和按钮交互卡片右上角的菜单时会弹出菜单选项，当用户点击具体选项的时候会触发相应的回调事件
回调示例

```json
{
    "msgid": "CAIQ16HMjQYY\/NGagIOAgAMgq4KM0AI=",
    "create_time": 1700000000,
    "aibotid": "AIBOTID",
    "from": {
        "userid": "USERID"
    },
    "chatid": "CHATID",
    "chattype": "group",
    "response_url": "RESPONSEURL",
    "msgtype": "event",
    "event": {
        "eventtype": "template_card_event",
        "template_card_event": {
            "card_type": "text_notice",
            "event_key": "button_replace_text",
            "task_id": "fBmjTL7ErRCQSNA6GZKMlcFiWX1shOvg"
        }
    }
}
```

| 参数     | 说明                                                         |
| -------- | ------------------------------------------------------------ |
| cardtype | 模版卡片的模版类型,此处可能为 `text_notice` , `news_notice` 和 `button_interaction` |

#### 用户反馈事件

开发者接收到智能机器人的消息事件后，可以将事件与即将回复的流式消息关联起来，并在被动回复/主动回复消息中设置反馈信息。当用户进行反馈时，可以收到用户反馈事件，复盘智能机器人的回复效果
![img](https://wework.qpic.cn/wwpic3az/951882_QQWUKEVcS-yEm3B_1764061549/0)

> 若智能机器人回复的消息中未设置反馈信息，则用户点击反馈后不会触发反馈事件
> 该事件目前仅**支持回复空包**, 不支持回复新的消息或者更新卡片消息

```json
{
    "msgid": "CAIQ16HMjQYY\/NGagIOAgAMgq4KM0AI=",
    "create_time": 1700000000,
    "aibotid": "AIBOTID",
    "chatid": "CHATID",
    "chattype": "group",
    "from": {
        "userid": "USERID"
    },
    "msgtype": "event",
    "event": {
        "eventtype": "feedback_event",
        "feedback_event": {
            "id": "FEEDBACKID",
            "type": 2,
            "content": "能再详细一些么",
            "inaccurate_reason_list": [
                2,
                4
            ]
        }
    }
}
```

**参数说明：**

| 参数                                  | 说明                                                         |
| ------------------------------------- | ------------------------------------------------------------ |
| eventtype                             | 事件类型，此处固定为feedback_event                           |
| feedback_event.id                     | [回复消息](https://developer.work.weixin.qq.com/document/path/101027#59068/流式消息回复)设置的反馈id |
| feedback_event.type                   | 反馈的类型： 1：准确  2：不准确  3：取消准确/不准确          |
| feedback_event.content                | 用户输入的反馈内容，仅不准确反馈类型支持返回                 |
| feedback_event.inaccurate_reason_list | 用户选择的负反馈原因列表，仅不准确反馈类型支持返回。负反馈原因有： 1: 与问题无关  2: 内容不完整  3: 内容有错误  4: 数据分析错误 |

## 被动回复消息

最后更新：2025/10/21

### 概述

当用户与智能机器人进行交互时，企业微信会将相关的交互事件回调到开发者设置的回调URL，开发者可根据事件类型做出相应的响应，实现丰富的自定义功能。

目前主要有以下场景支持回复消息：

1. 用户当天首次进入智能机器人单聊会话，回复欢迎语
2. 用户向智能机器人发送消息 ，回复消息
3. 用户点击模板卡片相关按钮等，回复消息更新模板卡片

 

> 具体可参考[接收回调与被动回复消息流程](https://developer.work.weixin.qq.com/document/path/101031#59174/接收回调与被动回复消息)

 

### 回复欢迎语

#### 文本消息

```json
{
  "msgtype": "text",
  "text": {
    "content": "hello\nI'm RobotA\n"
  }
}
```

**参数说明：**

| 参数         | 类型   | 必须 | 说明               |
| ------------ | ------ | ---- | ------------------ |
| msgtype      | String | 是   | 此时固定为text     |
| text         | Object | 是   | 文本消息的详细内容 |
| text.content | String | 是   | 文本内容           |

> 目前仅支持进入会话回调事件时，支持被动回复文本消息

#### 模板卡片消息

### 回复用户消息

#### 流式消息回复
                    }
                ]
            }
        ],
        "submit_button": {
            "text": "提交",
            "key": "submit_key"
        },
        "task_id": "task_id"
    }
}
```

**参数说明：**

| 参数          | 类型   | 必须 | 说明                                                         |
| ------------- | ------ | ---- | ------------------------------------------------------------ |
| msgtype       | String | 是   | 此时固定为template_card                                      |
| template_card | Object | 是   | 模板卡片TemplateCard结构体。参考[模板卡片类型](https://developer.work.weixin.qq.com/document/path/101031#59098)中类型说明参考[模板卡片类型](https://developer.work.weixin.qq.com/document/path/101031#59098)中类型说明 |

> 目前仅支持进入会话回调事件或者接收消息回调时，支持被动回复模板卡片消息

# 回复用户消息

## 流式消息回复

```javascript
{
    "msgtype": "stream",
    "stream": {
        "id": "STREAMID",
        "finish": false,
        "content": "**广州**今日天气：29度，大部分多云，降雨概率：60%",
        "msg_item": [
            {
                "msgtype": "image",
                "image": {
                    "base64": "BASE64",
                    "md5": "MD5"
                }
            }
        ],
        "feedback": {
            "id": "FEEDBACKID"
        }
    }
}
```

**参数说明：**

| 参数                    | 类型     | 是否必填                           | 说明                                                         |
| ----------------------- | -------- | ---------------------------------- | ------------------------------------------------------------ |
| msgtype                 | String   | 是                                 | 消息类型，此时固定为：stream                                 |
| stream.id               | String   | 否，流式消息首次回复的时候要求设置 | 自定义的唯一id，某次回调的首次回复第一次设置，后续回调会根据这个id来获取最新的流式消息 |
| stream.finish           | Bool     | 否                                 | 流式消息是否结束                                             |
| stream.content          | String   | 否                                 | 流式消息内容，最长不超过20480个字节，必须是utf8编码。 特殊的，第一次回复内容为"1"，第二次回复"123"，则此时消息展示内容"123" |
| stream.msg_item         | Object[] | 否                                 | 流式消息图文混排消息列表。                                   |
| stream.msg_item.msgtype | String   | 否                                 | 图文混排消息类型，目前支持：image 特殊的，目前image的消息类型仅当finish=true，即流式消息结束的最后一次回复中设置 |
| stream.msg_item.image   | Object   | 否                                 | 图片混排的图片资源。目前最多支持设置10个                     |
| image.base64            | String   | 是                                 | 图片内容的base64编码。 图片（base64编码前）最大不能超过10M，支持JPG,PNG格式 |
| image.md5               | String   | 是                                 | 图片内容（base64编码前）的md5值                              |
| stream.feedback.id      | String   | 否                                 | 流式消息首次回复时候若字段不为空值，回复的消息被用户反馈时候会触发[回调事件](https://developer.work.weixin.qq.com/document/path/101031#59058/用户反馈事件)。有效长度为 256 字节以内，必须是 utf-8 编码。 |

> 流式消息回复内容content字段支持常见的[markdown格式](https://developer.work.weixin.qq.com/document/path/101031#14404/markdown-v2类型)
> 若content中包含思考过程<think></think>标签，客户端会展示思考过程。
> 若回复内容包含图片，仅支持在最后一次回复时，即finish=true时支持包含msgtype为image的msg_item。

> 目前仅支持进入会话回调事件或者接收消息回调时，支持被动回复模板卡片消息

## 模板卡片消息

```json
{
    "msgtype": "template_card",
    "template_card": {
        "feedback": {
            "id": "FEEDBACKID"
        }
    }
}
```

**参数说明：**

| 参数                      | 类型   | 是否必填 | 说明                                                         |
| ------------------------- | ------ | -------- | ------------------------------------------------------------ |
| msgtype                   | String | 是       | 消息类型，此时固定为：template_card                          |
| template_card             | Object | 是       | 模板卡片结构体，参考[模板卡片类型](https://developer.work.weixin.qq.com/document/path/101031#59098)中类型说明 |
| template_card.feedback.id | String | 否       | 特殊的该回复场景支持设置反馈信息。若字段不为空值，回复的消息被用户反馈时候会触发[回调事件](https://developer.work.weixin.qq.com/document/path/101031#59058/用户反馈事件)。有效长度为 256 字节以内，必须是 utf-8 编码。 |

## 流式消息+模板卡片回复

若开发者需要回复流式消息外，还需要回复模板卡片，可回复该消息类型。

> 目前仅支持进入会话回调事件或者接收消息回调时，支持被动回复模板卡片消息
> 首次回复时必须返回stream的id。
> template_card可首次回复，也可在收到流式消息刷新事件时回复。但是同一个消息只能回复一次。

```json
{
    "msgtype": "stream_with_template_card",
    "stream": {
        "id": "STREAMID",
        "finish": false,
        "content": "**广州**今日天气：29度，大部分多云，降雨概率：60%",
        "msg_item": [
            {
                "msgtype": "image",
                "image": {
                    "base64": "BASE64",
                    "md5": "MD5"
                }
            }
        ],
        "feedback": {
            "id": "FEEDBACKID"
        }
    },
    "template_card": {
        "feedback": {
            "id": "FEEDBACKID"
        }
    }
}
```

**参数说明：**

| 参数                      | 类型   | 是否必填 | 说明                                                         |
| ------------------------- | ------ | -------- | ------------------------------------------------------------ |
| msgtype                   | String | 是       | 消息类型，此时固定为：stream_with_template_card              |
| stream                    | Object | 是       | 参考[流式回复消息](https://developer.work.weixin.qq.com/document/path/101031#59068/流式回复消息)说明。 |
| stream.feedback.id        | String | 否       | 特殊的该回复场景支持设置反馈信息。流式消息首次回复时候若字段不为空值，回复的消息被用户反馈时候会触发[回调事件](https://developer.work.weixin.qq.com/document/path/101031#59058/用户反馈事件)。有效长度为 256 字节以内，必须是 utf-8 编码。 |
| template_card             | Object | 否       | 模板卡片结构体，参考[模板卡片类型](https://developer.work.weixin.qq.com/document/path/101031#59098)中类型说明 |
| template_card.feedback.id | String | 否       | 特殊的该回复场景支持设置反馈信息。若字段不为空值，回复的消息被用户反馈时候会触发[回调事件](https://developer.work.weixin.qq.com/document/path/101031#59058/用户反馈事件)。有效长度为 256 字节以内，必须是 utf-8 编码。 |

 

### 回复消息更新模板卡片

#### 模版卡片更新消息

当机器人服务接收到[模版卡片事件](https://developer.work.weixin.qq.com/document/path/101031#59058/模板卡片事件)后，可以在该事件的返回包中添加消息进行即时响应。

```json
{
    "response_type": "update_template_card",
    "userids": [
        "USERID1",
        "USERID2"
    ],
    "template_card": {
        "feedback": {
            "id": "FEEDBACKID"
        }
    }
}
```

**参数说明**

| 参数                      | 类型     | 必须 | 说明                                                         |
| ------------------------- | -------- | ---- | ------------------------------------------------------------ |
| response_type             | String   | 是   | 响应类型，此处固定为 `update_template_card` - 替换部分用户的模版 |
| userids                   | String[] | 否   | 表示要替换模版卡片消息的userid列表。若不填，则表示替换当前消息涉及到的所有用户。开发者可以通过[模版卡片事件](https://developer.work.weixin.qq.com/document/path/101031#59058/模板卡片事件)中获取userid |
| template_card             | Object   | 是   | 要替换的模版卡片TemplateCard结构体。参考[模板卡片类型](https://developer.work.weixin.qq.com/document/path/101031#59098)中类型说明 |
| template_card.feedback.id | String   | 否   | 特殊的该回复场景支持设置反馈信息，替换用户模板会覆盖原先消息的反馈信息。若字段不为空值，回复的消息被用户反馈时候会触发[回调事件](https://developer.work.weixin.qq.com/document/path/101031#59058/用户反馈事件)。有效长度为 256 字节以内，必须是 utf-8 编码。 |

> 注：模板卡片中的task_id需跟回调收到的task_id一致

## 模板卡片类型

最后更新：2025/06/26

### 模版卡片类型

该文档主要说明各种类型模板卡片**TemplateCard结构体说明**。

> 其中，点击文本通知卡片以及图文通知卡片的“跳转指引”区域支持消息智能回复。

 

#### 文本通知模版卡片

文本通知模版卡片消息示例
![img](https://wework.qpic.cn/wwpic/262807_8RCSsMfbSAaGBYh_1633781903/0)完整文本通知模版卡片示例

```javascript
{
    "card_type": "text_notice",
    "source": {
        "icon_url": "https://wework.qpic.cn/wwpic/252813_jOfDHtcISzuodLa_1629280209/0",
        "desc": "企业微信",
        "desc_color": 0
    },
    "action_menu": {
        "desc": "消息气泡副交互辅助文本说明",
        "action_list": [
            {
                "text": "接收推送",
                "key": "action_key1"
            },
            {
                "text": "不再推送",
                "key": "action_key2"
            }
        ]
    },
    "main_title": {
        "title": "欢迎使用企业微信",
        "desc": "您的好友正在邀请您加入企业微信"
    },
    "emphasis_content": {
        "title": "100",
        "desc": "数据含义"
    },
    "quote_area": {
        "type": 1,
        "url": "https://work.weixin.qq.com/?from=openApi",
        "appid": "APPID",
        "pagepath": "PAGEPATH",
        "title": "引用文本标题",
        "quote_text": "Jack：企业微信真的很好用~\nBalian：超级好的一款软件！"
    },
    "sub_title_text": "下载企业微信还能抢红包！",
    "horizontal_content_list": [
        {
            "keyname": "邀请人",
            "value": "张三"
        },
        {
            "keyname": "企微官网",
            "value": "点击访问",
            "type": 1,
            "url": "https://work.weixin.qq.com/?from=openApi"
        },
        {
            "keyname": "企微下载",
            "value": "企业微信.apk",
            "type": 1,
            "url": "https://work.weixin.qq.com/?from=openApi"
        }
    ],
    "jump_list": [
        {
            "type": 1,
            "url": "https://work.weixin.qq.com/?from=openApi",
            "title": "企业微信官网"
        },
        {
            "type": 2,
            "appid": "APPID",
            "pagepath": "PAGEPATH",
            "title": "跳转小程序"
        },
        {
            "type": 3,
            "title": "企业微信官网",
            "question": "如何登录企业微信官网"
        }
    ],
    "card_action": {
        "type": 1,
        "url": "https://work.weixin.qq.com/?from=openApi",
        "appid": "APPID",
        "pagepath": "PAGEPATH"
    },
    "task_id": "task_id"
}
```

请求参数

| 参数                    | 类型     | 必须 | 说明                                                         |
| ----------------------- | -------- | ---- | ------------------------------------------------------------ |
| card_type               | String   | 是   | 模版卡片的模版类型，文本通知模版卡片的类型为`text_notice`    |
| source                  | Object   | 否   | 卡片来源样式信息，不需要来源样式可不填写。参考**Source结构体说明** |
| action_menu             | Object   | 否   | 卡片右上角更多操作按钮。参考**ActionMenu结构体说明**         |
| main_title              | Object   | 否   | 模版卡片的主要内容，包括一级标题和标题辅助信息。参考**MainTitle结构体说明** |
| emphasis_content        | Object   | 否   | 关键数据样式，建议不与引用样式共用。参考**EmphasisContent结构体说明** |
| quote_area              | Object   | 否   | 引用文献样式，建议不与关键数据共用。参考**QuoteArea结构体说明** |
| sub_title_text          | String   | 否   | 二级普通文本，建议不超过112个字。**模版卡片主要内容的一级标题main_title.title和二级普通文本sub_title_text必须有一项填写** |
| horizontal_content_list | Object[] | 否   | 二级标题+文本列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过6。参考**HorizontalContent结构体说明** |
| jump_list               | Object[] | 否   | 跳转指引样式的列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过3。参考**JumpAction结构体说明** |
| card_action             | Object   | 是   | 整体卡片的点击跳转事件，text_notice模版卡片中该字段为必填项。参考**CardAction结构体说明** |
| task_id                 | String   | 否   | 任务id，当文本通知模版卡片有action_menu字段的时候，该字段必填。同一个机器人任务id不能重复，只能由数字、字母和“_-@”组成，最长128字节。任务id只在发消息时候有效，更新消息的时候无效。任务id将会在相应的回调事件中返回 |

 

#### 图文展示模版卡片

图文展示模版卡片消息示例
![img](https://wework.qpic.cn/wwpic/602361_D5DSN3MBSFOqcGb_1633781666/0)完整图文展示模版卡片示例

```javascript
{
    "card_type": "news_notice",
    "source": {
        "icon_url": "https://wework.qpic.cn/wwpic/252813_jOfDHtcISzuodLa_1629280209/0",
        "desc": "企业微信",
        "desc_color": 0
    },
    "action_menu": {
        "desc": "消息气泡副交互辅助文本说明",
        "action_list": [
            {
                "text": "接收推送",
                "key": "action_key1"
            },
            {
                "text": "不再推送",
                "key": "action_key2"
            }
        ]
    },
    "main_title": {
        "title": "欢迎使用企业微信",
        "desc": "您的好友正在邀请您加入企业微信"
    },
    "card_image": {
        "url": "https://wework.qpic.cn/wwpic/354393_4zpkKXd7SrGMvfg_1629280616/0",
        "aspect_ratio": 2.25
    },
    "image_text_area": {
        "type": 1,
        "url": "https://work.weixin.qq.com",
        "title": "欢迎使用企业微信",
        "desc": "您的好友正在邀请您加入企业微信",
        "image_url": "https://wework.qpic.cn/wwpic/354393_4zpkKXd7SrGMvfg_1629280616/0"
    },
    "quote_area": {
        "type": 1,
        "url": "https://work.weixin.qq.com/?from=openApi",
        "appid": "APPID",
        "pagepath": "PAGEPATH",
        "title": "引用文本标题",
        "quote_text": "Jack：企业微信真的很好用~\nBalian：超级好的一款软件！"
    },
    "vertical_content_list": [
        {
            "title": "惊喜红包等你来拿",
            "desc": "下载企业微信还能抢红包！"
        }
    ],
    "horizontal_content_list": [
        {
            "keyname": "邀请人",
            "value": "张三"
        },
        {
            "keyname": "企微官网",
            "value": "点击访问",
            "type": 1,
            "url": "https://work.weixin.qq.com/?from=openApi"
        },
        {
            "keyname": "企微下载",
            "value": "企业微信.apk",
            "type": 1,
            "url": "https://work.weixin.qq.com/?from=openApi"
        }
    ],
    "jump_list": [
        {
            "type": 1,
            "url": "https://work.weixin.qq.com/?from=openApi",
            "title": "企业微信官网"
        },
        {
            "type": 2,
            "appid": "APPID",
            "pagepath": "PAGEPATH",
            "title": "跳转小程序"
        },
        {
            "type": 3,
            "title": "企业微信官网",
            "question": "如何登录企业微信官网"
        }
    ],
    "card_action": {
        "type": 1,
        "url": "https://work.weixin.qq.com/?from=openApi",
        "appid": "APPID",
        "pagepath": "PAGEPATH"
    },
    "task_id": "task_id"
}
```

| 参数                    | 类型     | 必须 | 说明                                                         |
| ----------------------- | -------- | ---- | ------------------------------------------------------------ |
| card_type               | String   | 是   | 模版卡片的模版类型，图文展示模版卡片的类型为`news_notice`    |
| source                  | Object   | 否   | 卡片来源样式信息，不需要来源样式可不填写。参考**Source结构体说明** |
| action_menu             | Object   | 否   | 卡片右上角更多操作按钮。参考**ActionMenu结构体说明**         |
| main_title              | Object   | 是   | 模版卡片的主要内容，包括一级标题和标题辅助信息。参考**MainTitle结构体说明** |
| card_image              | Object   | 否   | 图片样式，news_notice类型的卡片，card_image和image_text_area两者必填一个字段，不可都不填。参考**CardImage结构体说明** |
| image_text_area         | Object   | 否   | 左图右文样式。参考**ImageTextArea结构体说明**                |
| vertical_content_list   | Object[] | 否   | 卡片二级垂直内容，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过4。参考**VerticalContent结构体说明** |
| horizontal_content_list | Object[] | 否   | 二级标题+文本列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过6。参考**HorizontalContent结构体说明** |
| jump_list               | Object[] | 否   | 跳转指引样式的列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过3。参考**JumpAction结构体说明** |
| card_action             | Object   | 是   | 整体卡片的点击跳转事件，news_notice模版卡片中该字段为必填项。参考**CardAction结构体说明** |
| task_id                 | String   | 否   | 任务id，当图文展示模版卡片有action_menu字段的时候，该字段必填。同一个机器人任务id不能重复，只能由数字、字母和“_-@”组成，最长128字节。任务id只在发消息时候有效，更新消息的时候无效。任务id将会在相应的回调事件中返回 |

 

#### 按钮交互模版卡片

按钮交互模版卡片消息示例
![img](https://wework.qpic.cn/wwpic/93152_bcog7qR3R3qTsrr_1633786928/0)完整按钮交互模版卡片示例

```json
{
    "card_type": "button_interaction",
    "source": {
        "icon_url": "https://wework.qpic.cn/wwpic/252813_jOfDHtcISzuodLa_1629280209/0 ",
        "desc": "企业微信",
        "desc_color": 0
    },
    "action_menu": {
        "desc": "消息气泡副交互辅助文本说明",
        "action_list": [
            {
                "text": "接收推送",
                "key": "action_key1"
            },
            {
                "text": "不再推送",
                "key": "action_key2"
            }
        ]
    },
    "main_title": {
        "title": "欢迎使用企业微信",
        "desc": "您的好友正在邀请您加入企业微信"
    },
    "quote_area": {
        "type": 1,
        "url": "https://work.weixin.qq.com/?from=openApi ",
        "appid": "APPID",
        "pagepath": "PAGEPATH",
        "title": "引用文本标题",
        "quote_text": "Jack：企业微信真的很好用~\nBalian：超级好的一款软件！"
    },
    "sub_title_text": "下载企业微信还能抢红包！",
    "horizontal_content_list": [
        {
            "keyname": "邀请人",
            "value": "张三"
        },
        {
            "keyname": "企微官网",
            "value": "点击访问",
            "type": 1,
            "url": "https://work.weixin.qq.com/?from=openApi"
        },
        {
            "keyname": "企微下载",
            "value": "企业微信.apk",
            "type": 1,
            "url": "https://work.weixin.qq.com/?from=openApi"
        }
    ],
    "button_selection": {
        "question_key": "button_selection_key1",
        "title": "你的身份",
        "disable": false,
        "option_list": [
            {
                "id": "button_selection_id1",
                "text": "企业负责人"
            },
            {
                "id": "button_selection_id2",
                "text": "企业用户"
            }
        ],
        "selected_id": "button_selection_id1"
    },
    "button_list": [
        {
            "text": "按钮1",
            "style": 4,
            "key": "BUTTONKEYONE"
        },
        {
            "text": "按钮2",
            "style": 1,
            "key": "BUTTONKEYTWO"
        }
    ],
    "card_action": {
        "type": 1,
        "url": "https://work.weixin.qq.com/?from=openApi ",
        "appid": "APPID",
        "pagepath": "PAGEPATH"
    },
    "task_id": "task_id"
}
```

| 参数                    | 类型     | 必须 | 说明                                                         |
| ----------------------- | -------- | ---- | ------------------------------------------------------------ |
| card_type               | String   | 是   | 模版卡片的模版类型，按钮交互模版卡片的类型为`button_interaction`。当机器人设置了回调URL时，才能下发按钮交互模版卡片 |
| source                  | Object   | 否   | 卡片来源样式信息，不需要来源样式可不填写。参考**Source结构体说明** |
| action_menu             | Object   | 否   | 卡片右上角更多操作按钮。参考**ActionMenu结构体说明**         |
| main_title              | Object   | 是   | 模版卡片的主要内容，包括一级标题和标题辅助信息。参考**MainTitle结构体说明** |
| quote_area              | Object   | 否   | 引用文献样式，建议不与关键数据共用。参考**QuoteArea结构体说明** |
| sub_title_text          | String   | 否   | 二级普通文本，建议不超过112个字                              |
| horizontal_content_list | Object[] | 否   | 二级标题+文本列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过6。参考**HorizontalContent结构体说明** |
| button_selection        | Object   | 否   | 下拉式的选择器。参考**SelectionItem结构体说明**              |
| button_list             | Object[] | 是   | 按钮列表，列表长度不超过6。参考**Button结构体说明结构体说明** |
| card_action             | Object   | 否   | 整体卡片的点击跳转事件。参考**CardAction结构体说明**         |
| task_id                 | String   | 是   | 任务id，同一个机器人任务id不能重复，只能由数字、字母和“_-@”组成，最长128字节。任务id只在发消息时候有效，更新消息的时候无效。任务id将会在相应的回调事件中返回 |

 

## 投票选择模版卡片

投票选择模版卡片消息示例
![img](https://wework.qpic.cn/wwpic/713860_X2yj2lMbR2eBsCl_1629279992/0)

完整投票选择模版卡片示例

```json
{
    "card_type": "vote_interaction",
    "source": {
        "icon_url": "https://wework.qpic.cn/wwpic/252813_jOfDHtcISzuodLa_1629280209/0 ",
        "desc": "企业微信"
    },
    "main_title": {
        "title": "欢迎使用企业微信",
        "desc": "您的好友正在邀请您加入企业微信"
    },
    "checkbox": {
        "question_key": "question_key",
        "option_list": [
            {
                "id": "id_one",
                "text": "选择题选项1"
            },
            {
                "id": "id_two",
                "text": "选择题选项2",
                "is_checked": true
            }
        ],
        "disable": false,
        "mode": 1
    },
    "submit_button": {
        "text": "提交",
        "key": "submit_key"
    },
    "task_id": "task_id"
}
```

| 参数          | 类型   | 必须 | 说明                                                         |
| ------------- | ------ | ---- | ------------------------------------------------------------ |
| card_type     | String | 是   | 模版卡片的模版类型，投票选择模版卡片的类型为`vote_interaction`。当机器人设置了回调URL时，才能下发投票选择模版卡片 |
| source        | Object | 否   | 卡片来源样式信息，不需要来源样式可不填写。参考**Source结构体说明** |
| main_title    | Object | 是   | 模版卡片的主要内容，包括一级标题和标题辅助信息。参考**MainTitle结构体说明** |
| checkbox      | Object | 是   | 选择题样式。参考**CheckBox结构体说明**                       |
| submit_button | Object | 是   | 提交按钮样式。参考**SubmitButtion结构体说明**                |
| task_id       | String | 是   | 任务id，同一个机器人任务id不能重复，只能由数字、字母和“_-@”组成，最长128字节。任务id只在发消息时候有效，更新消息的时候无效。任务id将会在相应的回调事件中返回 |

 

#### 多项选择模版卡片

投票选择模版卡片消息示例
![img](https://wework.qpic.cn/wwpic/151585_1bMIFL0dQR-3cyd_1629280033/0)完整多项选择模版卡片示例

```json
{
    "card_type": "multiple_interaction",
    "source": {
        "icon_url": "https://wework.qpic.cn/wwpic/252813_jOfDHtcISzuodLa_1629280209/0 ",
        "desc": "企业微信"
    },
    "main_title": {
        "title": "欢迎使用企业微信",
        "desc": "您的好友正在邀请您加入企业微信"
    },
    "select_list": [
        {
            "question_key": "question_key_one",
            "title": "选择标签1",
            "disable": false,
            "selected_id": "id_one",
            "option_list": [
                {
                    "id": "id_one",
                    "text": "选择器选项1"
                },
                {
                    "id": "id_two",
                    "text": "选择器选项2"
                }
            ]
        },
        {
            "question_key": "question_key_two",
            "title": "选择标签2",
            "selected_id": "id_three",
            "option_list": [
                {
                    "id": "id_three",
                    "text": "选择器选项3"
                },
                {
                    "id": "id_four",
                    "text": "选择器选项4"
                }
            ]
        }
    ],
    "submit_button": {
        "text": "提交",
        "key": "submit_key"
    },
    "task_id": "task_id"
}
```

| 参数          | 类型     | 必须 | 说明                                                         |
| ------------- | -------- | ---- | ------------------------------------------------------------ |
| card_type     | String   | 是   | 模版卡片的模版类型，多项选择模版卡片的类型为`multiple_interaction`。当机器人设置了回调URL时，才能下发多项选择模版卡片 |
| source        | Object   | 否   | 卡片来源样式信息，不需要来源样式可不填写。参考**Source结构体说明** |
| main_title    | Object   | 是   | 模版卡片的主要内容，包括一级标题和标题辅助信息。参考**MainTitle结构体说明** |
| select_list   | Object[] | 是   | 下拉式的选择器列表，multiple_interaction类型的卡片该字段不可为空，一个消息最多支持 3 个选择器。参考**SelectionItem结构体说明** |
| submit_button | Object   | 是   | 提交按钮样式。参考**SubmitButton结构体说明**                 |
| task_id       | String   | 否   | 任务id，同一个机器人任务id不能重复，只能由数字、字母和“_-@”组成，最长128字节。任务id只在发消息时候有效，更新消息的时候无效。任务id将会在相应的回调事件中返回 |

 

### 结构体说明

#### Source结构体

卡片来源样式信息

| 参数       | 类型   | 必须 | 说明                                                         |
| ---------- | ------ | ---- | ------------------------------------------------------------ |
| icon_url   | String | 否   | 来源图片的url                                                |
| desc       | String | 否   | 来源图片的描述，建议不超过13个字                             |
| desc_color | Int    | 否   | 来源文字的颜色，目前支持：0(默认) 灰色，1 黑色，2 红色，3 绿色 |

#### ActionMenu结构体

卡片右上角更多操作按钮

| 参数             | 类型   | 必须 | 说明                                                         |
| ---------------- | ------ | ---- | ------------------------------------------------------------ |
| desc             | String | 是   | 更多操作界面的描述                                           |
| action_list      | Int    | 是   | 操作列表，列表长度取值范围为 [1, 3]                          |
| action_list.text | String | 是   | 操作的描述文案                                               |
| action_list.key  | String | 是   | 操作key值，用户点击后，会产生回调事件将本参数作为EventKey返回，回调事件会带上该key值，最长支持1024字节，不可重复 |

#### MainTitle结构体

模版卡片的主要内容，包括一级标题和标题辅助信息

| 参数  | 类型   | 必须 | 说明                                                         |
| ----- | ------ | ---- | ------------------------------------------------------------ |
| title | String | 否   | 一级标题，建议不超过26个字。**模版卡片主要内容的一级标题main_title.title和二级普通文本sub_title_text必须有一项填写** |
| desc  | String | 否   | 标题辅助信息，建议不超过30个字                               |

#### EmphasisContent结构体

关键数据样式

| 参数  | 类型   | 必须 | 说明                                         |
| ----- | ------ | ---- | -------------------------------------------- |
| title | String | 否   | 关键数据样式的数据内容，建议不超过10个字     |
| desc  | String | 否   | 关键数据样式的数据描述内容，建议不超过15个字 |

#### QuoteArea结构体

引用文献样式

| 参数       | 类型   | 必须 | 说明                                                         |
| ---------- | ------ | ---- | ------------------------------------------------------------ |
| type       | Int    | 否   | 引用文献样式区域点击事件，0或不填代表没有点击事件，1 代表跳转url，2 代表跳转小程序 |
| url        | String | 否   | 点击跳转的url，type是1时必填                                 |
| appid      | String | 否   | 点击跳转的小程序的appid，必须是与当前应用关联的小程序，type是2时必填 |
| pagepath   | String | 否   | 点击跳转的小程序的pagepath，type是2时选填                    |
| title      | String | 否   | 引用文献样式的标题                                           |
| quote_text | String | 否   | 引用文献样式的引用文案                                       |

#### HorizontalContent结构体

二级标题+文本列表

| 参数    | 类型   | 必须 | 说明                                                         |
| ------- | ------ | ---- | ------------------------------------------------------------ |
| type    | Int    | 否   | 链接类型，0或不填代表是普通文本，1 代表跳转url，3 代表点击跳转成员详情 |
| keyname | String | 是   | 二级标题，建议不超过5个字                                    |
| value   | String | 否   | 二级文本，建议不超过26个字                                   |
| url     | String | 否   | 链接跳转的url，type是1时必填                                 |
| userid  | String | 否   | 成员详情的userid，type是3时必填                              |

#### JumpAction结构体

跳转指引样式的列表

| 参数     | 类型   | 必须 | 说明                                                         |
| -------- | ------ | ---- | ------------------------------------------------------------ |
| type     | Int    | 否   | 跳转链接类型，0或不填代表不是链接，1 代表跳转url，2 代表跳转小程序，3 代表触发消息智能回复 |
| question | String | 否   | 智能问答问题，最长不超过200个字节。若type为3，必填           |
| title    | String | 是   | 跳转链接样式的文案内容，建议不超过13个字                     |
| url      | String | 否   | 跳转链接的url，type是1时必填                                 |
| appid    | String | 否   | 跳转链接的小程序的appid，type是2时必填                       |
| pagepath | String | 否   | 跳转链接的小程序的pagepath，type是2时选填                    |

#### CardAction结构体

整体卡片的点击跳转事件

| 参数     | 类型   | 必须 | 说明                                                         |
| -------- | ------ | ---- | ------------------------------------------------------------ |
| type     | Int    | 是   | 卡片跳转类型，0或不填代表不是链接，1 代表跳转url，2 代表打开小程序。text_notice模版卡片中该字段取值范围为[1,2] |
| url      | String | 否   | 跳转事件的url，type是1时必填                                 |
| appid    | String | 否   | 跳转事件的小程序的appid，type是2时必填                       |
| pagepath | String | 否   | 跳转事件的小程序的pagepath，type是2时选填                    |

#### VerticalContent结构体

卡片二级垂直内容

| 参数  | 类型   | 必须 | 说明                            |
| ----- | ------ | ---- | ------------------------------- |
| title | String | 是   | 卡片二级标题，建议不超过26个字  |
| desc  | String | 否   | 二级普通文本，建议不超过112个字 |

#### CardImage结构体

图片样式

| 参数         | 类型   | 必须 | 说明                                                       |
| ------------ | ------ | ---- | ---------------------------------------------------------- |
| url          | Object | 是   | 图片的url                                                  |
| aspect_ratio | Float  | 否   | 图片的宽高比，宽高比要小于2.25，大于1.3，不填该参数默认1.3 |

#### ImageTextArea结构体

左图右文样式

| 参数      | 类型   | 必须 | 说明                                                         |
| --------- | ------ | ---- | ------------------------------------------------------------ |
| type      | Int    | 否   | 左图右文样式区域点击事件，0或不填代表没有点击事件，1 代表跳转url，2 代表跳转小程序 |
| url       | String | 否   | 点击跳转的url，type是1时必填                                 |
| appid     | String | 否   | 点击跳转的小程序的appid，必须是与当前应用关联的小程序，type是2时必填 |
| pagepath  | String | 否   | 点击跳转的小程序的pagepath，type是2时选填                    |
| title     | String | 否   | 左图右文样式的标题                                           |
| desc      | String | 否   | 左图右文样式的描述                                           |
| image_url | String | 是   | 左图右文样式的图片url                                        |

#### SubmitButton结构体

提交按钮样式

| 参数 | 类型   | 必须 | 说明                                                         |
| ---- | ------ | ---- | ------------------------------------------------------------ |
| text | String | 是   | 按钮文案，建议不超过10个字                                   |
| key  | String | 是   | 提交按钮的key，会产生回调事件将本参数作为EventKey返回，最长支持1024字节 |

#### SelectionItem结构体

下拉式的选择器列表

| 参数             | 类型     | 必须 | 说明                                                         |
| ---------------- | -------- | ---- | ------------------------------------------------------------ |
| question_key     | String   | 是   | 下拉式的选择器题目的key，用户提交选项后，会产生回调事件，回调事件会带上该key值表示该题，最长支持1024字节，不可重复 |
| title            | String   | 否   | 选择器的标题，建议不超过13个字                               |
| disable          | Bool     | 否   | 下拉式的选择器是否不可选，false为可选，true为不可选。仅在更新模版卡片的时候该字段有效 |
| selected_id      | String   | 否   | 默认选定的id，不填或错填默认第一个                           |
| option_list      | Object[] | 是   | 选项列表，下拉选项不超过 10 个，最少1个                      |
| option_list.id   | String   | 是   | 下拉式的选择器选项的id，用户提交选项后，会产生回调事件，回调事件会带上该id值表示该选项，最长支持128字节，不可重复 |
| option_list.text | String   | 是   | 下拉式的选择器选项的文案，建议不超过10个字                   |

#### Button结构体

按钮列表

| 参数  | 类型   | 必须 | 说明                                                         |
| ----- | ------ | ---- | ------------------------------------------------------------ |
| text  | String | 是   | 按钮文案，建议不超过10个字                                   |
| style | Int    | 否   | 按钮样式，目前可填1~4，不填或错填默认1， 按钮样式如下所示： ![img](https://wework.qpic.cn/wwpic/805842_iKxTyYPiRBamTcX_1628665323/0) |
| key   | String | 是   | 按钮key值，用户点击后，会产生回调事件将本参数作为event_key返回，最长支持1024字节，不可重复 |

 

#### Checkbox结构体

选择题样式

| 参数                   | 类型     | 必须 | 说明                                                         |
| ---------------------- | -------- | ---- | ------------------------------------------------------------ |
| question_key           | String   | 是   | 选择题key值，用户提交选项后，会产生回调事件，回调事件会带上该key值表示该题，最长支持1024字节 |
| disable                | Bool     | 否   | 投票选择框的是否不可选，false为可选，true为不可选。仅在更新模版卡片的时候该字段有效 |
| mode                   | Int      | 否   | 选择题模式，单选：0，多选：1，不填默认0                      |
| option_list            | Object[] | 是   | 选项list，选项个数不超过 20 个，最少1个                      |
| option_list.id         | String   | 是   | 选项id，用户提交选项后，会产生回调事件，回调事件会带上该id值表示该选项，最长支持128字节，不可重复 |
| option_list.text       | String   | 是   | 选项文案描述，建议不超过11个字                               |
| option_list.is_checked | Bool     | 否   | 该选项是否要默认选中。                                       |

## 回调和回复的加解密方案

最后更新：2025/07/23

### 验证URL有效性

当点击“保存”提交开发配置信息时，企业微信会发送一条验证消息到填写的URL，发送方法为**GET**。
智能机器人的接收消息服务器接收到验证请求后，需要作出正确的响应才能通过URL验证。

> 获取请求参数时需要做Urldecode处理，否则会验证不成功

假设接收消息地址设置为：https://api.3dept.com/，企业微信将向该地址发送如下验证请求：

**请求方式：GET**

**请求地址**：https://api.3dept.com/?msg_signature=ASDFQWEXZCVAQFASDFASDFSS&timestamp=13500001234&nonce=123412323&echostr=ENCRYPT_STR
**参数说明**

| 参数          | 必须 | 说明                                                         |
| ------------- | ---- | ------------------------------------------------------------ |
| msg_signature | 是   | 企业微信加密签名，msg_signature结合了开发者填写的token、请求中的timestamp、nonce参数、加密的消息体 |
| timestamp     | 是   | 时间戳                                                       |
| nonce         | 是   | 随机数，两个小时内保证不重复                                 |
| echostr       | 是   | 加密的字符串。需要[解密得到消息内容明文](https://developer.work.weixin.qq.com/document/path/101033#12976/密文解密得到msg的过程)，解密后有random、msg_len、msg三个字段，其中msg即为消息内容明文 |

智能机器人后台收到请求后，需要做如下操作：

1. 对收到的请求做Urldecode处理
2. 通过参数msg_signature[对请求进行校验](https://developer.work.weixin.qq.com/document/path/101033#12976/消息体签名校验)，确认调用者的合法性。
3. [解密echostr](https://developer.work.weixin.qq.com/document/path/101033#12976/密文解密得到msg的过程)参数得到消息内容(即msg字段)
4. 在1秒内响应GET请求，响应内容为上一步得到的明文消息内容(不能加引号，不能带bom头，不能带换行符)

以上2~3步骤可以直接使用[验证URL函数](https://developer.work.weixin.qq.com/document/path/101033#12976/验证URL函数)一步到位。
之后接入验证生效，接收消息开启成功。

> 企业内部智能机器人场景中，ReceiveId为""

### 接收回调解密

**智能机器人的回调格式为json**，参考**接收数据格式**说明。开发者可以直接使用企业微信为应用提供的[加解密库](https://developer.work.weixin.qq.com/document/path/101033#12976)（目前已有c++/python/php/java/c#等语言版本）解密encrypt字段，获取事件明文json报文。**需要注意的是，加解密库要求传 receiveid 参数，企业自建智能机器人的使用场景里，receiveid直接传空字符串即可；。**

***\*加密数据格式 ：\****

```javascript
{
	"encrypt": "msg_encrypt"
}
```

***\*参数说明\****

| 参数    | 是否必填 | 说明                     |
| ------- | -------- | ------------------------ |
| encrypt | 是       | 消息结构体加密后的字符串 |



### 加密与被动回复

开发者解密数据得到用户消息内容后，可以选择直接回复空包，也可以在响应本次请求的时候直接回复消息。回复的消息需要先按[明文协议](https://developer.work.weixin.qq.com/document/path/101033#59068)构造json数据包，然后对明文消息进行加密，然后填充到下述协议中的encrypt字段中，之后再回复最终的密文json数据包。加密过程参见“[明文msg的加密过程](https://developer.work.weixin.qq.com/document/path/101033#12976/明文msg的加密过程)”。

**加密数据格式:**

```javascript
{
	"encrypt": "msg_encrypt",
	"msgsignature": "msg_signaturet",
	"timestamp": 1641002400,
	"nonce": "nonce"
}
```

**参数说明：**

| 参数         | 是否必须 | 说明                           |
| ------------ | -------- | ------------------------------ |
| encrypt      | 是       | 加密后的消息内容               |
| msgsignature | 是       | 消息签名                       |
| timestamp    | 是       | 时间戳，要求为秒级别的         |
| nonce        | 是       | 随机数，需要用回调url中的nonce |

## 主动回复消息

最后更新：2025/12/12

### 概述

当用户与智能机器人进行交互时，企业微信会将相关的[交互事件](https://developer.work.weixin.qq.com/document/path/101138#57141)回调到开发者设置的回调URL，回调中会返回一个 response_url 。开发者可根据事件类型先做出相应的[响应](https://developer.work.weixin.qq.com/document/path/101138#59068/回复用户消息)，待处理完业务逻辑后，使用response_url主动调用接口回复消息，实现丰富的自定义功能。

目前有以下场景回调会返回 response_url ，支持主动回复消息：

1. 用户向智能机器人发送消息，[前往查看](https://developer.work.weixin.qq.com/document/path/101138#57141)
2. 用户点击模板卡片相关按钮等，[前往查看](https://developer.work.weixin.qq.com/document/path/101138#59058/模板卡片事件)

> 请注意，每个 response_url 用户可以调用接口一次， 该 response_url 有效期为 1 个小时，超过有效期将无法使用。

交互流程如下图所示：
![img](https://wework.qpic.cn/wwpic3az/57968_09QYXvdcQX-TUnO_1761746792/0)

 

### 如何主动回复消息

开发者获取到response_url后，可以按以下说明向这个地址发起HTTP POST 请求，即可对相应的回调进行主动回复。下面举个简单的例子.

- 假设 response_url 是：https://qyapi.weixin.qq.com/cgi-bin/aibot/response?response_code=RESPONSE_CODE。以下是用curl工具往群组推送文本消息的示例（注意要将url替换成对应的response_url，content必须是utf8编码）：

```javascript
curl 'https://qyapi.weixin.qq.com/cgi-bin/aibot/response?response_code=RESPONSE_CODE' \
   -H 'Content-Type: application/json' \
   -d '
{
    "msgtype": "markdown",
    "markdown": {
        "content": "hello world"
    }
}'
```

### 消息类型及数据格式

#### markdown消息

```javascript
{
    "msgtype": "markdown",
    "markdown": {
        "content": "# 一、标题\n## 二级标题\n### 三级标题\n# 二、字体\n*斜体*\n\n**加粗**\n# 三、列表 \n- 无序列表 1 \n- 无序列表 2\n  - 无序列表 2.1\n  - 无序列表 2.2\n1. 有序列表 1\n2. 有序列表 2\n# 四、引用\n> 一级引用\n>>二级引用\n>>>三级引用\n# 五、链接\n[这是一个链接](https:work.weixin.qq.com\/api\/doc)\n![](https://res.mail.qq.com/node/ww/wwopenmng/images/independent/doc/test_pic_msg1.png)\n# 六、分割线\n\n---\n# 七、代码\n`这是行内代码`\n```\n这是独立代码块\n```\n\n# 八、表格\n| 姓名 | 文化衫尺寸 | 收货地址 |\n| :----- | :----: | -------: |\n| 张三 | S | 广州 |\n| 李四 | L | 深圳 |\n",
        "feedback": {
            "id": "FEEDBACKID"
        }
    }
}
```

**参数说明：**

| 参数                 | 类型   | 是否必填 | 说明                                                         |
| -------------------- | ------ | -------- | ------------------------------------------------------------ |
| msgtype              | String | 是       | 消息类型，此时固定为：markdown                               |
| markdown.content     | String | 是       | 消息内容，最长不超过20480个字节，必须是utf8编码。            |
| markdown.feedback.id | String | 否       | 若字段不为空值，回复的消息被用户反馈时候会触发[回调事件](https://developer.work.weixin.qq.com/document/path/101138#59058/用户反馈事件)。有效长度为 256 字节以内，必须是 utf-8 编码。 |

> 回复内容content字段支持常见的[markdown格式](https://developer.work.weixin.qq.com/document/path/101138#14404/markdown-v2类型)

#### 模板卡片消息
    "template_card": {
        "feedback": {
            "id": "FEEDBACKID"
        }
    }
}
```

**参数说明：**

| 参数                      | 类型   | 是否必填 | 说明                                                         |
| ------------------------- | ------ | -------- | ------------------------------------------------------------ |
| msgtype                   | String | 是       | 消息类型，此时固定为：template_card。当且仅当回调的会话类型为单聊的时候，支持该主动回复类型 |
| template_card             | Object | 是       | 模板卡片结构体，参考[模板卡片类型](https://developer.work.weixin.qq.com/document/path/101138#59098)中类型说明 |
| template_card.feedback.id | String | 否       | 特殊的该回复场景支持设置反馈信息。若字段不为空值，回复的消息被用户反馈时候会触发[回调事件](https://developer.work.weixin.qq.com/document/path/101138#59058/用户反馈事件)。有效长度为 256 字节以内，必须是 utf-8 编码。 |

> 特殊的，因主动回复可能跟用户交互触发的回调间隔较长时间，群聊中智能机器人主动回复消息的时候会引用触发回调的用户消息/被点击模板卡片消息。模板卡片消息不支持引用，会默认生成一条空消息进行引用。