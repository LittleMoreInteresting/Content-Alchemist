<template>
  <div class="markdown-editor">
    <div class="editor-toolbar">
      <el-button-group>
        <el-button size="small" @click="insertFormat('bold')">
          <b>B</b>
        </el-button>
        <el-button size="small" @click="insertFormat('italic')">
          <i>I</i>
        </el-button>
        <el-button size="small" @click="insertFormat('heading')">
          H
        </el-button>
        <el-button size="small" @click="insertFormat('list')">
          <el-icon><List /></el-icon>
        </el-button>
        <el-button size="small" @click="insertFormat('code')">
          <el-icon><Document /></el-icon>
        </el-button>
        <el-button size="small" @click="insertFormat('quote')">
          <el-icon><ChatDotSquare /></el-icon>
        </el-button>
      </el-button-group>
      
      <span class="word-count">{{ wordCount }} 字</span>
    </div>
    
    <textarea
      ref="editorRef"
      v-model="content"
      class="editor-content"
      placeholder="开始创作...&#10;支持 Markdown 语法&#10;&#10;输入 / 打开 AI 命令面板"
      @input="onInput"
      @keydown="onKeyDown"
      @mouseup="onMouseUp"
      @scroll="onScroll"
    />
    
    <!-- AI 斜杠命令菜单 -->
    <SlashCommand
      v-if="showSlashMenu"
      :position="slashMenuPosition"
      @select="onSlashSelect"
      @close="showSlashMenu = false"
    />
    
    <!-- AI 浮动工具条 -->
    <FloatingAIActions
      v-if="showFloatingAI && selectedText"
      :text="selectedText"
      :position="floatingAIPosition"
      @action="onAIAction"
      @close="showFloatingAI = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { List, Document, ChatDotSquare } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useEditorStore, useConfigStore } from '@/stores'
import { streamWriteWithAI, saveArticleToBackend } from '@/api'
import SlashCommand from './SlashCommand.vue'
import FloatingAIActions from './FloatingAIActions.vue'

const editorStore = useEditorStore()
const editorRef = ref<HTMLTextAreaElement>()

const content = computed({
  get: () => editorStore.article.content,
  set: (val) => editorStore.updateContent(val)
})

const wordCount = computed(() => editorStore.wordCount)
const selectedText = computed(() => editorStore.selectedText)

// 斜杠命令菜单
const showSlashMenu = ref(false)
const slashMenuPosition = ref({ x: 0, y: 0 })

// AI浮动工具条
const showFloatingAI = ref(false)
const floatingAIPosition = ref({ x: 0, y: 0 })
const selectionRange = ref({ start: 0, end: 0 })

let inputDebounceTimer: number | null = null

function onInput() {
  // 防抖自动保存
  if (inputDebounceTimer) {
    clearTimeout(inputDebounceTimer)
  }
  inputDebounceTimer = window.setTimeout(() => {
    autoSave()
  }, 3000)
  
  // 检查是否输入了 /
  checkSlashCommand()
}

function onKeyDown(e: KeyboardEvent) {
  if (e.key === '/') {
    // 延迟检查，确保 / 已输入
    setTimeout(() => checkSlashCommand(), 0)
  } else if (e.key === 'Escape') {
    showSlashMenu.value = false
    showFloatingAI.value = false
  }
}

function onMouseUp() {
  const textarea = editorRef.value
  if (!textarea) return
  
  const selection = window.getSelection()?.toString()
  if (selection && selection.trim()) {
    editorStore.setSelectedText(selection)
    // 记录选中文本的位置
    selectionRange.value = {
      start: textarea.selectionStart,
      end: textarea.selectionEnd
    }
    showFloatingAI.value = true
    // 计算位置
    const selectionObj = window.getSelection()
    if (selectionObj && selectionObj.rangeCount > 0) {
      const range = selectionObj.getRangeAt(0)
      const rect = range.getBoundingClientRect()
      floatingAIPosition.value = {
        x: rect.left + rect.width / 2,
        y: rect.top - 50
      }
    }
  } else {
    showFloatingAI.value = false
  }
}

function onScroll() {
  showSlashMenu.value = false
  showFloatingAI.value = false
}

function checkSlashCommand() {
  const textarea = editorRef.value
  if (!textarea) return
  
  const cursorPos = textarea.selectionStart
  const textBeforeCursor = content.value.substring(0, cursorPos)
  const lines = textBeforeCursor.split('\n')
  const currentLine = lines[lines.length - 1]
  
  // 检查当前行是否以 / 开头
  if (currentLine.trim() === '/' || currentLine.trim().startsWith('/')) {
    showSlashMenu.value = true
    // 计算菜单位置
    const lineHeight = 24
    const lineIndex = lines.length - 1
    slashMenuPosition.value = {
      x: 20,
      y: (lineIndex + 1) * lineHeight + 10
    }
  } else {
    showSlashMenu.value = false
  }
}

function onSlashSelect(command: string) {
  const textarea = editorRef.value
  if (!textarea) return
  
  const cursorPos = textarea.selectionStart
  const textBefore = content.value.substring(0, cursorPos)
  const textAfter = content.value.substring(cursorPos)
  
  // 替换 / 命令为对应内容
  const lines = textBefore.split('\n')
  lines[lines.length - 1] = `<!-- AI: ${command} -->`
  content.value = lines.join('\n') + textAfter
  
  showSlashMenu.value = false
}

async function onAIAction(action: string, customPrompt?: string) {
  const configStore = useConfigStore()
  editorStore.setAIGenerating(true)
  showFloatingAI.value = false
  
  try {
    // 构建上下文（前后各200字符）
    const textarea = editorRef.value
    const cursorPos = selectionRange.value.start
    const contextStart = Math.max(0, cursorPos - 200)
    const contextEnd = Math.min(content.value.length, cursorPos + 200)
    const context = content.value.substring(contextStart, contextEnd)
    
    console.log('AI Action:', action, 'customPrompt:', customPrompt)
    
    // 调用后端AI接口
    const result = await streamWriteWithAI(
      action,
      context,
      selectedText.value,
      'replace',
      configStore.styleDescription,
      customPrompt
    )
    
    // 使用记录的选中文本位置进行替换
    const start = selectionRange.value.start
    const end = selectionRange.value.end
    
    // 验证选中的文本是否匹配
    const actualSelected = content.value.substring(start, end)
    if (actualSelected !== selectedText.value) {
      console.warn('选中文本已变化，使用indexOf回退')
      const idx = content.value.indexOf(selectedText.value)
      if (idx === -1) {
        throw new Error('无法找到要替换的文本')
      }
      content.value = content.value.substring(0, idx) + result + content.value.substring(idx + selectedText.value.length)
    } else {
      content.value = content.value.substring(0, start) + result + content.value.substring(end)
    }
    
    ElMessage.success('AI处理完成')
  } catch (error) {
    console.error('AI处理失败:', error)
    ElMessage.error('AI处理失败：' + error)
  } finally {
    editorStore.setAIGenerating(false)
    // 清空选择
    editorStore.setSelectedText('')
    selectionRange.value = { start: 0, end: 0 }
  }
}

function insertFormat(format: string) {
  const textarea = editorRef.value
  if (!textarea) return
  
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selected = content.value.substring(start, end)
  
  let replacement = selected
  switch (format) {
    case 'bold':
      replacement = `**${selected || '粗体文字'}**`
      break
    case 'italic':
      replacement = `*${selected || '斜体文字'}*`
      break
    case 'heading':
      replacement = `## ${selected || '标题'}`
      break
    case 'list':
      replacement = `\n- ${selected || '列表项'}\n- 列表项\n- 列表项`
      break
    case 'code':
      replacement = `\`\`\`\n${selected || '代码块'}\n\`\`\``
      break
    case 'quote':
      replacement = `> ${selected || '引用文字'}`
      break
  }
  
  content.value = content.value.substring(0, start) + replacement + content.value.substring(end)
  
  // 聚焦并选中新插入的内容
  setTimeout(() => {
    textarea.focus()
    textarea.setSelectionRange(start + replacement.length, start + replacement.length)
  }, 0)
}

async function autoSave() {
  editorStore.markSaving()
  
  try {
    // 调用后端保存接口
    await saveArticleToBackend(editorStore.article)
    editorStore.markSaved()
  } catch (error) {
    editorStore.markUnsaved()
    console.error('自动保存失败:', error)
  }
}
</script>

<style scoped>
.markdown-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  border-bottom: 1px solid #e8e8e8;
  background: #fafafa;
}

.word-count {
  font-size: 12px;
  color: #8c8c8c;
}

.editor-content {
  flex: 1;
  width: 100%;
  padding: 24px;
  border: none;
  outline: none;
  resize: none;
  font-family: 'JetBrains Mono', 'Fira Code', Consolas, monospace;
  font-size: 15px;
  line-height: 1.8;
  color: #262626;
  background: #fff;
}

.editor-content::placeholder {
  color: #bfbfbf;
}
</style>
