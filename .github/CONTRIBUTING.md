# 贡献指南

感谢你对 Content Alchemist 的兴趣！以下是参与贡献的指南。

## 提交信息规范

本项目使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范。提交信息应遵循以下格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

### 类型 (Type)

| 类型 | 说明 |
|------|------|
| `feat` | 新功能 ✨ |
| `fix` | Bug 修复 🐛 |
| `docs` | 文档更新 📝 |
| `style` | 代码格式调整（不影响代码含义）🎨 |
| `refactor` | 重构（既不是新功能也不是修复）♻️ |
| `perf` | 性能优化 ⚡ |
| `test` | 添加测试 ✅ |
| `chore` | 构建/工具链更新 🔧 |

### 示例

```bash
feat(writing): 添加 AI 生成大纲功能

fix(preview): 修复移动端预览样式问题

docs(readme): 更新安装说明
```

## 发布流程

### 自动发布（推荐）

使用发布脚本创建新版本：

```bash
# macOS/Linux
./scripts/release.sh patch    # 修复版本 1.0.0 -> 1.0.1
./scripts/release.sh minor    # 次要版本 1.0.0 -> 1.1.0
./scripts/release.sh major    # 主要版本 1.0.0 -> 2.0.0
./scripts/release.sh 1.2.3    # 指定版本号

# Windows
.\scripts\release.ps1 patch
.\scripts\release.ps1 minor
.\scripts\release.ps1 major
```

脚本会自动：
1. 更新版本号（wails.json 和 package.json）
2. 创建并推送标签
3. 触发 GitHub Actions 构建
4. 自动发布到 GitHub Releases

### 手动发布

如果需要手动触发构建：

1. 访问 [Actions 页面](https://github.com/yourusername/Content-Alchemist/actions)
2. 选择 "Release" 工作流
3. 点击 "Run workflow"
4. 输入版本号（如 v1.0.1）

## 开发流程

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: add some feature'`)
4. 推送分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## Pull Request 检查清单

- [ ] 代码遵循项目编码规范
- [ ] 已测试本地构建通过
- [ ] 提交信息符合规范
- [ ] 已更新相关文档（如适用）
