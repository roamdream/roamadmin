# RoamAdmin（ra）

面向大型企业级物联网、高并发项目的 Go 开发框架。参考 [GoFrame](https://goframe.org) 的模块化思路，
把"通用能力"沉淀进内核，业务/设备能力以"增加模块"的方式不断扩展。

- 模块路径：`github.com/roamdream/roamadmin`
- 底层底座：GoFrame v2（ghttp / gdb ORM / gredis / gcfg / OpenAPI）
- 品牌前缀：facade `ra`、CLI `ra`、环境变量 `RA_`、日志前缀 `[ra]`

## 特性

- **模块化单体**：业务模块实现统一 `Module` SPI（`Name/Version/Init/Start/Stop`），编译期注册，增删模块 = 增删目录 + 改一行 import。
- **高并发提前内置**：协程池、分布式锁、幂等、限流/熔断、事件总线（含背压）、Prometheus metrics + traceId、统一错误码与响应体。
- **数据库多方言兼容**：全程 ORM、雪花主键、MySQL 8 / PostgreSQL 一套代码通吃。
- **多租户 + RBAC**：JWT 认证，用户/角色/权限/菜单/字典/租户/操作日志开箱即用。
- **全球化就绪**：后端接口一律英文输出，时间统一存 UTC、返回 RFC3339(UTC)，多语言与本地时区由前端实现。
- **代码生成**：`gen` CLI 一键生成标准 CRUD 模块（生成的模块自带英文 API 文档与校验）。

## 快速开始

```bash
# 1. 准备依赖（国内可先切代理）
go env -w GOPROXY=https://goproxy.cn,direct

# 2. 初始化数据库（MySQL 8 / PostgreSQL）
#    编辑 manifest/config/config.yaml 中的 database 配置
#    执行 manifest/sql/ 下的建表脚本

# 3. 编译并运行（默认 :8600）
go run ./cmd/server

# 3.1 需要 MQTT / 设备接入时，加 build tag（内嵌 broker 监听 :1883）
go run -tags mqtt ./cmd/server
```

> 完整的环境准备、数据库初始化、生产打包、部署上线、日常运维与常见问题排查，见
> **[部署与运维教程](docs/deploy.md)**；按需裁剪模块（含 MQTT 开关）见 [按需裁剪模块](#按需裁剪模块)。

启动后：

| 端点 | 说明 |
|---|---|
| `http://127.0.0.1:8600/admin` | 管理后台前端（Vue3 + Element Plus，内嵌、与后端同源） |
| `http://127.0.0.1:8600/api/v1/demo/ping` | 连通性测试 |
| `http://127.0.0.1:8600/metrics` | Prometheus 指标 |
| `http://127.0.0.1:8600/api.json` | OpenAPI 规范（机器可读） |
| `http://127.0.0.1:8600/swagger` | API 文档界面（Scalar，在线调试，内嵌、离线可用） |

统一响应体：

```json
{ "code": 0, "message": "ok", "data": { ... }, "traceId": "..." }
```

## 目录结构

```
roamadmin/
├── cmd/server/          # 入口：装配中间件 → 启动应用
├── kernel/              # 内核（通用能力，不依赖任何业务模块）
│   ├── module/ errcode/ config/ logger/ db/ cache/ auth/ event/ mq/ tsdb/ mqtt/
│   ├── lock/ idempotent/ pool/ limit/ observ/ utils/ swaggerui/
│   └── ra/              # 门面：业务模块面向 ra 编程
├── modules/             # 业务/设备模块（可插拔）
│   ├── system/          # 系统管理：登录/JWT、用户/角色/权限、菜单、字典、多租户、操作日志
│   ├── demo/            # 示例模块（模块模板）
│   ├── device/          # IoT 设备：物模型/设备/影子/遥测（依赖 tsdb）
│   └── gateway/         # MQTT 接入网关：内嵌 Broker 上行接入 → 事件总线分发
├── api/                 # 共享请求/响应结构
├── hack/gen/            # 代码生成 CLI
└── manifest/            # 配置 + 建表脚本
```

## 增加一个模块

1. 新建 `modules/<name>/`，实现 `Module` 接口，`init()` 里 `module.Register(New())`；
2. 在 `cmd/server/modules.go` 加一行 `import _ "github.com/roamdream/roamadmin/modules/<name>"`；
3. 或直接用代码生成器：`go run ./hack/gen -m <name> -t <table> -title Title`（`-title` 为英文标题，缺省用模块名）。

> 接口规范：所有接口输出与 API 文档（`summary`/`dc`/校验提示/错误信息）统一英文；
> 时间字段统一 RFC3339(UTC)。新增模块写好 `g.Meta` 后，API 文档自动生成，无需额外维护。

## 按需裁剪模块

模块化单体：增删业务模块 = 增删 `cmd/server/modules.go` 里的一行 import，未 import 的模块不注册、不进二进制。

```go
import (
    _ "github.com/roamdream/roamadmin/modules/demo"
    _ "github.com/roamdream/roamadmin/modules/device"   // 删掉即不带 IoT 设备能力
    _ "github.com/roamdream/roamadmin/modules/gateway"  // 删掉即不带 MQTT 接入
    _ "github.com/roamdream/roamadmin/modules/system"
)
```

**重依赖（MQTT）用 Go build tag 做成可选**，默认构建不引入 `mochi-mqtt`：

| 命令 | MQTT | 二进制是否含 mochi |
|---|---|---|
| `go build ./cmd/server` | ❌ 关 | 否 |
| `go build -tags mqtt ./cmd/server` | ✅ 开 | 是 |

机制：`kernel/mqtt` 拆成 `broker.go`（`//go:build mqtt`，真 broker）与 `broker_stub.go`
（`//go:build !mqtt`，同签名占位）。`ra.MQTT()` 门面两种构建都能编译；未加 tag 时 `Publish` 返回错误、
`/gateway/status` 显示 `running:false`。运行时再用 `mqtt.enabled` 二级开关控制是否真正监听 `:1883`。

> 一个"不用 MQTT 也不接设备"的后台：删掉 device/gateway 两行 import，直接 `go build ./cmd/server`，
> 产出仅含 `system`（+可选 `demo`）。详见 [部署与运维教程](docs/deploy.md)。

## 管理后台前端

前端位于 [web/](web/)，技术栈 Vue3 + TypeScript + Vite + Element Plus + Pinia。构建产物输出到
`kernel/webui/dist`，由内核 `go:embed` 内嵌，运行时挂在 `/admin`，单二进制、零 CORS。

```bash
cd web
npm install
npm run build      # 自动从后端 /api.json 重新生成接口类型，再打包
```

- **接口自动对接**：`npm run build`（或 `npm run gen:api`）会先从后端拉取 `/api.json`（OpenAPI），
  由 `scripts/gen-api.mjs` 生成 `src/api/types.ts` 类型；新增/修改后端接口后重新构建即可，前端类型自动同步，
  无需手写对接代码。
- **默认账号**：`admin` / `admin123`（内置超级管理员）。

## 路线图

| 阶段 | 内容 | 状态 |
|---|---|---|
| P0 | 框架底座：monorepo / kernel / Module SPI / 配置 / 日志 / DB / 缓存 / 错误码 / 限流熔断 / 协程池 / 可观测 / 优雅关闭 / OpenAPI | ✅ |
| P1 | `system` 模块：登录/JWT、RBAC、菜单、字典、多租户、操作日志、幂等、分布式锁 | 🚧 |
| P2 | `demo` 模块 + `gen` 代码生成 CLI | 🚧 |
| P3 | IoT 接入：MQTT / 设备接入 / 物模型 / `tsdb` 落地 | ✅ |
| P4 | 高并发增强：分布式 / 消息总线落库 / 连接池调优 / 拆微服务 | ⏳ |

设计文档见 [ARCHITECTURE.md](ARCHITECTURE.md)。
