// Package limit 限流 / 熔断 / 降级（保护机制，防雪崩）。
//
// 提供三种开箱即用能力：本地令牌桶、Redis 分布式固定窗口限流、本地熔断器。
package limit

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
)

// allowScript 原子限流：INCR + 首次 PEXPIRE 在一条 Lua 内完成，避免 INCR 与 EXPIRE
// 分两步导致首次计数后窗口永不重置；返回当前窗口计数。
const allowScript = `local n = redis.call("incr", KEYS[1]) if n == 1 then redis.call("pexpire", KEYS[1], ARGV[1]) end return n`

// TokenBucket 本地令牌桶限流器（单机）。
type TokenBucket struct {
	mu     sync.Mutex
	rate   float64 // 每秒补充令牌数
	burst  float64 // 桶容量
	tokens float64
	last   time.Time
}

// NewTokenBucket 创建令牌桶：rate 为每秒令牌数，burst 为桶容量。
func NewTokenBucket(rate, burst float64) *TokenBucket {
	if burst <= 0 {
		burst = rate
	}
	return &TokenBucket{rate: rate, burst: burst, tokens: burst, last: time.Now()}
}

// Allow 消耗 1 个令牌。
func (b *TokenBucket) Allow() bool { return b.AllowN(1) }

// AllowN 消耗 n 个令牌。
func (b *TokenBucket) AllowN(n float64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	b.tokens += now.Sub(b.last).Seconds() * b.rate
	if b.tokens > b.burst {
		b.tokens = b.burst
	}
	b.last = now
	if b.tokens >= n {
		b.tokens -= n
		return true
	}
	return false
}

// RedisLimiter Redis 分布式固定窗口限流（跨实例一致）。
type RedisLimiter struct {
	redis *gredis.Redis
}

// NewRedisLimiter 包装一个 gredis.Redis。
func NewRedisLimiter(r *gredis.Redis) *RedisLimiter { return &RedisLimiter{redis: r} }

// Default 返回默认 Redis 限流器。
func Default() *RedisLimiter { return &RedisLimiter{redis: g.Redis()} }

// Allow 固定窗口限流：window 时间窗口内最多 limit 次。
func (l *RedisLimiter) Allow(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	redisKey := "ra:limit:" + key
	v, err := l.redis.Do(ctx, "EVAL", allowScript, 1, redisKey, window.Milliseconds())
	if err != nil {
		return false, err
	}
	return v.Int64() <= limit, nil
}

// BreakerState 熔断器状态。
type BreakerState int

const (
	StateClosed   BreakerState = iota // 关闭（正常放行）
	StateOpen                          // 打开（熔断，拒绝）
	StateHalfOpen                      // 半开（放行一个探测请求）
)

// Breaker 本地熔断器：连续失败达阈值后熔断，超时后进入半开试探。
type Breaker struct {
	mu        sync.Mutex
	state     BreakerState
	failures  int
	threshold int
	timeout   time.Duration
	openedAt  time.Time
}

// NewBreaker 创建熔断器：连续 failureThreshold 次失败熔断，timeout 后进入半开。
func NewBreaker(failureThreshold int, timeout time.Duration) *Breaker {
	return &Breaker{state: StateClosed, threshold: failureThreshold, timeout: timeout}
}

// State 当前状态。
func (b *Breaker) State() BreakerState {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// Allow 检查是否允许请求通过；半开状态只放行一个探测请求。
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(b.openedAt) >= b.timeout {
			b.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return false
	}
	return false
}

// Success 上报成功，关闭熔断器。
func (b *Breaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// Failure 上报失败，达到阈值时打开熔断器。
func (b *Breaker) Failure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.state == StateHalfOpen || b.failures >= b.threshold {
		b.state = StateOpen
		b.openedAt = time.Now()
	}
}
