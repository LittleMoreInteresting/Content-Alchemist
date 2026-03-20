package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"content-alchemist/internal/model"
	"content-alchemist/internal/repository"
)

// ArticleService 文章服务
type ArticleService struct {
	db *repository.DB
}

// NewArticleService 创建文章服务
func NewArticleService() (*ArticleService, error) {
	db, err := repository.NewDB()
	if err != nil {
		return nil, fmt.Errorf("init db failed: %w", err)
	}

	return &ArticleService{db: db}, nil
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(title string) (*model.Article, error) {
	article := &model.Article{
		ID:        generateID(),
		Title:     title,
		Content:   fmt.Sprintf("# %s\n\n开始创作...", title),
		Outline:   []model.OutlineNode{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    "draft",
	}

	if err := s.db.SaveArticle(article); err != nil {
		return nil, err
	}

	return article, nil
}

// GetArticle 获取文章
func (s *ArticleService) GetArticle(id string) (*model.Article, error) {
	return s.db.GetArticle(id)
}

// SaveArticle 保存文章
func (s *ArticleService) SaveArticle(article *model.Article) error {
	article.UpdatedAt = time.Now()
	return s.db.SaveArticle(article)
}

// ListArticles 列出所有文章
func (s *ArticleService) ListArticles() ([]model.Article, error) {
	return s.db.ListArticles()
}

// DeleteArticle 删除文章
func (s *ArticleService) DeleteArticle(id string) error {
	return s.db.DeleteArticle(id)
}

// generateID 生成唯一ID
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
