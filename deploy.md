# LearnShare-FZU-Server — 部署指南

该文档给出在常见平台（Windows / Linux / Docker）上部署本项目的步骤、验证方法与故障排查建议。请优先参考仓库根目录的 `Makefile` 和 `docker/` 下的 compose 文件来启动依赖服务（MySQL / Redis）。

## 先决条件
- Go：以仓库根目录 `go.mod` 中的 `go` 指令为准（建议使用相同或更高版本）。
- MySQL、Redis（或使用 Docker 提供的依赖容器）。
- 推荐在开发/测试机安装 Docker 与 Docker Compose（Windows 用户可使用 Docker Desktop / WSL2）。


## 快速目录（快速跳转）
- Windows 本地部署（简洁）
- Linux / macOS 本地部署（简洁）
- Docker Compose（推荐用于本地集成测试）
- 数据库导入 / 迁移
- 验证服务是否启动
- 常见问题与排查
- 生产部署提示（systemd、日志、备份）


## Windows 本地部署（简洁）
适用于开发或在 Windows 机器上直接运行服务的场景。

1. 准备环境
   - 安装 Go（将 Go 添加到 PATH）。
   - 安装 MySQL/Redis，或使用 Docker（推荐）。

2. 获取源码并准备配置

在 cmd.exe 中：

```cmd
git clone https://github.com/2451965602/LearnShare-FZU-Server.git
cd LearnShare-FZU-Server
copy config\config.example.yaml config\config.yaml
notepad config\config.yaml
```

在 `config/config.yaml` 中配置 `mysql.dsn`、`redis.addr`、OSS、SMTP 等字段。

3. 启动依赖（推荐使用仓库的 Makefile）

如果系统有 `make`：

```cmd
make env-up    REM 使用仓库中的 docker/docker-compose-env.yml 启动 MySQL/Redis 等依赖
```

如果没有 `make`，在 Windows 上可以直接使用 Docker Compose（PowerShell）：

```powershell
docker-compose -f docker\docker-compose-env.yml up -d
```

4. 下载依赖并运行服务

```cmd
go mod download
go run main.go
```

说明：如果你希望以可执行文件运行（生产/长期运行）：

```cmd
go build -o learnshare-server.exe main.go
start /B learnshare-server.exe
```


## Linux / macOS 部署（简洁）

1. 准备 Go、MySQL、Redis（或使用 Docker）。
2. 克隆仓库并复制配置：

```bash
git clone https://github.com/2451965602/LearnShare-FZU-Server.git
cd LearnShare-FZU-Server
cp config/config.example.yaml config/config.yaml
vim config/config.yaml
```

3. 启动依赖服务（Makefile）：

```bash
make env-up
```

或直接用 docker-compose：

```bash
docker-compose -f docker/docker-compose-env.yml up -d
```

4. 运行服务（开发/测试）：

```bash
go mod download
nohup go run main.go > server.log 2>&1 &
```

构建二进制并用 systemd 启动（见生产部署示例）：

```bash
go build -o learnshare-server main.go
sudo cp learnshare-server /usr/local/bin/
```


## 使用 Docker Compose（推荐用于本地集成测试）
仓库内包含 `docker/docker-compose-env.yml`（仅用于依赖服务）和可选的 `docker/docker-compose.yml`（若已配置完整服务）。

常用命令：

```bash
# 启动依赖服务（MySQL/Redis）
make env-up

# 构建服务镜像（若仓库提供 Dockerfile / Makefile 的 build 任务）
make build

# 启动完整服务（如 docker/docker-compose.yml 已配置）
make run

# 停止并移除依赖
make env-down
```

如果没有 Makefile，直接使用 docker-compose：

```bash
# 启动依赖
docker-compose -f docker/docker-compose-env.yml up -d

# 查看容器
docker ps

# 查看容器日志
docker logs -f <container-name>
```


## 数据库导入 / 迁移

仓库包含 `sql/init.sql` 和 `sql/test_data.sql`。如果没有自动迁移工具，你可以手动导入初始化 SQL：

使用本地 MySQL 客户端：

```bash
mysql -u root -p < sql/init.sql
# 或者
cat sql/init.sql | mysql -u root -p your_database
```

如果你用的是 Docker 中的 mysql 容器（例如通过 `make env-up` 启动），先找到容器名再导入：

```bash
# 示例
docker exec -i <mysql-container> sh -c 'mysql -u root -p"root_password" your_database' < sql/init.sql
```

关于迁移：仓库中可能未包含自动迁移命令（例如 goose、migrate 等），建议：
- 查阅项目代码与 `sql/` 目录，按需要手动导入 SQL
- 或在部署流程中集成 goose / golang-migrate 并维护迁移文件


## 验证服务是否启动（常用检查）

默认端口以 `config/config.yaml` 中 `server.port` 为准，示例为 8080：

- HTTP 健康检查：

```bash
curl -v http://localhost:8080/ping
# 期望返回：pong
```

- 查看日志：

Linux:

```bash
tail -f server.log
# 或项目 logs 目录下的日志文件
```

Windows(cmd.exe)：

```cmd
type logs\app.log
```

- 检查端口占用（Windows）：

```cmd
netstat -ano | findstr :8080
```

- 检查容器与日志（Docker）：

```bash
docker ps
docker logs -f <container-name>
```


## 常见问题与排查

- 无法连接 MySQL：
  - 检查 `config/config.yaml` 中 `mysql.dsn` 是否正确；确认数据库服务已启动且网络可达。
  - 如果使用 Docker，确认容器名/端口与 DSN 中主机一致。

- Redis 连接失败：
  - 检查 `redis.addr`、认证信息和防火墙
  - 使用 `redis-cli`（或 docker exec 到 redis 容器）进行连接测试

- 端口被占用：
  - 修改 `server.port` 或停止占用进程（`netstat` / `ss` 查找进程并结束）。

- 日志无输出：
  - 检查 `logger.dir` 的写权限和磁盘空间。

- 缺少迁移脚本：
  - 使用仓库的 `sql/` 中的初始化脚本导入，或在部署流程中引入标准迁移工具。


## 生产部署示例（systemd）

下面是一个最小的 systemd unit 示例（Ubuntu/CentOS 等）：

创建 `/etc/systemd/system/learnshare.service`：

```ini
[Unit]
Description=LearnShare-FZU-Server
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
ExecStart=/usr/local/bin/learnshare-server --config /opt/learnshare/config/config.yaml
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启用并启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable learnshare
sudo systemctl start learnshare
sudo journalctl -u learnshare -f
```


## 备份与运维建议
- 为 MySQL/Redis 配置定期备份策略。
- 将日志写入到集中式日志系统（ELK / Promtail + Loki 等）。
- 在容器化部署中，务必为数据库配置持久化卷。
- 生产环境使用固定版本的二进制或镜像（不要直接使用 `latest`）。


---

如果你希望我把 `deploy.md` 中的某一部分扩展为更详细的操作脚本（例如：自动化导入 SQL 的 PowerShell / bash 脚本，或 `systemd` 单元的环境文件示例），告诉我需要哪个部分，我会继续完善并测试相应命令示例。
