# Content-Alchemist Pro

基于 **Eino ADK + Wails + Vue3** 的AI写作助手，专为一人公司创作者设计的本地优先写作工具。

## ✨ 核心特性

### 🎯 产品定位
面向一人公司创作者的本地优先AI写作助手，核心workflow是"输入标题→AI协作创作→公众号排版→一键发布"

### 🛠️ 技术栈
- **前端**: Vue 3 + TypeScript + Vite + Pinia + Element Plus
- **后端**: Go 1.21+ + Wails v2
- **存储**: SQLite (wails自带) + 本地文件系统
- **AI引擎**: DeepSeek/OpenAI API

### 🎨 主要功能

#### 1. 智能引导系统 (Onboarding)
- 首次启动自动检测配置
- 3步引导完成初始化：欢迎页 → AI配置 → 公众号定位
- 支持本地模式（跳过AI配置）

#### 2. 三栏编辑器架构
- 左栏：大纲导航 + 素材库
- 中栏：Markdown编辑器
- 右栏：实时预览（手机/桌面模式）

#### 3. 行内AI助手系统
- **斜杠命令**：输入 `/` 弹出AI命令面板
  - `/ai` - 召唤AI助手
  - `/outline` - 生成大纲
  - `/expand` - 扩写
  - `/polish` - 润色
  - `/shorter` - 精简
  - `/title` - 生成标题
- **选中文本AI菜单**：悬浮工具条提供润色、口语化、加数据、重写、续写等功能

#### 4. 智能大纲生成
- AI基于标题自动生成结构化大纲
- 大纲节点可拖拽排序
- 支持分段生成正文
- 字数进度追踪

#### 5. 手机预览与主题系统
- iPhone 14 Pro 模拟器预览
- 一键切换主题：科技蓝、简约白、活力橙
- Markdown → 公众号HTML转换
- 复制带样式的HTML到公众号后台

#### 6. 素材库管理
- 支持片段、数据、金句、历史四种类型
- 拖拽插入到编辑器
- 标签筛选和搜索
- 使用次数统计

#### 7. 发布助手
- 发布前检查清单
- 文章统计（字数、阅读时间）
- 复制带样式HTML
- 下载Markdown

#### 8. 命令面板与快捷键
- `Cmd/Ctrl + K` 唤起命令面板
- 全局快捷键支持
  - `Cmd/Ctrl + N` - 新建文章
  - `Cmd/Ctrl + S` - 保存
  - `Cmd/Ctrl + Shift + P` - 发布

## 📁 项目结构

```
content-alchemist/
├── build/                      # 构建输出
├── frontend/                   # Vue前端
│   ├── src/
│   │   ├── components/         # 组件
│   │   │   ├── common/         # 通用组件
│   │   │   ├── editor/         # 编辑器组件
│   │   │   ├── preview/        # 预览组件
│   │   │   └── ai/             # AI组件
│   │   ├── composables/        # 组合式函数
│   │   ├── stores/             # Pinia状态管理
│   │   ├── types/              # TypeScript类型
│   │   ├── utils/              # 工具函数
│   │   └── views/              # 页面视图
│   └── package.json
├── internal/                   # Go后端
│   ├── agent/                  # AI Agents
│   ├── service/                # 业务逻辑
│   ├── repository/             # 数据访问
│   └── model/                  # 数据模型
├── main.go                     # 应用入口
├── wails.json                  # Wails配置
└── go.mod                      # Go模块
```

## 🚀 开发

### 前置要求
- Go 1.21+
- Node.js 18+
- Wails CLI v2

### 安装依赖

```bash
# 安装前端依赖
cd frontend
npm install

# 安装 Go 依赖
go mod tidy
```

### 开发模式

```bash
# 启动开发服务器
wails dev
```

### 构建

```bash
# 构建生产版本
wails build

# 构建 Windows 安装包
wails build -platform windows/amd64 -nsis
```

## 📦 发布

构建完成后，可执行文件位于：
- Windows: `build/bin/Content-Alchemist.exe`
- macOS: `build/bin/Content-Alchemist.app`
- Linux: `build/bin/content-alchemist`

## 🔧 配置

应用配置存储在用户目录下的 `.content-alchemist/config.json`：
- API Base URL
- API Key
- 模型选择
- 温度参数
- 写作风格偏好

## 📝 开发任务清单

- [x] Task 1: 项目骨架与Eino集成
- [x] Task 2: Onboarding引导系统
- [x] Task 3: 核心编辑器骨架
- [x] Task 4: 行内AI助手系统
- [x] Task 5: 智能大纲生成
- [x] Task 6: 手机预览与主题系统
- [x] Task 7: 素材库管理
- [x] Task 8: 发布助手
- [x] Task 9: 快捷键与命令面板
- [x] Task 10: 性能优化与构建配置

## 📄 许可证

MIT License

## 🙏 致谢

- [Wails](https://wails.io/) - Go + Web 技术构建桌面应用
- [Vue.js](https://vuejs.org/) - 渐进式JavaScript框架
- [Element Plus](https://element-plus.org/) - Vue3组件库
- [Marked](https://marked.js.org/) - Markdown解析器
