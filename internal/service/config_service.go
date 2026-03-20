package service

import (
	"context"
	"fmt"
	"time"

	einoModel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/eino-ext/components/model/openai"

	"content-alchemist/internal/model"
	"content-alchemist/internal/repository"
)

// ConfigService 配置服务
type ConfigService struct {
	db *repository.DB
}

// NewConfigService 创建配置服务
func NewConfigService() (*ConfigService, error) {
	db, err := repository.NewDB()
	if err != nil {
		return nil, fmt.Errorf("init db failed: %w", err)
	}

	return &ConfigService{db: db}, nil
}

// SaveConfig 保存配置（自动加密API Key）
func (s *ConfigService) SaveConfig(config model.Config) error {
	config.UpdatedAt = time.Now()
	if config.CreatedAt.IsZero() {
		config.CreatedAt = time.Now()
	}
	return s.db.SaveConfig(&config)
}

// GetConfig 获取配置（自动解密API Key）
func (s *ConfigService) GetConfig() (*model.Config, error) {
	return s.db.GetConfig()
}

// HasConfig 检查是否已有配置
func (s *ConfigService) HasConfig() bool {
	return s.db.HasConfig()
}

// TestConnection 测试AI连接 - 使用 Eino
func (s *ConfigService) TestConnection(apiKey, baseURL, modelName string) error {
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	if modelName == "" {
		modelName = "deepseek-chat"
	}

	if apiKey == "" {
		return fmt.Errorf("API Key不能为空")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 使用 Eino 创建 ChatModel
	cm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   modelName,
	})
	if err != nil {
		return fmt.Errorf("init eino chat model failed: %w", err)
	}

	// 发送测试消息
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: "你好，这是一个测试消息，请回复\"连接成功\"",
		},
	}

	_, err = cm.Generate(ctx, messages, einoModel.WithTemperature(0.7))
	if err != nil {
		return fmt.Errorf("test connection failed: %w", err)
	}

	return nil
}
