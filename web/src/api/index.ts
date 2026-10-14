// 类型安全的接口封装：函数与后端 g.Meta 路由一一对应。
// 类型由 scripts/gen-api.mjs 从 /api.json 自动生成（src/api/types.ts）。
import { http } from './client'
import type {
  LoginReq,
  LoginRes,
  LogoutRes,
  UserInfoRes,
  UserListReq,
  UserListRes,
  UserGetRes,
  UserCreateReq,
  UserCreateRes,
  UserUpdateReq,
  UserUpdateRes,
  UserDeleteRes,
  UserResetPasswordRes,
  UserSetStatusRes,
  RoleListReq,
  RoleListRes,
  RoleGetRes,
  RoleCreateReq,
  RoleCreateRes,
  RoleUpdateReq,
  RoleUpdateRes,
  RoleDeleteRes,
  RoleAssignMenuRes,
  MenuListReq,
  MenuListRes,
  MenuGetRes,
  MenuCreateReq,
  MenuCreateRes,
  MenuUpdateReq,
  MenuUpdateRes,
  MenuDeleteRes,
  DictTypeListReq,
  DictTypeListRes,
  DictTypeCreateReq,
  DictTypeCreateRes,
  DictTypeUpdateReq,
  DictTypeUpdateRes,
  DictTypeDeleteRes,
  DictDataListReq,
  DictDataListRes,
  DictDataCreateReq,
  DictDataCreateRes,
  DictDataUpdateReq,
  DictDataUpdateRes,
  DictDataDeleteRes,
  TenantListReq,
  TenantListRes,
  TenantGetRes,
  TenantCreateReq,
  TenantCreateRes,
  TenantUpdateReq,
  TenantUpdateRes,
  TenantDeleteRes,
  OperLogListReq,
  OperLogListRes,
  OperLogClearRes,
} from './types'

// ---------- 认证 ----------
export const login = (body: LoginReq) => http.post<LoginRes>('/system/login', body)
export const getUserInfo = () => http.get<UserInfoRes>('/system/user/info')
export const logout = () => http.post<LogoutRes>('/system/logout')

// ---------- 用户 ----------
export const listUsers = (q: UserListReq) => http.get<UserListRes>('/system/user', q)
export const getUser = (id: number) => http.get<UserGetRes>(`/system/user/${id}`)
export const createUser = (b: UserCreateReq) => http.post<UserCreateRes>('/system/user', b)
export const updateUser = (id: number, b: Omit<UserUpdateReq, 'id'>) =>
  http.put<UserUpdateRes>(`/system/user/${id}`, b)
export const deleteUser = (id: number) => http.del<UserDeleteRes>(`/system/user/${id}`)
export const resetUserPassword = (id: number, password: string) =>
  http.put<UserResetPasswordRes>(`/system/user/${id}/password`, { password })
export const setUserStatus = (id: number, status: number) =>
  http.put<UserSetStatusRes>(`/system/user/${id}/status`, { status })

// ---------- 角色 ----------
export const listRoles = (q: RoleListReq) => http.get<RoleListRes>('/system/role', q)
export const getRole = (id: number) => http.get<RoleGetRes>(`/system/role/${id}`)
export const createRole = (b: RoleCreateReq) => http.post<RoleCreateRes>('/system/role', b)
export const updateRole = (id: number, b: Omit<RoleUpdateReq, 'id'>) =>
  http.put<RoleUpdateRes>(`/system/role/${id}`, b)
export const deleteRole = (id: number) => http.del<RoleDeleteRes>(`/system/role/${id}`)
export const assignRoleMenus = (id: number, menuIds: number[]) =>
  http.put<RoleAssignMenuRes>(`/system/role/${id}/menu`, { menuIds })

// ---------- 菜单 ----------
export const listMenus = (q?: MenuListReq) => http.get<MenuListRes>('/system/menu', q)
export const getMenu = (id: number) => http.get<MenuGetRes>(`/system/menu/${id}`)
export const createMenu = (b: MenuCreateReq) => http.post<MenuCreateRes>('/system/menu', b)
export const updateMenu = (id: number, b: Omit<MenuUpdateReq, 'id'>) =>
  http.put<MenuUpdateRes>(`/system/menu/${id}`, b)
export const deleteMenu = (id: number) => http.del<MenuDeleteRes>(`/system/menu/${id}`)

// ---------- 字典类型 ----------
export const listDictTypes = (q: DictTypeListReq) => http.get<DictTypeListRes>('/system/dict/type', q)
export const createDictType = (b: DictTypeCreateReq) => http.post<DictTypeCreateRes>('/system/dict/type', b)
export const updateDictType = (id: number, b: Omit<DictTypeUpdateReq, 'id'>) =>
  http.put<DictTypeUpdateRes>(`/system/dict/type/${id}`, b)
export const deleteDictType = (id: number) => http.del<DictTypeDeleteRes>(`/system/dict/type/${id}`)

// ---------- 字典数据 ----------
export const listDictData = (q: DictDataListReq) => http.get<DictDataListRes>('/system/dict/data', q)
export const createDictData = (b: DictDataCreateReq) => http.post<DictDataCreateRes>('/system/dict/data', b)
export const updateDictData = (id: number, b: Omit<DictDataUpdateReq, 'id'>) =>
  http.put<DictDataUpdateRes>(`/system/dict/data/${id}`, b)
export const deleteDictData = (id: number) => http.del<DictDataDeleteRes>(`/system/dict/data/${id}`)

// ---------- 租户 ----------
export const listTenants = (q: TenantListReq) => http.get<TenantListRes>('/system/tenant', q)
export const getTenant = (id: number) => http.get<TenantGetRes>(`/system/tenant/${id}`)
export const createTenant = (b: TenantCreateReq) => http.post<TenantCreateRes>('/system/tenant', b)
export const updateTenant = (id: number, b: Omit<TenantUpdateReq, 'id'>) =>
  http.put<TenantUpdateRes>(`/system/tenant/${id}`, b)
export const deleteTenant = (id: number) => http.del<TenantDeleteRes>(`/system/tenant/${id}`)

// ---------- 操作日志 ----------
export const listOperLogs = (q: OperLogListReq) => http.get<OperLogListRes>('/system/operlog', q)
export const clearOperLogs = () => http.del<OperLogClearRes>('/system/operlog')
