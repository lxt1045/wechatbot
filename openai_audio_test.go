package main

import (
	"io/ioutil"
	"testing"
)

// log "github.com/sirupsen/logrus"

/*
接口使用:
https://platform.openai.com/docs/api-reference/images/create
模型介绍：
https://platform.openai.com/docs/models/overview
*/
func TestAudio(t *testing.T) {
	Init()

	t.Run("openaiReplyAudio", func(t *testing.T) {

		file := `./test.mp3`

		bs, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		text, err := openaiReplyAudio(bs, "可爱风格")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("text:%+v", text)
	})
}
