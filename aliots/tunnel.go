package aliots

import (
	"github.com/aliyun/aliyun-tablestore-go-sdk/v5/tablestore"
	"github.com/aliyun/aliyun-tablestore-go-sdk/v5/tunnel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)
const (
	OTS_CONDITION_CHECK_FAIL       = "OTSConditionCheckFail"
)

type ChannelCallback  func(channelCtx *tunnel.ChannelContext, records []*tunnel.Record)error
type ChannelShutdown   func(channelCtx *tunnel.ChannelContext)

var (
	DefaultHeartbeatInterval = 30 * time.Second   //默认30
	DefaultHeartbeatTimeout  = 300 * time.Second  //默认300,注意超时时间一定要大于间隔时间
	DefaultChannelSize        = 10
	DefaultCheckpointInterval = 10 * time.Second //checkpoint的间隔时间，默认是10s
)

//自定义客户端配置--------------------------------------------------------------------------
//todo 在sword-core中已经采用默认值，不支持灵活配置，也不建议，默认值就是最优的
var DefaultTunnelConfig = &tunnel.TunnelConfig{
	MaxRetryElapsedTime: 75 * time.Second, //最大指数退避重试时间。  默认75
	RequestTimeout:     60 * time.Second,//HTTP请求超时时间。  默认60
	Transport:           http.DefaultTransport,//http.DefaultTransport。
}

//日志整体配置--------------------------------------------

//构建日志logger



var DefaultLogConfig = zap.Config{
	Level:       zap.NewAtomicLevelAt(zap.WarnLevel),//日志级别
	Development: false,
	Sampling: &zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	},
	Encoding: "console",// 输出格式 console 或 json, todo console这种格式不是说要输出到控制台
	EncoderConfig: zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		//CallerKey:      "caller", //关闭文件位置
		MessageKey:     "msg",
		//StacktraceKey:  "stacktrace", //关闭堆栈
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	},
	OutputPaths:      []string{"stdout"},
	ErrorOutputPaths: []string{"stderr"},
}

//日志切割
var DefaultSyncer = zapcore.AddSync(&lumberjack.Logger{
	Filename:   "./logs/tnl.log",//日志文件路径。
	MaxSize:    100, //MB //最大日志文件大小。每一百兆就进行切割一次
	MaxBackups: 5,  //压缩轮转的日志文件数。
	MaxAge:     7, //days //轮转日志文件保留的最大天数。
	Compress:   true, //是否压缩轮转日志文件。
})

//api请求回退配置---------------------------------------------------------------
var DefaultBackoffConfig = tunnel.ChannelBackoffConfig{
	MaxDelay:  5 * time.Second,
	//baseDelay: 20 * time.Millisecond,
	//factor:    5,
	//jitter:    0.25,
}





//lg, err := cloneConf.LogConfig.Build(ReplaceLogCore(cloneConf.LogWriteSyncer, *cloneConf.LogConfig))
//lg, err := cloneConf.LogConfig.Build()
//通过几个可配置参数，获取默认的tunnel worker config
//@author  sam@2020-11-28 09:16:50
func GetDefaultTunnelWorkerConfig(callback ChannelCallback,shutdown ChannelShutdown,customValue interface{})(*tunnel.TunnelWorkerConfig){
	logConfig:=&DefaultLogConfig
	logger,_:=DefaultLogConfig.Build(tunnel.ReplaceLogCore(DefaultSyncer, DefaultLogConfig)) //带日志轮转切割的
	//logger,_:=DefaultLogConfig.Build()
	workConfig := &tunnel.TunnelWorkerConfig{
		HeartbeatInterval:DefaultHeartbeatInterval,//worker发送心跳的频率
		HeartbeatTimeout: DefaultHeartbeatTimeout,//worker同Tunnel服务的心跳超时时间
		ChannelDialer:nil, //该值为nil,不需要配置,在创建tunnel worker的时候会初始化该实例的
		ProcessorFactory: &tunnel.SimpleProcessFactory{
			CustomValue:customValue,//如果客户端需要将某个值回传给消费函数来使用的，可以在这里设置，多个值可以设置成结构体
			CpInterval:DefaultCheckpointInterval,//processor checkpoint的间隔时间，默认就是10s
			ProcessFunc: callback,
			ShutdownFunc: shutdown,
			Logger:logger,//todo 注意这里的日志是processor专用的，我们可以在外部直接构造的额
		},
		LogConfig:logConfig, //todo 注意这里的日志是worker专用
		LogWriteSyncer:DefaultSyncer,
		BackoffConfig:&DefaultBackoffConfig,
	}
	return workConfig
}


//根据表名、通道名,获取通道ID
//@author sam@2020-11-28 11:08:00
func GetChannelIdByName(client tunnel.TunnelClient,tableName,channelName string)(string,error){
	req := &tunnel.DescribeTunnelRequest{
		TableName: tableName,
		TunnelName:channelName,
	}
	resp, err := client.DescribeTunnel(req)
	if err != nil {
		return "",err
	}
	return resp.Tunnel.TunnelId,nil
}

//通道消费的模拟欢送函数
//@author sam@2020-11-28 13:41:32
func Welcome(timestamp int64,count int64){
	var tell string
	if timestamp <=0{
		tell="[ BaseData without timestamp ]"
	}else{
		sec, _ := strconv.ParseInt(strconv.Itoa(int(timestamp))[0:10], 10, 64)
		tell=time.Unix(sec, 0).Format("2006-01-02 15:04:05")
	}
	log.Println(">>>>>  ",tell,"Record:",count)
}


//判断该错误类型是否不在重试范围内
//@author sam@2020-11-24 14:18:34
func RecordShouldRetry(err error) bool {

	if tErr, ok := err.(*tablestore.OtsError); ok {
		return  tErr.Code==OTS_CONDITION_CHECK_FAIL
	}
	if err == io.EOF || err == io.ErrUnexpectedEOF ||
		strings.Contains(err.Error(), io.EOF.Error()) || //retry on special net error contains EOF or reset
		strings.Contains(err.Error(), "Connection reset by peer") {
		return true
	}

	if nErr, ok := err.(net.Error); ok {
		return nErr.Temporary()
	}

	return false
}