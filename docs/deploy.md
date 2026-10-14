# RoamAdmin 部署与运维教程

本文是面向**开发、测试、部署、运维**的操作手册：从环境准备、依赖安装、数据库初始化，
到本地开发、生产打包、部署上线，以及日常运维与常见问题排查。按顺序照做即可跑通。

配套文档：[README.md](../README.md)（快速开始）、[ARCHITECTURE.md](../ARCHITECTURE.md)（架构设计）。

---

## 目录

1. [环境准备](#1-环境准备)
2. [依赖安装](#2-依赖安装)
3. [配置说明](#3-配置说明)
4. [数据库初始化](#4-数据库初始化)
5. [本地开发启动](#5-本地开发启动)
6. [前端构建](#6-前端构建)
7. [生产打包](#7-生产打包)
8. [部署上线](#8-部署上线)
9. [运行验证](#9-运行验证)
10. [日常运维](#10-日常运维)
11. [模块化打包（按需裁剪模块）](#11-模块化打包按需裁剪模块)
12. [常见问题（FAQ）](#12-常见问题faq)

---

## 1. 环境准备

| 依赖 | 版本要求 | 说明 |
|---|---|---|
| Go | ≥ 1.27 | 后端语言与工具链 |
| Node.js + npm | ≥ 20（建议 22+） | 前端构建 |
| MySQL | 8.x | 主数据库（也支持 PostgreSQL，见 [配置](#3-配置说明)） |
| Redis | 6+ | 缓存 / 分布式锁 / 限流 / 幂等（**写接口必需**） |
| Git | 任意 | 拉代码、打版本 |

> 本机（Windows 11）验证过的组合：Go 1.27.0、Node v25.2.1、npm 11.6.2、MySQL 8、Redis 8.8.0（Windows fork）。

---

## 2. 依赖安装

```bash
# 2.1 拉取代码
git clone <仓库地址> && cd roamadmin

# 2.2 国内加速 Go 模块代理（可选但推荐）
go env -w GOPROXY=https://goproxy.cn,direct

# 2.3 下载后端依赖
go mod download

# 2.4 前端依赖
cd web && npm install && cd ..
```

> `go.mod` 中声明了可选依赖（如 `mochi-mqtt`），但**只有开启对应 build tag 时才会编译进二进制**，
> 见 [模块化打包](#11-模块化打包按需裁剪模块)。

---

## 3. 配置说明

配置文件：[manifest/config/config.yaml](../manifest/config/config.yaml)。
多环境可叠加 `config.${RA_ENV}.yaml`（如 `config.prod.yaml`），环境变量用 `RA_` 前缀注入。

```yaml
server:
  address: ":8600"          # HTTP 监听地址

database:
  default:
    type: "mysql"           # mysql | pgsql，一套代码切换
    host: "127.0.0.1"
    port: "3306"
    name: "roamadmin"       # 库名（需先创建）
    user: "root"
    pass: "root"
    timezone: "UTC"         # 固定 UTC：库统一存 UTC，接口返回 RFC3339(UTC)

redis:
  default:
    address: "127.0.0.1:6379"

auth:
  jwt:
    secret: "please-change-me-in-production"  # 生产必须更换
    issuer: "roamadmin"
    expireSeconds: 7200

mqtt:
  enabled: true             # 运行时是否启动 broker（仅当用 -tags mqtt 编译时生效）
  address: ":1883"          # MQTT 监听地址
```

要点：

- **时区固定 UTC**：不要改成本地时区。数据层统一 UTC，接口返回 `2026-10-13T12:01:15Z` 这类 RFC3339(UTC)，由前端转本地时区显示。
- **库名需手动创建**：`CREATE DATABASE roamadmin DEFAULT CHARACTER SET utf8mb4;`（框架只建表、不建库）。
- **`mqtt.enabled` 是运行时开关**，与编译开关（build tag）是两级控制，见 [模块化打包](#11-模块化打包按需裁剪模块)。

---

## 4. 数据库初始化

建表脚本在 [manifest/sql/](../manifest/sql/)：

| 脚本 | 内容 | 是否必须手动执行 |
|---|---|---|
| `init.sql` | demo 模块示例表 | 可选（demo 模块自建） |
| `system.sql` | 系统管理：用户/角色/权限/菜单/字典/租户/操作日志 + 初始账号 | **必须** |
| `iot.sql` | IoT：物模型/设备/影子 | 可选（device 模块启动时自建） |

```bash
# 方式 A：用 mysql 客户端执行
mysql -uroot -proot roamadmin < manifest/sql/system.sql

# 方式 B：可视化工具（Navicat / DBeaver）打开 system.sql 直接运行
```

> 说明：
> - `system` 模块依赖 `system.sql` 预先建表（含内置超级管理员 `admin/admin123`）。
> - `device` 模块启动时会幂等建表（`ra_device_model / ra_device_model_property / ra_device / ra_device_shadow`），
>   `iot.sql` 只是部署基线，可执行可不执行；时序表 `ra_ts_point` 由 tsdb 引擎首次写入时自建。

---

## 5. 本地开发启动

```bash
# 5.1 先确保 MySQL、Redis 已启动
# MySQL：服务已装则直接可用；Redis 需手动起（Windows fork 不会开机自启）
redis-cli ping        # 期望 PONG；不通则先启动 redis-server

# 5.2 启动后端（不含 MQTT 的默认构建）
go run ./cmd/server

# 5.3 需要 MQTT 时加 -tags mqtt（内嵌 broker 监听 :1883）
go run -tags mqtt ./cmd/server
```

启动成功日志关键行：

```
[ra] RoamAdmin v0.1.0 started, 4 module(s) registered
http server started listening on [:8600]
openapi specification is serving at address: http://127.0.0.1:8600/api.json
[ra][mqtt] broker listening on :1883          # 仅 -tags mqtt 且 mqtt.enabled=true 时出现
```

---

## 6. 前端构建

前端在 [web/](../web/)，构建产物内嵌进后端二进制（`go:embed`），单二进制、零 CORS。

```bash
cd web
npm install
npm run build     # 先自动从后端 /api.json 重新生成接口类型，再打包
```

> 构建前请先启动后端，否则 `npm run build` 拉取 `/api.json` 会失败。
> 产物输出到 `kernel/webui/dist`，下次 `go build` 时自动内嵌。

---

## 7. 生产打包

### 7.1 本机打包

```bash
# 不含 MQTT
go build -o server ./cmd/server

# 含 MQTT
go build -tags mqtt -o server ./cmd/server
```

### 7.2 交叉编译（在 Windows 上打 Linux 包）

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags mqtt -o server-linux ./cmd/server
```

> 项目为纯 Go、无 CGO，可放心交叉编译；`CGO_ENABLED=0` 产出静态链接二进制。

### 7.3 减小体积 / 去掉调试信息（可选）

```bash
go build -tags mqtt -ldflags "-s -w" -o server ./cmd/server
```

### 7.4 注意事项

- **Windows 开发机杀软（360 等）会隔离本地编译的 `.exe`**：本地调试用 `go run ./cmd/server` 而不是 `go build` 后执行；生产部署在 Linux 上不受影响。
- **JWT 密钥**：生产必须改 `auth.jwt.secret`，建议通过环境变量 `RA_AUTH_JWT_SECRET=xxx` 注入而非写死。

---

## 8. 部署上线

### 8.1 部署物清单

1. 编译好的二进制 `server`（或 `server-linux`）
2. 配置文件 `manifest/config/config.yaml`（生产环境按需改库/密钥）
3. 日志目录（首次启动自动创建）

### 8.2 systemd 托管（Linux）

```ini
# /etc/systemd/system/roamadmin.service
[Unit]
Description=RoamAdmin
After=network.target mysql.service redis.service

[Service]
WorkingDirectory=/opt/roamadmin
ExecStart=/opt/roamadmin/server
Restart=always
RestartSec=3
Environment=RA_ENV=prod
# 生产建议用环境变量注入敏感配置
Environment=RA_AUTH_JWT_SECRET=change-me
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

```bash
systemctl daemon-reload
systemctl enable --now roamadmin
systemctl status roamadmin
```

### 8.3 Docker 示例（可选）

```dockerfile
FROM golang:1.27 AS build
WORKDIR /src
COPY . .
RUN go build -tags mqtt -o /out/server ./cmd/server

FROM debian:bookworm-slim
WORKDIR /opt/roamadmin
COPY --from=build /out/server .
COPY manifest/config/config.yaml manifest/config/
EXPOSE 8600 1883
ENTRYPOINT ["./server"]
```

> 容器内数据库/Redis 建议走外部网络地址，避免容器重启丢数据。

### 8.4 反向代理（Nginx 示例）

后端已内嵌前端与 API 文档，通常直连即可；如需域名/HTTPS：

```nginx
server {
    listen 443 ssl;
    server_name admin.example.com;
    # ssl_certificate ...
    location / {
        proxy_pass http://127.0.0.1:8600;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

> MQTT 走 TCP，不能用 HTTP 反向代理，直接放行 1883 端口即可。

---

## 9. 运行验证

| 端点 | 说明 |
|---|---|
| `http://127.0.0.1:8600/admin` | 管理后台前端（默认账号 `admin/admin123`） |
| `http://127.0.0.1:8600/api.json` | OpenAPI 规范 |
| `http://127.0.0.1:8600/swagger` | API 文档界面（在线调试） |
| `http://127.0.0.1:8600/metrics` | Prometheus 指标 |

```bash
# 登录拿 token（后续带 Authorization: Bearer <token>）
curl -s -X POST http://127.0.0.1:8600/api/v1/system/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}'
```

统一响应体：

```json
{ "code": 0, "message": "ok", "data": { ... }, "traceId": "..." }
```

`code == 0` 表示成功；非 0 为错误码（`10002` 未认证、`10003` 无权限、`10010` 数据库错误、`20004` 数据重复等）。

---

## 10. 日常运维

### 10.1 启动 / 停止

- **优雅关闭**：进程收到 `SIGINT`（Ctrl+C）或 `SIGTERM` 时会按模块逆序执行 `Stop`（如关闭 MQTT broker）、释放连接。
- 停止后 `:8600`（HTTP）与 `:1883`（MQTT）端口应立即释放。

### 10.2 健康检查

```bash
# HTTP 存活
curl -sf http://127.0.0.1:8600/api/v1/demo/ping || echo DOWN

# 端口监听
netstat -ano | grep -E ":(8600|1883)"

# 依赖状态（写接口依赖，必须 PONG）
redis-cli ping
mysqladmin -uroot -proot status
```

### 10.3 日志

- 后端日志输出到 `logs/`，按天滚动：`logs/{Y-m-d}.log`。
- 每次请求带 `traceId`，可在日志中按 `traceId` 串联排查。
- MQTT broker（mochi）日志与业务日志混在 stdout，关注 `level=WARN/ERROR` 行。

### 10.4 端口与资源

- 默认端口：HTTP `:8600`、MQTT `:1883`，冲突时改 [config.yaml](../manifest/config/config.yaml)。
- 高并发连接池参数在 `database.default.max*` 与 `redis.default.max*`，调优参考 [ARCHITECTURE.md §8.5](../ARCHITECTURE.md)。

---

## 11. 模块化打包（按需裁剪模块）

模块化单体：业务模块实现统一 `Module` SPI，在 `cmd/server/modules.go` 里**空导入即注册**。

### 11.1 增删业务模块（编译期选择）

增删模块 = 增删一行 import：

```go
// cmd/server/modules.go
import (
    _ "github.com/roamdream/roamadmin/modules/demo"
    _ "github.com/roamdream/roamadmin/modules/device"
    _ "github.com/roamdream/roamadmin/modules/gateway"
    _ "github.com/roamdream/roamadmin/modules/system"
)
```

- **去掉某模块**：注释/删除对应 import，该模块的 HTTP 路由、逻辑即不注册。
- **新增模块**：按 [README「增加一个模块」](../README.md#增加一个模块) 实现 SPI + 加 import。

### 11.2 可选重依赖：MQTT（build tag）

MQTT 是**重依赖**（mochi-mqtt），已用 Go build tag 做成**可选**：

| 命令 | 是否编译 MQTT | 二进制是否含 mochi |
|---|---|---|
| `go build ./cmd/server` | ❌ 否 | 否（0 个 mochi 包） |
| `go build -tags mqtt ./cmd/server` | ✅ 是 | 是 |

机制：`kernel/mqtt` 拆成两个文件——

- `broker.go`（`//go:build mqtt`）：真正的 mochi broker；
- `broker_stub.go`（`//go:build !mqtt`）：同签名占位实现，无任何 MQTT 依赖。

`ra.MQTT()` 门面在两种构建下都能编译：不加 tag 时返回占位实现，`Publish` 等返回错误，
`/gateway/status` 显示 `running:false`；加 `-tags mqtt` 后才是可用 broker。

> 说明：`go.mod` 仍保留 mochi 声明（`go mod tidy` 会扫描所有 build tag），
> 但**编译进二进制的部分**由 `-tags mqtt` 决定。若要让 mochi 彻底从 go.mod 消失，
> 需把它拆到独立 Go module（插件化），属于进一步演进，一般没必要。

### 11.3 完整裁剪示例：一个"不用 MQTT 也不接设备"的后台

```bash
# 1) 不注册 device / gateway 模块（编辑 cmd/server/modules.go 去掉这两行 import）
# 2) 不用 MQTT，直接默认构建即可（无需 -tags mqtt）
go build -o server ./cmd/server
```

这样产出的二进制只含 `system`（+可选 `demo`），MQTT、设备、时序能力全部不进入二进制。

---

## 12. 常见问题（FAQ）

| 现象 | 原因 | 处理 |
|---|---|---|
| 写接口返回 `{"code":10000,"message":"internal server error"}`，日志有 `Redis ... refused` | Redis 未启动（写路径走 `ra.Lock()` 分布式锁） | 启动 redis-server，`redis-cli ping` 确认 PONG |
| 偶发 `{"code":10010,"message":"database error"}` / `invalid connection` | Windows 杀软/OS 重置了空闲 MySQL 连接 | 重试即可；属瞬态，非代码缺陷 |
| 本地 `go build` 出的 `.exe` 被删/报毒 | 360 等杀软隔离 | 本地用 `go run ./cmd/server`；生产部署 Linux |
| MQTT 客户端连不上，日志 `bad username or password` | 老版本默认拒绝匿名连接 | 已内置 `AllowHook` 开放接入；若仍报错确认用了最新代码 |
| 登录成功但刷新后菜单消失 | 刷新后 `userInfo` 为空 → 权限过滤掉了菜单 | 前端 `router.beforeEach` 已自动 `fetchInfo`，确认前端已重新构建 |
| 接口时间比本地早/晚 8 小时 | 库按 UTC 存，返回 RFC3339(UTC) | 属预期设计，前端负责转本地时区 |
| 改代码不生效 | `go run` 缓存 / 前端未重构建 | 后端重启；前端改完 `npm run build` 后重新 `go build` |

---

### 附：一键起停（开发常用）

```bash
# 依赖（Redis）——Windows fork 需手动起一次
"/c/Users/<you>/AppData/Local/Microsoft/WinGet/Packages/...redis.../redis-server" &

# 后端（含 MQTT）
go run -tags mqtt ./cmd/server
# Ctrl+C 优雅停止
```
