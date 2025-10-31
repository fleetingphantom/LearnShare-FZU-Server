
MODULE = judgeMore_server

DIR = $(shell pwd)
CMD = $(DIR)/cmd
CONFIG_PATH = $(DIR)/config
IDL_PATH = $(DIR)/idl
OUTPUT_PATH = $(DIR)/output
API_PATH= $(DIR)/cmd/api

# 启动必要的环境，比如 etcd、mysql
.PHONY: env-up
env-up:
	@ docker compose -f ./docker/docker-compose.yml up -d

# 关闭必要的环境，但不清理 data（位于 docker/data 目录中）
.PHONY: env-down
env-down:
	@ cd ./docker && docker compose down

# 基于 idl 生成相关的 go 语言描述文件
.PHONY: kitex-gen-%
kitex-gen-%:
	@ kitex -module "${MODULE}" \
		-thrift no_default_serdes \
		${IDL_PATH}/$*.thrift
	@ go mod tidy

# 生成基于 Hertz 的脚手架
.PHONY: hz-%
hz-%:
	hz update -idl ${IDL_PATH}/api/$*.thrift

# 清除所有的构建产物
.PHONY: clean
clean:
	@find . -type d -name "output" -exec rm -rf {} + -print

# 清除所有构建产物、compose 环境和它的数据
.PHONY: clean-all
clean-all: clean
	@echo "$(PREFIX) Checking if docker-compose services are running..."
	@docker-compose -f ./docker/docker-compose.yml ps -q | grep '.' && docker-compose -f ./docker/docker-compose.yml down || echo "$(PREFIX) No services are running."
	@echo "$(PREFIX) Removing docker data..."
	rm -rf ./docker/data
