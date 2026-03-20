<template>
  <div class="settings-view">
    <div class="settings-header">
      <el-button :icon="ArrowLeft" @click="goBack">返回</el-button>
      <h1>设置</h1>
      <span></span>
    </div>
    
    <div class="settings-content">
      <el-tabs tab-position="left" class="settings-tabs">
        <el-tab-pane label="AI配置">
          <div class="tab-content">
            <h2>AI配置</h2>
            <el-form :model="config" label-width="120px">
              <el-form-item label="API Base URL">
                <el-input v-model="config.apiBaseUrl" />
              </el-form-item>
              <el-form-item label="API Key">
                <el-input v-model="config.apiKey" type="password" show-password />
              </el-form-item>
              <el-form-item label="模型">
                <el-select v-model="config.model">
                  <el-option label="DeepSeek Chat" value="deepseek-chat" />
                  <el-option label="DeepSeek Reasoner" value="deepseek-reasoner" />
                  <el-option label="GPT-4o" value="gpt-4o" />
                </el-select>
              </el-form-item>
              <el-form-item label="温度">
                <el-slider v-model="config.temperature" :min="0" :max="1" :step="0.1" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="saveConfig">保存</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="写作风格">
          <div class="tab-content">
            <h2>写作风格</h2>
            <el-form :model="config" label-width="120px">
              <el-form-item label="风格标签">
                <el-checkbox-group v-model="config.styleTags">
                  <el-checkbox label="技术干货" />
                  <el-checkbox label="职场成长" />
                  <el-checkbox label="产品思维" />
                  <el-checkbox label="生活随笔" />
                  <el-checkbox label="行业观察" />
                </el-checkbox-group>
              </el-form-item>
              <el-form-item label="目标受众">
                <el-input v-model="config.audience" type="textarea" :rows="3" />
              </el-form-item>
              <el-form-item label="写作人称">
                <el-radio-group v-model="config.persona">
                  <el-radio-button label="我">我</el-radio-button>
                  <el-radio-button label="我们">我们</el-radio-button>
                  <el-radio-button label="小编">小编</el-radio-button>
                  <el-radio-button label="笔者">笔者</el-radio-button>
                </el-radio-group>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="saveConfig">保存</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="快捷键">
          <div class="tab-content">
            <h2>快捷键</h2>
            <div class="shortcuts-list">
              <div v-for="shortcut in shortcuts" :key="shortcut.key" class="shortcut-item">
                <span class="shortcut-desc">{{ shortcut.desc }}</span>
                <kbd class="shortcut-key">{{ shortcut.key }}</kbd>
              </div>
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="关于">
          <div class="tab-content">
            <h2>关于 Content-Alchemist</h2>
            <p>版本: v1.0.0</p>
            <p>基于 Eino ADK + Wails + Vue3 构建</p>
            <p>本地优先的AI写作助手</p>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useConfigStore } from '@/stores'
import { saveConfigToBackend } from '@/api'

const router = useRouter()
const configStore = useConfigStore()

const config = reactive({ ...configStore.config })

const shortcuts = [
  { key: 'Cmd/Ctrl + N', desc: '新建文章' },
  { key: 'Cmd/Ctrl + S', desc: '手动保存' },
  { key: 'Cmd/Ctrl + Shift + P', desc: '发布文章' },
  { key: 'Cmd/Ctrl + K', desc: '命令面板' },
  { key: 'Cmd/Ctrl + /', desc: 'AI斜杠命令' },
  { key: 'Esc', desc: '关闭弹窗' }
]

function goBack() {
  router.push('/editor')
}

async function saveConfig() {
  const updatedConfig = {
    ...config,
    updatedAt: new Date().toISOString()
  }
  
  try {
    await saveConfigToBackend(updatedConfig)
    configStore.setConfig(updatedConfig)
    ElMessage.success('配置已保存')
  } catch (error) {
    ElMessage.error('保存配置失败：' + error)
  }
}
</script>

<style scoped>
.settings-view {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #fff;
}

.settings-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid #e8e8e8;
}

.settings-header h1 {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}

.settings-content {
  flex: 1;
  overflow: hidden;
}

.settings-tabs {
  height: 100%;
}

.settings-tabs :deep(.el-tabs__content) {
  height: 100%;
  overflow-y: auto;
}

.tab-content {
  padding: 24px;
  max-width: 600px;
}

.tab-content h2 {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 24px;
}

.shortcuts-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.shortcut-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #f5f5f5;
  border-radius: 8px;
}

.shortcut-desc {
  color: #595959;
}

.shortcut-key {
  padding: 4px 12px;
  background: #fff;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
}
</style>
