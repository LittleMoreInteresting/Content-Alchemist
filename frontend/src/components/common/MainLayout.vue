<template>
  <div class="main-layout">
    <header class="toolbar">
      <div class="toolbar-left">
        <input
          v-model="articleTitle"
          class="title-input"
          placeholder="请输入文章标题"
          @input="onTitleChange"
        />
        <span class="save-status" :class="saveStatus">
          {{ saveStatusText }}
        </span>
      </div>
      
      <div class="toolbar-center">
        <el-radio-group v-model="editorView" size="small">
          <el-radio-button label="edit">
            <el-icon><EditPen /></el-icon> 编辑
          </el-radio-button>
          <el-radio-button label="split">
            <el-icon><Monitor /></el-icon> 分屏
          </el-radio-button>
          <el-radio-button label="preview">
            <el-icon><View /></el-icon> 预览
          </el-radio-button>
        </el-radio-group>
      </div>
      
      <div class="toolbar-right">
        <el-button type="primary" :icon="Lightning" @click="showAIAssistant">
          AI助手
        </el-button>
        <el-button :icon="Setting" @click="goToSettings" />
        <el-button type="success" :icon="Position" @click="publish">
          发布
        </el-button>
      </div>
    </header>

    <main class="main-content" :class="[`view-${editorView}`]">
      <aside v-show="editorView !== 'preview'" class="sidebar-left">
        <OutlinePanel />
        <MaterialPanel />
      </aside>
      
      <section v-show="editorView !== 'preview'" class="editor-section">
        <MarkdownEditor />
      </section>
      
      <aside v-show="editorView !== 'edit'" class="sidebar-right">
        <PreviewPanel />
      </aside>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { EditPen, Monitor, View, Lightning, Setting, Position } from '@element-plus/icons-vue'
import { useEditorStore } from '@/stores'
import OutlinePanel from '../editor/OutlinePanel.vue'
import MaterialPanel from '../editor/MaterialPanel.vue'
import MarkdownEditor from '../editor/MarkdownEditor.vue'
import PreviewPanel from '../preview/PreviewPanel.vue'

const router = useRouter()
const editorStore = useEditorStore()

const articleTitle = computed({
  get: () => editorStore.article.title,
  set: (val) => editorStore.updateTitle(val)
})

const editorView = computed({
  get: () => editorStore.editorView,
  set: (val) => editorStore.setEditorView(val)
})

const saveStatus = computed(() => editorStore.saveStatus)

const saveStatusText = computed(() => {
  switch (saveStatus.value) {
    case 'saved': return '已保存'
    case 'saving': return '保存中...'
    case 'unsaved': return '未保存'
    default: return ''
  }
})

function onTitleChange() {
  editorStore.markUnsaved()
}

function showAIAssistant() {
  // 触发AI助手命令面板
  window.dispatchEvent(new KeyboardEvent('keydown', { key: '/', ctrlKey: true }))
}

function goToSettings() {
  router.push('/settings')
}

function publish() {
  router.push('/publish')
}
</script>

<style scoped>
.main-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: var(--bg-color, #f0f2f5);
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  background: #fff;
  border-bottom: 1px solid #e8e8e8;
  gap: 16px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.title-input {
  font-size: 18px;
  font-weight: 600;
  border: none;
  outline: none;
  background: transparent;
  padding: 4px 8px;
  width: 300px;
}

.title-input::placeholder {
  color: #bfbfbf;
}

.save-status {
  font-size: 12px;
  color: #8c8c8c;
}

.save-status.unsaved {
  color: #faad14;
}

.toolbar-center {
  display: flex;
  align-items: center;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar-left {
  width: 240px;
  min-width: 200px;
  max-width: 350px;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-right: 1px solid #e8e8e8;
  resize: horizontal;
  overflow: auto;
}

.editor-section {
  flex: 1;
  min-width: 400px;
  background: #fff;
  overflow: hidden;
}

.sidebar-right {
  width: 400px;
  min-width: 320px;
  max-width: 500px;
  background: #fff;
  border-left: 1px solid #e8e8e8;
  overflow: auto;
  resize: horizontal;
}

/* 视图模式调整 */
.view-edit .sidebar-left {
  width: 280px;
}

.view-edit .editor-section {
  flex: 1;
}

.view-preview .sidebar-right {
  flex: 1;
  width: 100%;
  max-width: none;
}
</style>
