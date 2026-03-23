package service

import (
	"context"
	"fmt"
	"time"

	"content-alchemist/internal/agent"
	"content-alchemist/internal/model"
	"content-alchemist/internal/repository"
)

// TopicService 选题服务
type TopicService struct {
	db         *repository.DB
	topicAgent *agent.TopicAgent
}

// NewTopicService 创建选题服务
func NewTopicService(db *repository.DB, topicAgent *agent.TopicAgent) *TopicService {
	return &TopicService{
		db:         db,
		topicAgent: topicAgent,
	}
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

// FetchHotTrends 获取热点趋势
func (s *TopicService) FetchHotTrends(ctx context.Context, platforms []string, limit int) ([]model.HotTrend, error) {
	// TODO: 实现热点获取逻辑
	// 这里使用模拟数据
	var trends []model.HotTrend
	
	for _, platform := range platforms {
		// 模拟从各平台获取热点
		for i := 1; i <= limit/len(platforms); i++ {
			trends = append(trends, model.HotTrend{
				ID:        fmt.Sprintf("%s_%d_%d", platform, i, time.Now().Unix()),
				Platform:  platform,
				Title:     fmt.Sprintf("%s热点 %d", platform, i),
				URL:       fmt.Sprintf("https://%s.com/trend/%d", platform, i),
				HotRank:   i,
				HotValue:  float64(1000000 - i*10000),
				Category:  "科技",
				CreatedAt: time.Now(),
			})
		}
	}
	
	return trends, nil
}

// SaveHotTrends 保存热点趋势
func (s *TopicService) SaveHotTrends(trends []model.HotTrend) error {
	return s.db.BatchSaveHotTrends(trends)
}

// GetHotTrends 获取热点趋势
func (s *TopicService) GetHotTrends(platform string, limit int) ([]model.HotTrend, error) {
	return s.db.GetHotTrends(platform, limit)
}

// GetHotTrendPlatforms 获取热点平台列表
func (s *TopicService) GetHotTrendPlatforms() ([]string, error) {
	return s.db.GetHotTrendPlatforms()
}

// ClearOldHotTrends 清理旧的热点数据
func (s *TopicService) ClearOldHotTrends() error {
	// 清理7天前的数据
	before := time.Now().AddDate(0, 0, -7)
	return s.db.DeleteOldHotTrends(before)
}

// generateTopicID 生成选题ID
func generateTopicID() string {
	return fmt.Sprintf("topic_%d", time.Now().UnixNano())
}
