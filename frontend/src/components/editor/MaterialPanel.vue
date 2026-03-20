<template>
  <div class="material-panel">
    <div class="panel-header">
      <h3>
        <el-icon><Collection /></el-icon>
        素材库
      </h3>
      <el-button 
        type="primary" 
        size="small" 
        :icon="Plus"
        @click="showAddDialog = true"
      >
        添加
      </el-button>
    </div>
    
    <el-tabs v-model="activeTab" class="material-tabs">
      <el-tab-pane name="snippet">
        <template #label>
          <span><el-icon><Document /></el-icon> 片段</span>
        </template>
        <MaterialList type="snippet" :materials="filteredMaterials" @insert="insertMaterial" />
      </el-tab-pane>
      
      <el-tab-pane name="data">
        <template #label>
          <span><el-icon><DataLine /></el-icon> 数据</span>
        </template>
        <MaterialList type="data" :materials="filteredMaterials" @insert="insertMaterial" />
      </el-tab-pane>
      
      <el-tab-pane name="quote">
        <template #label>
          <span><el-icon><ChatDotRound /></el-icon> 金句</span>
        </template>
        <MaterialList type="quote" :materials="filteredMaterials" @insert="insertMaterial" />
      </el-tab-pane>
      
      <el-tab-pane name="history">
        <template #label>
          <span><el-icon><Timer /></el-icon> 历史</span>
        </template>
        <MaterialList type="history" :materials="filteredMaterials" @insert="insertMaterial" />
      </el-tab-pane>
    </el-tabs>
    
    <!-- 添加素材对话框 -->
    <el-dialog v-model="showAddDialog" title="添加素材" width="500px">
      <el-form :model="newMaterial" label-width="80px">
        <el-form-item label="类型">
          <el-radio-group v-model="newMaterial.type">
            <el-radio-button label="snippet">片段</el-radio-button>
            <el-radio-button label="data">数据</el-radio-button>
            <el-radio-button label="quote">金句</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="标题">
          <el-input v-model="newMaterial.title" placeholder="素材标题" />
        </el-form-item>
        <el-form-item label="内容">
          <el-input 
            v-model="newMaterial.content" 
            type="textarea" 
            :rows="4"
            placeholder="素材内容"
          />
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="newMaterial.tags" placeholder="用逗号分隔多个标签" />
        </el-form-item>
        <el-form-item label="来源">
          <el-input v-model="newMaterial.source" placeholder="可选：素材来源" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="addMaterial">添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Collection, Plus, Document, DataLine, ChatDotRound, Timer } from '@element-plus/icons-vue'
import MaterialList from './MaterialList.vue'
import type { Material } from '@/types'

const activeTab = ref('snippet')
const showAddDialog = ref(false)

const materials = ref<Material[]>([
  {
    id: '1',
    type: 'snippet',
    title: '文章开头模板',
    content: '在数字化转型的浪潮中，...',
    tags: ['模板', '开头'],
    source: '',
    createdAt: new Date().toISOString(),
    usageCount: 5
  },
  {
    id: '2',
    type: 'quote',
    title: '关于创新的名言',
    content: '创新是区分领导者和追随者的唯一标准。——史蒂夫·乔布斯',
    tags: ['创新', '名言'],
    source: '史蒂夫·乔布斯',
    createdAt: new Date().toISOString(),
    usageCount: 3
  }
])

const newMaterial = ref({
  type: 'snippet' as const,
  title: '',
  content: '',
  tags: '',
  source: ''
})

const filteredMaterials = computed(() => {
  return materials.value.filter(m => m.type === activeTab.value)
})

function insertMaterial(material: Material) {
  // 触发插入事件，由父组件处理
  console.log('Insert material:', material)
}

function addMaterial() {
  const material: Material = {
    id: Date.now().toString(),
    type: newMaterial.value.type,
    title: newMaterial.value.title,
    content: newMaterial.value.content,
    tags: newMaterial.value.tags.split(',').map(t => t.trim()).filter(Boolean),
    source: newMaterial.value.source,
    createdAt: new Date().toISOString(),
    usageCount: 0
  }
  materials.value.push(material)
  showAddDialog.value = false
  
  // 重置表单
  newMaterial.value = {
    type: 'snippet',
    title: '',
    content: '',
    tags: '',
    source: ''
  }
}
</script>

<style scoped>
.material-panel {
  display: flex;
  flex-direction: column;
  height: 50%;
  background: #fff;
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

.material-tabs {
  flex: 1;
  overflow: hidden;
}

.material-tabs :deep(.el-tabs__header) {
  margin-bottom: 0;
}

.material-tabs :deep(.el-tabs__content) {
  height: calc(100% - 40px);
  overflow-y: auto;
}

.material-tabs :deep(.el-tab-pane) {
  height: 100%;
}
</style>
