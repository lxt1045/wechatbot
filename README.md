# wechatbot
wechat openAI bot

# 前言
当前代码比较潦草，还没时间整理，还请多担待

# 启动步骤
1. 把 apikey 写入新建的 配置文件 config.json 中。
```json
{
    "api_key": "sk-xxxx"
}
```
2. 启动后会弹出微信登录二维码，扫码即可登录。

# 开通 openai 账号
如果想购买外网验证码服务注册，可以看 [参考文档](https://blog.laoda.de/archives/play-with-chatgpt#2.-%E2%9C%88%EF%B8%8F%E6%B3%A8%E5%86%8C) 。


# openai 接口
[官方接口文档](https://platform.openai.com/docs/api-reference/edits)
跟着文档写接口就好了，这个库实现了大部分 openai 接口：
1. chat: 带上下文聊天，openai.go
2. edits: 文本编辑，openai.go
3. images: 根据文本画图，改图，openai_image.go
4. audio: 语音翻译，openai_audio.go

# 关于微信机器人
使用了 [eatmoreapple/openwecha](https://github.com/eatmoreapple/openwechat)库。

实现了以下功能：
1. “聊天模式”： 基于聊天对象实现聊天上下文，上下文是 10 分钟内的最多 10 条信息，对应 openai 的 chat 接口。
1. “文本编辑”： 根据描述修改文本内容，对应 openai 的 edits 接口。
1. “图片模式”： 实现通过文字实现画图、改图功能，对应 openai 的 images 接口。
1. 语音翻译：   “聊天模式”下发送语音后，再发送一段描述，返回语音翻译，对应 openai 的 audio 接口。


# 交流学习
![扫码加微信好友](https://github.com/lxt1045/wechatbot/blob/main/resource/Wechat-lxt.png "微信")
