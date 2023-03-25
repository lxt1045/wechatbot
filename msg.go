package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/eatmoreapple/openwechat"
)

func replyText(contextMgr ContextMgr, content string, msg *openwechat.Message) (retCtxMgr ContextMgr) {
	log.Printf("Received Msg : %v", msg.Content)
	// content := msg.Content
	// reply, err := replyByGPT(userFromName,content)
	reply, retCtxMgr, err := replyByGPT3_5(contextMgr, content)
	if err != nil {
		log.Println(err)
		// 如果文字超过4000个字会回错，截取前4000个文字进行回复
		if len(reply) > 4000 {
			reply = postReply(reply)
			reply = reply[:4000]
			_, err = msg.ReplyText(reply)
			if err != nil {
				log.Println("回复出错：", err.Error())
				return
			}
			log.Println("reply:", reply)
			return
		}

		msg.ReplyText(fmt.Sprintf("bot error: %s", err.Error()))
		log.Println("reply:", reply)
		return
	}
	reply = postReply(reply)
	msg.ReplyText(reply)
	log.Println("reply:", reply)
	// msg.ReplyText("以上回答由 ChatGPT 提供")
	return
}

func replyFriendAdd(msg *openwechat.Message) {
	log.Printf("Received Msg : %v", msg.Content)
	fm, err := msg.FriendAddMessageContent()
	if err != nil {
		log.Printf("err: %v", err)
		return
	}
	mFriend[fm.FromUserName] = &openwechat.Friend{
		User: &openwechat.User{
			NickName: fm.FromNickName,
		},
	}

	log.Printf("fm: %+v", fm)
	msg.Agree("OK!")
}

func GetTmpFileFromURL(url string, f func(*os.File) error) (err error) {
	// Create our Temp File
	tmpFile, err := ioutil.TempFile(os.TempDir(), "openai.*.png")
	if err != nil {
		log.Printf("Cannot create temporary file: %v", err)
		return
	}
	// defer os.Remove(tmpFile.Name())
	fmt.Println("Created File: " + tmpFile.Name())
	fmt.Println("url: " + url)

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
	return f(tmpFile)
}
func replyCreateImage(contextMgr ContextMgr, msg *openwechat.Message, n int, size, content string) (retCtxMgr ContextMgr) {
	retCtxMgr = contextMgr

	log.Printf("Received Msg : %v", msg.Content)

	var urls []string
	var err error
	if retCtxMgr.Mode == ImgMode && retCtxMgr.LastImg == "" {
		urls, err = replyImage(n, size, content)
	} else {
		var bs []byte
		bs, err = ioutil.ReadFile(retCtxMgr.LastImg)
		if err != nil {
			log.Printf("[replyImage]: %v", err)
			return
		}
		file2 := `/Users/bytedance/go/src/github.com/lxt1045/wechatbot/image.png`
		bsmask := []byte{}
		bsmask, err = ioutil.ReadFile(file2)
		if err != nil {
			log.Printf("[replyImage]: %v", err)
			return
		}
		urls, err = editImage2(bs, bsmask, n, size, content)
	}
	if err != nil {
		log.Printf("[replyImage]: %v", err)
		return
	}
	for _, u := range urls {
		err := GetTmpFileFromURL(u, func(file *os.File) error {
			_, err := msg.ReplyImage(file)
			if err != nil {
				log.Printf("[ReplyImage]: %v", err)
				return err
			}
			if retCtxMgr.Mode == ImgMode {
				if retCtxMgr.LastImg != "" {
					os.Remove(retCtxMgr.LastImg)
				}
				retCtxMgr.LastImg = file.Name()
			} else {
				os.Remove(file.Name())
			}
			return nil
		})
		if err != nil {
			log.Printf("[GetTmpFileFromURL]: %v", err)
			continue
		}
	}
	return
}
