<template>
  <div class="writing-chat-sidebar">
    <!-- 步骤指示器 -->
    <div class="step-indicator">
      <div 
        v-for="(step, index) in steps" 
        :key="index"
        class="step-item"
        :class="{ 
          active: currentStep === index, 
          completed: currentStep > index,
          clickable: currentStep > index
        }"
        @click="currentStep > index && goToStep(index)"
      >
        <div class="step-number">
          <span v-if="currentStep > index">✓</span>
          <span v-else>{{ index + 1 }}</span>
        </div>
        <div class="step-title">{{ step.title }}</div>
      </div>
    </div>

    <!-- 对话区域 -->
    <div class="chat-container" ref="chatContainer">
      <!-- 欢迎消息 -->
      <div v-if="messages.length === 0" class="welcome-message">
        <div class="welcome-icon">📝</div>
        <div class="welcome-title">开始创作</div>
        <div class="welcome-desc">输入文章标题和基本要求，AI 将协助您完成创作</div>
      </div>

      <!-- 消息列表 -->
      <div v-else class="messages-list">
        <div 
          v-for="(msg, index) in messages" 
          :key="index"
          class="message-item"
          :class="{ 'user-message': msg.role === 'user', 'ai-message': msg.role === 'assistant' }"
        >
          <div class="message-avatar">
            {{ msg.role === 'user' ? '👤' : '🤖' }}
          </div>
          <div class="message-content">
            <div class="message-text" v-html="formatMessage(msg.content)"></div>
            <div v-if="msg.isStreaming" class="streaming-indicator">
              <span class="dot"></span>
              <span class="dot"></span>
              <span class="dot"></span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 快捷指令按钮 -->
    <div class="quick-actions">
      <div class="quick-actions-title">快捷指令</div>
      <div class="quick-actions-grid">
        <button 
          v-for="action in availableQuickActions" 
          :key="action.key"
          class="quick-action-btn"
          :class="{ active: currentStep === action.step }"
          :disabled="!canUseQuickAction(action) || isProcessing"
          @click="handleQuickAction(action)"
        >
          <span class="action-icon">{{ action.icon }}</span>
          <span class="action-label">{{ action.label }}</span>
        </button>
      </div>
    </div>

    <!-- 输入区域 -->
    <div class="input-area">
      <textarea
        v-model="inputMessage"
        :placeholder="getInputPlaceholder()"
        class="message-input"
        rows="3"
        @keydown.enter.prevent="handleEnterKey"
      ></textarea>
      <div class="input-actions">
        <button 
          class="send-btn"
          :disabled="!canSendMessage || isProcessing"
          @click="sendMessage"
        >
          <span v-if="isProcessing" class="btn-spinner"></span>
          <span v-else>发送</span>
        </button>
      </div>
    </div>

    <!-- 当前文章信息 -->
    <div v-if="articleTitle" class="current-article-info">
      <div class="info-title">📄 {{ articleTitle }}</div>
    </div>

    <!-- 最近文章 -->
    <div class="recent-section">
      <div class="recent-header" @click="showRecent = !showRecent">
        <span>📂 最近文章</span>
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
import { ref, computed, nextTick, onMounted } from 'vue';
import type { Article } from '../types';
import { useWails } from '../composables/useWails';

// ============================================
// 类型定义
// ============================================
interface Message {
  role: 'user' | 'assistant';
  content: string;
  isStreaming?: boolean;
  step?: number;
}

interface QuickAction {
  key: string;
  label: string;
  icon: string;
  step: number;
  requires?: string[];
  autoInput?: string;
}

interface Step {
  title: string;
  description: string;
}

// ============================================
// Props & Emits
// ============================================
const props = defineProps<{
  articleTitle: string;
  isProcessing: boolean;
  streamingContent?: string;
}>();

const emit = defineEmits<{
  (e: 'send-message', message: string, step: number): void;
  (e: 'quick-action', action: string, step: number): void;
  (e: 'open-recent', filePath: string): void;
  (e: 'update:step', step: number): void;
}>();

// ============================================
// 步骤配置
// ============================================
const steps: Step[] = [
  { title: '设定主题', description: '输入标题和基本要求' },
  { title: '生成框架', description: '生成文章大纲' },
  { title: '生成内容', description: '根据大纲写文章' },
  { title: '优化调整', description: '多轮优化文章' },
  { title: '生成标题', description: '生成爆款标题' },
];

// ============================================
// 快捷指令配置
// ============================================
const quickActions: QuickAction[] = [
  { key: 'start', label: '开始创作', icon: '🚀', step: 0, autoInput: '我想写一篇关于【主题】的文章，要求：【要求】' },
  { key: 'generate-outline', label: '生成大纲', icon: '📋', step: 1, requires: ['title'] },
  { key: 'regenerate-outline', label: '重新生成', icon: '🔄', step: 1, requires: ['title'] },
  { key: 'generate-article', label: '生成文章', icon: '✨', step: 2, requires: ['outline'] },
  { key: 'continue-write', label: '继续写作', icon: '✍️', step: 2, requires: ['outline'] },
  { key: 'polish', label: '润色优化', icon: '💎', step: 3, requires: ['content'] },
  { key: 'expand', label: '扩写内容', icon: '📈', step: 3, requires: ['content'] },
  { key: 'simplify', label: '精简内容', icon: '📉', step: 3, requires: ['content'] },
  { key: 'add-example', label: '添加案例', icon: '📚', step: 3, requires: ['content'] },
  { key: 'generate-titles', label: '生成标题', icon: '🎯', step: 4, requires: ['content'] },
];

// ============================================
// 状态
// ============================================
const currentStep = ref(0);
const messages = ref<Message[]>([]);
const inputMessage = ref('');
const chatContainer = ref<HTMLDivElement | null>(null);
const recentArticles = ref<Article[]>([]);
const showRecent = ref(false);

const wails = useWails();

// ============================================
// 计算属性
// ============================================
const availableQuickActions = computed(() => {
  return quickActions.filter(action => action.step === currentStep.value);
});

const canSendMessage = computed(() => {
  return inputMessage.value.trim().length > 0 && !props.isProcessing;
});

// ============================================
// 方法
// ============================================
const getInputPlaceholder = () => {
  switch (currentStep.value) {
    case 0:
      return '输入文章标题和基本要求...\n例如：\n标题：如何提高编程效率\n要求：面向初学者，通俗易懂';
    case 1:
      return '查看生成的大纲，可以提出修改意见...\n例如：请增加一个关于工具的章节';
    case 2:
      return '查看生成的文章内容，可以要求继续写作...';
    case 3:
      return '输入优化要求...\n例如：\n- 润色这段文字\n- 增加更多例子\n- 简化表达方式';
    case 4:
      return '输入标题优化要求...\n例如：\n- 生成更吸引人的标题\n- 添加数字和悬念';
    default:
      return '输入消息...';
  }
};

const formatMessage = (content: string): string => {
  // 简单的 Markdown 格式转换
  return content
    .replace(/\n/g, '<br>')
    .replace(/#{1,6}\s+(.+)/g, '<strong>$1</strong>')
    .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.+?)\*/g, '<em>$1</em>');
};

const canUseQuickAction = (action: QuickAction): boolean => {
  // 根据当前步骤和已有内容判断
  switch (action.step) {
    case 0:
      return currentStep.value === 0;
    case 1:
      return currentStep.value >= 0;
    case 2:
      return currentStep.value >= 1;
    case 3:
      return currentStep.value >= 2;
    case 4:
      return currentStep.value >= 3;
    default:
      return false;
  }
};

const handleQuickAction = (action: QuickAction) => {
  if (action.autoInput) {
    inputMessage.value = action.autoInput;
  } else {
    // 触发对应的动作
    emit('quick-action', action.key, currentStep.value);
  }
};

const handleEnterKey = (e: KeyboardEvent) => {
  if (e.shiftKey) {
    // Shift+Enter 换行
    return;
  }
  sendMessage();
};

const sendMessage = () => {
  if (!canSendMessage.value) return;

  const message = inputMessage.value.trim();
  if (!message) return;

  // 添加用户消息
  addMessage('user', message);
  
  // 清空输入
  inputMessage.value = '';

  // 触发发送事件
  emit('send-message', message, currentStep.value);
};

const addMessage = (role: 'user' | 'assistant', content: string, isStreaming = false) => {
  messages.value.push({
    role,
    content,
    isStreaming,
    step: currentStep.value,
  });
  scrollToBottom();
};

const addStreamingMessage = (initialContent = '') => {
  messages.value.push({
    role: 'assistant',
    content: initialContent,
    isStreaming: true,
    step: currentStep.value,
  });
  scrollToBottom();
};

const updateStreamingMessage = (content: string) => {
  const lastMsg = messages.value[messages.value.length - 1];
  if (lastMsg && lastMsg.isStreaming) {
    lastMsg.content = content;
  }
};

const finishStreaming = () => {
  const lastMsg = messages.value[messages.value.length - 1];
  if (lastMsg && lastMsg.isStreaming) {
    lastMsg.isStreaming = false;
  }
};

const scrollToBottom = async () => {
  await nextTick();
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
  }
};

const goToStep = (step: number) => {
  currentStep.value = step;
  emit('update:step', step);
};

const nextStep = () => {
  if (currentStep.value < steps.length - 1) {
    currentStep.value++;
    emit('update:step', currentStep.value);
  }
};

// ============================================
// 最近文章
// ============================================
const loadRecentArticles = async () => {
  try {
    const articles = await wails.getRecentArticles(10);
    recentArticles.value = articles;
  } catch (err) {
    console.error('加载最近文章失败:', err);
  }
};

const handleOpenRecent = (filePath: string) => {
  emit('open-recent', filePath);
};

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const day = 24 * 60 * 60 * 1000;

  if (diff < day) return '今天';
  if (diff < 2 * day) return '昨天';
  if (diff < 7 * day) return `${Math.floor(diff / day)}天前`;
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' });
};

// ============================================
// 暴露方法
// ============================================
defineExpose({
  currentStep,
  messages,
  addMessage,
  addStreamingMessage,
  updateStreamingMessage,
  finishStreaming,
  nextStep,
  goToStep,
  loadRecentArticles,
  scrollToBottom,
});

onMounted(() => {
  loadRecentArticles();
});
</script>

<style scoped>
.writing-chat-sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--panel-bg, #fafafa);
  overflow: hidden;
}

/* 步骤指示器 */
.step-indicator {
  display: flex;
  padding: 12px 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  overflow-x: auto;
  gap: 4px;
}

.step-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  min-width: 60px;
  cursor: default;
  opacity: 0.6;
  transition: all 0.2s;
}

.step-item.clickable {
  cursor: pointer;
}

.step-item.active,
.step-item.completed {
  opacity: 1;
}

.step-item.clickable:hover {
  opacity: 0.9;
}

.step-number {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: white;
  transition: all 0.2s;
}

.step-item.active .step-number {
  background: white;
  color: #667eea;
  transform: scale(1.1);
}

.step-item.completed .step-number {
  background: #52c41a;
}

.step-title {
  font-size: 11px;
  color: white;
  white-space: nowrap;
  font-weight: 500;
}

/* 对话容器 */
.chat-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: #f5f7fa;
}

/* 欢迎消息 */
.welcome-message {
  text-align: center;
  padding: 40px 20px;
  color: var(--text-secondary, #666);
}

.welcome-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.welcome-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary, #333);
  margin-bottom: 8px;
}

.welcome-desc {
  font-size: 13px;
  line-height: 1.6;
}

/* 消息列表 */
.messages-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.message-item {
  display: flex;
  gap: 10px;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-item.user-message {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.message-content {
  max-width: 85%;
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 13px;
  line-height: 1.6;
  word-break: break-word;
}

.user-message .message-content {
  background: #667eea;
  color: white;
  border-bottom-right-radius: 4px;
}

.ai-message .message-content {
  background: white;
  color: var(--text-primary, #333);
  border-bottom-left-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

/* 流式指示器 */
.streaming-indicator {
  display: flex;
  gap: 4px;
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid rgba(0, 0, 0, 0.1);
}

.user-message .streaming-indicator {
  border-top-color: rgba(255, 255, 255, 0.3);
}

.streaming-indicator .dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  opacity: 0.5;
  animation: bounce 1.4s infinite ease-in-out both;
}

.user-message .streaming-indicator .dot {
  background: white;
}

.streaming-indicator .dot:nth-child(1) {
  animation-delay: -0.32s;
}

.streaming-indicator .dot:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes bounce {
  0%, 80%, 100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1);
  }
}

/* 快捷指令 */
.quick-actions {
  padding: 12px;
  background: white;
  border-top: 1px solid var(--border-color, #e8e8e8);
}

.quick-actions-title {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-secondary, #666);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 8px;
}

.quick-actions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 6px;
}

.quick-action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 10px;
  border: 1px solid var(--border-color, #e8e8e8);
  border-radius: 6px;
  background: #f8f9fa;
  cursor: pointer;
  font-size: 12px;
  color: var(--text-primary, #333);
  transition: all 0.2s;
}

.quick-action-btn:hover:not(:disabled) {
  border-color: #667eea;
  background: #f0f2ff;
}

.quick-action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.quick-action-btn.active {
  border-color: #667eea;
  background: linear-gradient(135deg, #667eea15 0%, #764ba215 100%);
}

.action-icon {
  font-size: 14px;
}

.action-label {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 输入区域 */
.input-area {
  padding: 12px;
  background: white;
  border-top: 1px solid var(--border-color, #e8e8e8);
}

.message-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid var(--border-color, #d9d9d9);
  border-radius: 8px;
  font-size: 13px;
  line-height: 1.5;
  resize: none;
  font-family: inherit;
  background: #f8f9fa;
  transition: all 0.2s;
  box-sizing: border-box;
}

.message-input:focus {
  outline: none;
  border-color: #667eea;
  background: white;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.input-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 8px;
}

.send-btn {
  padding: 8px 20px;
  border: none;
  border-radius: 6px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;
}

.send-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.send-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
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

/* 当前文章信息 */
.current-article-info {
  padding: 10px 12px;
  background: #f0f2ff;
  border-top: 1px solid var(--border-color, #e8e8e8);
}

.info-title {
  font-size: 12px;
  color: #667eea;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 最近文章 */
.recent-section {
  border-top: 1px solid var(--border-color, #e8e8e8);
  background: white;
}

.recent-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary, #666);
  user-select: none;
}

.recent-header:hover {
  color: #667eea;
  background: #f8f9fa;
}

.arrow {
  font-size: 10px;
  transition: transform 0.2s;
}

.arrow.is-open {
  transform: rotate(180deg);
}

.recent-content {
  max-height: 150px;
  overflow-y: auto;
  border-top: 1px solid var(--border-color, #e8e8e8);
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
  padding: 8px 12px;
  cursor: pointer;
  transition: background 0.15s;
}

.recent-item:hover {
  background: #f0f2ff;
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
  .writing-chat-sidebar {
    --panel-bg: #141414;
    --text-primary: #d9d9d9;
    --text-secondary: #8c8c8c;
    --text-muted: #595959;
    --border-color: #303030;
  }

  .chat-container {
    background: #0f0f0f;
  }

  .welcome-message {
    color: #8c8c8c;
  }

  .welcome-title {
    color: #d9d9d9;
  }

  .message-avatar {
    background: #1f1f1f;
  }

  .ai-message .message-content {
    background: #1f1f1f;
    color: #d9d9d9;
  }

  .quick-actions,
  .input-area,
  .recent-section {
    background: #141414;
    border-color: #303030;
  }

  .quick-action-btn {
    background: #1f1f1f;
    border-color: #303030;
    color: #d9d9d9;
  }

  .quick-action-btn:hover:not(:disabled) {
    background: #1a1a2e;
    border-color: #667eea;
  }

  .message-input {
    background: #1f1f1f;
    border-color: #303030;
    color: #d9d9d9;
  }

  .current-article-info {
    background: #1a1a2e;
  }

  .recent-header:hover,
  .recent-item:hover {
    background: #1a1a2e;
  }
}
</style>
