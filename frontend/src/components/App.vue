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
      <!-- 左侧边栏：对话式写作助手 -->
      <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
        <WritingChatSidebar
          ref="sidebarRef"
          :article-title="articleTitle"
          :is-processing="isProcessing"
          :streaming-content="streamingContent"
          @send-message="handleSendMessage"
          @quick-action="handleQuickAction"
          @open-recent="handleOpenRecent"
          @update:step="handleStepChange"
        />
      </aside>

      <!-- 编辑器和预览区域 -->
      <div class="editor-container">
        <!-- 编辑区域 -->
        <div class="editor-pane">
          <div class="pane-header">
            <span class="pane-title">✏️ 编辑</span>
            <div class="pane-actions">
              <button
                v-if="isStreaming"
                class="pane-btn stop-btn"
                @click="stopStreaming"
                title="停止生成"
              >
                ⏹ 停止
              </button>
              <span v-if="isStreaming" class="streaming-status">
                <span class="streaming-dot"></span>
                生成中...
              </span>
            </div>
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
          <div class="editor-statusbar">
            <span class="word-count">字数: {{ wordCount }}</span>
            <span class="step-info">{{ currentStepInfo }}</span>
            <span v-if="isDirty" class="save-status unsaved">未保存</span>
            <span v-else-if="isSaving" class="save-status saving">保存中...</span>
            <span v-else class="save-status saved">已保存</span>
          </div>
        </div>

        <!-- 预览区域 -->
        <div class="preview-pane" :class="{ 'mobile-mode': isMobilePreview }">
          <div class="pane-header">
            <span class="pane-title">👁️ 预览</span>
            <div class="pane-actions">
              <button
                class="pane-btn"
                :class="{ active: isMobilePreview }"
                @click="isMobilePreview = !isMobilePreview"
                title="切换手机端预览"
              >
                📱 手机端
              </button>
              <button
                class="pane-btn copy-btn"
                @click="copyToClipboard"
                title="复制公众号格式"
              >
                📋 复制到公众号
              </button>
            </div>
          </div>
          <div class="preview-content" ref="previewRef" :class="{ 'mobile-preview': isMobilePreview }">
            <div v-if="isMobilePreview" class="mobile-frame">
              <div class="mobile-screen">
                <div class="wechat-article" v-html="wechatFormattedContent"></div>
              </div>
            </div>
            <div v-else class="markdown-body" v-html="renderedContent"></div>
          </div>
        </div>
      </div>
    </div>

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
import WritingChatSidebar from './WritingChatSidebar.vue';
import SettingsModal, { type AIConfig } from './SettingsModal.vue';
import { useWails, FileError, type StreamCallback } from '../composables/useWails';

// 防抖工具函数
function debounce<T extends (...args: any[]) => void>(fn: T, delay: number): (...args: Parameters<T>) => void {
  let timeoutId: ReturnType<typeof setTimeout> | null = null;
  return (...args: Parameters<T>) => {
    if (timeoutId) clearTimeout(timeoutId);
    timeoutId = setTimeout(() => fn(...args), delay);
  };
}

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
const sidebarRef = ref<InstanceType<typeof WritingChatSidebar> | null>(null);
const editorRef = ref<HTMLTextAreaElement | null>(null);

// 文章状态
const articleTitle = ref('');
const articleOutline = ref('');
const editorContent = ref('');
const currentUUID = ref('');
const filePath = ref('');

// 写作流程状态
const currentStep = ref(0);
const isProcessing = ref(false);
const isStreaming = ref(false);
const streamingContent = ref('');
const streamUnsubscribe = ref<(() => void) | null>(null);

// UI 状态
const sidebarCollapsed = ref(false);
const showSettings = ref(false);
const errorMessage = ref('');
const isSaving = ref(false);
const isDirty = ref(false);
const isMobilePreview = ref(false);

// AI 配置
const aiConfig = ref<AIConfig>({
  baseUrl: 'https://api.deepseek.com/v1',
  token: '',
  temperature: 0.7,
  model: 'deepseek-chat',
});

const savedPositioning = ref('');

// ============================================
// 计算属性
// ============================================
const canSave = computed(() => {
  return isDirty.value && !isSaving.value && articleTitle.value !== '';
});

const wordCount = computed(() => {
  if (!editorContent.value) return 0;
  const chinese = (editorContent.value.match(/[\u4e00-\u9fa5]/g) || []).length;
  const english = editorContent.value
    .replace(/[\u4e00-\u9fa5]/g, '')
    .trim()
    .split(/\s+/)
    .filter((w) => w.length > 0).length;
  return chinese + english;
});

const renderedContent = computed(() => {
  return md.render(editorContent.value);
});

const wechatFormattedContent = computed(() => {
  let content = md.render(editorContent.value);
  content = content.replace(
    /<pre><code[^>]*>([\s\S]*?)<\/code><\/pre>/g,
    '<pre class="wechat-code-block"><code>$1</code></pre>'
  );
  content = content.replace(
    /<code>([^<]+)<\/code>(?![\s\S]*<\/pre>)/g,
    '<code class="wechat-inline-code">$1</code>'
  );
  return content;
});

const currentStepInfo = computed(() => {
  const stepNames = ['设定主题', '生成框架', '生成内容', '优化调整', '生成标题'];
  return `步骤 ${currentStep.value + 1}/5: ${stepNames[currentStep.value]}`;
});

// ============================================
// 监听编辑器内容变化（带防抖）
// ============================================
watch(editorContent, (newValue, oldValue) => {
  if (oldValue !== undefined && newValue !== oldValue) {
    isDirty.value = true;
    // 触发防抖自动保存
    debouncedAutoSave();
  }
});

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
      sidebarRef.value?.loadRecentArticles();
    }
  } catch (err) {
    if (err instanceof FileError && err.code === 'FILE_EXISTS_CONFIRM') {
      const filePathMatch = err.message.match(/文件已存在: (.+)$/);
      const existingFilePath = filePathMatch ? filePathMatch[1] : '';

      const confirmed = window.confirm(
        `文件 "${existingFilePath}" 已存在。\n\n是否覆盖保存？\n（旧版本将被替换）`
      );

      if (confirmed) {
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

// 防抖自动保存（3秒后自动保存）
const autoSave = async (overwrite = false): Promise<void> => {
  if (!articleTitle.value || !editorContent.value || !isDirty.value) return;

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
      sidebarRef.value?.loadRecentArticles();
    }
  } catch (err) {
    if (err instanceof FileError && err.code === 'FILE_EXISTS_CONFIRM') {
      await autoSave(true);
      return;
    }
    console.error('自动保存失败:', err);
  } finally {
    isSaving.value = false;
  }
};

const debouncedAutoSave = debounce(() => autoSave(), 3000);

// ============================================
// 对话消息处理
// ============================================
const handleSendMessage = async (message: string, step: number) => {
  // 解析用户输入
  const parsedInput = parseUserInput(message, step);
  
  if (step === 0) {
    // 步骤1: 设定主题
    await handleStep1_Setup(parsedInput);
  } else if (step === 1) {
    // 步骤2: 生成框架（大纲）
    await handleStep2_GenerateOutline(parsedInput);
  } else if (step === 2) {
    // 步骤3: 生成内容
    await handleStep3_GenerateArticle(parsedInput);
  } else if (step === 3) {
    // 步骤4: 优化调整
    await handleStep4_Optimize(parsedInput);
  } else if (step === 4) {
    // 步骤5: 生成标题
    await handleStep5_GenerateTitles(parsedInput);
  }
};

// 解析用户输入
const parseUserInput = (message: string, step: number) => {
  const result: {
    title?: string;
    requirements?: string;
    outline?: string;
    optimizationType?: string;
    optimizationRequest?: string;
  } = {};

  if (step === 0) {
    // 提取标题
    const titleMatch = message.match(/标题[：:]\s*(.+)/);
    if (titleMatch) {
      result.title = titleMatch[1].trim();
    } else {
      // 如果没有明确标记，使用第一行作为标题
      const firstLine = message.split('\n')[0].trim();
      if (firstLine.length < 50) {
        result.title = firstLine;
      }
    }

    // 提取要求
    const reqMatch = message.match(/要求[：:]\s*([\s\S]+)/);
    if (reqMatch) {
      result.requirements = reqMatch[1].trim();
    } else {
      result.requirements = message.replace(/标题[：:]\s*.+/, '').trim();
    }
  }

  return result;
};

// 步骤1: 设定主题
const handleStep1_Setup = async (input: { title?: string; requirements?: string }) => {
  if (input.title) {
    articleTitle.value = input.title;
  }

  // 添加 AI 回复
  sidebarRef.value?.addMessage(
    'assistant',
    `✅ 已记录您的创作需求\n\n**标题**: ${articleTitle.value || '待定'}\n\n**要求**: ${input.requirements || '无特殊要求'}\n\n接下来可以：\n1. 点击「生成大纲」按钮生成文章框架\n2. 继续补充更多要求`
  );

  // 自动进入下一步
  if (articleTitle.value) {
    currentStep.value = 1;
    sidebarRef.value?.nextStep();
  }
};

// 步骤2: 生成大纲（流式）
const handleStep2_GenerateOutline = async (input: any) => {
  if (!articleTitle.value) {
    sidebarRef.value?.addMessage('assistant', '⚠️ 请先输入文章标题');
    return;
  }

  isProcessing.value = true;
  isStreaming.value = true;
  streamingContent.value = '';

  // 添加流式消息
  sidebarRef.value?.addStreamingMessage('');

  // 清空编辑器并准备接收流式内容
  const separator = editorContent.value ? '\n\n---\n\n' : '';
  let receivedContent = '';

  const onChunk: StreamCallback = (chunk, done, errorMsg) => {
    if (errorMsg) {
      isStreaming.value = false;
      isProcessing.value = false;
      sidebarRef.value?.finishStreaming();
      sidebarRef.value?.addMessage('assistant', `❌ 生成大纲失败: ${errorMsg}`);
      return;
    }

    if (done) {
      isStreaming.value = false;
      isProcessing.value = false;
      articleOutline.value = receivedContent;
      sidebarRef.value?.finishStreaming();
      
      // 保存候选标题
      if ((chunk as any)?.titles) {
        (window as any).candidateTitles = (chunk as any).titles;
      }
      
      sidebarRef.value?.addMessage(
        'assistant',
        `✅ 大纲已生成完成！\n\n您可以：\n1. 在编辑区修改完善大纲\n2. 点击「生成文章」继续创作\n3. 提出修改意见重新生成`
      );

      // 自动进入下一步
      currentStep.value = 2;
      sidebarRef.value?.nextStep();
      
      // 自动保存
      autoSave();
      return;
    }

    // 累积内容
    receivedContent += chunk;
    streamingContent.value = receivedContent;
    
    // 实时更新编辑器
    editorContent.value = editorContent.value + separator + receivedContent;
    
    // 更新流式消息
    sidebarRef.value?.updateStreamingMessage(receivedContent);
    
    // 滚动到底部
    if (editorRef.value) {
      editorRef.value.scrollTop = editorRef.value.scrollHeight;
    }
  };

  // 启动流式生成
  streamUnsubscribe.value = wails.generateOutlineStream(
    articleTitle.value,
    input.requirements || '',
    savedPositioning.value,
    onChunk
  );
};

// 步骤3: 生成文章（流式）
const handleStep3_GenerateArticle = async (input: any) => {
  if (!articleTitle.value) {
    sidebarRef.value?.addMessage('assistant', '⚠️ 请先输入文章标题');
    return;
  }

  if (!articleOutline.value) {
    sidebarRef.value?.addMessage('assistant', '⚠️ 请先生成大纲');
    return;
  }

  isProcessing.value = true;
  isStreaming.value = true;
  streamingContent.value = '';

  // 添加流式消息
  sidebarRef.value?.addStreamingMessage('✍️ 正在撰写文章内容...\n\n');

  const combinedRequirements = [
    input.requirements,
    savedPositioning.value ? `公众号定位：${savedPositioning.value}` : '',
  ]
    .filter(Boolean)
    .join('\n\n');

  // 清空编辑器准备接收文章内容
  const savedOutline = articleOutline.value;
  editorContent.value = `# ${articleTitle.value}\n\n`;
  let receivedContent = '';

  const onChunk: StreamCallback = (chunk, done, errorMsg) => {
    if (errorMsg) {
      isStreaming.value = false;
      isProcessing.value = false;
      sidebarRef.value?.finishStreaming();
      sidebarRef.value?.addMessage('assistant', `❌ 生成文章失败: ${errorMsg}`);
      return;
    }

    if (done) {
      isStreaming.value = false;
      isProcessing.value = false;
      sidebarRef.value?.finishStreaming();
      
      sidebarRef.value?.addMessage(
        'assistant',
        `✅ 文章生成完成！\n\n您可以：\n1. 在编辑区修改完善\n2. 使用「优化调整」功能改进\n3. 继续生成更多内容`
      );

      // 自动进入下一步
      currentStep.value = 3;
      sidebarRef.value?.nextStep();
      
      // 自动保存
      autoSave();
      return;
    }

    // 累积内容
    receivedContent += chunk;
    streamingContent.value = receivedContent;
    
    // 实时更新编辑器（在标题后追加）
    editorContent.value = `# ${articleTitle.value}\n\n${receivedContent}`;
    
    // 更新流式消息（显示部分内容）
    const previewLength = Math.min(receivedContent.length, 200);
    const preview = receivedContent.substring(0, previewLength) + (receivedContent.length > 200 ? '...' : '');
    sidebarRef.value?.updateStreamingMessage(`✍️ 正在撰写文章内容...\n\n${preview}`);
    
    // 滚动到底部
    if (editorRef.value) {
      editorRef.value.scrollTop = editorRef.value.scrollHeight;
    }
  };

  // 启动流式生成
  streamUnsubscribe.value = wails.generateArticleStream(
    articleTitle.value,
    savedOutline,
    combinedRequirements,
    onChunk
  );
};

// 步骤4: 优化调整
const handleStep4_Optimize = async (input: any) => {
  const message = input.optimizationRequest || '请润色优化这篇文章';

  isProcessing.value = true;
  isStreaming.value = true;

  sidebarRef.value?.addStreamingMessage('💎 正在优化文章内容...');

  try {
    // 获取优化类型
    let optimizeType = 'polish';
    if (message.includes('精简')) optimizeType = 'simplify';
    else if (message.includes('扩写')) optimizeType = 'expand';
    else if (message.includes('案例')) optimizeType = 'example';

    // 调用优化 API（非流式）
    const optimizedContent = await wails.optimizeContent(
      editorContent.value,
      optimizeType,
      message
    );

    // 流式显示优化结果
    await streamTextToEditor('\n\n---\n\n**优化后**:\n\n' + optimizedContent, 'article');

    sidebarRef.value?.finishStreaming();
    sidebarRef.value?.addMessage(
      'assistant',
      `✅ 优化完成！\n\n您还可以：\n1. 继续提出优化意见\n2. 生成爆款标题\n3. 保存文章`
    );

    // 自动进入下一步
    currentStep.value = 4;
    sidebarRef.value?.nextStep();

    // 自动保存
    autoSave();
  } catch (err) {
    sidebarRef.value?.finishStreaming();
    sidebarRef.value?.addMessage('assistant', `❌ 优化失败: ${err instanceof Error ? err.message : '未知错误'}`);
  } finally {
    isProcessing.value = false;
    isStreaming.value = false;
  }
};

// 步骤5: 生成标题
const handleStep5_GenerateTitles = async (_input: any) => {
  isProcessing.value = true;

  sidebarRef.value?.addStreamingMessage('🎯 正在生成爆款标题...');

  try {
    const titles = await wails.generateViralTitles(
      editorContent.value,
      savedPositioning.value,
      5
    );

    const titlesText = titles.map((t: string, i: number) => `${i + 1}. ${t}`).join('\n');

    sidebarRef.value?.finishStreaming();
    sidebarRef.value?.addMessage(
      'assistant',
      `✅ 爆款标题生成完成！\n\n候选标题：\n${titlesText}\n\n您可以选择一个应用到文章中。`
    );
  } catch (err) {
    sidebarRef.value?.finishStreaming();
    sidebarRef.value?.addMessage('assistant', `❌ 生成标题失败: ${err instanceof Error ? err.message : '未知错误'}`);
  } finally {
    isProcessing.value = false;
  }
};

// ============================================
// 流式输出处理
// ============================================
const streamTextToEditor = async (text: string, _type: 'outline' | 'article') => {
  const chars = text.split('');
  let currentText = editorContent.value;
  
  for (let i = 0; i < chars.length; i++) {
    if (!isStreaming.value) {
      break;
    }
    
    currentText += chars[i];
    editorContent.value = currentText;
    
    if (editorRef.value) {
      editorRef.value.scrollTop = editorRef.value.scrollHeight;
    }
    
    await new Promise((resolve) => setTimeout(resolve, 5));
  }
  
  isDirty.value = true;
};

const stopStreaming = () => {
  isStreaming.value = false;
  isProcessing.value = false;

  // 取消事件订阅
  if (streamUnsubscribe.value) {
    try {
      streamUnsubscribe.value();
    } catch (e) {
      console.warn('取消流订阅失败:', e);
    }
    streamUnsubscribe.value = null;
  }

  // 完成流式消息
  sidebarRef.value?.finishStreaming();
};

// ============================================
// 快捷指令处理
// ============================================
const handleQuickAction = (action: string, _step: number) => {
  switch (action) {
    case 'generate-outline':
      handleStep2_GenerateOutline({});
      break;
    case 'generate-article':
      handleStep3_GenerateArticle({});
      break;
    case 'polish':
      handleStep4_Optimize({ optimizationRequest: '润色' });
      break;
    case 'expand':
      handleStep4_Optimize({ optimizationRequest: '扩写' });
      break;
    case 'simplify':
      handleStep4_Optimize({ optimizationRequest: '精简' });
      break;
    case 'add-example':
      handleStep4_Optimize({ optimizationRequest: '案例' });
      break;
    case 'generate-titles':
      handleStep5_GenerateTitles({});
      break;
  }
};

// ============================================
// 步骤变化处理
// ============================================
const handleStepChange = (step: number) => {
  currentStep.value = step;
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
      
      const outlineMatch = content.match(/##?\s+.+$/m);
      if (outlineMatch) {
        articleOutline.value = outlineMatch[0];
      }
      
      isDirty.value = false;
      currentStep.value = 0;
      
      if (sidebarRef.value) {
        sidebarRef.value.messages = [];
        sidebarRef.value.currentStep = 0;
      }
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
// 复制到公众号功能
// ============================================
const copyToClipboard = async (): Promise<void> => {
  try {
    const articleContent = wechatFormattedContent.value;
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = articleContent;
    applyInlineStyles(tempDiv);

    tempDiv.style.cssText = `
      font-family: 'JetBrains Mono', 'Fira Code', 'SF Mono', -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;
      font-size: 15px;
      line-height: 1.8;
      color: #1f2937;
      padding: 24px;
      background: #fafbfc;
      max-width: 100%;
    `;

    tempDiv.style.position = 'fixed';
    tempDiv.style.left = '-9999px';
    tempDiv.style.top = '-9999px';
    document.body.appendChild(tempDiv);

    const selection = window.getSelection();
    const range = document.createRange();
    range.selectNodeContents(tempDiv);
    selection?.removeAllRanges();
    selection?.addRange(range);

    document.execCommand('copy');

    selection?.removeAllRanges();
    document.body.removeChild(tempDiv);

    handleError('已复制到剪贴板，可直接粘贴到公众号编辑器');
  } catch (err) {
    console.error('复制失败:', err);
    handleError('复制失败，请手动复制');
  }
};

const applyInlineStyles = (element: HTMLElement): void => {
  const styleMap: Record<string, string> = {
    H1: 'font-size: 26px; font-weight: 700; margin: 30px 0 24px; line-height: 1.4; color: #1e40af; text-align: center; padding: 20px; background: linear-gradient(135deg, rgba(59, 130, 246, 0.08) 0%, rgba(139, 92, 246, 0.08) 100%); border: 1px solid rgba(59, 130, 246, 0.25); border-radius: 8px; position: relative; letter-spacing: 1px;',
    H2: 'font-size: 20px; font-weight: 600; margin: 32px 0 20px; padding: 12px 16px; background: linear-gradient(90deg, rgba(124, 58, 237, 0.1) 0%, transparent 100%); border-left: 4px solid #7c3aed; border-radius: 0 8px 8px 0; line-height: 1.4; color: #6d28d9; position: relative;',
    H3: 'font-size: 17px; font-weight: 600; margin: 26px 0 14px; padding: 8px 0 8px 14px; border-left: 3px solid #0891b2; color: #0e7490; background: linear-gradient(90deg, rgba(8, 145, 178, 0.08) 0%, transparent 100%);',
    H4: 'font-size: 15px; font-weight: 600; margin: 20px 0 12px; color: #ea580c; padding: 6px 0 6px 12px; border-left: 2px solid #ea580c; background: linear-gradient(90deg, rgba(234, 88, 12, 0.08) 0%, transparent 100%);',
    H5: 'font-size: 14px; font-weight: 600; margin: 16px 0 10px; color: #4b5563;',
    H6: 'font-size: 14px; font-weight: 600; margin: 16px 0 10px; color: #4b5563;',
    P: 'margin: 16px 0; text-align: justify; text-indent: 2em; color: #374151;',
    A: 'color: #2563eb; text-decoration: none; border-bottom: 1px solid rgba(37, 99, 235, 0.4);',
    IMG: 'max-width: 100%; height: auto; margin: 24px 0; border-radius: 8px; border: 1px solid rgba(59, 130, 246, 0.2); box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1); display: block;',
    BLOCKQUOTE: 'margin: 24px 0; padding: 20px 24px 20px 48px; background: #f3f4f6; border: 1px solid rgba(124, 58, 237, 0.2); border-radius: 8px; color: #4c1d95; position: relative; font-family: monospace;',
    UL: 'margin: 18px 0; padding-left: 24px; color: #374151;',
    OL: 'margin: 18px 0; padding-left: 24px; color: #374151;',
    LI: 'margin: 12px 0;',
    TABLE: 'width: 100%; border-collapse: separate; border-spacing: 0; margin: 24px 0; border-radius: 8px; overflow: hidden; border: 1px solid #e5e7eb; background: #ffffff; box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);',
    TH: 'background: linear-gradient(180deg, #f9fafb 0%, #f3f4f6 100%); color: #1e40af; font-weight: 600; padding: 14px 16px; text-align: left; border: none; border-bottom: 1px solid #e5e7eb; font-family: monospace; font-size: 13px;',
    TD: 'border: none; border-bottom: 1px solid #f3f4f6; padding: 12px 16px; text-align: left; color: #374151;',
    HR: 'border: none; height: 1px; background: linear-gradient(90deg, transparent 0%, rgba(59, 130, 246, 0.3) 20%, rgba(124, 58, 237, 0.4) 50%, rgba(59, 130, 246, 0.3) 80%, transparent 100%); margin: 36px 0; position: relative;',
    STRONG: 'font-weight: 700; color: #ea580c;',
    EM: 'font-style: italic; color: #6b7280;',
    CODE: 'background: #1e1e2e; padding: 2px 8px; border-radius: 4px; font-family: monospace; font-size: 0.88em; color: #a6e3a1; border: 1px solid #313244;',
    PRE: 'background: #0d1117; border: 1px solid #30363d; border-radius: 12px; padding: 20px; margin: 24px 0; overflow-x: auto; font-family: monospace; font-size: 13px; line-height: 1.7; color: #e6edf3; position: relative; box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);',
  };

  const tagName = element.tagName;
  if (styleMap[tagName]) {
    element.style.cssText = styleMap[tagName] + (element.style.cssText ? '; ' + element.style.cssText : '');
  }

  if (element.tagName === 'PRE' && element.classList.contains('wechat-code-block')) {
    element.style.cssText = 'background: #1e1e1e; border: none; border-radius: 8px; padding: 16px; margin: 20px 0; overflow-x: auto; font-family: "SF Mono", Monaco, Consolas, monospace; font-size: 14px; line-height: 1.6; color: #d4d4d4;';
  }

  if (element.tagName === 'CODE' && element.classList.contains('wechat-inline-code')) {
    element.style.cssText = 'background: #f1f3f4; padding: 2px 8px; border-radius: 4px; font-family: monospace; font-size: 0.9em; color: #e83e8c; border: 1px solid #e8e8e8;';
  }

  Array.from(element.children).forEach((child) => {
    applyInlineStyles(child as HTMLElement);
  });
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

// 页面关闭前警告未保存的更改
const handleBeforeUnload = (e: BeforeUnloadEvent) => {
  if (isDirty.value) {
    e.preventDefault();
    e.returnValue = '';
    return '';
  }
};

// ============================================
// 生命周期
// ============================================
onMounted(() => {
  document.addEventListener('keydown', handleKeyDown);
  window.addEventListener('beforeunload', handleBeforeUnload);
  loadAIConfig();
  loadPositioning();
});

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown);
  window.removeEventListener('beforeunload', handleBeforeUnload);

  // 停止任何进行中的流式生成
  if (isStreaming.value) {
    stopStreaming();
  }

  // 清理流式订阅
  if (streamUnsubscribe.value) {
    try {
      streamUnsubscribe.value();
    } catch (e) {
      console.warn('取消流订阅失败:', e);
    }
    streamUnsubscribe.value = null;
  }
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
  width: 340px;
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

.pane-actions {
  display: flex;
  gap: 6px;
  align-items: center;
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

.pane-btn.stop-btn {
  background: #ff4d4f;
  border-color: #ff4d4f;
  color: white;
}

.pane-btn.stop-btn:hover {
  background: #ff7875;
}

.streaming-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #667eea;
}

.streaming-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #667eea;
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.5;
    transform: scale(0.8);
  }
}

/* 编辑区域 */
.editor-content {
  flex: 1;
  overflow: auto;
}

/* 编辑器状态栏 */
.editor-statusbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 16px;
  background: var(--bg-hover, #f5f5f5);
  border-top: 1px solid var(--border-color, #e8e8e8);
  font-size: 12px;
  color: var(--text-secondary, #8c8c8c);
}

.word-count {
  font-weight: 500;
}

.step-info {
  padding: 2px 10px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 500;
}

.save-status {
  font-size: 11px;
}

.save-status.unsaved {
  color: var(--color-warning, #faad14);
}

.save-status.saving {
  color: var(--color-primary, #1890ff);
}

.save-status.saved {
  color: var(--text-muted, #bfbfbf);
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

.preview-content.mobile-preview {
  padding: 20px;
  background: #f5f5f5;
  display: flex;
  justify-content: center;
  align-items: flex-start;
}

.markdown-body {
  max-width: 800px;
  margin: 0 auto;
  font-size: 15px;
  line-height: 1.8;
}

/* 手机端预览框架 */
.mobile-frame {
  width: 375px;
  background: #1a1a1a;
  border-radius: 40px;
  padding: 15px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  position: relative;
}

.mobile-frame::before {
  content: '';
  display: block;
  width: 80px;
  height: 20px;
  background: #1a1a1a;
  border-radius: 10px;
  margin: 0 auto 15px;
}

.mobile-screen {
  width: 100%;
  background: #fff;
  border-radius: 30px;
  overflow: hidden;
  max-height: 70vh;
  overflow-y: auto;
}

/* 微信公众号文章样式 */
:deep(.wechat-article) {
  font-family: 'JetBrains Mono', 'Fira Code', 'SF Mono', -apple-system-font, BlinkMacSystemFont, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei UI', 'Microsoft YaHei', Arial, sans-serif;
  font-size: 15px;
  line-height: 1.8;
  color: #1f2937;
  padding: 24px;
  background: #fafbfc;
  background-image:
    linear-gradient(rgba(59, 130, 246, 0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(59, 130, 246, 0.04) 1px, transparent 1px);
  background-size: 24px 24px;
}

:deep(.wechat-article h1) {
  font-size: 26px;
  font-weight: 700;
  margin: 30px 0 24px;
  line-height: 1.4;
  color: #1e40af;
  text-align: center;
  padding: 20px;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.08) 0%, rgba(139, 92, 246, 0.08) 100%);
  border: 1px solid rgba(59, 130, 246, 0.25);
  border-radius: 8px;
  position: relative;
  letter-spacing: 1px;
}

:deep(.wechat-article h1::before) {
  content: '< ';
  color: #7c3aed;
}

:deep(.wechat-article h1::after) {
  content: ' />';
  color: #7c3aed;
}

:deep(.wechat-article h2) {
  font-size: 20px;
  font-weight: 600;
  margin: 32px 0 20px;
  padding: 12px 16px;
  background: linear-gradient(90deg, rgba(124, 58, 237, 0.1) 0%, transparent 100%);
  border-left: 4px solid #7c3aed;
  border-radius: 0 8px 8px 0;
  line-height: 1.4;
  color: #6d28d9;
  position: relative;
}

:deep(.wechat-article h2::before) {
  content: '#';
  color: #7c3aed;
  margin-right: 8px;
  opacity: 0.8;
}

:deep(.wechat-article h3) {
  font-size: 17px;
  font-weight: 600;
  margin: 26px 0 14px;
  padding: 8px 0 8px 14px;
  border-left: 3px solid #0891b2;
  color: #0e7490;
  background: linear-gradient(90deg, rgba(8, 145, 178, 0.08) 0%, transparent 100%);
}

:deep(.wechat-article h3::before) {
  content: '##';
  color: #0891b2;
  margin-right: 6px;
  opacity: 0.7;
  font-size: 0.85em;
}

:deep(.wechat-article h4) {
  font-size: 15px;
  font-weight: 600;
  margin: 20px 0 12px;
  color: #ea580c;
  padding: 6px 0 6px 12px;
  border-left: 2px solid #ea580c;
  background: linear-gradient(90deg, rgba(234, 88, 12, 0.08) 0%, transparent 100%);
}

:deep(.wechat-article h4::before) {
  content: '###';
  color: #ea580c;
  margin-right: 6px;
  opacity: 0.7;
  font-size: 0.8em;
}

:deep(.wechat-article h5),
:deep(.wechat-article h6) {
  font-size: 14px;
  font-weight: 600;
  margin: 16px 0 10px;
  color: #4b5563;
}

:deep(.wechat-article h5)::before {
  content: '####';
  color: #6b7280;
  margin-right: 6px;
  opacity: 0.6;
  font-size: 0.75em;
}

:deep(.wechat-article p) {
  margin: 16px 0;
  text-align: justify;
  text-indent: 2em;
  color: #374151;
}

:deep(.wechat-article a) {
  color: #2563eb;
  text-decoration: none;
  border-bottom: 1px solid rgba(37, 99, 235, 0.4);
}

:deep(.wechat-article img) {
  max-width: 100%;
  height: auto;
  margin: 24px 0;
  border-radius: 8px;
  border: 1px solid rgba(59, 130, 246, 0.2);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  display: block;
}

:deep(.wechat-article blockquote) {
  margin: 24px 0;
  padding: 20px 24px 20px 48px;
  background: #f3f4f6;
  border: 1px solid rgba(124, 58, 237, 0.2);
  border-radius: 8px;
  color: #4c1d95;
  position: relative;
  font-family: monospace;
}

:deep(.wechat-article blockquote::before) {
  content: '"';
  position: absolute;
  left: 16px;
  top: 12px;
  font-size: 36px;
  color: #7c3aed;
  font-family: Georgia, serif;
  line-height: 1;
  opacity: 0.6;
}

:deep(.wechat-article ul),
:deep(.wechat-article ol) {
  margin: 18px 0;
  padding-left: 24px;
  color: #374151;
}

:deep(.wechat-article li) {
  margin: 12px 0;
}

:deep(.wechat-article li::marker) {
  color: #7c3aed;
}

:deep(.wechat-article table) {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  margin: 24px 0;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e5e7eb;
  background: #ffffff;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

:deep(.wechat-article th) {
  background: linear-gradient(180deg, #f9fafb 0%, #f3f4f6 100%);
  color: #1e40af;
  font-weight: 600;
  padding: 14px 16px;
  text-align: left;
  border: none;
  border-bottom: 1px solid #e5e7eb;
  font-family: monospace;
  font-size: 13px;
}

:deep(.wechat-article td) {
  border: none;
  border-bottom: 1px solid #f3f4f6;
  padding: 12px 16px;
  text-align: left;
  color: #374151;
}

:deep(.wechat-article tr:hover td) {
  background: rgba(124, 58, 237, 0.03);
}

:deep(.wechat-article hr) {
  border: none;
  height: 1px;
  background: linear-gradient(90deg, transparent 0%, rgba(59, 130, 246, 0.3) 20%, rgba(124, 58, 237, 0.4) 50%, rgba(59, 130, 246, 0.3) 80%, transparent 100%);
  margin: 36px 0;
  position: relative;
}

:deep(.wechat-article hr::before) {
  content: '◆';
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  color: #7c3aed;
  font-size: 8px;
  background: #fafbfc;
  padding: 0 8px;
}

:deep(.wechat-article strong) {
  font-weight: 700;
  color: #ea580c;
}

:deep(.wechat-article em) {
  font-style: italic;
  color: #6b7280;
}

:deep(.wechat-article code) {
  background: #1e1e2e;
  padding: 2px 8px;
  border-radius: 4px;
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  font-size: 0.88em;
  color: #a6e3a1;
  border: 1px solid #313244;
}

:deep(.wechat-article pre) {
  background: #0d1117;
  border: 1px solid #30363d;
  border-radius: 12px;
  padding: 20px;
  margin: 24px 0;
  overflow-x: auto;
  font-family: 'SF Mono', Monaco, Consolas, monospace;
  font-size: 13px;
  line-height: 1.7;
  color: #e6edf3;
  position: relative;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
}

:deep(.wechat-article pre code) {
  background: transparent;
  padding: 0;
  border: none;
  color: inherit;
  font-size: inherit;
}

/* 错误提示 */
.error-toast {
  position: fixed;
  top: 60px;
  left: 50%;
  transform: translateX(-50%);
  padding: 12px 24px;
  background: #ff4d4f;
  color: white;
  border-radius: 6px;
  font-size: 13px;
  box-shadow: 0 4px 12px rgba(255, 77, 79, 0.3);
  z-index: 1000;
  cursor: pointer;
}

/* Toast 动画 */
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-20px);
}

/* 暗色主题适配 */
@media (prefers-color-scheme: dark) {
  .app {
    --bg-component: #141414;
    --text-primary: #d9d9d9;
    --text-secondary: #8c8c8c;
    --text-muted: #595959;
    --border-color: #303030;
    --bg-hover: #1f1f1f;
    --color-primary: #4ec9b0;
    --color-warning: #faad14;
  }

  .markdown-editor {
    background: #1f1f1f;
    color: #d9d9d9;
  }

  .step-info {
    background: linear-gradient(135deg, #4ec9b0 0%, #2d8a78 100%);
  }
}
</style>
