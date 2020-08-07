package rds

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/moka-mrp/sword-core/config"
	"sync"
	"time"
)


var (
	ErrNil = redis.ErrNil
	DefaultPool= "default"
)
//-------------------------多集合连接池结构体-----------------------

type MultiPool struct {
	mu         sync.RWMutex
	pools      map[string]*redis.Pool
}

// 注册单例对象
func (mp *MultiPool) SetPool(name string, pool  *redis.Pool) {
	mp.mu.Lock()
	mp.pools[name] = pool
	mp.mu.Unlock()
}

// 获取单例对象
func (mp *MultiPool) GetSingleton(name string) *redis.Pool {
	mp.mu.RLock()
	ins, _ := mp.pools[name]
	mp.mu.RUnlock()
	return ins
}



//当应用关闭的时候，可以释放所有池子里的所有链接额
//@author sam@2020-04-09 17:15:38
func (mp *MultiPool) Close() {
	for _, pool := range mp.pools {
		pool.Close()
	}
}

//捞取一个池子 *redis.Pool
//0是默认的池子 1是额外的第一个池子 2是额外的第二个池子 依次类推
//@author sam@2020-04-09 16:59:22
func (mp *MultiPool) PickPool(poolName string) (*redis.Pool,error) {
	// FIXME: check liveness and retry
	if pool,ok:=mp.pools[poolName];ok{
		return  pool,nil
	}else{
		return nil,fmt.Errorf("%s not register in pools",poolName)
	}
}

//捞取一个池子中的一个连接 redis.Conn
//@todo 注意外部调用者记得 defer conn.Close()
//@author sam@2020-04-09 16:59:22
func (rp *MultiPool) GetConn(poolName string) (redis.Conn,error) {
	 pool,err:=rp.PickPool(poolName)
	 if err !=nil{
	 	return nil,err
	 }
	 return pool.Get(),nil
}



//针对默认池子封装的方法糖
//@todo 注意该方法就是一次操作，自动放回了连接池了，所以外部使用就不需要再关闭了，但是使用额外连接池的时候需要调用者手动关闭
//@author sam@2020-07-29 17:54:56
func (rp *MultiPool) Do(poolName string,commandName string, args ...interface{}) (reply interface{}, err error) {
	conn,err:=rp.GetConn(poolName)
	if err !=nil{
		return  nil,err
	}
	defer conn.Close()
	return conn.Do(commandName, args...)
}

//----------------------------创建多集合连接池--------------------------------------------------------------
//@todo 单个连接池的使用步骤是 redis.Pool ---->  redis.Conn (defer conn.Close()) --->发送指令即可
//@author sam@2020-04-09 14:48:57
func NewMultiPool(conf config.RedisMultiConfig) (*MultiPool,error) {

	mp:=new(MultiPool)
	mp.pools=make(map[string]*redis.Pool,len(conf))
	for alias,redisConfig:=range conf{
		pool,_:=newConn(redisConfig)
		mp.SetPool(alias,pool)
	}
   return  mp,nil
}

//返回一个redis pool
//@author sam@2020-07-30 17:47:24
func newConn(conf config.RedisConfig)(*redis.Pool,error){
	pool:= &redis.Pool{
		Dial: func() (redis.Conn, error) { // 连接redis的方法
			addr := fmt.Sprintf("%s:%d", conf.Host, getPortOrDefault(conf.Port))
			conn, err :=redis.Dial("tcp", addr,
				redis.DialConnectTimeout(conf.ConnectTimeout * time.Second),
				redis.DialReadTimeout(conf.ReadTimeout * time.Second),
				redis.DialWriteTimeout(conf.WriteTimeout * time.Second))
			if err != nil {
				fmt.Println("Dial err=",err)
				return conn, err
			}
			if conf.Password != "" {
				if _, err := conn.Do("AUTH", conf.Password); err != nil {
					conn.Close()
					panic(err)
					return nil,err
				}
			}
			if conf.DB != 0 {
				if _, err := conn.Do("SELECT", conf.DB); err != nil {
					conn.Close()
					panic(err)
					return nil, err
				}
			}
			return conn, nil
		},
	}
	//设置连接参数
	if  conf.MaxIdle > 0 {
		pool.MaxIdle = conf.MaxIdle
	}
	if conf.Wait{
		pool.Wait = conf.Wait
	}
	if conf.MaxActive >0{
		pool.MaxActive = conf.MaxActive
	}
	if conf.IdleTimeout >0{
		pool.IdleTimeout = conf.IdleTimeout * time.Second
	}

	return pool,nil
}

//获取默认端口号
//@author sam@2020-07-30 17:46:51
func  getPortOrDefault(port int) int {
	if port == 0 {
		return 6379
	} else {
		return port
	}
}
