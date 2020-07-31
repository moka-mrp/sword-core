package rds

import (
	"fmt"
	"github.com/moka-mrp/sword-core/config"
	"testing"
)

var conf config.RedisMultiConfig

func init() {
	conf=make(config.RedisMultiConfig,2)
	conf["default"]=config.RedisConfig{
		Host:           "127.0.0.1",
		Port:           6379,
		Password:       "root",
		DB:             0,
		MaxIdle:        10,
		MaxActive:       50,
		Wait:           true,
		//IdleTimeout:    0,
		//ConnectTimeout: 0,
		//ReadTimeout:    0,
		//WriteTimeout:   0,
	}

	conf["arch"]=config.RedisConfig{
		Host:           "127.0.0.1",
		Port:           6379,
		Password:       "root",
		DB:             1,
		MaxIdle:        10,
		MaxActive:       50,
		Wait:           true,
		//IdleTimeout:    0,
		//ConnectTimeout: 0,
		//ReadTimeout:    0,
		//WriteTimeout:   0,
	}

}

//创建redis多集合连接池，然后进行set  get操作
//@todo 注意redis如果拒绝其它外网客户端连接的时候，可能会返回EOF
//@author sam@2020-07-29 15:40:34

func TestGetSet(t *testing.T) {
	//创建池子
	pools, err :=NewMultiPool(conf)
	if err != nil {
		t.Error("pools init failed")
		return
	}
	//默认池子的增查删
	b,err:=pools.Set("name2","default")
	fmt.Println(b,err)
     b2,err:=pools.Get("name2")
	fmt.Println(b2,err)

     //指明池子的增查删

	b3,err:=pools.AliasSet("arch","name3","alias")
	fmt.Println(b3,err)
	b4,err:=pools.AliasGet("arch","name3")
	fmt.Println(b4,err)



}

