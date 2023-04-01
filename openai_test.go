package main

import (
	"testing"
)

// log "github.com/sirupsen/logrus"

/*
接口使用:
https://platform.openai.com/docs/api-reference/images/create
模型介绍：
https://platform.openai.com/docs/models/overview
*/

func TestEditText(t *testing.T) {
	Init()

	str := `我转头，仔细往那里看，那里的手电暗了，有一个声音叫到“小三爷！”“潘子！”我惊了一下，但是没法靠过去看。对方道：“小三爷，快走。”声音相当微弱。接着我听到一连串的咳嗽声。“你怎么样？”我问道，“你怎么会在这儿？”潘子在黑暗中说道：“说来话长了，小三爷，你有烟吗？”“在这儿你还抽烟，不怕肺烧穿？”我听着潘子的语气，觉得他特别地淡定，忽然起了一种非常不详的预感。“哈哈哈，没关系了。”潘子道，“你看不到我现在是什么样子。”我心中的不详感越来越甚，道：“别磨蹭了，赶快过来，你不过来我就过去扶你。”说着，我用手电去照，隐约能照到他的样子，我就意识到为什么前几次我都看不到他。潘子似乎是卡在了岩层中，我扩大了光圈，一下子就看到，他的身子融在岩层里，成了人影。`
	t.Run("openaiEditText", func(t *testing.T) {
		text, _, err := openaiEditText(
			ContextMgr{
				EditTexTemperature: 1,
			},
			str,
			"改成古龙风格",
		)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("text:%+v", text)
	})
}
