<template>
  <div v-if="visible" class="modal-overlay" @click="handleClose">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>设置</h3>
        <button class="close-btn" @click="handleClose">✕</button>
      </div>
      <div class="modal-body">
        <p>设置功能开发中...</p>
      </div>
      <div class="modal-footer">
        <button class="btn-secondary" @click="handleClose">关闭</button>
        <button class="btn-primary" @click="handleSave">保存</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  visible: boolean;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void;
  (e: 'save', settings: Record<string, unknown>): void;
}>();

const handleClose = () => {
  emit('update:visible', false);
};

const handleSave = () => {
  emit('save', {});
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
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.modal-header {
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
}

.modal-body {
  padding: 20px;
  min-height: 100px;
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
}

.btn-secondary {
  background: var(--bg-component, #fff);
  color: var(--text-primary, #262626);
}

.btn-primary {
  background: var(--color-primary, #1890ff);
  border-color: var(--color-primary, #1890ff);
  color: white;
}
</style>
