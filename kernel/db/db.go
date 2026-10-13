// Package db 数据库封装（gdb ORM）。
//
// 内核在此定死"真兼容"策略：全程走 ORM 禁写方言 SQL、主键统一雪花 ID、
// 读写分离/多数据源由 gdb 配置组承载（换库只改一处配置）。
// 连接池高并发推荐值统一收敛到 manifest/config/config.yaml。
package db

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/roamdream/roamadmin/kernel/utils"
)

// Manager gdb 封装。
type Manager struct {
	inner gdb.DB
}

// New 包装一个 gdb.DB。
func New(inner gdb.DB) *Manager { return &Manager{inner: inner} }

// Default 返回默认数据源组（等价 g.DB()）。
func Default() *Manager { return &Manager{inner: g.DB()} }

// DB 返回底层 gdb.DB（供高级用法）。
func (m *Manager) DB() gdb.DB { return m.inner }

// Instance 返回命名数据源组（多数据源 / 读写分离）。
func (m *Manager) Instance(group ...string) gdb.DB { return g.DB(group...) }

// Model 返回指定表的 ORM Model，禁止在其中手写方言 SQL。
func (m *Manager) Model(table string) *gdb.Model { return m.inner.Model(table) }

// NextID 生成雪花主键（bigint），跨库一致、免自增。
func (m *Manager) NextID() int64 { return utils.NextID() }

// Transaction 统一事务入口：自动 commit/rollback，杜绝遗漏。
func (m *Manager) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return m.inner.Transaction(ctx, f)
}
