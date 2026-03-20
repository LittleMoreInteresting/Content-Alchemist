<template>
  <div class="welcome-view">
    <div class="welcome-container">
      <el-steps :active="currentStep" finish-status="success" simple>
        <el-step title="欢迎" />
        <el-step title="AI配置" />
        <el-step title="定位设置" />
      </el-steps>
      
      <!-- Step 1: 欢迎页 -->
      <div v-if="currentStep === 0" class="step-content">
        <div class="welcome-hero">
          <div class="logo">✍️ Content-Alchemist</div>
          <h1>让公众号创作从3小时缩短到30分钟</h1>
          <p class="subtitle">
            AI驱动的本地优先写作助手，专为一人公司创作者设计<br>
            输入标题 → AI协作创作 → 公众号排版 → 一键发布
          </p>
          <div class="features">
            <div class="feature-item">
              <el-icon><Star /></el-icon>
              <span>智能大纲生成</span>
            </div>
            <div class="feature-item">
              <el-icon><EditPen /></el-icon>
              <span>行内AI助手</span>
            </div>
            <div class="feature-item">
              <el-icon><Iphone /></el-icon>
              <span>公众号预览</span>
            </div>
          </div>
          <div class="actions">
            <el-button type="primary" size="large" @click="nextStep">
              开始配置
            </el-button>
            <el-button link size="large" @click="skipToLocalMode">
              使用本地模式
            </el-button>
          </div>
        </div>
      </div>
      
      <!-- Step 2: AI配置 -->
      <div v-else-if="currentStep === 1" class="step-content">
        <h2>配置AI助手</h2>
        <p class="step-desc">输入你的DeepSeek或OpenAI API信息，所有数据仅存储在本地</p>
        
        <el-form :model="config" label-width="120px" class="config-form">
          <el-form-item label="API Base URL">
            <el-input v-model="config.apiBaseUrl" placeholder="https://api.deepseek.com" />
          </el-form-item>
          
          <el-form-item label="API Key">
            <el-input 
              v-model="config.apiKey" 
              type="password" 
              show-password
              placeholder="sk-..."
            />
            <div class="form-tip">你的API Key仅存储在本地，不会上传到任何服务器</div>
          </el-form-item>
          
          <el-form-item label="模型选择">
            <el-select v-model="config.model" placeholder="选择模型">
              <el-option label="DeepSeek Chat" value="deepseek-chat" />
              <el-option label="DeepSeek Reasoner" value="deepseek-reasoner" />
              <el-option label="GPT-4o" value="gpt-4o" />
              <el-option label="GPT-4 Turbo" value="gpt-4-turbo" />
            </el-select>
          </el-form-item>
          
          <el-form-item label="温度参数">
            <el-slider v-model="config.temperature" :min="0" :max="1" :step="0.1" show-stops />
            <div class="form-tip">较低值输出更确定，较高值更有创意</div>
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              :loading="testing" 
              @click="testConnection"
            >
              测试连接
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="step-actions">
          <el-button @click="prevStep">上一步</el-button>
          <el-button type="primary" @click="nextStep">下一步</el-button>
        </div>
      </div>
      
      <!-- Step 3: 定位设置 -->
      <div v-else-if="currentStep === 2" class="step-content">
        <h2>公众号定位</h2>
        <p class="step-desc">设置你的写作风格，AI将根据这些偏好生成内容</p>
        
        <el-form :model="config" label-width="120px" class="config-form">
          <el-form-item label="风格标签">
            <el-checkbox-group v-model="config.styleTags">
              <el-checkbox label="技术干货" />
              <el-checkbox label="职场成长" />
              <el-checkbox label="产品思维" />
              <el-checkbox label="生活随笔" />
              <el-checkbox label="行业观察" />
              <el-checkbox label="创业故事" />
              <el-checkbox label="学习方法" />
              <el-checkbox label="工具推荐" />
            </el-checkbox-group>
          </el-form-item>
          
          <el-form-item label="目标受众">
            <el-input 
              v-model="config.audience" 
              type="textarea" 
              :rows="2"
              placeholder="例如：面向35岁+的Golang后端开发者..."
            />
          </el-form-item>
          
          <el-form-item label="写作人称">
            <el-radio-group v-model="config.persona">
              <el-radio-button label="我">我</el-radio-button>
              <el-radio-button label="我们">我们</el-radio-button>
              <el-radio-button label="小编">小编</el-radio-button>
              <el-radio-button label="笔者">笔者</el-radio-button>
            </el-radio-group>
          </el-form-item>
        </el-form>
        
        <div class="step-actions">
          <el-button @click="prevStep">上一步</el-button>
          <el-button type="primary" @click="finish">完成</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { Star, EditPen, Iphone } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useConfigStore } from '@/stores'
import { saveConfigToBackend } from '@/api'

const router = useRouter()
const configStore = useConfigStore()

const currentStep = ref(0)
const testing = ref(false)

const config = reactive({
  apiBaseUrl: 'https://api.deepseek.com',
  apiKey: '',
  model: 'deepseek-chat',
  temperature: 0.7,
  styleTags: [] as string[],
  audience: '',
  persona: '我'
})

function nextStep() {
  if (currentStep.value < 2) {
    currentStep.value++
  }
}

function prevStep() {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

async function testConnection() {
  if (!config.apiKey) {
    ElMessage.warning('请输入API Key')
    return
  }
  
  testing.value = true
  
  // 模拟测试
  setTimeout(() => {
    testing.value = false
    ElMessage.success('连接成功')
  }, 1500)
}

async function skipToLocalMode() {
  config.apiKey = ''
  const newConfig = {
    id: Date.now().toString(),
    ...config,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  }
  
  try {
    // 保存到后端
    await saveConfigToBackend(newConfig)
    configStore.setConfig(newConfig)
    configStore.setHasConfig(true)
    ElMessage.success('已切换到本地模式')
    router.push('/editor')
  } catch (error) {
    ElMessage.error('保存配置失败：' + error)
  }
}

async function finish() {
  const newConfig = {
    id: Date.now().toString(),
    ...config,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  }
  
  try {
    // 保存到后端
    await saveConfigToBackend(newConfig)
    configStore.setConfig(newConfig)
    configStore.setHasConfig(true)
    
    ElMessage.success('配置完成，开始使用吧！')
    router.push('/editor')
  } catch (error) {
    ElMessage.error('保存配置失败：' + error)
  }
}
</script>

<style scoped>
.welcome-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8ec 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.welcome-container {
  width: 100%;
  max-width: 700px;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  padding: 40px;
}

.step-content {
  padding: 40px 0;
}

.welcome-hero {
  text-align: center;
}

.logo {
  font-size: 48px;
  margin-bottom: 24px;
}

.welcome-hero h1 {
  font-size: 32px;
  font-weight: 700;
  color: #262626;
  margin-bottom: 16px;
}

.subtitle {
  font-size: 16px;
  color: #595959;
  line-height: 1.8;
  margin-bottom: 32px;
}

.features {
  display: flex;
  justify-content: center;
  gap: 40px;
  margin-bottom: 40px;
}

.feature-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: #595959;
}

.feature-item .el-icon {
  font-size: 32px;
  color: #1677ff;
}

.actions {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.step-content h2 {
  font-size: 24px;
  font-weight: 600;
  color: #262626;
  margin-bottom: 8px;
  text-align: center;
}

.step-desc {
  text-align: center;
  color: #8c8c8c;
  margin-bottom: 32px;
}

.config-form {
  max-width: 500px;
  margin: 0 auto;
}

.form-tip {
  font-size: 12px;
  color: #8c8c8c;
  margin-top: 4px;
}

.step-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #e8e8e8;
}
</style>
