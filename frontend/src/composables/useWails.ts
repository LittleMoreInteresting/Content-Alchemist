/**
 * Wails 后端调用封装
 * 使用 Wails 生成的绑定
 */

import { ref } from 'vue';
import * as App from '../../wailsjs/go/backend/App';
import type { models } from '../../wailsjs/go/models';

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

export type FileErrorCode =
  | 'FILE_NOT_FOUND'
  | 'PERMISSION_DENIED'
  | 'IS_DIRECTORY'
  | 'FILE_TOO_LARGE'
  | 'FILE_EXISTS'
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
  // 文件对话框操作
  // ============================================

  /**
   * 打开文件选择对话框
   */
  const openFileDialog = async (): Promise<string> => {
    return wrapAsync(async () => {
      const result = await App.OpenFileDialog();
      return result || '';
    });
  };

  /**
   * 打开保存文件对话框
   */
  const saveFileDialog = async (defaultFilename = ''): Promise<string> => {
    return wrapAsync(async () => {
      const result = await App.SaveFileDialog(defaultFilename);
      return result || '';
    });
  };

  /**
   * 打开目录选择对话框
   */
  const openDirectoryDialog = async (): Promise<string> => {
    return wrapAsync(async () => {
      const result = await App.OpenDirectoryDialog();
      return result || '';
    });
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

  /**
   * 创建新文章并指定初始内容
   */
  const createNewArticleWithContent = async (
    defaultFilename: string,
    content: string
  ): Promise<Article | null> => {
    return wrapAsync(async () => {
      return await App.CreateNewArticleWithContent(defaultFilename, content);
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
   * 更新文章元数据
   */
  const updateArticleMeta = async (
    uuid: string,
    tags: string[]
  ): Promise<void> => {
    return wrapAsync(async () => {
      await App.UpdateArticleMeta(uuid, tags);
    });
  };

  /**
   * 根据UUID获取文章
   */
  const getArticleByUUID = async (uuid: string): Promise<Article> => {
    return wrapAsync(async () => {
      return await App.GetArticleByUUID(uuid);
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

  /**
   * 删除文章记录和文件
   */
  const deleteArticleAndFile = async (uuid: string): Promise<void> => {
    return wrapAsync(async () => {
      await App.DeleteArticleAndFile(uuid);
    });
  };

  /**
   * 根据标题重命名文章文件
   * 返回新的文件路径
   */
  const renameArticleByTitle = async (
    uuid: string,
    newTitle: string,
    content: string
  ): Promise<string> => {
    return wrapAsync(async () => {
      return await App.RenameArticleByTitle(uuid, newTitle, content);
    });
  };

  // ============================================
  // 文件系统操作
  // ============================================

  /**
   * 检查文件是否存在
   */
  const checkFileExists = async (filePath: string): Promise<boolean> => {
    return wrapAsync(async () => {
      return await App.CheckFileExists(filePath);
    });
  };

  /**
   * 获取文件信息
   */
  const getFileInfo = async (filePath: string): Promise<FileInfo> => {
    return wrapAsync(async () => {
      return await App.GetFileInfo(filePath);
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
   * 智能保存文章
   * 如果未保存到本地文件，则根据标题自动生成文件名
   */
  const saveArticleWithSmartNaming = async (uuid: string, title: string, content: string): Promise<Article | null> => {
    return wrapAsync(async () => {
      return await App.SaveArticleWithSmartNaming(uuid, title, content);
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

    // 文件对话框
    openFileDialog,
    saveFileDialog,
    openDirectoryDialog,

    // 文件操作
    readArticle,
    saveArticle,
    saveArticleAs,
    createNewArticle,
    createNewArticleWithContent,

    // 文章元数据
    getRecentArticles,
    updateArticleMeta,
    getArticleByUUID,
    deleteArticle,
    deleteArticleAndFile,

    // 文件系统
    checkFileExists,
    getFileInfo,

    // 标题重命名
    renameArticleByTitle,

    // AI配置
    getAIConfig,
    saveAIConfig,

    // AI生成
    generateOutline,
    generateArticle,

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
