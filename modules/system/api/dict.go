package api

import "github.com/gogf/gf/v2/frame/g"

// DictTypeItem 字典类型列表项。
type DictTypeItem struct {
	Id        int64  `json:"id"        dc:"字典类型 ID"`
	TenantId  int64  `json:"tenantId"  dc:"租户 ID"`
	Name      string `json:"name"      dc:"字典名称"`
	Type      string `json:"type"      dc:"字典类型编码"`
	Status    int    `json:"status"    dc:"状态（1启用 0禁用）"`
	Remark    string `json:"remark"    dc:"备注"`
	CreatedAt string `json:"createdAt" dc:"创建时间"`
}

// DictTypeCreateReq 创建字典类型请求。
type DictTypeCreateReq struct {
	g.Meta `path:"/system/dict/type" method:"post" tags:"系统" summary:"创建字典类型"`
	Name   string `json:"name" v:"required|length:1,64#字典名称不能为空|字典名称长度1-64" dc:"字典名称"`
	Type   string `json:"type" v:"required|length:1,64#字典类型不能为空|字典类型长度1-64" dc:"字典类型编码"`
	Status int    `json:"status" d:"1" dc:"状态（1启用 0禁用）"`
	Remark string `json:"remark" dc:"备注"`
}

// DictTypeCreateRes 创建字典类型响应。
type DictTypeCreateRes struct {
	Id int64 `json:"id" dc:"新字典类型 ID"`
}

// DictTypeListReq 分页查询字典类型请求。
type DictTypeListReq struct {
	g.Meta `path:"/system/dict/type" method:"get" tags:"系统" summary:"分页查询字典类型"`
	Page   int    `json:"page" d:"1" v:"min:1#页码最小为1" dc:"页码"`
	Size   int    `json:"size" d:"20" v:"min:1|max:200#每页最小为1|每页最大为200" dc:"每页条数"`
	Name   string `json:"name" dc:"字典名称模糊搜索"`
}

// DictTypeListRes 分页查询字典类型响应。
type DictTypeListRes struct {
	List  []*DictTypeItem `json:"list"  dc:"字典类型列表"`
	Total int64           `json:"total" dc:"总数"`
	Page  int             `json:"page"  dc:"页码"`
	Size  int             `json:"size"  dc:"每页条数"`
}

// DictTypeGetReq 查询字典类型请求。
type DictTypeGetReq struct {
	g.Meta `path:"/system/dict/type/{id}" method:"get" tags:"系统" summary:"查询字典类型"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"字典类型 ID"`
}

// DictTypeGetRes 查询字典类型响应。
type DictTypeGetRes struct {
	DictType *DictTypeItem `json:"dictType" dc:"字典类型"`
}

// DictTypeUpdateReq 更新字典类型请求。
type DictTypeUpdateReq struct {
	g.Meta `path:"/system/dict/type/{id}" method:"put" tags:"系统" summary:"更新字典类型"`
	Id     int64  `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"字典类型 ID"`
	Name   string `json:"name" v:"required|length:1,64#字典名称不能为空|字典名称长度1-64" dc:"字典名称"`
	Status int    `json:"status" d:"1" dc:"状态（1启用 0禁用）"`
	Remark string `json:"remark" dc:"备注"`
}

// DictTypeUpdateRes 更新字典类型响应。
type DictTypeUpdateRes struct {
	Updated bool `json:"updated" dc:"是否已更新"`
}

// DictTypeDeleteReq 删除字典类型请求。
type DictTypeDeleteReq struct {
	g.Meta `path:"/system/dict/type/{id}" method:"delete" tags:"系统" summary:"删除字典类型"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"字典类型 ID"`
}

// DictTypeDeleteRes 删除字典类型响应。
type DictTypeDeleteRes struct {
	Deleted bool `json:"deleted" dc:"是否已删除"`
}

// DictDataItem 字典数据列表项。
type DictDataItem struct {
	Id       int64  `json:"id"       dc:"字典数据 ID"`
	TenantId int64  `json:"tenantId" dc:"租户 ID"`
	DictType string `json:"dictType" dc:"字典类型编码"`
	Label    string `json:"label"    dc:"标签"`
	Value    string `json:"value"    dc:"值"`
	Sort     int    `json:"sort"     dc:"排序"`
	Status   int    `json:"status"   dc:"状态（1启用 0禁用）"`
	Remark   string `json:"remark"   dc:"备注"`
}

// DictDataCreateReq 创建字典数据请求。
type DictDataCreateReq struct {
	g.Meta   `path:"/system/dict/data" method:"post" tags:"系统" summary:"创建字典数据"`
	DictType string `json:"dictType" v:"required|length:1,64#字典类型不能为空|字典类型长度1-64" dc:"字典类型编码"`
	Label    string `json:"label" v:"required|length:1,64#标签不能为空|标签长度1-64" dc:"标签"`
	Value    string `json:"value" v:"required|length:1,64#值不能为空|值长度1-64" dc:"值"`
	Sort     int    `json:"sort" d:"0" dc:"排序"`
	Status   int    `json:"status" d:"1" dc:"状态（1启用 0禁用）"`
	Remark   string `json:"remark" dc:"备注"`
}

// DictDataCreateRes 创建字典数据响应。
type DictDataCreateRes struct {
	Id int64 `json:"id" dc:"新字典数据 ID"`
}

// DictDataListReq 分页查询字典数据请求。
type DictDataListReq struct {
	g.Meta   `path:"/system/dict/data" method:"get" tags:"系统" summary:"分页查询字典数据"`
	Page     int    `json:"page" d:"1" v:"min:1#页码最小为1" dc:"页码"`
	Size     int    `json:"size" d:"20" v:"min:1|max:200#每页最小为1|每页最大为200" dc:"每页条数"`
	DictType string `json:"dictType" dc:"字典类型编码筛选"`
	Label    string `json:"label" dc:"标签模糊搜索"`
}

// DictDataListRes 分页查询字典数据响应。
type DictDataListRes struct {
	List  []*DictDataItem `json:"list"  dc:"字典数据列表"`
	Total int64           `json:"total" dc:"总数"`
	Page  int             `json:"page"  dc:"页码"`
	Size  int             `json:"size"  dc:"每页条数"`
}

// DictDataGetReq 查询字典数据请求。
type DictDataGetReq struct {
	g.Meta `path:"/system/dict/data/{id}" method:"get" tags:"系统" summary:"查询字典数据"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"字典数据 ID"`
}

// DictDataGetRes 查询字典数据响应。
type DictDataGetRes struct {
	DictData *DictDataItem `json:"dictData" dc:"字典数据"`
}

// DictDataUpdateReq 更新字典数据请求。
type DictDataUpdateReq struct {
	g.Meta   `path:"/system/dict/data/{id}" method:"put" tags:"系统" summary:"更新字典数据"`
	Id       int64  `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"字典数据 ID"`
	DictType string `json:"dictType" v:"required|length:1,64#字典类型不能为空|字典类型长度1-64" dc:"字典类型编码"`
	Label    string `json:"label" v:"required|length:1,64#标签不能为空|标签长度1-64" dc:"标签"`
	Value    string `json:"value" v:"required|length:1,64#值不能为空|值长度1-64" dc:"值"`
	Sort     int    `json:"sort" d:"0" dc:"排序"`
	Status   int    `json:"status" d:"1" dc:"状态（1启用 0禁用）"`
	Remark   string `json:"remark" dc:"备注"`
}

// DictDataUpdateRes 更新字典数据响应。
type DictDataUpdateRes struct {
	Updated bool `json:"updated" dc:"是否已更新"`
}

// DictDataDeleteReq 删除字典数据请求。
type DictDataDeleteReq struct {
	g.Meta `path:"/system/dict/data/{id}" method:"delete" tags:"系统" summary:"删除字典数据"`
	Id     int64 `json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"字典数据 ID"`
}

// DictDataDeleteRes 删除字典数据响应。
type DictDataDeleteRes struct {
	Deleted bool `json:"deleted" dc:"是否已删除"`
}
