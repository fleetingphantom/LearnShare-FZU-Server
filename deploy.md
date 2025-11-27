# LearnShare-FZU-Server 详细部署步骤

## 重要说明（先读）
- 请优先参考仓库根目录的 `Makefile` 和 `docker/` 下的 compose 文件来启动依赖服务（MySQL / Redis）。
- 项目 Go 版本以 `go.mod` 中的 `go` 指令为准（当前仓库显示 go 1.25.x）。

## Windows 本地部署（简洁版）

### 1. 环境准备
- 安装 Go（建议使用和 `go.mod` 一致或更高的版本）:
  - 安装时请把 Go 添加到 PATH。
- 安装并启动 MySQL 与 Redis：
  - 推荐使用 Docker（见下方 Docker 部署）或者按本文件中手动安装说明进行安装并启动。

### 2. 获取源码与配置
```cmd
//克隆仓库
git clone https://github.com/2451965602/LearnShare-FZU-Server.git
cd LearnShare-FZU-Server

//复制示例配置（Windows）
copy config\config.example.yaml config\config.yaml
notepad config\config.yaml
```

在 `config/config.yaml` 中填入数据库、Redis、OSS、SMTP 等配置。

### 3. 启动依赖服务（推荐使用 Makefile）
```cmd
// 启动 MySQL/Redis 等依赖（使用仓库中 docker-compose-env.yml）
make env-up

// 停止依赖
make env-down
```

如果不使用 Docker，请确保 MySQL/Redis 已经运行并与 `config.yaml` 中配置的一致。

### 4. 下载依赖并运行服务
```cmd
// 下载 go 依赖
go mod download

// 运行服务（开发模式）
go run main.go
```

> 可选：项目内可能没有数据库迁移脚本命令（例如 `cmd/migrate`），如果你的部署流程需要执行迁移，请在仓库中确认迁移脚本位置或使用你习惯的迁移工具（如 goose / golang-migrate）。

## Linux / macOS 部署（简洁版）

1. 准备 Go、MySQL、Redis（或使用 Docker）
2. 克隆仓库并复制配置：

```bash
git clone https://github.com/2451965602/LearnShare-FZU-Server.git
cd LearnShare-FZU-Server
cp config/config.example.yaml config/config.yaml
vim config/config.yaml
```

3. 启动依赖：

```bash
make env-up   # 启动 MySQL/Redis 等依赖
```

4. 下载依赖并后台启动服务：

```bash
go mod download
nohup go run main.go > server.log 2>&1 &
```

查看日志：

```bash
tail -f server.log
```

## 使用 Docker Compose（推荐用于本地集成测试）

仓库包含 `docker/docker-compose-env.yml`（仅用于依赖服务）和 `docker/docker-compose.yml`（可能用于完整服务部署）。

```bash
# 启动依赖
make env-up

# 构建镜像
make build

# 启动整套服务（若 docker-compose.yml 已配置）
make run
```

## 验证服务是否启动

- 健康检查：访问 `http://localhost:8080/ping`（或 `server.port` 指定的端口）应返回 `pong`。
- 日志：查看 `logs` 目录或控制台输出确认无致命错误。

## 常见问题与排查

- 无法连接 MySQL：检查 `config/config.yaml` 中 `mysql.dsn` 是否正确，数据库是否已启动，网络与防火墙设置。
- Redis 连接失败：检查 `redis.addr`、`redis.password` 配置与服务状态。
- 端口被占用：修改 `config/config.yaml` 中 `server.port` 或停止占用进程。
- 缺少迁移脚本：仓库中未必包含数据库迁移工具，请使用外部工具或联系项目维护者。
