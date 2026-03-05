<template>
  <div class="writing-sidebar">
    <!-- 标题输入区 -->
    <div class="sidebar-section">
      <label class="section-label">文章标题</label>
      <input
        v-model="localTitle"
        type="text"
        placeholder="输入文章标题..."
        class="title-input"
        @blur="handleTitleBlur"
      />
    </div>

    <!-- 写作要求 -->
    <div class="sidebar-section">
      <label class="section-label">写作要求</label>
      <textarea
        v-model="localRequirements"
        placeholder="例如：&#10;- 面向初学者&#10;- 通俗易懂&#10;- 包含实用案例&#10;- 语气亲切"
        class="requirements-input"
        rows="6"
      ></textarea>
    </div>

    <!-- 操作按钮区 -->
    <div class="sidebar-actions">
      <button
        class="action-btn outline-btn"
        :disabled="isGeneratingOutline || !localTitle"
        @click="handleGenerateOutline"
      >
        <span v-if="isGeneratingOutline" class="btn-spinner"></span>
        <span v-else>📋</span>
        <span>{{ isGeneratingOutline ? '生成中...' : '生成大纲' }}</span>
      </button>

      <button
        class="action-btn article-btn"
        :disabled="isGeneratingArticle || !localOutline || !localTitle"
        @click="handleGenerateArticle"
      >
        <span v-if="isGeneratingArticle" class="btn-spinner"></span>
        <span v-else>✨</span>
        <span>{{ isGeneratingArticle ? '写作中...' : '生成文章' }}</span>
      </button>
    </div>

    <!-- 提示信息 -->
    <div class="sidebar-tips">
      <p v-if="!localTitle" class="tip-text">请先输入文章标题</p>
      <p v-else-if="!localOutline" class="tip-text">点击「生成大纲」开始创作</p>
      <p v-else class="tip-text success">大纲已生成，可编辑后生成文章</p>
    </div>

    <!-- 最近文章折叠面板 -->
    <div class="recent-section">
      <div class="recent-header" @click="showRecent = !showRecent">
        <span>📄 最近文章</span>
        <span class="arrow" :class="{ 'is-open': showRecent }">▼</span>
      </div>
      
      <div v-show="showRecent" class="recent-content">
        <div v-if="recentArticles.length === 0" class="recent-empty">
          暂无最近文章
        </div>
        <ul v-else class="recent-list">
          <li
            v-for="article in recentArticles.slice(0, 5)"
            :key="article.uuid"
            class="recent-item"
            @click="handleOpenRecent(article.filePath)"
          >
            <span class="recent-title">{{ article.title || '无标题' }}</span>
            <span class="recent-date">{{ formatDate(article.lastOpenedAt) }}</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import type { Article } from '../types';
import { useWails } from '../composables/useWails';

/**
 * Props 定义
 */
interface Props {
  /** 当前标题 */
  title?: string;
  /** 当前大纲 */
  outline?: string;
  /** 当前写作要求 */
  requirements?: string;
  /** 公众号定位（从配置读取） */
  positioning?: string;
  /** 是否正在生成大纲 */
  isGeneratingOutline?: boolean;
  /** 是否正在生成文章 */
  isGeneratingArticle?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  outline: '',
  requirements: '',
  positioning: '',
  isGeneratingOutline: false,
  isGeneratingArticle: false,
});

/**
 * Emits 定义
 */
interface Emits {
  /** 更新标题 */
  (e: 'update:title', value: string): void;
  /** 更新大纲 */
  (e: 'update:outline', value: string): void;
  /** 更新写作要求 */
  (e: 'update:requirements', value: string): void;
  /** 生成大纲 */
  (e: 'generate-outline'): void;
  /** 生成文章 */
  (e: 'generate-article'): void;
  /** 打开最近文章 */
  (e: 'open-recent', filePath: string): void;
  /** 错误 */
  (e: 'error', message: string): void;
}

const emit = defineEmits<Emits>();

/**
 * 状态
 */
const localTitle = computed({
  get: () => props.title,
  set: (value) => emit('update:title', value),
});

const localOutline = computed({
  get: () => props.outline,
  set: (value) => emit('update:outline', value),
});

const localRequirements = computed({
  get: () => props.requirements,
  set: (value) => emit('update:requirements', value),
});

const recentArticles = ref<Article[]>([]);
const showRecent = ref(false);

const wails = useWails();

/**
 * 加载最近文章
 */
const loadRecentArticles = async (): Promise<void> => {
  try {
    const articles = await wails.getRecentArticles(10);
    recentArticles.value = articles;
  } catch (err) {
    console.error('加载最近文章失败:', err);
  }
};

/**
 * 处理标题失焦
 */
const handleTitleBlur = (): void => {
  // 标题变化时触发更新
  emit('update:title', localTitle.value);
};

/**
 * 生成大纲 - 触发父组件处理
 */
const handleGenerateOutline = (): void => {
  if (!localTitle.value) {
    emit('error', '请先输入文章标题');
    return;
  }

  emit('generate-outline');
};

/**
 * 生成文章 - 触发父组件处理
 */
const handleGenerateArticle = (): void => {
  if (!localTitle.value) {
    emit('error', '请先输入文章标题');
    return;
  }

  if (!localOutline.value) {
    emit('error', '请先生成大纲');
    return;
  }

  emit('generate-article');
};

/**
 * 打开最近文章
 */
const handleOpenRecent = (filePath: string): void => {
  emit('open-recent', filePath);
};

/**
 * 格式化日期
 */
const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr);
  const now = new Date();
  const diff = now.getTime() - date.getTime();

  const day = 24 * 60 * 60 * 1000;

  if (diff < day) {
    return '今天';
  }
  if (diff < 2 * day) {
    return '昨天';
  }
  if (diff < 7 * day) {
    return `${Math.floor(diff / day)}天前`;
  }

  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' });
};

/**
 * 暴露方法
 */
defineExpose({
  loadRecentArticles,
});

onMounted(() => {
  loadRecentArticles();
});
</script>

<style scoped>
.writing-sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--panel-bg, #fafafa);
  padding: 16px;
  box-sizing: border-box;
  overflow-y: auto;
}

/* 侧边栏区块 */
.sidebar-section {
  margin-bottom: 16px;
}

.section-label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary, #666);
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.title-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color, #d9d9d9);
  border-radius: 6px;
  font-size: 14px;
  background: var(--bg-component, #fff);
  color: var(--text-primary, #262626);
  box-sizing: border-box;
  transition: all 0.2s;
}

.title-input:focus {
  outline: none;
  border-color: var(--color-primary, #1890ff);
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.title-input::placeholder {
  color: var(--text-muted, #bfbfbf);
}

.requirements-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color, #d9d9d9);
  border-radius: 6px;
  font-size: 13px;
  line-height: 1.6;
  background: var(--bg-component, #fff);
  color: var(--text-primary, #262626);
  box-sizing: border-box;
  resize: vertical;
  min-height: 80px;
  font-family: inherit;
  transition: all 0.2s;
}

.requirements-input:focus {
  outline: none;
  border-color: var(--color-primary, #1890ff);
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.requirements-input::placeholder {
  color: var(--text-muted, #bfbfbf);
}

/* 操作按钮区 */
.sidebar-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 8px;
  margin-bottom: 12px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.outline-btn {
  background: var(--color-primary, #1890ff);
  color: white;
}

.outline-btn:hover:not(:disabled) {
  background: var(--color-primary-hover, #40a9ff);
}

.article-btn {
  background: var(--color-success, #52c41a);
  color: white;
}

.article-btn:hover:not(:disabled) {
  background: var(--color-success-hover, #73d13d);
}

.btn-spinner {
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

/* 提示信息 */
.sidebar-tips {
  margin-bottom: 16px;
  padding: 10px 12px;
  background: var(--bg-hover, #f0f0f0);
  border-radius: 6px;
  text-align: center;
}

.tip-text {
  margin: 0;
  font-size: 12px;
  color: var(--text-secondary, #666);
}

.tip-text.success {
  color: var(--color-success, #52c41a);
}

/* 最近文章区 */
.recent-section {
  margin-top: auto;
  border-top: 1px solid var(--border-color, #e8e8e8);
  padding-top: 12px;
}

.recent-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 0;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary, #333);
  user-select: none;
}

.recent-header:hover {
  color: var(--color-primary, #1890ff);
}

.arrow {
  font-size: 10px;
  transition: transform 0.2s;
}

.arrow.is-open {
  transform: rotate(180deg);
}

.recent-content {
  max-height: 200px;
  overflow-y: auto;
}

.recent-empty {
  padding: 16px;
  text-align: center;
  font-size: 12px;
  color: var(--text-muted, #999);
}

.recent-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.recent-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 4px;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.15s;
}

.recent-item:hover {
  background: var(--bg-hover, #f0f0f0);
}

.recent-title {
  font-size: 12px;
  color: var(--text-primary, #333);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 140px;
}

.recent-date {
  font-size: 11px;
  color: var(--text-muted, #999);
  flex-shrink: 0;
}

/* 暗色主题适配 */
@media (prefers-color-scheme: dark) {
  .writing-sidebar {
    --panel-bg: #141414;
    --text-primary: #d9d9d9;
    --text-secondary: #8c8c8c;
    --text-muted: #595959;
    --border-color: #303030;
    --bg-hover: #1f1f1f;
    --bg-component: #1f1f1f;
    --color-primary: #1890ff;
    --color-primary-hover: #40a9ff;
    --color-success: #52c41a;
    --color-success-hover: #73d13d;
  }
}
</style>
