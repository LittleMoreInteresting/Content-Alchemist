package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"content-alchemist/internal/model"
	"content-alchemist/internal/repository"
	"content-alchemist/internal/workflow"
	"gopkg.in/yaml.v3"
)

// WorkflowService 工作流服务
type WorkflowService struct {
	db     *repository.DB
	engine *workflow.Engine
}

// NewWorkflowService 创建工作流服务
func NewWorkflowService(db *repository.DB) *WorkflowService {
	engine := workflow.NewEngine(db)
	return &WorkflowService{
		db:     db,
		engine: engine,
	}
}

// GetEngine 获取工作流引擎
func (s *WorkflowService) GetEngine() *workflow.Engine {
	return s.engine
}

// CreateWorkflow 创建工作流
func (s *WorkflowService) CreateWorkflow(workflow *model.Workflow) error {
	if workflow.ID == "" {
		workflow.ID = generateWorkflowID()
	}
	
	// 默认状态
	if workflow.Status == "" {
		workflow.Status = model.WorkflowStatusActive
	}
	
	// 验证工作流
	if err := s.validateWorkflow(workflow); err != nil {
		return fmt.Errorf("validate workflow failed: %w", err)
	}
	
	return s.db.SaveWorkflow(workflow)
}

// UpdateWorkflow 更新工作流
func (s *WorkflowService) UpdateWorkflow(workflow *model.Workflow) error {
	// 检查工作流是否存在
	existing, err := s.db.GetWorkflow(workflow.ID)
	if err != nil {
		return fmt.Errorf("workflow not found: %w", err)
	}
	
	// 如果工作流正在运行，不能更新
	runs, err := s.db.ListActiveWorkflowRuns()
	if err != nil {
		return err
	}
	for _, run := range runs {
		if run.WorkflowID == workflow.ID {
			return fmt.Errorf("cannot update workflow with active runs")
		}
	}
	
	// 保留创建时间
	workflow.CreatedAt = existing.CreatedAt
	
	// 验证工作流
	if err := s.validateWorkflow(workflow); err != nil {
		return fmt.Errorf("validate workflow failed: %w", err)
	}
	
	return s.db.SaveWorkflow(workflow)
}

// GetWorkflow 获取工作流
func (s *WorkflowService) GetWorkflow(id string) (*model.Workflow, error) {
	return s.db.GetWorkflow(id)
}

// ListWorkflows 列出所有工作流
func (s *WorkflowService) ListWorkflows() ([]model.Workflow, error) {
	return s.db.ListWorkflows()
}

// DeleteWorkflow 删除工作流
func (s *WorkflowService) DeleteWorkflow(id string) error {
	// 检查是否有正在运行的实例
	runs, err := s.db.ListActiveWorkflowRuns()
	if err != nil {
		return err
	}
	for _, run := range runs {
		if run.WorkflowID == id {
			return fmt.Errorf("cannot delete workflow with active runs")
		}
	}
	
	return s.db.DeleteWorkflow(id)
}

// StartWorkflow 启动工作流
func (s *WorkflowService) StartWorkflow(ctx context.Context, workflowID string, input map[string]interface{}) (*model.WorkflowRun, error) {
	return s.engine.Start(ctx, workflowID, input)
}

// PauseWorkflow 暂停工作流
func (s *WorkflowService) PauseWorkflow(runID string) error {
	return s.engine.Pause(runID)
}

// ResumeWorkflow 恢复工作流
func (s *WorkflowService) ResumeWorkflow(ctx context.Context, runID string) error {
	return s.engine.Resume(ctx, runID)
}

// CancelWorkflow 取消工作流
func (s *WorkflowService) CancelWorkflow(runID string) error {
	return s.engine.Cancel(runID)
}

// GetWorkflowRun 获取工作流运行状态
func (s *WorkflowService) GetWorkflowRun(runID string) (*model.WorkflowRun, error) {
	return s.engine.GetRun(runID)
}

// ListWorkflowRuns 列出工作流运行实例
func (s *WorkflowService) ListWorkflowRuns(workflowID string) ([]model.WorkflowRun, error) {
	return s.engine.ListRuns(workflowID)
}

// ImportFromYAML 从YAML导入工作流
func (s *WorkflowService) ImportFromYAML(yamlContent string) (*model.Workflow, error) {
	var yamlWorkflow struct {
		Name        string                 `yaml:"name"`
		Description string                 `yaml:"description"`
		Trigger     map[string]interface{} `yaml:"trigger"`
		Steps       []struct {
			ID       string                 `yaml:"id"`
			Name     string                 `yaml:"name"`
			Type     string                 `yaml:"type"`
			Config   map[string]interface{} `yaml:"config"`
			NextStep string                 `yaml:"next"`
			OnError  string                 `yaml:"on_error"`
		} `yaml:"steps"`
		AutoPublish bool `yaml:"auto_publish"`
		NeedReview  bool `yaml:"need_review"`
	}
	
	if err := yaml.Unmarshal([]byte(yamlContent), &yamlWorkflow); err != nil {
		return nil, fmt.Errorf("parse yaml failed: %w", err)
	}
	
	// 转换为模型
	workflow := &model.Workflow{
		ID:          generateWorkflowID(),
		Name:        yamlWorkflow.Name,
		Description: yamlWorkflow.Description,
		AutoPublish: yamlWorkflow.AutoPublish,
		NeedReview:  yamlWorkflow.NeedReview,
		Status:      model.WorkflowStatusActive,
	}
	
	// 转换触发器
	if yamlWorkflow.Trigger != nil {
		triggerType, _ := yamlWorkflow.Trigger["type"].(string)
		workflow.Trigger = model.WorkflowTrigger{
			Type:       triggerType,
			CronExpr:   getString(yamlWorkflow.Trigger, "cron"),
			RSSURL:     getString(yamlWorkflow.Trigger, "rss_url"),
			WebhookURL: getString(yamlWorkflow.Trigger, "webhook_url"),
			Config:     yamlWorkflow.Trigger,
		}
	}
	
	// 转换步骤
	for _, step := range yamlWorkflow.Steps {
		workflow.Steps = append(workflow.Steps, model.WorkflowStep{
			ID:       step.ID,
			Name:     step.Name,
			Type:     step.Type,
			Config:   step.Config,
			NextStep: step.NextStep,
			OnError:  step.OnError,
		})
	}
	
	// 保存
	if err := s.CreateWorkflow(workflow); err != nil {
		return nil, err
	}
	
	return workflow, nil
}

// ExportToYAML 导出工作流到YAML
func (s *WorkflowService) ExportToYAML(workflowID string) (string, error) {
	workflow, err := s.db.GetWorkflow(workflowID)
	if err != nil {
		return "", err
	}
	
	// 转换为YAML结构
	yamlWorkflow := struct {
		Name        string                 `yaml:"name"`
		Description string                 `yaml:"description"`
		Trigger     map[string]interface{} `yaml:"trigger"`
		Steps       []map[string]interface{} `yaml:"steps"`
		AutoPublish bool                   `yaml:"auto_publish"`
		NeedReview  bool                   `yaml:"need_review"`
	}{
		Name:        workflow.Name,
		Description: workflow.Description,
		AutoPublish: workflow.AutoPublish,
		NeedReview:  workflow.NeedReview,
	}
	
	// 转换触发器
	if workflow.Trigger.Type != "" {
		yamlWorkflow.Trigger = map[string]interface{}{
			"type": workflow.Trigger.Type,
		}
		if workflow.Trigger.CronExpr != "" {
			yamlWorkflow.Trigger["cron"] = workflow.Trigger.CronExpr
		}
		if workflow.Trigger.RSSURL != "" {
			yamlWorkflow.Trigger["rss_url"] = workflow.Trigger.RSSURL
		}
		if workflow.Trigger.WebhookURL != "" {
			yamlWorkflow.Trigger["webhook_url"] = workflow.Trigger.WebhookURL
		}
	}
	
	// 转换步骤
	for _, step := range workflow.Steps {
		stepMap := map[string]interface{}{
			"id":   step.ID,
			"name": step.Name,
			"type": step.Type,
		}
		if step.Config != nil {
			stepMap["config"] = step.Config
		}
		if step.NextStep != "" {
			stepMap["next"] = step.NextStep
		}
		if step.OnError != "" {
			stepMap["on_error"] = step.OnError
		}
		yamlWorkflow.Steps = append(yamlWorkflow.Steps, stepMap)
	}
	
	bytes, err := yaml.Marshal(yamlWorkflow)
	if err != nil {
		return "", err
	}
	
	return string(bytes), nil
}

// LoadFromFile 从文件加载工作流
func (s *WorkflowService) LoadFromFile(filePath string) (*model.Workflow, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file failed: %w", err)
	}
	
	// 根据扩展名判断格式
	ext := filepath.Ext(filePath)
	switch ext {
	case ".yaml", ".yml":
		return s.ImportFromYAML(string(content))
	case ".json":
		var workflow model.Workflow
		if err := json.Unmarshal(content, &workflow); err != nil {
			return nil, fmt.Errorf("parse json failed: %w", err)
		}
		if err := s.CreateWorkflow(&workflow); err != nil {
			return nil, err
		}
		return &workflow, nil
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}

// validateWorkflow 验证工作流
func (s *WorkflowService) validateWorkflow(workflow *model.Workflow) error {
	if workflow.Name == "" {
		return fmt.Errorf("workflow name is required")
	}
	
	if len(workflow.Steps) == 0 {
		return fmt.Errorf("workflow must have at least one step")
	}
	
	// 验证步骤ID唯一性
	stepIDs := make(map[string]bool)
	for _, step := range workflow.Steps {
		if step.ID == "" {
			return fmt.Errorf("step ID is required")
		}
		if stepIDs[step.ID] {
			return fmt.Errorf("duplicate step ID: %s", step.ID)
		}
		stepIDs[step.ID] = true
		
		// 验证步骤类型
		if step.Type == "" {
			return fmt.Errorf("step type is required")
		}
	}
	
	return nil
}

// getString 从map中获取字符串值
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// generateWorkflowID 生成工作流ID
func generateWorkflowID() string {
	return fmt.Sprintf("wf_%d", time.Now().UnixNano())
}
