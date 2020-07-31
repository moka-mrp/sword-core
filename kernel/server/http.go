package server

//import (
//	"fmt"
//	"github.com/moka-mrp/sword-core/config"
//	"strconv"
//	"syscall"
//
//	"github.com/fvbock/endless"
//	"github.com/gin-gonic/gin"
//)
//
///**
// * 启动gin引擎
// * @wiki https://github.com/fvbock/endless#signals
// */
//func runEngine(engine *gin.Engine, addr string, pidPath string) error {
//	//设置gin调试模式
//	//if !GetDebug() {
//	//	gin.SetMode(gin.ReleaseMode)
//	//}
//	server := endless.NewServer(addr, engine)
//	server.BeforeBegin = func(add string) {
//		pid := syscall.Getpid()
//		if gin.Mode() != gin.ReleaseMode {
//			fmt.Printf("Actual pid is %d \r\n", pid)
//		}
//		WritePidFile(pidPath, pid)
//	}
//	err := server.ListenAndServe()
//	return err
//}
//
//// Start proxy with config file
////改成无中间件启动
////@reviser sam@2020-04-06 10:10:09
//func StartHttp(pidFile string, apiConf config.ApiConfig,registerRoute func(*gin.Engine)) error {
//	//设置运行模式
//	if !GetDebug() {
//		gin.SetMode(gin.ReleaseMode)
//	}
//	//配置路由引擎
//	//engine := gin.Default()
//	engine := gin.New()
//	registerRoute(engine)
//	addr := apiConf.Host + ":" + strconv.Itoa(apiConf.Port)
//	fmt.Printf("Start http server listening %s\r\n", addr)
//	err:=runEngine(engine, addr, pidFile)
//	if err !=nil{
//		return  err
//	}
//	//因为信号处理由endless接管实现平滑重启和关闭，这里模拟通用的结束信号
//	//@todo 下面这部分可以完全不要的额
//	go func() {
//		Stop()
//	}()
//
//	//等待停止信号
//	WaitStop()
//	return nil
//}
