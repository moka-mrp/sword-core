package ding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

var allLevel = []logrus.Level{
	logrus.DebugLevel,
	logrus.InfoLevel,
	logrus.WarnLevel,
	logrus.ErrorLevel,
	logrus.FatalLevel,
	logrus.PanicLevel,
}

type dingHook struct {
	apiUrl     string
	levels     []logrus.Level
	appName    string
	jsonBodies chan []byte
	closeChan  chan bool
}

func NewDingHook(url, app string, thresholdLevel logrus.Level) *dingHook {
	temp := []logrus.Level{}
	for _, v := range allLevel {
		if v <= thresholdLevel {
			temp = append(temp, v)
		}
	}
	hook := &dingHook{apiUrl: url, levels: temp, appName: app}
	hook.jsonBodies = make(chan []byte)
	hook.closeChan = make(chan bool)
	//开启chan 队列 执行post dingding hook api
	go hook.startDingHookQueueJob()
	//todo 完美ding logrus ding ding hook
	return hook
}

// 启用去监听
func (dh *dingHook) startDingHookQueueJob() {
	for {
		select {
		case <-dh.closeChan:
			return
		case bs := <-dh.jsonBodies:
			res, err := http.Post(dh.apiUrl, "application/json", bytes.NewBuffer(bs))
			if err != nil {
				log.Println(err)
			}
			if res != nil && res.StatusCode != 200 {
				log.Println("dingHook go chan http post error", res.StatusCode)
			}
		}
	}

}

// 设置日志级别
func (dh *dingHook) Levels() []logrus.Level {
	return dh.levels
}

//这个异步有可能导致 最后一条消息丢失,main goroutine 提前结束到导致 子线程http post 没有发送
//发送文本的信息
func (dh *dingHook) FireText(e *logrus.Entry) error {
	msg, err := e.String()
	if err != nil {
		return err
	}
	dm := dingMsg{MsgType: "text"}
	dm.Text.Content = fmt.Sprintf("%s \n %s", dh.appName, msg)
	bs, err := json.Marshal(dm)
	if err != nil {
		return err
	}
	dh.jsonBodies <- bs
	return nil
}

// 发送markdown格式
func (dh *dingHook) FireMarkDown(e *logrus.Entry) error {
	markdownString := ``
	markdownString += fmt.Sprintf("## %s %s %s\n", dh.appName, e.Level.String(), e.Time.Format("06-01-02T15:04"))
	markdownString += fmt.Sprintf("### %s\n", e.Message)
	for k, v := range e.Data {
		switch v := v.(type) {
		case error:
			markdownString += fmt.Sprintf("> %s \n", v)
		default:
			markdownString += fmt.Sprintf("## ***%s*** %s\n", k, v)
		}
	}
	if e.HasCaller() {
		funcVal := e.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", e.Caller.File, e.Caller.Line)
		if funcVal != "" {
			markdownString += fmt.Sprintf("> %s \n", funcVal)
		}
		if fileVal != "" {
			markdownString += fmt.Sprintf("> %s \n", fileVal)
		}
	}
	dm := dingMsg{MsgType: "markdown"}
	dm.Markdown.Title = dh.appName
	dm.Markdown.Text = markdownString

	bs, err := json.Marshal(dm)
	if err != nil {
		return err
	}
	res, err := http.Post(dh.apiUrl, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		return err
	}
	if res != nil && res.StatusCode != 200 {
		return fmt.Errorf("dingHook go chan http post error %d", res.StatusCode)
	}
	return nil
}

// 钉钉消息结构体
type dingMsg struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
}
