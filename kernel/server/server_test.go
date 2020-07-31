package server

import (
	"github.com/gin-gonic/gin"
	"github.com/moka-mrp/sword-core/command"
	"github.com/moka-mrp/sword-core/command/examples"
	"github.com/moka-mrp/sword-core/config"
	"net/http"
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

///快速体验的一个案例
//@author sam@2020-04-07 11:04:03
func HandleFast(c *gin.Context) {

	data:=make(map[string]interface{})
	data["name"]="sam"
	data["age"]=18
	c.JSON(http.StatusOK, gin.H{
		"errcode":         0,
		"errmsg":         "ok",
		"data":        data,
	})
	c.Abort()

	return
}

//api路由配置
//中间件是可以添加多个的，并按照添加的函数顺序执行
//入口->中间件1(before)->中间件N(before)->控制器->中间件N(after)->中间件1(after)—>结束

func RegisterRoute(router *gin.Engine) {
	t := router.Group("/test")
	{
		t.GET("/fast",HandleFast) //快速入门体验
	}

}

//测试http服务
func  TestHttpServer(t *testing.T){

	conf:=config.ApiConfig{
		Host: "0.0.0.0",
		Port: 8089,
		Debug:true,
	}
	SetDebug(conf.Debug) //告之统一的服务

	StartHttp(conf,RegisterRoute)
}
