package ding

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

//https://oapi.dingtalk.com/robot/send?access_token=a815d6cc1b875da762b8227f4dd75d455c1106c2cf720977926dcc950066d1b9
var (
	dingURL       = "https://oapi.dingtalk.com"
	dingActionURL = "/robot/send"
	//dingAccessToken = "4eb844844f228936cea289970c21aa8737bc1806f20479e7d6a41c74f9356d70"
	dingAccessToken = "a815d6cc1b875da762b8227f4dd75d455c1106c2cf720977926dcc950066d1b9"
)

func SendDing(content string) error {
	//请求钉钉webhook
	if content != "" {
		format := `
        {
            "msgtype": "text",
            "text": {
                "content": "%s"
            },
			"at":{
				"isAtAll":"true"
			}
        }`
		body := fmt.Sprintf(format, content)
		jsonValue := []byte(body)
		//发送消息到钉钉群使用webhook
		s := dingURL + dingActionURL + "?access_token=" + dingAccessToken
		resp, err := http.Post(s, "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			return err
		}
		log.Println(resp)
	}
	return nil
}

//异步发送消息
func SendSyncDing(app, messageType string, data map[string]interface{}) {
	urls := dingURL + dingActionURL + "?access_token=" + dingAccessToken
	ding := NewDingHook(urls, app, allLevel[1])
	msg := logrus.Entry{
		Data:    data,
		Time:    time.Now(),
		Level:   logrus.InfoLevel,
		Message: "发送请求SUCCESS",
	}
	switch messageType {
	case "markdown":
		_ = ding.FireMarkDown(&msg)
		break
	default:
		_ = ding.FireText(&msg)
	}
}
