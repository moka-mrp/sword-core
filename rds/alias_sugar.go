package rds
import (
	"github.com/garyburd/redigo/redis"
)
//@todo  注意，下面的所有方法都是针对默认连接池封装的方法糖
//@author sam@2020-04-09 17:18:34


//----------------set and  get 系列 -----------------------------------------------------
func (mp *MultiPool) AliasGet(poolName string,key string) (string, error) {
	//将返回的接口值转成字符串
	return redis.String(mp.Do(poolName,"GET", key))
}

//设置一个key
//@author sam@2020-04-09 16:48:38
func (mp *MultiPool) AliasSet(poolName string,key string, value interface{}) (bool, error) {
	//返回值的类型是interface{}
	//但本质其值是可以转成string ,直接类型断言成string不是很好，因为是两个返回参数，直接使用redis包提供的命令更加贴切
	return isOKString(redis.String(mp.Do(poolName,"SET", key, value)))
}

//删除一个key
//@author  sam@2020-07-31 09:52:59
func (mp *MultiPool) AliasDel(poolName string,keys ...interface{}) (int, error) {
	return redis.Int(mp.Do(poolName,"DEL", keys...))
}


//------------------------------------------------------------------------------------------




