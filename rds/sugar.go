package rds

import (
	"github.com/garyburd/redigo/redis"
)

//@todo  注意，下面的所有方法都是针对默认连接池封装的方法糖
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


//------------------------List链表类型操作----------------------------------------------------------------------------






//func (mp *MultiPool) GetRange(key string, start, end int64) (string, error) {
//	return redis.String(mp.Do(DefaultPool,"GETRANGE", key, start, end))
//}
//
//func (mp *MultiPool) SetRange(key string, offset int64, value string) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"SETRANGE", key, offset, value))
//}
//

//func (mp *MultiPool) SetEX(key string, value interface{}, seconds int64) (bool, error) {
//	return isOKString(redis.String(mp.Do(DefaultPool,"SET", key, value, "EX", seconds)))
//}
//
//func (mp *MultiPool) SetPX(key string, value interface{}, milliseconds int64) (bool, error) {
//	return isOKString(redis.String(mp.Do(DefaultPool,"SET", key, value, "PX", milliseconds)))
//}
//
//func (mp *MultiPool) SetNX(key string, value interface{}) (bool, error) {
//	return redis.Bool(mp.Do(DefaultPool,"SETNX", key, value))
//}
//
//func (mp *MultiPool) SetXX(key string, value interface{}) (bool, error) {
//	return isOKString(redis.String(mp.Do(DefaultPool,"SET", key, value, "XX")))
//}
//

//func (mp *MultiPool) StrLen(key string) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"STRLEN", key))
//}
//
//func (mp *MultiPool) GetBit(key string, offset int64) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"GETBIT", key, offset))
//}
//
//func (mp *MultiPool) SetBit(key string, offset int64, value int) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"SETBIT", key, offset, value))
//}
//
//func (mp *MultiPool) BitCount(key string, offsets ...int64) (int, error) {
//	switch len(offsets) {
//	case 0:
//		return redis.Int(mp.Do(DefaultPool,"BITCOUNT", key))
//	case 2:
//		return redis.Int(mp.Do(DefaultPool,"BITCOUNT", key, offsets[0], offsets[1]))
//	default:
//		return 0, errWrongArguments
//	}
//}
//
//func (mp *MultiPool) BitOpAnd(destKey string, keys ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, "AND", destKey)
//	args = append(args, keys...)
//	return redis.Int(mp.Do(DefaultPool,"BITOP", args...))
//}
//
//func (mp *MultiPool) BitOpOr(destKey string, keys ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, "OR", destKey)
//	args = append(args, keys...)
//	return redis.Int(mp.Do(DefaultPool,"BITOP", args...))
//}
//
//func (mp *MultiPool) BitOpXor(destKey string, keys ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, "XOR", destKey)
//	args = append(args, keys...)
//	return redis.Int(mp.Do(DefaultPool,"BITOP", args...))
//}
//
//func (mp *MultiPool) BitOpNot(destKey, key string) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"BITOP", "NOT", destKey, key))
//}
//
//func (mp *MultiPool) BitPos(key string, bit int64, offsets ...int64) (int, error) {
//	switch len(offsets) {
//	case 0:
//		return redis.Int(mp.Do(DefaultPool,"BITPOS", key, bit))
//	case 1:
//		return redis.Int(mp.Do(DefaultPool,"BITPOS", key, bit, offsets[0]))
//	case 2:
//		return redis.Int(mp.Do(DefaultPool,"BITPOS", key, bit, offsets[0], offsets[1]))
//	default:
//		return 0, errWrongArguments
//	}
//}
//
//func (mp *MultiPool) HDel(key string, fields ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, key)
//	args = append(args, fields...)
//	return redis.Int(mp.Do(DefaultPool,"HDEL", args...))
//}
//
//func (mp *MultiPool) HExists(key, field string) (bool, error) {
//	return redis.Bool(mp.Do(DefaultPool,"HEXISTS", key, field))
//}
//
//func (mp *MultiPool) HGet(key, field string) (string, error) {
//	return redis.String(mp.Do(DefaultPool,"HGET", key, field))
//}
//
//func (mp *MultiPool) HGetAll(key string) (map[string]string, error) {
//	return redis.StringMap(mp.Do(DefaultPool,"HGETALL", key))
//}
//
//func (mp *MultiPool) HIncrBy(key, field string, value int) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"HINCRBY", key, field, value))
//}
//
//func (mp *MultiPool) HIncrByFloat(key, field string, value float64) (float64, error) {
//	return redis.Float64(mp.Do(DefaultPool,"HINCRBYFLOAT", key, field, value))
//}
//
//func (mp *MultiPool) HMGet(key string, fields ...interface{}) ([]interface{}, error) {
//	args := make([]interface{}, 0)
//	args = append(args, key)
//	args = append(args, fields...)
//	return redis.Values(mp.Do(DefaultPool,"HMGET", args...))
//}
//
//func (mp *MultiPool) HMSet(key string, pairs ...interface{}) (bool, error) {
//	args := make([]interface{}, 0)
//	args = append(args, key)
//	args = append(args, pairs...)
//	return isOKString(redis.String(mp.Do(DefaultPool,"HMSET", args...)))
//}
//
//func (mp *MultiPool) HSet(key string, field string, value interface{}) (bool, error) {
//	return redis.Bool(mp.Do(DefaultPool,"HSET", key, field, value))
//}
//
//func (mp *MultiPool) HSetNX(key string, field string, value interface{}) (bool, error) {
//	return redis.Bool(mp.Do(DefaultPool,"HSETNX", key, field, value))
//}
//
//func (mp *MultiPool) HVals(key string) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"HVALS", key))
//}
//
//func (mp *MultiPool) HLen(key string) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"HLEN", key))
//}
//
//func (mp *MultiPool) BRPopLPush(source, destination string, timeout uint64) (string, error) {
//	return redis.String(mp.Do(DefaultPool,"BRPOPLPUSH", source, destination, timeout))
//}
//
//func (mp *MultiPool) LIndex(key string, index int64) (string, error) {
//	return redis.String(mp.Do(DefaultPool,"LINDEX", key, index))
//}
//
//func (mp *MultiPool) LInsert(key string, op string, pivot, value interface{}) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"LINSERT", key, op, pivot, value))
//}
//
//func (mp *MultiPool) LInsertBefore(key string, pivot, value interface{}) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"LINSERT", key, "BEFORE", pivot, value))
//}
//
//func (mp *MultiPool) LInsertAfter(key string, pivot, value interface{}) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"LINSERT", key, "AFTER", pivot, value))
//}
//
//func (mp *MultiPool) LLen(key string) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"LLEN", key))
//}
//
//func (mp *MultiPool) LPop(key string) (string, error) {
//	return redis.String(mp.Do(DefaultPool,"LPOP", key))
//}
//
//func (mp *MultiPool) BLPop(args ...interface{}) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"BLPOP", args...))
//}
//
//func (mp *MultiPool) LPush(key string, values ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, key)
//	args = append(args, values...)
//	return redis.Int(mp.Do(DefaultPool,"LPUSH", args...))
//}
//
//func (mp *MultiPool) LPushX(key string, value interface{}) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"LPUSHX", key, value))
//}
//
//func (mp *MultiPool) LRange(key string, start, stop int64) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"LRANGE", key, start, stop))
//}
//
//func (mp *MultiPool) LRem(key string, count int64, value interface{}) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"LREM", key, count, value))
//}
//
//func (mp *MultiPool) LSet(key string, index int64, value interface{}) (bool, error) {
//	return isOKString(redis.String(mp.Do(DefaultPool,"LSET", key, index, value)))
//}
//
//func (mp *MultiPool) LTrim(key string, start, stop int64) (bool, error) {
//	return isOKString(redis.String(mp.Do(DefaultPool,"LTRIM", key, start, stop)))
//}
//
//func (mp *MultiPool) RPop(key string) (string, error) {
//	return redis.String(mp.Do(DefaultPool,"RPOP", key))
//}
//
//func (mp *MultiPool) RPopLPush(source, destination string) (string, error) {
//	return redis.String(mp.Do(DefaultPool,"RPOPLPUSH", source, destination))
//}
//
//func (mp *MultiPool) RPush(key string, values ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, key)
//	args = append(args, values...)
//	return redis.Int(mp.Do(DefaultPool,"RPUSH", args...))
//}
//
//func (mp *MultiPool) RPushX(key string, value interface{}) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"RPUSHX", key, value))
//}
//
//func (mp *MultiPool) SAdd(key string, members ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, key)
//	args = append(args, members...)
//	return redis.Int(mp.Do(DefaultPool,"SADD", args...))
//}
//
//func (mp *MultiPool) SCard(key string) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"SCARD", key))
//}
//
//func (mp *MultiPool) SDiff(keys ...interface{}) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"SDIFF", keys...))
//}
//
//func (mp *MultiPool) SDiffStore(destination string, keys ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, destination)
//	args = append(args, keys...)
//	return redis.Int(mp.Do(DefaultPool,"SDIFFSTORE", args...))
//}
//
//func (mp *MultiPool) SInter(keys ...interface{}) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"SINTER", keys...))
//}
//
//func (mp *MultiPool) SInterStore(destination string, keys ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, destination)
//	args = append(args, keys...)
//	return redis.Int(mp.Do(DefaultPool,"SINTERSTORE", args...))
//}
//
//func (mp *MultiPool) SIsMember(key string, member interface{}) (bool, error) {
//	return redis.Bool(mp.Do(DefaultPool,"SISMEMBER", key, member))
//}
//
//func (mp *MultiPool) SMove(source, destination string, member interface{}) (bool, error) {
//	return redis.Bool(mp.Do(DefaultPool,"SMOVE", source, destination, member))
//}
//
//func (mp *MultiPool) SMembers(key string) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"SMEMBERS", key))
//}
//
//func (mp *MultiPool) SPop(key string) (string, error) {
//	return redis.String(mp.Do(DefaultPool,"SPOP", key))
//}
//
//func (mp *MultiPool) SPopN(key string, count int64) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"SPOP", key, count))
//}
//
//func (mp *MultiPool) SRandMember(key string) (string, error) {
//	return redis.String(mp.Do(DefaultPool,"SRANDMEMBER", key))
//}
//
//func (mp *MultiPool) SRandMemberN(key string, count int64) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"SRANDMEMBER", key, count))
//}
//
//func (mp *MultiPool) SRem(key string, members ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, key)
//	args = append(args, members...)
//	return redis.Int(mp.Do(DefaultPool,"SREM", args...))
//}
//
//func (mp *MultiPool) SUnion(keys ...interface{}) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"SUNION", keys...))
//}
//
//func (mp *MultiPool) SUnionStore(destionation string, keys ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, destionation)
//	args = append(args, keys...)
//	return redis.Int(mp.Do(DefaultPool,"SUNIONSTORE", args...))
//}
//
//func (mp *MultiPool) ZAdd(key string, pairs ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, key)
//	args = append(args, pairs...)
//	return redis.Int(mp.Do(DefaultPool,"ZADD", args...))
//}
//
//func (mp *MultiPool) ZCard(key string) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"ZCARD", key))
//}
//
//func (mp *MultiPool) ZCount(key string, min, max interface{}) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"ZCOUNT", key, min, max))
//}
//
//func (mp *MultiPool) ZIncrBy(key string, increment float64, member string) (float64, error) {
//	return redis.Float64(mp.Do(DefaultPool,"ZINCRBY", key, increment, member))
//}
//
//func (mp *MultiPool) ZInterStore(destination string, nkeys int, params ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, destination)
//	args = append(args, nkeys)
//	args = append(args, params...)
//	return redis.Int(mp.Do(DefaultPool,"ZINTERSTORE", args...))
//}
//
//func (mp *MultiPool) ZLexCount(key, min, max string) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"ZLEXCOUNT", key, min, max))
//}
//
//func (mp *MultiPool) ZRange(key string, start, stop int64) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"ZRANGE", key, start, stop))
//}
//
//func (mp *MultiPool) ZRangeWithScores(key string, start, stop int64) (map[string]string, error) {
//	return redis.StringMap(mp.Do(DefaultPool,"ZRANGE", key, start, stop, "WITHSCORES"))
//}
//
//func (mp *MultiPool) ZRangeByLex(key, min, max string, offset, count int64) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"ZRANGEBYLEX", key, min, max, "LIMIT", offset, count))
//}
//
//func (mp *MultiPool) ZRangeByScore(key, min, max interface{}, offset, count int64) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"ZRANGEBYSCORE", key, min, max, "LIMIT", offset, count))
//}
//
//func (mp *MultiPool) ZRank(key, member string) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"ZRANK", key, member))
//}
//
//func (mp *MultiPool) ZRem(key string, members ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, key)
//	args = append(args, members...)
//	return redis.Int(mp.Do(DefaultPool,"ZREM", args...))
//}
//
//func (mp *MultiPool) ZRemRangeByLex(key, min, max string) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"ZREMRANGEBYLEX", key, min, max))
//}
//
//func (mp *MultiPool) ZRemRangeByRank(key string, start, stop int64) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"ZREMRANGEByRANK", key, start, stop))
//}
//
//func (mp *MultiPool) ZRemRangeByScore(key, min, max interface{}) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"ZREMRANGEBYSCORE", key, min, max))
//}
//
//func (mp *MultiPool) ZRevRange(key string, start, stop int64) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"ZREVRANGE", key, start, stop))
//}
//
//func (mp *MultiPool) ZRevRangeWithScores(key string, start, stop int64) (map[string]string, error) {
//	return redis.StringMap(mp.Do(DefaultPool,"ZREVRANGE", key, start, stop, "WITHSCORES"))
//}
//
//func (mp *MultiPool) ZRevRangeByLex(key, max, min string, offset, count int64) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"ZREVRANGEBYLEX", key, max, min, "LIMIT", offset, count))
//}
//
//func (mp *MultiPool) ZRevRangeByScore(key, max, min interface{}, offset, count int64) ([]string, error) {
//	return redis.Strings(mp.Do(DefaultPool,"ZREVRANGEBYSCORE", key, max, min, "LIMIT", offset, count))
//}
//
//func (mp *MultiPool) ZRevRangeByScoreWithScores(key, max, min interface{}, offset, count int64) (map[string]string, error) {
//	return redis.StringMap(mp.Do(DefaultPool,"ZREVRANGEBYSCORE", key, max, min, "WITHSCORES", "LIMIT", offset, count))
//}
//
//func (mp *MultiPool) ZRevRank(key, member string) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"ZREVRANK", key, member))
//}
//
//func (mp *MultiPool) ZScore(key, member string) (float64, error) {
//	return redis.Float64(mp.Do(DefaultPool,"ZSCORE", key, member))
//}
//
//func (mp *MultiPool) ZUnionStore(destination string, nkeys int, params ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, destination)
//	args = append(args, nkeys)
//	args = append(args, params...)
//	return redis.Int(mp.Do(DefaultPool,"ZUNIONSTORE", args...))
//}
//
//func (mp *MultiPool) PFAdd(key string, els ...interface{}) (int, error) {
//	args := make([]interface{}, 0)
//	args = append(args, key)
//	args = append(args, els...)
//	return redis.Int(mp.Do(DefaultPool,"PFADD", args...))
//}
//
//func (mp *MultiPool) PFCount(keys ...interface{}) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"PFCOUNT", keys...))
//}
//
//func (mp *MultiPool) PFMerge(dest string, keys ...interface{}) (bool, error) {
//	args := make([]interface{}, 0)
//	args = append(args, dest)
//	args = append(args, keys...)
//	return isOKString(redis.String(mp.Do(DefaultPool,"PFMERGE", args...)))
//}
//
//func (mp *MultiPool) Publish(channel, msg string) (int, error) {
//	return redis.Int(mp.Do(DefaultPool,"PUBLISH", channel, msg))
//}
//
//func (mp *MultiPool) RawCommand(command string, args ...interface{}) (interface{}, error) {
//	return mp.Do(DefaultPool,command, args...)
//}