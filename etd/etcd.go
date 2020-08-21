package etd

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/moka-mrp/sword-core/config"
	"time"
)

const DefaultLock  ="/%s/lock/"

type Client struct {
	*clientv3.Client
	reqTimeout time.Duration
}


//创建etcd/clientv3客户端一般是分两步走
//1.config = clientv3.Config{}
//2.cli, err = clientv3.New(config)

func NewClient(etcdConf config.EtcdConfig) (c *Client, err error) {
	cli, err := clientv3.New(etcdConf.Copy())
	if err != nil {
		return
	}
	c = &Client{
		Client: cli,
		reqTimeout: time.Duration(etcdConf.ReqTimeout) * time.Second,
	}
	return
}



