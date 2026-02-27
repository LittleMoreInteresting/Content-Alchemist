/**
 * 类型定义统一导出
 */

export * from './article';
export * from './wails';

import type { Article, Result } from './article';

/** 可空类型 */
export type Nullable<T> = T | null;

/** 可选类型 */
export type Optional<T> = T | undefined;

/** 异步函数返回类型 */
export type AsyncResult<T> = Promise<Result<T>>;

/** 通用错误类型 */
export interface AppError {
  code: string;
  message: string;
  details?: unknown;
}

/** 编辑器状态 */
export interface EditorState {
  article: Article | null;
  content: string;
  isDirty: boolean;
  isSaving: boolean;
  saveError: string | null;
  lastSavedAt: Date | null;
}

/** 通知类型 */
export interface Notification {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  message: string;
  duration?: number;
}
