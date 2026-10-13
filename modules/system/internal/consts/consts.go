// Package consts system 模块内部常量：表名、权限码、缓存键、路由分组。
package consts

// 表名（统一 ra_ 前缀）。
const (
	TableUser     = "ra_user"
	TableRole     = "ra_role"
	TableUserRole = "ra_user_role"
	TableMenu     = "ra_menu"
	TableRoleMenu = "ra_role_menu"
	TableDictType = "ra_dict_type"
	TableDictData = "ra_dict_data"
	TableTenant   = "ra_tenant"
	TableOperLog  = "ra_oper_log"
)

// 权限码（RBAC）：登录后写入 Claims.Perms，由 auth.RequirePermission 中间件校验。
const (
	PermUserList     = "system:user:list"
	PermUserCreate   = "system:user:create"
	PermUserUpdate   = "system:user:update"
	PermUserDelete   = "system:user:delete"
	PermRoleList     = "system:role:list"
	PermRoleCreate   = "system:role:create"
	PermRoleUpdate   = "system:role:update"
	PermRoleDelete   = "system:role:delete"
	PermMenuList     = "system:menu:list"
	PermMenuCreate   = "system:menu:create"
	PermMenuUpdate   = "system:menu:update"
	PermMenuDelete   = "system:menu:delete"
	PermDictList     = "system:dict:list"
	PermDictCreate   = "system:dict:create"
	PermDictUpdate   = "system:dict:update"
	PermDictDelete   = "system:dict:delete"
	PermTenantList   = "system:tenant:list"
	PermTenantCreate = "system:tenant:create"
	PermTenantUpdate = "system:tenant:update"
	PermTenantDelete = "system:tenant:delete"
	PermOperLogList  = "system:operlog:list"
	PermOperLogClear = "system:operlog:clear"
)

// 缓存键前缀（分布式锁 / 幂等键）。
const (
	CachePrefixUser   = "ra:system:user:"
	CachePrefixRole   = "ra:system:role:"
	CachePrefixDict   = "ra:system:dict:"
	CachePrefixTenant = "ra:system:tenant:"
)

// 路由分组前缀（g.Meta 路径的前缀，仅作文档约定，实际路由写在 api 层）。
const (
	RoutePrefix = "/system"
)
