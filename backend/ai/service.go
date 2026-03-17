//go:build darwin || windows || linux

package ai

import (
	"Content-Alchemist/backend/models"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// StreamChunk 流式响应块
type StreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// OutlineResult 生成大纲的结果
type OutlineResult struct {
	Titles  []string `json:"titles"`  // 3个候选爆款标题
	Outline string   `json:"outline"` // 文章大纲
}

// StreamCallback 流式输出回调函数
type StreamCallback func(chunk string) error

// GenerateOutlineWithTitles 根据标题、写作要求和公众号定位生成大纲和候选标题
func (s *Service) GenerateOutlineWithTitles(title, requirements, positioning string) (*OutlineResult, error) {
	if s.config == nil || s.config.Token == "" {
		return nil, fmt.Errorf("AI 配置不完整，请先配置 API Token")
	}

	// 构建完整的提示信息
	var contextParts []string
	if title != "" {
		contextParts = append(contextParts, fmt.Sprintf("原始标题: %s", title))
	}
	if requirements != "" {
		contextParts = append(contextParts, fmt.Sprintf("写作要求:\n%s", requirements))
	}
	if positioning != "" {
		contextParts = append(contextParts, fmt.Sprintf("公众号定位:\n%s", positioning))
	}

	context := strings.Join(contextParts, "\n\n")

	prompt := fmt.Sprintf(`请根据以下信息，完成两个任务：

%s

任务1 - 生成3个候选爆款标题：
基于原始标题和公众号定位，生成3个更具吸引力的爆款标题。
标题要求：
- 具有吸引力和点击率
- 符合公众号定位
- 可以包含数字、疑问、悬念等元素
- 每个标题一行，前面标注序号

任务2 - 生成文章大纲：
基于原始标题，生成详细的文章大纲。
大纲要求：
1. 使用 Markdown 格式
2. 只生成大纲结构（一级、二级、三级标题）
3. 不要生成正文内容
4. 大纲应该逻辑清晰、层次分明
5. 根据主题合理规划章节
6. 返回纯 Markdown 文本，不需要代码块标记

请按以下格式返回：

===TITLES===
1. 标题1
2. 标题2
3. 标题3
===END TITLES===

===OUTLINE===
（大纲内容）
===END OUTLINE===`, context)

	reqBody := ChatRequest{
		Model: s.config.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个专业的公众号写作专家，擅长创作爆款标题和清晰的文章大纲。",
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
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", s.config.BaseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.Token)

	client := &http.Client{
		Timeout: 120 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回错误 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if chatResp.Error != nil {
		return nil, fmt.Errorf("API 错误: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("API 返回空结果")
	}

	// 解析响应内容
	result := parseOutlineResponse(chatResp.Choices[0].Message.Content)
	return result, nil
}

// parseOutlineResponse 解析AI返回的内容，提取标题和大纲
func parseOutlineResponse(content string) *OutlineResult {
	result := &OutlineResult{
		Titles:  []string{},
		Outline: "",
	}

	// 提取标题部分
	titlesStart := strings.Index(content, "===TITLES===")
	titlesEnd := strings.Index(content, "===END TITLES===")
	if titlesStart != -1 && titlesEnd != -1 && titlesEnd > titlesStart {
		titlesContent := content[titlesStart+len("===TITLES===") : titlesEnd]
		lines := strings.Split(titlesContent, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			// 移除序号前缀（如 "1. "、"1."、"1 " 等）
			if line == "" {
				continue
			}
			// 尝试移除序号前缀
			if idx := strings.Index(line, "."); idx != -1 && idx < 4 {
				line = strings.TrimSpace(line[idx+1:])
			} else if idx := strings.Index(line, " "); idx != -1 && idx < 2 {
				line = strings.TrimSpace(line[idx+1:])
			}
			if line != "" {
				result.Titles = append(result.Titles, line)
			}
		}
	}

	// 提取大纲部分
	outlineStart := strings.Index(content, "===OUTLINE===")
	outlineEnd := strings.Index(content, "===END OUTLINE===")
	if outlineStart != -1 && outlineEnd != -1 && outlineEnd > outlineStart {
		result.Outline = strings.TrimSpace(content[outlineStart+len("===OUTLINE===") : outlineEnd])
	} else if outlineStart != -1 {
		// 如果没有结束标记，取到末尾
		result.Outline = strings.TrimSpace(content[outlineStart+len("===OUTLINE==="):])
	} else {
		// 如果没有标记，整个内容作为大纲
		result.Outline = strings.TrimSpace(content)
	}

	return result
}

// GenerateArticleFromOutline 根据大纲生成完整文章
func (s *Service) GenerateArticleFromOutline(title string, outline string, requirements string) (string, error) {
	if s.config == nil || s.config.Token == "" {
		return "", fmt.Errorf("AI 配置不完整，请先配置 API Token")
	}

	prompt := fmt.Sprintf(`请根据以下信息生成一篇完整的公众号文章:

标题: %s

大纲:
%s

写作要求:
%s

请严格按照大纲结构生成文章内容，要求:
1. 文章风格适合公众号发布，通俗易懂，有吸引力
2. 每个大纲节点都要有对应的内容展开
3. 文章要有引人入胜的开头
4. 正文内容详实，逻辑清晰
5. 结尾要有总结或行动号召
6. 使用 Markdown 格式
7. 返回纯 Markdown 文本，不需要代码块标记

请直接返回完整的文章内容:`, title, outline, requirements)

	reqBody := ChatRequest{
		Model: s.config.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个专业的公众号文章写作专家，擅长根据大纲生成高质量、易读的文章内容。",
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
		Timeout: 120 * time.Second,
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

// GenerateArticleStream 流式生成文章
func (s *Service) GenerateArticleStream(title string, outline string, requirements string, callback StreamCallback) error {
	if s.config == nil || s.config.Token == "" {
		return fmt.Errorf("AI 配置不完整，请先配置 API Token")
	}

	prompt := fmt.Sprintf(`请根据以下信息生成一篇完整的公众号文章:

标题: %s

大纲:
%s

写作要求:
%s

请严格按照大纲结构生成文章内容，要求:
1. 文章风格适合公众号发布，通俗易懂，有吸引力
2. 每个大纲节点都要有对应的内容展开
3. 文章要有引人入胜的开头
4. 正文内容详实，逻辑清晰
5. 结尾要有总结或行动号召
6. 使用 Markdown 格式
7. 返回纯 Markdown 文本，不需要代码块标记

请直接返回完整的文章内容:`, title, outline, requirements)

	reqBody := ChatRequest{
		Model: s.config.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个专业的公众号文章写作专家，擅长根据大纲生成高质量、易读的文章内容。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: s.config.Temperature,
		Stream:      true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", s.config.BaseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.Token)
	req.Header.Set("Accept", "text/event-stream")

	client := &http.Client{
		Timeout: 300 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API 返回错误 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	// 读取 SSE 流
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("读取流失败: %w", err)
		}

		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if chunk.Error != nil {
			return fmt.Errorf("流式响应错误: %s", chunk.Error.Message)
		}

		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			if content != "" {
				if err := callback(content); err != nil {
					return err
				}
			}
			if chunk.Choices[0].FinishReason == "stop" {
				break
			}
		}
	}

	return nil
}

// OptimizeContent 优化文章内容
func (s *Service) OptimizeContent(content string, optimizeType string, requirements string) (string, error) {
	if s.config == nil || s.config.Token == "" {
		return "", fmt.Errorf("AI 配置不完整，请先配置 API Token")
	}

	var prompt string
	switch optimizeType {
	case "polish":
		prompt = fmt.Sprintf(`请对以下内容进行润色优化，使其更加流畅、专业：

%s

要求：
%s

请直接返回优化后的内容：`, content, requirements)
	case "expand":
		prompt = fmt.Sprintf(`请对以下内容进行扩写，增加更多细节和深度：

%s

要求：
%s

请直接返回扩写后的内容：`, content, requirements)
	case "simplify":
		prompt = fmt.Sprintf(`请对以下内容进行精简，保留核心要点：

%s

要求：
%s

请直接返回精简后的内容：`, content, requirements)
	case "example":
		prompt = fmt.Sprintf(`请为以下内容添加相关案例，增强说服力：

%s

要求：
%s

请直接返回添加案例后的内容：`, content, requirements)
	default:
		prompt = fmt.Sprintf(`请根据以下要求优化内容：

内容：
%s

优化要求：
%s

请直接返回优化后的内容：`, content, requirements)
	}

	reqBody := ChatRequest{
		Model: s.config.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个专业的文章编辑，擅长优化和改进文章内容。",
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
		Timeout: 120 * time.Second,
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

// GenerateViralTitles 生成爆款标题
func (s *Service) GenerateViralTitles(content string, positioning string, count int) ([]string, error) {
	if s.config == nil || s.config.Token == "" {
		return nil, fmt.Errorf("AI 配置不完整，请先配置 API Token")
	}

	if count <= 0 {
		count = 5
	}
	if count > 10 {
		count = 10
	}

	prompt := fmt.Sprintf(`请根据以下文章内容，生成%d个爆款标题：

文章内容：
%s

公众号定位：
%s

爆款标题要求：
1. 具有强烈的吸引力和点击率
2. 可以使用数字、疑问、悬念等元素
3. 符合公众号定位和目标受众
4. 简洁有力，突出重点
5. 每个标题单独一行，前面标注序号

请按以下格式返回：

1. 标题1
2. 标题2
3. 标题3
...`, count, content, positioning)

	reqBody := ChatRequest{
		Model: s.config.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "你是一个专业的标题创作专家，擅长创作高点击率的爆款标题。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: s.config.Temperature + 0.2, // 稍微提高温度以增加创意
		Stream:      false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", s.config.BaseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.Token)

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 返回错误 (状态码 %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if chatResp.Error != nil {
		return nil, fmt.Errorf("API 错误: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("API 返回空结果")
	}

	// 解析标题
	return parseTitles(chatResp.Choices[0].Message.Content), nil
}

// parseTitles 解析标题列表
func parseTitles(content string) []string {
	var titles []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 移除序号前缀
		if idx := strings.Index(line, "."); idx != -1 && idx < 4 {
			line = strings.TrimSpace(line[idx+1:])
		} else if idx := strings.Index(line, " "); idx != -1 && idx < 2 {
			line = strings.TrimSpace(line[idx+1:])
		}
		if line != "" {
			titles = append(titles, line)
		}
	}
	return titles
}
