package rds

import (
	"github.com/garyburd/redigo/redis"
)

//@todo  注意，下面的所有方法都是针对默认连接池封装的方法糖
//@link https://www.runoob.com/redis/redis-hyperloglog.html
//@author sam@2020-04-09 17:18:34


//---------------------公共系列---------------------------------------------

//1.删除一个key
//@author  sam@2020-07-31 09:54:50
func (mp *MultiPool) Del(keys ...interface{}) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"DEL", keys...))
}

//2.判断是否存在某个key，一次传递多个的时候，只要有一个存在就返回1
//@author sam@2020-07-31 10:04:48
func (mp *MultiPool) Exists(key ...interface{}) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"EXISTS", key...))
}
//3.设置某个key的过期时间
//@author sam@2020-07-31 10:05:46
func (mp *MultiPool) Expire(key string, expiration int64) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"EXPIRE", key, expiration))
}
//4.设置某个key在哪个时间戳点过期
//@author sam@2020-07-31 10:06:21
func (mp *MultiPool) ExpireAt(key string, tm int64) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"EXPIREAT", key, tm))
}


//5.当 key 不存在时，返回 -2 。 当 key 存在但没有设置剩余生存时间时，返回 -1 。 否则，以毫秒为单位，返回 key 的剩余生存时间。
//@author sam@2020-07-31 10:11:39
func (mp *MultiPool) PTTL(key string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"PTTL", key))
}


//6.当 key 不存在时，返回 -2 。 当 key 存在但没有设置剩余生存时间时，返回 -1 。 否则，以秒为单位，返回 key 的剩余生存时间。
//@author sam@2020-07-31 10:09:44
func (mp *MultiPool) TTL(key string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"TTL", key))
}


//7.发送原生命令,对服务器执行任何通用命令
//@author sam@2020-08-07 14:32:28
func (mp *MultiPool) RawCommand(command string, args ...interface{}) (interface{}, error) {
	return mp.Do(DefaultPool,command, args...)
}


//8.判断某个key的类型
//@author sam@2020-07-31 10:10:44
func (mp *MultiPool) Type(key string) (string, error) {
	return redis.String(mp.Do(DefaultPool,"TYPE", key))
}






//----------------String字符串类型操作 -----------------------------------------------------
//1.设置指定 key 的值
//@author sam@2020-04-09 16:48:38
func (mp *MultiPool) Set(key string, value interface{}) (bool, error) {
	//返回值的类型是interface{}
	//但本质其值是可以转成string ,直接类型断言成string不是很好，因为是两个返回参数，直接使用redis包提供的命令更加贴切
	return isOKString(redis.String(mp.Do(DefaultPool,"SET", key, value)))
}

//2.获取指定 key 的值。
//@author sam@2020-07-31 10:01:42
func (mp *MultiPool) Get(key string) (string, error) {
	//将返回的接口值转成字符串
	return redis.String(mp.Do(DefaultPool,"GET", key))
}

//3.返回 key 中字符串值的子字符
//@author sam@2020-08-07 17:04:35
func (mp *MultiPool) GetRange(key string, start, end int64) (string, error) {
	return redis.String(mp.Do(DefaultPool,"GETRANGE", key, start, end))
}

//4.将给定 key 的值设为 value ，并返回 key 的旧值(old value)
//@author sam@2020-07-31 10:01:24
func (mp *MultiPool) GetSet(key string, value interface{}) (string, error) {
	return redis.String(mp.Do(DefaultPool,"GETSET", key, value))
}
//5.对 key 所储存的字符串值，获取指定偏移量上的位(bit)
//@author sam@2020-08-07 17:06:22
func (mp *MultiPool) GetBit(key string, offset int64) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"GETBIT", key, offset))
}


//6.获取所有(一个或多个)给定 key 的值。
//@author sam@2020-07-31 10:18:06
func (mp *MultiPool) MGet(keys ...interface{}) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"MGET", keys...))
}

//7.对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)
func (mp *MultiPool) SetBit(key string, offset int64, value int) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"SETBIT", key, offset, value))
}

//8.将值 value 关联到 key ，并将 key 的过期时间设为 seconds (以秒为单位)。
//@author sam@2020-08-07 17:11:54
func (mp *MultiPool) SetEX(key string, value interface{}, seconds int64) (bool, error) {
	return isOKString(redis.String(mp.Do(DefaultPool,"SET", key, value, "EX", seconds)))
	//return isOKString(redis.String(mp.Do(DefaultPool,"SETEX", key,seconds,value)))
}

//9.只有在 key 不存在时设置 key 的值(分布式锁常用)
//@author sam@2020-08-07 17:14:35
func (mp *MultiPool) SetNX(key string, value interface{}) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"SETNX", key, value))
}
//10.用 value 参数覆写给定 key 所储存的字符串值，从偏移量 offset 开始。
//@author sam@2020-08-07 17:17:40
func (mp *MultiPool) SetRange(key string, offset int64, value string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"SETRANGE", key, offset, value))
}

//11.返回 key 所储存的字符串值的长度。
//@author sam@2020-08-07 17:18:18
func (mp *MultiPool) StrLen(key string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"STRLEN", key))
}

//12.同时设置一个或多个 key-value 对。
//@author sam@2020-07-31 10:13:35
func (mp *MultiPool) MSet(pairs ...interface{}) (bool, error) {
	return isOKString(redis.String(mp.Do(DefaultPool,"MSET", pairs...)))
}

//13.同时设置一个或多个 key-value 对，当且仅当所有给定 key 都不存在
//@author sam@2020-07-31 10:15:00
func (mp *MultiPool) MSetNX(pairs ...interface{}) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"MSETNX", pairs...))
}

//14.这个命令和 SETEX 命令相似，但它以毫秒为单位设置 key 的生存时间，而不是像 SETEX 命令那样，以秒为单位。

func (mp *MultiPool) SetPX(key string, value interface{}, milliseconds int64) (bool, error) {
	return isOKString(redis.String(mp.Do(DefaultPool,"SET", key, value, "PX", milliseconds)))
	//return isOKString(redis.String(mp.Do(DefaultPool,"PSETEX", key,milliseconds,value)))
}

//15.将 key 中储存的数字值增一。
//@author sam@2020-07-31 10:22:55
func (mp *MultiPool) Incr(key string) (int, error) {
	return mp.IncrBy(key, 1)
}
//16.将 key 所储存的值加上给定的增量值（increment）
//@author sam@2020-07-31 10:23:21
func (mp *MultiPool) IncrBy(key string, value int64) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"INCRBY", key, value))
}
//17.将 key 所储存的值加上给定的浮点增量值（increment）
//@author sam@2020-07-31 10:24:54
func (mp *MultiPool) IncrByFloat(key string, value float64) (float64, error) {
	return redis.Float64(mp.Do(DefaultPool,"INCRBYFLOAT", key, value))
}

//18.将 key 中储存的数字值减一
//@author sam@2020-07-31 10:25:30
func (mp *MultiPool) Decr(key string) (int, error) {
	return mp.DecrBy(key, 1)
}

//19.key 所储存的值减去给定的减量值（decrement）
//@author  sam@2020-07-31 10:26:09
func (mp *MultiPool) DecrBy(key string, value int64) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"DECRBY", key, value))
}
//20.如果 key 已经存在并且是一个字符串， APPEND 命令将指定的 value 追加到该 key 原来值（value）的末尾。
//@author sam@2020-07-31 10:26:48
func (mp *MultiPool) Append(key string, value string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"APPEND", key, value))
}

//21.统计字符串被设置为1的bit数
//@author sam@2020-08-08 11:09:12
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


//------------------------List链表类型操作--------------------------------------------------------------------------
//一般认为左头右尾

//1.移出并获取列表的第一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
//@author sam@2020-08-07 17:55:29
func (mp *MultiPool) BLPop(args ...interface{}) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"BLPOP", args...))
}

//2.移出并获取列表的最后一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
//@author sam@2020-08-07 20:18:06
func (mp *MultiPool) BRPop(args ...interface{}) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"BRPOP", args...))
}


//3.从列表中弹出一个值，将弹出的元素插入到另外一个列表中并返回它； 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
//@author sam@2020-08-07 20:30:57
func (mp *MultiPool) BRPopLPush(source, destination string, timeout uint64) (string, error) {
	return redis.String(mp.Do(DefaultPool,"BRPOPLPUSH", source, destination, timeout))
}

//4.通过索引获取列表中的元素
//@author sam@2020-08-07 20:31:09
func (mp *MultiPool) LIndex(key string, index int64) (string, error) {
	return redis.String(mp.Do(DefaultPool,"LINDEX", key, index))
}

//5.在列表的元素前或者后插入元素
//@author  sam@2020-08-07 20:32:32
func (mp *MultiPool) LInsert(key string, op string, pivot, value interface{}) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"LINSERT", key, op, pivot, value))
}

//6.获取列表长度
//@author sam@2020-08-07 16:33:05
func (mp *MultiPool) LLen(key string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"LLEN", key))
}

//7.移出并获取列表的第一个元素
//@author sam@2020-08-07 16:31:25
func (mp *MultiPool) LPop(key string) (string, error) {
	return redis.String(mp.Do(DefaultPool,"LPOP", key))
}

//8.将一个或多个值插入到列表头部
//@author sam@2020-08-07 16:30:40
func (mp *MultiPool) LPush(key string, values ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, values...)
	return redis.Int(mp.Do(DefaultPool,"LPUSH", args...))
}

//9.将一个值插入到已存在的列表头部
//@author sam@2020-08-07 20:42:43
func (mp *MultiPool) LPushX(key string, value interface{}) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"LPUSHX", key, value))
}

//10.获取列表指定范围内的元素
//@author  sam@2020-08-07 16:33:46
func (mp *MultiPool) LRange(key string, start, stop int64) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"LRANGE", key, start, stop))
}


//11.移除列表元素
func (mp *MultiPool) LRem(key string, count int64, value interface{}) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"LREM", key, count, value))
}

//12.通过索引设置列表元素的值
func (mp *MultiPool) LSet(key string, index int64, value interface{}) (bool, error) {
	return isOKString(redis.String(mp.Do(DefaultPool,"LSET", key, index, value)))
}

//13.对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
//@author sam@2020-08-07 16:35:30
func (mp *MultiPool) LTrim(key string, start, stop int64) (bool, error) {
	return isOKString(redis.String(mp.Do(DefaultPool,"LTRIM", key, start, stop)))
}

//14.移除列表的最后一个元素，返回值为移除的元素
//@author sam@2020-08-07 16:31:47
func (mp *MultiPool) RPop(key string) (string, error) {
	return redis.String(mp.Do(DefaultPool,"RPOP", key))
}

//15.移除列表的最后一个元素，并将该元素添加到另一个列表并返回
//@author sam@2020-08-07 20:57:15
func (mp *MultiPool) RPopLPush(source, destination string) (string, error) {
	return redis.String(mp.Do(DefaultPool,"RPOPLPUSH", source, destination))
}

//16.在列表中添加一个或多个值,在尾部（右边）添加
//@author sam@2020-08-07 16:31:03
func (mp *MultiPool) RPush(key string, values ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, values...)
	return redis.Int(mp.Do(DefaultPool,"RPUSH", args...))
}


//17.为已存在的列表添加值
func (mp *MultiPool) RPushX(key string, value interface{}) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"RPUSHX", key, value))
}


//-----------------------Set集合类型操作-----------------------------------------------------------------------------

//1.向集合添加一个或多个成员
//@author sam@2020-08-07 15:41:51
func (mp *MultiPool) SAdd(key string, members ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, members...)
	return redis.Int(mp.Do(DefaultPool,"SADD", args...))
}

//2.获取集合的成员数
//@author sam@2020-08-07 15:59:04
func (mp *MultiPool) SCard(key string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"SCARD", key))
}

//3.返回第一个集合与其他集合之间的差异。
//@author sam@2020-08-07 16:07:05
func (mp *MultiPool) SDiff(keys ...interface{}) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"SDIFF", keys...))
}

//4.返回给定所有集合的差集并存储在 destination 中
//@author sam@2020-08-08 09:44:20
func (mp *MultiPool) SDiffStore(destination string, keys ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, destination)
	args = append(args, keys...)
	return redis.Int(mp.Do(DefaultPool,"SDIFFSTORE", args...))
}

//5.返回给定所有集合的交集
//@author sam@2020-08-07 16:03:50
func (mp *MultiPool) SInter(keys ...interface{}) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"SINTER", keys...))
}

//6.返回给定所有集合的交集并存储在 destination 中
//@author sam@2020-08-08 09:46:56
func (mp *MultiPool) SInterStore(destination string, keys ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, destination)
	args = append(args, keys...)
	return redis.Int(mp.Do(DefaultPool,"SINTERSTORE", args...))
}


//7.判断 member 元素是否是集合 key 的成员
//@author sam@2020-08-07 16:01:52
func (mp *MultiPool) SIsMember(key string, member interface{}) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"SISMEMBER", key, member))
}

//8.返回集合中的所有成员
//@author sam@2020-08-07 16:07:49
func (mp *MultiPool) SMembers(key string) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"SMEMBERS", key))
}

//9.将 member 元素从 source 集合移动到 destination 集合
//@author sam@2020-08-07 15:59:52
func (mp *MultiPool) SMove(source, destination string, member interface{}) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"SMOVE", source, destination, member))
}

//10.移除并返回集合中的一个随机元素
//@author sam@2020-08-08 09:57:19
func (mp *MultiPool) SPop(key string) (string, error) {
	return redis.String(mp.Do(DefaultPool,"SPOP", key))
}

//11.返回集合中一个或多个随机数
//@author  sam@2020-08-08 10:00:56
func (mp *MultiPool) SRandMember(key string) (string, error) {
	return redis.String(mp.Do(DefaultPool,"SRANDMEMBER", key))
}

//12.移除集合中一个或多个成员
//@author sam@2020-08-07 15:52:45
func (mp *MultiPool) SRem(key string, members ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, members...)
	return redis.Int(mp.Do(DefaultPool,"SREM", args...))
}
//13.返回所有给定集合的并集
//@author sam@2020-08-07 16:06:37
func (mp *MultiPool) SUnion(keys ...interface{}) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"SUNION", keys...))
}

//14.所有给定集合的并集存储在 destination 集合中
func (mp *MultiPool) SUnionStore(destionation string, keys ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, destionation)
	args = append(args, keys...)
	return redis.Int(mp.Do(DefaultPool,"SUNIONSTORE", args...))
}

//15.迭代集合中的元素: SSCAN key cursor [MATCH pattern] [COUNT count]




//-----------------------ZSet有序集合类型操作-------------------------------------------------------------------------

//1.向有序集合添加一个或多个成员，或者更新已存在成员的分数
//@author sam@2020-08-07 16:13:40
func (mp *MultiPool) ZAdd(key string, pairs ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, pairs...)
	return redis.Int(mp.Do(DefaultPool,"ZADD", args...))
}

//2.获取有序集合的成员数
//@author sam@2020-08-07 16:17:32
func (mp *MultiPool) ZCard(key string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"ZCARD", key))
}

//3.计算在有序集合中指定区间分数的成员数
func (mp *MultiPool) ZCount(key string, min, max interface{}) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"ZCOUNT", key, min, max))
}

//4.有序集合中对指定成员的分数加上增量 increment
//@author sam@2020-08-07 16:19:08
func (mp *MultiPool) ZIncrBy(key string, increment float64, member string) (float64, error) {
	return redis.Float64(mp.Do(DefaultPool,"ZINCRBY", key, increment, member))
}

//5.计算给定的一个或多个有序集的交集并将结果集存储在新的有序集合 key 中
//@author sam@2020-08-08 10:19:02
func (mp *MultiPool) ZInterStore(destination string, nkeys int, params ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, destination)
	args = append(args, nkeys)
	args = append(args, params...)
	return redis.Int(mp.Do(DefaultPool,"ZINTERSTORE", args...))
}

//6.在有序集合中计算指定字典区间内成员数量
//@author sam@2020-08-08 10:23:37
func (mp *MultiPool) ZLexCount(key, min, max string) (int, error) {
	return redis.Int(mp.Do("DefaultPool", key, min, max))
}

//7.通过索引区间返回有序集合指定区间内的成员
//@author sam@2020-08-07 16:20:32
func (mp *MultiPool) ZRange(key string, start, stop int64) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"ZRANGE", key, start, stop))
}



//8.通过字典区间返回有序集合的成员
//@author sam@2020-08-08 10:41:55
func (mp *MultiPool) ZRangeByLex(key, min, max string, offset, count int64) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"ZRANGEBYLEX", key, min, max, "LIMIT", offset, count))
}

//9.通过分数返回有序集合指定区间内的成员
//@author sam@2020-08-08 10:45:41
func (mp *MultiPool) ZRangeByScore(key, min, max interface{}, offset, count int64) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"ZRANGEBYSCORE", key, min, max, "LIMIT", offset, count))
}


//10.返回有序集合中指定成员的索引
//@author sam@2020-08-07 16:23:17
func (mp *MultiPool) ZRank(key, member string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"ZRANK", key, member))
}

//11.移除有序集合中的一个或多个成员
//@author sam@2020-08-07 16:16:30
func (mp *MultiPool) ZRem(key string, members ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, members...)
	return redis.Int(mp.Do(DefaultPool,"ZREM", args...))
}


//12.移除有序集合中给定的字典区间的所有成员
//@author  sam@2020-08-08 10:51:31
func (mp *MultiPool) ZRemRangeByLex(key, min, max string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"ZREMRANGEBYLEX", key, min, max))
}
//13.移除有序集合中给定的排名区间的所有成员
//@author sam@2020-08-08 10:52:54
func (mp *MultiPool) ZRemRangeByRank(key string, start, stop int64) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"ZREMRANGEByRANK", key, start, stop))
}

//14.移除有序集合中给定的分数区间的所有成员
func (mp *MultiPool) ZRemRangeByScore(key, min, max interface{}) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"ZREMRANGEBYSCORE", key, min, max))
}

//15.返回有序集中指定区间内的成员，通过索引，分数从高到低
//@author sam@2020-08-07 16:47:34
func (mp *MultiPool) ZRevRange(key string, start, stop int64) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"ZREVRANGE", key, start, stop))
}


//16.返回有序集中指定分数区间内的成员，分数从高到低排序
//@author sam@2020-08-08 10:57:14
func (mp *MultiPool) ZRevRangeByScore(key, max, min interface{}, offset, count int64) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"ZREVRANGEBYSCORE", key, max, min, "LIMIT", offset, count))
}


//17.返回有序集合中指定成员的排名，有序集成员按分数值递减(从大到小)排序
//@author sam@2020-08-07 16:26:17
func (mp *MultiPool) ZRevRank(key, member string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"ZREVRANK", key, member))
}

//18.返回有序集中，成员的分数值
//@author sam@2020-08-07 16:28:02
func (mp *MultiPool) ZScore(key, member string) (float64, error) {
	return redis.Float64(mp.Do(DefaultPool,"ZSCORE", key, member))
}

//19.计算给定的一个或多个有序集的并集，并存储在新的 key 中
//@author sam@2020-08-08 10:59:09
func (mp *MultiPool) ZUnionStore(destination string, nkeys int, params ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, destination)
	args = append(args, nkeys)
	args = append(args, params...)
	return redis.Int(mp.Do(DefaultPool,"ZUNIONSTORE", args...))
}


//20	ZSCAN key cursor [MATCH pattern] [COUNT count]
//迭代有序集合中的元素（包括元素成员和元素分值）

//-----------------------Hash类型操作--------------------------------------------------------------------------------

//1.删除一个或多个哈希表字段
//@author sam@2020-08-07 15:31:09
func (mp *MultiPool) HDel(key string, fields ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, fields...)
	return redis.Int(mp.Do(DefaultPool,"HDEL", args...))
}

//2.查看哈希表 key 中，指定的字段是否存在。
//@author sam@2020-08-07 15:29:02
func (mp *MultiPool) HExists(key, field string) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"HEXISTS", key, field))
}

//3.获取存储在哈希表中指定字段的值。
//@author sam@2020-08-07 15:17:44
func (mp *MultiPool) HGet(key, field string) (string, error) {
	return redis.String(mp.Do(DefaultPool,"HGET", key, field))
}

//4.获取在哈希表中指定 key 的所有字段和值
//@author sam@2020-08-07 15:29:36
func (mp *MultiPool) HGetAll(key string) (map[string]string, error) {
	return redis.StringMap(mp.Do(DefaultPool,"HGETALL", key))
}

//5.为哈希表 key 中的指定字段的整数值加上增量 increment 。
//@author sam@2020-08-07 15:14:51
func (mp *MultiPool) HIncrBy(key, field string, value int) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"HINCRBY", key, field, value))
}

//6.为哈希表 key 中的指定字段的浮点数值加上增量 increment 。
//@author sam@2020-08-08 11:17:54
func (mp *MultiPool) HIncrByFloat(key, field string, value float64) (float64, error) {
	return redis.Float64(mp.Do(DefaultPool,"HINCRBYFLOAT", key, field, value))
}

//7.获取所有哈希表中的字段	HKEYS key
//@author sam@2020-08-07 15:30:12
func (mp *MultiPool) HKeys(key string) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"HKEYS", key))
}

//8.获取哈希表中字段的数量	HLEN key
//@author sam@2020-08-07 15:30:43
func (mp *MultiPool) HLen(key string) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"HLEN", key))
}

//9.获取所有给定字段的值
//@author sam@2020-08-07 15:28:21
func (mp *MultiPool) HMGet(key string, fields ...interface{}) ([]interface{}, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, fields...)
	return redis.Values(mp.Do(DefaultPool,"HMGET", args...))
}

//10.同时将多个 field-value (域-值)对设置到哈希表 key 中。
//@author sam@2020-08-07 15:25:12
func (mp *MultiPool) HMSet(key string, pairs ...interface{}) (bool, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, pairs...)
	return isOKString(redis.String(mp.Do(DefaultPool,"HMSET", args...)))
}


//11.将哈希表 key 中的字段 field 的值设为 value
//@author sam@2020-08-07 15:17:02
func (mp *MultiPool) HSet(key string, field string, value interface{}) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"HSET", key, field, value))
}


//12.只有在字段 field 不存在时，设置哈希表字段的值。
//@author sam@2020-08-08 11:26:15
func (mp *MultiPool) HSetNX(key string, field string, value interface{}) (bool, error) {
	return redis.Bool(mp.Do(DefaultPool,"HSETNX", key, field, value))
}

//13.获取哈希表中所有值。
//@author sam@2020-08-07 15:30:12
func (mp *MultiPool) HVals(key string) ([]string, error) {
	return redis.Strings(mp.Do(DefaultPool,"HVALS", key))
}
//14.迭代哈希表中的键值对。HSCAN key cursor [MATCH pattern] [COUNT count]


//---------------------  HyperLogLog  (基数统计)---------------------------

//1.添加指定元素到 HyperLogLog 中
//pfadd  nodes  1
//@author sam@2020-08-07 16:56:15
func (mp *MultiPool) PFAdd(key string, els ...interface{}) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, els...)
	return redis.Int(mp.Do(DefaultPool,"PFADD", args...))
}

//2.返回给定 HyperLogLog 的基数估算值
//@author sam@2020-08-07 16:56:38
func (mp *MultiPool) PFCount(keys ...interface{}) (int, error) {
	return redis.Int(mp.Do(DefaultPool,"PFCOUNT", keys...))
}

//3.将多个 HyperLogLog 合并为一个 HyperLogLog
//@author sam@2020-08-07 16:56:54
func (mp *MultiPool) PFMerge(dest string, keys ...interface{}) (bool, error) {
	args := make([]interface{}, 0)
	args = append(args, dest)
	args = append(args, keys...)
	return isOKString(redis.String(mp.Do(DefaultPool,"PFMERGE", args...)))
}





