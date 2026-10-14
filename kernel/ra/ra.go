// Package ra 是框架唯一的门面（facade）。
//
// 内核所有单例通过本包对外暴露，业务模块面向 ra 编程：
//
//	ra.DB()、ra.Log().Info(ctx, "...")、ra.Cfg().MustGet(ctx, "app.name")、ra.Auth() …
//
// 约定：内核各子包互不 import ra（避免环），ra 是依赖方向的最顶端。
// 需要生命周期管理的单例（事件总线/消息/观测/鉴权）采用懒初始化，首次使用时创建。
package ra

import (
	"sync"

	"github.com/roamdream/roamadmin/kernel/auth"
	"github.com/roamdream/roamadmin/kernel/cache"
	"github.com/roamdream/roamadmin/kernel/config"
	"github.com/roamdream/roamadmin/kernel/db"
	"github.com/roamdream/roamadmin/kernel/event"
	"github.com/roamdream/roamadmin/kernel/idempotent"
	"github.com/roamdream/roamadmin/kernel/limit"
	"github.com/roamdream/roamadmin/kernel/lock"
	"github.com/roamdream/roamadmin/kernel/logger"
	"github.com/roamdream/roamadmin/kernel/mq"
	"github.com/roamdream/roamadmin/kernel/mqtt"
	"github.com/roamdream/roamadmin/kernel/observ"
	"github.com/roamdream/roamadmin/kernel/pool"
	"github.com/roamdream/roamadmin/kernel/tsdb"
)

// DB 数据库（gdb 封装）。
func DB() *db.Manager { return db.Default() }

// Log 日志。
func Log() *logger.Logger { return logger.Default() }

// Cfg 配置。
func Cfg() *config.Manager { return config.Default() }

// Cache Redis 缓存。
func Cache() *cache.Cache { return cache.Default() }

// LocalCache 本地内存缓存。
func LocalCache() *cache.Local { return cache.NewLocal() }

// Lock 分布式锁。
func Lock() *lock.Locker { return lock.Default() }

// Idempotent 幂等管理器。
func Idempotent() *idempotent.Manager { return idempotent.Default() }

// Pool 协程池。
func Pool() *pool.Pool { return pool.Default() }

// Limiter Redis 分布式限流器。
func Limiter() *limit.RedisLimiter { return limit.Default() }

var (
	eventOnce sync.Once
	eventBus  *event.Bus
)

// Event 事件总线（懒初始化）。
func Event() *event.Bus {
	eventOnce.Do(func() { eventBus = event.New() })
	return eventBus
}

var (
	mqOnce  sync.Once
	mqMgr   mq.Manager
)

// MQ 消息管理器（默认本地实现）。
func MQ() mq.Manager {
	mqOnce.Do(func() { mqMgr = mq.NewLocal() })
	return mqMgr
}

var (
	observOnce sync.Once
	observMgr  *observ.Manager
)

// Observ 可观测管理器（metrics/tracing）。
func Observ() *observ.Manager {
	observOnce.Do(func() { observMgr = observ.New() })
	return observMgr
}

var (
	authOnce sync.Once
	authMgr  *auth.Manager
)

// Auth 认证鉴权管理器（懒初始化，读配置 auth.jwt.*）。
func Auth() *auth.Manager {
	authOnce.Do(func() { authMgr = auth.Default() })
	return authMgr
}

var (
	tsdbOnce  sync.Once
	tsdbEngine tsdb.Engine
)

// TSDB 时序存储引擎（懒初始化）。
func TSDB() tsdb.Engine {
	tsdbOnce.Do(func() { tsdbEngine = tsdb.Default() })
	return tsdbEngine
}

var (
	mqttOnce   sync.Once
	mqttBroker *mqtt.Broker
)

// MQTT 内嵌 MQTT Broker（懒初始化，由 gateway 模块启动）。
func MQTT() *mqtt.Broker {
	mqttOnce.Do(func() { mqttBroker = mqtt.New() })
	return mqttBroker
}
