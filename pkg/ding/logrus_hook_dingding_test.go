package ding

import (
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestNewDingHook(t *testing.T) {
	urls := dingURL + dingActionURL + "?access_token=" + dingAccessToken
	ding := NewDingHook(urls, "Crawl_api", allLevel[1])
	data := make(map[string]interface{})
	data["name"] = "James"
	data["sex"] = "male"
	data["age"] = "38"
	data["favourite"] = "have launch"
	log := logrus.Entry{
		Data:    data,
		Time:    time.Now(),
		Level:   logrus.InfoLevel,
		Message: "这边好像发送了一个请求",
	}
	err := ding.FireMarkDown(&log)
	//ding.startDingHookQueueJob()
	t.Log(err)
}
