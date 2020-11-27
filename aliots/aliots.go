package aliots

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/v5/tunnel"
	"github.com/moka-mrp/sword-core/config"
	"sync"
	"github.com/aliyun/aliyun-tablestore-go-sdk/v5/tablestore"
)

type MultiClient struct {
	mu  sync.RWMutex
	otsClients  map[string]*tablestore.TableStoreClient
	otsTunnelClients map[string]tunnel.TunnelClient
}

func (mc *MultiClient)SetOtsClient(name string,client *tablestore.TableStoreClient){
	mc.mu.Lock()
	mc.otsClients[name] = client
	mc.mu.Unlock()
}
func (mc *MultiClient)SetTunnelClient(name string,client tunnel.TunnelClient){
	mc.mu.Lock()
	mc.otsTunnelClients[name] = client
	mc.mu.Unlock()
}

//从多实例中获取单个ots client
func (mc *MultiClient)GetOtsClient(name string)(*tablestore.TableStoreClient){
	mc.mu.RLock()
	ins, _ := mc.otsClients[name]
	mc.mu.RUnlock()
	return ins
}

//从多实例中获取单个tunnel client
func (mc *MultiClient)GetTunnelClient(name string)(tunnel.TunnelClient){
	mc.mu.RLock()
	ins, _ := mc.otsTunnelClients[name]
	mc.mu.RUnlock()
	return ins
}

//创建多实例ots连接客户端
//@author sam@2020-11-27 15:31:45
func NewMultiClient(conf config.OtsMultiConfig) (*MultiClient,error) {

	mc:=new(MultiClient)
	mc.otsClients=make(map[string]*tablestore.TableStoreClient,len(conf))
	mc.otsTunnelClients=make(map[string]tunnel.TunnelClient,len(conf))

	for name,otsConfig:=range conf{
		otsClient:=tablestore.NewClient(otsConfig.EndPoint,otsConfig.InstanceName,otsConfig.AccessKeyId,otsConfig.AccessKeySecret)
		mc.SetOtsClient(name,otsClient)
		tunnelClient:=tunnel.NewTunnelClient(otsConfig.EndPoint,otsConfig.InstanceName,otsConfig.AccessKeyId,otsConfig.AccessKeySecret)
		mc.SetTunnelClient(name,tunnelClient)
	}
	return  mc,nil
}


