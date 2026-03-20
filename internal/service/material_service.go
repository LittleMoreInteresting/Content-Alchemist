package service

import (
	"fmt"
	"time"

	"content-alchemist/internal/model"
	"content-alchemist/internal/repository"
)

// MaterialService 素材服务
type MaterialService struct {
	db *repository.DB
}

// NewMaterialService 创建素材服务
func NewMaterialService() (*MaterialService, error) {
	db, err := repository.NewDB()
	if err != nil {
		return nil, fmt.Errorf("init db failed: %w", err)
	}

	return &MaterialService{db: db}, nil
}

// CreateMaterial 创建素材
func (s *MaterialService) CreateMaterial(material *model.Material) error {
	if material.ID == "" {
		material.ID = generateID()
	}
	material.CreatedAt = time.Now()
	material.UsageCount = 0

	return s.db.SaveMaterial(material)
}

// GetMaterial 获取素材（从列表中查找）
func (s *MaterialService) GetMaterial(id string) (*model.Material, error) {
	materials, err := s.db.GetMaterials("")
	if err != nil {
		return nil, err
	}
	for i := range materials {
		if materials[i].ID == id {
			return &materials[i], nil
		}
	}
	return nil, fmt.Errorf("material not found")
}

// UpdateMaterial 更新素材
func (s *MaterialService) UpdateMaterial(material *model.Material) error {
	return s.db.SaveMaterial(material)
}

// DeleteMaterial 删除素材
func (s *MaterialService) DeleteMaterial(id string) error {
	return s.db.DeleteMaterial(id)
}

// ListMaterials 列出售有素材
func (s *MaterialService) ListMaterials(materialType string) ([]model.Material, error) {
	return s.db.GetMaterials(materialType)
}

// SearchMaterials 搜索素材（简单实现）
func (s *MaterialService) SearchMaterials(keyword string) ([]model.Material, error) {
	materials, err := s.db.GetMaterials("")
	if err != nil {
		return nil, err
	}

	// 简单过滤
	var results []model.Material
	for _, m := range materials {
		if contains(m.Title, keyword) || contains(m.Content, keyword) {
			results = append(results, m)
		}
	}
	return results, nil
}

// IncrementUsage 增加使用次数
func (s *MaterialService) IncrementUsage(id string) error {
	material, err := s.GetMaterial(id)
	if err != nil {
		return err
	}
	material.UsageCount++
	return s.db.SaveMaterial(material)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(s[:len(substr)] == substr) || contains(s[1:], substr))
}
