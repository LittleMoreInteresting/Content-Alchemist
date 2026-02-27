# Content Alchemist - 文件管理模块实现文档

## 已实现的代码结构

### Go 后端代码

```
backend/
├── models/
│   └── models.go          # 数据结构和类型定义
├── db/
│   └── db.go              # SQLite 数据库操作
├── editor/
│   └── file.go            # 文件操作工具函数
├── app.go                 # Wails App 结构体和绑定方法
├── schema.sql             # 数据库 Schema
└── main.go                # 应用入口
```

### Vue 前端代码

```
frontend/
├── src/
│   ├── types/
│   │   ├── article.ts     # 文章相关类型
│   │   ├── wails.d.ts     # Wails 运行时类型
│   │   └── index.ts       # 类型统一导出
│   ├── composables/
│   │   ├── useWails.ts    # Wails 后端调用封装
│   │   └── useEditor.ts   # 编辑器状态管理
│   ├── components/
│   │   ├── App.vue        # 主应用组件
│   │   ├── FileToolbar.vue    # 文件工具栏
│   │   └── RecentArticles.vue # 最近文章列表
│   ├── styles/
│   │   └── global.css     # 全局样式
│   └── main.ts            # 前端入口
├── index.html
├── package.json
├── tsconfig.json
└── vite.config.ts
```

## 核心功能

### Go 后端绑定方法

| 方法 | 说明 | 约束符合 |
|------|------|----------|
| `OpenFileDialog()` | 打开文件选择对话框 | ✅ Wails 绑定 |
| `SaveFileDialog(defaultFilename)` | 打开保存对话框 | ✅ Wails 绑定 |
| `OpenDirectoryDialog()` | 打开目录选择对话框 | ✅ Wails 绑定 |
| `ReadArticle(filePath)` | 读取文章内容和元数据 | ✅ 绝对路径存储 |
| `SaveArticle(uuid, content)` | 保存文章 | ✅ 错误返回前端 |
| `SaveArticleAs(uuid, newPath, content)` | 另存为 | ✅ 绝对路径存储 |
| `CreateNewArticle()` | 创建新文章 | ✅ Wails 绑定 |
| `GetRecentArticles(limit)` | 获取最近文章 | ✅ SQLite 存储 |
| `UpdateArticleMeta(uuid, tags)` | 更新元数据 | ✅ |

### 前端 Composables

#### `useWails()`
- 类型安全的后端方法调用
- 统一的错误处理
- 加载状态管理
- 返回详细错误信息给调用方

#### `useEditor(saveDebounce?)`
- 编辑器状态管理
- 自动保存（防抖，默认2秒）
- 字数统计
- 保存错误重试（最多3次）

### 数据库 Schema

#### articles 表
```sql
CREATE TABLE articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT UNIQUE NOT NULL,
    file_path TEXT UNIQUE NOT NULL,  -- 绝对路径
    title TEXT,
    summary TEXT,
    tags TEXT DEFAULT '[]',           -- JSON 数组
    word_count INTEGER DEFAULT 0,
    created_at INTEGER,               -- Unix timestamp
    updated_at INTEGER,
    last_opened_at INTEGER
);
```

#### settings 表
```sql
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
```

## 技术约束检查

| 约束 | 实现状态 |
|------|----------|
| Wails v2 绑定机制（`//go:build` 注释） | ✅ 所有 Go 文件都有 `//go:build darwin \|\| windows \|\| linux` |
| 绝对路径存储在 SQLite | ✅ `file_path` 字段存储绝对路径，使用 `filepath.Abs()` 和 `filepath.Clean()` 处理 |
| Vue3 Composition API + TypeScript | ✅ 所有组件使用 `<script setup lang="ts">` |
| 禁止使用 Any 类型 | ✅ 所有函数和变量都有明确类型定义 |
| 错误处理返回给前端 | ✅ 所有错误都封装后返回，不在 Go 端吞掉 |
| 本地优先原则 | ✅ 所有功能离线可用，AI 调用除外（设计阶段） |

## 使用示例

### 打开文件
```typescript
const wails = useWails();

// 打开文件对话框并读取
const filePath = await wails.openFileDialog();
if (filePath) {
  const { article, content } = await wails.readArticle(filePath);
  // 处理文章数据
}
```

### 创建新文章
```typescript
const article = await wails.createNewArticle();
if (article) {
  // 新文章已创建，可以开始编辑
}
```

### 使用编辑器状态管理
```typescript
const editor = useEditor(2000); // 2秒防抖

// 设置内容（自动触发保存）
editor.setContent('# 新文章\n\n开始写作...');

// 监听状态
console.log(editor.isDirty.value);     // 是否有未保存更改
console.log(editor.wordCount.value);   // 字数
console.log(editor.isSaving.value);    // 是否正在保存
```

## 运行项目

### 1. 初始化 Go 依赖
```bash
cd Content-Alchemist
go mod tidy
```

### 2. 初始化前端依赖
```bash
cd frontend
npm install
```

### 3. 开发模式运行
```bash
# 在项目根目录
wails dev
```

### 4. 构建应用
```bash
wails build
```

## 文件说明

- **backend/app.go**: Wails 应用主入口，包含所有绑定到前端的方法
- **backend/db/db.go**: SQLite 数据库操作封装
- **backend/editor/file.go**: 文件系统操作和元数据提取
- **frontend/src/composables/useWails.ts**: 类型安全的前端调用封装
- **frontend/src/composables/useEditor.ts**: 编辑器状态管理和自动保存
- **frontend/src/components/FileToolbar.vue**: 文件操作工具栏 UI
- **frontend/src/components/RecentArticles.vue**: 最近文章侧边栏 UI
- **frontend/src/components/App.vue**: 主应用组件，整合所有功能

## 下一步开发建议

1. **集成 Milkdown 编辑器** 替换目前的 textarea
2. **实现 Markdown 实时预览** 使用 Milkdown 或 MarkdownIt + Shiki
3. **添加 AI 模块** 封装 DeepSeek API 调用
4. **实现标签管理** 添加标签输入和筛选功能
5. **添加设置界面** 配置 API Key 和主题
