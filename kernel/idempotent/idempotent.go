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

// Manager 幂等管理器。
type Manager struct {
	redis *gredis.Redis
}

// New 包装一个 gredis.Redis。
func New(r *gredis.Redis) *Manager { return &Manager{redis: r} }

// Default 返回默认 Manager。
func Default() *Manager { return &Manager{redis: g.Redis()} }

// Check 检查幂等键：首次出现则占用并返回 true，重复出现返回 false。
// ttl 控制幂等键的存活时间（即幂等窗口）。
func (m *Manager) Check(ctx context.Context, key string, ttl time.Duration) (first bool, err error) {
	idemKey := "ra:idem:" + key
	ok, err := m.redis.SetNX(ctx, idemKey, time.Now().UnixMilli())
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	if ttl > 0 {
		_, _ = m.redis.Expire(ctx, idemKey, int64(ttl.Seconds()))
	}
	return true, nil
}

// Remove 主动清除幂等键（如业务失败需回滚时）。
func (m *Manager) Remove(ctx context.Context, key string) error {
	_, err := m.redis.Del(ctx, "ra:idem:"+key)
	return err
}
