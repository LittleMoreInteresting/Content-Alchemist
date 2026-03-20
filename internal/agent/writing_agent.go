package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	einoModel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// WritingRequest 写作请求
type WritingRequest struct {
	Action       string  `json:"action"`       // "expand", "polish", "shorten", "continue", "title", "custom", "casual", "data", "rewrite"
	Context      string  `json:"context"`      // 上下文（前几段）
	SelectedText string  `json:"selectedText"` // 选中的文本
	Position     string  `json:"position"`     // "before", "after", "replace"
	Style        string  `json:"style"`        // 风格描述
	CustomPrompt string  `json:"customPrompt"` // 自定义提示词（当 action 为 "custom" 时使用）
}

// WritingResponse 写作响应
type WritingResponse struct {
	Content     string   `json:"content"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// WritingAgent 写作Agent - 基于 Eino
type WritingAgent struct {
	llm einoModel.ToolCallingChatModel
}

// NewWritingAgent 创建写作Agent - 使用 Eino
func NewWritingAgent(apiKey, baseURL, modelName string) (*WritingAgent, error) {
	baseAgent, err := NewBaseAgent(apiKey, baseURL, modelName)
	if err != nil {
		return nil, err
	}
	return &WritingAgent{llm: baseAgent.llm}, nil
}

// Execute 执行写作任务 - 使用 Eino Generate
func (a *WritingAgent) Execute(ctx context.Context, req WritingRequest) (*WritingResponse, error) {
	prompt := a.buildPrompt(req)
	
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	// 使用 Eino Generate 接口
	resp, err := a.llm.Generate(ctx, messages, einoModel.WithTemperature(0.7))
	if err != nil {
		return nil, fmt.Errorf("writing agent execute failed: %w", err)
	}

	// 如果是生成标题，尝试解析JSON数组
	if req.Action == "title" {
		var suggestions []string
		if err := json.Unmarshal([]byte(resp.Content), &suggestions); err == nil {
			return &WritingResponse{
				Content:     strings.Join(suggestions, "\n"),
				Suggestions: suggestions,
			}, nil
		}
	}

	return &WritingResponse{
		Content: resp.Content,
	}, nil
}

// StreamExecute 流式执行写作任务 - 使用 Eino Stream
func (a *WritingAgent) StreamExecute(ctx context.Context, req WritingRequest, handler func(chunk string)) error {
	prompt := a.buildPrompt(req)
	
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	// 使用 BaseAgent 的 GenerateStream 方法
	baseAgent := &BaseAgent{llm: a.llm}
	return baseAgent.GenerateStream(ctx, messages, handler)
}

// buildPrompt 构建提示词
func (a *WritingAgent) buildPrompt(req WritingRequest) string {
	var sb strings.Builder

	switch req.Action {
	case "expand":
		sb.WriteString("请对以下文本进行扩写，使其更加详细和有深度：\n\n")
		sb.WriteString(req.SelectedText)
		sb.WriteString("\n\n要求：")
		sb.WriteString("\n1. 保持原文的核心意思")
		sb.WriteString("\n2. 增加更多细节、例子或数据支撑")
		sb.WriteString("\n3. 保持段落连贯性")

	case "polish":
		sb.WriteString("请润色以下文本，使其更加专业和流畅：\n\n")
		sb.WriteString(req.SelectedText)
		sb.WriteString("\n\n要求：")
		sb.WriteString("\n1. 优化用词，提升专业度")
		sb.WriteString("\n2. 改善句子结构，增强可读性")
		sb.WriteString("\n3. 保持原意不变")

	case "shorten":
		sb.WriteString("请精简以下文本，保持核心信息的同时更加简洁：\n\n")
		sb.WriteString(req.SelectedText)
		sb.WriteString("\n\n要求：")
		sb.WriteString("\n1. 删除冗余词汇")
		sb.WriteString("\n2. 合并相似观点")
		sb.WriteString("\n3. 保留核心论点")

	case "continue":
		sb.WriteString("请基于以下上下文续写内容：\n\n")
		sb.WriteString("上下文：\n")
		sb.WriteString(req.Context)
		sb.WriteString("\n\n当前段落：\n")
		sb.WriteString(req.SelectedText)
		sb.WriteString("\n\n要求：")
		sb.WriteString("\n1. 自然衔接上文")
		sb.WriteString("\n2. 保持风格一致")
		sb.WriteString("\n3. 推进论述发展")

	case "title":
		sb.WriteString("请为以下文章生成5个吸引人的标题：\n\n")
		sb.WriteString(req.Context)
		sb.WriteString("\n\n要求：")
		sb.WriteString("\n1. 突出核心观点")
		sb.WriteString("\n2. 吸引读者点击")
		sb.WriteString("\n3. 长度适中（10-30字）")
		sb.WriteString("\n4. 以JSON数组格式返回")

	case "rewrite":
		sb.WriteString("请重写以下文本，保持核心意思但表达方式完全不同：\n\n")
		sb.WriteString(req.SelectedText)
		sb.WriteString("\n\n要求：")
		sb.WriteString("\n1. 使用不同的词汇和句式")
		sb.WriteString("\n2. 保持原意不变")
		sb.WriteString("\n3. 提升表达质量")

	case "casual":
		sb.WriteString("请将以下文本改写成更加口语化、通俗易懂的风格：\n\n")
		sb.WriteString(req.SelectedText)
		sb.WriteString("\n\n要求：")
		sb.WriteString("\n1. 使用日常用语")
		sb.WriteString("\n2. 增加亲和力")
		sb.WriteString("\n3. 保持内容准确")

	case "data":
		sb.WriteString("请为以下文本添加相关的数据支撑和统计信息：\n\n")
		sb.WriteString(req.SelectedText)
		sb.WriteString("\n\n要求：")
		sb.WriteString("\n1. 添加真实可信的数据")
		sb.WriteString("\n2. 使用具体数字增强说服力")
		sb.WriteString("\n3. 保持论述的逻辑性")

	case "custom":
		if req.CustomPrompt != "" {
			sb.WriteString(req.CustomPrompt)
			sb.WriteString("\n\n文本内容：\n")
			sb.WriteString(req.SelectedText)
		} else {
			sb.WriteString(req.SelectedText)
		}

	default:
		sb.WriteString(req.SelectedText)
	}

	if req.Style != "" {
		sb.WriteString(fmt.Sprintf("\n\n写作风格：%s", req.Style))
	}

	return sb.String()
}
