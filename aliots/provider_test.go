package aliots

import (
	"fmt"
	"github.com/aliyun/aliyun-tablestore-go-sdk/v5/tunnel"
	"testing"
)



//测试注册获取整个流程
//@author sam@2020-11-27 15:54:08
func TestProvider(t *testing.T) {
	//注册
	err := Pr.Register(SingletonMain, conf)
	if err != nil {
		t.Error(err)
		return
	}
	//获取ots多客户端
	clients := GetOts()
	if clients == nil {
		t.Error("clients is equal nil")
		return
	}
	//列出表名称================================================
    fmt.Println("tulong:")
	listtables, err := clients.GetOtsClient("tulong").ListTable()
	if err != nil {
		t.Error(err)
	} else {
		for _, table := range listtables.TableNames {
			fmt.Println("TableName: ", table)
		}
	}

	fmt.Println("xueyin:")
	listtables2, err := clients.GetOtsClient("xueyin").ListTable()
	if err != nil {
		t.Error(err)
	} else {
		for _, table := range listtables2.TableNames {
			fmt.Println("TableName: ", table)
		}
	}

	//列出通道====================================

	req := &tunnel.DescribeTunnelRequest{
		TableName: "sbm_demos",
		TunnelName: "sbm_demos_tunnel",
	}
	resp, err := clients.GetTunnelClient("tulong").DescribeTunnel(req)
	if err != nil {
		t.Error("describe test tunnel failed", err)
	}
	fmt.Println("tunnel id is", resp.Tunnel.TunnelId)

	req2 := &tunnel.DescribeTunnelRequest{
		TableName: "test",
		TunnelName: "test_tunnel",
	}
	resp2, err := clients.GetTunnelClient("xueyin").DescribeTunnel(req2)
	if err != nil {
		t.Error("describe test tunnel failed", err)
	}
	fmt.Println("tunnel id is", resp2.Tunnel.TunnelId)






}

