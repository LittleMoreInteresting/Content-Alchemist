<template>
  <div class="loading-view">
    <div class="loading-content">
      <div class="logo">✍️</div>
      <h1>Content-Alchemist</h1>
      
      <!-- 密钥初始化步骤 -->
      <div v-if="step === 'key-check'" class="step-section">
        <el-skeleton :rows="2" animated />
        <p class="loading-text">正在检查系统环境...</p>
      </div>
      
      <!-- 首次使用初始化 -->
      <div v-else-if="step === 'key-init'" class="step-section">
        <el-result
          icon="info"
          title="欢迎使用 Content-Alchemist"
          sub-title="首次使用需要初始化系统安全环境"
        >
          <template #extra>
            <el-button type="primary" @click="initKey" :loading="initLoading">
              初始化系统
            </el-button>
          </template>
        </el-result>
      </div>
      
      <!-- 配置检查 -->
      <div v-else-if="step === 'config-check'" class="step-section">
        <el-skeleton :rows="2" animated />
        <p class="loading-text">正在加载配置...</p>
      </div>
      
      <!-- 错误状态 -->
      <div v-else-if="step === 'error'" class="step-section">
        <el-result
          icon="error"
          title="初始化失败"
          :sub-title="errorMsg"
        >
          <template #extra>
            <el-button @click="retry">重试</el-button>
          </template>
        </el-result>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useConfigStore } from '@/stores'
import { checkHasEncryptionKey, initEncryptionKey, checkHasConfig, fetchConfig } from '@/api'
import { ElMessage } from 'element-plus'

const router = useRouter()
const configStore = useConfigStore()

const step = ref<'key-check' | 'key-init' | 'config-check' | 'error'>('key-check')
const initLoading = ref(false)
const errorMsg = ref('')

onMounted(async () => {
  await checkSystem()
})

async function checkSystem() {
  try {
    step.value = 'key-check'
    
    // 1. 检查加密密钥是否存在
    const hasKey = await checkHasEncryptionKey()
    
    if (!hasKey) {
      // 需要初始化密钥
      step.value = 'key-init'
      return
    }
    
    // 2. 密钥存在，继续检查配置
    step.value = 'config-check'
    await checkConfig()
    
  } catch (error) {
    console.error('系统检查失败:', error)
    errorMsg.value = '系统检查失败: ' + error
    step.value = 'error'
  }
}

async function initKey() {
  initLoading.value = true
  try {
    await initEncryptionKey()
    ElMessage.success('系统初始化成功')
    // 初始化后继续检查配置
    step.value = 'config-check'
    await checkConfig()
  } catch (error) {
    ElMessage.error('初始化失败: ' + error)
    errorMsg.value = '初始化失败: ' + error
    step.value = 'error'
  } finally {
    initLoading.value = false
  }
}

async function checkConfig() {
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
    console.error('配置检查失败:', error)
    errorMsg.value = '配置检查失败: ' + error
    step.value = 'error'
  }
}

function retry() {
  checkSystem()
}
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
  width: 400px;
  padding: 40px;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
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

.step-section {
  min-height: 200px;
}

.loading-text {
  margin-top: 24px;
  color: #8c8c8c;
  font-size: 14px;
}
</style>
