package server

import (
	"github.com/gin-gonic/gin"
	"github.com/moka-mrp/sword-core/command"
	"github.com/moka-mrp/sword-core/command/examples"
	"github.com/moka-mrp/sword-core/config"
	"github.com/moka-mrp/sword-core/kernel/close"
	"github.com/moka-mrp/sword-core/rds"
	"net/http"
	"testing"
)


var conf config.RedisMultiConfig

func init() {
	conf=make(config.RedisMultiConfig,2)
	conf["default"]=config.RedisConfig{
		Host:           "127.0.0.1",
		Port:           6379,
		Password:       "root",
		DB:             0,
		MaxIdle:        10,
		MaxActive:       50,
		Wait:           true,
		IdleTimeout:    180,
		ConnectTimeout: 3,
		ReadTimeout:    3,
		WriteTimeout:   3,
	}

}



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

	//for i:=0;i<5;i++{
	//	time.Sleep(1 * time.Second)
	//	fmt.Println(i)
	//}
	 pools:=rds.GetRedis()
	 pools.Set("http","sam")
     name,_:=pools.Get("http")

	data:=make(map[string]interface{})
	data["name"]="sam"
	data["age"]=18
	c.JSON(http.StatusOK, gin.H{
		"errcode":         0,
		"errmsg":         "ok",
		"data":        data,
		"name":name,
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
//todo 注意优雅重启:a.开启服务之后，记录pid  b.kill -1 pid(发送平滑重启信号) c.endless会等待所有的客户端连接结束之后重新启动的(pid已变)
func  TestHttpServer(t *testing.T){

	//注册redis服务   注入别名(string) + 对应配置  + 是否惰性加载(false)
	err:= rds.Pr.Register(rds.SingletonMain, conf,true)
	if err != nil {
		return
	}
	//注册应用停止时调用的关闭服务
	close.MultiRegister(rds.Pr)
	//-----------------------
	conf:=config.ApiConfig{
		Host: "0.0.0.0",
		Port: 8089,
		Debug:true,
	}
	SetDebug(conf.Debug) //告之统一的服务
	StartHttp(conf,RegisterRoute)
}
