package api

import "github.com/gogf/gf/v2/frame/g"

// UserItem 用户列表项。
type UserItem struct {
	Id        int64  `json:"id"        dc:"用户 ID"`
	TenantId  int64  `json:"tenantId"  dc:"租户 ID"`
	Username  string `json:"username"  dc:"用户名"`
	Nickname  string `json:"nickname"  dc:"昵称"`
	Email     string `json:"email"     dc:"邮箱"`
	Phone     string `json:"phone"     dc:"手机号"`
	Avatar    string `json:"avatar"    dc:"头像"`
	Status    int    `json:"status"    dc:"状态（1启用 0禁用）"`
	Remark    string `json:"remark"    dc:"备注"`
	CreatedAt string `json:"createdAt" dc:"创建时间"`
}

// UserCreateReq 创建用户请求。
type UserCreateReq struct {
	g.Meta   `path:"/system/user" method:"post" tags:"系统" summary:"创建用户"`
	Username string  `json:"username" v:"required|length:1,64#用户名不能为空|用户名长度1-64" dc:"用户名"`
	Password string  `json:"password" v:"required|length:6,64#密码不能为空|密码长度6-64" dc:"密码"`
	Nickname string  `json:"nickname" v:"required|length:1,64#昵称不能为空|昵称长度1-64" dc:"昵称"`
	Email    string  `json:"email"    dc:"邮箱"`
	Phone    string  `json:"phone"    dc:"手机号"`
	Avatar   string  `json:"avatar"   dc:"头像"`
	Status   int     `json:"status" d:"1" dc:"状态（1启用 0禁用）"`
	Remark   string  `json:"remark"   dc:"备注"`
	RoleIds  []int64 `json:"roleIds"  dc:"角色 ID 列表"`
}

// UserCreateRes 创建用户响应。
type UserCreateRes struct {
	Id int64 `json:"id" dc:"新用户 ID"`
}

// UserListReq 分页查询用户请求。
type UserListReq struct {
	g.Meta   `path:"/system/user" method:"get" tags:"系统" summary:"分页查询用户"`
	Page     int    `json:"page" d:"1" v:"min:1#页码最小为1" dc:"页码"`
	Size     int    `json:"size" d:"20" v:"min:1|max:200#每页最小为1|每页最大为200" dc:"每页条数"`
	Username string `json:"username" dc:"用户名模糊搜索"`
	Nickname string `json:"nickname" dc:"昵称模糊搜索"`
	Status   *int   `json:"status" dc:"状态筛选"`
}

// UserListRes 分页查询用户响应。
type UserListRes struct {
	List  []*UserItem `json:"list"  dc:"用户列表"`
	Total int64       `json:"total" dc:"总数"`
	Page  int         `json:"page"  dc:"页码"`
	Size  int         `json:"size"  dc:"每页条数"`
}

// UserGetReq 查询用户请求。
type UserGetReq struct {
	g.Meta `path:"/system/user/{id}" method:"get" tags:"系统" summary:"查询用户"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"用户 ID"`
}

// UserGetRes 查询用户响应。
type UserGetRes struct {
	User    *UserItem `json:"user"    dc:"用户信息"`
	RoleIds []int64   `json:"roleIds" dc:"角色 ID 列表"`
}

// UserUpdateReq 更新用户请求。
type UserUpdateReq struct {
	g.Meta   `path:"/system/user/{id}" method:"put" tags:"系统" summary:"更新用户"`
	Id       int64   `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"用户 ID"`
	Username string  `json:"username" v:"required|length:1,64#用户名不能为空|用户名长度1-64" dc:"用户名"`
	Nickname string  `json:"nickname" v:"required|length:1,64#昵称不能为空|昵称长度1-64" dc:"昵称"`
	Email    string  `json:"email"    dc:"邮箱"`
	Phone    string  `json:"phone"    dc:"手机号"`
	Avatar   string  `json:"avatar"   dc:"头像"`
	Status   int     `json:"status" d:"1" dc:"状态（1启用 0禁用）"`
	Remark   string  `json:"remark"   dc:"备注"`
	RoleIds  []int64 `json:"roleIds"  dc:"角色 ID 列表"`
}

// UserUpdateRes 更新用户响应。
type UserUpdateRes struct {
	Updated bool `json:"updated" dc:"是否已更新"`
}

// UserDeleteReq 删除用户请求。
type UserDeleteReq struct {
	g.Meta `path:"/system/user/{id}" method:"delete" tags:"系统" summary:"删除用户"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"用户 ID"`
}

// UserDeleteRes 删除用户响应。
type UserDeleteRes struct {
	Deleted bool `json:"deleted" dc:"是否已删除"`
}

// UserResetPasswordReq 重置用户密码请求。
type UserResetPasswordReq struct {
	g.Meta   `path:"/system/user/{id}/password" method:"put" tags:"系统" summary:"重置用户密码"`
	Id       int64  `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"用户 ID"`
	Password string `json:"password" v:"required|length:6,64#密码不能为空|密码长度6-64" dc:"新密码"`
}

// UserResetPasswordRes 重置用户密码响应。
type UserResetPasswordRes struct {
	Reset bool `json:"reset" dc:"是否已重置"`
}

// UserSetStatusReq 启用/禁用用户请求。
type UserSetStatusReq struct {
	g.Meta `path:"/system/user/{id}/status" method:"put" tags:"系统" summary:"启用/禁用用户"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"用户 ID"`
	Status int   `json:"status" v:"in:0,1#状态只能为0或1" dc:"状态（1启用 0禁用）"`
}

// UserSetStatusRes 启用/禁用用户响应。
type UserSetStatusRes struct {
	Updated bool `json:"updated" dc:"是否已更新"`
}
