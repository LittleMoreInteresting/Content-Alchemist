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

// App 应用程序结构
type App struct {
	ctx             context.Context
	configService   *service.ConfigService
	articleService  *service.ArticleService
	materialService *service.MaterialService
}

// NewApp 创建新应用
func NewApp() *App {
	return &App{}
}

// startup 应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 初始化服务
	configService, err := service.NewConfigService()
	if err != nil {
		fmt.Printf("Init config service failed: %v\n", err)
	}
	a.configService = configService

	articleService, err := service.NewArticleService()
	if err != nil {
		fmt.Printf("Init article service failed: %v\n", err)
	}
	a.articleService = articleService

	materialService, err := service.NewMaterialService()
	if err != nil {
		fmt.Printf("Init material service failed: %v\n", err)
	}
	a.materialService = materialService
}

// Greet 问候方法（示例）
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// ============ Config API ============

// GetConfig 获取配置
func (a *App) GetConfig() (*model.Config, error) {
	return a.configService.GetConfig()
}

// SaveConfig 保存配置
func (a *App) SaveConfig(config model.Config) error {
	return a.configService.SaveConfig(config)
}

// HasConfig 检查是否有配置
func (a *App) HasConfig() bool {
	return a.configService.HasConfig()
}

// TestConnection 测试AI连接
func (a *App) TestConnection(apiKey, baseURL, modelName string) error {
	return a.configService.TestConnection(apiKey, baseURL, modelName)
}

// ============ Article API ============

// CreateArticle 创建文章
func (a *App) CreateArticle(title string) (*model.Article, error) {
	return a.articleService.CreateArticle(title)
}

// GetArticle 获取文章
func (a *App) GetArticle(id string) (*model.Article, error) {
	return a.articleService.GetArticle(id)
}

// SaveArticle 保存文章
func (a *App) SaveArticle(article model.Article) error {
	return a.articleService.SaveArticle(&article)
}

// ListArticles 列出文章
func (a *App) ListArticles() ([]model.Article, error) {
	return a.articleService.ListArticles()
}

// DeleteArticle 删除文章
func (a *App) DeleteArticle(id string) error {
	return a.articleService.DeleteArticle(id)
}

// ============ Material API ============

// ListMaterials 列出素材
func (a *App) ListMaterials(materialType string) ([]model.Material, error) {
	return a.materialService.ListMaterials(materialType)
}

// CreateMaterial 创建素材
func (a *App) CreateMaterial(material model.Material) error {
	return a.materialService.CreateMaterial(&material)
}

// DeleteMaterial 删除素材
func (a *App) DeleteMaterial(id string) error {
	return a.materialService.DeleteMaterial(id)
}

// ============ AI API ============

// GenerateOutline 生成大纲
func (a *App) GenerateOutline(title, style, audience string) ([]model.OutlineNode, error) {
	config, err := a.configService.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("get config failed: %w", err)
	}

	outlineAgent, err := agent.NewOutlineAgent(config.APIKey, config.APIBaseURL, config.Model)
	if err != nil {
		return nil, err
	}

	return outlineAgent.GenerateOutline(a.ctx, agent.OutlineRequest{
		Title:    title,
		Style:    style,
		Audience: audience,
	})
}

// StreamWriting 流式写作 - 返回完整内容（简化版）
func (a *App) StreamWriting(action, context, selectedText, position, style string) (string, error) {
	config, err := a.configService.GetConfig()
	if err != nil {
		return "", fmt.Errorf("get config failed: %w", err)
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
	})
	if err != nil {
		return "", err
	}

	return resp.Content, nil
}

// StreamWritingEvent 流式写作 - 使用 Wails 事件
// 前端通过 EventsOn("ai-stream") 监听流式输出
// 通过 EventsOn("ai-stream-done") 监听完成事件
func (a *App) StreamWritingEvent(action, context, selectedText, position, style string) error {
	config, err := a.configService.GetConfig()
	if err != nil {
		return fmt.Errorf("get config failed: %w", err)
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
