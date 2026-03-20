<template>
  <div 
    class="floating-ai-actions" 
    :style="{ top: position.y + 'px', left: position.x + 'px' }"
  >
    <div class="actions-header">
      <span class="selected-text-preview">{{ textPreview }}</span>
    </div>
    <div class="actions-list">
      <el-button
        v-for="action in actions"
        :key="action.id"
        size="small"
        :type="action.type"
        :icon="action.icon"
        @click="onAction(action.id)"
      >
        {{ action.name }}
      </el-button>
    </div>
    <div class="custom-action">
      <el-input
        v-model="customPrompt"
        size="small"
        placeholder="输入自定义指令，如：用鲁迅风格改写"
        @keyup.enter="onCustomAction"
      >
        <template #append>
          <el-button :icon="ArrowRight" @click="onCustomAction" />
        </template>
      </el-input>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { 
  Brush, 
  ChatDotRound, 
  DataLine, 
  RefreshRight, 
  Plus,
  ArrowRight 
} from '@element-plus/icons-vue'

const props = defineProps<{
  text: string
  position: { x: number; y: number }
}>()

const emit = defineEmits<{
  action: [action: string, prompt?: string]
  close: []
}>()

const customPrompt = ref('')

const textPreview = computed(() => {
  if (props.text.length > 20) {
    return props.text.substring(0, 20) + '...'
  }
  return props.text
})

const actions = [
  { id: 'polish', name: '润色', type: 'primary' as const, icon: Brush },
  { id: 'casual', name: '口语化', type: 'default' as const, icon: ChatDotRound },
  { id: 'data', name: '加数据', type: 'default' as const, icon: DataLine },
  { id: 'rewrite', name: '重写', type: 'default' as const, icon: RefreshRight },
  { id: 'continue', name: '续写', type: 'success' as const, icon: Plus },
]

function onAction(actionId: string) {
  // 延迟关闭，确保事件被正确处理
  emit('action', actionId)
  setTimeout(() => emit('close'), 100)
}

function onCustomAction() {
  if (customPrompt.value.trim()) {
    emit('action', 'custom', customPrompt.value)
    emit('close')
  }
}
</script>

<style scoped>
.floating-ai-actions {
  position: absolute;
  transform: translateX(-50%);
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  padding: 12px;
  min-width: 280px;
}

.actions-header {
  padding-bottom: 8px;
  margin-bottom: 8px;
  border-bottom: 1px solid #e8e8e8;
}

.selected-text-preview {
  font-size: 12px;
  color: #8c8c8c;
  font-style: italic;
}

.actions-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}

.actions-list .el-button {
  flex: 1;
  min-width: 70px;
}

.custom-action {
  padding-top: 8px;
  border-top: 1px solid #e8e8e8;
}
</style>
