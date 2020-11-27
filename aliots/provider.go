package aliots

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
	SingletonMain = "ots"
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
 * 注册资源,不管是惰性还是非惰性，都要将对应的资源配置注入容器中
 * 非惰性则注入的时候就创建实例，惰性的话就在获取单例的时候现场注入资源即可
 * @param string 依赖注入别名 必选
 * @param config.LogConfig 配置 必选
 * @param bool 是否启用懒加载 可选
 * @author sam@2020-11-27 15:11:49
 */
func (p *provider) Register(args ...interface{}) (err error) {
	diName, lazy, err := helper.TransformArgs(args...)
	if err != nil {
		return
	}

	conf, ok := args[1].(config.OtsMultiConfig)
	if !ok {
		return errors.New("args[1] is not config.OtsConfig")
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
//@author sam@2020-11-27 16:04:21
func (p *provider) Close() error {
	return nil
}

//-----------------------------------------------------------------------------------------------------------------

//注入单例
//@author sam@2020-11-27 15:57:02
func setSingleton(diName string, conf config.OtsMultiConfig) (ins *MultiClient, err error) {
	//核心代码，生成ots多客户端
	ins, err = NewMultiClient(conf)
	if err == nil {
		container.App.SetSingleton(diName, ins)
	}
	return
}

//获取单例
func getSingleton(diName string, lazy bool) *MultiClient {
	rc := container.App.GetSingleton(diName)
	if rc != nil {
		return rc.(*MultiClient)
	}
	if lazy == false {
		return nil
	}

	Pr.mu.RLock()
	conf, ok := Pr.mp[diName].(config.OtsMultiConfig)
	Pr.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("ots di_name:%s not exist", diName))
	}

	ins, err := setSingleton(diName, conf)
	if err != nil {
		panic(fmt.Sprintf("ots di_name:%s err:%s", diName, err.Error()))
	}
	return ins
}

//外部通过注入别名获取资源，解耦资源的关系
//@author sam@2020-11-27 15:59:16
func GetOts(args ...string) *MultiClient {
	diName := helper.GetDiName(Pr.dn, args...)
	return getSingleton(diName, true)
}
