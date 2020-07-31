package helper

import (
	"errors"
)

//获取依赖中注入的名字  默认的名字
//todo  GetDb("db1")  ==>helper.GetDiName("db","db1") ===>最终结果指明别名
//todo  GetDb("db")  ==>helper.GetDiName("db","db") ==>最终结果取指明别名，这里恰好指明别名与默认别名相同罢了
//todo  GetDb()  ==>helper.GetDiName("db") ==>最终结果取默认别名
//@reviser sam@2020-04-02 14:10:54
func GetDiName(defaultName string, args ...string) string {
	var name string
	if len(args) > 0 {
		name = args[0]
	}
	if name == "" {
		return defaultName
	}
	return name
}


/**
 * 简易版的参数检测
 * @param string 依赖注入别名 必选
 * @param config.LogConfig 配置 必选
 * @param bool 是否启用懒加载 可选
 * @todo 提取 diName 和lazy值
 * @reviser sam@2020-04-02 14:48:51
 */
func TransformArgs(args ...interface{}) (diName string, lazy bool, err error) {
	//前两个参数必传
	if len(args) < 2 {
		err = errors.New("args is not enough")
		return
	}
    //提取第一个参数diName,必须是可以转成string的
	var ok bool
	diName, ok = args[0].(string)
	if !ok {
		err = errors.New("args[0] is not string")
		return
	}
    //提取第三个参数lazy
	if len(args) > 2 {
		lazy, _ = args[2].(bool)
	}
	return
}

/**
 *map中的key转字符arr
 *@author sam@2020-07-29 10:20:33
 */
func MapToArray(mp map[string]interface{}) []string {
	arr := make([]string, len(mp))
	i := 0
	for k := range mp {
		arr[i] = k
		i++
	}
	return arr
}
