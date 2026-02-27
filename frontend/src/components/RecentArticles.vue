<template>
  <div class="recent-articles">
    <div class="recent-header">
      <h3 class="recent-title">📄 最近打开</h3>
      <button
        class="refresh-btn"
        @click="loadRecent"
        :disabled="isLoading"
        title="刷新"
      >
        <span :class="{ 'is-spinning': isLoading }"">🔄</span>
      </button>
    </div>

    <div v-if="isLoading && articles.length === 0" class="loading-state">
      <span class="loading-spinner"></span>
      <span>加载中...</span>
    </div>

    <div v-else-if="articles.length === 0" class="empty-state">
      <div class="empty-icon">📭</div>
      <p>暂无最近打开的文章</p>
      <button class="action-btn" @click="handleOpen">打开文件</button>
    </div>

    <ul v-else class="article-list">
      <li
        v-for="article in articles"
        :key="article.uuid"
        class="article-item"
        :class="{ active: currentPath === article.filePath }"
        @click="handleClick(article)"
      >
        <div class="article-icon">📝</div>

        <div class="article-content">
          <div class="article-title" :title="article.title">
            {{ article.title || '无标题' }}
          </div>

          <div class="article-meta">
            <span class="meta-item" :title="article.filePath">
              {{ formatPath(article.filePath) }}
            </span>
            <span class="meta-item">{{ formatWordCount(article.wordCount) }}</span>
          </div>

          <div class="article-tags" v-if="article.tags?.length">
            <span
              v-for="tag in article.tags.slice(0, 3)"
              :key="tag"
              class="tag"
            >
              {{ tag }}
            </span>
          </div>
        </div>

        <div class="article-actions">
          <button
            class="action-icon-btn"
            @click.stop="handleRemove(article.uuid)"
            title="从列表移除"
          >
            ✕
          </button>
        </div>

        <div class="article-time" :title="formatFullDate(article.lastOpenedAt)">
          {{ formatRelativeDate(article.lastOpenedAt) }}
        </div>
      </li>
    </ul>

    <!-- 错误提示 -->
    <Transition name="slide">
      <div v-if="error" class="error-toast" @click="clearError">
        {{ error }}
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import type { Article } from '../types';
import { useWails } from '../composables/useWails';

/**
 * Props 定义
 */
interface Props {
  /** 当前打开的文件路径 */
  currentPath?: string;
}

const props = defineProps<Props>();

/**
 * Emits 定义
 */
interface Emits {
  /** 选择文章 */
  (e: 'select', filePath: string): void;
  /** 打开文件对话框 */
  (e: 'open'): void;
  /** 从列表移除 */
  (e: 'remove', uuid: string): void;
  /** 错误 */
  (e: 'error', message: string): void;
}

const emit = defineEmits<Emits>();

/**
 * 状态
 */
const articles = ref<Article[]>([]);
const isLoading = ref(false);
const error = ref<string | null>(null);

const wails = useWails();

/**
 * 加载最近文章
 */
const loadRecent = async (): Promise<void> => {
  isLoading.value = true;
  error.value = null;

  try {
    const result = await wails.getRecentArticles(20);
    articles.value = result;
  } catch (err) {
    const message = err instanceof Error ? err.message : '加载失败';
    error.value = message;
    emit('error', message);
  } finally {
    isLoading.value = false;
  }
};

/**
 * 处理点击
 */
const handleClick = (article: Article): void => {
  emit('select', article.filePath);
};

/**
 * 处理打开文件
 */
const handleOpen = (): void => {
  emit('open');
};

/**
 * 处理移除
 */
const handleRemove = async (uuid: string): Promise<void> => {
  try {
    await wails.deleteArticle(uuid);
    articles.value = articles.value.filter(a => a.uuid !== uuid);
  } catch (err) {
    const message = err instanceof Error ? err.message : '移除失败';
    error.value = message;
  }
};

/**
 * 清除错误
 */
const clearError = (): void => {
  error.value = null;
};

/**
 * 格式化路径
 */
const formatPath = (path: string): string => {
  if (!path) return '';

  // 显示父目录/文件名
  const parts = path.split(/[/\\]/);
  if (parts.length >= 2) {
    return `.../${parts[parts.length - 2]}/${parts[parts.length - 1]}`;
  }
  return path;
};

/**
 * 格式化字数
 */
const formatWordCount = (count: number): string => {
  if (!count) return '0 字';
  if (count < 1000) return `${count} 字`;
  return `${(count / 1000).toFixed(1)}k 字`;
};

/**
 * 格式化相对日期
 */
const formatRelativeDate = (dateStr: string): string => {
  const date = new Date(dateStr);
  const now = new Date();
  const diff = now.getTime() - date.getTime();

  const minute = 60 * 1000;
  const hour = 60 * minute;
  const day = 24 * hour;
  const week = 7 * day;

  if (diff < minute) return '刚刚';
  if (diff < hour) return `${Math.floor(diff / minute)}分钟前`;
  if (diff < day) return `${Math.floor(diff / hour)}小时前`;
  if (diff < week) return `${Math.floor(diff / day)}天前`;

  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' });
};

/**
 * 格式化完整日期
 */
const formatFullDate = (dateStr: string): string => {
  const date = new Date(dateStr);
  return date.toLocaleString('zh-CN');
};

/**
 * 暴露刷新方法
 */
defineExpose({
  refresh: loadRecent,
});

/**
 * 生命周期
 */
onMounted(() => {
  loadRecent();
});

/**
 * 监听当前路径变化，刷新高亮
 */
watch(() => props.currentPath, () => {
  // 可选：自动刷新列表
});
</script>

<style scoped>
.recent-articles {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--panel-bg, #fafafa);
}

.recent-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color, #e8e8e8);
}

.recent-title {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color, #333);
}

.refresh-btn {
  padding: 4px 8px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 14px;
  border-radius: 4px;
  transition: background 0.2s;
}

.refresh-btn:hover:not(:disabled) {
  background: var(--hover-bg, #e8e8e8);
}

.refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.is-spinning {
  display: inline-block;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 加载状态 */
.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 40px 20px;
  color: var(--text-secondary, #999);
  font-size: 13px;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid var(--border-color, #ddd);
  border-top-color: var(--primary-color, #1890ff);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
  color: var(--text-secondary, #999);
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state p {
  margin: 0 0 16px;
  font-size: 14px;
}

.action-btn {
  padding: 8px 16px;
  border: 1px solid var(--primary-color, #1890ff);
  border-radius: 4px;
  background: transparent;
  color: var(--primary-color, #1890ff);
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
}

.action-btn:hover {
  background: var(--primary-color, #1890ff);
  color: white;
}

/* 文章列表 */
.article-list {
  list-style: none;
  margin: 0;
  padding: 0;
  overflow-y: auto;
  flex: 1;
}

.article-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  border-bottom: 1px solid var(--border-light, #f0f0f0);
  transition: background 0.15s;
  position: relative;
}

.article-item:hover {
  background: var(--hover-bg, #f5f5f5);
}

.article-item.active {
  background: var(--active-bg, #e6f7ff);
  border-left: 3px solid var(--primary-color, #1890ff);
}

.article-icon {
  font-size: 18px;
  flex-shrink: 0;
}

.article-content {
  flex: 1;
  min-width: 0;
}

.article-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-color, #333);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.article-meta {
  display: flex;
  gap: 12px;
  font-size: 11px;
  color: var(--text-secondary, #999);
}

.meta-item {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.article-tags {
  display: flex;
  gap: 4px;
  margin-top: 6px;
  flex-wrap: wrap;
}

.tag {
  padding: 1px 6px;
  font-size: 10px;
  color: var(--text-secondary, #666);
  background: var(--tag-bg, #f0f0f0);
  border-radius: 3px;
}

.article-actions {
  opacity: 0;
  transition: opacity 0.15s;
}

.article-item:hover .article-actions {
  opacity: 1;
}

.action-icon-btn {
  padding: 4px 6px;
  border: none;
  background: transparent;
  color: var(--text-secondary, #999);
  cursor: pointer;
  font-size: 12px;
  border-radius: 3px;
}

.action-icon-btn:hover {
  background: var(--hover-bg, #e8e8e8);
  color: var(--error-color, #ff4d4f);
}

.article-time {
  position: absolute;
  right: 12px;
  bottom: 8px;
  font-size: 10px;
  color: var(--text-muted, #bfbfbf);
}

/* 错误提示 */
.error-toast {
  margin: 12px;
  padding: 10px 16px;
  background: var(--error-bg, #fff2f0);
  border: 1px solid var(--error-border, #ffccc7);
  border-radius: 4px;
  color: var(--error-color, #ff4d4f);
  font-size: 12px;
  text-align: center;
  cursor: pointer;
}

/* 过渡动画 */
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* 暗色主题适配 */
@media (prefers-color-scheme: dark) {
  .recent-articles {
    --panel-bg: #141414;
    --text-color: #d9d9d9;
    --text-secondary: #8c8c8c;
    --text-muted: #595959;
    --border-color: #303030;
    --border-light: #1f1f1f;
    --hover-bg: #1f1f1f;
    --active-bg: #111b26;
    --tag-bg: #1f1f1f;
    --error-bg: #2a1215;
    --error-border: #58181c;
    --error-color: #ff7875;
  }
}
</style>
