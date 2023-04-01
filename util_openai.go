package main

import (
	"strings"
	"time"
)

func postReply(reply string) string {
	reply = strings.TrimLeftFunc(reply, func(r rune) bool {
		switch r {
		case '?', '？':
			return true
		}
		return false
	})
	// 微信不支持markdown格式，所以把反引号直接去掉
	reply = strings.Replace(reply, "`", " ", -1)

	return reply
}

// ChatGPTResponseBody 响应体
type ChatGPTResponseBody struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int              `json:"created"`
	Model   string           `json:"model"`
	Choices []ResponseChoice `json:"choices"`
	Usage   ResponseUsage    `json:"usage"`
}

// ChatGPTRequestBody 请求体
type ChatGPTRequestBody struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ResponseUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatGPTErrorBody struct {
	Error map[string]interface{} `json:"error"`
}

type Context struct {
	Request  string
	Response string
	Time     int64
}

type ModeType int

const (
	ChatMode ModeType = 0
	ImgMode  ModeType = 1
	TextEdit ModeType = 2
)

type ContextMgr struct {
	contextList        []*Context // 聊天模式下，需保存上下文
	Mode               ModeType   // 当前聊天模式
	LastImg            string     // 图片编辑模式下，需要保存上一张图片
	EditText           string     // 文本编辑模式下，需要保存待编辑文本
	EditTexTemperature float32    // 文本编辑模式下: 介于 0 和 2 之间。较高的值（如 0.8）将使输出更加随机，而较低的值（如 0.2）将使输出更加集中和确定
	LastAudio          []byte     // 翻译语音时，需要先保存语音，在发送一次文本触发翻译
}

func (m *ContextMgr) Init() {
	m.contextList = make([]*Context, 10)
}

func (m *ContextMgr) checkExpire() {
	timeNow := time.Now().Unix()
	if len(m.contextList) > 10 {
		m.contextList = m.contextList[len(m.contextList)-10:]
	}
	if len(m.contextList) > 0 {
		startPos := len(m.contextList) - 1
		for i := 0; i < len(m.contextList); i++ {
			if timeNow-m.contextList[i].Time < 10*60 {
				startPos = i
				break
			}
		}

		m.contextList = m.contextList[startPos:]
	}
}

func (m *ContextMgr) AppendMsg(request string, response string) {
	m.checkExpire()
	context := &Context{Request: request, Response: response, Time: time.Now().Unix()}
	m.contextList = append(m.contextList, context)
}

func (m *ContextMgr) GetData() []*Context {
	m.checkExpire()
	return m.contextList
}
