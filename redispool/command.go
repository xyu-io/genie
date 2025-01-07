package redispool

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// Get a Redis Client from an initialized model
func (c *RedisClient) GetOriginCli() redis.UniversalClient {
	return c.originCli
}

// Close closes clients
func (c *RedisClient) Close() error {
	RedisInsMap.Range(func(key, value interface{}) bool {
		if client, ok := value.(*RedisClient); ok {
			_ = client.originCli.Close()
		}
		return true
	})
	return nil
}

// Get a Redis Client from an initialized model
func (c *RedisClient) FormatKey(key string) string {
	return MakeBSKey(c.Option.AppName, key)
}

// Get a Redis Client from an initialized model
func (c *RedisClient) FormatKeys(keys ...string) []string {
	bsKeys := make([]string, 0)
	for _, key := range keys {
		bsKeys = append(bsKeys, MakeBSKey(c.Option.AppName, key))
	}
	return bsKeys
}

// /// Native command rewrite /////
func (c *RedisClient) ClientGetName() *redis.StringCmd {
	return c.originCli.ClientGetName(c.ctx)
}

func (c *RedisClient) Echo(message interface{}) *redis.StringCmd {
	return c.originCli.Echo(c.ctx, message)
}

func (c *RedisClient) Ping() *redis.StatusCmd {
	return c.originCli.Ping(c.ctx)
}

func (c *RedisClient) Del(keys ...string) *redis.IntCmd {
	newKeys := c.FormatKeys(keys...)
	fmt.Println(newKeys)
	return c.originCli.Del(c.ctx, newKeys...)
}

func (c *RedisClient) Unlink(keys ...string) *redis.IntCmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.Unlink(c.ctx, newKeys...)
}

func (c *RedisClient) Dump(key string) *redis.StringCmd {
	newKey := c.FormatKey(key)
	return c.originCli.Dump(c.ctx, newKey)
}

func (c *RedisClient) Exists(keys ...string) *redis.IntCmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.Exists(c.ctx, newKeys...)
}

func (c *RedisClient) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	newKey := c.FormatKey(key)
	return c.originCli.Expire(c.ctx, newKey, expiration)
}

func (c *RedisClient) ExpireAt(key string, tm time.Time) *redis.BoolCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ExpireAt(c.ctx, newKey, tm)
}

func (c *RedisClient) Migrate(host, port, key string, db int, timeout time.Duration) *redis.StatusCmd {
	newKey := c.FormatKey(key)
	return c.originCli.Migrate(c.ctx, host, port, newKey, int(db), timeout)
}

func (c *RedisClient) Move(key string, db int) *redis.BoolCmd {
	newKey := c.FormatKey(key)
	return c.originCli.Move(c.ctx, newKey, db)
}

func (c *RedisClient) Persist(key string) *redis.BoolCmd {
	newKey := c.FormatKey(key)
	return c.originCli.Persist(c.ctx, newKey)
}

func (c *RedisClient) PExpire(key string, expiration time.Duration) *redis.BoolCmd {
	newKey := c.FormatKey(key)
	return c.originCli.PExpire(c.ctx, newKey, expiration)
}

func (c *RedisClient) PExpireAt(key string, tm time.Time) *redis.BoolCmd {
	newKey := c.FormatKey(key)
	return c.originCli.PExpireAt(c.ctx, newKey, tm)
}

func (c *RedisClient) PTTL(key string) *redis.DurationCmd {
	newKey := c.FormatKey(key)
	return c.originCli.PTTL(c.ctx, newKey)
}

func (c *RedisClient) RandomKey() *redis.StringCmd {
	return c.originCli.RandomKey(c.ctx)
}

func (c *RedisClient) Rename(key, newkey string) *redis.StatusCmd {
	newKey := c.FormatKey(key)
	return c.originCli.Rename(c.ctx, newKey, newkey)
}

func (c *RedisClient) RenameNX(key, newkey string) *redis.BoolCmd {
	newKey := c.FormatKey(key)
	return c.originCli.RenameNX(c.ctx, newKey, newkey)
}

func (c *RedisClient) Restore(key string, ttl time.Duration, value string) *redis.StatusCmd {
	newKey := c.FormatKey(key)
	return c.originCli.Restore(c.ctx, newKey, ttl, value)
}

func (c *RedisClient) RestoreReplace(key string, ttl time.Duration, value string) *redis.StatusCmd {
	newKey := c.FormatKey(key)
	return c.originCli.RestoreReplace(c.ctx, newKey, ttl, value)
}

func (c *RedisClient) Sort(key string, sort *redis.Sort) *redis.StringSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.Sort(c.ctx, newKey, sort)
}

func (c *RedisClient) SortStore(key, store string, sort *redis.Sort) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SortStore(c.ctx, newKey, store, sort)
}

func (c *RedisClient) SortInterfaces(key string, sort *redis.Sort) *redis.SliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SortInterfaces(c.ctx, newKey, sort)
}

func (c *RedisClient) Touch(keys ...string) *redis.IntCmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.Touch(c.ctx, newKeys...)
}

func (c *RedisClient) TTL(key string) *redis.DurationCmd {
	newKey := c.FormatKey(key)
	return c.originCli.TTL(c.ctx, newKey)
}

func (c *RedisClient) Type(key string) *redis.StatusCmd {
	newKey := c.FormatKey(key)
	return c.originCli.Type(c.ctx, newKey)
}

func (c *RedisClient) SScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SScan(c.ctx, newKey, cursor, match, count)
}

func (c *RedisClient) HScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	newKey := c.FormatKey(key)
	return c.originCli.HScan(c.ctx, newKey, cursor, match, count)
}

func (c *RedisClient) ZScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZScan(c.ctx, newKey, cursor, match, count)
}

func (c *RedisClient) Append(key, value string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.Append(c.ctx, newKey, value)
}

func (c *RedisClient) BitCount(key string, bitCount *redis.BitCount) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.BitCount(c.ctx, newKey, bitCount)
}

func (c *RedisClient) BitOpAnd(destKey string, keys ...string) *redis.IntCmd {
	newDestKey := c.FormatKey(destKey)
	newKeys := c.FormatKeys(keys...)
	return c.originCli.BitOpAnd(c.ctx, newDestKey, newKeys...)
}

func (c *RedisClient) BitOpOr(destKey string, keys ...string) *redis.IntCmd {
	newDestKey := c.FormatKey(destKey)
	newKeys := c.FormatKeys(keys...)
	return c.originCli.BitOpOr(c.ctx, newDestKey, newKeys...)
}

func (c *RedisClient) BitOpXor(destKey string, keys ...string) *redis.IntCmd {
	newDestKey := c.FormatKey(destKey)
	newKeys := c.FormatKeys(keys...)
	return c.originCli.BitOpXor(c.ctx, newDestKey, newKeys...)
}

func (c *RedisClient) BitOpNot(destKey string, key string) *redis.IntCmd {
	newDestKey := c.FormatKey(destKey)
	newKey := c.FormatKey(key)
	return c.originCli.BitOpNot(c.ctx, newDestKey, newKey)
}

func (c *RedisClient) BitPos(key string, bit int64, pos ...int64) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.BitPos(c.ctx, newKey, bit, pos...)
}

func (c *RedisClient) Decr(key string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.Decr(c.ctx, newKey)
}

func (c *RedisClient) DecrBy(key string, decrement int64) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.DecrBy(c.ctx, newKey, decrement)
}

func (c *RedisClient) Get(key string) *redis.StringCmd {
	newKey := c.FormatKey(key)
	return c.originCli.Get(c.ctx, newKey)
}

func (c *RedisClient) GetBit(key string, offset int64) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.GetBit(c.ctx, newKey, offset)
}

func (c *RedisClient) GetRange(key string, start, end int64) *redis.StringCmd {
	newKey := c.FormatKey(key)
	return c.originCli.GetRange(c.ctx, newKey, start, end)
}

func (c *RedisClient) GetSet(key string, value interface{}) *redis.StringCmd {
	newKey := c.FormatKey(key)
	return c.originCli.GetSet(c.ctx, newKey, value)
}

func (c *RedisClient) Incr(key string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.Incr(c.ctx, newKey)
}

func (c *RedisClient) IncrBy(key string, value int64) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.IncrBy(c.ctx, newKey, value)
}

func (c *RedisClient) IncrByFloat(key string, value float64) *redis.FloatCmd {
	newKey := c.FormatKey(key)
	return c.originCli.IncrByFloat(c.ctx, newKey, value)
}

func (c *RedisClient) MGet(keys ...string) *redis.SliceCmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.MGet(c.ctx, newKeys...)
}

func (c *RedisClient) MSet(pairs ...interface{}) *redis.StatusCmd {
	return c.originCli.MSet(c.ctx, pairs...)
}

func (c *RedisClient) MSetNX(pairs ...interface{}) *redis.BoolCmd {
	return c.originCli.MSetNX(c.ctx, pairs...)
}

func (c *RedisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	newKeys := c.FormatKey(key)
	return c.originCli.Set(c.ctx, newKeys, value, expiration)
}

func (c *RedisClient) SetBit(key string, offset int64, value int) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SetBit(c.ctx, newKey, offset, value)
}

func (c *RedisClient) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SetNX(c.ctx, newKey, value, expiration)
}

func (c *RedisClient) SetXX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SetXX(c.ctx, newKey, value, expiration)
}

func (c *RedisClient) SetRange(key string, offset int64, value string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SetRange(c.ctx, newKey, offset, value)
}

func (c *RedisClient) StrLen(key string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.StrLen(c.ctx, newKey)
}

func (c *RedisClient) HDel(key string, fields ...string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.HDel(c.ctx, newKey, fields...)
}

func (c *RedisClient) HExists(key, field string) *redis.BoolCmd {
	newKey := c.FormatKey(key)
	return c.originCli.HExists(c.ctx, newKey, field)
}

func (c *RedisClient) HGet(key, field string) *redis.StringCmd {
	newKey := c.FormatKey(key)
	return c.originCli.HGet(c.ctx, newKey, field)
}

func (c *RedisClient) HGetAll(key string) *redis.MapStringStringCmd {
	newKey := c.FormatKey(key)
	return c.originCli.HGetAll(c.ctx, newKey)
}

func (c *RedisClient) HIncrBy(key, field string, incr int64) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.HIncrBy(c.ctx, newKey, field, incr)
}

func (c *RedisClient) HIncrByFloat(key, field string, incr float64) *redis.FloatCmd {
	newKey := c.FormatKey(key)
	return c.originCli.HIncrByFloat(c.ctx, newKey, field, incr)
}

func (c *RedisClient) HKeys(key string) *redis.StringSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.HKeys(c.ctx, newKey)
}

func (c *RedisClient) HLen(key string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.HLen(c.ctx, newKey)
}

func (c *RedisClient) HMGet(key string, fields ...string) *redis.SliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.HMGet(c.ctx, newKey, fields...)
}

func (c *RedisClient) HMSet(key string, fields map[string]interface{}) *redis.BoolCmd {
	newKey := c.FormatKey(key)
	return c.originCli.HMSet(c.ctx, newKey, fields)
}

func (c *RedisClient) HSet(key, field string, value interface{}) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.HSet(c.ctx, newKey, field, value)
}

func (c *RedisClient) HSetNX(key, field string, value interface{}) *redis.BoolCmd {
	newKey := c.FormatKey(key)
	return c.originCli.HSetNX(c.ctx, newKey, field, value)
}

func (c *RedisClient) HVals(key string) *redis.StringSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.HVals(c.ctx, newKey)
}

func (c *RedisClient) BLPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.BLPop(c.ctx, timeout, newKeys...)
}

func (c *RedisClient) BRPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.BRPop(c.ctx, timeout, newKeys...)
}

func (c *RedisClient) BRPopLPush(source, destination string, timeout time.Duration) *redis.StringCmd {
	return c.originCli.BRPopLPush(c.ctx, source, destination, timeout)
}

func (c *RedisClient) LIndex(key string, index int64) *redis.StringCmd {
	newKey := c.FormatKey(key)
	return c.originCli.LIndex(c.ctx, newKey, index)
}

func (c *RedisClient) LInsert(key, op string, pivot, value interface{}) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.LInsert(c.ctx, newKey, op, pivot, value)
}

func (c *RedisClient) LInsertBefore(key string, pivot, value interface{}) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.LInsertBefore(c.ctx, newKey, pivot, value)
}

func (c *RedisClient) LInsertAfter(key string, pivot, value interface{}) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.LInsertAfter(c.ctx, newKey, pivot, value)
}

func (c *RedisClient) LLen(key string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.LLen(c.ctx, newKey)
}

func (c *RedisClient) LPop(key string) *redis.StringCmd {
	newKey := c.FormatKey(key)
	return c.originCli.LPop(c.ctx, newKey)
}

func (c *RedisClient) LPush(key string, values ...interface{}) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.LPush(c.ctx, newKey, values)
}

func (c *RedisClient) LPushX(key string, value interface{}) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.LPushX(c.ctx, newKey, value)
}

func (c *RedisClient) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.LRange(c.ctx, newKey, start, stop)
}

func (c *RedisClient) LRem(key string, count int64, value interface{}) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.LRem(c.ctx, newKey, count, value)
}

func (c *RedisClient) LSet(key string, index int64, value interface{}) *redis.StatusCmd {
	newKey := c.FormatKey(key)
	return c.originCli.LSet(c.ctx, newKey, index, value)
}

func (c *RedisClient) LTrim(key string, start, stop int64) *redis.StatusCmd {
	newKey := c.FormatKey(key)
	return c.originCli.LTrim(c.ctx, newKey, start, stop)
}

func (c *RedisClient) RPop(key string) *redis.StringCmd {
	newKey := c.FormatKey(key)
	return c.originCli.RPop(c.ctx, newKey)
}

func (c *RedisClient) RPopLPush(source, destination string) *redis.StringCmd {
	return c.originCli.RPopLPush(c.ctx, source, destination)
}

func (c *RedisClient) RPush(key string, values ...interface{}) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.RPush(c.ctx, newKey, values...)
}

func (c *RedisClient) RPushX(key string, value interface{}) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.RPushX(c.ctx, newKey, value)
}

func (c *RedisClient) SAdd(key string, members ...interface{}) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SAdd(c.ctx, newKey, members...)
}

func (c *RedisClient) SCard(key string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SCard(c.ctx, newKey)
}

func (c *RedisClient) SDiff(keys ...string) *redis.StringSliceCmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.SDiff(c.ctx, newKeys...)
}

func (c *RedisClient) SDiffStore(destination string, keys ...string) *redis.IntCmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.SDiffStore(c.ctx, destination, newKeys...)
}

func (c *RedisClient) SInter(keys ...string) *redis.StringSliceCmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.SInter(c.ctx, newKeys...)
}

func (c *RedisClient) SInterStore(destination string, keys ...string) *redis.IntCmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.SInterStore(c.ctx, destination, newKeys...)
}

func (c *RedisClient) SIsMember(key string, member interface{}) *redis.BoolCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SIsMember(c.ctx, newKey, member)
}

func (c *RedisClient) SMembers(key string) *redis.StringSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SMembers(c.ctx, newKey)
}

func (c *RedisClient) SMembersMap(key string) *redis.StringStructMapCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SMembersMap(c.ctx, newKey)
}

func (c *RedisClient) SMove(source, destination string, member interface{}) *redis.BoolCmd {
	return c.originCli.SMove(c.ctx, source, destination, member)
}

func (c *RedisClient) SPop(key string) *redis.StringCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SPop(c.ctx, newKey)
}

func (c *RedisClient) SPopN(key string, count int64) *redis.StringSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SPopN(c.ctx, newKey, count)
}

func (c *RedisClient) SRandMember(key string) *redis.StringCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SRandMember(c.ctx, newKey)
}

func (c *RedisClient) SRandMemberN(key string, count int64) *redis.StringSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SRandMemberN(c.ctx, newKey, count)
}

func (c *RedisClient) SRem(key string, members ...interface{}) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.SRem(c.ctx, newKey, members...)
}

func (c *RedisClient) SUnion(keys ...string) *redis.StringSliceCmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.SUnion(c.ctx, newKeys...)
}

func (c *RedisClient) SUnionStore(destination string, keys ...string) *redis.IntCmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.SUnionStore(c.ctx, destination, newKeys...)
}

func (c *RedisClient) ZAdd(key string, members ...redis.Z) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZAdd(c.ctx, newKey, members...)
}

func (c *RedisClient) ZAddNX(key string, members ...redis.Z) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZAddNX(c.ctx, newKey, members...)
}

func (c *RedisClient) ZAddXX(key string, members ...redis.Z) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZAddXX(c.ctx, newKey, members...)
}

func (c *RedisClient) ZAddGT(key string, members ...redis.Z) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZAddGT(c.ctx, newKey, members...)
}

func (c *RedisClient) ZAddLT(key string, members ...redis.Z) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZAddLT(c.ctx, newKey, members...)
}

func (c *RedisClient) ZAddArgs(key string, members redis.ZAddArgs) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZAddArgs(c.ctx, newKey, members)
}

func (c *RedisClient) ZAddArgsIncr(key string, member redis.ZAddArgs) *redis.FloatCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZAddArgsIncr(c.ctx, newKey, member)
}

func (c *RedisClient) ZCard(key string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZCard(c.ctx, newKey)
}

func (c *RedisClient) ZCount(key, min, max string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZCount(c.ctx, newKey, min, max)
}

func (c *RedisClient) ZLexCount(key, min, max string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZLexCount(c.ctx, newKey, min, max)
}

func (c *RedisClient) ZIncrBy(key string, increment float64, member string) *redis.FloatCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZIncrBy(c.ctx, newKey, increment, member)
}

func (c *RedisClient) ZInterStore(destination string, store *redis.ZStore) *redis.IntCmd {
	return c.originCli.ZInterStore(c.ctx, destination, store)
}

func (c *RedisClient) ZPopMax(key string, count ...int64) *redis.ZSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZPopMax(c.ctx, newKey, count...)
}

func (c *RedisClient) ZPopMin(key string, count ...int64) *redis.ZSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZPopMin(c.ctx, newKey, count...)
}

func (c *RedisClient) ZRange(key string, start, stop int64) *redis.StringSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRange(c.ctx, newKey, start, stop)
}

func (c *RedisClient) ZRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRangeWithScores(c.ctx, newKey, start, stop)
}

func (c *RedisClient) ZRangeByScore(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRangeByScore(c.ctx, newKey, opt)
}

func (c *RedisClient) ZRangeByLex(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRangeByLex(c.ctx, newKey, opt)
}

func (c *RedisClient) ZRangeByScoreWithScores(key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRangeByScoreWithScores(c.ctx, newKey, opt)
}

func (c *RedisClient) ZRank(key, member string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRank(c.ctx, newKey, member)
}

func (c *RedisClient) ZRem(key string, members ...interface{}) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRem(c.ctx, newKey, members...)
}

func (c *RedisClient) ZRemRangeByRank(key string, start, stop int64) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRemRangeByRank(c.ctx, newKey, start, stop)
}

func (c *RedisClient) ZRemRangeByScore(key, min, max string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRemRangeByScore(c.ctx, newKey, min, max)
}

func (c *RedisClient) ZRemRangeByLex(key, min, max string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRemRangeByLex(c.ctx, newKey, min, max)
}

func (c *RedisClient) ZRevRange(key string, start, stop int64) *redis.StringSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRevRange(c.ctx, newKey, start, stop)
}

func (c *RedisClient) ZRevRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRevRangeWithScores(c.ctx, newKey, start, stop)
}

func (c *RedisClient) ZRevRangeByScore(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRevRangeByScore(c.ctx, newKey, opt)
}

func (c *RedisClient) ZRevRangeByLex(key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRevRangeByLex(c.ctx, newKey, opt)
}

func (c *RedisClient) ZRevRangeByScoreWithScores(key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRevRangeByScoreWithScores(c.ctx, newKey, opt)
}

func (c *RedisClient) ZRevRank(key, member string) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZRevRank(c.ctx, newKey, member)
}

func (c *RedisClient) ZScore(key, member string) *redis.FloatCmd {
	newKey := c.FormatKey(key)
	return c.originCli.ZScore(c.ctx, newKey, member)
}

func (c *RedisClient) ZUnionStore(dest string, store *redis.ZStore) *redis.IntCmd {
	return c.originCli.ZUnionStore(c.ctx, dest, store)
}

func (c *RedisClient) PFAdd(key string, els ...interface{}) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.PFAdd(c.ctx, newKey, els...)
}

func (c *RedisClient) PFCount(keys ...string) *redis.IntCmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.PFCount(c.ctx, newKeys...)
}

func (c *RedisClient) PFMerge(dest string, keys ...string) *redis.StatusCmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.PFMerge(c.ctx, dest, newKeys...)
}

func (c *RedisClient) SlaveOf(host, port string) *redis.StatusCmd {
	return c.originCli.SlaveOf(c.ctx, host, port)
}

func (c *RedisClient) Time() *redis.TimeCmd {
	return c.originCli.Time(c.ctx)
}

func (c *RedisClient) Eval(script string, keys []string, args ...interface{}) *redis.Cmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.Eval(c.ctx, script, newKeys, args...)
}

func (c *RedisClient) EvalSha(sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	newKeys := c.FormatKeys(keys...)
	return c.originCli.EvalSha(c.ctx, sha1, newKeys, args...)
}

// 用于订阅指定的单个 Redis 频道（channel）
func (c *RedisClient) Subscribe(channels ...string) *redis.PubSub {
	return c.originCli.Subscribe(c.ctx, channels...)
}

// 代表的是模式订阅（Pattern Subscribe）。
//
//	它允许客户端按照一定的模式来订阅多个频道，只要频道名称符合所设定的模式规则，就能接收到对应频道发布的消息。例如，使用模式 "news_*" 进行订阅
func (c *RedisClient) PSubscribe(channels ...string) *redis.PubSub {
	return c.originCli.PSubscribe(c.ctx, channels...)
}

// 在 Redis 集群环境下用于订阅单个频道
func (c *RedisClient) SSubscribe(channels ...string) *redis.PubSub {
	return c.originCli.SSubscribe(c.ctx, channels...)
}

func (c *RedisClient) Publish(channel string, message interface{}) *redis.IntCmd {
	return c.originCli.Publish(c.ctx, channel, message)
}

func (c *RedisClient) PubSubChannels(pattern string) *redis.StringSliceCmd {
	return c.originCli.PubSubChannels(c.ctx, pattern)
}

func (c *RedisClient) PubSubNumSub(channels ...string) *redis.MapStringIntCmd {
	return c.originCli.PubSubNumSub(c.ctx, channels...)
}

func (c *RedisClient) PubSubNumPat() *redis.IntCmd {
	return c.originCli.PubSubNumPat(c.ctx)
}

func (c *RedisClient) GeoAdd(key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.GeoAdd(c.ctx, newKey, geoLocation...)
}

func (c *RedisClient) GeoPos(key string, members ...string) *redis.GeoPosCmd {
	newKey := c.FormatKey(key)
	return c.originCli.GeoPos(c.ctx, newKey, members...)
}

func (c *RedisClient) GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	newKey := c.FormatKey(key)
	return c.originCli.GeoRadius(c.ctx, newKey, longitude, latitude, query)
}

func (c *RedisClient) GeoGeoRadiusStore(key string, longitude float64, latitude float64, query *redis.GeoRadiusQuery) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.GeoRadiusStore(c.ctx, newKey, longitude, latitude, query)
}

func (c *RedisClient) GeoRadiusByMember(key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	newKey := c.FormatKey(key)
	return c.originCli.GeoRadiusByMember(c.ctx, newKey, member, query)
}

func (c *RedisClient) GeoRadiusByMemberRO(key string, member string, query *redis.GeoRadiusQuery) *redis.IntCmd {
	newKey := c.FormatKey(key)
	return c.originCli.GeoRadiusByMemberStore(c.ctx, newKey, member, query)
}

func (c *RedisClient) GeoDist(key string, member1, member2, unit string) *redis.FloatCmd {
	newKey := c.FormatKey(key)
	return c.originCli.GeoDist(c.ctx, newKey, member1, member2, unit)
}

func (c *RedisClient) GeoHash(key string, members ...string) *redis.StringSliceCmd {
	newKey := c.FormatKey(key)
	return c.originCli.GeoHash(c.ctx, newKey, members...)
}
