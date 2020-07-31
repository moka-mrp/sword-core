package config

import "time"


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

//------------------- Log相关结构体配置 -----------------------------------
type LogConfig struct {
	Handler  string
	Level    string
	Dir      string
	Name string //文件名
	EnableFileName bool //日志中是否显示调用文件的名字
	EnableFuncName bool //日志中是否显示调用的函数的名字
}
//------------------- Api相关结构体配置 -----------------------------------
type ApiConfig struct {
	Host string
	Port int
	Debug bool //是否是debug模式
}
