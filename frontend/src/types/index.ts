// 文章类型
export interface Article {
  id: string;
  title: string;
  content: string;
  outline: OutlineNode[];
  createdAt: string;
  updatedAt: string;
  status: 'draft' | 'published';
}

// 大纲节点
export interface OutlineNode {
  id: string;
  level: 1 | 2 | 3;
  title: string;
  content?: string;
  parentId?: string;
  status: 'empty' | 'draft' | 'done';
  wordCount: number;
  targetWords: number;
}

// 配置类型
export interface Config {
  id: string;
  apiBaseUrl: string;
  apiKey: string;
  model: string;
  temperature: number;
  styleTags: string[];
  audience: string;
  persona: string;
  createdAt: string;
  updatedAt: string;
}

// 素材类型
export interface Material {
  id: string;
  type: 'snippet' | 'data' | 'quote' | 'history';
  title: string;
  content: string;
  tags: string[];
  source: string;
  createdAt: string;
  usageCount: number;
}

// 版本类型
export interface Version {
  id: string;
  articleId: string;
  content: string;
  snapshot: string;
  createdAt: string;
  type: 'auto' | 'manual';
}

// AI 写作请求
export interface WritingRequest {
  action: 'expand' | 'polish' | 'shorten' | 'continue' | 'title';
  context: string;
  selectedText: string;
  position: 'before' | 'after' | 'replace';
  style: string;
}

// AI 写作响应
export interface WritingResponse {
  content: string;
  suggestions?: string[];
}

// 大纲生成请求
export interface OutlineRequest {
  title: string;
  style: string;
  audience: string;
}

// 命令面板命令
export interface Command {
  id: string;
  name: string;
  description: string;
  icon: string;
  shortcut?: string;
  action: () => void;
}
