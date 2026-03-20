
# 🎯 Content-Alchemist Pro 开发任务书
## 基于 Eino ADK + Wails + Vue3 的AI写作助手

---

## 📋 项目概述

**产品定位**：面向一人公司创作者的本地优先AI写作助手，核心workflow是"输入标题→AI协作创作→公众号排版→一键发布"

**技术栈**：
- **前端**：Vue 3 + TypeScript + Vite + Pinia + Element Plus
- **后端**：Go 1.21+ + Wails v2 + Eino ADK (字节跳动AI应用开发框架)
- **存储**：SQLite (wails自带) + 本地文件系统
- **AI引擎**：DeepSeek/OpenAI API (通过Eino封装)

**架构原则**：本地优先(Local-first)、流式交互、渐进式披露

---

## 🏗️ 整体架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                        Wails App                             │
├──────────────────────┬──────────────────────────────────────┤
│    Frontend (Vue3)   │        Backend (Go)                  │
│  ┌────────────────┐  │   ┌──────────────────────────────┐   │
│  │  Editor Core   │  │   │   Eino ADK Agents            │   │
│  │  ├─ Monaco     │  │   │   ├─ OutlineAgent            │   │
│  │  ├─ Slash Cmd  │  │   │   ├─ WritingAgent            │   │
│  │  └─ Outline Nav│  │   │   ├─ PolishAgent             │   │
│  ├────────────────┤  │   │   └─ MaterialAgent           │   │
│  │  Preview Panel │  │   └──────────────────────────────┘   │
│  │  ├─ Desktop    │  │              │                       │
│  │  └─ Mobile     │  │   ┌──────────────────────────────┐   │
│  ├────────────────┤  │   │   Service Layer              │   │
│  │  Material Lib  │◄─┼───┤   ├─ ConfigService           │   │
│  │  └─ Local DB   │◄─┼───┤   ├─ ArticleService          │   │
│  ├────────────────┤  │   │   ├─ MaterialService         │   │
│  │  Onboarding    │◄─┼───┤   └─ HistoryService          │   │
│  └────────────────┘  │   └──────────────────────────────┘   │
└──────────────────────┴──────────────────────────────────────┘
```

---

## 📦 任务清单（10个垂直切片）

---

### **Task 1: 项目骨架与Eino集成**
**优先级**：P0 | **预估工时**：2h

#### 目标
搭建Wails项目基础结构，集成Eino ADK，完成基础路由和状态管理设计。

#### 详细需求

**1.1 项目初始化**
```bash
# 创建wails项目
wails init -n content-alchemist -t vue-ts
cd content-alchemist
```

**1.2 目录结构设计**
```
content-alchemist/
├── frontend/                    # Vue前端
│   ├── src/
│   │   ├── components/          # 组件
│   │   │   ├── editor/          # 编辑器相关
│   │   │   ├── preview/         # 预览相关
│   │   │   ├── ai/              # AI交互组件
│   │   │   └── common/          # 通用组件
│   │   ├── composables/         # 组合式函数
│   │   ├── stores/              # Pinia状态
│   │   ├── types/               # TS类型定义
│   │   └── utils/               # 工具函数
│   └── package.json
├── internal/                    # Go后端
│   ├── agent/                   # Eino Agents
│   ├── service/                 # 业务逻辑
│   ├── repository/              # 数据访问
│   └── model/                   # 数据模型
└── main.go
```

**1.3 Eino ADK集成（Backend）**
- 安装Eino: `go get github.com/cloudwego/eino`
- 创建 `internal/agent/base.go`：
  - 封装LLM客户端（支持DeepSeek/OpenAI）
  - 实现流式输出封装（SSE）
  - 统一Prompt模板管理

**1.4 基础类型定义**
```go
// internal/model/article.go
type Article struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Outline   []OutlineNode `json:"outline"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

type OutlineNode struct {
    ID       string `json:"id"`
    Level    int    `json:"level"`    // 1,2,3
    Title    string `json:"title"`
    Content  string `json:"content"`  // 该节点对应正文
    ParentID string `json:"parentId"`
}
```

**1.5 前端基础配置**
- 安装依赖：`pinia`, `vue-router@4`, `@element-plus/icons-vue`
- 创建基础布局组件：`MainLayout.vue`（三栏布局容器）
- 配置路由：/editor, /settings, /welcome

#### 验收标准
- [ ] `wails dev` 能正常启动，显示基础界面
- [ ] Eino能正常调用DeepSeek API并返回流式数据
- [ ] 项目结构符合上述规范

---

### **Task 2: Onboarding引导系统**
**优先级**：P0 | **预估工时**：1.5h

#### 目标
设计首次用户体验，将配置门槛转化为引导流程，确保用户5分钟内完成初始化。

#### 详细需求

**2.1 首次启动检测**
- 在 `App.vue` 中检查本地配置是否存在（`~/.content-alchemist/config.json`）
- 如无配置，自动跳转到 `/welcome` 路由
- 配置项：AI API Key、模型选择、公众号定位（风格标签）

**2.2 引导页面设计（3步法）**

**Step 1: 欢迎页**
- 文案："让公众号创作从3小时缩短到30分钟"
- 按钮："开始配置"（主按钮）、"使用本地模式"（文字按钮，跳过AI配置）

**Step 2: AI配置**
- 表单字段：
  - API Base URL（默认：https://api.deepseek.com）
  - API Key（密码输入框）
  - 模型选择（下拉：deepseek-chat/deepseek-reasoner/gpt-4o）
  - 温度参数（滑块：0-1，默认0.7）
- 实时验证：点击"测试连接"，调用Eino发送"你好"验证连通性

**Step 3: 公众号定位**
- 风格标签选择（多选）：技术干货 | 职场成长 | 产品思维 | 生活随笔 | 行业观察
- 受众描述（textarea）："面向35岁+的Golang后端开发者..."
- 写作人称（单选）：我/我们/小编/笔者

**2.3 数据存储**
- Go端创建 `ConfigService`：
  - `SaveConfig(config Config) error` - 保存到SQLite
  - `GetConfig() (Config, error)` - 读取配置
- 前端Pinia store：`stores/config.ts`，初始化时从后端加载

#### UI/UX细节
```vue
<!-- WelcomeView.vue 伪代码 -->
<template>
  <div class="welcome-container">
    <el-steps :active="currentStep" finish-status="success">
      <el-step title="欢迎" />
      <el-step title="AI配置" />
      <el-step title="定位设置" />
    </el-steps>
    
    <div v-show="currentStep === 0">...</div>
    <div v-show="currentStep === 1">
      <el-form>
        <el-form-item label="API Key">
          <el-input v-model="config.apiKey" show-password />
          <el-button @click="testConnection">测试连接</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>
```

#### 验收标准
- [ ] 首次打开应用自动进入引导页
- [ ] 配置验证失败时显示具体错误（如API Key无效）
- [ ] 配置保存后跳转到主编辑器
- [ ] 支持"本地模式"（跳过AI配置，后续可再开启）

---

### **Task 3: 核心编辑器骨架（Markdown+大纲联动）**
**优先级**：P0 | **预估工时**：3h

#### 目标
构建三栏布局的编辑器骨架，实现大纲导航与编辑区的联动，支持基本的Markdown编辑。

#### 详细需求

**3.1 布局架构（MainLayout.vue）**
```
┌─────────────────────────────────────────────────────────┐
│  Header (工具栏)                                         │
├──────────┬───────────────────────────────┬──────────────┤
│ 大纲面板  │      编辑器 (中间)             │  预览面板    │
│ (15%)    │      (50%)                    │  (35%)       │
│          │                               │              │
│ 树形结构  │   Monaco Editor /             │  手机模拟器   │
│ 点击跳转  │   ContentEditable             │  或          │
│          │                               │  桌面预览    │
└──────────┴───────────────────────────────┴──────────────┘
```

**3.2 大纲导航系统（OutlinePanel.vue）**
- 实时解析编辑器内容中的H1/H2/H3标题
- 树形展示，支持点击跳转对应位置
- 当前编辑区域高亮对应大纲节点（滚动监听）
- 支持大纲拖拽排序（调整后会同步调整正文顺序）

**3.3 编辑器选型与封装**
推荐使用 `milkdown` 或 `vditor`（Vue3友好的Markdown编辑器），或轻量级的 `monaco-editor` 配合Markdown语法高亮。

**3.4 状态管理设计（Pinia）**
```typescript
// stores/editor.ts
export const useEditorStore = defineStore('editor', {
  state: () => ({
    article: {
      id: '',
      title: '',
      content: '# 请输入标题\n\n开始创作...',
      outline: []
    } as Article,
    currentOutlineId: '', // 当前聚焦的大纲节点
    editorView: 'edit',   // 'edit' | 'preview' | 'split'
    isAIGenerating: false,
    selectedText: ''      // 当前选中文本（用于AI操作）
  }),
  
  actions: {
    updateContent(content: string) {
      this.article.content = content
      this.parseOutline() // 自动解析大纲
    },
    async generateOutline(title: string) {
      // 调用后端Agent
    }
  }
})
```

**3.5 工具栏设计**
- 左侧：标题输入框（大字号，无边框）、保存状态指示（自动保存）
- 中间：编辑器模式切换（编辑/预览/分屏）
- 右侧：AI助手按钮（spark图标）、设置按钮、发布按钮

#### 技术细节
- **大纲解析算法**：使用正则或Markdown AST解析器提取标题
- **双向绑定**：编辑器滚动时高亮大纲，点击大纲时编辑器平滑滚动到对应位置
- **自动保存**：内容变化后3秒自动保存到SQLite

#### 验收标准
- [ ] 三栏布局响应式，支持拖拽调整宽度
- [ ] 输入Markdown后大纲实时更新
- [ ] 点击大纲项，编辑器滚动到对应位置
- [ ] 自动保存功能正常（断网/崩溃后可恢复）

---

### **Task 4: 行内AI助手系统（核心创新）**
**优先级**：P0 | **预估工时**：4h

#### 目标
实现"斜杠命令"和"选中文本AI操作"，让AI从"生成整篇文章"转变为"写作伙伴"。

#### 详细需求

**4.1 斜杠命令系统（Slash Commands）**
在编辑器输入 `/` 时弹出命令面板：

| 命令 | 功能 | 快捷键 |
|------|------|--------|
| `/ai` | 召唤AI助手（通用） | Cmd+K |
| `/outline` | 基于当前标题生成大纲 | - |
| `/expand` | 扩写当前段落 | - |
| `/polish` | 润色当前段落 | - |
| `/shorter` | 精简内容 | - |
| `/title` | 生成5个标题建议 | - |

**4.2 选中文本AI菜单**
鼠标选中文本后，浮动工具条显示：
- ✨ 润色（更专业）
- 🗣️ 口语化（更自然）
- 📊 加数据（补充案例/数据）
- 🔄 换个角度（重写）
- ➕ 续写

**4.3 Eino Agent设计（Backend）**

**WritingAgent** - 负责文本生成与润色
```go
// internal/agent/writing_agent.go
type WritingAgent struct {
    llm *openai.ChatModel // Eino封装的模型
}

type WritingRequest struct {
    Action      string   // "expand", "polish", "shorten", "continue"
    Context     string   // 上下文（前几段）
    SelectedText string  // 选中的文本
    Position    string   // "before", "after", "replace"
    Style       string   // 从配置读取的风格描述
}

type WritingResponse struct {
    Content     string  // AI生成的内容
    Suggestions []string // 备选建议（如果是标题生成）
}
```

**流式输出实现**：
- 前端通过Wails Events接收流式数据：`EventsOn("ai-stream", handler)`
- 后端使用Eino的流式ChatModel，逐字发送给前端

**4.4 前端交互实现**

**SlashCommand组件**：
```vue
<template>
  <div v-if="show" class="slash-menu" :style="position">
    <div 
      v-for="cmd in commands" 
      :key="cmd.id"
      @click="execute(cmd)"
      :class="{ active: selectedIndex === cmd.id }"
    >
      <el-icon><component :is="cmd.icon" /></el-icon>
      <span>{{ cmd.name }}</span>
      <span class="desc">{{ cmd.desc }}</span>
    </div>
  </div>
</template>
```

**AI浮动工具条（FloatingAIActions.vue）**：
- 监听文本选择事件（mouseup）
- 计算选中位置，显示浮动按钮组
- 点击后弹出Inline对话框（输入额外指令，如"用鲁迅风格改写"）

**4.5 内容替换逻辑**
- 生成内容后，显示"diff对比"（删除线+新增高亮）
- 用户可选择：✅ 接受 | ❌ 拒绝 | 🔄 重试 | ✏️ 修改指令

#### 验收标准
- [ ] 输入 `/` 弹出命令菜单，支持键盘上下选择
- [ ] 选中文本后显示AI浮动工具条
- [ ] AI生成内容时显示打字机效果（流式输出）
- [ ] 支持"撤销"AI操作（Cmd+Z）
- [ ] 生成后显示接受/拒绝选项

---

### **Task 5: 智能大纲生成与文章组装**
**优先级**：P1 | **预估工时**：3h

#### 目标
实现"标题→大纲→正文"的workflow，大纲可编辑并与正文双向绑定。

#### 详细需求

**5.1 大纲生成Agent（OutlineAgent）**
```go
type OutlineAgent struct {
    llm model.ChatModel
}

func (a *OutlineAgent) GenerateOutline(ctx context.Context, title string, style string) (*Outline, error) {
    prompt := fmt.Sprintf(`请为公众号文章"%s"生成详细大纲。
要求：
1. 受众：%s
2. 风格：%s
3. 大纲结构：包含引言、3-5个核心论点（每个论点配案例）、总结
4. 每个节点注明建议字数

输出JSON格式：
{
  "nodes": [
    {"id": "1", "level": 1, "title": "引言", "suggestedWords": 200},
    {"id": "2", "level": 1, "title": "核心观点一", "suggestedWords": 500},
    {"id": "2.1", "level": 2, "title": "案例：xxx", "suggestedWords": 200}
  ]
}`, title, audience, style)
    // 使用Eino调用LLM并解析JSON
}
```

**5.2 大纲编辑器（OutlineEditor.vue）**
- 树形结构展示大纲（Element Plus的Tree组件或自定义）
- 支持操作：
  - 拖拽调整顺序
  - 双击编辑标题
  - 右键菜单：生成此节正文 | 删除 | 添加子节点
  - 每个节点显示建议字数和实际字数对比（进度条）

**5.3 大纲→正文映射**
- 每个大纲节点对应正文中的一个H2/H3标题
- 数据结构：
```typescript
interface OutlineNode {
  id: string;
  title: string;
  level: number;
  content?: string;  // 该节点下的正文内容
  status: 'empty' | 'draft' | 'done';
  wordCount: number;
  targetWords: number;
}
```

**5.4 分段生成正文**
- 点击大纲节点的"生成正文"按钮：
  - 发送该节点标题+上下文（前后节点）给WritingAgent
  - 生成内容插入到对应位置
  - 节点状态变为'done'

**5.5 快速操作栏**
在大纲面板顶部提供：
- 🎲 重新生成大纲（基于当前标题）
- 📋 应用模板（选择预设大纲模板：如"问题解决型"、"观点论证型"）
- 📊 查看全文结构（饼图展示各节点字数占比）

#### 验收标准
- [ ] 输入标题点击"生成大纲"，10秒内返回结构化大纲
- [ ] 大纲支持拖拽排序，排序后正文同步调整
- [ ] 点击"生成本节"，仅生成该H2下的内容，不影响其他部分
- [ ] 大纲节点显示字数进度（实际/建议）

---

### **Task 6: 手机预览与主题系统**
**优先级**：P1 | **预估工时**：2.5h

#### 目标
实现微信公众号风格的手机预览，支持多主题切换和实时同步。

#### 详细需求

**6.1 手机模拟器组件（MobilePreview.vue）**
- 视觉：iPhone 14 Pro外框（CSS绘制），包含刘海屏、Home Indicator
- 内部：iframe或div渲染公众号文章内容
- 顶部状态栏：时间、信号、电量（装饰性）
- 底部：公众号菜单栏模拟（增强真实感）

**6.2 主题系统**
定义CSS变量方案，支持一键切换：
```css
/* 默认主题（科技蓝） */
:root {
  --primary-color: #1677ff;
  --heading-font: 'PingFang SC', sans-serif;
  --body-font: 'PingFang SC', sans-serif;
  --code-bg: #f6f8fa;
  --quote-bg: linear-gradient(120deg, #e0c3fc 0%, #8ec5fc 100%);
  --h1-style: center double-border;
}

/* 简约白 */
[data-theme="minimal"] {
  --primary-color: #333;
  --h1-style: left border-bottom;
  --quote-bg: #f5f5f5;
}

/* 活力橙 */
[data-theme="vibrant"] {
  --primary-color: #ff6b35;
  --quote-bg: linear-gradient(120deg, #fccb90 0%, #d57eeb 100%);
}
```

**6.3 Markdown→公众号HTML转换器**
创建 `utils/wechatRenderer.ts`：
- H1: 居中 + 双彩边框 + 首字下沉
- H2: 绿色左边框 + 背景色
- H3: 蓝色左边框
- 代码块：macOS风格窗口按钮（红黄绿）+ 深色主题
- 引用块：渐变背景 + 引号装饰
- 表格：绿色表头 + 斑马纹

**6.4 实时同步**
- 编辑器使用防抖（300ms），停止输入后触发渲染
- 预览区平滑滚动（如果用户正在查看预览，保持滚动位置比例）
- 支持"跟随模式"：编辑器滚动时预览区同步滚动（根据标题锚点）

**6.5 预览模式切换**
顶部工具栏提供：
- 💻 桌面预览（普通Markdown渲染，适合技术文档）
- 📱 手机预览（默认）
- 👁️ 纯编辑模式（隐藏预览）

#### 验收标准
- [ ] 手机预览视觉效果接近真实公众号
- [ ] 切换主题即时生效，无需刷新
- [ ] 编辑器输入后300ms内预览更新
- [ ] 支持复制带样式的HTML（用于粘贴到公众号后台）

---

### **Task 7: 素材库与本地知识管理**
**优先级**：P1 | **预估工时**：2h

#### 目标
构建创作者的个人素材库，支持金句、数据、历史文章的复用。

#### 详细需求

**7.1 素材库侧边栏（MaterialPanel.vue）**
位于大纲面板下方（可折叠），包含标签页：
- 📝 **片段**：常用开头/结尾、过渡句
- 📊 **数据**：常用统计数据、报告链接
- 💡 **金句**：收藏的名人名言、金句
- 📄 **历史**：最近发布的文章，可快速引用

**7.2 数据模型**
```go
type Material struct {
    ID        string   `json:"id"`
    Type      string   `json:"type"` // snippet, data, quote, history
    Title     string   `json:"title"`
    Content   string   `json:"content"`
    Tags      []string `json:"tags"`
    Source    string   `json:"source"` // 来源（可选）
    CreatedAt time.Time `json:"createdAt"`
    UsageCount int     `json:"usageCount"` // 使用次数统计
}
```

**7.3 快速插入功能**
- 素材库中的条目支持拖拽到编辑器
- 或点击"插入"按钮，插入到当前光标位置
- 插入格式：`{{素材标题}}` 或直接替换为内容

**7.4 AI素材推荐（智能）**
当AI检测到用户输入特定关键词时（如"根据数据显示"），自动推荐相关数据素材。

**7.5 素材收集**
- 选中文本后，右键菜单"添加到素材库"
- 弹出对话框选择分类和标签

#### 验收标准
- [ ] 素材库支持CRUD操作
- [ ] 支持拖拽插入到编辑器
- [ ] 素材支持标签筛选和搜索
- [ ] 常用素材排序靠前（基于UsageCount）

---

### **Task 8: 发布助手与版本管理**
**优先级**：P1 | **预估工时**：2h

#### 目标
打通"创作→发布"的最后一公里，支持历史版本回溯。

#### 详细需求

**8.1 发布助手弹窗（PublishDialog.vue）**
点击"发布"按钮后弹出：
- **内容检查清单**：
  - [ ] 标题字数检查（建议≤30字）
  - [ ] 封面图提示（建议尺寸900×383）
  - [ ] 敏感词检测（本地基础词库）
  - [ ] 字数统计（建议800-3000字）
- **操作按钮**：
  - 📋 复制到公众号（带样式，使用clipboard API）
  - 🔗 打开公众号后台（打开浏览器 mp.weixin.qq.com）
  - ⏰ 定时提醒（设置发布时间，本地通知）

**8.2 版本管理（VersionHistory.vue）**
- 自动保存：每5分钟创建一个版本（保留最近30个）
- 手动保存：Cmd+S创建命名版本（如"完成大纲"）
- 版本对比：左右分屏diff视图（使用diff-match-patch库）
- 回滚：一键恢复到某个版本（当前内容自动保存到临时版本）

**8.3 后端实现**
```go
type Version struct {
    ID        string    `json:"id"`
    ArticleID string    `json:"articleId"`
    Content   string    `json:"content"`
    Snapshot  string    `json:"snapshot"` // 前100字摘要
    CreatedAt time.Time `json:"createdAt"`
    Type      string    `json:"type"` // auto, manual
}

func (s *HistoryService) SaveVersion(v Version) error
func (s *HistoryService) GetVersions(articleID string) ([]Version, error)
func (s *HistoryService) RestoreVersion(versionID string) (string, error) // 返回内容
```

#### 验收标准
- [ ] 发布前自动检查并给出优化建议
- [ ] 复制到剪贴板保留公众号样式（可通过微信测试）
- [ ] 支持查看历史版本并对比差异
- [ ] 崩溃后可从自动保存版本恢复

---

### **Task 9: 快捷键与命令面板**
**优先级**：P2 | **预估工时**：1.5h

#### 目标
提升效率用户的工作流，支持键盘操作和全局搜索。

#### 详细需求

**9.1 命令面板（CommandPalette.vue）**
按 `Cmd/Ctrl + K` 唤起，支持：
- 搜索功能："生成大纲"、"切换主题"、"打开设置"
- 最近文件：快速打开最近编辑的文章
- AI命令：直接输入指令如"写一个关于Go的引言"

**9.2 快捷键映射**
| 快捷键 | 功能 |
|--------|------|
| Cmd+N | 新建文章 |
| Cmd+S | 手动保存/创建版本 |
| Cmd+Shift+P | 发布文章 |
| Cmd+B | 加粗选中文本 |
| Cmd+/ | 注释/取消注释 |
| Cmd+K | 打开命令面板 |
| Esc | 关闭弹窗/退出全屏 |

**9.3 快捷键提示**
在菜单项和按钮hover时显示对应快捷键（如"保存 (Cmd+S)"）

#### 验收标准
- [ ] Cmd+K能唤起命令面板并搜索
- [ ] 所有快捷键在Windows/Mac上正常工作
- [ ] 快捷键支持自定义（可选，高级功能）

---

### **Task 10: 性能优化与构建配置**
**优先级**：P2 | **预估工时**：2h

#### 目标
确保应用流畅运行（>60fps），并完成跨平台构建配置。

#### 详细需求

**10.1 性能优化**
- **虚拟滚动**：大纲和素材库使用虚拟列表（如果条目>100）
- **防抖节流**：
  - 编辑器输入：防抖300ms保存
  - AI流式输出：节流50ms渲染
  - 预览渲染：防抖300ms
- **懒加载**：预览组件和图片懒加载

**10.2 错误处理**
- AI调用失败时显示友好提示（带重试按钮）
- 网络断开时自动转为本地模式（禁用AI功能，提示用户）
- 全局错误边界（ErrorBoundary）

**10.3 构建配置（wails.json）**
```json
{
  "name": "Content-Alchemist",
  "outputfilename": "Content-Alchemist",
  "frontend": {
    "dir": "./frontend",
    "install": "npm install",
    "build": "npm run build",
    "package": "npm run package"
  },
  "WailsVersion": "2.x",
  "build": {
    "productName": "Content-Alchemist",
    "appId": "com.yourdomain.content-alchemist",
    "asar": true,
    "directories": {
      "output": "build"
    },
    "mac": {
      "target": ["dmg", "zip"]
    },
    "win": {
      "target": ["nsis", "portable"]
    },
    "linux": {
      "target": ["AppImage", "deb"]
    }
  }
}
```

**10.4 自动更新（可选）**
集成wails的自动更新功能（使用GitHub Releases）

#### 验收标准
- [ ] 大文档（1万字）编辑不卡顿
- [ ] AI生成时界面保持响应
- [ ] 能成功构建Windows/Mac/Linux安装包
- [ ] 安装包体积<100MB（优化前端代码分割）

---

## 🎨 设计规范（UI/UX）

### 色彩系统
- **主色**：#1677ff（科技蓝）
- **成功**：#52c41a
- **警告**：#faad14
- **错误**：#f5222d
- **文本**：rgba(0,0,0,0.85)（主文本）、rgba(0,0,0,0.45)（次要文本）
- **背景**：#f0f2f5（页面背景）、#ffffff（卡片背景）

### 字体
- **界面字体**：-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto
- **编辑器字体**："JetBrains Mono", "Fira Code", monospace
- **预览字体**："PingFang SC", "Microsoft YaHei", sans-serif

### 间距
- 基础单位：8px
- 编辑器内边距：24px
- 卡片圆角：8px
- 按钮圆角：4px

---

## 📝 AI编程助手执行建议

### 每个任务的输入输出规范

**输入格式**：
```
任务ID: Task X
任务名称: XXX
技术栈: Go Eino + Vue3 + Wails
详细需求: [上述对应任务的内容]
相关接口: [如果有前置依赖]
```

**输出要求**：
1. 完整可运行的代码（不要省略关键部分）
2. 必要的注释说明设计决策
3. 简单测试用例或手动测试步骤
4. 如果涉及前后端联调，提供API契约（Request/Response格式）

### 依赖顺序
```
Task 1 → Task 2 → Task 3 → Task 4 → Task 5 → Task 6 → Task 7 → Task 8 → Task 9 → Task 10
```

**并行可能性**：
- Task 2 (Onboarding) 和 Task 6 (主题系统) 可并行
- Task 7 (素材库) 和 Task 8 (版本管理) 可并行

### 代码提交规范
每个Task完成后提交Git，提交信息格式：
```
feat(task-X): 简短描述

- 详细变更点1
- 详细变更点2
```

---

## 🔧 Eino ADK 参考示例

```go
// internal/agent/base.go
package agent

import (
    "context"
    "github.com/cloudwego/eino/components/model"
    "github.com/cloudwego/eino/flow/agent"
    "github.com/cloudwego/eino/schema"
)

type BaseAgent struct {
    model model.ChatModel
}

func NewBaseAgent(apiKey string) (*BaseAgent, error) {
    // 初始化DeepSeek/OpenAI模型
    llm, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
        APIKey: apiKey,
        Model:  "deepseek-chat",
    })
    if err != nil {
        return nil, err
    }
    return &BaseAgent{model: llm}, nil
}

func (a *BaseAgent) StreamChat(ctx context.Context, prompt string, handler func(chunk string)) error {
    stream, err := a.model.Stream(ctx, []*schema.Message{
        {Role: schema.User, Content: prompt},
    })
    if err != nil {
        return err
    }
    
    for chunk := range stream {
        if chunk.Err != nil {
            return chunk.Err
        }
        handler(chunk.Content)
    }
    return nil
}
```