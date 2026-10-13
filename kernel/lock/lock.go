// Package lock 分布式锁（Redis）。
//
// 基于 SETNX + 过期时间；释放用 Lua 比较 token 保证"只有持有者能释放"，
// 并发写热点数据时保证互斥。
package lock

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/roamdream/roamadmin/kernel/utils"
)

// releaseScript 原子释放：仅当锁值等于本持有者 token 时才删除。
const releaseScript = `if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`

// Locker Redis 分布式锁。
type Locker struct {
	redis *gredis.Redis
}

// New 包装一个 gredis.Redis。
func New(r *gredis.Redis) *Locker { return &Locker{redis: r} }

// Default 返回默认 Locker（等价 g.Redis()）。
func Default() *Locker { return &Locker{redis: g.Redis()} }

// Lock 阻塞获取锁，直到成功或 ctx 超时；返回释放函数。
func (l *Locker) Lock(ctx context.Context, key string, ttl time.Duration) (unlock func(), err error) {
	return l.LockWithInterval(ctx, key, ttl, 50*time.Millisecond)
}

// LockWithInterval 阻塞获取锁，可指定重试间隔。
func (l *Locker) LockWithInterval(ctx context.Context, key string, ttl time.Duration, interval time.Duration) (unlock func(), err error) {
	for {
		var acquired bool
		unlock, acquired, err = l.TryLock(ctx, key, ttl)
		if err != nil {
			return nil, err
		}
		if acquired {
			return unlock, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(interval):
		}
	}
}

// TryLock 非阻塞尝试获取锁；返回释放函数与是否成功。
func (l *Locker) TryLock(ctx context.Context, key string, ttl time.Duration) (unlock func(), acquired bool, err error) {
	token := utils.UUID()
	lockKey := "ra:lock:" + key
	ok, err := l.redis.SetNX(ctx, lockKey, token)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	_, _ = l.redis.Expire(ctx, lockKey, int64(ttl.Seconds()))
	return l.releaseFunc(ctx, lockKey, token), true, nil
}

func (l *Locker) releaseFunc(ctx context.Context, lockKey, token string) func() {
	return func() {
		_, _ = l.redis.Do(ctx, "EVAL", releaseScript, 1, lockKey, token)
	}
}
