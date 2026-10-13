// Package model system 模块数据模型。
package model

import "github.com/gogf/gf/v2/os/gtime"

// User 用户实体（对应表 sys_user，租户隔离）。
type User struct {
	Id        int64       `json:"id"        orm:"id,primary"`
	TenantId  int64       `json:"tenantId"  orm:"tenant_id"`
	Username  string      `json:"username"  orm:"username"`
	Password  string      `json:"-"         orm:"password"`
	Nickname  string      `json:"nickname"  orm:"nickname"`
	Email     string      `json:"email"     orm:"email"`
	Phone     string      `json:"phone"     orm:"phone"`
	Avatar    string      `json:"avatar"    orm:"avatar"`
	Status    int         `json:"status"    orm:"status"`
	Remark    string      `json:"remark"    orm:"remark"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}

// UserInfo 用户聚合信息（基本信息 + 角色编码 + 权限码）。
type UserInfo struct {
	User  *User    `json:"user"`
	Roles []string `json:"roles"`
	Perms []string `json:"perms"`
}

// Role 角色实体（对应表 sys_role，租户隔离）。
type Role struct {
	Id        int64       `json:"id"        orm:"id,primary"`
	TenantId  int64       `json:"tenantId"  orm:"tenant_id"`
	Name      string      `json:"name"      orm:"name"`
	Code      string      `json:"code"      orm:"code"`
	Status    int         `json:"status"    orm:"status"`
	Remark    string      `json:"remark"    orm:"remark"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}

// Menu 菜单/权限实体（对应表 sys_menu，租户隔离，树形）。
type Menu struct {
	Id        int64       `json:"id"        orm:"id,primary"`
	TenantId  int64       `json:"tenantId"  orm:"tenant_id"`
	ParentId  int64       `json:"parentId"  orm:"parent_id"`
	Name      string      `json:"name"      orm:"name"`
	Type      string      `json:"type"      orm:"type"`
	Perm      string      `json:"perm"      orm:"perm"`
	Path      string      `json:"path"      orm:"path"`
	Component string      `json:"component" orm:"component"`
	Icon      string      `json:"icon"      orm:"icon"`
	Sort      int         `json:"sort"      orm:"sort"`
	Status    int         `json:"status"    orm:"status"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}

// MenuNode 菜单树节点。
type MenuNode struct {
	*Menu
	Children []*MenuNode `json:"children"`
}

// DictType 字典类型实体（对应表 sys_dict_type，租户隔离）。
type DictType struct {
	Id        int64       `json:"id"        orm:"id,primary"`
	TenantId  int64       `json:"tenantId"  orm:"tenant_id"`
	Name      string      `json:"name"      orm:"name"`
	Type      string      `json:"type"      orm:"type"`
	Status    int         `json:"status"    orm:"status"`
	Remark    string      `json:"remark"    orm:"remark"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}

// DictData 字典数据实体（对应表 sys_dict_data，租户隔离）。
type DictData struct {
	Id        int64       `json:"id"        orm:"id,primary"`
	TenantId  int64       `json:"tenantId"  orm:"tenant_id"`
	DictType  string      `json:"dictType"  orm:"dict_type"`
	Label     string      `json:"label"     orm:"label"`
	Value     string      `json:"value"     orm:"value"`
	Sort      int         `json:"sort"      orm:"sort"`
	Status    int         `json:"status"    orm:"status"`
	Remark    string      `json:"remark"    orm:"remark"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}

// Tenant 租户实体（对应表 sys_tenant，租户注册表本身，无 tenant_id）。
type Tenant struct {
	Id        int64       `json:"id"        orm:"id,primary"`
	Name      string      `json:"name"      orm:"name"`
	Code      string      `json:"code"      orm:"code"`
	Status    int         `json:"status"    orm:"status"`
	Remark    string      `json:"remark"    orm:"remark"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"`
}

// OperLog 操作日志实体（对应表 sys_oper_log，租户隔离）。
type OperLog struct {
	Id        int64       `json:"id"        orm:"id,primary"`
	TenantId  int64       `json:"tenantId"  orm:"tenant_id"`
	UserId    int64       `json:"userId"    orm:"user_id"`
	Username  string      `json:"username"  orm:"username"`
	Module    string      `json:"module"    orm:"module"`
	Action    string      `json:"action"    orm:"action"`
	Method    string      `json:"method"    orm:"method"`
	Path      string      `json:"path"      orm:"path"`
	IP        string      `json:"ip"        orm:"ip"`
	Status    int         `json:"status"    orm:"status"`
	ErrMsg    string      `json:"errMsg"    orm:"err_msg"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
