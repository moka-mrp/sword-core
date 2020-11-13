package config

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)


//-------------------  redis 相关结构体配置 -----------------------------------
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int //第几个库，默认0

	MaxIdle        int // 最大的空闲连接数，即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
	MaxActive       int //表示和数据库的最大链接数， 0 表示没有限制
	Wait           bool //如果该值为true,则当存在获取连接恰好超过MaxActive限制的时候，则将等待有连接回到连接池之后才返回给当前的程序
	IdleTimeout    time.Duration  //空闲连接等待时间，超过此时间后，空闲连接将被关闭
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

type RedisMultiConfig  map[string]RedisConfig
//------------------- Db相关结构体配置 -----------------------------------
//数据库基本配置项
//@author sam@2020-07-01 10:38:11
type DbBaseConfig struct {
	Host     string //IP地址
	Port     int   //端口号
	User     string //用户名
	Password string  //密码
	Name   string  //数据库名
}

//数据库资源配置项
//@author sam@2020-07-01 10:38:33
type DbOptionConfig struct {
	MaxIdleConns   int //  连接池-最大的空闲连接数，即使没有mysql连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
	MaxOpenConns   int //连接池-表示和数据库的最大链接数， 0 表示没有限制
	ConnMaxLifetime   time.Duration //连接池-空闲连接等待时间，超过此时间后，空闲连接将被关闭,即 设置可以重用连接的最长时间。
	ConnectTimeout time.Duration //服务器连接超时时间
	Charset        string  //服务器连接字符编码
}

type DbConfig struct {
	Driver string //驱动类型，目前支持mysql
	Master DbBaseConfig
	Slaves []DbBaseConfig
	Option DbOptionConfig
}

//------------------- Etcd相关结构体配置 -----------------------------------
type EtcdConfig struct {
	Endpoints   []string //["http://127.0.0.1:2379"]
	Username    string // ""
	Password    string // ""
	DialTimeout int64 // 2 单位秒
	ReqTimeout  int   //// etcd客户端的,请求超时时间，单位秒  之所以放在外层，是因为是通过控制上下文时间来达成效果的
	conf clientv3.Config
}

func (e EtcdConfig) Copy() clientv3.Config {
	 e.conf.Endpoints=e.Endpoints
	 e.conf.Username=e.Username
	 e.conf.Password=e.Password
     if e.DialTimeout > 0{
     	e.conf.DialTimeout= time.Duration(e.DialTimeout) * time.Second
	 }
	return e.conf
}


//------------------- Log相关结构体配置 -----------------------------------
type LogConfig struct {
	Handler  string
	Level    string
	Dir      string
	Name string //文件名
	EnableFileName bool //日志中是否显示调用文件的名字
	EnableFuncName bool //日志中是否显示调用的函数的名字
}

//---------------------JWTConf 签名方法配置-------------------------------------------------------------------
//todo  algo = "HS512" 签名方式没必要配置,我们默认就采用HS512,至于对方采用什么我们完全可以分析token的Header部分是能够解析出来的
type JwtConfig struct {
	Secret string
	Exp    int     //token 有效期(小时)
	Algo   string
}


//------------------- Api相关结构体配置 -----------------------------------
type ApiConfig struct {
	Host string
	Port int
	Debug bool //是否是debug模式
}

//-------------Ding相关的结构体配置-------
type DingConfig struct {
	HookLink string
	AtMobiles []string
}