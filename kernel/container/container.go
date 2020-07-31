package container

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var (
	ErrFactoryNotFound = errors.New("factory not found")
)

type factory = func() (interface{}, error)

// 容器
type Container struct {
	mu         sync.RWMutex
	singletons map[string]interface{} //资源别名->资源实例
	factories  map[string]factory
}

// 容器实例化
func NewContainer() *Container {
	return &Container{
		singletons: make(map[string]interface{}),
		factories:  make(map[string]factory),
	}
}

// 注册单例对象
//@author sam@2020-07-29 10:55:12
func (p *Container) SetSingleton(name string, singleton interface{}) {
	p.mu.Lock()
	p.singletons[name] = singleton
	p.mu.Unlock()
}

// 获取单例对象
//@author sam@2020-07-29 13:34:07
func (p *Container) GetSingleton(name string) interface{} {
	p.mu.RLock()
	ins, _ := p.singletons[name]
	p.mu.RUnlock()
	return ins
}

// 获取实例对象,并执行
//@author sam@2020-07-31 14:25:51
func (p *Container) GetPrototype(name string) (interface{}, error) {
	p.mu.RLock()
	factory, ok := p.factories[name]
	p.mu.RUnlock()
	if !ok {
		return nil, ErrFactoryNotFound
	}
	return factory()
}

// 设置实例对象工厂
// @author sam@2020-07-31 14:26:01
func (p *Container) SetPrototype(name string, factory factory) {
	p.mu.Lock()
	p.factories[name] = factory
	p.mu.Unlock()
}

//通过结构体的tag机制巧妙实现注入依赖
//todo kind是分类，如struct，map... | Type是具体的类型，如Student
//todo 依赖注入功能待后期完善
//@author sam@2020-07-31 14:20:27
func (p *Container) Ensure(instance interface{}) error {

	//通过反射获取的传入的变量的 type , kind, 值
	//以防传递的是指针,则通过Elem()方法获取本身的元素
	elemType := reflect.TypeOf(instance).Elem()
	ele := reflect.ValueOf(instance).Elem()

	//变量实例的每个属性字段，分析tag情况
	for i := 0; i < elemType.NumField(); i++ { // 遍历字段
		//判断是否有tag-di标识
		fieldType := elemType.Field(i)
		tag := fieldType.Tag.Get("di") // 获取tag
		diName := p.injectName(tag)
		if diName == "" {
			continue
		}
		//fmt.Println(tag) //rick,tom,sam
		//fmt.Println(diName) //rick
		var (
			diInstance interface{}
			err        error
		)
		if p.isSingleton(tag) { //注入的依赖是否在容器池中
			diInstance = p.GetSingleton(diName)
		}
		if p.isPrototype(tag) { //tag中含有prototype字样就是实例依赖
			diInstance, err = p.GetPrototype(diName)
		}
		if err != nil {
			return err
		}
		if diInstance == nil {
			return errors.New(diName + " dependency not found")
		}
		ele.Field(i).Set(reflect.ValueOf(diInstance))
	}
	return nil
}

// 获取需要注入的依赖名称
func (p *Container) injectName(tag string) string {
	tags := strings.Split(tag, ",")
	if len(tags) == 0 {
		return ""
	}
	return tags[0]
}

// 检测是否单例依赖
// tag字符串中含有prototype字样就不是单例依赖了
//@author  sam@2020-07-31 11:44:47
func (p *Container) isSingleton(tag string) bool {
	tags := strings.Split(tag, ",")
	for _, name := range tags {
		if name == "prototype" {
			return false
		}
	}
	return true
}

// 检测是否实例依赖
// tag中含有prototype字样就是实例依赖
// @author sam@2020-07-31 11:46:08
func (p *Container) isPrototype(tag string) bool {
	tags := strings.Split(tag, ",")
	for _, name := range tags {
		if name == "prototype" {
			return true
		}
	}
	return false
}

// 打印容器内部实例
//singletons:
//di1: 0xc0000524f0 string
//di2: 0xc0000524f0 string
//factories:
//@author sam@2020-07-31 11:28:25
func (p *Container) String() string {
	//将要输出字符先临时存入切片字符串中
	lines := make([]string, 0, len(p.singletons)+len(p.factories)+2)

	//singletons:
	lines = append(lines, "singletons:")
	for name, item := range p.singletons {
		if item == nil { //item为nil
			line := fmt.Sprintf("  %s: %s %s", name, "<nil>", "<nil>")
			lines = append(lines, line)
			continue
		}

		line := fmt.Sprintf("  %s: %p %s", name, &item, reflect.TypeOf(item).String())
		lines = append(lines, line)
	}
	//factories:
	lines = append(lines, "factories:")
	for name, item := range p.factories {
		if item == nil {
			line := fmt.Sprintf("  %s: %s %s", name, "<nil>", "<nil>")
			lines = append(lines, line)
			continue
		}

		line := fmt.Sprintf("  %s: %p %s", name, &item, reflect.TypeOf(item).String())
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
