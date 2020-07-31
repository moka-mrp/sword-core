package close

import "sync"

//我们封装的启动资源，如果需要在服务结束的时候释放对应资源，则建议在启动bootstrap之后，随手注入对应的关闭
//比如db、redis资源
//close.MultiRegister(db.Pr, redis.Pr)

var (
	closeSet []Closeable
	lock     sync.RWMutex
)

//--------------Closeable 接口------------------------------------------------
//资源提供者必然也实现了该接口
//@author sam@2020-07-31 14:57:16
type Closeable interface {
	Close() error
}
//----------------- 往全局切片中注入要关闭的资源-----------------------------------

//注册应用停止时需要释放链接的服务
func Register(closeable Closeable) {
	lock.Lock()
	defer lock.Unlock()
	closeSet = append(closeSet, closeable)
}

//批量注册应用停止时需要释放链接的服务 (资源启动的时候会添加进来)
//@author sam@2020-04-14 15:15:27
func MultiRegister(closeableSet ...Closeable) {
	lock.Lock()
	defer lock.Unlock()
	closeSet = append(closeSet, closeableSet...)
}

//-------------------将全局切片中的资源释放------------------------------------------------------------

//释放链接
//@author sam@2020-04-14 15:14:29
func Free() {
	for _, v := range closeSet {
		if v != nil {
			v.Close()
		}
	}
}



