package db

import (
	"errors"
	"fmt"
	"github.com/moka-mrp/sword-core/config"
	"github.com/moka-mrp/sword-core/helper"
	"github.com/moka-mrp/sword-core/kernel/container"
	"xorm.io/xorm"
	"sync"
)
//---------------------------- init ------------------------------------------------------------------------------------
const (
	SingletonMain = "db" //供bootstrap中使用，默认推荐的注入别名
)
var Pr *provider

func init() {
	Pr = new(provider)
	Pr.mp = make(map[string]interface{})
}

//------------------------------- provider struct ----------------------------------------------------------------------
type provider struct {
	mu sync.RWMutex
	mp map[string]interface{} //配置   别名1==>配置1   别名2==>配置2
	dn string                 //default name
}

/**
 * 注入资源 如 Pr.Register("db", dbConf, lazyBool)
 * @param string 依赖注入别名 必选
 * @param config.LogConfig 配置 必选
 * @param bool 是否启用懒加载 可选
 * @author sam@2020-07-29 09:41:30
 */
func (p *provider) Register(args ...interface{}) (err error) {
	//1.提取注入的别名以及是否惰性加载(本质就是提取第一和第三个参数)
	diName, lazy, err := helper.TransformArgs(args...)
	if err != nil {
		return
	}
    //2.对第二个参数进行类型断言，或者对应的配置信息
	conf, ok := args[1].(config.DbConfig)
	if !ok {
		return errors.New("args[1] is not config.DbConfig")
	}
	//3.别名绑定配置信息
	p.mu.Lock()
	p.mp[diName] = args[1]
	if len(p.mp) == 1 {
		p.dn = diName
	}
	p.mu.Unlock()
	//4.是否启动的时候就注入资源单例
	if !lazy {
		_, err = setSingleton(diName, conf)
	}
	return
}

//注册过的别名
//@author sam@2020-07-29 10:49:20
func (p *provider) Provides() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return helper.MapToArray(p.mp)
}

//释放资源
//@author sam@2020-07-29 10:49:49
func (p *provider) Close() error {
	arr := p.Provides()
	for _, k := range arr {
		c := getSingleton(k, false)
		if c != nil {
			c.Close()
		}
	}
	return nil
}

//------------------------------------------------- 单例 ---------------------------------------------------------------

//注入单例
//@author sam@2020-07-01 10:06:01
func setSingleton(diName string, conf config.DbConfig) (ins *xorm.EngineGroup, err error) {
	//核心代码，创建db实例
	ins, err = NewEngineGroup(conf) //db资源连接、释放
	//注入到容器中
	if err == nil {
		container.App.SetSingleton(diName, ins)
	}
	return
}
//获取单例
//@author sam@2020-07-01 14:46:48
func getSingleton(diName string, lazy bool) *xorm.EngineGroup {
	//从容器中打捞一下,打捞了就返回好了
	rc := container.App.GetSingleton(diName)
	if rc != nil {
		return rc.(*xorm.EngineGroup)
	}
	//如果未打捞到，但是惰性的也直接返回好了
	if lazy == false {
		return nil
	}
    //如果是未打捞到，非惰性的，我们当场现注入现取即可
	Pr.mu.RLock()
	conf, ok := Pr.mp[diName].(config.DbConfig)
	Pr.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("db di_name:%s not exist", diName))
	}

	ins, err := setSingleton(diName, conf)
	if err != nil {
		panic(fmt.Sprintf("db di_name:%s err:%s", diName, err.Error()))
	}
	return ins
}
//外部通过注入别名获取资源，解耦资源的关系
//@author sam@2020-07-01 14:45:28
func GetDb(args ...string) *xorm.EngineGroup {
	//获取注入的容器别名，未传递取Pr.dn默认别名
	diName := helper.GetDiName(Pr.dn, args...)
	return getSingleton(diName, true)
}
