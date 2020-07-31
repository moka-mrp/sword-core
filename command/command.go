package command

import (
	"errors"
	"sync"
)
//本质是通过别名注入一堆方法，后面可以通过别名指明调用某个方法
//@reviser sam@2020-04-14 11:49:34

var (
	ErrUnknownName = errors.New("The current command does not exist")
)

//--------------------------------- 一次性任务脚本存储结构体 -----------------------------------------
type Command struct {
	mu        sync.RWMutex
	container map[string]func()
}


//绑定name与函数的关系
func (c *Command) AddFunc(name string, f func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.container[name] = f
}

//通过name执行函数
//@reviser sam@2020-04-14 14:00:21
func (c *Command) Execute(name string) (err error) {
	c.mu.RLock()
	f, ok := c.container[name]
	c.mu.RUnlock()
	if ok {
		f()
	} else {
		return  ErrUnknownName
	}
	return
}

//----------------------------------------------------------------------------
//new实例
func New() *Command {
	c := new(Command)
	c.container = make(map[string]func())
	return c
}
