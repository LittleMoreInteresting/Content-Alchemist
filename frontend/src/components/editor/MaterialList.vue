<template>
  <div class="material-list">
    <div 
      v-for="material in materials" 
      :key="material.id"
      class="material-item"
      draggable="true"
      @dragstart="onDragStart($event, material)"
      @click="$emit('insert', material)"
    >
      <div class="material-header">
        <span class="material-title">{{ material.title }}</span>
        <span class="material-count">使用 {{ material.usageCount }} 次</span>
      </div>
      <p class="material-content">{{ material.content }}</p>
      <div class="material-footer">
        <el-tag 
          v-for="tag in material.tags.slice(0, 2)" 
          :key="tag"
          size="small"
          effect="plain"
        >
          {{ tag }}
        </el-tag>
        <span v-if="material.source" class="material-source">
          来源: {{ material.source }}
        </span>
      </div>
    </div>
    
    <el-empty v-if="materials.length === 0" description="暂无素材" />
  </div>
</template>

<script setup lang="ts">
import type { Material } from '@/types'

defineProps<{
  type: string
  materials: Material[]
}>()

defineEmits<{
  insert: [material: Material]
}>()

function onDragStart(e: DragEvent, material: Material) {
  if (e.dataTransfer) {
    e.dataTransfer.setData('application/json', JSON.stringify(material))
    e.dataTransfer.effectAllowed = 'copy'
  }
}
</script>

<style scoped>
.material-list {
  padding: 8px;
}

.material-item {
  padding: 12px;
  margin-bottom: 8px;
  background: #f6ffed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.material-item:hover {
  background: #d9f7be;
  transform: translateX(4px);
}

.material-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.material-title {
  font-weight: 500;
  color: #262626;
  font-size: 13px;
}

.material-count {
  font-size: 11px;
  color: #8c8c8c;
}

.material-content {
  font-size: 12px;
  color: #595959;
  line-height: 1.5;
  margin: 0 0 8px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.material-footer {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.material-footer .el-tag {
  font-size: 10px;
}

.material-source {
  font-size: 10px;
  color: #8c8c8c;
  margin-left: auto;
}
</style>
