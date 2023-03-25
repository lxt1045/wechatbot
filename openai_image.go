package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	// log "github.com/sirupsen/logrus"
)

/*
接口使用:
https://platform.openai.com/docs/api-reference/images/create
模型介绍：
https://platform.openai.com/docs/models/overview
*/

type ImageReq struct {
	Prompt         string `json:"prompt"`                    // A text description of the desired image(s). The maximum length is 1000 characters.
	N              int    `json:"n,omitempty"`               // Optional, Defaults to 1; The number of images to generate. Must be between 1 and 10.
	Size           string `json:"size,omitempty"`            // Optional, Defaults to 1024x1024; The size of the generated images. Must be one of 256x256, 512x512, or 1024x1024
	ResponseFormat string `json:"response_format,omitempty"` // Optional, Defaults to url; The format in which the generated images are returned. Must be one of url or b64_json
	Image          string `json:"image,omitempty"`
}

type ImageResp struct {
	Created int64 `json:"created"`
	Data    []struct {
		URL string `json:"url"`
	} `json:"data"`
	Err struct {
		Code    *int        `json:"code"`
		Message string      `json:"message"`
		Param   interface{} `json:"param"`
		Type    string      `json:"type"`
	} `json:"error"`
}

func replyImage(n int, size, content string) (urls []string, err error) {
	requestBody := ImageReq{
		Prompt: content,
		N:      n,
		Size:   size,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("request openai json string : %v", string(requestData))
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewBuffer(requestData))
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer "+apiKey))
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	gptResp := &ImageResp{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResp)
	if err != nil {
		log.Println(err)
		return
	}

	for _, u := range gptResp.Data {
		urls = append(urls, u.URL)
	}

	return urls, nil
}

func editImage(img string, n int, size, content string) (urls []string, err error) {
	requestBody := ImageReq{
		Prompt: content,
		N:      n,
		Size:   size,
		Image:  img,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("request openai json string : %v", string(requestData))
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/edits", bytes.NewBuffer(requestData))
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer "+apiKey))
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	gptResp := &ImageResp{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResp)
	if err != nil {
		log.Println(err)
		return
	}

	for _, u := range gptResp.Data {
		urls = append(urls, u.URL)
	}

	return urls, nil
}

func editImage2(img, mask []byte, n int, size, content string) (urls []string, err error) {
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
	fw, err := w.CreateFormFile("image", "test.png")
	if err != nil {
		log.Println(err)
	}
	_, err = fw.Write(img)
	if err != nil {
		log.Println(err)
	}
	if len(mask) > 0 {
		fw, err = w.CreateFormFile("mask", "mask.png")
		if err != nil {
			log.Println(err)
		}
		_, err = fw.Write(mask)
		if err != nil {
			log.Println(err)
		}
	}

	fw, err = w.CreateFormField("prompt")
	if err != nil {
		log.Println(err)
	}
	_, err = fw.Write([]byte(content))
	if err != nil {
		log.Println(err)
	}

	fw, err = w.CreateFormField("n")
	if err != nil {
		log.Println(err)
	}
	_, err = fw.Write([]byte(strconv.Itoa(n)))
	if err != nil {
		log.Println(err)
	}
	fw, err = w.CreateFormField("size")
	if err != nil {
		log.Println(err)
	}
	_, err = fw.Write([]byte("1024x1024"))
	if err != nil {
		log.Println(err)
	}

	w.Close()

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/edits", &b)
	if err != nil {
		log.Println(err)
	}
	req.SetBasicAuth("api", "key-3ax6xnjp29jd6fds4gc373sgvjxteol0")

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

	gptResp := &ImageResp{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResp)
	if err != nil {
		log.Println(err)
		return
	}

	for _, u := range gptResp.Data {
		urls = append(urls, u.URL)
	}

	return urls, nil
}

func PostTmpFileFromURL(url string, f func(*os.File) error) (err error) {
	// Create our Temp File
	tmpFile, err := ioutil.TempFile(os.TempDir(), "openai.*.png")
	if err != nil {
		log.Printf("Cannot create temporary file: %v", err)
		return
	}
	defer os.Remove(tmpFile.Name())
	fmt.Println("Created File: " + tmpFile.Name())

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Cannot create temporary file: %v", err)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		log.Printf("Cannot create temporary file: %v", err)
		return
	}
	if _, err = tmpFile.Seek(0, 0); err != nil {
		return err
	}

	return nil
}
