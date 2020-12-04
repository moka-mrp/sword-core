package server

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/v5/tunnel"
	"github.com/moka-mrp/sword-core/aliots"
	"github.com/moka-mrp/sword-core/kernel/event"
	"log"
)




var DefaultDaemon *tunnel.TunnelWorkerDaemon
func ExitTunnel(i interface{}){
	DefaultDaemon.Close()
}

//开启定时任务
//@author sam@2020-12-04 09:53:22
func StartTunnel(client tunnel.TunnelClient,tableName string,tunnelName string,
	callback aliots.ChannelCallback,shutdown aliots.ChannelShutdown,
	customValue interface{} ) error {

	//1.获取通道ID
	chId,err:=aliots.GetChannelIdByName(client,tableName,tunnelName)
	if err !=nil{
		return err
	}
	//2.获取worker配置
	tunnelWorkConfig := aliots.GetDefaultTunnelWorkerConfig(callback,shutdown,customValue)
	//3.获取daemon并异步运行
	DefaultDaemon = tunnel.NewTunnelDaemon(client,chId,tunnelWorkConfig)
	go func() {
		err := DefaultDaemon.Run()
		if err != nil {
			log.Fatal("tunnel worker fatal error: ", err)
		}
	}()
	//4.守护
	event.On(event.EXIT, ExitTunnel)
	event.Wait()
	event.Emit(event.EXIT, nil)

	return nil
}

