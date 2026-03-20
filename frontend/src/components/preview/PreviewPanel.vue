<template>
  <div class="preview-panel">
    <div class="preview-header">
      <el-radio-group v-model="previewMode" size="small">
        <el-radio-button label="mobile">
          <el-icon><Iphone /></el-icon> 手机预览
        </el-radio-button>
        <el-radio-button label="desktop">
          <el-icon><Monitor /></el-icon> 桌面预览
        </el-radio-button>
      </el-radio-group>
      
      <el-dropdown @command="onThemeChange">
        <el-button size="small">
          主题 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="default">科技蓝</el-dropdown-item>
            <el-dropdown-item command="minimal">简约白</el-dropdown-item>
            <el-dropdown-item command="vibrant">活力橙</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
    
    <div class="preview-content" :class="previewMode">
      <!-- 手机模拟器 -->
      <div v-if="previewMode === 'mobile'" class="phone-mockup" :data-theme="currentTheme">
        <div class="phone-frame">
          <div class="phone-notch"></div>
          <div class="phone-screen">
            <div class="wechat-header">
              <span class="time">{{ currentTime }}</span>
              <div class="wechat-title">公众号文章</div>
            </div>
            <div class="article-content" v-html="renderedContent"></div>
            <div class="wechat-footer">
              <div class="action-buttons">
                <span>👍 赞</span>
                <span>💬 在看</span>
                <span>➡️ 分享</span>
              </div>
            </div>
          </div>
          <div class="phone-home-indicator"></div>
        </div>
      </div>
      
      <!-- 桌面预览 -->
      <div v-else class="desktop-preview" :data-theme="currentTheme">
        <article class="markdown-body" v-html="renderedContent"></article>
      </div>
    </div>
    
    <div class="preview-actions">
      <el-button type="primary" :icon="CopyDocument" @click="copyHTML">
        复制HTML
      </el-button>
      <el-button :icon="Download" @click="downloadMarkdown">
        下载Markdown
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Iphone, Monitor, ArrowDown, CopyDocument, Download } from '@element-plus/icons-vue'
import { useEditorStore } from '@/stores'
import { renderWechatHTML } from '@/utils/wechatRenderer'
import { ElMessage } from 'element-plus'

const editorStore = useEditorStore()
const previewMode = ref<'mobile' | 'desktop'>('mobile')
const currentTheme = computed(() => editorStore.currentTheme)

const currentTime = computed(() => {
  const now = new Date()
  return `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`
})

const renderedContent = computed(() => {
  return renderWechatHTML(editorStore.article.content)
})

function onThemeChange(theme: string) {
  editorStore.setTheme(theme)
}

function copyHTML() {
  const html = renderedContent.value
  navigator.clipboard.writeText(html).then(() => {
    ElMessage.success('HTML已复制到剪贴板')
  })
}

function downloadMarkdown() {
  const blob = new Blob([editorStore.article.content], { type: 'text/markdown' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${editorStore.article.title || 'untitled'}.md`
  a.click()
  URL.revokeObjectURL(url)
  ElMessage.success('Markdown已下载')
}
</script>

<style scoped>
.preview-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.preview-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #e8e8e8;
  background: #fafafa;
}

.preview-content {
  flex: 1;
  overflow: auto;
  padding: 20px;
  background: #f0f2f5;
}

/* 手机模拟器 */
.phone-mockup {
  display: flex;
  justify-content: center;
}

.phone-frame {
  width: 375px;
  height: 812px;
  background: #000;
  border-radius: 40px;
  padding: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  position: relative;
}

.phone-notch {
  position: absolute;
  top: 12px;
  left: 50%;
  transform: translateX(-50%);
  width: 150px;
  height: 30px;
  background: #000;
  border-radius: 0 0 20px 20px;
  z-index: 10;
}

.phone-screen {
  width: 100%;
  height: 100%;
  background: #fff;
  border-radius: 32px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.wechat-header {
  background: #ededed;
  padding: 44px 16px 12px;
  text-align: center;
  position: relative;
}

.time {
  position: absolute;
  left: 16px;
  top: 12px;
  font-size: 14px;
  font-weight: 600;
}

.wechat-title {
  font-size: 16px;
  font-weight: 600;
  color: #000;
}

.article-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px 16px;
}

.wechat-footer {
  padding: 12px 16px;
  border-top: 1px solid #e5e5e5;
  background: #f9f9f9;
}

.action-buttons {
  display: flex;
  justify-content: space-around;
  font-size: 14px;
  color: #576b95;
}

.phone-home-indicator {
  position: absolute;
  bottom: 8px;
  left: 50%;
  transform: translateX(-50%);
  width: 134px;
  height: 5px;
  background: #fff;
  border-radius: 3px;
}

/* 桌面预览 */
.desktop-preview {
  max-width: 800px;
  margin: 0 auto;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 40px;
}

.preview-actions {
  display: flex;
  gap: 12px;
  padding: 16px;
  border-top: 1px solid #e8e8e8;
  background: #fafafa;
}

.preview-actions .el-button {
  flex: 1;
}
</style>
