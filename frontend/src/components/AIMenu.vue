<template>
  <Teleport to="body">
    <div
      class="ai-menu"
      :style="{ left: position.x + 'px', top: position.y + 'px' }"
      @click.stop
    >
      <div class="ai-menu-header">
        <span class="ai-icon">✨</span>
        <span class="ai-title">AI 助手</span>
      </div>

      <div class="ai-actions">
        <button
          v-for="action in actions"
          :key="action.key"
          class="ai-action-btn"
          :disabled="isLoading"
          @click="handleAction(action.key)"
        >
          <span class="action-icon">{{ action.icon }}</span>
          <span class="action-label">{{ action.label }}</span>
        </button>
      </div>

      <div v-if="isLoading" class="ai-loading">
        <span class="loading-spinner"></span>
        <span>AI 处理中...</span>
      </div>

      <div v-else-if="result" class="ai-result">
        <div class="result-header">
          <span>AI 建议</span>
          <button class="close-btn" @click="closeResult">✕</button>
        </div>
        <div class="result-content">{{ result }}</div>
        <div class="result-actions">
          <button class="btn-secondary" @click="closeResult">取消</button>
          <button class="btn-primary" @click="applyResult">应用</button>
        </div>
      </div>

      <div v-if="error" class="ai-error">
        {{ error }}
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

interface Props {
  selectedText: string;
  position: { x: number; y: number };
}

const props = defineProps<Props>();

interface Emits {
  (e: 'apply', result: string): void;
  (e: 'close'): void;
}

const emit = defineEmits<Emits>();

// AI 操作选项
const actions = [
  { key: 'polish', label: '润色', icon: '✨' },
  { key: 'shorter', label: '精简', icon: '📉' },
  { key: 'continue', label: '续写', icon: '✍️' },
  { key: 'code', label: '解释代码', icon: '💻' },
];

// 状态
const isLoading = ref(false);
const result = ref('');
const error = ref('');

// 处理 AI 操作
const handleAction = async (action: string) => {
  isLoading.value = true;
  error.value = '';
  result.value = '';

  try {
    // 调用后端 AI 接口 (开发阶段使用模拟数据)
    // const response = await window.go?.main?.App?.CallDeepSeek?.(props.selectedText, action);

    // 开发阶段：直接使用模拟数据
    await simulateAIResponse(action);
  } catch (err) {
    // 模拟结果（开发阶段）
    await simulateAIResponse(action);
  } finally {
    isLoading.value = false;
  }
};

// 模拟 AI 响应（开发阶段使用）
const simulateAIResponse = async (action: string): Promise<void> => {
  return new Promise((resolve) => {
    setTimeout(() => {
      const responses: Record<string, string> = {
        polish: `【润色后的文本】\n\n${props.selectedText}\n\n（这里是 AI 润色后的结果，会使文字更加通顺专业）`,
        shorter: `【精简版本】\n\n${props.selectedText.slice(0, Math.floor(props.selectedText.length * 0.7))}...`,
        continue: `${props.selectedText}\n\n【续写内容】\n\n这里是 AI 根据上文继续写作的内容，保持风格一致...`,
        code: `【代码解释】\n\n这段代码的主要功能是：\n1. 实现了 xxx 功能\n2. 使用了 xxx 设计模式\n3. 注意 xxx 边界情况`,
      };
      result.value = responses[action] || 'AI 处理完成';
      resolve();
    }, 1000);
  });
};

// 应用结果
const applyResult = () => {
  emit('apply', result.value);
  emit('close');
};

// 关闭结果
const closeResult = () => {
  result.value = '';
  error.value = '';
};

// 点击外部关闭
const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement;
  if (!target.closest('.ai-menu')) {
    emit('close');
  }
};

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>

<style scoped>
.ai-menu {
  position: fixed;
  z-index: 10000;
  min-width: 200px;
  background: var(--bg-component, #fff);
  border: 1px solid var(--border-color, #d9d9d9);
  border-radius: var(--radius-lg, 6px);
  box-shadow: var(--shadow-3, 0 8px 16px rgba(0, 0, 0, 0.15));
  overflow: hidden;
}

.ai-menu-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 13px;
  font-weight: 500;
}

.ai-icon {
  font-size: 14px;
}

.ai-actions {
  display: flex;
  flex-direction: column;
  padding: 4px;
}

.ai-action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 13px;
  color: var(--text-primary, #262626);
  border-radius: var(--radius-sm, 2px);
  transition: background var(--transition-fast, 0.15s ease);
}

.ai-action-btn:hover:not(:disabled) {
  background: var(--bg-hover, #f5f5f5);
}

.ai-action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-icon {
  font-size: 14px;
}

.action-label {
  flex: 1;
  text-align: left;
}

.ai-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 16px;
  font-size: 13px;
  color: var(--text-secondary, #595959);
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid var(--border-color, #d9d9d9);
  border-top-color: var(--color-primary, #1890ff);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.ai-result {
  border-top: 1px solid var(--border-color, #d9d9d9);
  max-width: 400px;
  max-height: 300px;
  overflow: auto;
}

.result-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: var(--bg-hover, #f5f5f5);
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary, #595959);
  position: sticky;
  top: 0;
}

.close-btn {
  padding: 2px 6px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 12px;
  color: var(--text-muted, #8c8c8c);
}

.close-btn:hover {
  color: var(--text-primary, #262626);
}

.result-content {
  padding: 12px;
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-primary, #262626);
  white-space: pre-wrap;
}

.result-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  padding: 12px;
  border-top: 1px solid var(--border-light, #f0f0f0);
  position: sticky;
  bottom: 0;
  background: var(--bg-component, #fff);
}

.btn-secondary,
.btn-primary {
  padding: 6px 16px;
  border: 1px solid var(--border-color, #d9d9d9);
  border-radius: var(--radius-sm, 2px);
  cursor: pointer;
  font-size: 13px;
  transition: all var(--transition-fast, 0.15s ease);
}

.btn-secondary {
  background: var(--bg-component, #fff);
  color: var(--text-primary, #262626);
}

.btn-secondary:hover {
  border-color: var(--color-primary, #1890ff);
  color: var(--color-primary, #1890ff);
}

.btn-primary {
  background: var(--color-primary, #1890ff);
  border-color: var(--color-primary, #1890ff);
  color: white;
}

.btn-primary:hover {
  background: var(--color-primary-hover, #40a9ff);
  border-color: var(--color-primary-hover, #40a9ff);
}

.ai-error {
  padding: 12px;
  font-size: 13px;
  color: var(--color-error, #ff4d4f);
  background: #fff2f0;
  border-top: 1px solid #ffccc7;
}

/* 暗色主题适配 */
@media (prefers-color-scheme: dark) {
  .ai-menu {
    --bg-component: #141414;
    --text-primary: #d9d9d9;
    --text-secondary: #8c8c8c;
    --text-muted: #595959;
    --border-color: #303030;
    --border-light: #1f1f1f;
    --bg-hover: #1f1f1f;
    --bg-active: #111b26;
    --color-primary: #4ec9b0;
    --color-primary-hover: #5ddbc1;
    --color-error: #ff7875;
  }

  .ai-error {
    background: #2a1215;
    border-top-color: #58181c;
  }
}
</style>
