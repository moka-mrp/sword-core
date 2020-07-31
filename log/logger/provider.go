package logger

import (
	"errors"
	"fmt"
	"github.com/moka-mrp/sword-core/config"
	"github.com/moka-mrp/sword-core/helper"
	"github.com/moka-mrp/sword-core/kernel/container"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

//-------------------------------- 初始化 -------------------------------------------------
const SingletonMain = "logger" //这个推荐使用的默认别名，当然注入的时候可以无视该别名，自己定义即可，第一个注入的别名会自动设置为默认别名的

var Pr *provider

func init() {
	Pr = new(provider)
	Pr.mp = make(map[string]interface{})
}

//--------------------------  服务提供者 -------------------------------------------------------------
//@author sam@2020-04-02 14:59:34
type provider struct {
	mu sync.RWMutex
	mp map[string]interface{} //配置,可以添加多个别名以及映射对应的配置额
	dn string                 //default name,一般是注入的别名,如logger
}

/**
 * @param string 依赖注入别名 必选
 * @param config.LogConfig 配置 必选
 * @param bool 是否启用懒加载 可选
 * @todo  如果是惰性加载，则仅仅完成mp和dn字段的初始化工作而已
 * @todo  如果是非惰性加载，则还会完成底层log的创建配置并注入到容器中的工作额
 * @author sam@2020-03-27 10:04:17
 */
func (p *provider) Register(args ...interface{}) (err error) {
	//(1)检测参数并提取主要参数值
	diName, lazy, err := helper.TransformArgs(args...)
	if err != nil {
		return
	}
	//(2)类型断言提取conf(此处起到检测的作用，只有非惰性加载才会需要conf)
	conf, ok := args[1].(config.LogConfig)
	if !ok {
		return errors.New("args[1] is not config.LogConfig")
	}
	//(3) 将该服务设置的别名与配置放置到mp字段上
	p.mu.Lock() //上写锁
	p.mp[diName] = args[1]
	if len(p.mp) == 1 { //第一个别名设置为默认别名
		p.dn = diName
	}
	p.mu.Unlock()
	//(4)判断是否注入的时候就加载，类似阻塞的感觉，可以巧妙的控制启动资源的顺序额
	//todo 如何注入的时候是惰性的，那么该段核心程序会在调用端使用的时候执行
	if !lazy {
		_, err = setSingleton(diName, conf)
	}
	return
}

//注册过的别名
//@author sam@2020-07-30 11:42:01
func (p *provider) Provides() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return helper.MapToArray(p.mp)
}

//释放资源
//@author sam@2020-07-30 11:42:26
func (p *provider) Close() error {
	arr := p.Provides()
	for _, k := range arr {
		logger := getSingleton(k, false)
		if logger != nil {
			log, ok := logger.Out.(*os.File)
			if ok {
				log.Sync()
				log.Close()
			}
		}
	}
	return nil
}

//------------------------------------------ 与底层容器交互的几个公共方法 ------------------------------
//注入单例
//@todo 追加文件名
//@reviser sam@2020-04-02 14:30:28
func setSingleton(diName string, conf config.LogConfig) (ins *logrus.Logger, err error) {
	//log驱动核心程序
	ins, err = InitLog(conf)
	if err == nil {
		container.App.SetSingleton(diName, ins)
	}
	return
}

//获取单例
//@reviser sam@2020-04-02 14:30:12
func getSingleton(diName string, lazy bool) *logrus.Logger {
	//(1)先从容器中打捞一下
	rc := container.App.GetSingleton(diName)
	if rc != nil {
		return rc.(*logrus.Logger)
	}
	//(2)非惰性直接返回，即从容器中打捞就可以了
	if lazy == false {
		return nil
	}
	//(3)需要动态加载一下logger
	Pr.mu.RLock()                                //设置读锁
	conf, ok := Pr.mp[diName].(config.LogConfig) //读取要加载的资源的配置信息
	Pr.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("logger conf di_name:%s not exist", diName))
	}

	ins, err := setSingleton(diName, conf) //设置到容器中
	if err != nil {
		panic(fmt.Sprintf("logger di_name:%s err:%s", diName, err.Error()))
	}
	return ins
}

//------------------------------------  外部想使用该服务的方法--------------------------
//外部通过注入别名获取资源，解耦资源的关系
//@reviser sam@2020-04-02 14:20:47
func GetLogger(args ...string) *logrus.Logger {
	diName := helper.GetDiName(Pr.dn, args...)
	return getSingleton(diName, true)
}

