// Package cache 缓存封装：Redis（gredis）+ 本地内存缓存（gcache）。
package cache

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
)

// Cache Redis 缓存。
type Cache struct {
	inner *gredis.Redis
}

// New 包装一个 gredis.Redis。
func New(inner *gredis.Redis) *Cache { return &Cache{inner: inner} }

// Default 返回默认 Redis（等价 g.Redis()）。
func Default() *Cache { return &Cache{inner: g.Redis()} }

// Redis 返回底层 gredis.Redis（供高级用法，如发布订阅、脚本）。
func (c *Cache) Redis() *gredis.Redis { return c.inner }

// Get 读取缓存。
func (c *Cache) Get(ctx context.Context, key string) (*gvar.Var, error) {
	return c.inner.Get(ctx, key)
}

// Set 写入缓存；ttl <= 0 表示不过期。
func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if ttl > 0 {
		return c.inner.SetEX(ctx, key, value, int64(ttl.Seconds()))
	}
	_, err := c.inner.Set(ctx, key, value)
	return err
}

// Del 删除缓存。
func (c *Cache) Del(ctx context.Context, keys ...string) (int64, error) {
	return c.inner.Del(ctx, keys...)
}

// Exists 判断键是否存在。
func (c *Cache) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.inner.Exists(ctx, keys...)
}

// Expire 设置过期时间。
func (c *Cache) Expire(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	return c.inner.Expire(ctx, key, int64(ttl.Seconds()))
}

// Local 本地内存缓存（gcache，内置容量上限与 LRU 淘汰），热数据优先走本地。
type Local struct {
	inner *gcache.Cache
}

// NewLocal 创建本地缓存。
func NewLocal() *Local { return &Local{inner: gcache.New()} }

// Get 读取本地缓存。
func (l *Local) Get(ctx context.Context, key any) (*gvar.Var, error) {
	return l.inner.Get(ctx, key)
}

// Set 写入本地缓存；duration <= 0 表示不过期。
func (l *Local) Set(ctx context.Context, key, value any, duration time.Duration) error {
	return l.inner.Set(ctx, key, value, duration)
}

// Remove 删除本地缓存。
func (l *Local) Remove(ctx context.Context, keys ...any) (*gvar.Var, error) {
	return l.inner.Remove(ctx, keys...)
}

// Contains 判断本地缓存是否存在。
func (l *Local) Contains(ctx context.Context, key any) (bool, error) {
	return l.inner.Contains(ctx, key)
}
