#!/bin/bash

# Content Alchemist 版本发布脚本
# 使用方法: ./scripts/release.sh [version_type]
# version_type: patch (默认), minor, major, alpha, beta

set -e

VERSION_TYPE=${1:-patch}
CURRENT_BRANCH=$(git branch --show-current)

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查是否在 main 或 master 分支
if [[ "$CURRENT_BRANCH" != "main" && "$CURRENT_BRANCH" != "master" ]]; then
    echo -e "${RED}❌ 错误: 必须在 main 或 master 分支上发布版本${NC}"
    echo "当前分支: $CURRENT_BRANCH"
    exit 1
fi

# 检查工作目录是否干净
if [[ -n $(git status --porcelain) ]]; then
    echo -e "${RED}❌ 错误: 工作目录不干净，请先提交或暂存更改${NC}"
    git status --short
    exit 1
fi

# 获取当前版本
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
echo -e "${BLUE}当前版本: $CURRENT_VERSION${NC}"

# 解析版本号
CURRENT_VERSION=${CURRENT_VERSION#v}
MAJOR=$(echo $CURRENT_VERSION | cut -d. -f1)
MINOR=$(echo $CURRENT_VERSION | cut -d. -f2)
PATCH=$(echo $CURRENT_VERSION | cut -d. -f3 | cut -d- -f1)

# 计算新版本
case $VERSION_TYPE in
    major)
        NEW_VERSION="$((MAJOR + 1)).0.0"
        ;;
    minor)
        NEW_VERSION="${MAJOR}.$((MINOR + 1)).0"
        ;;
    patch)
        NEW_VERSION="${MAJOR}.${MINOR}.$((PATCH + 1))"
        ;;
    alpha|beta)
        NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}-${VERSION_TYPE}.1"
        ;;
    *)
        # 如果输入是完整的版本号，直接使用
        if [[ "$VERSION_TYPE" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+ ]]; then
            NEW_VERSION=${VERSION_TYPE#v}
        else
            echo -e "${RED}❌ 错误: 无效的版本类型 '$VERSION_TYPE'${NC}"
            echo "支持的类型: major, minor, patch, alpha, beta, 或完整版本号 (如 1.2.3)"
            exit 1
        fi
        ;;
esac

NEW_TAG="v${NEW_VERSION}"

# 确认
echo ""
echo -e "${YELLOW}即将创建新版本: ${GREEN}${NEW_TAG}${NC}"
echo ""
read -p "确认发布? (y/N): " CONFIRM

if [[ "$CONFIRM" != "y" && "$CONFIRM" != "Y" ]]; then
    echo -e "${YELLOW}已取消发布${NC}"
    exit 0
fi

echo ""
echo -e "${BLUE}🚀 开始发布流程...${NC}"

# 更新 wails.json 中的版本
if [ -f "wails.json" ]; then
    echo -e "${BLUE}📦 更新 wails.json 版本...${NC}"
    sed -i.bak "s/\"productVersion\": \".*\"/\"productVersion\": \"${NEW_VERSION}\"/" wails.json
    rm -f wails.json.bak
fi

# 更新 frontend/package.json 中的版本
if [ -f "frontend/package.json" ]; then
    echo -e "${BLUE}📦 更新 package.json 版本...${NC}"
    sed -i.bak "s/\"version\": \".*\"/\"version\": \"${NEW_VERSION}\"/" frontend/package.json
    rm -f frontend/package.json.bak
fi

# 提交版本更新（如果有变更）
if [[ -n $(git status --porcelain) ]]; then
    echo -e "${BLUE}📝 提交版本更新...${NC}"
    git add wails.json frontend/package.json
    git commit -m "chore(release): bump version to ${NEW_TAG}"
    git push origin $CURRENT_BRANCH
fi

# 创建并推送标签
echo -e "${BLUE}🏷️ 创建标签 ${NEW_TAG}...${NC}"
git tag -a "${NEW_TAG}" -m "Release ${NEW_TAG}"
git push origin "${NEW_TAG}"

echo ""
echo -e "${GREEN}✅ 版本 ${NEW_TAG} 发布成功!${NC}"
echo ""
echo -e "${BLUE}🔄 GitHub Actions 正在构建发布...${NC}"
echo -e "查看进度: https://github.com/$(git remote get-url origin | sed 's/.*github.com[:/]//' | sed 's/\.git$//')/actions"
echo ""
echo -e "${YELLOW}注意:${NC}"
echo -e "  - Release 初始为 Draft 状态，构建完成后会自动发布"
echo -e "  - 你可以在 GitHub 上查看并编辑 Release 说明"
