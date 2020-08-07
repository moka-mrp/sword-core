package rds

import (
	"github.com/garyburd/redigo/redis"
)

//@todo  注意，下面的所有方法都是针对默认连接池封装的方法糖
//@link https://www.runoob.com/redis/redis-hyperloglog.html
//@author sam@2020-04-09 17:18:34


//---------------------公共系列---------------------------------------------
//判断是否存在某个key，一次传递多个的时候，只要有一个存在就返回1
//@author sam@2020-07-31 10:04:48
func (mp *MultiPool) Exists(key ...interface{}) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"EXISTS", key...))
}
//设置某个key的过期时间
//@author sam@2020-07-31 10:05:46
func (mp *MultiPool) Expire(key string, expiration int64) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"EXPIRE", key, expiration))
}
//设置某个key在哪个时间戳点过期
//@author sam@2020-07-31 10:06:21
func (mp *MultiPool) ExpireAt(key string, tm int64) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"EXPIREAT", key, tm))
}

//当 key 不存在时，返回 -2 。 当 key 存在但没有设置剩余生存时间时，返回 -1 。 否则，以秒为单位，返回 key 的剩余生存时间。
//@author sam@2020-07-31 10:09:44
func (mp *MultiPool) TTL(key string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"TTL", key))
}
//当 key 不存在时，返回 -2 。 当 key 存在但没有设置剩余生存时间时，返回 -1 。 否则，以毫秒为单位，返回 key 的剩余生存时间。
//@author sam@2020-07-31 10:11:39
func (mp *MultiPool) PTTL(key string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"PTTL", key))
}

//判断某个key的类型
//@author sam@2020-07-31 10:10:44
func (mp *MultiPool) Type(key string) (string, error) {
	return redis.String(mp.Do(DefaultPool,"TYPE", key))
}


//删除一个key
//@author  sam@2020-07-31 09:54:50
func (mp *MultiPool) Del(keys ...interface{}) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"DEL", keys...))
}

//发送原生命令,对服务器执行任何通用命令
//@author sam@2020-08-07 14:32:28
func (mp *MultiPool) RawCommand(command string, args ...interface{}) (interface{}, error) {
	return mp.Do(DefaultPool,command, args...)
}


//----------------String字符串类型操作 -----------------------------------------------------
//设置一个key
//@author sam@2020-04-09 16:48:38
func (mp *MultiPool) Set(key string, value interface{}) (bool, error) {
	//返回值的类型是interface{}
	//但本质其值是可以转成string ,直接类型断言成string不是很好，因为是两个返回参数，直接使用redis包提供的命令更加贴切
	return isOKString(redis.String(mp.Do(DefaultPool,"SET", key, value)))
}

//获取某个key当前的值，然后再赋值新的值
//@author sam@2020-07-31 10:01:24
func (mp *MultiPool) GetSet(key string, value interface{}) (string, error) {
	return redis.String(mp.Do(DefaultPool,"GETSET", key, value))
}

//mset  key1 value1   key2 value2 ...  一次设置多个key的值
//@author sam@2020-07-31 10:13:35
func (mp *MultiPool) MSet(pairs ...interface{}) (bool, error) {
	return isOKString(redis.String(mp.Do(DefaultPool,"MSET", pairs...)))
}

//Redis Msetnx 命令用于所有给定 key 都不存在时，同时设置一个或多个 key-value 对。
//当所有 key 都成功设置，返回 1 。 如果所有给定 key 都设置失败(至少有一个 key 已经存在)，那么返回 0 。
//@author sam@2020-07-31 10:15:00
func (mp *MultiPool) MSetNX(pairs ...interface{}) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"MSETNX", pairs...))
}


//get 一个key
//@author sam@2020-07-31 10:01:42
func (mp *MultiPool) Get(key string) (string, error) {
	//将返回的接口值转成字符串
	return redis.String(mp.Do(DefaultPool,"GET", key))
}

//mget   key1   key2  ... 一次获取多个key的值
//@author sam@2020-07-31 10:18:06
func (mp *MultiPool) MGet(keys ...interface{}) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"MGET", keys...))
}

// 对key的值做++操作，并返回新值
//@author sam@2020-07-31 10:22:55
func (mp *MultiPool) Incr(key string) (int, error) {
	return mp.IncrBy(key, 1)
}
// 每次加指定的值
//@author sam@2020-07-31 10:23:21
func (mp *MultiPool) IncrBy(key string, value int64) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"INCRBY", key, value))
}
//为 key 中所储存的值加上指定的浮点数增量值。
//如果 key 不存在，那么 INCRBYFLOAT 会先将 key 的值设为 0 ，再执行加法操作。
//@author sam@2020-07-31 10:24:54
func (mp *MultiPool) IncrByFloat(key string, value float64) (float64, error) {
	return redis.Float64(mp.Do(DefaultPool,"INCRBYFLOAT", key, value))
}

//对key的值做--操作，并返回新值
//@author sam@2020-07-31 10:25:30
func (mp *MultiPool) Decr(key string) (int, error) {
	return mp.DecrBy(key, 1)
}

//每次减指定的值
//@author  sam@2020-07-31 10:26:09
func (mp *MultiPool) DecrBy(key string, value int64) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"DECRBY", key, value))
}
//给指定的key的字符串追加value，就是拼接啦，类似php中的.
//@author sam@2020-07-31 10:26:48
func (mp *MultiPool) Append(key string, value string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"APPEND", key, value))
}


//------------------------List链表类型操作--------------------------------------------------------------------------
//在key对应list的头部（左边）添加字符串元素
//@author sam@2020-08-07 16:30:40
func (mp *MultiPool) LPush(key string, values ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, values...)
	return redis.Int(mp.Do(DefaultPool,"LPUSH", args...))
}
//在尾部（右边）添加
//@author sam@2020-08-07 16:31:03
func (mp *MultiPool) RPush(key string, values ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, values...)
	return redis.Int(mp.Do(DefaultPool,"RPUSH", args...))
}

//在头部（左边）删除
//@author sam@2020-08-07 16:31:25
func (mp *MultiPool) LPop(key string) (string, error) {
	return redis.String(mp.Do(DefaultPool,"LPOP", key))
}
//从key对应的list的尾部（右边）删除一个元素，并返回删除的元素
//@author sam@2020-08-07 16:31:47
func (mp *MultiPool) RPop(key string) (string, error) {
	return redis.String(mp.Do(DefaultPool,"RPOP", key))
}
//返回key对应list的长度，key不存在返回0，如果key对应类型不是list返回错误
//@author sam@2020-08-07 16:33:05
func (mp *MultiPool) LLen(key string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"LLEN", key))
}

//返回指定区间内的元素，下标从0开始
//@author  sam@2020-08-07 16:33:46
func (mp *MultiPool) LRange(key string, start, stop int64) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"LRANGE", key, start, stop))
}
//截取list,保留指定区间内元素,类似址操作
//@author sam@2020-08-07 16:35:30
func (mp *MultiPool) LTrim(key string, start, stop int64) (bool, error) {
	return isOKString(redis.String(mp.Do(DefaultPool,"LTRIM", key, start, stop)))
}

//-----------------------Set集合类型操作-----------------------------------------------------------------------------

//添加一个或多个string元素到key对应的set集合中，成功返回1
//如果元素已经在集合中，返回0，key对应的set不存在返回错误哦
//@author sam@2020-08-07 15:41:51
func (mp *MultiPool) SAdd(key string, members ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, members...)
	return redis.Int(mp.Do(DefaultPool,"SADD", args...))
}

//从key对应set中移除给定单个或者多个元素，成功返回1
//@author sam@2020-08-07 15:52:45
func (mp *MultiPool) SRem(key string, members ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, members...)
	return redis.Int(mp.Do(DefaultPool,"SREM", args...))
}

//从source对应set中移除member并添加到destination对应的set中
//@author sam@2020-08-07 15:59:52
func (mp *MultiPool) SMove(source, destination string, member interface{}) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"SMOVE", source, destination, member))
}

//返回set的元素个数
//@author sam@2020-08-07 15:59:04
func (mp *MultiPool) SCard(key string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"SCARD", key))
}

//判断member是否在set中
//@author sam@2020-08-07 16:01:52
func (mp *MultiPool) SIsMember(key string, member interface{}) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"SISMEMBER", key, member))
}


//返回所有给定key的交集  sinter     key1     key2  ...
//@author sam@2020-08-07 16:03:50
func (mp *MultiPool) SInter(keys ...interface{}) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"SINTER", keys...))
}

//返回所有给定key的并集 sunion    key1    key2 ...
//@author sam@2020-08-07 16:06:37
func (mp *MultiPool) SUnion(keys ...interface{}) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"SUNION", keys...))
}

//返回所有给定key的差集  sdiff        key1     key2 ...
//@author sam@2020-08-07 16:07:05
func (mp *MultiPool) SDiff(keys ...interface{}) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"SDIFF", keys...))
}

//返回key对应的set的所有元素，结果是无序的
//@author sam@2020-08-07 16:07:49
func (mp *MultiPool) SMembers(key string) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"SMEMBERS", key))
}

//-----------------------ZSet有序集合类型操作-------------------------------------------------------------------------

//添加元素到集合，元素在集合中存在则更新对应的权重值score
//zadd   key   score   member
//@author sam@2020-08-07 16:13:40
func (mp *MultiPool) ZAdd(key string, pairs ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, pairs...)
	return redis.Int(mp.Do(DefaultPool,"ZADD", args...))
}

//删除指定元素，1表示成功，如果元素不存在返回0
//zrem    key     member
//@author sam@2020-08-07 16:16:30
func (mp *MultiPool) ZRem(key string, members ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, members...)
	return redis.Int(mp.Do(DefaultPool,"ZREM", args...))
}

//返回集合中元素个数
//@author sam@2020-08-07 16:17:32
func (mp *MultiPool) ZCard(key string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"ZCARD", key))
}

//按照incr幅度增加对应member的权重值，返回新的权重值
//@author sam@2020-08-07 16:19:08
func (mp *MultiPool) ZIncrBy(key string, increment float64, member string) (float64, error) {
	return redis.Float64(mp.Do(DefaultPool,"ZINCRBY", key, increment, member))
}

//类似lrange操作从集合中去指定区间的元素。返回的是有序结果
//@author sam@2020-08-07 16:20:32
func (mp *MultiPool) ZRange(key string, start, stop int64) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"ZRANGE", key, start, stop))
}

//按照权重从大到小排元素
//@author sam@2020-08-07 16:47:34
func (rp *MultiPool) ZRevRange(key string, start, stop int64) ([]string, error) {
	return redis.Strings(rp.Do("ZREVRANGE", key, start, stop))
}

//返回指定元素在集合中的排名（下标） 第一个元素的下标是0 ...依次类推
//@author sam@2020-08-07 16:23:17
func (mp *MultiPool) ZRank(key, member string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"ZRANK", key, member))
}

//同上，但是集合中元素是按照权重从大到小排序
//@author sam@2020-08-07 16:26:17
func (mp *MultiPool) ZRevRank(key, member string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"ZREVRANK", key, member))
}


//返回给定元素对应的权重值score
//@author sam@2020-08-07 16:28:02
func (mp *MultiPool) ZScore(key, member string) (float64, error) {
	return redis.Float64(mp.Do(DefaultPool,"ZSCORE", key, member))
}

//删除集合中排名在给定区间的元素    默认是从小到大的，所以删除的是权重最小的的
//@author sam@2020-08-07 16:29:13
func (mp *MultiPool) ZRemRangeByRank(key string, start, stop int64) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"ZREMRANGEByRANK", key, start, stop))
}

//-----------------------Hash类型操作--------------------------------------------------------------------------------


//设置hash  field为指定值，如果key不存在，则先创建
//@author sam@2020-08-07 15:17:02
func (mp *MultiPool) HSet(key string, field string, value interface{}) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"HSET", key, field, value))
}

//hmset     key    field1  value1   field2   value2
//@author sam@2020-08-07 15:25:12
func (mp *MultiPool) HMSet(key string, pairs ...interface{}) (bool, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, pairs...)
	return isOKString(redis.String(mp.Do(DefaultPool,"HMSET", args...)))
}

//获取指定的hash field
//@author sam@2020-08-07 15:17:44
func (mp *MultiPool) HGet(key, field string) (string, error) {
	return redis.String(mp.Do(DefaultPool,"HGET", key, field))
}
//hmget     key   field1  field2 ...
//@author sam@2020-08-07 15:28:21
func (mp *MultiPool) HMGet(key string, fields ...interface{}) ([]interface{}, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, fields...)
	return redis.Values(mp.Do(DefaultPool,"HMGET", args...))
}

// 将指定的hash   field加上给定值
//@author sam@2020-08-07 15:14:51
func (mp *MultiPool) HIncrBy(key, field string, value int) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"HINCRBY", key, field, value))
}
//测试指定的field是否存在
//@author sam@2020-08-07 15:29:02
func (mp *MultiPool) HExists(key, field string) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"HEXISTS", key, field))
}

//删除指定的hash     field
//@author sam@2020-08-07 15:31:09
func (mp *MultiPool) HDel(key string, fields ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, fields...)
	return redis.Int(mp.Do(DefaultPool,"HDEL", args...))
}

//返回hash的所有field和value
//@author sam@2020-08-07 15:29:36
func (mp *MultiPool) HGetAll(key string) (map[string]string, error) {
	return redis.StringMap(mp.Do(DefaultPool,"HGETALL", key))
}

//返回hash的所有field
//@author sam@2020-08-07 15:30:12
func (mp *MultiPool) HKeys(key string) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"HKEYS", key))
}

//返回hash的所有value
//@author sam@2020-08-07 15:30:12
func (mp *MultiPool) HVals(key string) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"HVALS", key))
}

//返回指定的hash的field数量
//@author sam@2020-08-07 15:30:43
func (mp *MultiPool) HLen(key string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"HLEN", key))
}

//---------------------- 位操作---------------------------------------------------------------------------------------
func (mp *MultiPool) GetBit(key string, offset int64) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"GETBIT", key, offset))
}
func (mp *MultiPool) SetBit(key string, offset int64, value int) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"SETBIT", key, offset, value))
}

func (mp *MultiPool) BitCount(key string, offsets ...int64) (int, error) {
	switch len(offsets) {
	case 0:
		return redis.Int(mp.Do(DefaultPool,"BITCOUNT", key))
	case 2:
		return redis.Int(mp.Do(DefaultPool,"BITCOUNT", key, offsets[0], offsets[1]))
	default:
		return 0, errWrongArguments
	}
}

//---------------------  HyperLogLog  (基数统计)---------------------------

//添加指定元素到 HyperLogLog 中
//pfadd  nodes  1
//@author sam@2020-08-07 16:56:15
func (mp *MultiPool) PFAdd(key string, els ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, els...)
	return redis.Int(mp.Do(DefaultPool,"PFADD", args...))
}

//返回给定 HyperLogLog 的基数估算值
//@author sam@2020-08-07 16:56:38
func (mp *MultiPool) PFCount(keys ...interface{}) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"PFCOUNT", keys...))
}

//将多个 HyperLogLog 合并为一个 HyperLogLog
//@author sam@2020-08-07 16:56:54
func (mp *MultiPool) PFMerge(dest string, keys ...interface{}) (bool, error) {
	args := make([]interface{}, 0)
	args = append(args, dest)
	args = append(args, keys...)
	return isOKString(redis.String(mp.Do(DefaultPool,"PFMERGE", args...)))
}