/**
 * 编辑器状态管理 Composable
 * 管理当前编辑的文章状态和自动保存逻辑
 */

import { ref, computed, watch, onUnmounted } from 'vue';
import type { Article, EditorState } from '../types';
import { useWails, FileError } from './useWails';

/** 防抖时间（毫秒） */
const DEFAULT_SAVE_DEBOUNCE = 2000;
/** 最大自动重试次数 */
const MAX_RETRY_COUNT = 3;

/**
 * 编辑器状态管理
 */
export function useEditor(saveDebounce = DEFAULT_SAVE_DEBOUNCE) {
  const wails = useWails();

  // ============================================
  // 状态
  // ============================================

  /** 当前文章 */
  const article = ref<Article | null>(null);
  /** 当前内容 */
  const content = ref('');
  /** 是否有未保存更改 */
  const isDirty = ref(false);
  /** 是否正在保存 */
  const isSaving = ref(false);
  /** 保存错误 */
  const saveError = ref<string | null>(null);
  /** 最后保存时间 */
  const lastSavedAt = ref<Date | null>(null);
  /** 保存重试计数 */
  const retryCount = ref(0);

  // ============================================
  // 计算属性
  // ============================================

  /** 编辑器状态 */
  const state = computed<EditorState>(() => ({
    article: article.value,
    content: content.value,
    isDirty: isDirty.value,
    isSaving: isSaving.value,
    saveError: saveError.value,
    lastSavedAt: lastSavedAt.value,
  }));

  /** 标题 */
  const title = computed(() => article.value?.title || '未命名');

  /** 文件路径 */
  const filePath = computed(() => article.value?.filePath || '');

  /** 字数统计 */
  const wordCount = computed(() => {
    if (!content.value) return 0;
    // 中文字符 + 英文单词
    const chinese = (content.value.match(/[\u4e00-\u9fa5]/g) || []).length;
    const english = content.value
      .replace(/[\u4e00-\u9fa5]/g, '')
      .trim()
      .split(/\s+/)
      .filter(w => w.length > 0).length;
    return chinese + english;
  });

  /** 是否可保存 */
  const canSave = computed(() => isDirty.value && !isSaving.value && !!article.value);

  // ============================================
  // 保存逻辑
  // ============================================

  let saveTimer: ReturnType<typeof setTimeout> | null = null;

  /**
   * 执行保存
   */
  const doSave = async (): Promise<void> => {
    if (!article.value || !isDirty.value || isSaving.value) {
      return;
    }

    isSaving.value = true;
    saveError.value = null;

    try {
      await wails.saveArticle(article.value.uuid, content.value);

      // 保存成功
      isDirty.value = false;
      lastSavedAt.value = new Date();
      retryCount.value = 0;

      // 更新文章元数据
      article.value.wordCount = wordCount.value;
      article.value.updatedAt = new Date().toISOString();
    } catch (error) {
      console.error('保存失败:', error);

      let errorMessage = '保存失败';

      if (error instanceof FileError) {
        errorMessage = error.message;

        // 特殊处理：文件被外部修改
        if (error.code === 'FILE_MODIFIED_EXTERNALLY') {
          errorMessage = '文件已被外部程序修改，请使用"另存为"或刷新后重试';
        }
      } else if (error instanceof Error) {
        errorMessage = error.message;
      }

      saveError.value = errorMessage;

      // 自动重试
      if (retryCount.value < MAX_RETRY_COUNT) {
        retryCount.value++;
        scheduleSave(saveDebounce * retryCount.value);
      }
    } finally {
      isSaving.value = false;
    }
  };

  /**
   * 调度保存（防抖）
   */
  const scheduleSave = (delay = saveDebounce): void => {
    if (saveTimer) {
      clearTimeout(saveTimer);
    }

    if (isDirty.value) {
      saveTimer = setTimeout(() => {
        doSave();
      }, delay);
    }
  };

  /**
   * 立即保存
   */
  const saveNow = async (): Promise<void> => {
    if (saveTimer) {
      clearTimeout(saveTimer);
      saveTimer = null;
    }
    await doSave();
  };

  /**
   * 另存为
   */
  const saveAs = async (): Promise<Article | null> => {
    if (!article.value) return null;

    isSaving.value = true;
    saveError.value = null;

    try {
      const newArticle = await wails.saveArticleAs(
        article.value.uuid,
        '', // 由对话框选择
        content.value
      );

      article.value = newArticle;
      isDirty.value = false;
      lastSavedAt.value = new Date();

      return newArticle;
    } catch (error) {
      console.error('另存为失败:', error);
      saveError.value = error instanceof Error ? error.message : '另存为失败';
      return null;
    } finally {
      isSaving.value = false;
    }
  };

  // ============================================
  // 内容操作
  // ============================================

  /**
   * 设置内容
   */
  const setContent = (newContent: string): void => {
    if (content.value !== newContent) {
      content.value = newContent;
      isDirty.value = true;
      saveError.value = null;
      scheduleSave();
    }
  };

  /**
   * 加载文章
   */
  const loadArticle = async (filePath: string): Promise<void> => {
    isSaving.value = true;
    saveError.value = null;

    try {
      const { article: loadedArticle, content: loadedContent } =
        await wails.readArticle(filePath);

      article.value = loadedArticle;
      content.value = loadedContent;
      isDirty.value = false;
      lastSavedAt.value = new Date();
      retryCount.value = 0;
    } catch (error) {
      console.error('加载文章失败:', error);
      saveError.value = error instanceof Error ? error.message : '加载文章失败';
      throw error;
    } finally {
      isSaving.value = false;
    }
  };

  /**
   * 创建新文章
   */
  const createNew = async (): Promise<boolean> => {
    try {
      const newArticle = await wails.createNewArticle();
      if (newArticle) {
        article.value = newArticle;
        content.value = '';
        isDirty.value = false;
        lastSavedAt.value = new Date();
        return true;
      }
      return false;
    } catch (error) {
      console.error('创建新文章失败:', error);
      saveError.value = error instanceof Error ? error.message : '创建新文章失败';
      return false;
    }
  };

  /**
   * 关闭文章
   */
  const closeArticle = (): void => {
    if (saveTimer) {
      clearTimeout(saveTimer);
      saveTimer = null;
    }

    article.value = null;
    content.value = '';
    isDirty.value = false;
    isSaving.value = false;
    saveError.value = null;
    lastSavedAt.value = null;
    retryCount.value = 0;
  };

  // ============================================
  // 监听
  // ============================================

  // 监听内容变化（字数统计更新）
  watch(content, () => {
    if (article.value) {
      article.value.wordCount = wordCount.value;
    }
  });

  // 组件卸载时清理
  onUnmounted(() => {
    if (saveTimer) {
      clearTimeout(saveTimer);
    }
  });

  // ============================================
  // 导出
  // ============================================

  return {
    // 状态
    article,
    content,
    isDirty,
    isSaving,
    saveError,
    lastSavedAt,
    state,

    // 计算属性
    title,
    filePath,
    wordCount,
    canSave,

    // 方法
    setContent,
    loadArticle,
    createNew,
    closeArticle,
    saveNow,
    saveAs,
    scheduleSave,
  };
}

export default useEditor;
