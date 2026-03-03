<template>
  <div class="file-toolbar">
    <!-- 主要操作 -->
    <div class="toolbar-group">
      <button
        class="toolbar-btn primary"
        @click="handleNew"
        :disabled="isLoading"
        title="新建文章 (Ctrl+N)"
      >
        <span class="icon">+</span>
        <span>新建</span>
      </button>

      <button
        class="toolbar-btn"
        @click="handleOpen"
        :disabled="isLoading"
        title="打开文件 (Ctrl+O)"
      >
        <span class="icon">📂</span>
        <span>打开</span>
      </button>

      <button
        class="toolbar-btn"
        @click="handleSave"
        :disabled="!canSave"
        :class="{ 'has-changes': isDirty }"
        title="保存 (Ctrl+S)"
      >
        <span class="icon">💾</span>
        <span>保存</span>
        <span v-if="isDirty" class="dirty-indicator">●</span>
      </button>

      <button
        class="toolbar-btn"
        @click="handleSaveAs"
        :disabled="!article"
        title="另存为 (Ctrl+Shift+S)"
      >
        <span>另存为</span>
      </button>
    </div>

    <!-- 分隔线 -->
    <div class="toolbar-separator"></div>

    <!-- 设置按钮 -->
    <button
      class="toolbar-btn"
      @click="handleSettings"
      title="设置 (Ctrl+,)"
    >
      <span class="icon">⚙️</span>
      <span>设置</span>
    </button>

    <!-- 分隔线 -->
    <div class="toolbar-separator"></div>

    <!-- 最近文件 -->
    <div class="toolbar-group recent-dropdown" v-if="recentArticles.length > 0">
      <button
        class="toolbar-btn dropdown-trigger"
        @click="toggleRecent"
        :disabled="isLoading"
      >
        <span>最近文件</span>
        <span class="dropdown-arrow" :class="{ 'is-open': showRecent }">▼</span>
      </button>

      <!-- 最近文件下拉菜单 -->
      <Transition name="dropdown">
        <div v-show="showRecent" class="dropdown-menu">
          <div
            v-for="item in recentArticles"
            :key="item.uuid"
            class="dropdown-item"
            @click="handleOpenRecent(item.filePath)"
          >
            <div class="item-title">{{ item.title }}</div>
            <div class="item-meta">
              <span>{{ formatWordCount(item.wordCount) }}</span>
              <span>{{ formatDate(item.lastOpenedAt) }}</span>
            </div>
          </div>

          <div class="dropdown-divider"></div>

          <div class="dropdown-item" @click="refreshRecent">
            <span class="icon">🔄</span>
            <span>刷新列表</span>
          </div>
        </div>
      </Transition>
    </div>

    <!-- 分隔线 -->
    <div class="toolbar-separator"></div>

    <!-- 文章信息 -->
    <div class="toolbar-group article-info" v-if="article">
      <span class="info-item" :title="article.filePath">
        {{ article.title || formatFileName(article.filePath) }}
      </span>
    </div>

    <!-- 加载状态 -->
    <div v-if="isLoading" class="toolbar-status">
      <span class="loading-spinner"></span>
      <span>处理中...</span>
    </div>

    <!-- 错误提示 -->
    <Transition name="fade">
      <div v-if="error" class="toolbar-error" @click="clearError">
        <span class="error-icon">⚠️</span>
        <span class="error-message">{{ error }}</span>
        <span class="error-close">✕</span>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import type { Article } from '../types';
import { useWails } from '../composables/useWails';

/**
 * Props 定义
 */
interface Props {
  /** 当前文章 */
  article?: Article | null;
  /** 是否有未保存更改 */
  isDirty?: boolean;
  /** 是否正在保存 */
  isSaving?: boolean;
  /** 是否可保存 */
  canSave?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  article: null,
  isDirty: false,
  isSaving: false,
  canSave: false,
});

/**
 * Emits 定义
 */
interface Emits {
  /** 请求新建 */
  (e: 'new'): void;
  /** 请求打开 */
  (e: 'open', filePath: string): void;
  /** 请求保存 */
  (e: 'save'): void;
  /** 请求另存为 */
  (e: 'saveAs'): void;
  /** 打开设置 */
  (e: 'settings'): void;
  /** 错误 */
  (e: 'error', message: string): void;
}

const emit = defineEmits<Emits>();

/**
 * 组合式函数
 */
const wails = useWails();

/**
 * 状态
 */
const recentArticles = ref<Article[]>([]);
const showRecent = ref(false);
const isLoading = computed(() => wails.isLoading.value || props.isSaving);
const error = computed(() => wails.error.value);

/**
 * 加载最近文章列表
 */
const loadRecent = async (): Promise<void> => {
  try {
    const articles = await wails.getRecentArticles(10);
    recentArticles.value = articles;
  } catch (err) {
    console.error('加载最近文章失败:', err);
  }
};

/**
 * 刷新最近文章列表
 */
const refreshRecent = async (): Promise<void> => {
  await loadRecent();
};

/**
 * 切换最近文件下拉菜单
 */
const toggleRecent = (): void => {
  showRecent.value = !showRecent.value;
  if (showRecent.value) {
    loadRecent();
  }
};

/**
 * 关闭最近文件下拉菜单
 */
const closeRecent = (event: MouseEvent): void => {
  const target = event.target as HTMLElement;
  if (!target.closest('.recent-dropdown')) {
    showRecent.value = false;
  }
};

/**
 * 处理新建
 */
const handleNew = (): void => {
  if (props.isDirty) {
    // 如果有未保存更改，可以在这里添加确认对话框
    const confirm = window.confirm('当前文章有未保存的更改，是否继续？');
    if (!confirm) return;
  }
  emit('new');
};

/**
 * 处理打开
 */
const handleOpen = async (): Promise<void> => {
  if (props.isDirty) {
    const confirm = window.confirm('当前文章有未保存的更改，是否继续？');
    if (!confirm) return;
  }

  try {
    const filePath = await wails.openFileDialog();
    if (filePath) {
      emit('open', filePath);
    }
  } catch (err) {
    const message = err instanceof Error ? err.message : '打开文件失败';
    emit('error', message);
  }
};

/**
 * 处理打开最近文件
 */
const handleOpenRecent = async (filePath: string): Promise<void> => {
  showRecent.value = false;

  if (props.isDirty) {
    const confirm = window.confirm('当前文章有未保存的更改，是否继续？');
    if (!confirm) return;
  }

  emit('open', filePath);
};

/**
 * 处理保存
 */
const handleSave = (): void => {
  emit('save');
};

/**
 * 处理另存为
 */
const handleSaveAs = (): void => {
  emit('saveAs');
};

/**
 * 处理设置
 */
const handleSettings = (): void => {
  emit('settings');
};

/**
 * 清除错误
 */
const clearError = (): void => {
  wails.clearError();
};

/**
 * 格式化字数
 */
const formatWordCount = (count: number): string => {
  if (count < 1000) return `${count} 字`;
  return `${(count / 1000).toFixed(1)}k 字`;
};

/**
 * 格式化日期
 */
const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr);
  const now = new Date();
  const diff = now.getTime() - date.getTime();

  // 小于1小时
  if (diff < 60 * 60 * 1000) {
    const minutes = Math.floor(diff / (60 * 1000));
    return minutes < 1 ? '刚刚' : `${minutes}分钟前`;
  }

  // 小于24小时
  if (diff < 24 * 60 * 60 * 1000) {
    const hours = Math.floor(diff / (60 * 60 * 1000));
    return `${hours}小时前`;
  }

  // 小于7天
  if (diff < 7 * 24 * 60 * 60 * 1000) {
    const days = Math.floor(diff / (24 * 60 * 60 * 1000));
    return `${days}天前`;
  }

  // 默认显示日期
  return date.toLocaleDateString('zh-CN');
};

/**
 * 格式化文件名
 */
const formatFileName = (filePath: string): string => {
  if (!filePath) return '';
  const parts = filePath.split(/[/\\]/);
  return parts[parts.length - 1] || filePath;
};

/**
 * 键盘快捷键处理
 */
const handleKeyDown = (event: KeyboardEvent): void => {
  // Ctrl/Cmd + N: 新建
  if ((event.ctrlKey || event.metaKey) && event.key === 'n') {
    event.preventDefault();
    handleNew();
  }

  // Ctrl/Cmd + O: 打开
  if ((event.ctrlKey || event.metaKey) && event.key === 'o') {
    event.preventDefault();
    handleOpen();
  }

  // Ctrl/Cmd + S: 保存
  if ((event.ctrlKey || event.metaKey) && event.key === 's' && !event.shiftKey) {
    event.preventDefault();
    handleSave();
  }

  // Ctrl/Cmd + Shift + S: 另存为
  if ((event.ctrlKey || event.metaKey) && event.shiftKey && event.key === 'S') {
    event.preventDefault();
    handleSaveAs();
  }
};

/**
 * 生命周期
 */
onMounted(() => {
  loadRecent();
  document.addEventListener('click', closeRecent);
  document.addEventListener('keydown', handleKeyDown);
});

onUnmounted(() => {
  document.removeEventListener('click', closeRecent);
  document.removeEventListener('keydown', handleKeyDown);
});
</script>

<style scoped>
.file-toolbar {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  background: var(--toolbar-bg, #f5f5f5);
  border-bottom: 1px solid var(--border-color, #ddd);
  gap: 12px;
  position: relative;
}

.toolbar-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  border: 1px solid transparent;
  border-radius: 4px;
  background: transparent;
  cursor: pointer;
  font-size: 13px;
  color: var(--text-color, #333);
  transition: all 0.2s;
}

.toolbar-btn:hover:not(:disabled) {
  background: var(--btn-hover-bg, #e8e8e8);
  border-color: var(--border-color, #ddd);
}

.toolbar-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.toolbar-btn.primary {
  background: var(--primary-color, #1890ff);
  color: white;
}

.toolbar-btn.primary:hover:not(:disabled) {
  background: var(--primary-hover, #40a9ff);
}

.toolbar-btn.has-changes {
  font-weight: 600;
}

.dirty-indicator {
  color: var(--warning-color, #faad14);
  font-size: 10px;
}

.toolbar-separator {
  width: 1px;
  height: 24px;
  background: var(--border-color, #ddd);
}

/* 最近文件下拉菜单 */
.recent-dropdown {
  position: relative;
}

.dropdown-arrow {
  font-size: 10px;
  transition: transform 0.2s;
}

.dropdown-arrow.is-open {
  transform: rotate(180deg);
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  margin-top: 4px;
  min-width: 280px;
  max-height: 400px;
  overflow-y: auto;
  background: var(--dropdown-bg, white);
  border: 1px solid var(--border-color, #ddd);
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1000;
}

.dropdown-item {
  padding: 10px 14px;
  cursor: pointer;
  border-bottom: 1px solid var(--border-light, #f0f0f0);
}

.dropdown-item:last-of-type {
  border-bottom: none;
}

.dropdown-item:hover {
  background: var(--hover-bg, #f5f5f5);
}

.item-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-color, #333);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-meta {
  display: flex;
  gap: 12px;
  margin-top: 4px;
  font-size: 11px;
  color: var(--text-secondary, #666);
}

.dropdown-divider {
  height: 1px;
  background: var(--border-color, #ddd);
  margin: 4px 0;
}

/* 文章信息 */
.article-info {
  flex: 1;
  justify-content: flex-end;
}

.info-item {
  font-size: 12px;
  color: var(--text-secondary, #666);
  max-width: 300px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 状态指示 */
.toolbar-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-secondary, #666);
}

.loading-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid var(--border-color, #ddd);
  border-top-color: var(--primary-color, #1890ff);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 错误提示 */
.toolbar-error {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--error-bg, #fff2f0);
  border: 1px solid var(--error-border, #ffccc7);
  border-radius: 4px;
  color: var(--error-color, #ff4d4f);
  font-size: 13px;
  cursor: pointer;
  white-space: nowrap;
  z-index: 100;
}

.error-close {
  margin-left: 8px;
  font-size: 11px;
}

/* 过渡动画 */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 暗色主题适配 */
@media (prefers-color-scheme: dark) {
  .file-toolbar {
    --toolbar-bg: #1f1f1f;
    --text-color: #d9d9d9;
    --text-secondary: #8c8c8c;
    --border-color: #434343;
    --btn-hover-bg: #2c2c2c;
    --dropdown-bg: #262626;
    --hover-bg: #2c2c2c;
    --error-bg: #2a1215;
    --error-border: #58181c;
    --error-color: #ff7875;
  }
}
</style>
