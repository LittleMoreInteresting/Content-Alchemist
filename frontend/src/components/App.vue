<template>
  <div class="app">
    <!-- 顶部工具栏 -->
    <FileToolbar
      :title="articleTitle"
      :is-dirty="isDirty"
      :is-saving="isSaving"
      :can-save="canSave"
      @save="handleSave"
      @settings="showSettings = true"
      @error="handleError"
    />

    <!-- 主体区域 -->
    <div class="main-container">
      <!-- 左侧边栏：写作助手 -->
      <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
        <WritingSidebar
          ref="sidebarRef"
          v-model:title="articleTitle"
          v-model:outline="articleOutline"
          v-model:requirements="writingRequirements"
          :positioning="savedPositioning"
          :is-generating-outline="isGeneratingOutline"
          :is-generating-article="isGeneratingArticle"
          @generate-outline="handleGenerateOutline"
          @generate-article="handleGenerateArticle"
          @open-recent="handleOpenRecent"
          @error="handleError"
        />
      </aside>

      <!-- 编辑器和预览区域 -->
      <div class="editor-container">
        <!-- 编辑区域 -->
        <div class="editor-pane">
          <div class="pane-header">
            <span class="pane-title">✏️ 编辑</span>
            <span class="pane-subtitle" v-if="!articleOutline">输入大纲或文章内容</span>
            <span class="pane-subtitle" v-else-if="!articleContent">基于大纲生成文章</span>
            <span class="pane-subtitle" v-else>编辑文章内容</span>
          </div>
          <div class="editor-content">
            <textarea
              ref="editorRef"
              v-model="editorContent"
              placeholder="开始写作..."
              class="markdown-editor"
              spellcheck="false"
            ></textarea>
          </div>
        </div>

        <!-- 预览区域 -->
        <div class="preview-pane">
          <div class="pane-header">
            <span class="pane-title">👁️ 预览</span>
            <div class="pane-actions">
              <button
                class="pane-btn"
                :class="{ active: showOutline }"
                @click="showOutline = !showOutline"
                title="显示/隐藏目录"
              >
                目录
              </button>
            </div>
          </div>
          <div class="preview-content" ref="previewRef">
            <div class="markdown-body" v-html="renderedContent"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 全局 Loading 遮罩 -->
    <div v-if="isLoading" class="global-loading">
      <div class="loading-spinner-large"></div>
      <p class="loading-text">{{ loadingText }}</p>
    </div>

    <!-- 标题选择弹窗 -->
    <TitleSelectorModal
      v-model:visible="showTitleSelector"
      :titles="candidateTitles"
      :original-title="articleTitle"
      @select="handleTitleSelected"
    />

    <!-- 设置弹窗 -->
    <SettingsModal
      v-model:visible="showSettings"
      :initial-config="aiConfig"
      :initial-positioning="savedPositioning"
      @save="handleSettingsSave"
      @save-positioning="handlePositioningSave"
    />

    <!-- 错误提示 -->
    <Transition name="toast">
      <div v-if="errorMessage" class="error-toast" @click="errorMessage = ''">
        {{ errorMessage }}
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import MarkdownIt from 'markdown-it';
import FileToolbar from './FileToolbar.vue';
import WritingSidebar from './WritingSidebar.vue';
import TitleSelectorModal from './TitleSelectorModal.vue';
import SettingsModal, { type AIConfig } from './SettingsModal.vue';
import { useWails, FileError } from '../composables/useWails';

// ============================================
// Markdown 渲染
// ============================================
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  breaks: true,
});

// ============================================
// 状态
// ============================================
const wails = useWails();
const sidebarRef = ref<InstanceType<typeof WritingSidebar> | null>(null);

// 文章状态
const articleTitle = ref('');
const articleOutline = ref('');
const articleContent = ref('');
const editorContent = ref('');
const writingRequirements = ref('');
const savedPositioning = ref('');
const currentUUID = ref('');
const filePath = ref('');
const candidateTitles = ref<string[]>([]);

// UI 状态
const sidebarCollapsed = ref(false);
const showOutline = ref(false);
const showSettings = ref(false);
const showTitleSelector = ref(false);
const errorMessage = ref('');
const isSaving = ref(false);
const isDirty = ref(false);
const isGeneratingOutline = ref(false);
const isGeneratingArticle = ref(false);
const isLoading = ref(false);
const loadingText = ref('');

// AI 配置
const aiConfig = ref<AIConfig>({
  baseUrl: 'https://api.deepseek.com/v1',
  token: '',
  temperature: 0.7,
  model: 'deepseek-chat'
});

// ============================================
// 计算属性
// ============================================
const canSave = computed(() => {
  return isDirty.value && !isSaving.value && articleTitle.value !== '';
});

const renderedContent = computed(() => {
  return md.render(editorContent.value);
});

// 监听编辑器内容变化
watch(editorContent, (newValue, oldValue) => {
  if (oldValue !== undefined && newValue !== oldValue) {
    isDirty.value = true;
    
    // 更新大纲或文章内容
    if (!articleContent.value && articleOutline.value) {
      // 如果还没有正式生成文章，编辑的是大纲
      articleOutline.value = newValue;
    } else {
      // 编辑的是文章内容
      articleContent.value = newValue;
    }
  }
});

// ============================================
// Loading 控制
// ============================================
const showLoading = (text: string) => {
  loadingText.value = text;
  isLoading.value = true;
};

const hideLoading = () => {
  isLoading.value = false;
  loadingText.value = '';
};

// ============================================
// 保存功能
// ============================================
const handleSave = async (overwrite = false): Promise<void> => {
  if (!canSave.value && !overwrite) return;

  isSaving.value = true;

  try {
    const result = await wails.saveArticleWithSmartNaming(
      currentUUID.value,
      articleTitle.value,
      editorContent.value,
      overwrite
    );

    if (result) {
      currentUUID.value = result.uuid;
      filePath.value = result.filePath;
      isDirty.value = false;
      // 刷新最近文章列表
      sidebarRef.value?.loadRecentArticles();
    }
  } catch (err) {
    // 检查是否是文件存在需要确认的错误
    if (err instanceof FileError && err.code === 'FILE_EXISTS_CONFIRM') {
      // 提取文件路径
      const filePathMatch = err.message.match(/文件已存在: (.+)$/);
      const existingFilePath = filePathMatch ? filePathMatch[1] : '';
      
      // 显示确认对话框
      const confirmed = window.confirm(
        `文件 "${existingFilePath}" 已存在。\n\n是否覆盖保存？\n（旧版本将被替换）`
      );
      
      if (confirmed) {
        // 用户确认覆盖，重新调用保存并设置 overwrite 为 true
        await handleSave(true);
        return;
      }
    } else {
      console.error('保存失败:', err);
      handleError(err instanceof Error ? err.message : '保存失败');
    }
  } finally {
    isSaving.value = false;
  }
};

// ============================================
// 生成大纲
// ============================================
const handleGenerateOutline = async (): Promise<void> => {
  if (!articleTitle.value) {
    handleError('请先输入文章标题');
    return;
  }

  isGeneratingOutline.value = true;
  showLoading('正在生成大纲和候选标题...');

  try {
    // 调用生成大纲接口，传入标题、写作要求、公众号定位
    const result = await wails.generateOutline(
      articleTitle.value,
      writingRequirements.value,
      savedPositioning.value
    );
    
    if (result) {
      // 保存候选标题
      candidateTitles.value = result.titles;
      
      // 保存大纲
      articleOutline.value = result.outline;
      
      // 显示标题选择弹窗
      if (candidateTitles.value.length > 0) {
        showTitleSelector.value = true;
      } else {
        // 如果没有候选标题，直接显示大纲
        editorContent.value = articleOutline.value;
        isDirty.value = true;
        // 自动保存
        await autoSave();
      }
    }
  } catch (err) {
    console.error('生成大纲失败:', err);
    handleError(err instanceof Error ? err.message : '生成大纲失败');
  } finally {
    isGeneratingOutline.value = false;
    hideLoading();
  }
};

// ============================================
// 标题选择处理
// ============================================
const handleTitleSelected = async (selectedTitle: string): Promise<void> => {
  // 更新文章标题
  articleTitle.value = selectedTitle;
  
  // 显示大纲到编辑器
  editorContent.value = articleOutline.value;
  isDirty.value = true;
  
  // 自动保存
  await autoSave();
};

// ============================================
// 自动保存
// ============================================
const autoSave = async (overwrite = false): Promise<void> => {
  if (!articleTitle.value || !editorContent.value) return;
  
  isSaving.value = true;
  
  try {
    const result = await wails.saveArticleWithSmartNaming(
      currentUUID.value,
      articleTitle.value,
      editorContent.value,
      overwrite
    );

    if (result) {
      currentUUID.value = result.uuid;
      filePath.value = result.filePath;
      isDirty.value = false;
      // 刷新最近文章列表
      sidebarRef.value?.loadRecentArticles();
    }
  } catch (err) {
    // 自动保存时如果文件存在，自动覆盖（不打扰用户）
    if (err instanceof FileError && err.code === 'FILE_EXISTS_CONFIRM') {
      await autoSave(true);
      return;
    }
    
    console.error('自动保存失败:', err);
    // 自动保存失败不提示错误，避免打扰用户
  } finally {
    isSaving.value = false;
  }
};

// ============================================
// 文章生成处理
// ============================================
const handleGenerateArticle = async (): Promise<void> => {
  if (!articleTitle.value) {
    handleError('请先输入文章标题');
    return;
  }

  if (!articleOutline.value) {
    handleError('请先生成大纲');
    return;
  }

  isGeneratingArticle.value = true;
  showLoading('正在生成文章内容，请稍候...');

  try {
    // 合并写作要求和公众号定位
    const combinedRequirements = [
      writingRequirements.value,
      savedPositioning.value ? `公众号定位：${savedPositioning.value}` : '',
    ].filter(Boolean).join('\n\n');

    const content = await wails.generateArticle(
      articleTitle.value,
      articleOutline.value,
      combinedRequirements
    );

    if (content) {
      // 添加标题到文章内容
      const fullContent = `# ${articleTitle.value}\n\n${content}`;
      articleContent.value = fullContent;
      editorContent.value = fullContent;
      isDirty.value = true;
      
      // 生成文章后自动保存
      await autoSave();
    }
  } catch (err) {
    console.error('生成文章失败:', err);
    handleError(err instanceof Error ? err.message : '生成文章失败');
  } finally {
    isGeneratingArticle.value = false;
    hideLoading();
  }
};

// ============================================
// 打开最近文章
// ============================================
const handleOpenRecent = async (path: string): Promise<void> => {
  if (isDirty.value) {
    const confirmed = window.confirm('当前文章有未保存的更改，是否继续打开其他文章？');
    if (!confirmed) return;
  }

  try {
    const { article, content } = await wails.readArticle(path);
    
    if (article) {
      currentUUID.value = article.uuid;
      filePath.value = article.filePath;
      articleTitle.value = article.title || '';
      editorContent.value = content;
      articleContent.value = content;
      
      // 尝试从内容中提取大纲
      const outlineMatch = content.match(/##?\s+.+$/m);
      if (outlineMatch) {
        articleOutline.value = outlineMatch[0];
      }
      
      isDirty.value = false;
    }
  } catch (err) {
    console.error('打开文章失败:', err);
    handleError(err instanceof Error ? err.message : '打开文章失败');
  }
};

// ============================================
// AI配置和公众号定位
// ============================================
const loadAIConfig = async (): Promise<void> => {
  try {
    const config = await wails.getAIConfig();
    aiConfig.value = config;
  } catch (err) {
    console.error('加载AI配置失败:', err);
  }
};

const loadPositioning = async (): Promise<void> => {
  try {
    const positioning = await wails.getPositioning();
    savedPositioning.value = positioning;
  } catch (err) {
    console.error('加载公众号定位失败:', err);
  }
};

const handleSettingsSave = async (config: AIConfig): Promise<void> => {
  try {
    await wails.saveAIConfig(config);
    aiConfig.value = config;
    console.log('AI配置已保存');
  } catch (err) {
    console.error('保存AI配置失败:', err);
    handleError(err instanceof Error ? err.message : '保存配置失败');
  }
};

const handlePositioningSave = async (positioning: string): Promise<void> => {
  try {
    await wails.savePositioning(positioning);
    savedPositioning.value = positioning;
    console.log('公众号定位已保存');
  } catch (err) {
    console.error('保存公众号定位失败:', err);
    handleError(err instanceof Error ? err.message : '保存公众号定位失败');
  }
};

// ============================================
// 工具函数
// ============================================
const handleError = (message: string): void => {
  errorMessage.value = message;
  setTimeout(() => {
    errorMessage.value = '';
  }, 5000);
};

// ============================================
// 键盘快捷键
// ============================================
const handleKeyDown = (e: KeyboardEvent): void => {
  if ((e.ctrlKey || e.metaKey) && e.key === ',') {
    e.preventDefault();
    showSettings.value = true;
  }
};

// ============================================
// 生命周期
// ============================================
onMounted(() => {
  document.addEventListener('keydown', handleKeyDown);
  loadAIConfig();
  loadPositioning();
});

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown);
});
</script>

<style scoped>
.app {
  display: flex;
  flex-direction: column;
  height: 100vh;
  font-family: var(--font-family-base, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif);
  background: var(--bg-component, #fff);
  color: var(--text-primary, #262626);
}

/* 主体区域 */
.main-container {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* 侧边栏 */
.sidebar {
  width: 300px;
  flex-shrink: 0;
  border-right: 1px solid var(--border-color, #e8e8e8);
  transition: width 0.3s;
  overflow: hidden;
  background: var(--panel-bg, #fafafa);
}

.sidebar.collapsed {
  width: 0;
  border-right: none;
}

/* 编辑器容器 */
.editor-container {
  flex: 1;
  display: flex;
  overflow: hidden;
}

/* 编辑和预览面板 */
.editor-pane,
.preview-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
}

.editor-pane {
  border-right: 1px solid var(--border-color, #e8e8e8);
}

/* 面板头部 */
.pane-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  background: var(--bg-hover, #f5f5f5);
  border-bottom: 1px solid var(--border-color, #e8e8e8);
  font-size: 13px;
}

.pane-title {
  font-weight: 600;
  color: var(--text-primary, #262626);
}

.pane-subtitle {
  font-size: 12px;
  color: var(--text-secondary, #8c8c8c);
  margin-left: 8px;
}

.pane-actions {
  display: flex;
  gap: 6px;
}

.pane-btn {
  padding: 4px 10px;
  border: 1px solid var(--border-color, #d9d9d9);
  border-radius: 4px;
  background: var(--bg-component, #fff);
  cursor: pointer;
  font-size: 12px;
  color: var(--text-secondary, #595959);
  transition: all 0.2s;
}

.pane-btn:hover {
  border-color: var(--color-primary, #1890ff);
  color: var(--color-primary, #1890ff);
}

.pane-btn.active {
  background: var(--color-primary, #1890ff);
  border-color: var(--color-primary, #1890ff);
  color: white;
}

/* 编辑区域 */
.editor-content {
  flex: 1;
  overflow: auto;
}

.markdown-editor {
  width: 100%;
  height: 100%;
  padding: 20px;
  border: none;
  outline: none;
  resize: none;
  font-family: var(--font-family-mono, 'SF Mono', Monaco, 'Cascadia Code', monospace);
  font-size: 14px;
  line-height: 1.8;
  background: var(--bg-component, #fff);
  color: var(--text-primary, #262626);
  tab-size: 2;
  box-sizing: border-box;
}

.markdown-editor::placeholder {
  color: var(--text-muted, #bfbfbf);
}

/* 预览区域 */
.preview-content {
  flex: 1;
  overflow: auto;
  padding: 20px;
  background: var(--bg-component, #fff);
}

.markdown-body {
  max-width: 800px;
  margin: 0 auto;
  font-size: 15px;
  line-height: 1.8;
}

/* 全局 Loading 遮罩 */
.global-loading {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.loading-spinner-large {
  width: 60px;
  height: 60px;
  border: 4px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.loading-text {
  margin-top: 20px;
  font-size: 16px;
  color: white;
  font-weight: 500;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 错误提示 */
.error-toast {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  padding: 12px 24px;
  background: var(--error-bg, #fff2f0);
  border: 1px solid var(--error-border, #ffccc7);
  border-radius: 6px;
  color: var(--error-color, #ff4d4f);
  font-size: 13px;
  z-index: 1000;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 过渡动画 */
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(20px);
}

/* 暗色主题适配 */
@media (prefers-color-scheme: dark) {
  .app {
    --bg-component: #1f1f1f;
    --bg-hover: #2c2c2c;
    --panel-bg: #141414;
    --text-primary: #d9d9d9;
    --text-secondary: #8c8c8c;
    --text-muted: #595959;
    --border-color: #434343;
    --color-primary: #1890ff;
    --error-bg: #2a1215;
    --error-border: #58181c;
    --error-color: #ff7875;
  }
}
</style>
