/**
 * 文章相关类型定义
 * 对应后端 models/models.go
 */

/** 文章元数据 */
export interface Article {
  id: number;
  uuid: string;
  filePath: string;
  title: string;
  summary: string;
  tags: string[];
  wordCount: number;
  createdAt: string; // ISO 8601 format
  updatedAt: string;
  lastOpenedAt: string;
}

/** 文章及内容 */
export interface ArticleWithContent {
  article: Article;
  content: string;
}

/** 编辑器设置 */
export interface EditorSettings {
  deepseekApiKey: string;
  editorTheme: string;
  fontSize: number;
  aiModel: string;
}

/** AI 请求 */
export interface AIRequest {
  selectedText: string;
  action: 'rewrite' | 'polish' | 'continue' | 'shorter' | 'code';
  context?: string;
}

/** AI 响应 */
export interface AIResponse {
  result: string;
  error?: string;
  tokensUsed?: number;
}

/** 文件过滤器 */
export interface FileFilter {
  displayName: string;
  pattern: string;
}

/** 文件对话框选项 */
export interface FileDialogOptions {
  defaultDirectory?: string;
  defaultFilename?: string;
  title?: string;
  filters?: FileFilter[];
}

/** 文件错误 */
export interface FileError {
  code: FileErrorCode;
  message: string;
}

/** 文件错误代码 */
export type FileErrorCode =
  | 'FILE_NOT_FOUND'
  | 'PERMISSION_DENIED'
  | 'IS_DIRECTORY'
  | 'FILE_TOO_LARGE'
  | 'FILE_EXISTS'
  | 'FILE_MODIFIED_EXTERNALLY'
  | 'ARTICLE_NOT_FOUND'
  | 'UNKNOWN_ERROR';

/** 通用结果包装器 */
export interface Result<T> {
  data: T;
  error?: string;
  success: boolean;
}

/** 文件信息 */
export interface FileInfo {
  exists: boolean;
  isDir?: boolean;
  size?: number;
  modTime?: number;
  path?: string;
}

/** 最近文章列表项 */
export interface RecentArticle {
  uuid: string;
  title: string;
  filePath: string;
  wordCount: number;
  lastOpenedAt: string;
}

/** 文章统计信息 */
export interface ArticleStats {
  wordCount: number;
  charCount: number;
  lineCount: number;
  readingTime: number; // 预计阅读时间（分钟）
}
