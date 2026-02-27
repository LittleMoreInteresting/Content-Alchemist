# Content Alchemist - Phase 1 功能设计文档
**版本**：v1.0  
**阶段**：MVP - 核心编辑器  
**技术栈**：Wails v2 (Go) + Vue 3 (TypeScript) + SQLite  
**运行模式**：纯本地优先 (Local-First)

---

## 1. 产品定义与边界

### 1.1 核心目标
构建一个**本地优先的技术写作编辑器**，支持Markdown实时渲染、代码高亮、AI辅助改写，所有数据存储在用户本地磁盘，无需联网即可编辑（AI功能除外）。

### 1.2 功能边界（Phase 1不做）
- ❌ 暂不实现多平台发布（Phase 2）
- ❌ 暂不实现素材库管理（Phase 3）
- ❌ 暂不实现用户系统/License验证
- ❌ 暂不实现云端同步

### 1.3 成功标准
- [ ] 创建/打开/保存 Markdown 文件到本地任意路径
- [ ] 编辑时实时预览（分屏模式），代码块支持语法高亮
- [ ] 支持 Mermaid 图表渲染、LaTeX 公式
- [ ] 选中文字后调用 DeepSeek API 进行改写/续写/翻译
- [ ] 文章元数据（标签、创建时间）存储在 SQLite

---

## 2. 技术架构详解

### 2.1 架构分层
```
┌─────────────────────────────────────┐
│  Frontend (Vue 3 + TypeScript)      │  ← Wails Runtime注入window.go
│  - components/Editor.vue            │
│  - components/Preview.vue           │
│  - services/ai.ts                   │
├─────────────────────────────────────┤
│  Wails Bridge (WailsJS)             │  ← 自动生成的绑定层
├─────────────────────────────────────┤
│  Backend (Go 1.21+)                 │
│  - app.go (Wails App)               │
│  - editor/ (文件操作)                │
│  - ai/ (DeepSeek封装)                │
│  - db/ (SQLite元数据)                │
│  - models/ (数据结构)                │
└─────────────────────────────────────┘
```

### 2.2 数据存储策略（Local-First）
- **文章内容**：存为原始 `.md` 文件，使用标准Markdown格式，**不加密**，方便用户用其他编辑器打开
- **元数据**：SQLite 存储文件路径、标签、字数统计、创建/修改时间、AI使用记录
- **配置**：SQLite `config` 表（主题、字体、DeepSeek API Key）

### 2.3 关键技术选型
| 模块 | 技术 | 理由 |
|------|------|------|
| 富文本编辑 | Milkdown (ProseMirror) | 原生Markdown AST支持，插件生态丰富，支持React/Vue |
| 代码高亮 | Shiki | VS Code同款引擎，支持50+语言，可加载VS主题 |
| Mermaid | mermaid.js (CDN/本地) | 技术文章画图标配 |
| LaTeX | KaTeX (CDN/本地) | 渲染快，无需完整TeX环境 |
| SQLite驱动 | mattn/go-sqlite3 | CGO驱动，Wails支持，支持并发读 |
| HTTP Client | resty/v2 | 简洁的DeepSeek API调用 |
| 配置管理 | os.UserConfigDir | 遵循XDG规范，Windows/Mac/Linux自动适配 |

---

## 3. 数据模型设计

### 3.1 SQLite Schema
```sql
-- 文章元数据表
CREATE TABLE articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT UNIQUE NOT NULL,           -- 业务ID，前端识别用
    file_path TEXT UNIQUE NOT NULL,      -- 本地绝对路径
    title TEXT,                          -- 缓存标题（H1或文件名）
    summary TEXT,                        -- 文章摘要（前100字）
    tags TEXT,                           -- JSON数组 ["Go", "架构"]
    word_count INTEGER DEFAULT 0,
    created_at INTEGER,                  -- Unix timestamp
    updated_at INTEGER,                  -- Unix timestamp
    last_opened_at INTEGER               -- 用于"最近打开"列表
);

-- 应用配置表
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT
);

-- 初始化默认配置
INSERT INTO settings (key, value) VALUES 
('deepseek_api_key', ''),
('editor_theme', 'github-dark'),
('font_size', '16'),
('ai_model', 'deepseek-chat');
```

### 3.2 Go 数据结构 (models/models.go)
```go
package models

import "time"

type Article struct {
    ID           int64     `json:"id"`
    UUID         string    `json:"uuid"`
    FilePath     string    `json:"filePath"`
    Title        string    `json:"title"`
    Summary      string    `json:"summary"`
    Tags         []string  `json:"tags"` // SQLite存储为JSON字符串，Go端自动序列化
    WordCount    int       `json:"wordCount"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
    LastOpenedAt time.Time `json:"lastOpenedAt"`
}

type EditorSettings struct {
    DeepseekAPIKey string `json:"deepseekApiKey"`
    EditorTheme    string `json:"editorTheme"`
    FontSize       int    `json:"fontSize"`
    AIModel        string `json:"aiModel"`
}

// AI请求/响应结构
type AIRequest struct {
    SelectedText string `json:"selectedText"`
    Action       string `json:"action"` // "rewrite", "polish", "continue", "shorter"
    Context      string `json:"context"` // 全文上下文（可选）
}

type AIResponse struct {
    Result  string `json:"result"`
    Error   string `json:"error,omitempty"`
    TokensUsed int `json:"tokensUsed,omitempty"`
}
```

---

## 4. 功能模块详细设计

### 4.1 文件管理模块 (editor/file.go)

**职责**：本地文件系统操作，路径管理

**暴露给前端的绑定方法**：
```go
// App struct 方法 (Wails要求)

// OpenFileDialog 打开系统文件选择器，返回选中的路径
func (a *App) OpenFileDialog() (string, error)

// ReadArticle 读取文件内容 + 元数据（如果不存在则创建DB记录）
func (a *App) ReadArticle(filePath string) (*models.Article, string, error)

// SaveArticle 保存内容到文件，更新DB元数据
func (a *App) SaveArticle(uuid string, content string) error

// CreateNewArticle 创建新文件（弹出保存对话框）
func (a *App) CreateNewArticle() (*models.Article, error)

// GetRecentArticles 获取最近打开的文章列表（按last_opened_at倒序）
func (a *App) GetRecentArticles(limit int) ([]*models.Article, error)

// UpdateArticleMeta 更新元数据（标题、标签等，不触碰文件内容）
func (a *App) UpdateArticleMeta(uuid string, tags []string) error
```

**业务逻辑**：
1. 文件保存使用 `os.WriteFile`，覆盖写入，**不自动备份**（Phase 1简化）
2. 首次打开文件时，解析H1标题（`# Title`）或提取前50字作为标题存入DB
3. 文件被外部修改检测（可选Phase 1实现或标记TODO）：通过对比`os.Stat().ModTime`与DB中的`updated_at`

### 4.2 编辑器前端组件 (frontend/src/components/)

#### EditorPane.vue（左侧编辑区）
```vue
<template>
  <div class="editor-pane">
    <MilkdownProvider>
      <MilkdownEditor 
        v-model="content"
        :theme="editorTheme"
        @change="handleChange"
      />
    </MilkdownProvider>
    <!-- 底部状态栏：字数、编码、文件路径 -->
    <StatusBar :wordCount="wordCount" :filePath="currentFilePath" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { MilkdownEditor } from '@milkdown/vue'
import { useWails } from '../composables/useWails'

const props = defineProps<{ initialContent: string, articleUUID: string }>()
const content = ref(props.initialContent)
const wordCount = ref(0)
const { SaveArticle } = useWails() // 封装window.go.main.App.SaveArticle

// 防抖保存：停止输入2秒后自动保存到本地
let saveTimer: NodeJS.Timeout
watch(content, (newVal) => {
  wordCount.value = newVal.split(/\s+/).length // 简单分词统计
  clearTimeout(saveTimer)
  saveTimer = setTimeout(() => {
    SaveArticle(props.articleUUID, newVal)
  }, 2000)
})

// AI菜单：选中文字后弹出浮动菜单
const handleTextSelect = (selectedText: string) => {
  if (selectedText.length > 5) {
    showAIMenu(selectedText)
  }
}
</script>
```

#### PreviewPane.vue（右侧预览区）
```vue
<template>
  <div class="preview-pane markdown-body" ref="previewRef">
    <!-- 使用MarkdownIt渲染，配合Shiki高亮 -->
    <div v-html="renderedHTML"></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import MarkdownIt from 'markdown-it'
import Shiki from '@shikijs/markdown-it'
import { createHighlighter } from 'shiki'

const props = defineProps<{ markdown: string }>()

// 初始化MarkdownIt（带代码高亮、Mermaid、KaTeX）
const md = MarkdownIt()
// 具体初始化逻辑在onMounted，需异步加载Shiki

const renderedHTML = computed(() => md.render(props.markdown))
</script>
```

#### AIMenu.vue（浮动AI助手）
```vue
<template>
  <div v-if="visible" class="ai-menu" :style="{ top: y + 'px', left: x + 'px' }">
    <button @click="handleAction('polish')">润色</button>
    <button @click="handleAction('shorter')">精简</button>
    <button @click="handleAction('continue')">续写</button>
    <button @click="handleAction('code')">解释代码</button>
  </div>
</template>

<script setup>
const emit = defineEmits(['apply'])
const handleAction = async (action: string) => {
  const result = await window.go.main.App.CallDeepSeek(props.selectedText, action)
  emit('apply', result)
}
</script>
```

### 4.3 AI模块 (ai/deepseek.go)

**配置**：从Settings表读取API Key，无Key时前端提示"请配置API Key"

**API封装**：
```go
package ai

import (
    "context"
    "fmt"
    "github.com/go-resty/resty/v2"
)

type DeepSeekClient struct {
    APIKey string
    Client *resty.Client
}

func NewClient(apiKey string) *DeepSeekClient {
    return &DeepSeekClient{
        APIKey: apiKey,
        Client: resty.New().SetBaseURL("https://api.deepseek.com"),
    }
}

func (c *DeepSeekClient) Rewrite(ctx context.Context, text, action string) (string, int, error) {
    prompts := map[string]string{
        "polish":   "润色以下文字，使其更通顺专业，保持原意：",
        "shorter":  "将以下文字精简到原来的70%长度，保留核心信息：",
        "continue": "基于以下上下文继续写作，保持风格一致：",
        "code":     "解释以下代码的功能，用中文说明：",
    }
    
    prompt := prompts[action] + "\n\n" + text
    
    // 调用DeepSeek Chat API
    // 使用 deepseek-chat 模型，温度0.7
    // 返回结果和token消耗
}
```

**Wails绑定方法**：
```go
func (a *App) CallDeepSeek(selectedText, action string) models.AIResponse {
    apiKey := a.getSetting("deepseek_api_key")
    if apiKey == "" {
        return models.AIResponse{Error: "API Key未配置"}
    }
    
    client := ai.NewClient(apiKey)
    result, tokens, err := client.Rewrite(context.Background(), selectedText, action)
    if err != nil {
        return models.AIResponse{Error: err.Error()}
    }
    
    return models.AIResponse{Result: result, TokensUsed: tokens}
}
```

---

## 5. 界面布局与交互

### 5.1 主窗口布局 (App.vue)
```
┌─────────────────────────────────────────────┐
│  [Logo]  新建  打开  最近文件 ▼  |  [设置]   │  ← 顶部工具栏
├──────────────┬──────────────────────────────┤
│              │                              │
│   Editor     │      Preview                 │
│   (Milkdown) │   (Markdown渲染)             │
│              │                              │
│  - 行号显示   │   - 代码高亮(Shiki)          │
│  - 当前行高亮 │   - Mermaid图表渲染          │
│  - AI触发按钮 │   - LaTeX公式                │
│              │                              │
├──────────────┴──────────────────────────────┤
│  就绪 | 字数: 1,240 | 路径: ~/docs/test.md   │  ← 底部状态栏
└─────────────────────────────────────────────┘
```

### 5.2 交互流程图

**打开文件**：
1. 用户点击"打开" → `OpenFileDialog()` → 返回路径
2. `ReadArticle(path)` → 检查DB：存在则读取元数据，不存在则插入新记录
3. 返回内容给Vue → 渲染Editor和Preview
4. 更新`last_opened_at`

**编辑保存**：
1. 用户输入 → Vue的v-model更新
2. 防抖2秒 → 调用`SaveArticle(uuid, content)`
3. Go端：`os.WriteFile`写入磁盘 + `UPDATE articles SET updated_at=?, word_count=?`
4. 前端显示"已保存"提示（Toast）

**AI改写**：
1. 用户选中文字 → 弹出浮动菜单（计算选区坐标）
2. 点击"润色" → `CallDeepSeek(selectedText, "polish")`
3. 显示Loading → 返回结果 → 弹出Diff对比框（接受/拒绝）
4. 接受则替换Editor中选中的文本

---

## 6. 开发顺序与检查点

### Week 1：基础框架
- [ ] **Day 1-2**：Wails项目初始化，配置Vue3+TypeScript，集成Milkdown基础编辑器
- [ ] **Day 3**：实现文件对话框绑定（`OpenFileDialog`, `CreateNewArticle`）
- [ ] **Day 4**：SQLite初始化，实现`ReadArticle`和`SaveArticle`基础逻辑
- [ ] **Day 5**：左右分屏布局，实现基础实时预览（MarkdownIt基础渲染）
- [ ] **Checkpoint**：能创建文件、编辑、保存，重启后通过"最近打开"列表找回

### Week 2：增强功能
- [ ] **Day 6**：集成Shiki代码高亮，支持切换主题（Light/Dark）
- [ ] **Day 7**：集成Mermaid和KaTeX，测试技术文章渲染效果
- [ ] **Day 8**：实现AI模块基础封装，配置界面（设置API Key）
- [ ] **Day 9**：实现选中文本浮动菜单，AI改写Diff对比UI
- [ ] **Day 10**：标签管理、文章元数据编辑、搜索过滤最近文件列表
- [ ] **Final Checkpoint**：完整的技术文章编辑体验，AI辅助可用，数据本地存储可靠

---

## 7. 异常处理与边界情况

| 场景 | 处理方案 |
|------|----------|
| 文件被外部删除 | 保存时检测`os.Stat`，如果不存在则提示"文件已被删除，是否另存为？" |
| 磁盘权限不足 | 捕获`os.WriteFile`错误，前端提示"无写入权限，请更换保存路径" |
| AI API超时 | 设置30秒超时，返回错误"AI服务响应超时，请检查网络" |
| 大文件性能 | 超过10MB的Markdown文件，禁用实时预览，改为手动刷新按钮 |
| 路径含中文/空格 | Go端使用`filepath.Clean`处理，确保跨平台兼容 |

---
