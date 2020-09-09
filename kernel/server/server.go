package server

import (
	"github.com/moka-mrp/sword-core/kernel/close"
	"log"
)


type serverInfo struct {
	stop  chan bool  //这个通道中有值，则服务结束
	debug bool  //是否测试模式启动
}

var srv *serverInfo

func init() {
	srv = new(serverInfo)
	srv.stop = make(chan bool, 0)
}

//关闭整个服务启动过程加载的资源
//@author  sam@2020-09-09 11:45:07
func CloseService() {
	if srv.debug {
		log.Println("close boot  resources ...")
	}
	close.Free()
}


//--------------------------------------------------Debug-----------------------------------------------------
//给外部设置启动服务的debug模式
//@author sam@2020-04-13 17:37:46
func SetDebug(debug bool) {
	srv.debug = debug
	return
}
//给外部获取启动服务的debug模式
//@author  sam@2020-04-13 17:38:26
func GetDebug() bool {
	return srv.debug
}
