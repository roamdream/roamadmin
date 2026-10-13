// gen 是 RoamAdmin 的一键 CRUD 模块生成器（P2 阶段）。
//
// 依据给定的表名，在 modules/<moduleName>/ 下生成一个完整可编译、可直接注册的
// 业务模块（含标准 CRUD），结构与 modules/demo 完全一致：
//
//	<moduleName>.go                    —— 实现 Module 接口 + init() 里 Register
//	api/<moduleName>.go                —— 请求/响应 DTO（带 g.Meta 与 v 校验）
//	internal/model/model.go            —— 实体（orm tag + TableName()）
//	internal/service/<moduleName>.go   —— Create/List/Get/Update/Delete
//	internal/controller/<moduleName>.go —— handler (ctx, req) (res, error)
//
// 用法：
//
//	gen -m <moduleName> -t <TableName> [-title <中文标题>] [-force]
//
// 示例：
//
//	gen -m device -t device_info -title 设备
//
// 生成完成后，在 cmd/server/modules.go 里加一行 import 即可注册该模块：
//
//	_ "github.com/roamdream/roamadmin/modules/device"
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

// tmplFuncs 提供在模板里输出反引号的能力（Go 结构体 tag 需要反引号，
// 而模板本体是反引号包裹的原始字符串，无法直接书写反引号）。
var tmplFuncs = template.FuncMap{
	"bt": func() string { return "`" },
}

// data 模板渲染所需的全部命名结果。
type data struct {
	Module   string // 模块名（目录/包名），如 device
	Table    string // 表名，如 device_info
	Title    string // 中文标题，如 设备
	Entity   string // 实体结构体名，如 DeviceInfo
	Resource string // 路由资源段（kebab-case），如 device-info
}

var identRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func main() {
	var (
		moduleName string
		tableName  string
		title      string
		force      bool
	)
	flag.StringVar(&moduleName, "m", "", "模块名（生成目录 modules/<moduleName>，如 device）")
	flag.StringVar(&tableName, "t", "", "表名（如 device_info）")
	flag.StringVar(&title, "title", "", "中文标题（用于 API 文档 summary/dc，如 设备；缺省用模块名）")
	flag.BoolVar(&force, "force", false, "目标目录已存在时强制覆盖")
	flag.Parse()

	if moduleName == "" || tableName == "" {
		usage()
		os.Exit(2)
	}
	if !validIdent(moduleName) {
		fmt.Fprintf(os.Stderr, "错误: 模块名 %q 不是合法的 Go 标识符\n", moduleName)
		os.Exit(2)
	}
	if !validIdent(tableName) {
		fmt.Fprintf(os.Stderr, "错误: 表名 %q 不是合法的标识符\n", tableName)
		os.Exit(2)
	}
	if title == "" {
		title = moduleName
	}

	targetDir := filepath.Join("modules", moduleName)
	if gfile.IsDir(targetDir) {
		if !force {
			fmt.Fprintf(os.Stderr, "错误: 目标目录 %s 已存在，如需覆盖请加 -force\n", targetDir)
			os.Exit(1)
		}
	}

	d := &data{
		Module:   moduleName,
		Table:    tableName,
		Title:    title,
		Entity:   gstr.CaseCamel(tableName),
		Resource: gstr.CaseKebab(gstr.CaseCamel(tableName)),
	}

	files := []struct {
		name string
		tmpl string
	}{
		{filepath.Join(targetDir, moduleName+".go"), tmplModule},
		{filepath.Join(targetDir, "api", moduleName+".go"), tmplAPI},
		{filepath.Join(targetDir, "internal", "model", "model.go"), tmplModel},
		{filepath.Join(targetDir, "internal", "service", moduleName+".go"), tmplService},
		{filepath.Join(targetDir, "internal", "controller", moduleName+".go"), tmplController},
	}

	for _, f := range files {
		if err := render(d, f.name, f.tmpl); err != nil {
			fmt.Fprintf(os.Stderr, "生成 %s 失败: %v\n", f.name, err)
			os.Exit(1)
		}
		fmt.Println("已生成:", f.name)
	}

	fmt.Printf("\n完成。请在 cmd/server/modules.go 里加一行 import 完成注册：\n")
	fmt.Printf("\t_ \"github.com/roamdream/roamadmin/modules/%s\"\n", d.Module)
}

func usage() {
	fmt.Fprintln(os.Stderr, "用法: gen -m <moduleName> -t <TableName> [-title <中文标题>] [-force]")
	fmt.Fprintln(os.Stderr, "示例: gen -m device -t device_info -title 设备")
}

func validIdent(s string) bool {
	return identRe.MatchString(s)
}

// render 渲染模板并写入目标文件（目录不存在时递归创建）。
func render(d *data, path, tmplText string) error {
	t, err := template.New(path).Funcs(tmplFuncs).Parse(tmplText)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, d); err != nil {
		return err
	}
	return gfile.PutContents(path, buf.String())
}

// tmplModule 模块入口：实现 Module 接口，init() 里 Register，Start() 里绑定路由。
const tmplModule = `// Package {{.Module}} {{.Title}}模块：由 gen 一键生成的 CRUD 业务模块。
//
// 标准动作 = 新建目录 + 实现 Module 接口 + init() 里 Register + 在
// cmd/server/modules.go 里加一行 import。内核自动接管 Init/Start/Stop。
package {{.Module}}

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/roamdream/roamadmin/kernel/module"
	"github.com/roamdream/roamadmin/modules/{{.Module}}/internal/controller"
)

func init() {
	module.Register(New())
}

// Module {{.Module}} 模块。
type Module struct{}

// New 创建 {{.Module}} 模块实例。
func New() *Module { return &Module{} }

// Name 模块名。
func (m *Module) Name() string { return "{{.Module}}" }

// Version 版本号。
func (m *Module) Version() string { return "0.1.0" }

// Init 依赖注入阶段：{{.Module}} 无内部依赖，直接返回。
func (m *Module) Init(ctx context.Context) error { return nil }

// Start 启动阶段：注册 HTTP 路由。
func (m *Module) Start(ctx context.Context) error {
	m.registerRoutes(ctx)
	return nil
}

// Stop 优雅关闭阶段。
func (m *Module) Stop(ctx context.Context) error { return nil }

func (m *Module) registerRoutes(ctx context.Context) {
	s := g.Server()
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		// 控制器方法的路由由 api 中请求结构体的 g.Meta 定义。
		group.Bind(controller.New())
	})
}
`

// tmplAPI 对外契约：请求/响应 DTO，带 g.Meta 与 v 校验。
const tmplAPI = `// Package api {{.Module}} 模块对外契约（请求/响应 DTO）。
package api

import (
	"github.com/gogf/gf/v2/frame/g"
)

// {{.Entity}} {{.Title}}数据。
type {{.Entity}} struct {
	Id        int64  {{bt}}json:"id" dc:"ID"{{bt}}
	Name      string {{bt}}json:"name" dc:"名称"{{bt}}
	Remark    string {{bt}}json:"remark" dc:"备注"{{bt}}
	CreatedAt string {{bt}}json:"createdAt" dc:"创建时间"{{bt}}
	UpdatedAt string {{bt}}json:"updatedAt" dc:"更新时间"{{bt}}
}

// CreateReq 创建{{.Title}}请求。
type CreateReq struct {
	g.Meta {{bt}}path:"/{{.Module}}/{{.Resource}}" method:"post" tags:"{{.Module}}" summary:"创建{{.Title}}"{{bt}}
	Name   string {{bt}}json:"name" v:"required|length:1,64#名称不能为空|名称长度1-64" dc:"名称"{{bt}}
	Remark string {{bt}}json:"remark" dc:"备注"{{bt}}
}

// CreateRes 创建{{.Title}}响应。
type CreateRes struct {
	Id int64 {{bt}}json:"id" dc:"新数据 ID"{{bt}}
}

// ListReq 分页查询{{.Title}}请求。
type ListReq struct {
	g.Meta {{bt}}path:"/{{.Module}}/{{.Resource}}" method:"get" tags:"{{.Module}}" summary:"分页查询{{.Title}}"{{bt}}
	Page   int    {{bt}}json:"page" d:"1" v:"min:1#页码最小为1" dc:"页码"{{bt}}
	Size   int    {{bt}}json:"size" d:"20" v:"min:1|max:200#每页最小为1|每页最大为200" dc:"每页条数"{{bt}}
	Name   string {{bt}}json:"name" dc:"名称模糊搜索"{{bt}}
}

// ListRes 分页查询{{.Title}}响应。
type ListRes struct {
	List  []*{{.Entity}} {{bt}}json:"list" dc:"数据列表"{{bt}}
	Total int64         {{bt}}json:"total" dc:"总数"{{bt}}
	Page  int           {{bt}}json:"page" dc:"页码"{{bt}}
	Size  int           {{bt}}json:"size" dc:"每页条数"{{bt}}
}

// GetReq 查询单条{{.Title}}请求。
type GetReq struct {
	g.Meta {{bt}}path:"/{{.Module}}/{{.Resource}}/{id}" method:"get" tags:"{{.Module}}" summary:"查询单条{{.Title}}"{{bt}}
	Id     int64 {{bt}}json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"ID"{{bt}}
}

// UpdateReq 更新{{.Title}}请求。
type UpdateReq struct {
	g.Meta {{bt}}path:"/{{.Module}}/{{.Resource}}/{id}" method:"put" tags:"{{.Module}}" summary:"更新{{.Title}}"{{bt}}
	Id     int64  {{bt}}json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"ID"{{bt}}
	Name   string {{bt}}json:"name" v:"required|length:1,64#名称不能为空|名称长度1-64" dc:"名称"{{bt}}
	Remark string {{bt}}json:"remark" dc:"备注"{{bt}}
}

// UpdateRes 更新{{.Title}}响应。
type UpdateRes struct {
	Updated bool {{bt}}json:"updated" dc:"是否已更新"{{bt}}
}

// DeleteReq 删除{{.Title}}请求。
type DeleteReq struct {
	g.Meta {{bt}}path:"/{{.Module}}/{{.Resource}}/{id}" method:"delete" tags:"{{.Module}}" summary:"删除{{.Title}}"{{bt}}
	Id     int64 {{bt}}json:"id" in:"path" v:"required|min:1#ID不能为空|ID不合法" dc:"ID"{{bt}}
}

// DeleteRes 删除{{.Title}}响应。
type DeleteRes struct {
	Deleted bool {{bt}}json:"deleted" dc:"是否已删除"{{bt}}
}
`

// tmplModel 数据模型：实体（orm tag + TableName()）。
const tmplModel = `// Package model {{.Module}} 模块数据模型。
package model

import "github.com/gogf/gf/v2/os/gtime"

// {{.Entity}} {{.Title}}实体（对应表 {{.Table}}）。
type {{.Entity}} struct {
	Id        int64       {{bt}}json:"id"        orm:"id,primary"{{bt}}
	Name      string      {{bt}}json:"name"      orm:"name"{{bt}}
	Remark    string      {{bt}}json:"remark"    orm:"remark"{{bt}}
	CreatedAt *gtime.Time {{bt}}json:"createdAt" orm:"created_at"{{bt}}
	UpdatedAt *gtime.Time {{bt}}json:"updatedAt" orm:"updated_at"{{bt}}
}

// TableName 返回数据库表名。
func ({{.Entity}}) TableName() string { return "{{.Table}}" }
`

// tmplService 业务逻辑：Create/List/Get/Update/Delete。
const tmplService = `// Package service {{.Module}} 模块业务逻辑。
package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/roamdream/roamadmin/kernel/ra"
	"github.com/roamdream/roamadmin/modules/{{.Module}}/internal/model"
)

const table = "{{.Table}}"

// Create 创建{{.Title}}（雪花主键）。
func Create(ctx context.Context, name, remark string) (int64, error) {
	id := ra.DB().NextID()
	now := gtime.Now()
	_, err := ra.DB().Model(table).Data(g.Map{
		"id":         id,
		"name":       name,
		"remark":     remark,
		"created_at": now,
		"updated_at": now,
	}).Insert()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// List 分页查询。
func List(ctx context.Context, page, size int, name string) ([]*model.{{.Entity}}, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 200 {
		size = 20
	}
	m := ra.DB().Model(table)
	if name != "" {
		m = m.WhereLike("name", "%"+name+"%")
	}
	total, err := m.Count()
	if err != nil || total == 0 {
		return []*model.{{.Entity}}{}, int64(total), err
	}
	var list []*model.{{.Entity}}
	if err = m.Page(page, size).OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, err
	}
	return list, int64(total), nil
}

// Get 查询单条。
func Get(ctx context.Context, id int64) (*model.{{.Entity}}, error) {
	var item model.{{.Entity}}
	if err := ra.DB().Model(table).Where("id", id).Scan(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

// Update 更新单条。
func Update(ctx context.Context, id int64, name, remark string) error {
	_, err := ra.DB().Model(table).Where("id", id).Data(g.Map{
		"name":       name,
		"remark":     remark,
		"updated_at": gtime.Now(),
	}).Update()
	return err
}

// Delete 删除单条。
func Delete(ctx context.Context, id int64) error {
	_, err := ra.DB().Model(table).Where("id", id).Delete()
	return err
}
`

// tmplController HTTP 控制器：handler (ctx, req) (res, error)。
const tmplController = `// Package controller {{.Module}} 模块 HTTP 控制器。
package controller

import (
	"context"

	"github.com/roamdream/roamadmin/modules/{{.Module}}/api"
	"github.com/roamdream/roamadmin/modules/{{.Module}}/internal/service"
)

// Controller {{.Module}} 控制器。
type Controller struct{}

// New 创建控制器。
func New() *Controller { return &Controller{} }

// Create 创建{{.Title}}。
func (c *Controller) Create(ctx context.Context, req *api.CreateReq) (*api.CreateRes, error) {
	id, err := service.Create(ctx, req.Name, req.Remark)
	if err != nil {
		return nil, err
	}
	return &api.CreateRes{Id: id}, nil
}

// List 分页查询{{.Title}}。
func (c *Controller) List(ctx context.Context, req *api.ListReq) (*api.ListRes, error) {
	items, total, err := service.List(ctx, req.Page, req.Size, req.Name)
	if err != nil {
		return nil, err
	}
	out := make([]*api.{{.Entity}}, 0, len(items))
	for _, it := range items {
		out = append(out, &api.{{.Entity}}{
			Id:        it.Id,
			Name:      it.Name,
			Remark:    it.Remark,
			CreatedAt: it.CreatedAt.String(),
			UpdatedAt: it.UpdatedAt.String(),
		})
	}
	return &api.ListRes{List: out, Total: total, Page: req.Page, Size: req.Size}, nil
}

// Get 查询单条{{.Title}}。
func (c *Controller) Get(ctx context.Context, req *api.GetReq) (*api.{{.Entity}}, error) {
	it, err := service.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &api.{{.Entity}}{
		Id:        it.Id,
		Name:      it.Name,
		Remark:    it.Remark,
		CreatedAt: it.CreatedAt.String(),
		UpdatedAt: it.UpdatedAt.String(),
	}, nil
}

// Update 更新{{.Title}}。
func (c *Controller) Update(ctx context.Context, req *api.UpdateReq) (*api.UpdateRes, error) {
	if err := service.Update(ctx, req.Id, req.Name, req.Remark); err != nil {
		return nil, err
	}
	return &api.UpdateRes{Updated: true}, nil
}

// Delete 删除{{.Title}}。
func (c *Controller) Delete(ctx context.Context, req *api.DeleteReq) (*api.DeleteRes, error) {
	if err := service.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &api.DeleteRes{Deleted: true}, nil
}
`
