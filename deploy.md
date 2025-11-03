# LearnShare-FZU-Server 详细部署步骤

## Windows 部署

### 1. 环境准备
- **安装 Go**
    1. 下载 Go 1.23.6+ 安装包：[https://go.dev/dl/](https://go.dev/dl/)
    2. 安装时勾选 "Add Go to PATH"，或手动配置环境变量：
        - GOROOT: `C:\Program Files\Go`
        - GOPATH: `%USERPROFILE%\go`
        - 将 `%GOROOT%\bin` 和 `%GOPATH%\bin` 添加到系统 PATH

- **安装 MySQL**
    1. 下载 MySQL 9.0.1+：[https://dev.mysql.com/downloads/mysql/](https://dev.mysql.com/downloads/mysql/)
    2. 安装时选择 "Server only"，设置 root 密码
    3. 配置 MySQL 服务自启动：
       ```bash
       sc config mysql start= auto
       net start mysql
       ```

- **安装 Redis**

    1. 下载 Redis 6.0+：[https://github.com/tporadowski/redis/releases](https://github.com/tporadowski/redis/releases)
    2. 解压到 `C:\Redis`，执行以下命令安装服务：
       ```bash
       cd C:\Redis
       redis-server --service-install redis.windows.conf --loglevel verbose
       net start redis
       ```

### 2. 源码获取与配置
```bash
# 克隆代码仓库
git clone https://github.com/2451965602/LearnShare-FZU-Server.git
cd LearnShare-FZU-Server

# 复制配置文件模板
copy config\config.example.yaml config\config.yaml

# 编辑配置文件（需修改以下内容）
notepad config\config.yaml
```

配置文件关键项：
```yaml
mysql:
  dsn: "root:password@tcp(127.0.0.1:3306)/LearnShare?charset=utf8mb4&parseTime=True&loc=Local"
redis:
  addr: "127.0.0.1:6379"
  password: ""
server:
  port: 8080
```

### 3. 依赖安装与启动
```bash
# 安装项目依赖
go mod download

# 初始化数据库（如提供迁移脚本）
go run cmd/migrate/main.go

# 启动服务
go run main.go
```

## Linux 部署

### 1. 环境准备
- **安装 Go**
  ```bash
  # 下载安装包
  wget https://go.dev/dl/go1.23.6.linux-amd64.tar.gz
  sudo tar -C /usr/local -xzf go1.23.6.linux-amd64.tar.gz
  
  # 配置环境变量
  echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
  echo 'export GOPATH=$HOME/go' >> ~/.bashrc
  source ~/.bashrc
  ```

- **安装 MySQL**
  ```bash
  # Ubuntu/Debian
  sudo apt update
  sudo apt install mysql-server-9.0
  
  # CentOS/RHEL
  sudo dnf install mysql-server
  
  # 启动服务并设置自启
  sudo systemctl enable --now mysql
  
  # 安全配置
  sudo mysql_secure_installation
  ```

- **安装 Redis**
  ```bash
  # Ubuntu/Debian
  sudo apt install redis-server
  
  # CentOS/RHEL
  sudo dnf install redis
  
  # 启动服务
  sudo systemctl enable --now redis
  ```

### 2. 源码获取与配置
```bash
# 克隆代码仓库
git clone https://github.com/2451965602/LearnShare-FZU-Server.git
cd LearnShare-FZU-Server

# 复制配置文件
cp config/config.example.yaml config/config.yaml

# 编辑配置
vim config/config.yaml
```

### 3. 服务启动
```bash
# 安装依赖
go mod download

# 启动服务（后台运行）
nohup go run main.go > server.log 2>&1 &

# 查看日志
tail -f server.log
```

## 验证部署
1. 访问 API 服务：[http://localhost:8080/health](http://localhost:8080/health)
2. 检查返回状态：`{"status":"ok"}`
3. 查看服务日志确认无错误信息
