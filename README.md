# LearnShare-FZU-Server

基于 CloudWeGo Hertz 的 Go 后端服务，作为福州大学学生学习分享平台的 API 层实现。该服务集成了 MySQL、Redis、七牛云 OSS、SMTP、异步写入池等常见后端组件，面向高并发场景并便于二次开发与部署。


## 目录

- [项目简介](#项目简介)
- [主要特性](#主要特性)
- [技术栈](#技术栈)
- [环境要求](#环境要求)
- [快速启动](#快速启动)
- [配置说明](#配置说明)
- [部署（Docker / 本地）](#部署docker--本地)
- [开发者指南](#开发者指南)
- [性能优化](#性能优化)
- [常见问题与诊断](#常见问题与诊断)
- [贡献](#贡献)
- [License](#license)


## 项目简介

LearnShare-FZU-Server 是一个后端 API 服务，基于 CloudWeGo 的 Hertz 框架实现，提供用户、资源、课程、收藏、审核等一系列后台能力，便于前端或移动端调用。项目内置常用中间件（JWT 鉴权、权限控制、请求日志、Turnstile 验证等），并提供异步写入池以提高写操作吞吐。


## 主要特性

- 高性能 HTTP 服务（基于 CloudWeGo Hertz）
- MySQL + GORM 数据持久化
- Redis 缓存与会话支持
- 七牛云（或其他 OSS）文件存储抽象
- SMTP 邮件通知（验证码、通知）
- 异步数据库写入池（减少写入延迟） — 见 `biz/dal/db/ASYNC_USAGE.md`
- 完整的路由与权限中间件（参见 `router/` 与 `pkg/permissions`）
- 测试覆盖若干业务模块（部分测试文件位于 `biz/` 与 `service/` 目录）


## 技术栈

- Go（以 `go.mod` 为准，建议使用 go >= 1.20+）
- CloudWeGo Hertz（HTTP 框架）
- GORM（MySQL ORM）
- Redis（go-redis）
- 七牛云 OSS（或其他 OSS，请在配置中替换）
- Viper（配置管理）


## 环境要求

- Go（以 `go.mod` 为准）
- MySQL（建议 5.7+ / 8.x）
- Redis（或 Kvrocks/Upstash 根据需要）
- Docker & Docker Compose（推荐用于本地依赖服务）


## 快速启动

下面给出本地开发与 Docker 两种常用启动方式。以下命令在仓库根目录执行。

- 使用 Makefile（推荐，用于启动本地依赖服务）：

Windows(cmd.exe)：
```
make env-up
```

停止依赖服务：
```
make env-down
```

- 本地直接运行（需先配置 `config/config.yaml`）：

```bash
# 拉取依赖
go mod tidy

# 运行服务（开发模式）
go run main.go
```

- 使用 Docker Compose（仓库内提供 compose 文件或参考 `docker/` 目录）：

```bash
# 使用项目提供的 docker compose 启动（含服务）
make run
```

更多部署细节请查看 `deploy.md`。


## 配置说明

配置文件位于 `config/config.yaml`（仓库中提供 `config/config.example.yaml` 作为示例）。启动前请复制并修改：

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

常用配置示例（摘录）：

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
oss:
  provider: qiniu
  bucket: your-bucket
  access_key: xxxxx
  secret_key: xxxxx
mail:
  smtp_host: smtp.example.com
  smtp_port: 587
  username: noreply@example.com
  password: xxxxx
```

配置项说明：
- mysql.dsn：MySQL 连接字符串（包含用户名/密码/地址/数据库名）
- redis.addr：Redis 地址
- server.port：服务监听端口
- logger.dir/level：日志目录与日志级别
- oss.*：OSS 存储相关配置（七牛/阿里云等，根据 provider 不同字段略有差异）
- mail.*：SMTP 相关配置


## 部署（Docker / 本地）

- 使用 Docker Compose：仓库内包含 `docker/` 目录与相关 compose 示例，可通过 `make run` 启动（看项目 Makefile）。
- 生产环境建议：
  - 为 MySQL/Redis 配置持久化卷
  - 使用固定版本的镜像或构建的二进制而非 `latest`
  - 配置合理的日志轮转/监控

请参阅 `deploy.md` 获取按平台的详细部署步骤与排查指南。


## 开发者指南

- 运行单元测试：

```bash
# 在项目根目录运行
go test ./... -v
```

- 代码生成 / 路由更新：

本项目使用 Hertz 路由工具部分自动生成（参考 `router_gen.go`、`idl/` 目录）。如需更新路由：

```bash
# 例：更新 idl 后生成路由
hz update -idl idl/your_service.thrift
```

- 静态检查：

仓库可能包含 `.golangci.yml`，建议使用 `golangci-lint run` 进行静态检查。

- 常见开发命令：

```bash
# 整理模块依赖
go mod tidy

# 构建可执行文件
go build -o learnshare-server main.go
```


## 性能优化

项目包含一个基于 goroutine + channel 的异步写入池，用于减少写操作对请求响应的影响。详见 `biz/dal/db/ASYNC_USAGE.md`。


## 常见问题与诊断

- 无法连接数据库：检查 `config/config.yaml` 中的 MySQL DSN 与网络连通性。
- Redis 连接失败：确认 `redis.addr` 与 Redis 服务是否已启动。
- 日志无输出：确认 `logger.dir` 存在且进程有写权限。

更多排查步骤请参见 `deploy.md`。


## 贡献

欢迎贡献代码、issue 与 PR。贡献前请阅读项目根目录下的 `CONTRIBUTING.md`（若无该文件，请在 PR 中说明开发与测试步骤）。


## License

本项目采用 MIT 许可证（详见 LICENSE），欢迎自由使用与贡献。
