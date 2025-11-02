# LearnShare-FZU-Server

基于 CloudWeGo Hertz 框架的 Go 语言后端服务项目，为福大学生学习分享平台提供 API 服务。

## 项目简介

本项目采用 CloudWeGo 生态的 Hertz 框架开发，提供高性能的 HTTP 服务。项目集成了 MySQL、Redis、阿里云 OSS、SMTP 邮件服务等功能模块。

## 技术栈

- **框架**: CloudWeGo Hertz
- **数据库**: MySQL 9.0.1
- **缓存**: Redis
- **对象存储**: 七牛云 OSS
- **邮件服务**: SMTP
- **配置管理**: Viper
- **ORM**: GORM

## 环境要求

- Go 1.23.6+
- Docker & Docker Compose
- MySQL 9.0.1+ (如果本地部署)
- Redis 6.0+ (如果本地部署)

## 快速启动

### 方式一：使用 Docker Compose (推荐)

#### Windows 系统

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd LearnShare-FZU-Server
   ```

2. **配置环境变量**
   ```bash
   # 复制配置文件模板
   copy config\config.example.yaml config\config.yaml

   # 根据需要修改 config/config.yaml 中的配置
   ```

3. **启动依赖服务**
   ```bash
   make env-up
   ```
   或者直接使用：
   ```bash
   docker compose -f ./docker/docker-compose.yml up -d
   ```

4. **安装项目依赖**
   ```bash
   go mod download
   ```

5. **启动服务**
   ```bash
   go run main.go
   ```

#### Linux 系统

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd LearnShare-FZU-Server
   ```

2. **配置环境变量**
   ```bash
   # 复制配置文件模板
   cp config/config.example.yaml config/config.yaml

   # 根据需要修改 config/config.yaml 中的配置
   vim config/config.yaml
   ```

3. **启动依赖服务**
   ```bash
   make env-up
   ```
   或者直接使用：
   ```bash
   docker compose -f ./docker/docker-compose.yml up -d
   ```

4. **安装项目依赖**
   ```bash
   go mod download
   ```

5. **启动服务**
   ```bash
   go run main.go
   ```

### 方式二：本地部署

#### Windows 系统

1. **安装并配置 MySQL**
   - 下载并安装 MySQL 9.0.1+
   - 创建数据库：`CREATE DATABASE LearnShare;`
   - 配置用户权限

2. **安装并配置 Redis**
   - 下载并安装 Redis 6.0+
   - 启动 Redis 服务

3. **配置项目**
   ```bash
   copy config\config.example.yaml config\config.yaml
   # 修改配置文件中的数据库和 Redis 连接信息
   ```

4. **启动服务**
   ```bash
   go mod download
   go run main.go
   ```

#### Linux 系统

1. **安装 MySQL**
   ```bash
   # Ubuntu/Debian
   sudo apt update
   sudo apt install mysql-server-9.0

   # CentOS/RHEL
   sudo yum install mysql-server

   # 启动服务
   sudo systemctl start mysql
   sudo systemctl enable mysql
   ```

2. **安装 Redis**
   ```bash
   # Ubuntu/Debian
   sudo apt install redis-server

   # CentOS/RHEL
   sudo yum install redis

   # 启动服务
   sudo systemctl start redis
   sudo systemctl enable redis
   ```

3. **配置数据库**
   ```bash
   mysql -u root -p
   CREATE DATABASE LearnShare;
   CREATE USER 'learnshare'@'localhost' IDENTIFIED BY 'your_password';
   GRANT ALL PRIVILEGES ON LearnShare.* TO 'learnshare'@'localhost';
   FLUSH PRIVILEGES;
   ```

4. **配置项目**
   ```bash
   cp config/config.example.yaml config/config.yaml
   vim config/config.yaml
   # 修改数据库和 Redis 连接配置
   ```

5. **启动服务**
   ```bash
   go mod download
   go run main.go
   ```

## 配置文件说明

主要配置文件位于 `config/config.yaml`，包含以下配置项：

- `mysql`: 数据库连接配置
- `redis`: Redis 连接配置
- `oss`: 阿里云 OSS 配置
- `smtp`: 邮件服务配置
- `verify`: 验证码配置

## 常用命令

```bash
# 启动依赖环境
make env-up

# 停止依赖环境
make env-down

# 清理构建产物
make clean

# 完全清理（包括数据）
make clean-all

# 运行服务
go run main.go

# 构建二进制文件
go build -o output/bin/main main.go
```

## 服务端口

默认服务端口为 `8080`，可在配置文件中修改。

## API 文档

启动服务后，可通过以下地址访问：
- API 服务: `http://localhost:8080`

## 开发说明

- 项目使用 CloudWeGo Hertz 框架
- 数据库操作使用 GORM ORM 框架
- 配置管理使用 Viper
- 支持 Docker 容器化部署

## 故障排除

1. **端口被占用**: 检查 8080、3306、6379 端口是否被其他服务占用
2. **数据库连接失败**: 检查 MySQL 服务是否启动，配置是否正确
3. **Redis 连接失败**: 检查 Redis 服务是否启动
4. **依赖下载失败**: 检查网络连接，或使用代理

## 许可证

[请根据项目实际情况添加许可证信息]
