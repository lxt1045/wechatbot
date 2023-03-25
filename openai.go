package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

func replyByGPT(user, msg string) (string, error) {

	type ChatGPTResponseBody struct {
		ID      string                   `json:"id"`
		Object  string                   `json:"object"`
		Created int                      `json:"created"`
		Model   string                   `json:"model"`
		Choices []map[string]interface{} `json:"choices"`
		Usage   map[string]interface{}   `json:"usage"`
	}

	type ChatGPTRequestBody struct {
		Model            string  `json:"model"`
		Prompt           string  `json:"prompt"`
		MaxTokens        int     `json:"max_tokens"`
		Temperature      float32 `json:"temperature"`
		TopP             int     `json:"top_p"`
		FrequencyPenalty int     `json:"frequency_penalty"`
		PresencePenalty  int     `json:"presence_penalty"`
	}

	requestBody := ChatGPTRequestBody{
		Model:            "text-davinci-003",
		Prompt:           msg,
		MaxTokens:        256,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	log.Println(requestBody)
	requestData, err := json.Marshal(requestBody)
	log.Println(string(requestData))
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewBuffer(requestData))
	if err != nil {
		log.Println(err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var reply string
	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			reply = v["text"].(string)
			break
		}
	}
	log.Printf("gpt response text: %s \n", reply)

	return reply, nil
}

var apiKey = ""

var (
	mContextMgr sync.Map
)

func GetContextMgr(user string) (contextMgr ContextMgr) {
	i, ok := mContextMgr.Load(user)
	if ok {
		contextMgr, ok = i.(ContextMgr)
		if ok {
			return
		}
	}
	return
}

func SetContextMgr(user string, contextMgr ContextMgr) {
	mContextMgr.Store(user, contextMgr)
	return
}

// replyByGPT3_5 sendMsg
func replyByGPT3_5(contextMgr ContextMgr, msg string) (reply string, retCtxMgr ContextMgr, err error) {
	retCtxMgr = contextMgr
	var messages []ChatMessage
	messages = append(messages, ChatMessage{
		Role:    "system",
		Content: "You are a helpful assistant.",
	})
	list := retCtxMgr.GetData()
	for i := 0; i < len(list); i++ {
		messages = append(messages, ChatMessage{
			Role:    "user",
			Content: list[i].Request,
		})

		messages = append(messages, ChatMessage{
			Role:    "assistant",
			Content: list[i].Response,
		})
	}

	messages = append(messages, ChatMessage{
		Role:    "user",
		Content: msg,
	})

	requestBody := ChatGPTRequestBody{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("request openai json string : %v", string(requestData))
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestData))
	if err != nil {
		log.Error(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer "+apiKey))
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	gptResponseBody := &ChatGPTResponseBody{}
	log.Debug(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		log.Error(err)
		return
	}

	if len(gptResponseBody.Choices) > 0 {
		for _, v := range gptResponseBody.Choices {
			reply += "\n"
			reply += v.Message.Content
		}

		retCtxMgr.AppendMsg(msg, reply)
	}

	if len(reply) == 0 {
		gptErrorBody := &ChatGPTErrorBody{}
		err = json.Unmarshal(body, gptErrorBody)
		if err != nil {
			log.Error(err)
			return
		}

		reply += "Error: "
		reply += gptErrorBody.Error["message"].(string)
	}

	return reply, retCtxMgr, nil
}
