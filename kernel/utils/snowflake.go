package utils

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// 雪花算法：41bit 毫秒时间戳 + 10bit 节点 ID + 12bit 序列号。
// 用于生成分布式友好的全局唯一 bigint 主键（跨库一致、免自增）。
const (
	snowflakeEpoch     = int64(1704067200000) // 2024-01-01T00:00:00Z
	snowflakeNodeBits  = 10
	snowflakeSeqBits   = 12
	snowflakeMaxNode   = -1 ^ (-1 << snowflakeNodeBits)
	snowflakeMaxSeq    = -1 ^ (-1 << snowflakeSeqBits)
	snowflakeNodeShift = snowflakeSeqBits
	snowflakeTimeShift = snowflakeSeqBits + snowflakeNodeBits
)

// Snowflake 线程安全的雪花 ID 生成器。
type Snowflake struct {
	mu     sync.Mutex
	nodeID int64
	seq    int64
	lastTs int64
}

// NewSnowflake 创建生成器，nodeID 超范围时自动取模。
func NewSnowflake(nodeID int64) *Snowflake {
	if nodeID < 0 {
		nodeID = -nodeID
	}
	return &Snowflake{nodeID: nodeID & snowflakeMaxNode}
}

// Next 返回下一个全局唯一 ID。
func (s *Snowflake) Next() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()
	// 时钟回拨保护：不回退，沿用上次时间戳。
	if now < s.lastTs {
		now = s.lastTs
	}

	if now == s.lastTs {
		s.seq = (s.seq + 1) & snowflakeMaxSeq
		if s.seq == 0 {
			// 本毫秒序列号耗尽，自旋等待到下一毫秒。
			for now <= s.lastTs {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.seq = 0
	}
	s.lastTs = now

	return ((now - snowflakeEpoch) << snowflakeTimeShift) |
		(s.nodeID << snowflakeNodeShift) |
		s.seq
}

var (
	defaultSnowflakeOnce sync.Once
	defaultSnowflake     *Snowflake
)

// NextID 返回默认实例生成的 ID（懒初始化：节点号读配置 snowflake.nodeId，缺省 1）。
//
// 多实例部署时为避免主键冲突，每个实例需配置不同的 snowflake.nodeId（0..1023）。
func NextID() int64 {
	defaultSnowflakeOnce.Do(func() {
		node := g.Cfg().MustGet(context.Background(), "snowflake.nodeId", 1).Int()
		defaultSnowflake = NewSnowflake(int64(node))
	})
	return defaultSnowflake.Next()
}
