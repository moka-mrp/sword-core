package ding

import "testing"

func TestSendDing(t *testing.T) {
	content := "这个消息发送成功了吗？"
	SendDing(content)
}
