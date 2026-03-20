<template>
  <div 
    class="slash-command-menu" 
    :style="{ top: position.y + 'px', left: position.x + 'px' }"
  >
    <div class="menu-header">AI 命令</div>
    <div 
      v-for="(cmd, index) in commands" 
      :key="cmd.id"
      class="command-item"
      :class="{ active: selectedIndex === index }"
      @click="selectCommand(cmd)"
      @mouseenter="selectedIndex = index"
    >
      <el-icon class="command-icon">
        <component :is="cmd.icon" />
      </el-icon>
      <div class="command-info">
        <span class="command-name">{{ cmd.name }}</span>
        <span class="command-desc">{{ cmd.desc }}</span>
      </div>
      <span v-if="cmd.shortcut" class="command-shortcut">{{ cmd.shortcut }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { 
  Star, 
  List, 
  Expand, 
  Brush, 
  Scissor, 
  TopRight,
  ChatDotRound
} from '@element-plus/icons-vue'

interface Command {
  id: string
  name: string
  desc: string
  icon: any
  shortcut?: string
}

defineProps<{
  position: { x: number; y: number }
}>()

const emit = defineEmits<{
  select: [command: string]
  close: []
}>()

const commands: Command[] = [
  { id: 'ai', name: '/ai', desc: '召唤AI助手', icon: Star, shortcut: '⌘K' },
  { id: 'outline', name: '/outline', desc: '基于标题生成大纲', icon: List },
  { id: 'expand', name: '/expand', desc: '扩写当前段落', icon: Expand },
  { id: 'polish', name: '/polish', desc: '润色当前段落', icon: Brush },
  { id: 'shorter', name: '/shorter', desc: '精简内容', icon: Scissor },
  { id: 'title', name: '/title', desc: '生成5个标题建议', icon: ChatDotRound },
  { id: 'continue', name: '/continue', desc: '续写内容', icon: TopRight },
]

const selectedIndex = ref(0)

function selectCommand(cmd: Command) {
  emit('select', cmd.id)
}

function handleKeydown(e: KeyboardEvent) {
  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault()
      selectedIndex.value = (selectedIndex.value + 1) % commands.length
      break
    case 'ArrowUp':
      e.preventDefault()
      selectedIndex.value = (selectedIndex.value - 1 + commands.length) % commands.length
      break
    case 'Enter':
      e.preventDefault()
      selectCommand(commands[selectedIndex.value])
      break
    case 'Escape':
      emit('close')
      break
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.slash-command-menu {
  position: absolute;
  width: 280px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  overflow: hidden;
}

.menu-header {
  padding: 8px 12px;
  font-size: 12px;
  color: #8c8c8c;
  background: #f5f5f5;
  border-bottom: 1px solid #e8e8e8;
}

.command-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.command-item:hover,
.command-item.active {
  background: #f0f7ff;
}

.command-icon {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #1677ff;
  font-size: 16px;
}

.command-info {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.command-name {
  font-size: 14px;
  font-weight: 500;
  color: #262626;
}

.command-desc {
  font-size: 12px;
  color: #8c8c8c;
}

.command-shortcut {
  font-size: 11px;
  color: #bfbfbf;
  padding: 2px 6px;
  background: #f5f5f5;
  border-radius: 4px;
}
</style>
