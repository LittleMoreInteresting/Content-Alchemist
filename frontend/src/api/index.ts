/**
 * Wails API 封装
 * 提供与 Go 后端交互的方法
 */

import {
  // 密钥管理
  HasEncryptionKey,
  InitEncryptionKey,
  // 配置
  GetConfig,
  SaveConfig,
  HasConfig,
  TestConnection,
  // 文章
  CreateArticle,
  GetArticle,
  SaveArticle,
  ListArticles,
  DeleteArticle,
  // 素材
  ListMaterials,
  CreateMaterial,
  DeleteMaterial,
  // AI
  GenerateOutline,
  StreamWriting
} from '../../wailsjs/go/main/App'

import type * as GoModels from '../../wailsjs/go/models'

// ===== 密钥管理 API =====

export async function checkHasEncryptionKey(): Promise<boolean> {
  return await HasEncryptionKey()
}

export async function initEncryptionKey(): Promise<void> {
  return await InitEncryptionKey()
}

// ===== Config API =====

export async function fetchConfig(): Promise<GoModels.model.Config> {
  return await GetConfig()
}

export async function saveConfigToBackend(config: any): Promise<void> {
  return await SaveConfig(config as GoModels.model.Config)
}

export async function checkHasConfig(): Promise<boolean> {
  return await HasConfig()
}

export async function testAPIConnection(
  apiKey: string,
  baseURL: string,
  modelName: string
): Promise<void> {
  return await TestConnection(apiKey, baseURL, modelName)
}

// ===== Article API =====

export async function createNewArticle(title: string): Promise<GoModels.model.Article> {
  return await CreateArticle(title)
}

export async function fetchArticle(id: string): Promise<GoModels.model.Article> {
  return await GetArticle(id)
}

export async function saveArticleToBackend(article: any): Promise<void> {
  return await SaveArticle(article as GoModels.model.Article)
}

export async function fetchAllArticles(): Promise<GoModels.model.Article[]> {
  return await ListArticles()
}

export async function removeArticle(id: string): Promise<void> {
  return await DeleteArticle(id)
}

// ===== Material API =====

export async function fetchMaterials(type: string): Promise<GoModels.model.Material[]> {
  return await ListMaterials(type)
}

export async function addMaterial(material: any): Promise<void> {
  return await CreateMaterial(material as GoModels.model.Material)
}

export async function removeMaterial(id: string): Promise<void> {
  return await DeleteMaterial(id)
}

// ===== AI API =====

export async function generateOutlineFromAI(
  title: string,
  style: string,
  audience: string
): Promise<GoModels.model.OutlineNode[]> {
  return await GenerateOutline(title, style, audience)
}

export async function streamWriteWithAI(
  action: string,
  context: string,
  selectedText: string,
  position: string,
  style: string
): Promise<string> {
  return await StreamWriting(action, context, selectedText, position, style)
}
