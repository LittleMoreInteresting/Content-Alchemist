<template>
  <div class="loading-view">
    <div class="loading-content">
      <div class="logo">✍️</div>
      <h1>Content-Alchemist</h1>
      <el-skeleton :rows="3" animated />
      <p class="loading-text">正在初始化...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useConfigStore } from '@/stores'
import { checkHasConfig, fetchConfig } from '@/api'

const router = useRouter()
const configStore = useConfigStore()

onMounted(async () => {
  try {
    // 检查是否有配置
    const hasConfig = await checkHasConfig()
    
    if (hasConfig) {
      // 加载配置
      const config = await fetchConfig()
      configStore.setConfig(config)
      // 跳转到编辑器
      router.replace('/editor')
    } else {
      // 没有配置，跳转到欢迎页
      router.replace('/welcome')
    }
  } catch (error) {
    console.error('初始化失败:', error)
    // 出错时跳转到欢迎页
    router.replace('/welcome')
  }
})
</script>

<style scoped>
.loading-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8ec 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-content {
  text-align: center;
  width: 300px;
}

.logo {
  font-size: 64px;
  margin-bottom: 16px;
}

h1 {
  font-size: 24px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 32px;
}

.loading-text {
  margin-top: 24px;
  color: #8c8c8c;
  font-size: 14px;
}
</style>
