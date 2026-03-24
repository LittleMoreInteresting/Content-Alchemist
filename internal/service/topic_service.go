package service

import (
	"context"
	"fmt"
	"time"

	"content-alchemist/internal/agent"
	"content-alchemist/internal/model"
	"content-alchemist/internal/platform/hot"
	"content-alchemist/internal/repository"
)

// TopicService 选题服务
type TopicService struct {
	db          *repository.DB
	topicAgent  *agent.TopicAgent
	hotManager  *hot.Manager
}

// NewTopicService 创建选题服务
func NewTopicService(db *repository.DB, topicAgent *agent.TopicAgent) *TopicService {
	svc := &TopicService{
		db:         db,
		topicAgent: topicAgent,
		hotManager: hot.NewManager(),
	}
	
	// 注册热点获取器
	svc.hotManager.Register(hot.NewWeiboFetcher())
	svc.hotManager.Register(hot.NewZhihuFetcher())
	svc.hotManager.Register(hot.NewBaiduFetcher())
	svc.hotManager.Register(hot.NewToutiaoFetcher())
	
	return svc
}

// CreateTopic 创建选题
func (s *TopicService) CreateTopic(topic *model.Topic) error {
	if topic.ID == "" {
		topic.ID = generateTopicID()
	}
	if topic.Status == "" {
		topic.Status = model.TopicStatusPending
	}
	return s.db.SaveTopic(topic)
}

// GetTopic 获取选题
func (s *TopicService) GetTopic(id string) (*model.Topic, error) {
	return s.db.GetTopic(id)
}

// ListTopics 列出选题
func (s *TopicService) ListTopics(status string, limit int) ([]model.Topic, error) {
	return s.db.ListTopics(status, limit)
}

// ListTopicsBySource 根据来源列出选题
func (s *TopicService) ListTopicsBySource(source string, limit int) ([]model.Topic, error) {
	return s.db.ListTopicsBySource(source, limit)
}

// UpdateTopic 更新选题
func (s *TopicService) UpdateTopic(topic *model.Topic) error {
	existing, err := s.db.GetTopic(topic.ID)
	if err != nil {
		return err
	}
	topic.CreatedAt = existing.CreatedAt
	return s.db.SaveTopic(topic)
}

// DeleteTopic 删除选题
func (s *TopicService) DeleteTopic(id string) error {
	return s.db.DeleteTopic(id)
}

// ApproveTopic 批准选题
func (s *TopicService) ApproveTopic(id string) error {
	topic, err := s.db.GetTopic(id)
	if err != nil {
		return err
	}
	topic.Status = model.TopicStatusApproved
	return s.db.SaveTopic(topic)
}

// RejectTopic 拒绝选题
func (s *TopicService) RejectTopic(id string) error {
	topic, err := s.db.GetTopic(id)
	if err != nil {
		return err
	}
	topic.Status = model.TopicStatusRejected
	return s.db.SaveTopic(topic)
}

// SearchTopics 搜索选题
func (s *TopicService) SearchTopics(keyword string) ([]model.Topic, error) {
	return s.db.SearchTopics(keyword)
}

// GenerateTopics 生成选题
func (s *TopicService) GenerateTopics(ctx context.Context, trends []model.HotTrend, limit int) ([]model.TopicResult, error) {
	if s.topicAgent == nil {
		return nil, fmt.Errorf("topic agent not initialized")
	}
	
	req := agent.TopicRequest{
		Limit: limit,
	}
	
	results, err := s.topicAgent.GenerateTopics(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return results, nil
}

// FetchHotTrends 获取热点趋势（使用Manager）
func (s *TopicService) FetchHotTrends(ctx context.Context, platforms []string, limit int) ([]model.HotTrend, error) {
	results := s.hotManager.FetchPlatforms(platforms)
	
	var allTrends []model.HotTrend
	now := time.Now()
	
	for _, result := range results {
		if result.Error != "" {
			continue
		}
		
		for _, item := range result.Items {
			if len(allTrends) >= limit {
				break
			}
			
			trend := model.HotTrend{
				ID:        fmt.Sprintf("%s_%d_%d", result.Platform, item.Rank, now.Unix()),
				Platform:  result.Platform,
				Title:     item.Title,
				URL:       item.URL,
				HotRank:   item.Rank,
				HotValue:  item.HotValue,
				Category:  item.Category,
				CreatedAt: now,
			}
			
			allTrends = append(allTrends, trend)
		}
		
		if len(allTrends) >= limit {
			break
		}
	}
	
	return allTrends, nil
}

// FetchAllHotTrends 获取所有平台热点
func (s *TopicService) FetchAllHotTrends(ctx context.Context, limitPerPlatform int) ([]model.HotTrend, error) {
	platforms := s.hotManager.GetPlatforms()
	return s.FetchHotTrends(ctx, platforms, limitPerPlatform*len(platforms))
}

// GetHotPlatforms 获取热点平台列表
func (s *TopicService) GetHotPlatforms() []string {
	return s.hotManager.GetPlatforms()
}

// SaveHotTrends 保存热点趋势
func (s *TopicService) SaveHotTrends(trends []model.HotTrend) error {
	return s.db.BatchSaveHotTrends(trends)
}

// GetHotTrends 获取热点趋势
func (s *TopicService) GetHotTrends(platform string, limit int) ([]model.HotTrend, error) {
	return s.db.GetHotTrends(platform, limit)
}

// GetHotTrendPlatforms 获取热点平台列表（从数据库）
func (s *TopicService) GetHotTrendPlatforms() ([]string, error) {
	return s.db.GetHotTrendPlatforms()
}

// ClearOldHotTrends 清理旧的热点数据
func (s *TopicService) ClearOldHotTrends() error {
	// 清理7天前的数据
	before := time.Now().AddDate(0, 0, -7)
	return s.db.DeleteOldHotTrends(before)
}

// AIGenerateTopicsFromHot 基于热点AI生成选题
func (s *TopicService) AIGenerateTopicsFromHot(ctx context.Context, platforms []string, topicLimit int) ([]*model.Topic, error) {
	if s.topicAgent == nil {
		return nil, fmt.Errorf("topic agent not initialized")
	}
	
	// 1. 获取热点
	hotTrends, err := s.FetchHotTrends(ctx, platforms, 20)
	if err != nil {
		return nil, err
	}
	
	// 2. 保存热点
	if err := s.SaveHotTrends(hotTrends); err != nil {
		// 非致命错误，继续
	}
	
	// 3. AI生成选题
	results, err := s.topicAgent.GenerateTopics(ctx, agent.TopicRequest{
		Limit: topicLimit,
	})
	if err != nil {
		return nil, err
	}
	
	// 4. 转换为Topic模型并保存
	var topics []*model.Topic
	for _, result := range results {
		topic := &model.Topic{
			ID:       generateTopicID(),
			Title:    result.Title,
			Source:   model.TopicSourceAI,
			Score:    result.Score,
			Reason:   result.Reason,
			Angles:   result.Angles,
			Keywords: result.Keywords,
			Status:   model.TopicStatusPending,
		}
		
		if err := s.CreateTopic(topic); err != nil {
			continue
		}
		
		topics = append(topics, topic)
	}
	
	return topics, nil
}

// generateTopicID 生成选题ID
func generateTopicID() string {
	return fmt.Sprintf("topic_%d", time.Now().UnixNano())
}
