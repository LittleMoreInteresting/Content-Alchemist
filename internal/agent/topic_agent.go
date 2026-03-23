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

// TopicAgent 选题Agent
type TopicAgent struct {
	llm einoModel.ToolCallingChatModel
}

// TopicRequest 选题请求
type TopicRequest struct {
	Keywords      []string `json:"keywords"`      // 关注领域
	ExcludeTopics []string `json:"excludeTopics"` // 排除的选题
	Limit         int      `json:"limit"`         // 生成数量
}

// TopicResult 选题结果 (使用model中的定义)
type TopicResult = model.TopicResult

// NewTopicAgent 创建选题Agent
func NewTopicAgent(apiKey, baseURL, modelName string) (*TopicAgent, error) {
	baseAgent, err := NewBaseAgent(apiKey, baseURL, modelName)
	if err != nil {
		return nil, err
	}
	return &TopicAgent{llm: baseAgent.llm}, nil
}

// GenerateTopics 基于热点生成选题
func (a *TopicAgent) GenerateTopics(ctx context.Context, req TopicRequest) ([]TopicResult, error) {
	if req.Limit <= 0 {
		req.Limit = 5
	}
	
	prompt := fmt.Sprintf(`请作为一位拥有10万+阅读量爆款经验的公众号资深主编，基于当前热点数据生成%d个优质选题。

要求：
1. 选题要有传播潜力，能引发读者共鸣
2. 每个选题需包含：标题、评分(0-100)、推荐理由、3个切入角度、关键词
3. 标题要有吸引力，符合公众号爆款标题特征
4. 切入角度要多样化，覆盖不同读者群体
5. 评分标准：热度(40%) + 深度(30%) + 共鸣度(30%)

请严格按照以下JSON格式输出，不要包含其他内容：
{
  "topics": [
    {
      "title": "选题标题",
      "score": 85.5,
      "reason": "推荐理由",
      "angles": ["角度1", "角度2", "角度3"],
      "keywords": ["关键词1", "关键词2"]
    }
  ]
}`, req.Limit)

	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	// 使用 Eino Generate 接口
	resp, err := a.llm.Generate(ctx, messages, einoModel.WithTemperature(0.8))
	if err != nil {
		return nil, fmt.Errorf("generate topics failed: %w", err)
	}

	// 解析JSON响应
	return a.parseTopicResults(resp.Content)
}

// AnalyzeTopic 深度分析选题
func (a *TopicAgent) AnalyzeTopic(ctx context.Context, title string) (*model.TopicAnalysis, error) {
	prompt := fmt.Sprintf(`请对以下选题进行深度分析：

选题：%s

请从以下维度进行分析：
1. 热度评分(0-100)：当前话题的流行程度
2. 竞争度评分(0-100)：同类内容的数量和质量
3. 匹配度评分(0-100)：与公众号定位的匹配程度
4. 传播潜力(0-100)：成为爆款的概率
5. 最佳发布时间：建议的发布时间
6. 竞品分析：列举3-5篇同类热门文章
7. 关键词分析：相关关键词的搜索量和竞争度

请严格按照以下JSON格式输出：
{
  "hotScore": 80,
  "compScore": 60,
  "fitScore": 90,
  "potential": 85,
  "timeline": "建议本周内发布",
  "competitors": [
    {"title": "竞品标题", "platform": "公众号", "readCount": 100000}
  ],
  "keywords": [
    {"word": "关键词", "searchVol": 10000, "competition": 0.5}
  ]
}`, title)

	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	resp, err := a.llm.Generate(ctx, messages, einoModel.WithTemperature(0.7))
	if err != nil {
		return nil, fmt.Errorf("analyze topic failed: %w", err)
	}

	return a.parseTopicAnalysis(resp.Content, title)
}

// parseTopicResults 解析选题结果
func (a *TopicAgent) parseTopicResults(content string) ([]TopicResult, error) {
	// 提取JSON部分
	startIdx := strings.Index(content, "{")
	endIdx := strings.LastIndex(content, "}")
	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		content = content[startIdx : endIdx+1]
	}

	var result struct {
		Topics []TopicResult `json:"topics"`
	}

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		// 如果解析失败，返回默认结果
		return a.generateDefaultTopics(), nil
	}

	return result.Topics, nil
}

// parseTopicAnalysis 解析选题分析
func (a *TopicAgent) parseTopicAnalysis(content string, title string) (*model.TopicAnalysis, error) {
	// 提取JSON部分
	startIdx := strings.Index(content, "{")
	endIdx := strings.LastIndex(content, "}")
	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		content = content[startIdx : endIdx+1]
	}

	var analysis model.TopicAnalysis
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		// 如果解析失败，返回默认分析
		return &model.TopicAnalysis{
			TopicID:   "",
			Title:     title,
			HotScore:  70,
			CompScore: 50,
			FitScore:  80,
			Potential: 70,
			Timeline:  "建议尽快发布",
		}, nil
	}

	analysis.Title = title
	return &analysis, nil
}

// generateDefaultTopics 生成默认选题
func (a *TopicAgent) generateDefaultTopics() []TopicResult {
	return []TopicResult{
		{
			Title:    "AI时代，内容创作者如何保持竞争力",
			Score:    85.0,
			Reason:   "AI话题持续火热，与创作者群体高度相关",
			Angles:   []string{"工具使用", "能力升级", "差异化策略"},
			Keywords: []string{"AI", "创作", "竞争力"},
		},
		{
			Title:    "从0到10万粉：我的公众号运营实战总结",
			Score:    80.0,
			Reason:   "成长路径类内容一直受读者欢迎",
			Angles:   []string{"内容策略", "涨粉技巧", "变现方式"},
			Keywords: []string{"运营", "涨粉", "变现"},
		},
	}
}
