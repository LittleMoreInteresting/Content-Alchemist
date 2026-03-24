import { Topic, HotTrend } from '../types';

// 创建选题
export async function CreateTopic(topic: Topic): Promise<Topic> {
  return await window.go.main.App.CreateTopic(topic);
}

// 获取选题
export async function GetTopic(id: string): Promise<Topic> {
  return await window.go.main.App.GetTopic(id);
}

// 列出选题
export async function ListTopics(status: string = '', limit: number = 50): Promise<Topic[]> {
  return await window.go.main.App.ListTopics(status, limit);
}

// 更新选题
export async function UpdateTopic(topic: Topic): Promise<void> {
  return await window.go.main.App.UpdateTopic(topic);
}

// 删除选题
export async function DeleteTopic(id: string): Promise<void> {
  return await window.go.main.App.DeleteTopic(id);
}

// 批准选题
export async function ApproveTopic(id: string): Promise<void> {
  return await window.go.main.App.ApproveTopic(id);
}

// 拒绝选题
export async function RejectTopic(id: string): Promise<void> {
  return await window.go.main.App.RejectTopic(id);
}

// 搜索选题
export async function SearchTopics(keyword: string): Promise<Topic[]> {
  return await window.go.main.App.SearchTopics(keyword);
}

// 获取热点趋势
export async function GetHotTrends(platform: string, limit: number = 20): Promise<HotTrend[]> {
  return await window.go.main.App.GetHotTrends(platform, limit);
}

// 实时获取热点（不保存）
export async function FetchHotTrendsRealtime(platforms: string[], limit: number = 20): Promise<HotTrend[]> {
  return await window.go.main.App.FetchHotTrendsRealtime(platforms, limit);
}

// 获取并保存热点趋势
export async function FetchAndSaveHotTrends(platforms: string[], limit: number = 20): Promise<HotTrend[]> {
  return await window.go.main.App.FetchAndSaveHotTrends(platforms, limit);
}

// 获取热点平台列表
export async function GetHotPlatforms(): Promise<string[]> {
  return await window.go.main.App.GetHotPlatforms();
}

// AI基于热点生成选题
export async function AIGenerateTopicsFromHot(platforms: string[], limit: number = 5): Promise<Topic[]> {
  return await window.go.main.App.AIGenerateTopicsFromHot(platforms, limit);
}

// 基于选题创建文章并启动工作流
export async function CreateArticleFromTopic(topicId: string, workflowId: string): Promise<{article: any, run: any}> {
  return await window.go.main.App.CreateArticleFromTopic(topicId, workflowId);
}
