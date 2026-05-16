# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此仓库中工作时提供指引。

## 常用命令

```bash
# 启动服务（通过脚本构建并运行）
make run

# 生成 Swagger 文档（需要安装 swag CLI）
make swagger

# 运行所有测试
go test ./...

# 运行单个测试
go test ./services/... -run TestServiceUserSearch

# 手动构建
go build -o ./build/twitta cmd/main.go

# Docker（含全部依赖服务）
docker-compose up -d
```

服务启动后，Swagger UI 地址：`http://localhost:6565/api/twitta/swagger/index.html`

## 架构

Twitta 是基于 Gin 构建的仿 Twitter REST API，采用分层架构：

```
cmd/main.go          → 入口：初始化全局对象、向 Consul 注册服务、启动 HTTP 服务器
global/              → 包级单例（DB、Redis、gRPC 客户端、配置）
global/initialize/   → 启动时调用的初始化函数（配置、日志、Mongo、外部 gRPC 服务）
configs/             → 配置结构体 + configs/config.yml（通过 Viper 加载）
routers/             → Gin 路由注册与 HTTP 服务器生命周期管理
services/            → 业务逻辑层（由路由处理函数调用）
pkg/dao/             → MongoDB 数据访问（泛型通用层 + 各模型专用层）
forms/               → 跨层使用的请求/响应结构体
pkg/middlewares/     → JWT 认证、RBAC（Casbin）、请求日志
pkg/response/        → 统一 JSON 响应封装
pkg/cache/           → Redis 客户端工厂
pkg/utils/           → 通用工具函数
```

### 核心设计模式

**泛型 DAO 层**（`pkg/dao/generic.go`）：基于 `TableNamer` 约束使用 Go 泛型。所有 MongoDB 操作均通过 `GInsertOne[T]`、`GWhereFirst[T]`、`GPaginatorOrder[T]` 等函数完成，无需为每个模型编写重复代码。模型结构体通过实现 `TableName() string` 来指定对应的集合名。

**全局单例**（`global/global.go`）：`global.DB`（MongoDB）、`global.RedisMapPool`（按 DB 索引区分的 Redis）以及各 gRPC 客户端存根（`EsSrvClient`、`PushSrvClient`、`FaceRecognitionSrvClient`、`OssSrvClient`）均为包级变量，在启动时初始化。各层直接导入 `twitta/global` 使用。

**通过 gRPC + Consul 接入外部服务**：Elasticsearch、OSS、消息推送、人脸识别均为独立的 gRPC 微服务，通过 Consul 进行服务发现。其 proto 生成的客户端位于 `github.com/solost23/protopb`。

**认证流程**：JWT token 通过请求头 `token` 字段传递（非 `Authorization: Bearer`）。中间件验证 token 后，会与 Redis 中的记录（device+userId 为 key）进行比对，以实现单设备登录。验证通过的 `*dao.User` 以 `"user"` 为键存入 Gin context。

**RBAC**：Casbin 中间件（`configs/rbac_model.conf` + `configs/rbac_policy.csv`）负责角色权限控制。所有需认证的路由要求角色为 `admin` 或 `user`（二选一）。

**路由结构**：公开路由（注册、登录、搜索、推文列表）挂载在 `api/twitta` 下。需认证的路由应用 JWT + Casbin 中间件，并按资源分组（users、tweets、friends、chats、fans、comments）。

**服务层**（`services/`）：单一 `Service` 结构体，按业务域划分方法。`services/*_test.go` 中的测试直接调用服务方法，使用真实的 Gin 测试 context，依赖真实的基础设施（MongoDB、Redis、gRPC 服务）。

### 响应格式

所有接口均返回 HTTP 200，响应体为统一 JSON 结构：
```json
{"code": 0, "message": "", "success": true, "data": ...}
```
在 handler 中统一使用 `response.Success(c, data)` 或 `response.Error(c, code, err)`，不要直接调用 `c.JSON`。
