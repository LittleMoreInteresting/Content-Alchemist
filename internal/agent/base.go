package agent

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	einoModel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

// ChatRequest 聊天请求
type ChatRequest struct {
	SystemPrompt string  `json:"systemPrompt"`
	UserPrompt   string  `json:"userPrompt"`
	Temperature  float32 `json:"temperature"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Content     string   `json:"content"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// BaseAgent 基于 Eino 的基础Agent封装
type BaseAgent struct {
	llm einoModel.ToolCallingChatModel
}

// NewBaseAgent 创建基于 Eino 的基础Agent
func NewBaseAgent(apiKey string, baseURL string, modelName string) (*BaseAgent, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}

	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	if modelName == "" {
		modelName = "deepseek-chat"
	}

	ctx := context.Background()
	
	// 创建自定义 HTTP 客户端，禁用 HTTP/2 以避免流错误
	httpClient := &http.Client{
		Timeout: 120 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     false, // 禁用 HTTP/2
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}
	
	// 使用 Eino 的 OpenAI 组件创建 ChatModel
	// 支持 DeepSeek（兼容 OpenAI API 格式）
	cm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:     apiKey,
		BaseURL:    baseURL,
		Model:      modelName,
		HTTPClient: httpClient,
	})
	if err != nil {
		return nil, fmt.Errorf("init eino chat model failed: %w", err)
	}

	return &BaseAgent{
		llm: cm,
	}, nil
}

// StreamChat 流式聊天 - 使用 Eino 的 Stream 接口
func (a *BaseAgent) StreamChat(ctx context.Context, req ChatRequest, handler func(chunk string)) error {
	messages := []*schema.Message{}
	
	if req.SystemPrompt != "" {
		messages = append(messages, &schema.Message{
			Role:    schema.System,
			Content: req.SystemPrompt,
		})
	}
	
	messages = append(messages, &schema.Message{
		Role:    schema.User,
		Content: req.UserPrompt,
	})

	// 准备选项
	opts := []einoModel.Option{}
	if req.Temperature > 0 {
		opts = append(opts, einoModel.WithTemperature(req.Temperature))
	}

	// 调用 Eino 的 Stream 接口
	streamReader, err := a.llm.Stream(ctx, messages, opts...)
	if err != nil {
		return fmt.Errorf("eino stream failed: %w", err)
	}
	defer streamReader.Close()

	// 读取流式输出
	for {
		chunk, err := streamReader.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("stream recv error: %w", err)
		}
		if chunk != nil && chunk.Content != "" {
			handler(chunk.Content)
		}
	}

	return nil
}

// Chat 普通聊天（非流式）- 使用 Eino 的 Generate 接口
func (a *BaseAgent) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	messages := []*schema.Message{}
	
	if req.SystemPrompt != "" {
		messages = append(messages, &schema.Message{
			Role:    schema.System,
			Content: req.SystemPrompt,
		})
	}
	
	messages = append(messages, &schema.Message{
		Role:    schema.User,
		Content: req.UserPrompt,
	})

	// 准备选项
	opts := []einoModel.Option{}
	if req.Temperature > 0 {
		opts = append(opts, einoModel.WithTemperature(req.Temperature))
	}

	// 调用 Eino 的 Generate 接口
	response, err := a.llm.Generate(ctx, messages, opts...)
	if err != nil {
		return nil, fmt.Errorf("eino generate failed: %w", err)
	}

	return &ChatResponse{
		Content: response.Content,
	}, nil
}

// GenerateStream 生成流式响应的辅助函数
func (a *BaseAgent) GenerateStream(ctx context.Context, messages []*schema.Message, handler func(chunk string)) error {
	streamReader, err := a.llm.Stream(ctx, messages)
	if err != nil {
		return err
	}
	defer streamReader.Close()

	for {
		chunk, err := streamReader.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if chunk != nil && chunk.Content != "" {
			handler(chunk.Content)
		}
	}
	return nil
}
