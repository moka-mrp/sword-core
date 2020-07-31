package provider


//所有启动资源的注入都要创建一个实现该接口的提供者
//@author sam@2020-07-29 09:30:48
type Provider interface {
	Register(args ...interface{}) error
	Provides() []string   //注册过的别名
	Close() error //释放资源
}
