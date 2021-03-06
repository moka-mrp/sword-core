package etd

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/moka-mrp/sword-core/config"
	"log"
	"testing"
	"time"
)


var client *Client
//-------------------------------- init -----------------------------------------------
func init() {
	//初始化的时候会注入的
	etcdInit(true)
	//从容器中获取资源
	client= GetEtcd()
}

//注入容器以及从容器中快速取出来
func etcdInit(lazyBool bool) {

	etcdConf := config.EtcdConfig{
		Endpoints:  []string{"127.0.0.1:2379"},
		Username:    "",
		Password:    "",
		DialTimeout: 2,
		ReqTimeout:  3,
	}
	//测试容器注入功能(容器本身已经自动在kernel/container/app.go中初始化好了)
	err := Pr.Register(SingletonMain, etcdConf, lazyBool)
	if err != nil {
		fmt.Println(err)
	}
}

//-----------------------------------------添加----------------------------------------------
//测试添加一个key的操作
//@author sam@2020-08-21 16:00:13
func TestPut(t *testing.T){
	if putResp, err := client.Put("/test/food", "apple", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Revision:", putResp.Header.Revision) //整体的版本号
		if putResp.PrevKv != nil {	// 打印修改之前的值
			fmt.Println("PrevValue:", string(putResp.PrevKv.Value))
		}
	}
}


//测试添加一个key的操作(持乐观锁)
//@author sam@2020-08-21 16:37:26
func TestPutWithModRev(t *testing.T){
	if putResp, err := client.PutWithModRev("/test/food", "apple2",1164, clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Revision:", putResp.Header.Revision)
		if putResp.PrevKv != nil {	// 打印修改之前的值
			fmt.Println("PrevValue:", string(putResp.PrevKv.Value))
		}
	}
}

//-----------------------------------------查看----------------------------------------------

//测试添加一个key的操作(持乐观锁)
//@author sam@2020-08-21 16:37:26
func TestGet(t *testing.T){
	if getResp, err := client.Get( "/test/food", /*clientv3.WithCountOnly()*/); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(getResp.Kvs, getResp.Count)
		//[key:"/test/food" create_revision:1140 mod_revision:1163 version:12 value:"apple" ] 1
	}


}

//-----------------------------------------删除----------------------------------------------
//测试删除
//@author  sam@2020-08-22 09:38:17
func TestDelete(t *testing.T)  {
	// 删除KV   clientv3.WithFromKey()
	if delResp, err := client.Delete("/test/food", clientv3.WithFromKey()); err != nil {
		fmt.Println(err)
		return
	}else{
		fmt.Printf("%+v\r\n",delResp)
		// 被删除之前的value是什么
		if len(delResp.PrevKvs) != 0 {
			for _, kvpair := range delResp.PrevKvs {
				fmt.Println("删除了:", string(kvpair.Key), string(kvpair.Value))
			}
		}
	}

}

//-----------------------------------------监听----------------------------------------------
//测试watch机制
//@author sam@2020-08-22 09:41:51
func  TestWatch(t *testing.T){

	//1.模拟etcd中key-value的变化情况,比如我们模拟一个定时任务的变化情况
	key:="/crontab/jobs/jb001"
	go func() {
		for {
			client.Put(key, "put001")
			time.Sleep(5 * time.Second)
			client.Put(key, "put002")
			time.Sleep(5 * time.Second)
			client.Put(key, "put003")
			time.Sleep(5 * time.Second)
			client.Put(key, "put004")
			time.Sleep(5 * time.Second)

			client.Delete(key)
			time.Sleep(10 * time.Second)
			fmt.Println("------------------------")
		}
	}()
	//2.先GET到当前的值，并监听后续变化
	if getResp, err := client.Get(key); err != nil {
		t.Error(err)
		return
	}else{
		// 当前etcd集群事务ID, 单调递增的
		watchStartRevision := getResp.Header.Revision + 1
		// 启动监听
		fmt.Println("从该版本向后监听:", watchStartRevision)
		watchRespChan := client.Watch(key, clientv3.WithRev(watchStartRevision))
		// 处理kv变化事件
		for watchResp := range watchRespChan {
			for _, event := range watchResp.Events {
				switch event.Type {
				case mvccpb.PUT:
					fmt.Println("修改为:", string(event.Kv.Value), "CreateRevision:", event.Kv.CreateRevision,"ModRevision", event.Kv.ModRevision)
				case mvccpb.DELETE:
					fmt.Println("删除了","CreateRevision:", event.Kv.CreateRevision,"ModRevision", event.Kv.ModRevision)
				}
			}
		}



	}



}

//-----------------------------------------租约----------------------------------------------

func TestGrant(t *testing.T){
	var err error
	//1.申请一个10秒的租约   etcdctl lease grant  10
    var leaseGrantResp *clientv3.LeaseGrantResponse
	if leaseGrantResp, err = client.Grant(10); err != nil {
		fmt.Println(err)
		return
	}
	//2.租约续约(代码中还需要应答额)  etcdctl lease keep-alive  694d69eb2b36657f
	if keepRespChan, err := client.KeepAlive(leaseGrantResp.ID); err != nil {
		fmt.Println(err)
		return
	}else {
		go func() {
			for {
				select {
				case keepResp := <-keepRespChan:
					if keepRespChan == nil {
						fmt.Println("租约已经失效了")
						goto END
					} else { // 每秒会续租一次, 所以就会收到一次应答
						log.Println("收到自动续租应答:", keepResp.ID)
					}
				}
			}
		END:
		}()
	}

	//3. 使用上面申请的租约ID授权租约  etcdctl put --lease=694d69eb2b36657f  /crontab/lock/job1  ""  #
		if putResp, err := client.Put("/crontab/lock/jb001", "locked", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
			fmt.Println(err)
			return
		}else{
			fmt.Println("写入成功:", putResp.Header.Revision)
			// 定时的看一下key过期了没有,死循环阻塞中...
			for {
				var getResp *clientv3.GetResponse
				if getResp, err = client.Get("/crontab/lock/jb001"); err != nil {
					fmt.Println(err)
					return
				}
				if getResp.Count == 0 {
					fmt.Println("kv过期了")
					break
				}
				log.Println("还没过期:", getResp.Kvs)
				time.Sleep(30 * time.Second)
			}

		}
}

//-----------------------------------------分布式锁----------------------------------------------