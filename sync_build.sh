#!/bin/bash
# ============================================
# 🚀 同步 exposing_intranet 服务端到远程服务器
# 用法：
#   ./sync_build.sh server   # 只同步服务端
#   ./sync_build.sh all      # 同步服务端（默认）
# ============================================

# 获取当前脚本所在的绝对路径
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# 项目根目录
BASE_DIR="$SCRIPT_DIR"

# 本地目录
SERVER_BUILD="$BASE_DIR/build"

# 远程服务器信息
REMOTE_USER="root"
REMOTE_HOST="server"
REMOTE_DIR="/usr/local/go_project/intranet"

# 颜色定义
GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[1;33m"
RESET="\033[0m"

# 同步类型参数
TARGET=${1:-all}

# 检查本地目录存在性
check_dir() {
    local dir="$1"
    if [ ! -d "$dir" ]; then
        echo -e "${RED}❌ 目录不存在:${RESET} $dir"
        exit 1
    fi
}

sync_server() {
    check_dir "$SERVER_BUILD"
    echo -e "${YELLOW}🚀 正在同步服务端...${RESET}"
    rsync -avz --delete "$SERVER_BUILD/" "$REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR/server/"
    echo -e "${GREEN}✅ 服务端同步完成！${RESET}"
}

# 主逻辑
case "$TARGET" in
    server)
        sync_server
        ;;
    all)
        sync_server
        ;;
    *)
        echo -e "${RED}❌ 无效参数: '$TARGET'${RESET}"
        echo "正确用法: $0 [server|all]"
        exit 1
        ;;
esac

echo -e "${GREEN}🎉 同步任务完成${RESET}"
