<template>
  <div class="app">
    <!-- 顶部工具栏 -->
    <FileToolbar
      :article="editor.article.value"
      :is-dirty="editor.isDirty.value"
      :is-saving="editor.isSaving.value"
      :can-save="editor.canSave.value"
      @new="handleNew"
      @open="handleOpen"
      @save="handleSave"
      @saveAs="handleSaveAs"
      @settings="showSettings = true"
      @error="handleError"
    />

    <!-- 主体区域 -->
    <div class="main-container">
      <!-- 左侧边栏：最近文章 -->
      <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
        <RecentArticles
          ref="recentArticlesRef"
          :current-path="editor.filePath.value"
          @select="handleOpen"
          @open="handleOpenClick"
          @error="handleError"
        />
      </aside>

      <!-- 编辑器和预览区域 -->
      <div class="editor-container">
        <!-- 标题编辑区 -->
        <div v-if="editor.article.value" class="title-editor">
          <div class="title-input-wrapper">
            <input
              type="text"
              v-model="editableTitle"
              @blur="handleTitleBlur"
              @keydown.enter="handleTitleEnter"
              placeholder="输入文章标题..."
              class="title-input"
            />
            <button
              class="generate-outline-btn"
              :disabled="isGeneratingOutline || !editableTitle"
              @click="handleGenerateOutline"
              title="根据标题生成大纲"
            >
              <span v-if="isGeneratingOutline" class="loading-spinner-small"></span>
              <span v-else>✨ 生成大纲</span>
            </button>
          </div>
        </div>

        <div class="editor-layout" :class="{ 'preview-only': previewOnly }">
          <!-- 编辑区域 -->
          <div class="editor-pane" v-show="!previewOnly">
            <div class="pane-header">
              <span class="pane-title">编辑</span>
              <div class="pane-actions">
                <button
                  class="pane-btn"
                  :class="{ active: syncScroll }"
                  @click="syncScroll = !syncScroll"
                  title="同步滚动"
                >
                  同步
                </button>
                <button class="pane-btn" @click="previewOnly = true" title="隐藏编辑">
                  →
                </button>
              </div>
            </div>
            <div class="editor-content">
              <textarea
                ref="editorRef"
                :value="editor.content.value"
                @input="e => editor.setContent((e.target as HTMLTextAreaElement).value)"
                @scroll="handleEditorScroll"
                placeholder="开始写作..."
                class="markdown-editor"
                spellcheck="false"
              ></textarea>
            </div>
          </div>

          <!-- 预览区域 -->
          <div class="preview-pane" :class="{ expanded: previewOnly }">
            <div class="pane-header">
              <div class="pane-actions">
                <button class="pane-btn" @click="previewOnly = false" title="显示编辑" v-if="previewOnly">
                  ←
                </button>
              </div>
              <span class="pane-title">预览</span>
              <div class="pane-actions">
                <button
                  class="pane-btn"
                  :class="{ active: showOutline }"
                  @click="showOutline = !showOutline"
                  title="目录"
                >
                  目录
                </button>
              </div>
            </div>
            <div class="preview-content" ref="previewRef" @scroll="handlePreviewScroll">
              <div class="markdown-body" v-html="renderedContent"></div>
            </div>
          </div>
        </div>

        <!-- 底部状态栏 -->
        <div class="status-bar">
          <div class="status-left">
            <span class="status-item" :class="{ dirty: editor.isDirty.value }">
              {{ editor.isSaving.value ? '保存中...' : editor.isDirty.value ? '未保存' : '已保存' }}
            </span>
            <span class="status-separator">|</span>
            <span class="status-item">
              字数: {{ editor.wordCount.value }}
            </span>
            <span class="status-separator">|</span>
            <span class="status-item" v-if="editor.article.value">
              {{ formatFileName(editor.article.value.filePath) }}
            </span>
          </div>
          <div class="status-right">
            <button class="status-btn" @click="sidebarCollapsed = !sidebarCollapsed">
              {{ sidebarCollapsed ? '显示侧边栏' : '隐藏侧边栏' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- AI 浮动菜单 -->
    <AIMenu
      v-if="showAIMenu"
      :selected-text="selectedText"
      :position="aiMenuPosition"
      @apply="handleAIApply"
      @close="showAIMenu = false"
    />

    <!-- 设置弹窗 -->
    <SettingsModal
      v-model:visible="showSettings"
      :initial-config="aiConfig"
      @save="handleSettingsSave"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import MarkdownIt from 'markdown-it';
import FileToolbar from './FileToolbar.vue';
import RecentArticles from './RecentArticles.vue';
import AIMenu from './AIMenu.vue';
import SettingsModal, { type AIConfig } from './SettingsModal.vue';
import { useEditor } from '../composables/useEditor';
import { useWails } from '../composables/useWails';

// ============================================
// 编辑器状态
// ============================================
const editor = useEditor(2000);
const wails = useWails();
const recentArticlesRef = ref<InstanceType<typeof RecentArticles> | null>(null);

// ============================================
// Markdown 渲染
// ============================================
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true,
});

const renderedContent = computed(() => {
  return md.render(editor.content.value);
});

// ============================================
// UI 状态
// ============================================
const sidebarCollapsed = ref(false);
const previewOnly = ref(false);
const syncScroll = ref(true);
const showOutline = ref(false);
const showSettings = ref(false);
const aiConfig = ref<AIConfig>({
  baseUrl: 'https://api.deepseek.com/v1',
  token: '',
  temperature: 0.7,
  model: 'deepseek-chat'
});
const showAIMenu = ref(false);
const selectedText = ref('');
const aiMenuPosition = ref({ x: 0, y: 0 });
const previewRef = ref<HTMLElement | null>(null);
const editorRef = ref<HTMLTextAreaElement | null>(null);
const editableTitle = ref('');
const isGeneratingOutline = ref(false);

/**
 * 聚焦到编辑器
 */
const focusEditor = (): void => {
  editorRef.value?.focus();
};

// ============================================
// 生成大纲功能
// ============================================

/**
 * 根据标题生成大纲
 */
const handleGenerateOutline = async (): Promise<void> => {
  if (!editableTitle.value) {
    handleError('请先输入文章标题');
    return;
  }

  isGeneratingOutline.value = true;

  try {
    const outline = await wails.generateOutline(editableTitle.value);

    // 将生成的大纲设置到编辑器内容中
    if (outline) {
      // 如果有标题，在第一行添加标题
      const content = `# ${editableTitle.value}\n\n${outline}`;
      editor.setContent(content);
    }
  } catch (err) {
    console.error('生成大纲失败:', err);
    handleError(err instanceof Error ? err.message : '生成大纲失败');
  } finally {
    isGeneratingOutline.value = false;
  }
};

// ============================================
// 滚动同步
// ============================================
const handleEditorScroll = (e: Event) => {
  if (!syncScroll.value || !previewRef.value) return;

  const textarea = e.target as HTMLTextAreaElement;
  const scrollRatio = textarea.scrollTop / (textarea.scrollHeight - textarea.clientHeight || 1);
  const previewScrollTop = scrollRatio * (previewRef.value.scrollHeight - previewRef.value.clientHeight);

  previewRef.value.scrollTop = previewScrollTop;
};

const handlePreviewScroll = () => {
  // 目前只支持编辑器到预览的单向同步
  // 双向同步会导致循环，需要更复杂的处理
};

// ============================================
// 文件操作
// ============================================
const handleNew = async () => {
  const success = await editor.createNew();
  if (success) {
    recentArticlesRef.value?.refresh();
    // 同步标题到编辑框
    editableTitle.value = editor.article.value?.title || '';
  }
};

const handleOpen = async (filePath: string) => {
  try {
    await editor.loadArticle(filePath);
    // 同步标题到编辑框
    editableTitle.value = editor.article.value?.title || '';
  } catch (err) {
    handleError(err instanceof Error ? err.message : '打开文件失败');
  }
};

const handleOpenClick = async () => {
  // 打开文件对话框
  try {
    const filePath = await wails.openFileDialog();
    if (filePath) {
      await handleOpen(filePath);
    }
  } catch (err) {
    handleError(err instanceof Error ? err.message : '打开文件失败');
  }
};

const handleSave = async () => {
  await editor.saveNow();
};

const handleSaveAs = async () => {
  const newArticle = await editor.saveAs();
  if (newArticle) {
    recentArticlesRef.value?.refresh();
  }
};

// ============================================
// 标题编辑
// ============================================

/**
 * 同步标题到文章
 */
const syncTitleToArticle = async (): Promise<void> => {
  if (!editor.article.value || editableTitle.value === editor.article.value.title) {
    return;
  }

  try {
    const newPath = await wails.renameArticleByTitle(
      editor.article.value.uuid,
      editableTitle.value,
      editor.content.value
    );

    if (newPath) {
      // 更新本地文章路径
      editor.article.value.filePath = newPath;
      editor.article.value.title = editableTitle.value;
    }
  } catch (err) {
    console.error('更新标题失败:', err);
    handleError(err instanceof Error ? err.message : '更新标题失败');
    // 重新加载文章以获取正确的状态
    if (editor.article.value) {
      try {
        await editor.loadArticle(editor.article.value.filePath);
        editableTitle.value = editor.article.value?.title || '';
      } catch (reloadErr) {
        console.error('重新加载文章失败:', reloadErr);
      }
    }
  }
};

/**
 * 处理标题失焦
 */
const handleTitleBlur = (): void => {
  syncTitleToArticle();
};

/**
 * 处理标题回车
 */
const handleTitleEnter = (): void => {
  syncTitleToArticle();
  focusEditor();
};

// ============================================
// AI 功能
// ============================================
const handleAIApply = (result: string) => {
  console.log('AI 结果:', result);
  showAIMenu.value = false;
};

// ============================================
// AI配置
// ============================================
const loadAIConfig = async () => {
  try {
    const config = await wails.getAIConfig();
    aiConfig.value = config;
  } catch (err) {
    console.error('加载AI配置失败:', err);
  }
};

const handleSettingsSave = async (config: AIConfig) => {
  try {
    await wails.saveAIConfig(config);
    aiConfig.value = config;
    console.log('AI配置已保存');
  } catch (err) {
    console.error('保存AI配置失败:', err);
    handleError(err instanceof Error ? err.message : '保存配置失败');
  }
};

// ============================================
// 工具函数
// ============================================
const handleError = (message: string) => {
  console.error('错误:', message);
};

const formatFileName = (filePath: string): string => {
  if (!filePath) return '未命名';
  const parts = filePath.split(/[/\\]/);
  return parts[parts.length - 1] || filePath;
};

// ============================================
// 键盘快捷键
// ============================================
const handleKeyDown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key === ',') {
    e.preventDefault();
    showSettings.value = true;
  }
};

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown);
  loadAIConfig();
});

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown);
});
</script>

<style scoped>
.app {
  display: flex;
  flex-direction: column;
  height: 100vh;
  font-family: var(--font-family-base);
  background: var(--bg-component);
  color: var(--text-primary);
}

/* 主体区域 */
.main-container {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* 侧边栏 */
.sidebar {
  width: 280px;
  flex-shrink: 0;
  border-right: 1px solid var(--border-color);
  transition: width 0.3s;
  overflow: hidden;
}

.sidebar.collapsed {
  width: 0;
  border-right: none;
}

/* 编辑器容器 */
.editor-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
}

/* 编辑器布局 */
.editor-layout {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.editor-layout.preview-only .editor-pane {
  display: none;
}

/* 编辑和预览面板 */
.editor-pane,
.preview-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
}

.editor-pane {
  border-right: 1px solid var(--border-color);
}

.preview-pane.expanded {
  flex: 1;
}

/* 面板头部 */
.pane-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: var(--bg-hover);
  border-bottom: 1px solid var(--border-color);
  font-size: 12px;
}

.pane-title {
  font-weight: 500;
  color: var(--text-secondary);
}

.pane-actions {
  display: flex;
  gap: 4px;
}

.pane-btn {
  padding: 2px 8px;
  border: 1px solid transparent;
  border-radius: 3px;
  background: transparent;
  cursor: pointer;
  font-size: 11px;
  color: var(--text-secondary);
}

.pane-btn:hover {
  background: var(--bg-component);
}

.pane-btn.active {
  background: var(--bg-active);
  color: var(--color-primary);
  border-color: var(--color-primary);
}

/* 编辑区域 */
.editor-content {
  flex: 1;
  overflow: auto;
}

.markdown-editor {
  width: 100%;
  height: 100%;
  padding: 20px;
  border: none;
  outline: none;
  resize: none;
  font-family: var(--font-family-mono);
  font-size: 14px;
  line-height: 1.8;
  background: var(--bg-component);
  color: var(--text-primary);
  tab-size: 2;
}

/* 预览区域 */
.preview-content {
  flex: 1;
  overflow: auto;
  padding: 20px;
  background: var(--bg-component);
}

.markdown-body {
  max-width: 800px;
  margin: 0 auto;
}

/* 标题编辑区 */
.title-editor {
  padding: 12px 20px 8px;
  background: var(--bg-component);
  border-bottom: 1px solid var(--border-color);
}

.title-input-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-component);
  color: var(--text-primary);
  font-size: 16px;
  font-weight: 500;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.title-input:focus {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.title-input::placeholder {
  color: var(--text-secondary);
}

.generate-outline-btn {
  padding: 8px 16px;
  border: 1px solid var(--color-primary);
  border-radius: 4px;
  background: var(--color-primary);
  color: white;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}

.generate-outline-btn:hover:not(:disabled) {
  background: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
}

.generate-outline-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-spinner-small {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 状态栏 */
.status-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 16px;
  background: var(--bg-hover);
  border-top: 1px solid var(--border-color);
  font-size: 12px;
  color: var(--text-secondary);
}

.status-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-item.dirty {
  color: var(--color-warning);
}

.status-separator {
  color: var(--border-color);
}

.status-btn {
  padding: 2px 8px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 11px;
  color: var(--text-secondary);
}

.status-btn:hover {
  color: var(--color-primary);
}
</style>
