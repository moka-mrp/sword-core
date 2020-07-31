package rds

import (
	"fmt"
	"testing"
)

func TestGetSingleton(t *testing.T) {
	c := getSingleton("", false)
	if c != nil {
		t.Error("client is not equal nil")
		return
	}
}

//测试注册获取整个流程
//@author sam@2020-07-31 10:35:32
func TestProvider(t *testing.T) {
	//注册
	err := Pr.Register("redis", conf)
	if err != nil {
		t.Error(err)
		return
	}
	//获取提供者别名
	arr := Pr.Provides()
	fmt.Println(arr)

	//获取连接池子
	pools := GetRedis()
	if pools == nil {
		t.Error("client is equal nil")
		return
	}
	b,err:=pools.Set("testproviser","good")
	fmt.Println(b,err)
	//关闭池子
	err = Pr.Close()
	if err != nil {
		t.Error(err)
		return
	}
}
