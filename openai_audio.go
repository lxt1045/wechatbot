package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	// log "github.com/sirupsen/logrus"
)

/*
接口使用:
https://platform.openai.com/docs/api-reference/images/create
模型介绍：
https://platform.openai.com/docs/models/overview
*/

type AudioResp struct {
	Created int64  `json:"created"`
	Text    string `json:"text"`
	Err     struct {
		Code    *int        `json:"code"`
		Message string      `json:"message"`
		Param   interface{} `json:"param"`
		Type    string      `json:"type"`
	} `json:"error"`
}

func openaiReplyAudio(audio []byte, content string) (text string, err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// fw, err := w.CreateFormField("image")
	// if err != nil {
	// 	log.Println(err)
	// }
	// _, err = fw.Write(img)
	// if err != nil {
	// 	log.Println(err)
	// }
	fw, err := w.CreateFormFile("file", "test.mp3")
	if err == nil {
		_, err = fw.Write(audio)
	}

	if err == nil {
		fw, err = w.CreateFormField("model")
	}
	if err == nil {
		_, err = fw.Write([]byte("whisper-1"))
	}

	if content != "" {
		if err == nil {
			fw, err = w.CreateFormField("prompt")
		}
		if err == nil {
			_, err = fw.Write([]byte(content))
		}
	}

	if err == nil {
		fw, err = w.CreateFormField("language")
	}
	if err == nil {
		// ISO-639-1 https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
		languages := []string{"zh", "en", "ja"}
		_, err = fw.Write([]byte(languages[0]))
	}
	if err != nil {
		log.Println(err)
		return
	}

	w.Close()

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", &b)
	if err != nil {
		log.Println(err)
	}
	req.SetBasicAuth("api", apiKey)

	req.Header.Add("Content-Type", w.FormDataContentType())
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.Header.Set("Authorization", fmt.Sprintf("Bearer "+apiKey))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	gptResp := &AudioResp{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResp)
	if err != nil {
		log.Println(err)
		return
	}

	return gptResp.Text, nil
}
