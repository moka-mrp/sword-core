package etd

import (
	"fmt"
	"context"
	"github.com/coreos/etcd/clientv3"
)

//-----------------------------------------添加--------------------------------------------
//添加一个key
//todo 相对原生的好处是不用指明上下文取消时间的设置
//@author sam@2020-08-21 15:47:52
func (c *Client) Put(key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	ctx, cancel := NewEtcdTimeoutContext(c)
	defer cancel()
	return c.Client.Put(ctx, key, val, opts...)
}


//持乐观锁进行put操作
//Revision 是全局的版本，不针对key,但某个key的创建和修改的版本号也是全局版本号范畴
//CreateRevision  是当前key创建的时候的版本号
//ModRevision 是当前key最后一次修改的版本(针对单个key)
//version 是某个key的修改次数递增
//todo 比如某个key的Revision是1150,则不管接下来整体版本号是否增加了，只要对该key的修改必须持1150版本号才可以额
//@author sam@2020-08-21 16:12:12
func (c *Client) PutWithModRev(key, val string, rev int64,opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	//如果为0表示首次添加key值，或者是强制不持乐观锁，直接input即可
	if rev == 0 {
		return c.Put(key, val,opts ...)
	}
	//自定义一个超时上下文
	ctx, cancel := NewEtcdTimeoutContext(c)
	defer cancel()
	// 创建事务-定义事务(版本号此刻未变，才可以提交事务)-提交事务
	txnResp, err:=c.Txn(ctx).If(clientv3.Compare(clientv3.ModRevision(key), "=", rev)).
				            Then(clientv3.OpPut(key, val,opts...)).
						    Else(clientv3.OpGet(key)).
				            Commit()
	if err != nil {
		return nil, err
	}
	//判断事务是否提交成功
	if !txnResp.Succeeded {
		//fmt.Println(string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return nil, ErrValueMayChanged
	}
	resp := clientv3.PutResponse(*txnResp.Responses[0].GetResponsePut())
	return &resp, nil
}

//-------------------------------------------查看-------------------------------------------------------------

//获取key
//@author sam@2020-08-21 17:03:18
func (c *Client) Get(key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	ctx, cancel := NewEtcdTimeoutContext(c)
	defer cancel()
	return c.Client.Get(ctx, key, opts...)
}


//-----------------------------------------删除-----------------
//删除某个key
//@author sam@2020-08-21 17:11:04
func (c *Client) Delete(key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	ctx, cancel := NewEtcdTimeoutContext(c)
	defer cancel()
	return c.Client.Delete(ctx, key, opts...)
}

//-----------------------------------------监听-----------------

//监听某部分key的变化情况
//一般监听都是长期监听的，所以context可以不单独传递
//@author  sam@2020-08-21 17:11:57
func (c *Client) Watch(key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	return c.Client.Watch(context.Background(), key, opts...)
}

//--------------------------------租约 -----------------------------------------------

//直接创建一个指明时间的租约
//@author  sam@2020-08-21 17:24:02
func (c *Client) Grant(ttl int64) (*clientv3.LeaseGrantResponse, error) {
	ctx, cancel := NewEtcdTimeoutContext(c)
	defer cancel()
	return c.Client.Grant(ctx, ttl)
}
//销毁租约
//@author sam@2020-08-21 17:25:04
func (c *Client) Revoke(id clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	defer cancel()
	return c.Client.Revoke(ctx, id)
}

//给该租约ID续约一次
//@author  sam@2020-08-21 17:25:32
func (c *Client) KeepAliveOnce(id clientv3.LeaseID) (*clientv3.LeaseKeepAliveResponse, error) {
	//创建一个超时上下文
	ctx, cancel := NewEtcdTimeoutContext(c)
	defer cancel()
	return c.Client.KeepAliveOnce(ctx, id)
}


//给该租约ID长期续租
//@author  sam@2020-08-21 17:25:32
func (c *Client) KeepAlive(id clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	return c.Client.KeepAlive(context.TODO(), id)
}



//----------------------分布式锁--------------------------------

//分布式抢锁
//@author sam@2020-08-21 17:25:48
func (c *Client) GetLock(resource,key string, id clientv3.LeaseID) (bool, error) {
	key = fmt.Sprintf(DefaultLock,resource) + key
	ctx, cancel := NewEtcdTimeoutContext(c)
	defer  cancel()
	resp, err := c.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "locked", clientv3.WithLease(id))).
		Commit()
	if err != nil {
		return false, err
	}
	return resp.Succeeded, nil
}
//主动删除分布式锁
//@author sam@2020-08-21 17:44:02
func (c *Client) DelLock(resource,key string) error {
	_, err := c.Delete(fmt.Sprintf(DefaultLock,resource) + key)
	return err
}



