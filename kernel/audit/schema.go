package audit

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/roamdream/roamadmin/kernel/db"
)

// ensureSchema 幂等建表 + 补列（兼容旧表升级）。DDL 走方言中立语法，MySQL/PostgreSQL 通用。
func (m *Manager) ensureSchema(ctx context.Context) error {
	if err := m.ensureTable(ctx); err != nil {
		return err
	}
	for _, c := range []struct{ name, ddl string }{
		{"query", "VARCHAR(1024) NOT NULL DEFAULT ''"},
		{"user_agent", "VARCHAR(512) NOT NULL DEFAULT ''"},
		{"referer", "VARCHAR(512) NOT NULL DEFAULT ''"},
		{"resp_code", "INT NOT NULL DEFAULT 0"},
		{"duration_ms", "BIGINT NOT NULL DEFAULT 0"},
		{"target_type", "VARCHAR(64) NOT NULL DEFAULT ''"},
		{"target_id", "BIGINT NOT NULL DEFAULT 0"},
		{"trace_id", "VARCHAR(64) NOT NULL DEFAULT ''"},
		{"req_body", "TEXT NULL"},
		{"detail", "TEXT NULL"},
	} {
		m.ensureColumn(ctx, c.name, c.ddl)
	}
	return nil
}

func (m *Manager) ensureTable(ctx context.Context) error {
	ddl := `CREATE TABLE IF NOT EXISTS ` + tableOperLog + ` (
		id BIGINT NOT NULL PRIMARY KEY,
		tenant_id BIGINT NOT NULL DEFAULT 0,
		user_id BIGINT NOT NULL DEFAULT 0,
		username VARCHAR(64) NOT NULL DEFAULT '',
		module VARCHAR(64) NOT NULL DEFAULT '',
		action VARCHAR(64) NOT NULL DEFAULT '',
		method VARCHAR(16) NOT NULL DEFAULT '',
		path VARCHAR(255) NOT NULL DEFAULT '',
		query VARCHAR(1024) NOT NULL DEFAULT '',
		ip VARCHAR(64) NOT NULL DEFAULT '',
		user_agent VARCHAR(512) NOT NULL DEFAULT '',
		referer VARCHAR(512) NOT NULL DEFAULT '',
		status INT NOT NULL DEFAULT 0,
		resp_code INT NOT NULL DEFAULT 0,
		err_msg VARCHAR(512) NOT NULL DEFAULT '',
		duration_ms BIGINT NOT NULL DEFAULT 0,
		target_type VARCHAR(64) NOT NULL DEFAULT '',
		target_id BIGINT NOT NULL DEFAULT 0,
		trace_id VARCHAR(64) NOT NULL DEFAULT '',
		req_body TEXT NULL,
		detail TEXT NULL,
		created_at TIMESTAMP NOT NULL
	)`
	if _, err := db.Default().DB().Exec(ctx, ddl); err != nil {
		return fmt.Errorf("create table: %w", err)
	}
	return nil
}

// ensureColumn 补列：列已存在时忽略（方言中立的 try-and-ignore，避免 information_schema 差异）。
func (m *Manager) ensureColumn(ctx context.Context, name, ddlType string) {
	stmt := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableOperLog, name, ddlType)
	if _, err := db.Default().DB().Exec(ctx, stmt); err != nil && !isDuplicateColumn(err) {
		g.Log().Errorf(ctx, "[ra][audit] add column %s failed: %v", name, err)
	}
}

func isDuplicateColumn(err error) bool {
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "duplicate column") || strings.Contains(s, "already exists")
}
