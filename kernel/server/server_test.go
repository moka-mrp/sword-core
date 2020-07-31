package server

import (
	"github.com/moka-mrp/sword-core/command"
	"github.com/moka-mrp/sword-core/command/examples"
	"testing"
)

//----------------------------------- 测试command服务 -----------------------------------------------------------------


//命令行的注入在这里，整体模式与使用robfig/cron一样一样的
//@author sam@2020-04-14 13:45:21
func RegisterCommand(c *command.Command) {
	c.AddFunc("sam", examples.Sam)
	c.AddFunc("tom", examples.Tom)
	c.AddFunc("rick", examples.Rick)
}

//go run main.go  -a command -m sam
func  TestCommandServer(t *testing.T){
	ExecuteCommand("sam", RegisterCommand)
	ExecuteCommand("tom", RegisterCommand)
	ExecuteCommand("rick", RegisterCommand)

}




//---------------------------------------  测试http服务 --------------------------------------------------------------
//测试http服务
func  TestHttpServer(t *testing.T){

}
