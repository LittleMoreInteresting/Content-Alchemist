# Content Alchemist 📝✨

一款现代化的智能写作工具，专为内容创作者设计。支持 AI 辅助写作、微信公众号排版美化、一键发布等功能。

![wails](https://img.shields.io/badge/Wails-v2.11.0%20-blue)
![Vue3](https://img.shields.io/badge/Vue-3.x-4FC08D?style=flat-square&logo=vue.js)
![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?style=flat-square&logo=typescript)
![Vite](https://img.shields.io/badge/Vite-5.x-646CFF?style=flat-square&logo=vite)

## ✨ 核心功能

### 🖊️ 智能写作
- **AI 生成大纲**：输入标题，自动生成文章大纲和候选标题
- **AI 生成文章**：基于大纲一键生成完整文章内容
- **公众号定位**：自定义写作风格和内容定位
- **智能保存**：自动保存，支持历史版本管理

### 📱 手机端预览
- **真实手机模拟**：iPhone 风格预览框架
- **微信公众号排版**：完美还原公众号阅读体验
- **实时预览**：编辑内容即时同步预览

### 🎨 微信公众号排版
- **标题样式**：H1 居中双彩边框、H2 绿色左边框、H3 蓝色左边框
- **首字下沉**：文章首段自动首字下沉效果
- **代码块**：深色主题 + macOS 风格窗口按钮
- **引用块**：渐变背景 + 引号装饰
- **表格**：绿色表头 + 斑马纹 + 悬停效果
- **分隔线**：渐变 + 菱形装饰

### 📋 一键复制到公众号
- **富文本复制**：保留完整样式直接粘贴到公众号编辑器
- **内联样式转换**：自动将 CSS 转换为元素内联样式
- **多平台兼容**：支持微信公众号、知乎、掘金等平台

## 🚀 快速开始

### 环境要求
- Node.js >= 18
- Go >= 1.21（后端运行需要）

### 安装依赖

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

### 构建项目

```bash
# 构建前端
npm run build

# 构建 Wails 应用（桌面端）
wails build
```

## 📁 项目结构

```
Content-Alchemist/
├── frontend/                 # 前端项目
│   ├── src/
│   │   ├── components/      # Vue 组件
│   │   │   ├── App.vue      # 主应用组件
│   │   │   ├── FileToolbar.vue
│   │   │   ├── WritingSidebar.vue
│   │   │   └── ...
│   │   ├── composables/     # 组合式函数
│   │   ├── styles/          # 全局样式
│   │   └── types/           # TypeScript 类型定义
│   ├── package.json
│   └── vite.config.ts
├── internal/                 # 后端代码
├── main.go                   # 程序入口
└── README.md
```

## 🛠️ 技术栈

### 前端
- **Vue 3** - 渐进式 JavaScript 框架
- **TypeScript** - 类型安全的 JavaScript
- **Vite** - 下一代前端构建工具
- **Markdown-It** - Markdown 渲染引擎

### 后端（Wails）
- **Go** - 高性能后端语言
- **Wails** - 使用 Go + Web 技术构建桌面应用

## 📝 使用指南

### 1. 创建文章
在左侧边栏输入文章标题，点击"生成大纲"获取 AI 生成的大纲和候选标题。

### 2. 编辑内容
在中间编辑区编写或修改文章内容，支持标准 Markdown 语法。

### 3. 预览效果
- 点击"👁️ 预览"切换普通预览模式
- 点击"📱 手机端"查看微信公众号排版效果

### 4. 发布文章
点击"📋 复制到公众号"按钮，将美化后的内容复制到剪贴板，直接粘贴到微信公众号编辑器即可。

## 🎨 Markdown 支持

| 元素 | 语法 | 效果 |
|------|------|------|
| 标题 | `# 标题` | ✅ 微信公众号主题样式 |
| 加粗 | `**文字**` | ✅ 绿色加粗 |
| 斜体 | `*文字*` | ✅ 灰色斜体 |
| 代码块 | `
代码
` | ✅ 深色主题 + macOS 按钮 |
| 行内代码 | `` `代码` `` | ✅ 粉色边框样式 |
| 引用 | `> 引用` | ✅ 渐变背景 + 引号装饰 |
| 列表 | `- 项目` | ✅ 自定义彩色圆点 |
| 表格 | `\| 表头 \|` | ✅ 绿色表头 + 阴影 |
| 分隔线 | `---` | ✅ 渐变 + 菱形装饰 |
| 图片 | `![描述](url)` | ✅ 圆角 + 阴影 |

## ⚙️ 配置说明

### AI 配置
按 `Ctrl + ,` 打开设置面板，配置：
- API Base URL
- API Token
- 模型选择（默认 deepseek-chat）
- 温度参数

### 公众号定位
在设置面板中配置公众号定位，AI 将根据定位生成符合风格的内容。

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

[MIT](LICENSE) © Content Alchemist

## 🙏 致谢

- [Vue.js](https://vuejs.org/)
- [Vite](https://vitejs.dev/)
- [Wails](https://wails.io/)
- [Markdown-It](https://github.com/markdown-it/markdown-it)

---

**Content Alchemist** - 让内容创作更高效、更专业 🚀
