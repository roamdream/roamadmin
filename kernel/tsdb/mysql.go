package tsdb

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/roamdream/roamadmin/kernel/utils"
)

// tablePoint 时序数据点表（方言无关 DDL，主键统一雪花 ID）。
const tablePoint = "ra_ts_point"

// MySQL 基于关系库（MySQL 8 / PostgreSQL）的时序存储参考实现。
//
// P3 阶段用它打通端到端闭环；高并发/大规模场景换 TDengine / ClickHouse 时，
// 内核与业务模块仍只面向 Engine 接口编程，零改动切换。
type MySQL struct {
	db gdb.DB
}

// NewMySQL 创建基于关系库的时序引擎。
func NewMySQL(db gdb.DB) *MySQL { return &MySQL{db: db} }

// Default 返回默认时序引擎（等价 g.DB()）。
func Default() Engine { return NewMySQL(g.DB()) }

// ensureSchema 幂等建表（DDL 不包含反引号 / 方言特有语法）。
func (m *MySQL) ensureSchema(ctx context.Context) error {
	_, err := m.db.Exec(ctx, `
CREATE TABLE IF NOT EXISTS ra_ts_point (
    id          BIGINT       NOT NULL,
    measurement VARCHAR(255) NOT NULL,
    tags        TEXT,
    fields      TEXT         NOT NULL,
    ts          DATETIME     NOT NULL,
    PRIMARY KEY (id)
)`)
	return err
}

// Write 批量写入数据点（单条多行 INSERT，提升高并发写入吞吐）。
func (m *MySQL) Write(ctx context.Context, points ...*Point) error {
	if len(points) == 0 {
		return nil
	}
	if err := m.ensureSchema(ctx); err != nil {
		return err
	}
	rows := make([]g.Map, 0, len(points))
	for _, p := range points {
		if p == nil {
			continue
		}
		tagsJSON, _ := json.Marshal(p.Tags)
		fieldsJSON, err := json.Marshal(p.Fields)
		if err != nil {
			return err
		}
		rows = append(rows, g.Map{
			"id":          utils.NextID(),
			"measurement": p.Measurement,
			"tags":        string(tagsJSON),
			"fields":      string(fieldsJSON),
			"ts":          gtime.New(p.Time.UTC()),
		})
	}
	if len(rows) == 0 {
		return nil
	}
	_, err := m.db.Model(tablePoint).Data(rows).Insert()
	return err
}

// Query 执行原生查询，返回 gdb.Result（实现相关原始结果）。
func (m *MySQL) Query(ctx context.Context, statement string, args ...any) (any, error) {
	if err := m.ensureSchema(ctx); err != nil {
		return nil, err
	}
	return m.db.GetAll(ctx, statement, args...)
}

// QueryRecent 查询某 measurement 最近 limit 个数据点（时间倒序）。
func (m *MySQL) QueryRecent(ctx context.Context, measurement string, limit int) ([]*Point, error) {
	if err := m.ensureSchema(ctx); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	var rows []struct {
		Tags   string      `orm:"tags"`
		Fields string      `orm:"fields"`
		Ts     *gtime.Time `orm:"ts"`
	}
	if err := m.db.Model(tablePoint).
		Where("measurement", measurement).
		OrderDesc("ts").
		Limit(limit).
		Scan(&rows); err != nil {
		return nil, err
	}
	out := make([]*Point, 0, len(rows))
	for _, r := range rows {
		var tags map[string]string
		var fields map[string]any
		_ = json.Unmarshal([]byte(r.Tags), &tags)
		_ = json.Unmarshal([]byte(r.Fields), &fields)
		out = append(out, &Point{
			Measurement: measurement,
			Tags:        tags,
			Fields:      fields,
			Time:        r.Ts.Time,
		})
	}
	return out, nil
}

// Close 关闭连接（关系库由 gdb 统一管理，无需单独关闭）。
func (m *MySQL) Close(ctx context.Context) error { return nil }
