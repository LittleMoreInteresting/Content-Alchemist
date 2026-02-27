/**
 * Wails 运行时类型声明
 * 对应 Wails v2 自动生成的绑定
 */

import type { Article, ArticleWithContent, AIResponse, FileInfo } from './article';

/** Wails 运行时接口 */
export interface WailsRuntime {
  /** 打开文件对话框 */
  OpenFileDialog(): Promise<string>;

  /** 保存文件对话框 */
  SaveFileDialog(defaultFilename: string): Promise<string>;

  /** 打开目录对话框 */
  OpenDirectoryDialog(): Promise<string>;

  /** 读取文章 */
  ReadArticle(filePath: string): Promise<[Article, string]>;

  /** 保存文章 */
  SaveArticle(uuid: string, content: string): Promise<void>;

  /** 另存为 */
  SaveArticleAs(uuid: string, newPath: string, content: string): Promise<Article>;

  /** 创建新文章 */
  CreateNewArticle(): Promise<Article | null>;

  /** 创建新文章并指定初始内容 */
  CreateNewArticleWithContent(defaultFilename: string, content: string): Promise<Article | null>;

  /** 获取最近文章列表 */
  GetRecentArticles(limit: number): Promise<Article[]>;

  /** 更新文章元数据 */
  UpdateArticleMeta(uuid: string, tags: string[]): Promise<void>;

  /** 根据UUID获取文章 */
  GetArticleByUUID(uuid: string): Promise<Article>;

  /** 删除文章记录 */
  DeleteArticle(uuid: string): Promise<void>;

  /** 删除文章记录和文件 */
  DeleteArticleAndFile(uuid: string): Promise<void>;

  /** 检查文件是否存在 */
  CheckFileExists(filePath: string): Promise<boolean>;

  /** 获取文件信息 */
  GetFileInfo(filePath: string): Promise<FileInfo>;
}

/** Wails 事件接口 */
export interface WailsEvents {
  /** 监听事件 */
  On(eventName: string, callback: (data?: unknown) => void): () => void;

  /** 触发事件 */
  Emit(eventName: string, data?: unknown): void;
}

/** 扩展全局 Window 接口 */
declare global {
  interface Window {
    /** Wails Go 绑定 */
    go?: {
      main?: {
        App: WailsRuntime;
      };
    };

    /** Wails 运行时 */
    runtime?: {
      EventsOn: (eventName: string, callback: (data?: unknown) => void) => (() => void);
      EventsEmit: (eventName: string, data?: unknown) => void;
      LogInfo: (message: string) => void;
      LogWarning: (message: string) => void;
      LogError: (message: string) => void;
      LogDebug: (message: string) => void;
    };
  }
}

export {};
