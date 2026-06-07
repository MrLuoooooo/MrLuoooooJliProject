# Community Forum

论坛后端服务，Go + Gin + GORM + JWT + Uber Fx。

实时消息推送对接 JuggleIM。

## 快速开始

```bash
cd server
go run .
```

服务启动在 `:1807`。

## API 文档

[Swagger UI](https://editor.swagger.io/?url=https://raw.githubusercontent.com/... ) — 或者直接看 `server/docs/swagger.yaml`。

把 `swagger.yaml` 贴到 [Swagger Editor](https://editor.swagger.io) 即可预览。

## 部署

```bash
docker compose up -d
```

启动 MySQL + Redis + IM Server + Community Server 全套。

## 项目结构

```
server/
├── main.go              # 入口，Uber Fx DI
├── config.yaml           # 配置
├── docker-compose.yaml   # 一键部署
├── docs/
│   └── swagger.yaml      # API 文档
├── internal/
│   ├── config/           # 配置加载
│   ├── db/
│   │   ├── mysql/        # GORM 模型 + 自动迁移
│   │   └── redis/        # Redis 连接（可选）
│   ├── di/               # Fx 依赖注入
│   ├── handler/          # HTTP 处理器
│   ├── im/               # JuggleIM 客户端
│   ├── middleware/        # JWT / CORS / 限流
│   ├── model/            # 请求/响应 DTO
│   ├── repository/       # 数据访问接口 + GORM 实现
│   ├── router/           # 路由注册
│   └── service/          # 业务逻辑
└── pkg/
    ├── jwt/              # JWT 签发/解析
    └── response/         # 统一响应格式
```

## 技术栈

| 组件 | 选型 |
|------|------|
| Web 框架 | Gin |
| ORM | GORM v2 + MySQL |
| 依赖注入 | Uber Fx |
| 认证 | JWT (golang-jwt/v5) |
| 密码 | bcrypt |
| 日志 | Zap + Lumberjack |
| IM | JuggleIM REST API |
| 缓存 | Redis (go-redis/v9) |
| AI | SiliconFlow (Qwen3-8B) |
