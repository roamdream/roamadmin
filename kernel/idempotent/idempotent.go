// Package idempotent 幂等键 / 幂等去重。
//
// 基于 Redis SETNX：首次出现的幂等键占用成功（返回 true），重复键返回 false。
// 用于接口幂等与消息消费幂等（重复上报、消息重投不重复落库）。
package idempotent

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
)

// acquireScript 原子占用幂等键：SET NX PX 一次完成"占键 + 设过期"，避免 SETNX 与 EXPIRE
// 分两步导致进程崩溃时幂等键永不回收。
const acquireScript = `if redis.call("set", KEYS[1], ARGV[1], "NX", "PX", ARGV[2]) then return 1 else return 0 end`

// Manager 幂等管理器。
type Manager struct {
	redis *gredis.Redis
}

// New 包装一个 gredis.Redis。
func New(r *gredis.Redis) *Manager { return &Manager{redis: r} }

// Default 返回默认 Manager。
func Default() *Manager { return &Manager{redis: g.Redis()} }

// Check 检查幂等键：首次出现则占用并返回 true，重复出现返回 false。
// ttl 控制幂等键的存活时间（即幂等窗口）；ttl<=0 表示不过期（单条 SETNX 本身原子）。
func (m *Manager) Check(ctx context.Context, key string, ttl time.Duration) (first bool, err error) {
	idemKey := "ra:idem:" + key
	val := time.Now().UnixMilli()
	if ttl <= 0 {
		ok, err := m.redis.SetNX(ctx, idemKey, val)
		if err != nil {
			return false, err
		}
		return ok, nil
	}
	v, err := m.redis.Do(ctx, "EVAL", acquireScript, 1, idemKey, val, ttl.Milliseconds())
	if err != nil {
		return false, err
	}
	return v.Int() == 1, nil
}

// Remove 主动清除幂等键（如业务失败需回滚时）。
func (m *Manager) Remove(ctx context.Context, key string) error {
	_, err := m.redis.Del(ctx, "ra:idem:"+key)
	return err
}
