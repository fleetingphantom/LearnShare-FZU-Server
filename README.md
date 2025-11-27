# LearnShare-FZU-Server

基于 CloudWeGo Hertz 框架的 Go 后端服务，为福大学生学习分享平台提供 API 服务。

## 一句话说明
高性能、模块化的学习资源分享后端（Hertz + GORM + Redis + OSS）。

## 目录

- [项目简介](#项目简介)
- [主要特性](#主要特性)
- [技术栈](#技术栈)
- [环境要求](#环境要求)
- [快速启动](#快速启动)
- [配置说明](#配置说明)
- [部署文档](#部署文档)
- [开发者指南](#开发者指南)
- [性能优化](#性能优化)
- [贡献](#贡献)

## 项目简介

本项目基于 CloudWeGo Hertz 开发，目标是为福大学生提供学习资源分享与协作的 API 服务。项目集成了 MySQL、Redis、七牛云 OSS、SMTP 邮件、异步任务池等常见后端模块，方便二次开发与扩展。

## 主要特性

- 基于 Hertz 的高并发 HTTP 服务
- MySQL + GORM 数据持久化
- Redis 缓存与会话支持
- 七牛云 OSS 文件存储支持
- 邮件验证码/通知（SMTP）
- 异步数据库写入池，减少写操作延迟
- 完整的路由与权限中间件

## 技术栈

- Go (see go.mod) — 当前模块声明的 Go 版本见 `go.mod`，建议使用 go >= 1.25
- CloudWeGo Hertz
- GORM (MySQL)
- Redis (go-redis)
- 七牛云 OSS
- Viper 配置管理

## 环境要求

- Go >= 1.25（以项目根目录 `go.mod` 为准）
- Docker & Docker Compose（推荐用于本地快速启动依赖服务）
- MySQL（生产/本地部署）
- Redis（生产/本地部署）

## 快速启动

下面给出本地与 Docker 两种常用启动方式。

- 使用 Makefile（推荐）
  - 启动依赖服务（MySQL/Redis 等）：

```bash
make env-up
```

  - 停止依赖服务：

```bash
make env-down
```

- 本地直接运行（需先配置 `config/config.yaml`）：

```bash
# 下载依赖
go mod tidy

# 运行服务
go run main.go
```

- 使用 Docker Compose（仓库内提供 compose 文件）：

```bash
# 使用项目提供的 docker compose 启动（含服务）
make run
```

更多部署细节请查看 `deploy.md`。

## 配置说明

配置文件位于 `config/config.yaml`（仓库内有 `config/config.example.yaml` 作为示例）。启动前请复制并修改：

Windows:

```cmd
copy config\config.example.yaml config\config.yaml
notepad config\config.yaml
```

Linux/macOS:

```bash
cp config/config.example.yaml config/config.yaml
vim config/config.yaml
```

常用配置项：

```yaml
mysql:
  dsn: "root:password@tcp(127.0.0.1:3306)/LearnShare?charset=utf8mb4&parseTime=True&loc=Local"
redis:
  addr: "127.0.0.1:6379"
server:
  port: 8080
logger:
  dir: ./logs
  level: info
```

## 部署文档

详细部署步骤请参见 `deploy.md`（包含 Windows / Linux / Docker 部署流程，及常见问题排查）。

## 开发者指南

- 运行单元测试：

```bash
# 在项目根目录运行
go test ./... -v
```

- 生成/更新 Hertz 路由（如果需要）：

```bash
# 依赖仓库中的 idl 文件，示例：
hz update -idl idl/your_service.thrift
```

- 代码风格与 lint：仓库包含 `.golangci.yml`，推荐使用 `golangci-lint run` 进行静态检查。

## 性能优化

本项目包含一个基于 goroutine + channel 的异步写入池，用于减少写操作对请求响应的影响。详见 `biz/dal/db/ASYNC_USAGE.md`。

## 贡献

欢迎贡献代码与 issue。请先阅读 `CONTRIBUTING.md` 获取提交流程与代码规范。
