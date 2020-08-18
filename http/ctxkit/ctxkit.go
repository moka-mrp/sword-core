package ctxkit

import (
	"context"
	"github.com/gin-gonic/gin"
)

const (
	TraceId  = "x-trace-id"
	ClientIp = "x-cip"
	ServerIp = "x-sip"
	HOST     = "x-host"
	UIN      = "x-uin"
	SessionKey = "x-session-key"
	Appid    = "x-appid"
	JwtSecret = "x-jwt-secret"
)


//用户追踪一个用户触发的多个请求(日志追踪) 或者用来标识客户端的唯一请求序列，防止重复执行
func SetTraceId(ctx *gin.Context, value string) {
	ctx.Set(TraceId, value)
}

func GetTraceId(ctx context.Context) string {
	s, _ := ctx.Value(TraceId).(string)
	return s
}

//客户端IP
func SetClientId(ctx *gin.Context, value string) {
	ctx.Set(ClientIp, value)
}

func GetClientId(ctx context.Context) string {
	s, _ := ctx.Value(ClientIp).(string)
	return s
}
//服务端IP
func SetServerId(ctx *gin.Context, value string) {
	ctx.Set(ServerIp, value)
}

func GetServerId(ctx context.Context) string {
	s, _ := ctx.Value(ServerIp).(string)
	return s
}
//请求host
func SetHost(ctx *gin.Context, value string) {
	ctx.Set(HOST, value)
}

func GetHost(ctx context.Context) string {
	s, _ := ctx.Value(HOST).(string)
	return s
}

//jwt相关  sam@2020-08-18 09:23:25

func SetJwtSecret(ctx *gin.Context, value string) {
	ctx.Set(JwtSecret, value)
}

func GetJwtSecret(ctx context.Context) string {
	s, _ := ctx.Value(JwtSecret).(string)
	return s
}

//------密文相关 sam@2020-06-10 08:46:52
func SetUin(ctx *gin.Context, value string) {
	ctx.Set(UIN, value)
}

func GetUin(ctx context.Context) string {
	s, _ := ctx.Value(UIN).(string)
	return s
}

func SetSessionKey(ctx *gin.Context, value string) {
	ctx.Set(SessionKey, value)
}

func GetSessionKey(ctx context.Context) string {
	s, _ := ctx.Value(SessionKey).(string)
	return s
}

func SetAppid(ctx *gin.Context, value string) {
	ctx.Set(Appid, value)
}

func GetAppid(ctx context.Context) string {
	s, _ := ctx.Value(Appid).(string)
	return s
}

