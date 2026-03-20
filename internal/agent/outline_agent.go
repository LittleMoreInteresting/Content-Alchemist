package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	einoModel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"content-alchemist/internal/model"
)

// OutlineRequest 大纲生成请求
type OutlineRequest struct {
	Title    string `json:"title"`
	Style    string `json:"style"`
	Audience string `json:"audience"`
}

// Outline 大纲结构
type Outline struct {
	Nodes []OutlineNodeData `json:"nodes"`
}

// OutlineNodeData 大纲节点数据
type OutlineNodeData struct {
	ID             string `json:"id"`
	Level          int    `json:"level"`
	Title          string `json:"title"`
	SuggestedWords int    `json:"suggestedWords"`
}

// OutlineAgent 大纲生成Agent - 基于 Eino
type OutlineAgent struct {
	llm einoModel.ToolCallingChatModel
}

// NewOutlineAgent 创建大纲Agent - 使用 Eino
func NewOutlineAgent(apiKey, baseURL, modelName string) (*OutlineAgent, error) {
	baseAgent, err := NewBaseAgent(apiKey, baseURL, modelName)
	if err != nil {
		return nil, err
	}
	return &OutlineAgent{llm: baseAgent.llm}, nil
}

// GenerateOutline 生成大纲 - 使用 Eino Generate
func (a *OutlineAgent) GenerateOutline(ctx context.Context, req OutlineRequest) ([]model.OutlineNode, error) {
	prompt := fmt.Sprintf(`请为公众号文章"%s"生成详细大纲。

要求：
1. 受众：%s
2. 风格：%s
3. 大纲结构：包含引言、3-5个核心论点（每个论点配案例）、总结
4. 每个节点注明建议字数
5. 层级深度不超过3级

请严格按照以下JSON格式输出，不要包含其他内容：
{
  "nodes": [
    {"id": "1", "level": 1, "title": "引言", "suggestedWords": 200},
    {"id": "2", "level": 1, "title": "核心观点一", "suggestedWords": 500},
    {"id": "2.1", "level": 2, "title": "案例说明", "suggestedWords": 200},
    {"id": "3", "level": 1, "title": "总结", "suggestedWords": 150}
  ]
}`, req.Title, req.Audience, req.Style)

	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	// 使用 Eino Generate 接口
	resp, err := a.llm.Generate(ctx, messages, einoModel.WithTemperature(0.7))
	if err != nil {
		return nil, fmt.Errorf("generate outline failed: %w", err)
	}

	// 解析JSON响应
	content := resp.Content
	// 提取JSON部分
	startIdx := strings.Index(content, "{")
	endIdx := strings.LastIndex(content, "}")
	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		content = content[startIdx : endIdx+1]
	}

	var outline Outline
	if err := json.Unmarshal([]byte(content), &outline); err != nil {
		// 如果解析失败，返回一个默认大纲
		return a.generateDefaultOutline(req.Title), nil
	}

	// 转换为model.OutlineNode
	result := make([]model.OutlineNode, len(outline.Nodes))
	for i, node := range outline.Nodes {
		parentID := ""
		if strings.Contains(node.ID, ".") {
			parts := strings.Split(node.ID, ".")
			if len(parts) > 1 {
				parentID = parts[0]
			}
		}

		result[i] = model.OutlineNode{
			ID:          node.ID,
			Level:       node.Level,
			Title:       node.Title,
			ParentID:    parentID,
			Status:      "empty",
			TargetWords: node.SuggestedWords,
		}
	}

	return result, nil
}

// GenerateSectionContent 生成章节内容 - 使用 Eino
func (a *OutlineAgent) GenerateSectionContent(ctx context.Context, title string, section model.OutlineNode, context string, style string) (string, error) {
	prompt := fmt.Sprintf(`请为公众号文章撰写以下章节内容。

文章标题：%s
章节标题：%s
建议字数：%d字

上下文：
%s

写作风格：%s

要求：
1. 紧扣章节主题，内容充实
2. 与上下文自然衔接
3. 语言流畅，适合公众号阅读
4. 适当使用小标题、列表等排版元素

请直接输出正文内容，不需要标题。`, title, section.Title, section.TargetWords, context, style)

	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	// 使用 Eino Generate 接口
	resp, err := a.llm.Generate(ctx, messages, einoModel.WithTemperature(0.7))
	if err != nil {
		return "", fmt.Errorf("generate section content failed: %w", err)
	}

	return resp.Content, nil
}

// StreamGenerateSection 流式生成章节内容 - 使用 Eino Stream
func (a *OutlineAgent) StreamGenerateSection(ctx context.Context, title string, section model.OutlineNode, context string, style string, handler func(chunk string)) error {
	prompt := fmt.Sprintf(`请为公众号文章撰写以下章节内容。

文章标题：%s
章节标题：%s
建议字数：%d字

上下文：
%s

写作风格：%s

要求：
1. 紧扣章节主题，内容充实
2. 与上下文自然衔接
3. 语言流畅，适合公众号阅读
4. 适当使用小标题、列表等排版元素

请直接输出正文内容，不需要标题。`, title, section.Title, section.TargetWords, context, style)

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

// generateDefaultOutline 生成默认大纲
func (a *OutlineAgent) generateDefaultOutline(title string) []model.OutlineNode {
	return []model.OutlineNode{
		{ID: "1", Level: 1, Title: "引言", Status: "empty", TargetWords: 200},
		{ID: "2", Level: 1, Title: "核心观点", Status: "empty", TargetWords: 600},
		{ID: "2.1", Level: 2, Title: "案例说明", ParentID: "2", Status: "empty", TargetWords: 300},
		{ID: "3", Level: 1, Title: "实践建议", Status: "empty", TargetWords: 400},
		{ID: "4", Level: 1, Title: "总结", Status: "empty", TargetWords: 150},
	}
}
