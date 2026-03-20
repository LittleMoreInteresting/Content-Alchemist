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

// ErrKeyNotFound 加密密钥未找到错误
var ErrKeyNotFound = repository.ErrKeyNotFound

// ConfigService 配置服务
type ConfigService struct {
	db *repository.DB
}

// NewConfigService 创建配置服务
func NewConfigService() (*ConfigService, error) {
	db, err := repository.NewDB()
	if err != nil {
		return nil, err
	}

	return &ConfigService{db: db}, nil
}

// InitEncryptionKey 初始化加密密钥（首次使用）
func InitEncryptionKey() error {
	return repository.InitWithNewKey()
}

// HasEncryptionKey 检查是否存在加密密钥
func HasEncryptionKey() bool {
	return repository.HasKey()
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

// HasConfig 检查是否有配置
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

	// 记录日志（隐藏部分 API Key）
	maskedKey := apiKey
	if len(apiKey) > 8 {
		maskedKey = apiKey[:4] + "****" + apiKey[len(apiKey)-4:]
	}
	fmt.Printf("Testing connection with API Key: %s, BaseURL: %s, Model: %s\n", maskedKey, baseURL, modelName)

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
			Content: "你好",
		},
	}

	_, err = cm.Generate(ctx, messages, einoModel.WithTemperature(0.7))
	if err != nil {
		return fmt.Errorf("test connection failed: %w", err)
	}

	return nil
}
