package utils

import (
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// TimeRFC3339 将 *gtime.Time 序列化为 RFC3339(UTC) 字符串（如 2026-10-13T12:01:15Z）。
//
// 全局约定：数据库存储 UTC，接口统一返回 RFC3339(UTC)，由前端按用户时区转换显示。
// nil 或零值返回空串。
func TimeRFC3339(t *gtime.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Time.UTC().Format(time.RFC3339)
}
