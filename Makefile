# Makefile for file_syn project

# 项目名称
BINARY_NAME=file_syn
CMD_PATH=./cmd/file_syn

# 版本信息
VERSION?=1.0.0
BUILD_TIME=$(shell date +%Y-%m-%d\ %H:%M:%S)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 编译标志（如果需要版本信息，可以在代码中定义相应的变量）
LDFLAGS=

# 构建目录
BUILD_DIR=./bin

# 支持的平台
PLATFORMS=linux windows darwin
ARCHITECTURES=amd64 arm64

# 默认目标
.DEFAULT_GOAL := help

.PHONY: help
help: ## 显示帮助信息
	@echo "可用的命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## 编译当前平台的可执行文件
	@echo "正在编译 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "编译完成: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: run
run: build ## 编译并运行程序（使用 config/config.json 或指定配置文件）
	@if [ -n "$(CONFIG)" ]; then \
		$(BUILD_DIR)/$(BINARY_NAME) $(CONFIG); \
	else \
		$(BUILD_DIR)/$(BINARY_NAME); \
	fi

.PHONY: test
test: ## 运行测试
	@echo "正在运行测试..."
	@go test -v ./...

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "正在运行测试并生成覆盖率报告..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

.PHONY: clean
clean: ## 清理编译产物
	@echo "正在清理..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean
	@echo "清理完成"

.PHONY: fmt
fmt: ## 格式化代码
	@echo "正在格式化代码..."
	@go fmt ./...
	@echo "格式化完成"

.PHONY: vet
vet: ## 运行 go vet 检查代码
	@echo "正在运行 go vet..."
	@go vet ./...
	@echo "检查完成"

.PHONY: lint
lint: fmt vet ## 运行代码检查（fmt + vet）

.PHONY: build-all
build-all: ## 编译所有平台的可执行文件
	@echo "正在编译所有平台..."
	@mkdir -p $(BUILD_DIR)
	@for GOOS in $(PLATFORMS); do \
		for GOARCH in $(ARCHITECTURES); do \
			if [ "$$GOOS" = "windows" ] && [ "$$GOARCH" = "arm64" ]; then \
				continue; \
			fi; \
			EXT=""; \
			if [ "$$GOOS" = "windows" ]; then \
				EXT=".exe"; \
			fi; \
			echo "编译 $$GOOS/$$GOARCH..."; \
			GOOS=$$GOOS GOARCH=$$GOARCH go build $(LDFLAGS) \
				-o $(BUILD_DIR)/$(BINARY_NAME)-$$GOOS-$$GOARCH$$EXT $(CMD_PATH) || exit 1; \
		done; \
	done
	@echo "所有平台编译完成，文件位于 $(BUILD_DIR)/ 目录"

.PHONY: build-linux
build-linux: ## 编译 Linux 版本
	@echo "正在编译 Linux 版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_PATH)
	@echo "Linux 版本编译完成"

.PHONY: build-windows
build-windows: ## 编译 Windows 版本
	@echo "正在编译 Windows 版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)
	@echo "Windows 版本编译完成"

.PHONY: build-darwin
build-darwin: ## 编译 macOS 版本
	@echo "正在编译 macOS 版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_PATH)
	@echo "macOS 版本编译完成"

.PHONY: install
install: build ## 安装到系统路径（需要 sudo 权限）
	@echo "正在安装到系统路径..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "安装完成: /usr/local/bin/$(BINARY_NAME)"

.PHONY: deps
deps: ## 下载依赖
	@echo "正在下载依赖..."
	@go mod download
	@go mod tidy
	@echo "依赖下载完成"

.PHONY: mod-verify
mod-verify: ## 验证依赖
	@echo "正在验证依赖..."
	@go mod verify
	@echo "验证完成"

.PHONY: all
all: clean lint test build-all ## 执行完整流程：清理、检查、测试、编译所有平台

