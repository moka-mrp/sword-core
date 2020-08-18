package ctxkit

import (
	"testing"
	"github.com/gin-gonic/gin"
)

var c *gin.Context

func init() {
	c = &gin.Context{}
}


//------------------------------------http基本相关------------------------

//用户追踪一个用户触发的多个请求(日志追踪) 或者用来标识客户端的唯一请求序列，防止重复执行
//@author sam@2020-08-18 14:34:09
func TestGetTraceId(t *testing.T) {
	v := "000000000000000000001"
	SetTraceId(c, v)
	v1 := GetTraceId(c)
	if v1 != v {
		t.Error("TraceId miss match")
		return
	}
}


//客户端IP
//@author sam@2020-08-18 14:27:32
func TestGetClientId(t *testing.T) {
	v := "218.4.157.186"
	SetClientId(c, v)
	v1 := GetClientId(c)
	if v1 != v {
		t.Error("ClientId miss match")
		return
	}
}

//服务端IP
//@author sam@2020-08-18 14:36:31
func TestGetServerId(t *testing.T) {
	v := "112.80.248.76"
	SetServerId(c, v)
	v1 := GetServerId(c)
	if v1 != v {
		t.Error("ServerId miss match")
		return
	}
}
//host
//@author sam@2020-08-18 14:35:41
func TestGetHost(t *testing.T) {
	v := "www.baidu.com"
	SetHost(c, v)
	v1 := GetHost(c)
	if v1 != v {
		t.Error("Host miss match")
		return
	}
}


//------------------------------------jwt基本相关------------------------
//jwt secret
//@author sam@2020-08-18 14:39:10
func TestGetJwtSecret(t *testing.T) {
	v := "aakdsfjklasdjflasdjccc"
	SetJwtSecret(c, v)
	v1 := GetJwtSecret(c)
	if v1 != v {
		t.Error("Secret miss match")
		return
	}
}


//------------------------------------aes与rsa基本相关------------------------

//测试用户id存储
//@author sam@2020-08-18 14:41:07
func TestGetUin(t *testing.T) {
	v := "100000000000000000001"
	SetUin(c, v)
	v1 := GetUin(c)
	if v1 != v {
		t.Error("Secret miss match")
		return
	}
}


//测试aes的对称加密令牌
//@author sam@2020-08-18 14:41:07
func TestGetSessionkey(t *testing.T) {
	v := "bbbbbbbbbbbbbbbbbbbbbb"
	SetSessionKey(c, v)
	v1 := GetSessionKey(c)
	if v1 != v {
		t.Error("SessionKey miss match")
		return
	}
}



//测试appid
//@author sam@2020-08-18 14:41:07
func TestGetAppid(t *testing.T) {
	v := "wx00000000001"
	SetAppid(c, v)
	v1 := GetAppid(c)
	if v1 != v {
		t.Error("Appid miss match")
		return
	}
}


