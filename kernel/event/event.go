package event

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

//Golang信号处理和优雅退出守护进程
const (
	EXIT = "exit"
	WAIT = "wait"
)

var (
	Events = make(map[string][]func(interface{}), 2)
)

//分门别类的进行一组回调函数的注入
//@author sam@2020-09-08 16:43:39
func On(name string, fs ...func(interface{})) error {
	evs, ok := Events[name]
	if !ok { //如果里面一个元素都不存在的时候需要初始化
		evs = make([]func(interface{}), 0, len(fs))
	}
	//遍历追加的函数，分别注入到Events中
	for _, f := range fs {
		if f == nil {
			continue
		}
		fp := reflect.ValueOf(f).Pointer()
		for i := 0; i < len(evs); i++ {
			if reflect.ValueOf(evs[i]).Pointer() == fp {
				return fmt.Errorf("func[%v] already exists in event[%s]", fp, name)
			}
		}
		evs = append(evs, f)
	}
	//追加
	Events[name] = evs
	return nil
}

//指明回调某组回调函数，可以统一传递参数
//todo  Emit(EXIT, nil)
//todo  Emit(WAIT, nil)
//@author sam@2020-09-08 16:50:06
func Emit(name string, arg interface{}) {
	evs, ok := Events[name]
	if !ok {
		return
	}

	for _, f := range evs {
		f(arg)
	}
}

func EmitAll(arg interface{}) {
	for _, fs := range Events {
		for _, f := range fs {
			f(arg)
		}
	}
	return
}

//从Events中移除某组注入的某个回调函数
//@author sam@2020-09-08 17:05:58
func Off(name string, f func(interface{})) error {
	evs, ok := Events[name]
	if !ok || len(evs) == 0 {
		return fmt.Errorf("envet[%s] doesn't have any funcs", name)
	}
	fp := reflect.ValueOf(f).Pointer()
	for i := 0; i < len(evs); i++ {
		if reflect.ValueOf(evs[i]).Pointer() == fp {
			evs = append(evs[:i], evs[i+1:]...)
			Events[name] = evs
			return nil
		}
	}

	return fmt.Errorf("%v func dones't exist in event[%s]", fp, name)
}

//清空所有的回调函数
//@author  sam@2020-09-08 17:07:29
func OffAll(name string) error {
	Events[name] = nil
	return nil
}

//监听信号，未指明则默认监控 syscall.SIGINT, syscall.SIGTERM 两个信号
//SIGINT	2	Term	中断信号:用户发送INTR字符(Ctrl+C)触发
//SIGTERM	15	Term	结束程序(可以被捕获、阻塞或忽略)
//@author sam@2020-09-08 17:09:12
func Wait(sig ...os.Signal) os.Signal {
	//创建chan
	c := make(chan os.Signal, 1)
	if len(sig) == 0 {
		signal.Notify(c, syscall.SIGINT,syscall.SIGTERM)
	} else {
		signal.Notify(c, sig...)
	}
	return <-c
}
