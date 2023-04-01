package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/eatmoreapple/openwechat"
)

func handler(msg *openwechat.Message) {
	if msg.IsSendBySelf() {
		return
	}

	sender, err := msg.Sender()
	if err != nil {
		log.Println(err)
		return
	}

	preMsg := ""
	userFromName := msg.FromUserName
	if msg.IsSendByGroup() {
		group := openwechat.Group{User: sender}
		log.Println(group.NickName)

		g, ok := mGroup[msg.FromUserName]
		if !ok || (g.NickName != "机器人聊天群" && g.NickName != "疯狂动物城") {
			log.Printf("msg.Content:%+v", msg.Content)
			return
		}
		log.Printf("from group:%+v", g.NickName)
		preMsg = "@机器人"
	} else {
		f, ok := mFriend[msg.FromUserName]
		if ok {
			log.Printf("from friend:%+v", f.NickName)
		}
	}
	// receiver, err := msg.Receiver()
	// if err != nil || receiver.NickName != "小影" {
	// 	log.Printf("receiver : %v", receiver.NickName)
	// 	return
	// }
	switch {
	case msg.IsText():
		if !strings.HasPrefix(msg.Content, preMsg) {
			return
		}
		content := msg.Content[len(preMsg):]
		content = strings.TrimSpace(content)
		contextMgr := GetContextMgr(userFromName)
		if content == "图片模式" {
			contextMgr.Mode = ImgMode
			contextMgr.LastImg = ""
			SetContextMgr(userFromName, contextMgr)
			msg.ReplyText(`好的，开始图片模式。
你可以发"新图片"跳出当前图片的编辑。
也可以发"聊天模式"切换到文本聊天模式。
也可以发"文本编辑"切换到文本编辑模式。`)
			return
		} else if content == "聊天模式" {
			contextMgr.Mode = ChatMode
			if len(contextMgr.contextList) > 0 {
				contextMgr.contextList = contextMgr.contextList[:0]
			}
			SetContextMgr(userFromName, contextMgr)
			msg.ReplyText(`好的，开始文本聊天模式
也可以发"图片模式"切换到图片模式。
也可以发"文本编辑"切换到文本编辑模式。`)
			return
		} else if content == "新图片" {
			if contextMgr.Mode == ImgMode {
				contextMgr.LastImg = ""
				SetContextMgr(userFromName, contextMgr)
				msg.ReplyText(`好的，准备好创建新图片，请说出你对图片的描述。
也可以发"聊天模式"切换到文本聊天模式。
也可以发"文本编辑"切换到文本编辑模式。`)
				return
			}
		} else if content == "文本编辑" {
			contextMgr.Mode = TextEdit
			SetContextMgr(userFromName, contextMgr)
			msg.ReplyText(`好的，开始文本编辑模式。
也可以发"图片模式"切换到图片模式。
也可以发"聊天模式"切换到文本聊天模式。`)
			return
		}

		switch contextMgr.Mode {
		case ChatMode:
			go func() {
				if len(contextMgr.LastAudio) > 0 {
					replyAudio(contextMgr, msg, content)
					return
				}
				contextMgr = replyText(contextMgr, content, msg)
				SetContextMgr(userFromName, contextMgr)
			}()
		case TextEdit:
			go func() {
				contextMgr = editText(contextMgr, content, msg)
				SetContextMgr(userFromName, contextMgr)
			}()
		case ImgMode:
			go func() {
				contextMgr = replyCreateImage(contextMgr, msg, 1, "1024x1024", content)
				SetContextMgr(userFromName, contextMgr)
			}()
			return
		}

	case msg.IsFriendAdd():
		replyFriendAdd(msg)
	case msg.MsgType == openwechat.MsgTypeApp &&
		msg.AppMsgType == openwechat.AppMsgTypeTransfers:
		log.Printf("红包:%+v", msg)
	case msg.MsgType == openwechat.MsgTypeSys &&
		msg.AppMsgType == openwechat.AppMsgTypeRedEnvelopes:
		log.Printf("红包:%+v", msg)
	case msg.MsgType == openwechat.MsgTypeVoice:
		resp, err := msg.GetVoice()
		if err != nil {
			log.Printf("err:%v", err)
			return
		}
		defer resp.Body.Close()
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("err:%v", err)
			return
		}
		contextMgr := GetContextMgr(userFromName)
		contextMgr.LastAudio = bs
		SetContextMgr(userFromName, contextMgr)
		// err = os.WriteFile("test.mp3", bs, 0666)
		// if err != nil {
		// 	log.Printf("err:%v", err)
		// 	return
		// }
	default:
		log.Printf("msg:%+v", msg)

	}
	return
}

var imgType = `请求格式出错，正确格式是：
	img:n:size:content 或 img:::content
	img：固定
	n：1、2、3 ... 8
	size：256、512、1024
	content：图片描述
`
