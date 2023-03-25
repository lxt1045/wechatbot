package main

import (
	"log"
	"strconv"
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
			msg.ReplyText("好的，开始图片模式。\n你可以发\"新图片\"跳出当前图片的编辑。\n也可以发\"文本模式\"切换到文本聊天模式。")
			return
		} else if content == "文本模式" {
			contextMgr.Mode = TextMode
			if len(contextMgr.contextList) > 0 {
				contextMgr.contextList = contextMgr.contextList[:0]
			}
			SetContextMgr(userFromName, contextMgr)
			msg.ReplyText("好的，开始文本聊天模式。\n你可以发\"图片模式\"切换到图片模式。")
			return
		} else if content == "新图片" {
			if contextMgr.Mode == ImgMode {
				contextMgr.LastImg = ""
				SetContextMgr(userFromName, contextMgr)
				msg.ReplyText("好的，准备好创建新图片，请说出你对图片的描述。\n你可以发\"文本模式\"切换到文本聊天模式。")
				return
			}
		}

		switch contextMgr.Mode {
		case TextMode:
		case ImgMode:
			go func() {
				contextMgr = replyCreateImage(contextMgr, msg, 1, "1024x1024", content)
				SetContextMgr(userFromName, contextMgr)
			}()
			return
		}

		imgPre := "img"
		if strings.HasPrefix(content, imgPre) {
			// img:n:size:content
			strs := strings.SplitN(content, ":", 4)
			if len(strs) != 4 {
				strs = strings.SplitN(content, "：", 4)
			}
			if len(strs) != 4 {
				_, err = msg.ReplyText(imgType)
				if err != nil {
					log.Println("回复出错：", err.Error())
					return
				}
				break
			}

			x := strings.TrimSpace(strs[1])
			n := 3
			if x != "" {
				n, err = strconv.Atoi(x)
				if err != nil {
					_, err = msg.ReplyText(imgType)
					if err != nil {
						log.Println("回复出错：", err.Error())
						return
					}
					return
				}
			}
			if n <= 0 {
				n = 1
			} else if n > 8 {
				n = 8
			}

			y := strings.TrimSpace(strs[2])
			size := 0
			if y != "" {
				size, err = strconv.Atoi(y)
				if err != nil {
					return
				}
			}
			sizeStr := "1024x1024"
			if size == 256 || size == 1024 {
				s := strconv.Itoa(size)
				sizeStr = s + "x" + s
			}

			content := strs[3]
			go func() {
				contextMgr = replyCreateImage(contextMgr, msg, n, sizeStr, content)
				SetContextMgr(userFromName, contextMgr)
			}()
			break
		}
		if strings.HasSuffix(content, "的图片") ||
			strings.HasSuffix(content, "的照片") ||
			strings.HasSuffix(content, "的画") {
			go func() {
				contextMgr = replyCreateImage(contextMgr, msg, 3, "1024x1024", content)
				SetContextMgr(userFromName, contextMgr)
			}()
			break
		}
		go func() {
			contextMgr = replyText(contextMgr, content, msg)
			SetContextMgr(userFromName, contextMgr)
		}()
	case false && msg.IsFriendAdd():
		replyFriendAdd(msg)
	case msg.MsgType == openwechat.MsgTypeApp &&
		msg.AppMsgType == openwechat.AppMsgTypeTransfers:
		log.Printf("红包:%+v", msg)
	case msg.MsgType == openwechat.MsgTypeSys &&
		msg.AppMsgType == openwechat.AppMsgTypeRedEnvelopes:
		log.Printf("红包:%+v", msg)

	default:
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
