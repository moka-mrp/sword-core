package etd

import (
	"context"
	"fmt"
)



//单独为etcd客户端请求延迟情况封装的上下文
//仅仅比context.Context多了个etcdEndpoints字段,同时改写了Err()
//@author sam@2020-08-21 15:55:14
type etcdTimeoutContext struct {
	context.Context
	etcdEndpoints []string
}

func (c *etcdTimeoutContext) Err() error {
	err := c.Context.Err()
	if err == context.DeadlineExceeded {
		err = fmt.Errorf("%s: etcd(%v) maybe lost", err, c.etcdEndpoints)
	}
	return err
}

// NewEtcdTimeoutContext return a new etcdTimeoutContext
func NewEtcdTimeoutContext(c *Client) (context.Context, context.CancelFunc) {
	//设置一个超时的context
	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	etcdCtx := &etcdTimeoutContext{}
	etcdCtx.Context = ctx
	etcdCtx.etcdEndpoints = c.Endpoints()
	return etcdCtx, cancel
}