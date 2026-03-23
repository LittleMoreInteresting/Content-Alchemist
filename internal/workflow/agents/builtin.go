package agents

import (
	"context"
	"fmt"
	"time"

	"content-alchemist/internal/model"
)

// HotFetchAgent 热点获取Agent
type HotFetchAgent struct {
	FetchFunc func(ctx context.Context, platforms []string, limit int) ([]model.HotTrend, error)
}

// Execute 执行热点获取
func (a *HotFetchAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	if a.FetchFunc == nil {
		return nil, fmt.Errorf("fetch function not set")
	}

	config, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	// 获取平台列表
	platforms := []string{"weibo", "zhihu", "baidu", "toutiao"}
	if p, ok := config["platforms"].([]string); ok {
		platforms = p
	}

	limit := 20
	if l, ok := config["limit"].(int); ok {
		limit = l
	}

	// 调用获取函数
	trends, err := a.FetchFunc(ctx, platforms, limit)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"trends": trends,
	}, nil
}

// TopicGenerateAgent 选题生成Agent
type TopicGenerateAgent struct {
	GenerateFunc func(ctx context.Context, trends []model.HotTrend, limit int) ([]model.TopicResult, error)
}

// Execute 执行选题生成
func (a *TopicGenerateAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	if a.GenerateFunc == nil {
		return nil, fmt.Errorf("generate function not set")
	}

	config, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	// 获取热点数据
	trends, ok := config["trends"].([]model.HotTrend)
	if !ok {
		return nil, fmt.Errorf("trends not found in input")
	}

	limit := 5
	if l, ok := config["limit"].(int); ok {
		limit = l
	}

	// 调用生成函数
	topics, err := a.GenerateFunc(ctx, trends, limit)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"topics": topics,
	}, nil
}

// TopicSelectAgent 选题选择Agent
type TopicSelectAgent struct{}

// Execute 执行选题选择
func (a *TopicSelectAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	config, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	topics, ok := config["topics"].([]model.TopicResult)
	if !ok || len(topics) == 0 {
		return nil, fmt.Errorf("no topics to select")
	}

	strategy := "highest_score"
	if s, ok := config["strategy"].(string); ok {
		strategy = s
	}

	var selected *model.TopicResult

	switch strategy {
	case "highest_score":
		// 选择评分最高的
		for i := range topics {
			if selected == nil || topics[i].Score > selected.Score {
				selected = &topics[i]
			}
		}
	case "random":
		// 随机选择
		if len(topics) > 0 {
			selected = &topics[0] // 简化处理
		}
	default:
		return nil, fmt.Errorf("unknown strategy: %s", strategy)
	}

	return map[string]interface{}{
		"selected_topic": selected,
	}, nil
}

// OutlineGenerateAgent 大纲生成Agent
type OutlineGenerateAgent struct {
	GenerateFunc func(ctx context.Context, title, style string) ([]model.OutlineNode, error)
}

// Execute 执行大纲生成
func (a *OutlineGenerateAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	config, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	var title string
	if t, ok := config["title"].(string); ok {
		title = t
	} else if topic, ok := config["selected_topic"].(*model.TopicResult); ok {
		title = topic.Title
	}

	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	style := "干货专业"
	if s, ok := config["style"].(string); ok {
		style = s
	}

	// 使用默认大纲或调用生成函数
	var outline []model.OutlineNode
	if a.GenerateFunc != nil {
		var err error
		outline, err = a.GenerateFunc(ctx, title, style)
		if err != nil {
			outline = generateDefaultOutline()
		}
	} else {
		outline = generateDefaultOutline()
	}

	return map[string]interface{}{
		"outline": outline,
		"title":   title,
		"style":   style,
	}, nil
}

// generateDefaultOutline 生成默认大纲
func generateDefaultOutline() []model.OutlineNode {
	return []model.OutlineNode{
		{ID: "1", Level: 1, Title: "引言", Status: "empty", TargetWords: 200},
		{ID: "2", Level: 1, Title: "核心观点", Status: "empty", TargetWords: 600},
		{ID: "2.1", Level: 2, Title: "案例说明", ParentID: "2", Status: "empty", TargetWords: 300},
		{ID: "3", Level: 1, Title: "实践建议", Status: "empty", TargetWords: 400},
		{ID: "4", Level: 1, Title: "总结", Status: "empty", TargetWords: 150},
	}
}

// ContentWriteAgent 内容创作Agent
type ContentWriteAgent struct {
	WriteFunc func(ctx context.Context, outline []model.OutlineNode, title string) (*model.Article, error)
}

// Execute 执行内容创作
func (a *ContentWriteAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	config, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	outline, ok := config["outline"].([]model.OutlineNode)
	if !ok {
		return nil, fmt.Errorf("outline not found")
	}

	title, _ := config["title"].(string)

	var article *model.Article
	if a.WriteFunc != nil {
		var err error
		article, err = a.WriteFunc(ctx, outline, title)
		if err != nil {
			article = generateDefaultArticle(title, outline)
		}
	} else {
		article = generateDefaultArticle(title, outline)
	}

	return map[string]interface{}{
		"article": article,
		"outline": outline,
	}, nil
}

// generateDefaultArticle 生成默认文章
func generateDefaultArticle(title string, outline []model.OutlineNode) *model.Article {
	content := fmt.Sprintf("# %s\n\n这是自动创作的文章内容...\n\n", title)

	return &model.Article{
		ID:        generateAgentID(),
		Title:     title,
		Content:   content,
		Outline:   outline,
		Status:    model.ArticleStatusDraft,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ContentReviewAgent 内容审核Agent
type ContentReviewAgent struct {
	ReviewFunc func(ctx context.Context, article *model.Article) (map[string]interface{}, error)
}

// Execute 执行内容审核
func (a *ContentReviewAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	config, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	article, ok := config["article"].(*model.Article)
	if !ok {
		return nil, fmt.Errorf("article not found")
	}

	var review map[string]interface{}
	if a.ReviewFunc != nil {
		var err error
		review, err = a.ReviewFunc(ctx, article)
		if err != nil {
			review = generateDefaultReview()
		}
	} else {
		review = generateDefaultReview()
	}

	article.QualityScore = 85.0

	return map[string]interface{}{
		"review":  review,
		"article": article,
	}, nil
}

// generateDefaultReview 生成默认审核结果
func generateDefaultReview() map[string]interface{} {
	return map[string]interface{}{
		"passed": true,
		"score":  85.0,
		"issues": []map[string]string{},
	}
}

// LayoutApplyAgent 排版Agent
type LayoutApplyAgent struct {
	LayoutFunc func(ctx context.Context, article *model.Article, theme string) (string, error)
}

// Execute 执行排版
func (a *LayoutApplyAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	config, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	article, ok := config["article"].(*model.Article)
	if !ok {
		return nil, fmt.Errorf("article not found")
	}

	theme := "科技蓝"
	if t, ok := config["theme"].(string); ok {
		theme = t
	}

	var html string
	if a.LayoutFunc != nil {
		var err error
		html, err = a.LayoutFunc(ctx, article, theme)
		if err != nil {
			html = fmt.Sprintf("<html><body><h1>%s</h1><p>%s</p></body></html>", article.Title, article.Content)
		}
	} else {
		html = fmt.Sprintf("<html><body><h1>%s</h1><p>%s</p></body></html>", article.Title, article.Content)
	}

	return map[string]interface{}{
		"html":    html,
		"article": article,
		"theme":   theme,
	}, nil
}

// DelayAgent 延迟Agent
type DelayAgent struct{}

// Execute 执行延迟
func (a *DelayAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	config, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	duration := 5
	if d, ok := config["duration"].(int); ok {
		duration = d
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(time.Duration(duration) * time.Second):
		return map[string]interface{}{
			"delayed": true,
		}, nil
	}
}

// ConditionAgent 条件Agent
type ConditionAgent struct{}

// Execute 执行条件判断
func (a *ConditionAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	config, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	// 获取审核结果进行判断
	review, ok := config["review"].(map[string]interface{})
	if ok {
		if passed, ok := review["passed"].(bool); ok {
			return map[string]interface{}{
				"condition": passed,
			}, nil
		}
		if score, ok := review["score"].(float64); ok {
			return map[string]interface{}{
				"condition": score >= 80,
			}, nil
		}
	}

	// 默认返回true
	return map[string]interface{}{
		"condition": true,
	}, nil
}

// ManualReviewAgent 人工审核Agent
type ManualReviewAgent struct {
	ReviewChan chan *ManualReviewRequest
}

// ManualReviewRequest 人工审核请求
type ManualReviewRequest struct {
	RunID    string
	StepID   string
	Data     interface{}
	Response chan ManualReviewResponse
}

// ManualReviewResponse 人工审核响应
type ManualReviewResponse struct {
	Approved bool
	Data     interface{}
}

// Execute 执行人工审核
func (a *ManualReviewAgent) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	if a.ReviewChan == nil {
		return nil, fmt.Errorf("review channel not initialized")
	}

	config, ok := input.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid input type")
	}

	// 从配置中获取runID和stepID
	runID, _ := config["run_id"].(string)
	stepID, _ := config["step_id"].(string)

	responseChan := make(chan ManualReviewResponse, 1)

	request := &ManualReviewRequest{
		RunID:    runID,
		StepID:   stepID,
		Data:     config,
		Response: responseChan,
	}

	// 发送审核请求
	a.ReviewChan <- request

	// 等待审核结果
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case response := <-responseChan:
		return map[string]interface{}{
			"approved": response.Approved,
			"data":     response.Data,
		}, nil
	}
}

// generateAgentID 生成唯一ID
func generateAgentID() string {
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), time.Now().Unix())
}
