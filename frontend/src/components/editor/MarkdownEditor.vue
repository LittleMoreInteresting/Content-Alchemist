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
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { List, Document, ChatDotSquare } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useEditorStore, useConfigStore } from '@/stores'
import { streamWriteWithAI, saveArticleToBackend, generateOutlineFromAI } from '@/api'
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
    // 计算位置（带边界检查）
    const selectionObj = window.getSelection()
    if (selectionObj && selectionObj.rangeCount > 0) {
      const range = selectionObj.getRangeAt(0)
      const rect = range.getBoundingClientRect()
      const editorRect = textarea.getBoundingClientRect()
      
      // 工具条尺寸（估算）
      const toolbarWidth = 320
      const toolbarHeight = 120
      
      // 计算理想位置（在选中文本上方居中）
      let x = rect.left + rect.width / 2
      let y = rect.top - 10
      
      // 水平边界检查：确保不超出编辑器左右边界
      const minX = editorRect.left + toolbarWidth / 2 + 10
      const maxX = editorRect.right - toolbarWidth / 2 - 10
      x = Math.max(minX, Math.min(x, maxX))
      
      // 垂直边界检查：如果上方空间不足，则显示在下方
      if (y < toolbarHeight + 10) {
        y = rect.bottom + toolbarHeight + 10
      }
      
      // 确保不超出视口底部
      const viewportHeight = window.innerHeight
      if (y > viewportHeight - 10) {
        y = viewportHeight - toolbarHeight - 10
      }
      
      floatingAIPosition.value = { x, y }
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
    // 计算菜单位置（带边界检查）
    const lineHeight = 24
    const lineIndex = lines.length - 1
    const cursorLineY = (lineIndex + 1) * lineHeight + 10
    
    // 菜单尺寸（估算）
    const menuHeight = 280
    
    // 获取编辑器容器信息
    const editorRect = textarea.getBoundingClientRect()
    
    // 默认位置：在光标下方
    let x = 20
    let y = cursorLineY + lineHeight
    
    // 如果靠近底部，显示在光标上方
    const viewportHeight = window.innerHeight
    const spaceBelow = viewportHeight - editorRect.top - y
    if (spaceBelow < menuHeight + 20 && editorRect.top + cursorLineY > menuHeight + 20) {
      y = cursorLineY - menuHeight - 10
    }
    
    // 确保 Y 不超出编辑器上方
    y = Math.max(10, y)
    
    slashMenuPosition.value = { x, y }
  } else {
    showSlashMenu.value = false
  }
}

async function onSlashSelect(command: string) {
  const textarea = editorRef.value
  if (!textarea) return
  
  const cursorPos = textarea.selectionStart
  const textBefore = content.value.substring(0, cursorPos)
  const textAfter = content.value.substring(cursorPos)
  
  // 删除 / 命令行
  const lines = textBefore.split('\n')
  lines.pop() // 删除当前命令行
  const newTextBefore = lines.join('\n')
  content.value = newTextBefore + textAfter
  
  showSlashMenu.value = false
  
  // 执行对应的 AI 命令
  const configStore = useConfigStore()
  const article = editorStore.article
  
  try {
    editorStore.setAIGenerating(true)
    
    switch (command) {
      case 'ai':
        // 召唤 AI 助手 - 显示浮动工具条
        // 获取光标所在位置的文本
        const currentPara = getCurrentParagraph(cursorPos)
        if (currentPara.text.trim()) {
          editorStore.setSelectedText(currentPara.text)
          selectionRange.value = { start: currentPara.start, end: currentPara.end }
          // 计算浮动工具条位置
          const lineHeight = 24
          const lines = newTextBefore.split('\n').length
          floatingAIPosition.value = {
            x: 100,
            y: lines * lineHeight + 50
          }
          showFloatingAI.value = true
          ElMessage.info('已选中当前段落，请选择 AI 操作')
        } else {
          ElMessage.warning('请先输入一些文本')
        }
        break
        
      case 'outline':
        // 基于标题生成大纲
        if (!article.title) {
          ElMessage.warning('请先设置文章标题')
          return
        }
        await generateOutlineFromAIAndSync(article.title, configStore.styleDescription, configStore.config.audience)
        break
        
      case 'expand':
      case 'polish':
      case 'shorter':
      case 'continue':
        // 对当前段落进行 AI 处理
        await processCurrentParagraph(command, cursorPos, configStore.styleDescription)
        break
        
      case 'title':
        // 生成标题建议
        await generateTitleSuggestions(configStore.styleDescription)
        break
        
      default:
        ElMessage.warning(`未知命令: ${command}`)
    }
  } catch (error) {
    console.error('AI 命令执行失败:', error)
    ElMessage.error('AI 处理失败：' + error)
  } finally {
    editorStore.setAIGenerating(false)
  }
}

// 获取当前段落信息
function getCurrentParagraph(cursorPos: number): { text: string, start: number, end: number } {
  const text = content.value
  // 找到段落开始（上一个换行符或文本开始）
  let start = cursorPos
  while (start > 0 && text[start - 1] !== '\n') {
    start--
  }
  // 找到段落结束（下一个换行符或文本结束）
  let end = cursorPos
  while (end < text.length && text[end] !== '\n') {
    end++
  }
  return { text: text.substring(start, end), start, end }
}

// 处理当前段落
async function processCurrentParagraph(action: string, cursorPos: number, style: string) {
  const para = getCurrentParagraph(cursorPos)
  if (!para.text.trim()) {
    ElMessage.warning('当前段落为空')
    return
  }
  
  // 构建上下文（前后各200字符）
  const contextStart = Math.max(0, para.start - 200)
  const contextEnd = Math.min(content.value.length, para.end + 200)
  const context = content.value.substring(contextStart, contextEnd)
  
  // 根据 action 类型确定操作
  const actionMap: Record<string, string> = {
    'expand': 'expand',
    'polish': 'polish', 
    'shorter': 'shorten',
    'continue': 'continue'
  }
  
  const result = await streamWriteWithAI(
    actionMap[action] || action,
    context,
    para.text,
    'replace',
    style,
    ''
  )
  
  // 替换段落
  content.value = content.value.substring(0, para.start) + result + content.value.substring(para.end)
  ElMessage.success('AI处理完成')
}

// 生成标题建议
async function generateTitleSuggestions(style: string) {
  const article = editorStore.article
  if (!article.content.trim() && !article.title) {
    ElMessage.warning('请先输入文章内容')
    return
  }
  
  const context = article.content.substring(0, 1000) // 取前1000字符
  const result = await streamWriteWithAI(
    'title',
    context,
    article.title || '',
    'replace',
    style,
    ''
  )
  
  // 显示标题建议对话框
  showTitleSuggestions(result)
}

// 显示标题建议
function showTitleSuggestions(result: string) {
  // 解析标题（按行分割，过滤空行）
  const titles = result.split('\n').filter(t => t.trim()).map(t => t.replace(/^\d+[.、\s]*/, '').trim())
  
  if (titles.length === 0) {
    ElMessage.warning('未生成标题建议')
    return
  }
  
  // 存储标题建议到 store，供其他组件显示
  editorStore.setTitleSuggestions(titles)
  ElMessage.success(`已生成 ${titles.length} 个标题建议，请查看大纲面板`)
}

// 生成大纲并同步到编辑器
async function generateOutlineFromAIAndSync(title: string, style: string, audience: string) {
  ElMessage.info('正在生成大纲...')
  const nodes = await generateOutlineFromAI(title, style, audience)
  editorStore.updateOutline(nodes as any)
  
  // 同步到编辑器
  let markdownContent = `# ${title}\n\n`
  for (const node of nodes) {
    const prefix = '#'.repeat(node.level)
    markdownContent += `${prefix} ${node.title}\n\n`
  }
  
  // 保留原有正文内容
  const existingBody = content.value.replace(/^# .+[\r\n]+/, '').trim()
  if (existingBody) {
    markdownContent += existingBody
  }
  
  content.value = markdownContent
  ElMessage.success('大纲生成成功')
}

async function onAIAction(action: string, customPrompt?: string) {
  const configStore = useConfigStore()
  editorStore.setAIGenerating(true)
  showFloatingAI.value = false
  
  try {
    // 构建上下文（前后各200字符）
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

// 监听 AI 助手按钮点击事件
function handleShowAISlashMenu() {
  // 获取编辑器焦点
  const textarea = editorRef.value
  if (textarea) {
    textarea.focus()
    // 在当前光标位置显示 slash 菜单
    const cursorPos = textarea.selectionStart || 0
    const textBefore = content.value.substring(0, cursorPos)
    const lines = textBefore.split('\n')
    const currentLine = lines[lines.length - 1]
    const lineIndex = lines.length - 1
    
    // 计算位置（带边界检查）
    const lineHeight = 24
    const menuHeight = 280
    const editorRect = textarea.getBoundingClientRect()
    const viewportHeight = window.innerHeight
    
    // 默认位置
    let x = 20
    let y = (lineIndex + 1) * lineHeight + 10 + lineHeight
    
    // 垂直边界检查
    const spaceBelow = viewportHeight - editorRect.top - y
    if (spaceBelow < menuHeight + 20 && editorRect.top + (lineIndex * lineHeight) > menuHeight + 20) {
      y = lineIndex * lineHeight - menuHeight - 10
    }
    y = Math.max(10, y)
    
    // 如果当前行没有 /，则在行首插入
    if (!currentLine.trim().startsWith('/')) {
      // 在行首插入 / 并显示菜单
      const lineStart = textBefore.lastIndexOf('\n') + 1
      content.value = content.value.substring(0, lineStart) + '/ ' + content.value.substring(lineStart)
      
      // 更新光标位置
      setTimeout(() => {
        const newPos = lineStart + 2
        textarea.setSelectionRange(newPos, newPos)
        // 显示菜单
        showSlashMenu.value = true
        slashMenuPosition.value = { x, y }
      }, 0)
    } else {
      // 已经输入了 /，直接显示菜单
      showSlashMenu.value = true
      slashMenuPosition.value = { x, y }
    }
  }
}

onMounted(() => {
  window.addEventListener('show-ai-slash-menu', handleShowAISlashMenu)
})

onUnmounted(() => {
  window.removeEventListener('show-ai-slash-menu', handleShowAISlashMenu)
})
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
