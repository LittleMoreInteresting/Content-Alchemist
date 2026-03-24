package main

import (
	"context"
	"embed"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"content-alchemist/internal/agent"
	"content-alchemist/internal/model"
	"content-alchemist/internal/service"
)

//go:embed all:frontend/dist
var assets embed.FS
var version = "1.0.0"

// App 应用程序结构
type App struct {
	ctx             context.Context
	configService   *service.ConfigService
	articleService  *service.ArticleService
	materialService *service.MaterialService
	workflowService *service.WorkflowService
	topicService    *service.TopicService
}

// NewApp 创建新应用
func NewApp() *App {
	return &App{}
}

// startup 应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 初始化服务（如果密钥存在）
	configService, err := service.NewConfigService()
	if err != nil {
		fmt.Printf("Init config service failed: %v\n", err)
		// 密钥不存在，不中断，前端会检查并引导初始化
	} else {
		a.configService = configService
	}

	articleService, err := service.NewArticleService()
	if err != nil {
		fmt.Printf("Init article service failed: %v\n", err)
	} else {
		a.articleService = articleService
	}

	materialService, err := service.NewMaterialService()
	if err != nil {
		fmt.Printf("Init material service failed: %v\n", err)
	} else {
		a.materialService = materialService
	}

	// 初始化新服务（如果数据库已初始化）
	if articleService != nil {
		// 获取数据库连接
		db := articleService.GetDB()

		// 初始化工作流服务
		a.workflowService = service.NewWorkflowService(db)

		// 初始化选题服务（需要配置）
		if configService != nil {
			config, _ := configService.GetConfig()
			if config != nil && config.APIKey != "" {
				topicAgent, err := agent.NewTopicAgent(config.APIKey, config.APIBaseURL, config.Model)
				if err != nil {
					fmt.Printf("Init topic agent failed: %v\n", err)
				} else {
					a.topicService = service.NewTopicService(db, topicAgent)
				}
			}
		}
	}
}

// Greet 问候方法（示例）
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// ============ 密钥管理 API ============

// HasEncryptionKey 检查是否存在加密密钥
func (a *App) HasEncryptionKey() bool {
	return service.HasEncryptionKey()
}

// InitEncryptionKey 初始化加密密钥（首次使用）
func (a *App) InitEncryptionKey() error {
	if err := service.InitEncryptionKey(); err != nil {
		return err
	}

	// 初始化成功后，重新创建服务
	configService, err := service.NewConfigService()
	if err != nil {
		return fmt.Errorf("reinit config service failed: %w", err)
	}
	a.configService = configService

	articleService, err := service.NewArticleService()
	if err != nil {
		return fmt.Errorf("reinit article service failed: %w", err)
	}
	a.articleService = articleService

	materialService, err := service.NewMaterialService()
	if err != nil {
		return fmt.Errorf("reinit material service failed: %w", err)
	}
	a.materialService = materialService

	return nil
}

// ============ Config API ============

// GetConfig 获取配置
func (a *App) GetConfig() (*model.Config, error) {
	if a.configService == nil {
		return nil, fmt.Errorf("service not initialized")
	}
	config, err := a.configService.GetConfig()
	if err != nil {
		return nil, err
	}
	// 隐藏 API Key 的前端显示
	config.APIKey = maskAPIKey(config.APIKey)
	return config, nil
}

// maskAPIKey 隐藏 API Key 中间部分
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return key
	}
	return key[:4] + "****" + key[len(key)-4:]
}

// SaveConfig 保存配置
func (a *App) SaveConfig(config model.Config) error {
	if a.configService == nil {
		return fmt.Errorf("service not initialized")
	}

	// 如果传入的 API Key 是掩码格式，保留原有的 API Key
	if len(config.APIKey) > 4 && config.APIKey[4:8] == "****" {
		existingConfig, err := a.configService.GetConfig()
		if err == nil && existingConfig != nil {
			config.APIKey = existingConfig.APIKey
		}
	}

	return a.configService.SaveConfig(config)
}

// HasConfig 检查是否有配置
func (a *App) HasConfig() bool {
	if a.configService == nil {
		return false
	}
	return a.configService.HasConfig()
}

// TestConnection 测试AI连接
func (a *App) TestConnection(apiKey, baseURL, modelName string) error {
	if a.configService == nil {
		return fmt.Errorf("service not initialized")
	}
	return a.configService.TestConnection(apiKey, baseURL, modelName)
}

// ============ Article API ============

// CreateArticle 创建文章
func (a *App) CreateArticle(title string) (*model.Article, error) {
	if a.articleService == nil {
		return nil, fmt.Errorf("service not initialized")
	}
	return a.articleService.CreateArticle(title)
}

// GetArticle 获取文章
func (a *App) GetArticle(id string) (*model.Article, error) {
	if a.articleService == nil {
		return nil, fmt.Errorf("service not initialized")
	}
	return a.articleService.GetArticle(id)
}

// SaveArticle 保存文章
func (a *App) SaveArticle(article model.Article) error {
	if a.articleService == nil {
		return fmt.Errorf("service not initialized")
	}
	return a.articleService.SaveArticle(&article)
}

// ListArticles 列出文章
func (a *App) ListArticles() ([]model.Article, error) {
	if a.articleService == nil {
		return nil, fmt.Errorf("service not initialized")
	}
	return a.articleService.ListArticles()
}

// DeleteArticle 删除文章
func (a *App) DeleteArticle(id string) error {
	if a.articleService == nil {
		return fmt.Errorf("service not initialized")
	}
	return a.articleService.DeleteArticle(id)
}

// ============ Material API ============

// ListMaterials 列出素材
func (a *App) ListMaterials(materialType string) ([]model.Material, error) {
	if a.materialService == nil {
		return nil, fmt.Errorf("service not initialized")
	}
	return a.materialService.ListMaterials(materialType)
}

// CreateMaterial 创建素材
func (a *App) CreateMaterial(material model.Material) error {
	if a.materialService == nil {
		return fmt.Errorf("service not initialized")
	}
	return a.materialService.CreateMaterial(&material)
}

// DeleteMaterial 删除素材
func (a *App) DeleteMaterial(id string) error {
	if a.materialService == nil {
		return fmt.Errorf("service not initialized")
	}
	return a.materialService.DeleteMaterial(id)
}

// ============ AI API ============

// GenerateOutline 生成大纲
func (a *App) GenerateOutline(title, style, audience string) ([]model.OutlineNode, error) {
	if a.configService == nil {
		return nil, fmt.Errorf("service not initialized")
	}

	config, err := a.configService.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("get config failed: %w", err)
	}

	// 检查 API Key
	if config.APIKey == "" {
		return nil, fmt.Errorf("api key is empty, please configure in settings")
	}

	// 日志记录（隐藏部分 key）
	maskedKey := config.APIKey
	if len(config.APIKey) > 8 {
		maskedKey = config.APIKey[:4] + "****" + config.APIKey[len(config.APIKey)-4:]
	}
	fmt.Printf("GenerateOutline: using API Key=%s, BaseURL=%s, Model=%s\n", maskedKey, config.APIBaseURL, config.Model)

	outlineAgent, err := agent.NewOutlineAgent(config.APIKey, config.APIBaseURL, config.Model)
	if err != nil {
		return nil, fmt.Errorf("create outline agent failed: %w", err)
	}

	return outlineAgent.GenerateOutline(a.ctx, agent.OutlineRequest{
		Title:    title,
		Style:    style,
		Audience: audience,
	})
}

// StreamWriting 流式写作 - 返回完整内容（简化版）
func (a *App) StreamWriting(action, context, selectedText, position, style, customPrompt string) (string, error) {
	if a.configService == nil {
		return "", fmt.Errorf("service not initialized")
	}

	config, err := a.configService.GetConfig()
	if err != nil {
		return "", fmt.Errorf("get config failed: %w", err)
	}

	if config.APIKey == "" {
		return "", fmt.Errorf("api key is empty, please configure in settings")
	}

	writingAgent, err := agent.NewWritingAgent(config.APIKey, config.APIBaseURL, config.Model)
	if err != nil {
		return "", err
	}

	// 使用 Execute 获取完整内容
	resp, err := writingAgent.Execute(a.ctx, agent.WritingRequest{
		Action:       action,
		Context:      context,
		SelectedText: selectedText,
		Position:     position,
		Style:        style,
		CustomPrompt: customPrompt,
	})
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}

// StreamWritingEvent 流式写作 - 使用 Wails 事件
// 前端通过 EventsOn("ai-stream") 监听流式输出
// 通过 EventsOn("ai-stream-done") 监听完成事件
func (a *App) StreamWritingEvent(action, context, selectedText, position, style, customPrompt string) error {
	if a.configService == nil {
		return fmt.Errorf("service not initialized")
	}

	config, err := a.configService.GetConfig()
	if err != nil {
		return fmt.Errorf("get config failed: %w", err)
	}

	if config.APIKey == "" {
		return fmt.Errorf("api key is empty, please configure in settings")
	}

	writingAgent, err := agent.NewWritingAgent(config.APIKey, config.APIBaseURL, config.Model)
	if err != nil {
		return err
	}

	// 生成唯一请求ID
	requestID := fmt.Sprintf("stream-%d", time.Now().UnixNano())

	// 启动协程进行流式输出
	go func() {
		err := writingAgent.StreamExecute(a.ctx, agent.WritingRequest{
			Action:       action,
			Context:      context,
			SelectedText: selectedText,
			Position:     position,
			Style:        style,
			CustomPrompt: customPrompt,
		}, func(chunk string) {
			// 发送流式数据到前端
			runtime.EventsEmit(a.ctx, "ai-stream", map[string]string{
				"requestId": requestID,
				"chunk":     chunk,
			})
		})

		if err != nil {
			runtime.EventsEmit(a.ctx, "ai-stream-error", map[string]string{
				"requestId": requestID,
				"error":     err.Error(),
			})
		} else {
			runtime.EventsEmit(a.ctx, "ai-stream-done", map[string]string{
				"requestId": requestID,
			})
		}
	}()

	return nil
}

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "Content-Alchemist",
		Width:     1440,
		Height:    900,
		MinWidth:  1024,
		MinHeight: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// ============ Workflow API ============

// CreateWorkflow 创建工作流
func (a *App) CreateWorkflow(workflow model.Workflow) (*model.Workflow, error) {
	if a.workflowService == nil {
		return nil, fmt.Errorf("workflow service not initialized")
	}
	if err := a.workflowService.CreateWorkflow(&workflow); err != nil {
		return nil, err
	}
	return &workflow, nil
}

// GetWorkflow 获取工作流
func (a *App) GetWorkflow(id string) (*model.Workflow, error) {
	if a.workflowService == nil {
		return nil, fmt.Errorf("workflow service not initialized")
	}
	return a.workflowService.GetWorkflow(id)
}

// ListWorkflows 列出所有工作流
func (a *App) ListWorkflows() ([]model.Workflow, error) {
	if a.workflowService == nil {
		return nil, fmt.Errorf("workflow service not initialized")
	}
	return a.workflowService.ListWorkflows()
}

// UpdateWorkflow 更新工作流
func (a *App) UpdateWorkflow(workflow model.Workflow) error {
	if a.workflowService == nil {
		return fmt.Errorf("workflow service not initialized")
	}
	return a.workflowService.UpdateWorkflow(&workflow)
}

// DeleteWorkflow 删除工作流
func (a *App) DeleteWorkflow(id string) error {
	if a.workflowService == nil {
		return fmt.Errorf("workflow service not initialized")
	}
	return a.workflowService.DeleteWorkflow(id)
}

// StartWorkflow 启动工作流
func (a *App) StartWorkflow(workflowID string, input map[string]interface{}) (*model.WorkflowRun, error) {
	if a.workflowService == nil {
		return nil, fmt.Errorf("workflow service not initialized")
	}
	return a.workflowService.StartWorkflow(a.ctx, workflowID, input)
}

// GetWorkflowRun 获取工作流运行状态
func (a *App) GetWorkflowRun(runID string) (*model.WorkflowRun, error) {
	if a.workflowService == nil {
		return nil, fmt.Errorf("workflow service not initialized")
	}
	return a.workflowService.GetWorkflowRun(runID)
}

// ListWorkflowRuns 列出工作流运行实例
func (a *App) ListWorkflowRuns(workflowID string) ([]model.WorkflowRun, error) {
	if a.workflowService == nil {
		return nil, fmt.Errorf("workflow service not initialized")
	}
	return a.workflowService.ListWorkflowRuns(workflowID)
}

// CancelWorkflow 取消工作流运行
func (a *App) CancelWorkflow(runID string) error {
	if a.workflowService == nil {
		return fmt.Errorf("workflow service not initialized")
	}
	return a.workflowService.CancelWorkflow(runID)
}

// GetAgentRegistry 获取Agent注册表
func (a *App) GetAgentRegistry() []map[string]interface{} {
	if a.workflowService == nil {
		return nil
	}
	registry := a.workflowService.GetEngine().GetRegistry()
	agents := registry.List()

	var result []map[string]interface{}
	for _, agent := range agents {
		result = append(result, map[string]interface{}{
			"type":        agent.Type,
			"name":        agent.Name,
			"description": agent.Description,
			"icon":        agent.Icon,
			"category":    agent.Category,
		})
	}
	return result
}

// ============ Topic API ============

// CreateTopic 创建选题
func (a *App) CreateTopic(topic model.Topic) (*model.Topic, error) {
	if a.topicService == nil {
		return nil, fmt.Errorf("topic service not initialized")
	}
	if err := a.topicService.CreateTopic(&topic); err != nil {
		return nil, err
	}
	return &topic, nil
}

// GetTopic 获取选题
func (a *App) GetTopic(id string) (*model.Topic, error) {
	if a.topicService == nil {
		return nil, fmt.Errorf("topic service not initialized")
	}
	return a.topicService.GetTopic(id)
}

// ListTopics 列出选题
func (a *App) ListTopics(status string, limit int) ([]model.Topic, error) {
	if a.topicService == nil {
		return nil, fmt.Errorf("topic service not initialized")
	}
	return a.topicService.ListTopics(status, limit)
}

// UpdateTopic 更新选题
func (a *App) UpdateTopic(topic model.Topic) error {
	if a.topicService == nil {
		return fmt.Errorf("topic service not initialized")
	}
	return a.topicService.UpdateTopic(&topic)
}

// DeleteTopic 删除选题
func (a *App) DeleteTopic(id string) error {
	if a.topicService == nil {
		return fmt.Errorf("topic service not initialized")
	}
	return a.topicService.DeleteTopic(id)
}

// ApproveTopic 批准选题
func (a *App) ApproveTopic(id string) error {
	if a.topicService == nil {
		return fmt.Errorf("topic service not initialized")
	}
	return a.topicService.ApproveTopic(id)
}

// RejectTopic 拒绝选题
func (a *App) RejectTopic(id string) error {
	if a.topicService == nil {
		return fmt.Errorf("topic service not initialized")
	}
	return a.topicService.RejectTopic(id)
}

// GetHotTrends 获取热点趋势
func (a *App) GetHotTrends(platform string, limit int) ([]model.HotTrend, error) {
	if a.topicService == nil {
		return nil, fmt.Errorf("topic service not initialized")
	}
	return a.topicService.GetHotTrends(platform, limit)
}

// FetchAndSaveHotTrends 获取并保存热点趋势
func (a *App) FetchAndSaveHotTrends(platforms []string, limit int) ([]model.HotTrend, error) {
	if a.topicService == nil {
		return nil, fmt.Errorf("topic service not initialized")
	}

	trends, err := a.topicService.FetchHotTrends(a.ctx, platforms, limit)
	if err != nil {
		return nil, err
	}

	if err := a.topicService.SaveHotTrends(trends); err != nil {
		return nil, err
	}

	return trends, nil
}

// GetHotTrendPlatforms 获取热点平台列表
func (a *App) GetHotTrendPlatforms() ([]string, error) {
	if a.topicService == nil {
		return nil, fmt.Errorf("topic service not initialized")
	}
	return a.topicService.GetHotTrendPlatforms()
}

// SearchTopics 搜索选题
func (a *App) SearchTopics(keyword string) ([]model.Topic, error) {
	if a.topicService == nil {
		return nil, fmt.Errorf("topic service not initialized")
	}
	return a.topicService.SearchTopics(keyword)
}

// GetHotPlatforms 获取热点平台列表
func (a *App) GetHotPlatforms() []string {
	if a.topicService == nil {
		return nil
	}
	return a.topicService.GetHotPlatforms()
}

// FetchHotTrendsRealtime 实时获取热点（不保存）
func (a *App) FetchHotTrendsRealtime(platforms []string, limit int) ([]model.HotTrend, error) {
	if a.topicService == nil {
		return nil, fmt.Errorf("topic service not initialized")
	}
	return a.topicService.FetchHotTrends(a.ctx, platforms, limit)
}

// AIGenerateTopicsFromHot AI基于热点生成选题
func (a *App) AIGenerateTopicsFromHot(platforms []string, limit int) ([]*model.Topic, error) {
	if a.topicService == nil {
		return nil, fmt.Errorf("topic service not initialized")
	}
	return a.topicService.AIGenerateTopicsFromHot(a.ctx, platforms, limit)
}

// CreateArticleFromTopic 基于选题创建文章并启动工作流
func (a *App) CreateArticleFromTopic(topicID string, workflowID string) (*model.Article, *model.WorkflowRun, error) {
	if a.topicService == nil || a.workflowService == nil || a.articleService == nil {
		return nil, nil, fmt.Errorf("service not initialized")
	}
	
	// 1. 获取选题
	topic, err := a.topicService.GetTopic(topicID)
	if err != nil {
		return nil, nil, err
	}
	
	// 2. 创建文章
	article, err := a.articleService.CreateArticle(topic.Title)
	if err != nil {
		return nil, nil, err
	}
	
	// 3. 关联选题
	article.TopicID = topicID
	article.SourceType = model.ArticleSourceWorkflow
	if err := a.articleService.SaveArticle(article); err != nil {
		return nil, nil, err
	}
	
	// 4. 更新选题状态
	topic.Status = model.TopicStatusProcessing
	if err := a.topicService.UpdateTopic(topic); err != nil {
		// 非致命错误
	}
	
	// 5. 启动工作流
	input := map[string]interface{}{
		"topic_id":    topicID,
		"article_id":  article.ID,
		"title":       topic.Title,
		"keywords":    topic.Keywords,
		"angles":      topic.Angles,
	}
	
	run, err := a.workflowService.StartWorkflow(a.ctx, workflowID, input)
	if err != nil {
		return article, nil, err
	}
	
	// 6. 更新文章关联工作流
	article.WorkflowRunID = run.ID
	if err := a.articleService.SaveArticle(article); err != nil {
		// 非致命错误
	}
	
	return article, run, nil
}
