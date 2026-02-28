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
import SettingsModal from './SettingsModal.vue';
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
const showAIMenu = ref(false);
const selectedText = ref('');
const aiMenuPosition = ref({ x: 0, y: 0 });
const previewRef = ref<HTMLElement | null>(null);

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
  }
};

const handleOpen = async (filePath: string) => {
  try {
    await editor.loadArticle(filePath);
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
// AI 功能
// ============================================
const handleAIApply = (result: string) => {
  console.log('AI 结果:', result);
  showAIMenu.value = false;
};

// ============================================
// 设置
// ============================================
const handleSettingsSave = (settings: Record<string, unknown>) => {
  console.log('保存设置:', settings);
  showSettings.value = false;
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
