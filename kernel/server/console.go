package server

import (
	"github.com/moka-mrp/sword-core/kernel/event"
	"github.com/robfig/cron"
)

var DefaultCronEngine *cron.Cron
func ExitCron(i interface{}){
	DefaultCronEngine.Stop()
}

//开启定时任务
//@author sam@2020-12-04 09:53:22
func StartConsole(registerSchedule func(*cron.Cron)) error {
	//注册Cron执行计划
	DefaultCronEngine= cron.New()
	registerSchedule(DefaultCronEngine)
	DefaultCronEngine.Start()
	//守护进程
	event.On(event.EXIT,ExitCron)
	event.Wait()
	event.Emit(event.EXIT, nil)
	return nil
}
