// Package model demo 模块数据模型。
package model

import "github.com/gogf/gf/v2/os/gtime"

// Item 示例数据实体（对应表 ra_demo_item）。
type Item struct {
	Id        int64       `json:"id"        orm:"id,primary"`
	Name      string      `json:"name"      orm:"name"`
	Remark    string      `json:"remark"    orm:"remark"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
