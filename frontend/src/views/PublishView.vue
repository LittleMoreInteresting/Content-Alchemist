<template>
  <div class="publish-view">
    <div class="publish-header">
      <el-button :icon="ArrowLeft" @click="goBack">返回</el-button>
      <h1>发布助手</h1>
      <span></span>
    </div>
    
    <div class="publish-content">
      <div class="checklist-section">
        <h2>📋 发布前检查清单</h2>
        <el-card class="checklist-card">
          <div 
            v-for="item in checklist" 
            :key="item.id"
            class="check-item"
            :class="{ checked: item.checked }"
          >
            <el-checkbox v-model="item.checked" size="large">
              <span class="check-label">{{ item.label }}</span>
              <span class="check-status" :class="item.status">{{ statusText(item) }}</span>
            </el-checkbox>
          </div>
        </el-card>
      </div>
      
      <div class="stats-section">
        <h2>📊 文章统计</h2>
        <el-card>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="标题字数">
              {{ articleTitle.length }} 字
            </el-descriptions-item>
            <el-descriptions-item label="正文字数">
              {{ wordCount }} 字
            </el-descriptions-item>
            <el-descriptions-item label="阅读时间">
              {{ Math.ceil(wordCount / 300) }} 分钟
            </el-descriptions-item>
            <el-descriptions-item label="大纲节点">
              {{ outlineCount }} 个
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </div>
      
      <div class="actions-section">
        <h2>🚀 发布操作</h2>
        <el-card>
          <div class="action-buttons">
            <el-button type="primary" size="large" :icon="CopyDocument" @click="copyHTML">
              复制带样式HTML
            </el-button>
            <el-button size="large" :icon="Link" @click="openWechat">
              打开公众号后台
            </el-button>
            <el-button size="large" :icon="Clock" @click="setReminder">
              设置定时提醒
            </el-button>
          </div>
          <el-divider />
          <div class="export-buttons">
            <el-button :icon="Document" @click="downloadMarkdown">下载Markdown</el-button>
            <el-button :icon="Picture" @click="exportImage">导出长图</el-button>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { 
  ArrowLeft, 
  CopyDocument, 
  Link, 
  Clock, 
  Document, 
  Picture 
} from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useEditorStore } from '@/stores'
import { renderWechatHTML } from '@/utils/wechatRenderer'

const router = useRouter()
const editorStore = useEditorStore()

const articleTitle = computed(() => editorStore.article.title)
const wordCount = computed(() => editorStore.wordCount)
const outlineCount = computed(() => editorStore.outlineFromContent.length)

const checklist = ref([
  { 
    id: 'title', 
    label: '标题字数检查', 
    checked: false, 
    status: computed(() => articleTitle.value.length <= 30 ? 'pass' : 'warn'),
    message: computed(() => articleTitle.value.length <= 30 ? '建议≤30字' : `当前${articleTitle.value.length}字，建议精简`)
  },
  { 
    id: 'cover', 
    label: '封面图检查', 
    checked: false, 
    status: 'info' as const,
    message: '建议尺寸900×383'
  },
  { 
    id: 'words', 
    label: '字数统计', 
    checked: false, 
    status: computed(() => {
      const count = wordCount.value
      if (count < 800) return 'warn'
      if (count > 3000) return 'warn'
      return 'pass'
    }),
    message: computed(() => {
      const count = wordCount.value
      if (count < 800) return `当前${count}字，建议扩充到800-3000字`
      if (count > 3000) return `当前${count}字，建议精简到3000字以内`
      return '建议800-3000字'
    })
  },
  { 
    id: 'preview', 
    label: '手机预览', 
    checked: false, 
    status: 'info' as const,
    message: '建议在手机上预览效果'
  }
])

function statusText(item: any) {
  const statusMap: Record<string, string> = {
    pass: '✅',
    warn: '⚠️',
    info: 'ℹ️'
  }
  return `${statusMap[item.status as string] || ''} ${item.message}`
}

function goBack() {
  router.push('/editor')
}

async function copyHTML() {
  const html = renderWechatHTML(editorStore.article.content)
  await navigator.clipboard.writeText(html)
  ElMessage.success('HTML已复制到剪贴板，可直接粘贴到公众号后台')
}

function openWechat() {
  window.open('https://mp.weixin.qq.com', '_blank')
}

function setReminder() {
  ElMessage.info('定时提醒功能开发中...')
}

function downloadMarkdown() {
  const blob = new Blob([editorStore.article.content], { type: 'text/markdown' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${editorStore.article.title || 'article'}.md`
  a.click()
  URL.revokeObjectURL(url)
  ElMessage.success('Markdown已下载')
}

function exportImage() {
  ElMessage.info('导出长图功能开发中...')
}
</script>

<style scoped>
.publish-view {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}

.publish-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
}

.publish-header h1 {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.publish-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
}

.publish-content h2 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
  color: #262626;
}

.checklist-section,
.stats-section,
.actions-section {
  margin-bottom: 24px;
}

.checklist-card :deep(.el-card__body) {
  padding: 0;
}

.check-item {
  padding: 16px 20px;
  border-bottom: 1px solid #e8e8e8;
}

.check-item:last-child {
  border-bottom: none;
}

.check-item :deep(.el-checkbox__label) {
  display: flex;
  align-items: center;
  gap: 12px;
}

.check-label {
  font-weight: 500;
}

.check-status {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
}

.check-status.pass {
  background: #f6ffed;
  color: #52c41a;
}

.check-status.warn {
  background: #fff7e6;
  color: #faad14;
}

.check-status.info {
  background: #e6f7ff;
  color: #1677ff;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-buttons .el-button {
  width: 100%;
}

.export-buttons {
  display: flex;
  gap: 12px;
}

.export-buttons .el-button {
  flex: 1;
}
</style>
