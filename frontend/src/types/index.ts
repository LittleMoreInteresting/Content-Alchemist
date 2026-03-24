
// 文章类型
export interface Article {
  id: string;
  title: string;
  content: string;
  outline: OutlineNode[];
  createdAt: string;
  updatedAt: string;
  status: 'draft' | 'reviewing' | 'ready' | 'published' | 'archived';
  sourceType?: 'manual' | 'workflow';
  workflowRunId?: string;
  topicId?: string;
  qualityScore?: number;
  readTime?: number;
  wordCount?: number;
  publishTaskId?: string;
  publishedAt?: string;
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

// ==================== 新增：工作流相关类型 ====================

// 工作流触发器
export interface WorkflowTrigger {
  type: 'manual' | 'schedule' | 'webhook' | 'rss';
  cronExpr?: string;
  rssUrl?: string;
  webhookUrl?: string;
  config?: Record<string, any>;
}

// 工作流步骤
export interface WorkflowStep {
  id: string;
  name: string;
  type: string;
  config?: Record<string, any>;
  nextStep?: string;
  onError?: 'retry' | 'skip' | 'abort' | 'manual';
}

// 工作流
export interface Workflow {
  id: string;
  name: string;
  description?: string;
  trigger: WorkflowTrigger;
  steps: WorkflowStep[];
  autoPublish: boolean;
  needReview: boolean;
  status: 'active' | 'paused' | 'archived';
  createdAt: string;
  updatedAt: string;
}

// 工作流运行步骤状态
export interface WorkflowRunStep {
  stepId: string;
  status: 'pending' | 'running' | 'completed' | 'failed' | 'skipped';
  input?: any;
  output?: any;
  startedAt?: string;
  completedAt?: string;
  error?: string;
  retryCount: number;
}

// 工作流运行
export interface WorkflowRun {
  id: string;
  workflowId: string;
  status: 'pending' | 'running' | 'paused' | 'completed' | 'failed' | 'cancelled';
  currentStep?: string;
  input?: Record<string, any>;
  output?: Record<string, any>;
  steps: WorkflowRunStep[];
  startedAt: string;
  completedAt?: string;
  error?: string;
}

// ==================== 新增：选题相关类型 ====================

// 热点趋势
export interface HotTrend {
  id: string;
  platform: string;
  title: string;
  url?: string;
  hotRank: number;
  hotValue: number;
  category?: string;
  createdAt: string;
}

// 选题
export interface Topic {
  id: string;
  title: string;
  category?: string;
  source: 'ai' | 'rss' | 'manual' | 'hot';
  sourceUrl?: string;
  score: number;
  hotScore?: number;
  compScore?: number;
  fitScore?: number;
  keywords?: string[];
  summary?: string;
  references?: string[];
  angles?: string[];
  reason?: string;
  status: 'pending' | 'approved' | 'rejected' | 'processing' | 'published' | 'archived';
  workflowRunId?: string;
  createdAt: string;
  updatedAt: string;
}

// ==================== 新增：发布相关类型 ====================

// 发布账号
export interface PublishAccount {
  id: string;
  platform: string;
  name: string;
  accountId?: string;
  accountType?: string;
  appId?: string;
  status: 'active' | 'expired' | 'disabled';
  lastUsedAt?: string;
  createdAt: string;
}

// 发布任务
export interface PublishTask {
  id: string;
  articleId: string;
  accountId: string;
  title?: string;
  content?: string;
  coverImage?: string;
  summary?: string;
  author?: string;
  tags?: string[];
  category?: string;
  originalUrl?: string;
  scheduleType: 'immediate' | 'scheduled';
  scheduleAt?: string;
  status: 'pending' | 'scheduled' | 'publishing' | 'published' | 'failed' | 'cancelled';
  platformId?: string;
  platformUrl?: string;
  error?: string;
  publishedAt?: string;
  retryCount: number;
  createdAt: string;
  updatedAt: string;
}
