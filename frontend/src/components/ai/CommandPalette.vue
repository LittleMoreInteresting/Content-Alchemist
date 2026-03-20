<template>
  <Teleport to="body">
    <div v-if="visible" class="command-palette-overlay" @click="close">
      <div class="command-palette" @click.stop>
        <div class="search-box">
          <el-icon><Search /></el-icon>
          <input
            v-model="searchQuery"
            ref="inputRef"
            type="text"
            placeholder="输入命令或搜索..."
            @keydown="onKeyDown"
          />
          <span class="shortcut-hint">ESC 关闭</span>
        </div>
        
        <div class="command-list" v-if="filteredCommands.length > 0">
          <div
            v-for="(cmd, index) in filteredCommands"
            :key="cmd.id"
            class="command-item"
            :class="{ active: selectedIndex === index }"
            @click="executeCommand(cmd)"
            @mouseenter="selectedIndex = index"
          >
            <el-icon class="command-icon">
              <component :is="cmd.icon" />
            </el-icon>
            <div class="command-info">
              <span class="command-name">{{ cmd.name }}</span>
              <span class="command-desc">{{ cmd.description }}</span>
            </div>
            <span v-if="cmd.shortcut" class="command-shortcut">{{ cmd.shortcut }}</span>
          </div>
        </div>
        
        <div v-else class="no-results">
          <el-icon><InfoFilled /></el-icon>
          <span>未找到匹配的命令</span>
        </div>
        
        <div class="command-footer">
          <div class="footer-hint">
            <kbd>↑</kbd> <kbd>↓</kbd> 选择
            <kbd>Enter</kbd> 执行
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, watch, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search, InfoFilled } from '@element-plus/icons-vue'
import { useEditorStore } from '@/stores'
import type { Command } from '@/types'

const router = useRouter()
const editorStore = useEditorStore()

const visible = ref(false)
const searchQuery = ref('')
const selectedIndex = ref(0)
const inputRef = ref<HTMLInputElement>()

// 所有可用命令
const allCommands = computed<Command[]>(() => [
  // 文件操作
  { id: 'new', name: '新建文章', description: '创建一篇新文章', icon: 'DocumentAdd', shortcut: '⌘N', action: () => { editorStore.reset(); close() } },
  { id: 'save', name: '保存文章', description: '手动保存当前文章', icon: 'Upload', shortcut: '⌘S', action: () => { editorStore.markSaved(); close() } },
  
  // AI操作
  { id: 'ai-outline', name: '生成大纲', description: '基于标题AI生成文章大纲', icon: 'List', shortcut: '', action: () => { close() } },
  { id: 'ai-expand', name: '扩写内容', description: '扩写当前段落', icon: 'Expand', shortcut: '', action: () => { close() } },
  { id: 'ai-polish', name: '润色文本', description: '润色当前段落', icon: 'MagicStick', shortcut: '', action: () => { close() } },
  { id: 'ai-title', name: '生成标题', description: '生成5个标题建议', icon: 'ChatDotRound', shortcut: '', action: () => { close() } },
  
  // 视图切换
  { id: 'view-edit', name: '编辑模式', description: '切换到纯编辑模式', icon: 'EditPen', shortcut: '', action: () => { editorStore.setEditorView('edit'); close() } },
  { id: 'view-split', name: '分屏模式', description: '切换到分屏预览模式', icon: 'Monitor', shortcut: '', action: () => { editorStore.setEditorView('split'); close() } },
  { id: 'view-preview', name: '预览模式', description: '切换到纯预览模式', icon: 'View', shortcut: '', action: () => { editorStore.setEditorView('preview'); close() } },
  
  // 主题切换
  { id: 'theme-default', name: '科技蓝主题', description: '切换到默认科技蓝主题', icon: 'Brush', shortcut: '', action: () => { editorStore.setTheme('default'); close() } },
  { id: 'theme-minimal', name: '简约白主题', description: '切换到简约白主题', icon: 'Brush', shortcut: '', action: () => { editorStore.setTheme('minimal'); close() } },
  { id: 'theme-vibrant', name: '活力橙主题', description: '切换到活力橙主题', icon: 'Brush', shortcut: '', action: () => { editorStore.setTheme('vibrant'); close() } },
  
  // 导航
  { id: 'go-settings', name: '打开设置', description: '进入设置页面', icon: 'Setting', shortcut: '', action: () => { router.push('/settings'); close() } },
  { id: 'go-publish', name: '发布文章', description: '打开发布助手', icon: 'Position', shortcut: '⌘⇧P', action: () => { router.push('/publish'); close() } },
  
  // 帮助
  { id: 'help-shortcuts', name: '快捷键帮助', description: '查看所有快捷键', icon: 'QuestionFilled', shortcut: '', action: () => { close() } },
])

// 过滤后的命令
const filteredCommands = computed(() => {
  const query = searchQuery.value.toLowerCase().trim()
  if (!query) return allCommands.value
  
  return allCommands.value.filter(cmd => 
    cmd.name.toLowerCase().includes(query) ||
    cmd.description.toLowerCase().includes(query)
  )
})

// 监听过滤结果变化，重置选中索引
watch(filteredCommands, () => {
  selectedIndex.value = 0
})

function open() {
  visible.value = true
  searchQuery.value = ''
  selectedIndex.value = 0
  nextTick(() => {
    inputRef.value?.focus()
  })
}

function close() {
  visible.value = false
}

function executeCommand(cmd: Command) {
  cmd.action()
}

function onKeyDown(e: KeyboardEvent) {
  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault()
      selectedIndex.value = (selectedIndex.value + 1) % filteredCommands.value.length
      break
    case 'ArrowUp':
      e.preventDefault()
      selectedIndex.value = (selectedIndex.value - 1 + filteredCommands.value.length) % filteredCommands.value.length
      break
    case 'Enter':
      e.preventDefault()
      if (filteredCommands.value[selectedIndex.value]) {
        executeCommand(filteredCommands.value[selectedIndex.value])
      }
      break
    case 'Escape':
      e.preventDefault()
      close()
      break
  }
}

// 全局快捷键监听
function handleGlobalKeyDown(e: KeyboardEvent) {
  // Cmd/Ctrl + K 打开命令面板
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
    e.preventDefault()
    open()
  }
  
  // Cmd/Ctrl + N 新建文章
  if ((e.metaKey || e.ctrlKey) && e.key === 'n') {
    e.preventDefault()
    editorStore.reset()
  }
  
  // Cmd/Ctrl + S 保存
  if ((e.metaKey || e.ctrlKey) && e.key === 's') {
    e.preventDefault()
    editorStore.markSaved()
  }
  
  // Cmd/Ctrl + Shift + P 发布
  if ((e.metaKey || e.ctrlKey) && e.shiftKey && e.key === 'P') {
    e.preventDefault()
    router.push('/publish')
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleGlobalKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleGlobalKeyDown)
})

// 暴露打开方法给父组件
defineExpose({ open, close })
</script>

<style scoped>
.command-palette-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 100px;
  z-index: 9999;
}

.command-palette {
  width: 600px;
  max-width: 90vw;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  overflow: hidden;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid #e8e8e8;
}

.search-box .el-icon {
  font-size: 20px;
  color: #8c8c8c;
}

.search-box input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 16px;
  background: transparent;
}

.search-box input::placeholder {
  color: #bfbfbf;
}

.shortcut-hint {
  font-size: 12px;
  color: #8c8c8c;
  background: #f5f5f5;
  padding: 4px 8px;
  border-radius: 4px;
}

.command-list {
  max-height: 400px;
  overflow-y: auto;
}

.command-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  cursor: pointer;
  transition: background 0.15s;
}

.command-item:hover,
.command-item.active {
  background: #f0f7ff;
}

.command-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  border-radius: 6px;
  color: #1677ff;
  font-size: 16px;
}

.command-info {
  flex: 1;
  display: flex;
  flex-direction: column;
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
  color: #8c8c8c;
  background: #f5f5f5;
  padding: 4px 8px;
  border-radius: 4px;
  font-family: 'JetBrains Mono', monospace;
}

.no-results {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 40px;
  color: #8c8c8c;
}

.no-results .el-icon {
  font-size: 32px;
}

.command-footer {
  padding: 12px 20px;
  background: #fafafa;
  border-top: 1px solid #e8e8e8;
}

.footer-hint {
  font-size: 12px;
  color: #8c8c8c;
  display: flex;
  align-items: center;
  gap: 8px;
}

.footer-hint kbd {
  padding: 2px 6px;
  background: #fff;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px;
}
</style>
