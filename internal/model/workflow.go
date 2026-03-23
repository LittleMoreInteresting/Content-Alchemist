package model

import (
	"encoding/json"
	"time"
)

// ==================== 工作流相关模型 ====================

// Workflow 工作流定义
type Workflow struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Trigger     WorkflowTrigger `json:"trigger"`     // 触发器配置
	Steps       []WorkflowStep  `json:"steps"`       // 执行步骤
	AutoPublish bool            `json:"autoPublish"` // 是否自动发布
	NeedReview  bool            `json:"needReview"`  // 是否需要人工审核
	Status      string          `json:"status"`      // active, paused, archived
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// WorkflowTrigger 工作流触发器
type WorkflowTrigger struct {
	Type       string         `json:"type"`       // manual, schedule, webhook, rss
	CronExpr   string         `json:"cronExpr"`   // 定时表达式
	RSSURL     string         `json:"rssUrl"`     // RSS源地址
	WebhookURL string         `json:"webhookUrl"` // Webhook地址
	Config     map[string]any `json:"config"`     // 额外配置
}

// WorkflowStep 工作流步骤
type WorkflowStep struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Type     string         `json:"type"`     // topic_fetch, topic_select, outline_generate, content_write, content_review, layout_apply, image_generate, publish_draft, condition, delay, manual_review
	Config   map[string]any `json:"config"`   // 步骤配置
	NextStep string         `json:"nextStep"` // 下一步ID
	OnError  string         `json:"onError"`  // 错误处理: retry, skip, abort, manual
}

// WorkflowRun 工作流运行实例
type WorkflowRun struct {
	ID          string            `json:"id"`
	WorkflowID  string            `json:"workflowId"`
	Status      string            `json:"status"`      // pending, running, paused, completed, failed
	CurrentStep string            `json:"currentStep"`
	Input       map[string]any    `json:"input"`       // 输入参数
	Output      map[string]any    `json:"output"`      // 输出结果
	Steps       []WorkflowRunStep `json:"steps"`       // 各步骤执行状态
	StartedAt   time.Time         `json:"startedAt"`
	CompletedAt *time.Time        `json:"completedAt"`
	Error       string            `json:"error"`
}

// WorkflowRunStep 工作流步骤执行状态
type WorkflowRunStep struct {
	StepID      string     `json:"stepId"`
	Status      string     `json:"status"`    // pending, running, completed, failed, skipped
	Input       any        `json:"input"`
	Output      any        `json:"output"`
	StartedAt   time.Time  `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt"`
	Error       string     `json:"error"`
	RetryCount  int        `json:"retryCount"`
}

// ErrorAction 错误处理动作
type ErrorAction string

const (
	ErrorActionRetry  ErrorAction = "retry"
	ErrorActionSkip   ErrorAction = "skip"
	ErrorActionAbort  ErrorAction = "abort"
	ErrorActionManual ErrorAction = "manual"
)

// WorkflowStatus 工作流状态
const (
	WorkflowStatusActive    = "active"
	WorkflowStatusPaused    = "paused"
	WorkflowStatusArchived  = "archived"
	WorkflowStatusRunning   = "running"
	WorkflowStatusCompleted = "completed"
	WorkflowStatusFailed    = "failed"
)

// StepStatus 步骤状态
const (
	StepStatusPending   = "pending"
	StepStatusRunning   = "running"
	StepStatusCompleted = "completed"
	StepStatusFailed    = "failed"
	StepStatusSkipped   = "skipped"
)

// StepType 步骤类型
const (
	StepTypeHotFetch       = "hot_fetch"
	StepTypeTopicGenerate  = "topic_generate"
	StepTypeTopicSelect    = "topic_select"
	StepTypeOutlineGenerate = "outline_generate"
	StepTypeContentWrite   = "content_write"
	StepTypeContentReview  = "content_review"
	StepTypeContentRewrite = "content_rewrite"
	StepTypeImageGenerate  = "image_generate"
	StepTypeLayoutApply    = "layout_apply"
	StepTypePublishDraft   = "publish_draft"
	StepTypeCondition      = "condition"
	StepTypeDelay          = "delay"
	StepTypeManualReview   = "manual_review"
)

// TriggerType 触发器类型
const (
	TriggerTypeManual  = "manual"
	TriggerTypeSchedule = "schedule"
	TriggerTypeWebhook = "webhook"
	TriggerTypeRSS     = "rss"
)

// ToJSON 将工作流转换为JSON字符串
func (w *Workflow) ToJSON() (string, error) {
	bytes, err := json.Marshal(w)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// WorkflowFromJSON 从JSON字符串解析工作流
func WorkflowFromJSON(data string) (*Workflow, error) {
	var w Workflow
	if err := json.Unmarshal([]byte(data), &w); err != nil {
		return nil, err
	}
	return &w, nil
}

// GetStepByID 根据ID获取步骤
func (w *Workflow) GetStepByID(id string) *WorkflowStep {
	for i := range w.Steps {
		if w.Steps[i].ID == id {
			return &w.Steps[i]
		}
	}
	return nil
}

// GetFirstStep 获取第一个步骤
func (w *Workflow) GetFirstStep() *WorkflowStep {
	if len(w.Steps) == 0 {
		return nil
	}
	return &w.Steps[0]
}

// GetNextStep 根据当前步骤获取下一步
func (w *Workflow) GetNextStep(currentStepID string, conditionResult bool) *WorkflowStep {
	currentStep := w.GetStepByID(currentStepID)
	if currentStep == nil {
		return nil
	}
	
	// 如果是条件步骤，根据结果选择分支
	if currentStep.Type == StepTypeCondition {
		// 条件步骤的config中应该包含true_next和false_next
		config := currentStep.Config
		if config != nil {
			if conditionResult {
				if nextID, ok := config["true_next"].(string); ok {
					return w.GetStepByID(nextID)
				}
			} else {
				if nextID, ok := config["false_next"].(string); ok {
					return w.GetStepByID(nextID)
				}
			}
		}
	}
	
	if currentStep.NextStep == "" {
		return nil
	}
	
	return w.GetStepByID(currentStep.NextStep)
}
