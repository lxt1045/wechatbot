package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/eatmoreapple/openwechat"
)

func init() {
	log.SetFlags(log.Flags() | log.Lmicroseconds | log.Lshortfile) //log.Llongfile
}

var mGroup = make(map[string]*openwechat.Group)
var mFriend = make(map[string]*openwechat.Friend)

func Init() {
	cfg := struct {
		APIKey string `json:"api_key"`
	}{}
	bs, err := os.ReadFile(`./config.json`)
	if err != nil {
		log.Panicln(err)
	}
	err = json.Unmarshal(bs, &cfg)
	if err != nil || cfg.APIKey == "" {
		log.Panicln(err)
	}
	apiKey = cfg.APIKey

}

func main() {
	Init()

	// bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 注册消息处理函数
	bot.MessageHandler = handler

	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	/*
		reloadStorage := openwechat.NewJsonFileHotReloadStorage("token.json")
		err := bot.HotLogin(reloadStorage)
		if err != nil {
			err := os.Remove("token.json")
			if err != nil {
				return
			}
			reloadStorage = openwechat.NewJsonFileHotReloadStorage("token.json")
			err = bot.HotLogin(reloadStorage)
			if err != nil {
				return
			}
		}
	*/

	loginFile := `./.login.json`
	hot, err := os.ReadFile(loginFile)
	if err != nil {
		log.Println(err)
	}
	buf := bytes.NewBuffer(hot)
	err = bot.HotLogin(buf)
	if err != nil {
		buf.Reset()
		err = bot.HotLogin(buf)
		if err != nil {
			log.Println(err)
			return
		}
		buf.Reset()
		if err := bot.DumpHotReloadStorage(); err != nil {
			log.Println(err)
			return
		}
	}
	if len(hot) == 0 {
		buf.Reset()
		if err := bot.DumpHotReloadStorage(); err != nil {
			log.Println(err)
			return
		}
	}
	// log.Printf("hot:%s", buf.String())
	err = os.WriteFile(loginFile, buf.Bytes(), 0666)
	if err != nil {
		log.Println(err)
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		log.Println(err)
		return
	}

	// 获取所有的好友
	friends, err := self.Friends()
	log.Println(friends, err)
	for _, f := range friends {
		mFriend[f.UserName] = f
	}

	// 获取所有的群组
	groups, err := self.Groups()
	log.Println(groups, err)
	for _, g := range groups {
		mGroup[g.UserName] = g
	}

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
