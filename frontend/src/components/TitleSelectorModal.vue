<template>
  <div v-if="visible" class="modal-overlay" @click="handleClose">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>🎯 选择爆款标题</h3>
        <p class="modal-subtitle">AI 为您生成了3个候选标题，请选择最喜欢的一个</p>
        <button class="close-btn" @click="handleClose">✕</button>
      </div>
      
      <div class="modal-body">
        <div class="titles-list">
          <div
            v-for="(title, index) in titles"
            :key="index"
            class="title-option"
            :class="{ selected: selectedIndex === index }"
            @click="selectedIndex = index"
          >
            <div class="title-number">{{ index + 1 }}</div>
            <div class="title-text">{{ title }}</div>
            <div class="title-check">
              <span v-if="selectedIndex === index" class="check-icon">✓</span>
            </div>
          </div>
        </div>
        
        <!-- 自定义标题输入 -->
        <div class="custom-title-section">
          <label class="custom-label">或者自定义标题：</label>
          <input
            v-model="customTitle"
            type="text"
            placeholder="输入您自己的标题..."
            class="custom-input"
            @focus="selectedIndex = -1"
          />
        </div>
      </div>
      
      <div class="modal-footer">
        <button class="btn-secondary" @click="handleClose">取消</button>
        <button class="btn-primary" :disabled="!canConfirm" @click="handleConfirm">
          确认选择
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

interface Props {
  visible: boolean;
  titles: string[];
  originalTitle?: string;
}

const props = withDefaults(defineProps<Props>(), {
  titles: () => [],
  originalTitle: '',
});

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'select', title: string): void;
}>();

const selectedIndex = ref(0);
const customTitle = ref('');

// 当弹窗显示时，默认选择第一个标题
watch(() => props.visible, (newVal) => {
  if (newVal) {
    selectedIndex.value = 0;
    customTitle.value = '';
  }
});

const canConfirm = computed(() => {
  if (selectedIndex.value >= 0 && selectedIndex.value < props.titles.length) {
    return true;
  }
  return customTitle.value.trim() !== '';
});

const handleClose = () => {
  emit('update:visible', false);
};

const handleConfirm = () => {
  let selectedTitle = '';
  
  if (selectedIndex.value >= 0 && selectedIndex.value < props.titles.length) {
    selectedTitle = props.titles[selectedIndex.value];
  } else if (customTitle.value.trim()) {
    selectedTitle = customTitle.value.trim();
  }
  
  if (selectedTitle) {
    emit('select', selectedTitle);
    emit('update:visible', false);
  }
};
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: var(--bg-component, #fff);
  border-radius: 12px;
  width: 500px;
  max-width: 90vw;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
}

.modal-header {
  position: relative;
  padding: 24px 24px 16px;
  border-bottom: 1px solid var(--border-color, #e8e8e8);
}

.modal-header h3 {
  margin: 0 0 8px 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary, #262626);
}

.modal-subtitle {
  margin: 0;
  font-size: 13px;
  color: var(--text-secondary, #8c8c8c);
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  padding: 4px 8px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 16px;
  color: var(--text-secondary, #8c8c8c);
  border-radius: 4px;
  transition: all 0.2s;
}

.close-btn:hover {
  background: var(--bg-hover, #f5f5f5);
  color: var(--text-primary, #262626);
}

.modal-body {
  padding: 20px 24px;
  overflow-y: auto;
}

.titles-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 20px;
}

.title-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border: 2px solid var(--border-color, #e8e8e8);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  background: var(--bg-component, #fff);
}

.title-option:hover {
  border-color: var(--color-primary, #1890ff);
  background: var(--bg-hover, #f5f5f5);
}

.title-option.selected {
  border-color: var(--color-primary, #1890ff);
  background: rgba(24, 144, 255, 0.05);
}

.title-number {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary, #1890ff);
  color: white;
  font-size: 13px;
  font-weight: 600;
  border-radius: 50%;
  flex-shrink: 0;
}

.title-text {
  flex: 1;
  font-size: 15px;
  color: var(--text-primary, #262626);
  line-height: 1.5;
}

.title-check {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.check-icon {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary, #1890ff);
  color: white;
  font-size: 14px;
  border-radius: 50%;
}

.custom-title-section {
  padding-top: 16px;
  border-top: 1px dashed var(--border-color, #e8e8e8);
}

.custom-label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary, #595959);
  margin-bottom: 8px;
}

.custom-input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border-color, #d9d9d9);
  border-radius: 6px;
  font-size: 14px;
  background: var(--bg-component, #fff);
  color: var(--text-primary, #262626);
  box-sizing: border-box;
  transition: all 0.2s;
}

.custom-input:focus {
  outline: none;
  border-color: var(--color-primary, #1890ff);
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid var(--border-color, #e8e8e8);
}

.btn-secondary,
.btn-primary {
  padding: 8px 20px;
  border: 1px solid var(--border-color, #d9d9d9);
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
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

.btn-primary:hover:not(:disabled) {
  background: var(--color-primary-hover, #40a9ff);
  border-color: var(--color-primary-hover, #40a9ff);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 暗色主题适配 */
@media (prefers-color-scheme: dark) {
  .modal-content {
    --bg-component: #1f1f1f;
    --bg-hover: #2c2c2c;
    --text-primary: #d9d9d9;
    --text-secondary: #8c8c8c;
    --border-color: #434343;
    --color-primary: #1890ff;
    --color-primary-hover: #40a9ff;
  }
}
</style>
