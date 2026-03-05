<template>
  <div class="file-toolbar">
    <!-- 主要操作 -->
    <div class="toolbar-group">
      <button
        class="toolbar-btn primary"
        @click="handleSave"
        :disabled="!canSave"
        :class="{ 'has-changes': isDirty }"
        title="保存 (Ctrl+S)"
      >
        <span class="icon">💾</span>
        <span>保存</span>
        <span v-if="isDirty" class="dirty-indicator">●</span>
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

    <!-- 文章信息 -->
    <div class="toolbar-group article-info" v-if="title">
      <span class="info-item" :title="title">
        {{ title }}
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
import { computed, onMounted, onUnmounted } from 'vue';
import { useWails } from '../composables/useWails';

/**
 * Props 定义
 */
interface Props {
  /** 文章标题 */
  title?: string;
  /** 是否有未保存更改 */
  isDirty?: boolean;
  /** 是否正在保存 */
  isSaving?: boolean;
  /** 是否可保存 */
  canSave?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  isDirty: false,
  isSaving: false,
  canSave: false,
});

/**
 * Emits 定义
 */
interface Emits {
  /** 请求保存 */
  (e: 'save'): void;
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
const isLoading = computed(() => wails.isLoading.value || props.isSaving);
const error = computed(() => wails.error.value);

/**
 * 处理保存
 */
const handleSave = (): void => {
  emit('save');
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
 * 键盘快捷键处理
 */
const handleKeyDown = (event: KeyboardEvent): void => {
  // Ctrl/Cmd + S: 保存
  if ((event.ctrlKey || event.metaKey) && event.key === 's' && !event.shiftKey) {
    event.preventDefault();
    handleSave();
  }
};

/**
 * 生命周期
 */
onMounted(() => {
  document.addEventListener('keydown', handleKeyDown);
});

onUnmounted(() => {
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

/* 文章信息 */
.article-info {
  flex: 1;
  justify-content: flex-end;
}

.info-item {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-color, #333);
  max-width: 400px;
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
