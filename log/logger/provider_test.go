package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/moka-mrp/sword-core/config"
	"testing"
)

var contextTest, contextTest1 *gin.Context

func init() {
	contextTest = &gin.Context{}
	contextTest1 = &gin.Context{}
}

//测试一下从容器中获取单例
//@author sam@2020-07-30 13:47:35
func TestGetSingleton(t *testing.T) {
	//非惰性直接获取未注册的资源必然panic,因为非惰性，必然主动帮你注册，但是你没有提供配置啊，所以注入不了
	c := getSingleton("logger",false)
	fmt.Printf("%+v\r\n",c)
	if c != nil {
		t.Error("client is not equal nil")
		return
	}
}

//测试整个注入到获取的一个过程
//@author sam@2020-07-30 13:54:27
func TestProvider(t *testing.T) {
	//一、注入一个标准输出配置的logger
	err := Pr.Register("logger1", config.LogConfig{
		Handler: "stdout",
		Level:   "debug",
		Dir:     "./logs",
		Name : "lumen",
		EnableFileName:true,
		EnableFuncName:true,
	})
	if err != nil {
		t.Error(err)
		return
	}
    //todo 其实不建议使用语法糖方法调用，这样很难追踪报错位置了
	Debug("a debug")
	Info("a info")
	Warn("a warn")
	Error("a error")
	//Fatal("a fatal")  //程序会结束
	//Panic("a panic")  //程序会结束

	GetLogger("logger1").Info("d001")
	GetLogger().Info("d002")

	//二、注入一个写入到日志文件

	err = Pr.Register("logger2",config.LogConfig{
		Handler:  "file",
		Level:    "debug",
		Dir:      "./logs",
		Name: "sword",
		EnableFileName:false,
		EnableFuncName:false,
	}, true)
	if err != nil {
		t.Error(err)
		return
	}

	GetLogger("logger2").Info("x001")
	GetLogger("logger2").Info("x002")
	GetLogger("logger2").Info("x003")






	//关闭所有的日志资源
	err = Pr.Close()
	if err != nil {
		t.Error(err)
		return
	}

}



