//go:build darwin || windows || linux

package ai

import (
	"Content-Alchemist/backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Service AI 服务
type Service struct {
	config *models.AIConfig
}

// NewService 创建 AI 服务
func NewService(config *models.AIConfig) *Service {
	return &Service{
		config: config,
	}
}

// UpdateConfig 更新配置
func (s *Service) UpdateConfig(config *models.AIConfig) {
	s.config = config
}

// Message 聊天消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	Stream      bool      `json:"stream"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// GenerateOutlineByTitle 根据标题生成文章大纲
func (s *Service) GenerateOutlineByTitle(title string) (string, error) {
	if s.config == nil || s.config.Token == "" {
		return "", fmt.Errorf("AI 配置不完整，请先配置 API Token")
	}

	prompt := fmt.Sprintf(`请为以下文章标题生成一个详细的大纲:

标题: %s

要求:
1. 使用 Markdown 格式
2. 只生成大纲结构（一级、二级、三级标题）
3. 不要生成正文内容
4. 大纲应该逻辑清晰、层次分明
5. 根据标题主题合理规划章节
6. 返回纯 Markdown 文本，不需要代码块标记

请直接返回大纲内容:`, title)

	reqBody := ChatRequest{
		Model: s.config.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个专业的写作助手，擅长根据标题生成清晰、有条理的文章大纲。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: s.config.Temperature,
		Stream:      false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", s.config.BaseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.Token)

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 返回错误 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API 错误: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("API 返回空结果")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// GenerateOutlineRequest 生成大纲请求
type GenerateOutlineRequest struct {
	Title string `json:"title"`
}

// GenerateOutlineResponse 生成大纲响应
type GenerateOutlineResponse struct {
	Outline string `json:"outline"`
	Error   string `json:"error,omitempty"`
}
