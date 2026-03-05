<template>
  <div v-if="visible" class="modal-overlay" @click="handleClose">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>设置</h3>
        <button class="close-btn" @click="handleClose">✕</button>
      </div>
      <div class="modal-body">
        <!-- AI 配置 -->
        <div class="settings-section">
          <h4 class="section-title">AI 配置</h4>
          <div class="form-group">
            <label>Base URL</label>
            <input
              v-model="aiConfig.baseUrl"
              type="text"
              placeholder="https://api.deepseek.com/v1"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label>API Token</label>
            <input
              v-model="aiConfig.token"
              type="password"
              placeholder="输入你的 API Token"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label>Model</label>
            <input
              v-model="aiConfig.model"
              type="text"
              placeholder="deepseek-chat"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label>Temperature: {{ aiConfig.temperature }}</label>
            <input
              v-model.number="aiConfig.temperature"
              type="range"
              min="0"
              max="2"
              step="0.1"
              class="form-slider"
            />
            <div class="slider-hints">
              <span>精确 (0)</span>
              <span>平衡 (1)</span>
              <span>创意 (2)</span>
            </div>
          </div>
        </div>

        <!-- 公众号定位配置 -->
        <div class="settings-section">
          <h4 class="section-title">公众号定位</h4>
          <div class="form-group">
            <label>定位描述</label>
            <textarea
              v-model="positioningConfig"
              placeholder="例如：&#10;- 科技类公众号&#10;- 面向25-35岁职场人群&#10;- 专业但有温度&#10;&#10;这个配置会在生成大纲和文章时自动使用"
              class="form-textarea"
              rows="5"
            ></textarea>
            <p class="form-hint">保存后，生成大纲和文章时会自动应用此定位</p>
          </div>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn-secondary" @click="handleClose">关闭</button>
        <button class="btn-primary" @click="handleSave">保存</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue';

export interface AIConfig {
  baseUrl: string;
  token: string;
  temperature: number;
  model: string;
}

interface Props {
  visible: boolean;
  initialConfig?: AIConfig;
  initialPositioning?: string;
}

const props = withDefaults(defineProps<Props>(), {
  initialConfig: () => ({
    baseUrl: 'https://api.deepseek.com/v1',
    token: '',
    temperature: 0.7,
    model: 'deepseek-chat'
  }),
  initialPositioning: ''
});

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'save', config: AIConfig): void;
  (e: 'save-positioning', positioning: string): void;
}>();

const aiConfig = reactive<AIConfig>({
  baseUrl: props.initialConfig.baseUrl,
  token: props.initialConfig.token,
  temperature: props.initialConfig.temperature,
  model: props.initialConfig.model
});

const positioningConfig = ref(props.initialPositioning);

// 当外部配置变化时更新
watch(() => props.initialConfig, (newConfig) => {
  if (newConfig) {
    aiConfig.baseUrl = newConfig.baseUrl;
    aiConfig.token = newConfig.token;
    aiConfig.temperature = newConfig.temperature;
    aiConfig.model = newConfig.model;
  }
}, { deep: true });

// 当外部定位配置变化时更新
watch(() => props.initialPositioning, (newValue) => {
  positioningConfig.value = newValue;
});

const handleClose = () => {
  emit('update:visible', false);
};

const handleSave = () => {
  // 保存 AI 配置
  emit('save', { ...aiConfig });
  // 保存公众号定位
  emit('save-positioning', positioningConfig.value);
  handleClose();
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
  border-radius: 8px;
  width: 500px;
  max-width: 90vw;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.modal-header {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color, #e8e8e8);
}

.modal-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
}

.close-btn {
  padding: 4px 8px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 14px;
  color: var(--text-secondary, #8c8c8c);
  border-radius: 4px;
  transition: all 0.2s;
}

.close-btn:hover {
  background: var(--bg-hover, #f5f5f5);
  color: var(--text-primary, #262626);
}

.modal-body {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 20px;
  border-top: 1px solid var(--border-color, #e8e8e8);
}

.btn-secondary,
.btn-primary {
  padding: 6px 16px;
  border: 1px solid var(--border-color, #d9d9d9);
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
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

.btn-primary:hover {
  background: var(--color-primary-hover, #40a9ff);
  border-color: var(--color-primary-hover, #40a9ff);
}

/* 设置分区 */
.settings-section {
  margin-bottom: 24px;
}

.settings-section:last-child {
  margin-bottom: 0;
}

.section-title {
  margin: 0 0 16px 0;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary, #262626);
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-color, #e8e8e8);
}

/* 表单样式 */
.form-group {
  margin-bottom: 16px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 13px;
  color: var(--text-secondary, #595959);
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--border-color, #d9d9d9);
  border-radius: 4px;
  font-size: 14px;
  background: var(--bg-component, #fff);
  color: var(--text-primary, #262626);
  box-sizing: border-box;
  transition: all 0.2s;
}

.form-textarea {
  resize: vertical;
  font-family: inherit;
  line-height: 1.5;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: var(--color-primary, #1890ff);
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.form-slider {
  width: 100%;
  margin: 8px 0;
}

.slider-hints {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
  color: var(--text-secondary, #8c8c8c);
}

.form-hint {
  margin: 6px 0 0 0;
  font-size: 12px;
  color: var(--text-muted, #8c8c8c);
}

/* 暗色主题适配 */
@media (prefers-color-scheme: dark) {
  .modal-content {
    --bg-component: #1f1f1f;
    --bg-hover: #2c2c2c;
    --text-primary: #d9d9d9;
    --text-secondary: #8c8c8c;
    --text-muted: #595959;
    --border-color: #434343;
    --color-primary: #1890ff;
    --color-primary-hover: #40a9ff;
  }
}
</style>
