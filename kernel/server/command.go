package server

import (
	"github.com/moka-mrp/sword-core/command"
)

//注册并执行某个name对应的脚本
//@reviser sam@2020-04-14 13:58:56
func ExecuteCommand(name string, registerCommand func(*command.Command)) error {
	//创建一个命令行容器
	c := command.New()
	//注入要定义的命令
	registerCommand(c)
	//执行某个具体的命令
	err := c.Execute(name)
	return err
}
