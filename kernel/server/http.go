package server

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/moka-mrp/sword-core/config"
	"strconv"
)

/**
* 优雅重启或停止
* endless内部会打印  2020/09/09 13:48:26 8015 0.0.0.0:8089
* @author sam@2020-07-31 17:37:26
* @wiki https://github.com/fvbock/endless#signals
*/
func runEngine(engine *gin.Engine, addr string) error {
	server := endless.NewServer(addr, engine)
	err := server.ListenAndServe()
	return err
}


//改成无中间件启动
//@author sam@2020-04-06 10:10:09
func StartHttp(apiConf config.ApiConfig,registerRoute func(*gin.Engine)) error {
	//设置运行模式
	if !GetDebug() {
		gin.SetMode(gin.ReleaseMode)
	}
	//配置路由引擎
	engine := gin.New()
	registerRoute(engine)
	addr := apiConf.Host + ":" + strconv.Itoa(apiConf.Port)
	err:=runEngine(engine, addr)
	if err !=nil{
		CloseService() 	//关闭启动资源
		fmt.Println(err)
		return err
	}
	//因为信号处理由endless接管实现平滑重启和关闭，这里不需要使用server包中统一封装的阻塞以及关闭了
	return nil
}
