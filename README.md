# Twitta

## 项目介绍

一个仿 Twitter 的后端项目，基于 Gin + MongoDB + Redis 构建，集成 Elasticsearch 全文搜索、WebSocket 实时聊天、人脸识别登录、OSS 文件存储等能力。

## 功能列表

- [x] 用户注册 / 登录（账号密码 + 人脸识别）
- [x] JWT 单点登录（支持多平台 token 隔离）
- [x] RBAC 权限验证（Casbin）
- [x] 用户资料管理
- [x] 推文发布 / 删除（支持图片/视频附件）
- [x] 推文列表（分页 + 按时间倒序）
- [x] 全文搜索（用户 / 推文，基于 Elasticsearch）
- [x] 关注 / 粉丝
- [x] 好友申请 / 通过 / 拒绝 / 删除
- [x] **好友列表**
- [x] **WebSocket 实时聊天**（历史消息持久化到 MongoDB）
- [x] 点赞 / 取消点赞
- [x] 评论（支持多级嵌套）
- [x] 收藏推文
- [x] gRPC 微服务（Elasticsearch / OSS / 推送 / 人脸识别）
- [x] Prometheus 监控 + pprof 性能分析
- [x] 注册 / 登录邮件通知

## 快速启动

### Docker（推荐）

```bash
docker-compose up -d
```

### 本地运行

```bash
make run
```

## WebSocket 聊天

连接地址：`ws://localhost:6565/api/twitta/chats/{对方用户ID}/ws`

请求头需携带 `token`。

**发送消息格式（JSON）：**
```json
{ "content": "你好" }
```

**接收消息格式（JSON）：**
```json
{
  "roomId": "uid1:uid2",
  "fromId": "发送方用户ID",
  "content": "你好",
  "createdAt": "2024-01-01T12:00:00Z"
}
```

历史消息通过 HTTP 接口查询：`GET /api/twitta/chats/{对方用户ID}?page=1&size=20`

## Swagger

### 生成文档

```bash
make swagger
```

### 访问地址

```
http://localhost:6565/api/twitta/swagger/index.html
```

## 技术栈

| 组件 | 说明 |
|------|------|
| Gin | HTTP 框架 |
| MongoDB | 主数据库 |
| Redis | Token 存储 / 缓存 |
| Elasticsearch | 全文搜索（通过 gRPC 微服务） |
| gorilla/websocket | 实时聊天 |
| Consul | 服务注册与发现 |
| Casbin | RBAC 权限控制 |
| JWT | 身份认证 |
| Prometheus | 指标监控 |
| OSS（MinIO） | 文件存储（通过 gRPC 微服务） |
