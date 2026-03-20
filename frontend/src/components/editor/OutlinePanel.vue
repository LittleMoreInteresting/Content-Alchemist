<template>
  <div class="outline-panel">
    <div class="panel-header">
      <h3>
        <el-icon><List /></el-icon>
        文章大纲
      </h3>
      <el-button 
        type="primary" 
        size="small" 
        :icon="Star"
        @click="generateOutline"
        :loading="isGenerating"
      >
        生成
      </el-button>
    </div>
    
    <div class="outline-tree">
      <div
        v-for="node in outline"
        :key="node.id"
        class="outline-node"
        :class="{
          'level-1': node.level === 1,
          'level-2': node.level === 2,
          'level-3': node.level === 3,
          'active': currentOutlineId === node.id
        }"
        @click="onNodeClick(node)"
      >
        <div class="node-content">
          <span class="node-title">{{ node.title }}</span>
          <span class="node-status" :class="node.status">
            {{ statusText(node.status) }}
          </span>
        </div>
        <div v-if="node.targetWords > 0" class="word-progress">
          <el-progress 
            :percentage="(node.wordCount / node.targetWords) * 100" 
            :show-text="false"
            :stroke-width="4"
          />
          <span class="word-count">{{ node.wordCount }}/{{ node.targetWords }}字</span>
        </div>
      </div>
    </div>
    
    <div class="outline-actions">
      <el-button size="small" :icon="Refresh" @click="refreshOutline">
        刷新大纲
      </el-button>
      <el-button size="small" :icon="Plus" @click="addNode">
        添加节点
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { List, Star, Refresh, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useEditorStore, useConfigStore } from '@/stores'
import { generateOutlineFromAI } from '@/api'
import type { OutlineNode } from '@/types'

const editorStore = useEditorStore()
const configStore = useConfigStore()
const outline = computed(() => editorStore.article.outline || [])
const currentOutlineId = computed(() => editorStore.currentOutlineId)
const isGenerating = computed(() => editorStore.isAIGenerating)

function statusText(status?: string) {
  const map: Record<string, string> = {
    'empty': '未写',
    'draft': '草稿',
    'done': '完成'
  }
  return map[status || 'empty'] || '未写'
}

function onNodeClick(node: OutlineNode) {
  editorStore.setCurrentOutlineId(node.id)
  // 滚动到对应位置
  const headings = document.querySelectorAll('.editor-content h1, .editor-content h2, .editor-content h3')
  const index = outline.value.findIndex(n => n.id === node.id)
  if (headings[index]) {
    headings[index].scrollIntoView({ behavior: 'smooth' })
  }
}

async function generateOutline() {
  if (!editorStore.article.title) {
    ElMessage.warning('请先输入文章标题')
    return
  }
  
  editorStore.setAIGenerating(true)
  
  try {
    const nodes = await generateOutlineFromAI(
      editorStore.article.title,
      configStore.styleDescription,
      configStore.config.audience
    )
    
    // 更新大纲
    editorStore.updateOutline(nodes as any)
    
    // 将大纲同步到编辑器内容
    syncOutlineToEditor(nodes as OutlineNode[])
    
    ElMessage.success('大纲生成成功')
  } catch (error) {
    ElMessage.error('生成大纲失败：' + error)
  } finally {
    editorStore.setAIGenerating(false)
  }
}

// 将大纲同步到编辑器内容
function syncOutlineToEditor(nodes: OutlineNode[]) {
  if (nodes.length === 0) return
  
  // 构建 Markdown 大纲
  let markdownContent = `# ${editorStore.article.title}\n\n`
  
  for (const node of nodes) {
    const prefix = '#'.repeat(node.level)
    markdownContent += `${prefix} ${node.title}\n\n`
  }
  
  // 保留原有内容（除了标题）
  const existingContent = editorStore.article.content
  const bodyMatch = existingContent.match(/# .+[\r\n]+([\s\S]*)/)
  if (bodyMatch && bodyMatch[1].trim()) {
    markdownContent += bodyMatch[1].trim()
  } else {
    markdownContent += '开始创作...'
  }
  
  editorStore.updateContent(markdownContent)
}

function refreshOutline() {
  // 从内容重新解析大纲并更新
  const newOutline = editorStore.outlineFromContent
  if (newOutline.length > 0) {
    editorStore.updateOutline(newOutline)
    ElMessage.success('大纲已刷新')
  } else {
    ElMessage.info('未检测到标题')
  }
}

function addNode() {
  // 添加一个新的 H2 节点
  const newNode: OutlineNode = {
    id: Date.now().toString(),
    level: 2,
    title: '新章节',
    status: 'empty',
    wordCount: 0,
    targetWords: 300
  }
  
  // 更新大纲
  const currentOutline = [...editorStore.article.outline]
  currentOutline.push(newNode)
  editorStore.updateOutline(currentOutline)
  
  // 同步到编辑器内容
  const newContent = editorStore.article.content + '\n\n## 新章节\n\n'
  editorStore.updateContent(newContent)
  
  ElMessage.success('已添加新章节')
}
</script>

<style scoped>
.outline-panel {
  display: flex;
  flex-direction: column;
  height: 50%;
  border-bottom: 1px solid #e8e8e8;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #e8e8e8;
}

.panel-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
}

.outline-tree {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.outline-node {
  padding: 8px 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.outline-node:hover {
  background: #f5f5f5;
}

.outline-node.active {
  background: #e6f7ff;
  border-right: 3px solid #1677ff;
}

.level-1 {
  font-weight: 600;
}

.level-2 {
  padding-left: 28px;
  font-size: 13px;
}

.level-3 {
  padding-left: 40px;
  font-size: 12px;
  color: #595959;
}

.node-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.node-title {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.node-status {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
  background: #f5f5f5;
  color: #8c8c8c;
}

.node-status.draft {
  background: #fff7e6;
  color: #faad14;
}

.node-status.done {
  background: #f6ffed;
  color: #52c41a;
}

.word-progress {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
}

.word-progress :deep(.el-progress) {
  flex: 1;
}

.word-count {
  font-size: 10px;
  color: #8c8c8c;
  white-space: nowrap;
}

.outline-actions {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid #e8e8e8;
}

.outline-actions .el-button {
  flex: 1;
}
</style>
