package rds

import (
	"errors"
	"fmt"
	"github.com/moka-mrp/sword-core/config"
	"github.com/moka-mrp/sword-core/helper"
	"github.com/moka-mrp/sword-core/kernel/container"
	"sync"
)

//------------------初始化---------------------------------
const (
	SingletonMain = "redis"
)

var Pr *provider

func init() {
	Pr = new(provider)
	Pr.mp = make(map[string]interface{})
}

//------------------------provider结构体 ------------------------------------------
type provider struct {
	mu sync.RWMutex
	mp map[string]interface{} //配置
	dn string                 //default name
}

/**
 * 注册资源
 * @param string 依赖注入别名 必选
 * @param config.LogConfig 配置 必选
 * @param bool 是否启用懒加载 可选
 * @author sam@2020-07-29 14:33:36
 */
func (p *provider) Register(args ...interface{}) (err error) {
	diName, lazy, err := helper.TransformArgs(args...)
	if err != nil {
		return
	}

	conf, ok := args[1].(config.RedisMultiConfig)
	if !ok {
		return errors.New("args[1] is not config.RedisConfig")
	}

	p.mu.Lock()
	p.mp[diName] = args[1]
	if len(p.mp) == 1 {
		p.dn = diName
	}
	p.mu.Unlock()

	if !lazy {
		_, err = setSingleton(diName, conf)
	}
	return
}

//注册过的别名
//@author sam@2020-07-29 14:36:45
func (p *provider) Provides() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return helper.MapToArray(p.mp)
}
//释放资源
//@author sam@2020-07-29 14:36:55
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

//-----------------------------------------------------------------------------------------------------------------

//注入单例
//@author sam@2020-04-10 10:47:52
func setSingleton(diName string, conf config.RedisMultiConfig) (ins *MultiPool, err error) {
	//核心代码，生成redis多连接池
	ins, err = NewMultiPool(conf)
	if err == nil {
		container.App.SetSingleton(diName, ins)
	}
	return
}

//获取单例
func getSingleton(diName string, lazy bool) *MultiPool {
	rc := container.App.GetSingleton(diName)
	if rc != nil {
		return rc.(*MultiPool)
	}
	if lazy == false {
		return nil
	}

	Pr.mu.RLock()
	conf, ok := Pr.mp[diName].(config.RedisMultiConfig)
	Pr.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("redis di_name:%s not exist", diName))
	}

	ins, err := setSingleton(diName, conf)
	if err != nil {
		panic(fmt.Sprintf("redis di_name:%s err:%s", diName, err.Error()))
	}
	return ins
}

//外部通过注入别名获取资源，解耦资源的关系
//@author sam@2020-07-29 14:32:39
func GetRedis(args ...string) *MultiPool {
	diName := helper.GetDiName(Pr.dn, args...)
	return getSingleton(diName, true)
}
