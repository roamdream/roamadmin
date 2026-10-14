// Package lock 分布式锁（Redis）。
//
// 基于 SETNX + 过期时间；释放用 Lua 比较 token 保证"只有持有者能释放"，
// 并发写热点数据时保证互斥。
package lock

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/roamdream/roamadmin/kernel/utils"
)

const (
	// acquireScript 原子获取：SET NX PX 一次性完成"占锁 + 设过期"，避免 SETNX 与 EXPIRE
	// 分两步导致进程崩溃时锁永不回收；PX 用毫秒，支持亚秒 TTL。
	acquireScript = `if redis.call("set", KEYS[1], ARGV[1], "NX", "PX", ARGV[2]) then return 1 else return 0 end`

	// renewScript 原子续租：仅当锁仍属本持有者时才延长过期时间。
	renewScript = `if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pexpire", KEYS[1], ARGV[2]) else return 0 end`

	// releaseScript 原子释放：仅当锁值等于本持有者 token 时才删除。
	releaseScript = `if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`
)

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

// LockWithRenew 阻塞获取锁并持续续租（watchdog），直到返回的 unlock 被调用或 ctx 超时。
// 适用于临界区可能超过 ttl 的长任务：周期把过期时间续回 ttl，任务结束调用 unlock 停续并释放。
func (l *Locker) LockWithRenew(ctx context.Context, key string, ttl, renewInterval time.Duration) (unlock func(), err error) {
	for {
		lockKey, token, acquired, err := l.acquire(ctx, key, ttl)
		if err != nil {
			return nil, err
		}
		if acquired {
			return l.renewFunc(ctx, lockKey, token, ttl, renewInterval), nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(50 * time.Millisecond):
		}
	}
}

// TryLock 非阻塞尝试获取锁；返回释放函数与是否成功。
func (l *Locker) TryLock(ctx context.Context, key string, ttl time.Duration) (unlock func(), acquired bool, err error) {
	lockKey, token, acquired, err := l.acquire(ctx, key, ttl)
	if err != nil || !acquired {
		return nil, acquired, err
	}
	return l.releaseFunc(ctx, lockKey, token), true, nil
}

// acquire 原子获取锁，返回锁 key 与 token 供释放/续租使用。
func (l *Locker) acquire(ctx context.Context, key string, ttl time.Duration) (lockKey, token string, acquired bool, err error) {
	token = utils.UUID()
	lockKey = "ra:lock:" + key
	v, err := l.redis.Do(ctx, "EVAL", acquireScript, 1, lockKey, token, ttl.Milliseconds())
	if err != nil {
		return "", "", false, err
	}
	return lockKey, token, v.Int() == 1, nil
}

func (l *Locker) releaseFunc(ctx context.Context, lockKey, token string) func() {
	return func() {
		_, _ = l.redis.Do(ctx, "EVAL", releaseScript, 1, lockKey, token)
	}
}

// renewFunc 启动续租 goroutine，返回幂等 unlock（停止续租 + 释放）。
func (l *Locker) renewFunc(ctx context.Context, lockKey, token string, ttl, renewInterval time.Duration) func() {
	if renewInterval <= 0 {
		renewInterval = ttl / 3
	}
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(renewInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_, _ = l.redis.Do(ctx, "EVAL", renewScript, 1, lockKey, token, ttl.Milliseconds())
			case <-done:
				return
			}
		}
	}()
	var once sync.Once
	return func() {
		once.Do(func() {
			close(done)
			_, _ = l.redis.Do(ctx, "EVAL", releaseScript, 1, lockKey, token)
		})
	}
}
