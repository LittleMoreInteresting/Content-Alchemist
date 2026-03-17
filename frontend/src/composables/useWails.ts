/**
 * Wails 后端调用封装
 * 使用 Wails 生成的绑定
 */

import { ref } from 'vue';
import * as App from '../../wailsjs/go/backend/App';
import type { models } from '../../wailsjs/go/models';
import { EventsOn } from '../../wailsjs/runtime';

// 导出 Article 类型
export type Article = models.Article;
export type ReadArticleResponse = models.ReadArticleResponse;
export type GenerateOutlineResult = {
  titles: string[];
  outline: string;
};
export type FileInfo = Record<string, any>;

// AI配置类型
export interface AIConfig {
  baseUrl: string;
  token: string;
  temperature: number;
  model: string;
}

// 流式回调类型
export type StreamCallback = (chunk: string, done: boolean, error?: string) => void;

export type FileErrorCode =
  | 'FILE_NOT_FOUND'
  | 'PERMISSION_DENIED'
  | 'IS_DIRECTORY'
  | 'FILE_TOO_LARGE'
  | 'FILE_EXISTS'
  | 'FILE_EXISTS_CONFIRM'
  | 'FILE_MODIFIED_EXTERNALLY'
  | 'ARTICLE_NOT_FOUND'
  | 'UNKNOWN_ERROR';

/** 文件错误类 */
export class FileError extends Error {
  code: FileErrorCode;

  constructor(code: FileErrorCode, message: string) {
    super(message);
    this.name = 'FileError';
    this.code = code;
  }
}

/** 包装错误 */
const wrapError = (error: unknown): FileError => {
  if (error instanceof FileError) {
    return error;
  }

  const message = error instanceof Error ? error.message : String(error);

  // 尝试解析Go错误
  if (message.includes('FILE_NOT_FOUND')) {
    return new FileError('FILE_NOT_FOUND', '文件不存在');
  }
  if (message.includes('PERMISSION_DENIED')) {
    return new FileError('PERMISSION_DENIED', '无权限访问文件');
  }
  if (message.includes('FILE_EXISTS_CONFIRM')) {
    return new FileError('FILE_EXISTS_CONFIRM', message);
  }
  if (message.includes('FILE_EXISTS')) {
    return new FileError('FILE_EXISTS', '文件已存在');
  }
  if (message.includes('FILE_MODIFIED_EXTERNALLY')) {
    return new FileError('FILE_MODIFIED_EXTERNALLY', '文件已被外部修改');
  }
  if (message.includes('ARTICLE_NOT_FOUND')) {
    return new FileError('ARTICLE_NOT_FOUND', '文章不存在');
  }

  return new FileError('UNKNOWN_ERROR', message);
};

/**
 * Wails 后端调用 Composable
 */
export function useWails() {
  /** 加载状态 */
  const isLoading = ref(false);
  /** 错误信息 */
  const error = ref<string | null>(null);

  /**
   * 包装异步操作
   */
  const wrapAsync = async <T>(fn: () => Promise<T>): Promise<T> => {
    isLoading.value = true;
    error.value = null;

    try {
      const result = await fn();
      return result;
    } catch (err) {
      const wrappedError = wrapError(err);
      error.value = wrappedError.message;
      throw wrappedError;
    } finally {
      isLoading.value = false;
    }
  };

  // ============================================
  // 文件操作
  // ============================================

  /**
   * 读取文章
   * 返回文章元数据和内容
   */
  const readArticle = async (
    filePath: string
  ): Promise<{ article: Article; content: string }> => {
    return wrapAsync(async () => {
      const response = await App.ReadArticle(filePath);
      if (!response.article) {
        throw new FileError('ARTICLE_NOT_FOUND', '文章不存在');
      }
      return {
        article: response.article,
        content: response.content,
      };
    });
  };

  /**
   * 保存文章
   */
  const saveArticle = async (uuid: string, content: string): Promise<void> => {
    return wrapAsync(async () => {
      await App.SaveArticle(uuid, content);
    });
  };

  /**
   * 另存为
   */
  const saveArticleAs = async (
    uuid: string,
    newPath: string,
    content: string
  ): Promise<Article> => {
    return wrapAsync(async () => {
      return await App.SaveArticleAs(uuid, newPath, content);
    });
  };

  /**
   * 创建新文章
   */
  const createNewArticle = async (): Promise<Article | null> => {
    return wrapAsync(async () => {
      return await App.CreateNewArticle();
    });
  };

  // ============================================
  // 文章元数据操作
  // ============================================

  /**
   * 获取最近文章列表
   */
  const getRecentArticles = async (limit = 10): Promise<Article[]> => {
    return wrapAsync(async () => {
      const articles = await App.GetRecentArticles(limit);
      return articles || [];
    });
  };

  /**
   * 删除文章记录（不删除文件）
   */
  const deleteArticle = async (uuid: string): Promise<void> => {
    return wrapAsync(async () => {
      await App.DeleteArticle(uuid);
    });
  };

  // ============================================
  // AI配置操作
  // ============================================

  /**
   * 获取AI配置
   */
  const getAIConfig = async (): Promise<AIConfig> => {
    return wrapAsync(async () => {
      return await App.GetAIConfig();
    });
  };

  /**
   * 保存AI配置
   */
  const saveAIConfig = async (config: AIConfig): Promise<void> => {
    return wrapAsync(async () => {
      await App.SaveAIConfig(config);
    });
  };

  /**
   * 根据标题生成大纲和候选标题
   * 返回包含3个候选标题和大纲的结果
   */
  const generateOutline = async (
    title: string,
    requirements?: string,
    positioning?: string
  ): Promise<{ titles: string[]; outline: string }> => {
    return wrapAsync(async () => {
      const result = await App.GenerateOutline(title, requirements || '', positioning || '');
      return {
        titles: result.titles || [],
        outline: result.outline || '',
      };
    });
  };

  /**
   * 根据大纲生成文章
   */
  const generateArticle = async (title: string, outline: string, requirements?: string): Promise<string> => {
    return wrapAsync(async () => {
      return await App.GenerateArticle(title, outline, requirements || '');
    });
  };

  /**
   * 优化文章内容
   * optimizeType: 优化类型 (polish:润色, expand:扩写, simplify:精简, example:添加案例)
   * requirements: 额外要求
   */
  const optimizeContent = async (
    content: string,
    optimizeType: string,
    requirements?: string
  ): Promise<string> => {
    return wrapAsync(async () => {
      return await App.OptimizeContent(content, optimizeType, requirements || '');
    });
  };

  /**
   * 生成爆款标题
   * content: 文章内容
   * positioning: 公众号定位
   * count: 生成数量（默认5个）
   */
  const generateViralTitles = async (
    content: string,
    positioning?: string,
    count?: number
  ): Promise<string[]> => {
    return wrapAsync(async () => {
      return await App.GenerateViralTitles(content, positioning || '', count || 5);
    });
  };

  // ============================================
  // 流式生成
  // ============================================

  /**
   * 流式生成大纲
   * @param title 标题
   * @param requirements 要求
   * @param positioning 定位
   * @param onChunk 数据块回调
   * @returns 取消订阅函数
   */
  const generateOutlineStream = (
    title: string,
    requirements: string,
    positioning: string,
    onChunk: StreamCallback
  ): (() => void) => {
    let isCancelled = false;
    let unsubscribe: (() => void) | null = null;

    // 包装回调以检查取消状态
    const wrappedCallback: StreamCallback = (chunk, done, error) => {
      if (!isCancelled) {
        onChunk(chunk, done, error);
      }
    };

    // 设置事件监听
    unsubscribe = EventsOn('outline:stream', (data: any) => {
      const { chunk, done, error: errMsg, titles, outline } = data;

      if (errMsg) {
        wrappedCallback('', true, errMsg);
      } else if (done) {
        // 完成时如果有完整数据，一并返回
        if (outline && titles) {
          wrappedCallback(outline, true);
        } else {
          wrappedCallback('', true);
        }
      } else {
        wrappedCallback(chunk || '', false);
      }
    });

    // 启动流式生成
    App.GenerateOutlineStream(title, requirements, positioning).catch((err: any) => {
      wrappedCallback('', true, err?.message || '流式生成失败');
    });

    // 返回取消函数
    return () => {
      isCancelled = true;
      if (unsubscribe) {
        unsubscribe();
        unsubscribe = null;
      }
    };
  };

  /**
   * 流式生成文章
   * @param title 标题
   * @param outline 大纲
   * @param requirements 要求
   * @param onChunk 数据块回调
   * @returns 取消订阅函数
   */
  const generateArticleStream = (
    title: string,
    outline: string,
    requirements: string,
    onChunk: StreamCallback
  ): (() => void) => {
    let isCancelled = false;
    let unsubscribe: (() => void) | null = null;

    // 包装回调以检查取消状态
    const wrappedCallback: StreamCallback = (chunk, done, error) => {
      if (!isCancelled) {
        onChunk(chunk, done, error);
      }
    };

    // 设置事件监听
    unsubscribe = EventsOn('article:stream', (data: any) => {
      const { chunk, done, error: errMsg } = data;

      if (errMsg) {
        wrappedCallback('', true, errMsg);
      } else {
        wrappedCallback(chunk || '', done);
      }
    });

    // 启动流式生成
    App.GenerateArticleStream(title, outline, requirements).catch((err: any) => {
      wrappedCallback('', true, err?.message || '流式生成失败');
    });

    // 返回取消函数
    return () => {
      isCancelled = true;
      if (unsubscribe) {
        unsubscribe();
        unsubscribe = null;
      }
    };
  };

  /**
   * 智能保存文章
   * 如果未保存到本地文件，则根据标题自动生成文件名并弹出对话框
   * overwrite: 是否强制覆盖已存在的文件
   */
  const saveArticleWithSmartNaming = async (
    uuid: string,
    title: string,
    content: string,
    overwrite = false
  ): Promise<Article | null> => {
    return wrapAsync(async () => {
      return await App.SaveArticleWithSmartNaming(uuid, title, content, overwrite);
    });
  };

  /**
   * 获取公众号定位配置
   */
  const getPositioning = async (): Promise<string> => {
    return wrapAsync(async () => {
      return await App.GetPositioning();
    });
  };

  /**
   * 保存公众号定位配置
   */
  const savePositioning = async (positioning: string): Promise<void> => {
    return wrapAsync(async () => {
      await App.SavePositioning(positioning);
    });
  };

  /**
   * 清除错误
   */
  const clearError = (): void => {
    error.value = null;
  };

  return {
    // 状态
    isLoading,
    error,

    // 文件操作
    readArticle,
    saveArticle,
    saveArticleAs,
    createNewArticle,

    // 文章元数据
    getRecentArticles,
    deleteArticle,

    // AI配置
    getAIConfig,
    saveAIConfig,

    // AI生成（非流式）
    generateOutline,
    generateArticle,
    optimizeContent,
    generateViralTitles,

    // 流式生成
    generateOutlineStream,
    generateArticleStream,

    // 智能保存
    saveArticleWithSmartNaming,

    // 公众号定位配置
    getPositioning,
    savePositioning,

    // 工具
    clearError,
  };
}

export default useWails;
