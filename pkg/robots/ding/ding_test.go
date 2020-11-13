package ding

import "testing"

const token  = "ea38ad31f9654bbbd41a83dedc530e782965aeb67c255672dde301654c5eb57b"
const  who  = "13071771599"


func TestSendTex(t *testing.T) {
	webhook := "https://oapi.dingtalk.com/robot/send?access_token="+token
	robot := NewRobot(webhook)
	content := "[ots]我就是我,  @"+who+"是不一样的烟火"
	atMobiles := []string{"18565056618",who}
	isAtAll := false
	err := robot.SendText(content, atMobiles, isAtAll)
	if err != nil {
		t.Error(err)
	}


}


func TestSendLink(t *testing.T) {
	webhook := "https://oapi.dingtalk.com/robot/send?access_token="+token
	robot := NewRobot(webhook)
	title := "[ots]热烈祝贺“梦嘉集团”获园区2019年度准独角兽企业殊荣"
	text := "近日，苏州工业园区2019年度独角兽、瞪羚企业培育工程入库企业名单正式公布。江苏梦嘉控股集团有限公司在众多企业中脱颖而出，荣获“准独角兽企业”殊荣。"
	messageUrl := "https://www.mokasz.com/"
	picUrl := "https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png"
	err := robot.SendLink(title, text, messageUrl, picUrl)
	if err != nil {
		t.Error(err)
	}


}



func TestSendMarkdown(t *testing.T) {

	webhook := "https://oapi.dingtalk.com/robot/send?access_token="+token
	robot := NewRobot(webhook)
	title := "苏州天气"
	text := "[ots]#### 苏州天气  \n > 9度， 西北风1级，空气良89，相对温度73%\n\n > ![screenshot](http://i01.lw.aliimg.com/media/lALPBbCc1ZhJGIvNAkzNBLA_1200_588.png)\n  > ###### 10点20分发布 [天气](http://www.thinkpage.cn/) "
	atMobiles := []string{who}
	isAtAll := false
	err := robot.SendMarkdown(title, text, atMobiles, isAtAll)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSendActionCard(t *testing.T) {

	webhook := "https://oapi.dingtalk.com/robot/send?access_token="+token
	robot := NewRobot(webhook)

	title := "乔布斯 20 年前想打造一间苹果咖啡厅，而它正是 Apple Store 的前身"
	text := "[ots]近日，苏州工业园区2019年度独角兽、瞪羚企业培育工程入库企业名单正式公布。江苏梦嘉控股集团有限公司在众多企业中脱颖而出，荣获“准独角兽企业”殊荣。"

	singleTitle := "阅读全文"
	singleURL := "http://www.mokasz.com/news/89-cn.html"
	btnOrientation := "0"
	hideAvatar := "0"
	err := robot.SendActionCard(title, text, singleTitle, singleURL, btnOrientation, hideAvatar)
	if err != nil {
		t.Fatal(err)
	}
}

