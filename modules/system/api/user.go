package api

import "github.com/gogf/gf/v2/frame/g"

// UserItem 用户列表项。
type UserItem struct {
	Id        int64  `json:"id"        dc:"User ID"`
	TenantId  int64  `json:"tenantId"  dc:"Tenant ID"`
	Username  string `json:"username"  dc:"Username"`
	Nickname  string `json:"nickname"  dc:"Nickname"`
	Email     string `json:"email"     dc:"Email"`
	Phone     string `json:"phone"     dc:"Phone"`
	Avatar    string `json:"avatar"    dc:"Avatar"`
	Status    int    `json:"status"    dc:"Status (1 enabled, 0 disabled)"`
	Remark    string `json:"remark"    dc:"Remark"`
	CreatedAt string `json:"createdAt" dc:"Creation time"`
}

// UserCreateReq 创建用户请求。
type UserCreateReq struct {
	g.Meta   `path:"/system/user" method:"post" tags:"system" summary:"Create user"`
	Username string  `json:"username" v:"required|length:1,64#username is required|username length must be 1-64" dc:"Username"`
	Password string  `json:"password" v:"required|length:6,64#password is required|password length must be 6-64" dc:"Password"`
	Nickname string  `json:"nickname" v:"required|length:1,64#nickname is required|nickname length must be 1-64" dc:"Nickname"`
	Email    string  `json:"email"    dc:"Email"`
	Phone    string  `json:"phone"    dc:"Phone"`
	Avatar   string  `json:"avatar"   dc:"Avatar"`
	Status   int     `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
	Remark   string  `json:"remark"   dc:"Remark"`
	RoleIds  []int64 `json:"roleIds"  dc:"Role ID list"`
}

// UserCreateRes 创建用户响应。
type UserCreateRes struct {
	Id int64 `json:"id" dc:"New user ID"`
}

// UserListReq 分页查询用户请求。
type UserListReq struct {
	g.Meta   `path:"/system/user" method:"get" tags:"system" summary:"List users"`
	Page     int    `json:"page" d:"1" v:"min:1#page must be at least 1" dc:"Page number"`
	Size     int    `json:"size" d:"20" v:"min:1|max:200#page size must be at least 1|page size must be at most 200" dc:"Page size"`
	Username string `json:"username" dc:"Fuzzy search by username"`
	Nickname string `json:"nickname" dc:"Fuzzy search by nickname"`
	Status   *int   `json:"status" dc:"Filter by status"`
}

// UserListRes 分页查询用户响应。
type UserListRes struct {
	List  []*UserItem `json:"list"  dc:"User list"`
	Total int64       `json:"total" dc:"Total count"`
	Page  int         `json:"page"  dc:"Page number"`
	Size  int         `json:"size"  dc:"Page size"`
}

// UserGetReq 查询用户请求。
type UserGetReq struct {
	g.Meta `path:"/system/user/{id}" method:"get" tags:"system" summary:"Get user"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"User ID"`
}

// UserGetRes 查询用户响应。
type UserGetRes struct {
	User    *UserItem `json:"user"    dc:"User info"`
	RoleIds []int64   `json:"roleIds" dc:"Role ID list"`
}

// UserUpdateReq 更新用户请求。
type UserUpdateReq struct {
	g.Meta   `path:"/system/user/{id}" method:"put" tags:"system" summary:"Update user"`
	Id       int64   `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"User ID"`
	Username string  `json:"username" v:"required|length:1,64#username is required|username length must be 1-64" dc:"Username"`
	Nickname string  `json:"nickname" v:"required|length:1,64#nickname is required|nickname length must be 1-64" dc:"Nickname"`
	Email    string  `json:"email"    dc:"Email"`
	Phone    string  `json:"phone"    dc:"Phone"`
	Avatar   string  `json:"avatar"   dc:"Avatar"`
	Status   int     `json:"status" d:"1" dc:"Status (1 enabled, 0 disabled)"`
	Remark   string  `json:"remark"   dc:"Remark"`
	RoleIds  []int64 `json:"roleIds"  dc:"Role ID list"`
}

// UserUpdateRes 更新用户响应。
type UserUpdateRes struct {
	Updated bool `json:"updated" dc:"Whether updated"`
}

// UserDeleteReq 删除用户请求。
type UserDeleteReq struct {
	g.Meta `path:"/system/user/{id}" method:"delete" tags:"system" summary:"Delete user"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"User ID"`
}

// UserDeleteRes 删除用户响应。
type UserDeleteRes struct {
	Deleted bool `json:"deleted" dc:"Whether deleted"`
}

// UserResetPasswordReq 重置用户密码请求。
type UserResetPasswordReq struct {
	g.Meta   `path:"/system/user/{id}/password" method:"put" tags:"system" summary:"Reset user password"`
	Id       int64  `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"User ID"`
	Password string `json:"password" v:"required|length:6,64#password is required|password length must be 6-64" dc:"New password"`
}

// UserResetPasswordRes 重置用户密码响应。
type UserResetPasswordRes struct {
	Reset bool `json:"reset" dc:"Whether reset"`
}

// UserSetStatusReq 启用/禁用用户请求。
type UserSetStatusReq struct {
	g.Meta `path:"/system/user/{id}/status" method:"put" tags:"system" summary:"Enable or disable user"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID is required|ID is invalid" dc:"User ID"`
	Status int   `json:"status" v:"in:0,1#status must be 0 or 1" dc:"Status (1 enabled, 0 disabled)"`
}

// UserSetStatusRes 启用/禁用用户响应。
type UserSetStatusRes struct {
	Updated bool `json:"updated" dc:"Whether updated"`
}
