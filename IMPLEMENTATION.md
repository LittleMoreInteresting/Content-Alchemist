# Content Alchemist - 实现文档

**版本**: v0.0.1
**技术栈**: Wails v2 (Go 1.23 + Vue 3 + TypeScript) + SQLite
**模式**: 本地优先 (Local-First) 桌面应用

---

## 1. 项目概述

Content Alchemist 是一个 AI 辅助的内容创作工具，帮助用户通过 AI 生成文章大纲和完整内容。所有数据存储在本地，支持 Markdown 编辑和实时预览。

### 核心功能
- 🤖 AI 辅助生成文章大纲和候选标题
- ✨ 根据大纲生成完整文章内容
- 📝 Markdown 编辑器与实时预览
- 💾 本地文件存储与元数据管理
- 📚 最近文章列表
- ⚙️ 可配置的 AI 服务（支持 DeepSeek/OpenAI 兼容 API）

---

## 2. 项目结构

```
Content-Alchemist/
├── main.go                    # 应用入口，Wails 配置
├── wails.json                 # Wails 项目配置
├── go.mod                     # Go 依赖
├── backend/                   # Go 后端
│   ├── app.go                 # 主应用结构，Wails 绑定方法
│   ├── schema.sql             # 数据库 Schema
│   ├── db/
│   │   └── db.go              # SQLite 数据库操作
│   ├── editor/
│   │   └── file.go            # 文件管理器
│   ├── ai/
│   │   └── service.go         # AI 服务封装
│   └── models/
│       └── models.go          # 数据模型定义
├── frontend/                  # Vue3 前端
│   ├── package.json           # Node.js 依赖
│   ├── vite.config.ts         # Vite 配置
│   ├── index.html
│   └── src/
│       ├── main.ts            # 应用入口
│       ├── types/
│       │   ├── article.ts     # 文章类型定义
│       │   └── index.ts       # 类型导出
│       ├── composables/
│       │   ├── useWails.ts    # Wails 后端调用封装
│       │   └── useEditor.ts   # 编辑器逻辑
│       ├── components/
│       │   ├── App.vue        # 主应用组件
│       │   ├── FileToolbar.vue       # 顶部工具栏
│       │   ├── WritingSidebar.vue    # 左侧写作助手
│       │   ├── SettingsModal.vue     # 设置弹窗
│       │   ├── TitleSelectorModal.vue # 标题选择弹窗
│       │   └── RecentArticles.vue    # 最近文章列表
│       └── styles/
│           └── global.css     # 全局样式
└── build/                     # 构建资源
    ├── appicon.png
    └── windows/
```

---

## 3. 技术架构

### 3.1 架构分层

```
┌─────────────────────────────────────────┐
│  Frontend (Vue 3 + TypeScript)          │
│  - components/App.vue                   │
│  - components/WritingSidebar.vue        │
│  - composables/useWails.ts              │
├─────────────────────────────────────────┤
│  Wails Runtime (自动生成的绑定)          │
│  - window.go.backend.App.*              │
├─────────────────────────────────────────┤
│  Backend (Go 1.23)                      │
│  - app.go (Wails App)                   │
│  - editor/file.go (文件操作)            │
│  - ai/service.go (AI API封装)           │
│  - db/db.go (SQLite元数据)              │
└─────────────────────────────────────────┘
```

### 3.2 数据存储策略

| 数据类型 | 存储方式 | 说明 |
|---------|---------|------|
| 文章内容 | `.md` 文件 | 本地文件系统，标准 Markdown 格式 |
| 文章元数据 | SQLite | 标题、摘要、字数、标签、打开时间等 |
| AI 配置 | SQLite | API URL、Token、模型、Temperature |
| 公众号定位 | SQLite | 定位描述，用于 AI 生成上下文 |

---

## 4. 数据模型

### 4.1 SQLite Schema

```sql
-- 文章元数据表
CREATE TABLE articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT UNIQUE NOT NULL,
    file_path TEXT UNIQUE NOT NULL,
    title TEXT,
    summary TEXT,
    tags TEXT DEFAULT '[]',
    word_count INTEGER DEFAULT 0,
    created_at INTEGER,
    updated_at INTEGER,
    last_opened_at INTEGER
);

-- 应用配置表
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT,
    updated_at INTEGER DEFAULT (strftime('%s', 'now'))
);
```

### 4.2 Go 数据模型 (backend/models/models.go)

```go
// Article 文章元数据
type Article struct {
    ID           int64       `json:"id"`
    UUID         string      `json:"uuid"`
    FilePath     string      `json:"filePath"`
    Title        string      `json:"title"`
    Summary      string      `json:"summary"`
    Tags         StringSlice `json:"tags"`  // JSON 序列化
    WordCount    int         `json:"wordCount"`
    CreatedAt    time.Time   `json:"createdAt"`
    UpdatedAt    time.Time   `json:"updatedAt"`
    LastOpenedAt time.Time   `json:"lastOpenedAt"`
}

// AIConfig AI 配置
type AIConfig struct {
    BaseURL     string  `json:"baseUrl"`
    Token       string  `json:"token"`
    Temperature float64 `json:"temperature"`
    Model       string  `json:"model"`
}

// GenerateOutlineResult 生成大纲结果
type GenerateOutlineResult struct {
    Titles  []string `json:"titles"`
    Outline string   `json:"outline"`
}
```

---

## 5. 后端实现

### 5.1 主应用结构 (backend/app.go)

```go
type App struct {
    ctx       context.Context
    db        *db.DB
    fileMgr   *editor.FileManager
    aiService *ai.Service
    configDir string
}
```

### 5.2 Wails 绑定方法

| 方法 | 说明 |
|-----|------|
| `SaveFileDialog(defaultFilename)` | 打开保存文件对话框 |
| `ReadArticle(filePath)` | 读取文章（自动创建元数据） |
| `SaveArticle(uuid, content)` | 保存文章 |
| `SaveArticleAs(uuid, newPath, content)` | 另存为 |
| `CreateNewArticle()` | 创建新文章 |
| `SaveArticleWithSmartNaming(...)` | 智能保存（自动生成文件名） |
| `GetRecentArticles(limit)` | 获取最近文章列表 |
| `GetAIConfig()` | 获取 AI 配置 |
| `SaveAIConfig(config)` | 保存 AI 配置 |
| `GetPositioning()` | 获取公众号定位 |
| `SavePositioning(positioning)` | 保存公众号定位 |
| `GenerateOutline(title, requirements, positioning)` | 生成大纲和候选标题 |
| `GenerateArticle(title, outline, requirements)` | 生成文章内容 |

### 5.3 AI 服务 (backend/ai/service.go)

支持 OpenAI 兼容的 API 接口：

```go
// 生成大纲和候选标题
func (s *Service) GenerateOutlineWithTitles(title, requirements, positioning string) (*OutlineResult, error)

// 根据大纲生成文章
func (s *Service) GenerateArticleFromOutline(title, outline, requirements string) (string, error)
```

**提示词设计**：
- 生成大纲时同时生成 3 个候选爆款标题
- 返回格式使用 `===TITLES===` 和 `===OUTLINE===` 标记
- 支持公众号定位作为上下文

---

## 6. 前端实现

### 6.1 主组件结构 (frontend/src/components/App.vue)

```
┌───────────────────────────────────────────┐
│  FileToolbar (顶部工具栏)                  │
├─────────────────────┬─────────────────────┤
│                     │                     │
│  WritingSidebar     │   编辑区  │  预览区  │
│  (写作助手)          │   (textarea)  (HTML)│
│                     │                     │
│  - 标题输入          ├─────────────────────┤
│  - 写作要求          │  状态栏 (字数/保存状态)│
│  - 生成大纲按钮       │                     │
│  - 生成文章按钮       │                     │
│  - 最近文章列表       │                     │
│                     │                     │
└─────────────────────┴─────────────────────┘
```

### 6.2 核心 Composables

**useWails.ts** - Wails 后端调用封装：

```typescript
export function useWails() {
  // 文件操作
  const readArticle = async (filePath: string): Promise<{article, content}>
  const saveArticle = async (uuid: string, content: string): Promise<void>
  const saveArticleWithSmartNaming = async (...): Promise<Article | null>

  // AI 操作
  const generateOutline = async (title, requirements?, positioning?): Promise<{titles, outline}>
  const generateArticle = async (title, outline, requirements?): Promise<string>

  // 配置
  const getAIConfig = async (): Promise<AIConfig>
  const saveAIConfig = async (config: AIConfig): Promise<void>
}
```

### 6.3 Markdown 渲染

使用 `markdown-it` 进行实时预览渲染：

```typescript
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true,
});

const renderedContent = computed(() => md.render(editorContent.value));
```

---

## 7. 开发指南

### 7.1 环境要求

- Go 1.23+
- Node.js 18+
- Wails CLI v2.10+

### 7.2 安装依赖

```bash
# 安装 Go 依赖
go mod tidy

# 安装前端依赖
cd frontend
npm install
```

### 7.3 开发模式

```bash
# 启动开发服务器（带热重载）
wails dev

# 仅启动前端开发服务器
cd frontend
npm run dev
```

### 7.4 构建

```bash
# 构建生产版本
wails build

# 构建 Windows 安装包
wails build -platform windows/amd64 -nsis
```

---

## 8. 关键实现细节

### 8.1 字数统计

支持中英文混合统计：
- 中文：每个汉字算 1 个字
- 英文：按空格分词

```go
func (fm *FileManager) CountWords(content string) int {
    // 移除代码块
    // 统计中文字符数 + 英文单词数
}
```

### 8.2 智能文件名生成

```go
func (fm *FileManager) SanitizeFilename(name string) string {
    // 移除非法字符
    // 限制长度 200
    // 避免重复文件名
}
```

### 8.3 外部文件修改检测

保存前检查文件是否被外部修改：

```go
modified, _, err := fm.CheckFileModified(article.FilePath, article.UpdatedAt)
if modified {
    return FileError{Code: "FILE_MODIFIED_EXTERNALLY"}
}
```

### 8.4 自动保存流程

1. 用户编辑内容触发 `watch`
2. 设置 `isDirty = true`
3. 调用 `saveArticleWithSmartNaming`
4. 如果文件未保存过，根据标题生成文件名并弹出对话框
5. 保存成功更新 `isDirty = false`

---

## 9. 配置文件

### 9.1 wails.json

```json
{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "Content-Alchemist",
  "frontend": {
    "dir": "./frontend",
    "install": "npm install",
    "build": "npm run build",
    "dev": "npm run dev"
  },
  "bindings": {
    "output": "frontend/src/generated/bindings"
  }
}
```

### 9.2 vite.config.ts

```typescript
export default defineConfig({
  plugins: [vue()],
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
});
```

---

## 10. 依赖列表

### Go 依赖

| 包 | 用途 |
|---|------|
| github.com/wailsapp/wails/v2 | Wails 框架 |
| modernc.org/sqlite | SQLite 驱动（纯 Go） |
| github.com/google/uuid | UUID 生成 |

### Node.js 依赖

| 包 | 用途 |
|---|------|
| vue | 前端框架 |
| markdown-it | Markdown 渲染 |
| vite | 构建工具 |
| typescript | 类型支持 |

---

## 11. 待办事项

- [ ] 支持图片插入和素材库
- [ ] 多平台发布（微信公众号、知乎等）
- [ ] 代码块语法高亮
- [ ] Mermaid 图表渲染
- [ ] 深色主题完整适配
- [ ] AI 改写/续写/润色功能
- [ ] 全文搜索
- [ ] 标签管理
