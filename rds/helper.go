package rds

import "strings"

/*
127.0.0.1:6379> set age 11
OK
127.0.0.1:6379>
我们通过读取每次redis的返回值来判断当前命令的执行对错
@todo  这里我们将每次redis执行命令的结果归结为true|false
@author sam@2020-04-09 18:03:34
*/
func isOKString(str string, err error) (bool, error) {
	if strings.ToUpper(str) == "OK" {
		return true, err
	}
	return false, err
}