GOCMD = go
GOBUILD = $(GOCMD) build
GOMOD = $(GOCMD) mod
BUILD_DIR = ./build
SERVICE_NAME = linux_service

# 构建标志
LDFLAGS = -ldflags="-s -w"
CGO = CGO_ENABLED=0
GOARCH = amd64

.PHONY: all build build-darwin build-linux build-windows clean install prepare run scripts

all: clean prepare build

# 安装依赖
install:
	$(GOMOD) tidy

# 准备构建目录和配置文件
prepare:
	@echo "准备构建目录..."
	@if [ ! -d "$(BUILD_DIR)" ]; then \
		mkdir -p $(BUILD_DIR); \
		echo "创建 $(BUILD_DIR) 目录"; \
	else \
		echo "$(BUILD_DIR) 目录已存在"; \
	fi
	@echo "复制配置文件和资源..."
	@if [ -d "website" ]; then \
		cp -r website $(BUILD_DIR)/; \
		if [ -d "$(BUILD_DIR)/website/logs" ]; then \
			find $(BUILD_DIR)/website/logs -type f -name "*.log" -delete; \
			echo "已排除日志文件"; \
		fi; \
		if [ -f "$(BUILD_DIR)/website/configs/client.key" ]; then \
			rm -f "$(BUILD_DIR)/website/configs/client.key"; \
			echo "已排除客户端凭据"; \
		fi; \
		if [ -f "$(BUILD_DIR)/website/configs/config.yaml" ]; then \
			sed -i 's/user: ".*"/user: "car"/' $(BUILD_DIR)/website/configs/config.yaml; \
			sed -i 's/pwd: ".*"/pwd: "Aa123098.."/' $(BUILD_DIR)/website/configs/config.yaml; \
			echo "已修改数据库配置: user=car, pwd=Aa123098.."; \
		fi; \
		 echo "website 目录复制完成"; \
	fi
	@echo "配置文件复制完成"

# 编译所有平台
build: build-darwin build-linux build-windows
	@echo "所有平台编译完成！"

# 编译 macOS
build-darwin:
	@echo "编译 macOS 版本..."
	GOOS=darwin $(CGO) GOARCH=$(GOARCH) $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/mac_service main.go
	@echo "macOS 版本编译完成: $(BUILD_DIR)/mac_service"

# 编译 Linux
build-linux:
	@echo "编译 Linux 版本..."
	GOOS=linux $(CGO) GOARCH=$(GOARCH) $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(SERVICE_NAME) main.go
	@echo "Linux 版本编译完成: $(BUILD_DIR)/$(SERVICE_NAME)"
	@$(MAKE) scripts
	@echo "✅ Linux 启动、停止、重启脚本已生成"

# 编译 Windows
build-windows:
	@echo "编译 Windows 版本..."
	GOOS=windows $(CGO) GOARCH=$(GOARCH) $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/windows_service.exe main.go
	@echo "Windows 版本编译完成: $(BUILD_DIR)/windows_service.exe"

# 清理构建目录
clean:
	@echo "清理构建目录..."
	@if [ -d "$(BUILD_DIR)" ]; then \
		rm -rf $(BUILD_DIR); \
		echo "$(BUILD_DIR) 目录已删除"; \
	else \
		echo "$(BUILD_DIR) 目录不存在"; \
	fi

# 运行程序
run:
	$(GOCMD) run main.go

# 生成启动、停止、重启脚本
scripts:
	@echo "生成启动脚本 start.sh ..."
	@echo '#!/bin/bash' > $(BUILD_DIR)/start.sh
	@echo 'APP_NAME="$(basename $(SERVICE_NAME))"' >> $(BUILD_DIR)/start.sh
	@echo 'SCRIPT_DIR="$$(cd "$$(dirname "$$0")" && pwd)"' >> $(BUILD_DIR)/start.sh
	@echo 'APP_PATH="$$SCRIPT_DIR/$$APP_NAME"' >> $(BUILD_DIR)/start.sh
	@echo 'LOG_DIR="$$SCRIPT_DIR/log"' >> $(BUILD_DIR)/start.sh
	@echo 'PID_FILE="$$SCRIPT_DIR/$$APP_NAME.pid"' >> $(BUILD_DIR)/start.sh
	@echo 'mkdir -p "$$LOG_DIR"' >> $(BUILD_DIR)/start.sh
	@echo 'if [ -f "$$PID_FILE" ] && kill -0 "$$(cat $$PID_FILE)" 2>/dev/null; then' >> $(BUILD_DIR)/start.sh
	@echo '  echo "⚠️ 服务已在运行 (PID=$$(cat $$PID_FILE))"; exit 1;' >> $(BUILD_DIR)/start.sh
	@echo 'fi' >> $(BUILD_DIR)/start.sh
	@echo 'nohup "$$APP_PATH" > "$$LOG_DIR/server.log" 2>&1 &' >> $(BUILD_DIR)/start.sh
	@echo 'echo $$! > "$$PID_FILE"' >> $(BUILD_DIR)/start.sh
	@echo 'echo "✅ 服务已启动 (PID=$$(cat $$PID_FILE))"' >> $(BUILD_DIR)/start.sh
	@echo 'echo "📄 日志文件: $$LOG_DIR/server.log"' >> $(BUILD_DIR)/start.sh
	@chmod +x $(BUILD_DIR)/start.sh

	@echo "生成停止脚本 stop.sh ..."
	@echo '#!/bin/bash' > $(BUILD_DIR)/stop.sh
	@echo 'APP_NAME="$(basename $(SERVICE_NAME))"' >> $(BUILD_DIR)/stop.sh
	@echo 'SCRIPT_DIR="$$(cd "$$(dirname "$$0")" && pwd)"' >> $(BUILD_DIR)/stop.sh
	@echo 'PID_FILE="$$SCRIPT_DIR/$$APP_NAME.pid"' >> $(BUILD_DIR)/stop.sh
	@echo 'if [ ! -f "$$PID_FILE" ]; then echo "❌ 未找到 PID 文件"; exit 1; fi' >> $(BUILD_DIR)/stop.sh
	@echo 'PID=$$(cat $$PID_FILE)' >> $(BUILD_DIR)/stop.sh
	@echo 'if kill -0 $$PID 2>/dev/null; then' >> $(BUILD_DIR)/stop.sh
	@echo '  echo "🛑 正在停止服务 (PID=$$PID)..."' >> $(BUILD_DIR)/stop.sh
	@echo '  kill $$PID' >> $(BUILD_DIR)/stop.sh
	@echo '  sleep 1' >> $(BUILD_DIR)/stop.sh
	@echo '  if kill -0 $$PID 2>/dev/null; then' >> $(BUILD_DIR)/stop.sh
	@echo '    echo "⚠️ 服务未正常停止，强制终止..."' >> $(BUILD_DIR)/stop.sh
	@echo '    kill -9 $$PID' >> $(BUILD_DIR)/stop.sh
	@echo '  fi' >> $(BUILD_DIR)/stop.sh
	@echo '  rm -f "$$PID_FILE"' >> $(BUILD_DIR)/stop.sh
	@echo '  echo "✅ 服务已停止"' >> $(BUILD_DIR)/stop.sh
	@echo 'else' >> $(BUILD_DIR)/stop.sh
	@echo '  echo "⚠️ 进程 $$PID 不存在，清理 PID 文件"' >> $(BUILD_DIR)/stop.sh
	@echo '  rm -f "$$PID_FILE"' >> $(BUILD_DIR)/stop.sh
	@echo 'fi' >> $(BUILD_DIR)/stop.sh
	@chmod +x $(BUILD_DIR)/stop.sh

	@echo "生成重启脚本 restart.sh ..."
	@echo '#!/bin/bash' > $(BUILD_DIR)/restart.sh
	@echo 'SCRIPT_DIR="$$(cd "$$(dirname "$$0")" && pwd)"' >> $(BUILD_DIR)/restart.sh
	@echo 'bash "$$SCRIPT_DIR/stop.sh"' >> $(BUILD_DIR)/restart.sh
	@echo 'sleep 1' >> $(BUILD_DIR)/restart.sh
	@echo 'bash "$$SCRIPT_DIR/start.sh"' >> $(BUILD_DIR)/restart.sh
	@chmod +x $(BUILD_DIR)/restart.sh
