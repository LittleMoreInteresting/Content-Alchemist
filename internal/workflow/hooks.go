package workflow

import (
	"content-alchemist/internal/model"
)

// HookType 钩子类型
type HookType string

const (
	// 工作流钩子
	BeforeWorkflow HookType = "before:workflow"
	AfterWorkflow  HookType = "after:workflow"

	// 步骤钩子
	BeforeStep HookType = "before:step"
	AfterStep  HookType = "after:step"

	// 错误钩子
	OnStepError HookType = "on:step:error"
	OnRunError  HookType = "on:run:error"

	// 人工审核钩子
	OnManualReview HookType = "on:manual:review"
)

// HookHandler 钩子处理器
type HookHandler func(run *model.WorkflowRun, step *model.WorkflowStep, data interface{}) error

// Hooks 钩子系统
type Hooks struct {
	handlers map[HookType][]HookHandler
}

// NewHooks 创建钩子系统
func NewHooks() *Hooks {
	return &Hooks{
		handlers: make(map[HookType][]HookHandler),
	}
}

// Register 注册钩子
func (h *Hooks) Register(hookType HookType, handler HookHandler) {
	h.handlers[hookType] = append(h.handlers[hookType], handler)
}

// Trigger 触发钩子
func (h *Hooks) Trigger(hookType HookType, run *model.WorkflowRun, step *model.WorkflowStep, data interface{}) error {
	handlers := h.handlers[hookType]
	for _, handler := range handlers {
		if err := handler(run, step, data); err != nil {
			return err
		}
	}
	return nil
}

// TriggerWithResult 触发钩子并返回结果
func (h *Hooks) TriggerWithResult(hookType HookType, run *model.WorkflowRun, step *model.WorkflowStep, data interface{}) (interface{}, error) {
	handlers := h.handlers[hookType]
	var result interface{}
	for _, handler := range handlers {
		r, err := handlerWithResult(handler, run, step, data)
		if err != nil {
			return nil, err
		}
		if r != nil {
			result = r
		}
	}
	return result, nil
}

// handlerWithResult 带结果的处理器
func handlerWithResult(handler HookHandler, run *model.WorkflowRun, step *model.WorkflowStep, data interface{}) (interface{}, error) {
	err := handler(run, step, data)
	return nil, err
}

// DefaultHooks 创建默认钩子
func DefaultHooks() *Hooks {
	hooks := NewHooks()

	// 注册日志钩子
	hooks.Register(BeforeWorkflow, func(run *model.WorkflowRun, step *model.WorkflowStep, data interface{}) error {
		// TODO: 记录日志
		return nil
	})

	hooks.Register(AfterWorkflow, func(run *model.WorkflowRun, step *model.WorkflowStep, data interface{}) error {
		// TODO: 记录日志
		return nil
	})

	hooks.Register(BeforeStep, func(run *model.WorkflowRun, step *model.WorkflowStep, data interface{}) error {
		// TODO: 记录日志
		return nil
	})

	hooks.Register(AfterStep, func(run *model.WorkflowRun, step *model.WorkflowStep, data interface{}) error {
		// TODO: 记录日志
		return nil
	})

	return hooks
}

// ManualReviewHook 人工审核钩子
type ManualReviewHook struct {
	callback func(run *model.WorkflowRun, step *model.WorkflowStep, data interface{}) (bool, error)
}

// NewManualReviewHook 创建人工审核钩子
func NewManualReviewHook(callback func(run *model.WorkflowRun, step *model.WorkflowStep, data interface{}) (bool, error)) *ManualReviewHook {
	return &ManualReviewHook{callback: callback}
}

// Execute 执行人工审核
func (h *ManualReviewHook) Execute(run *model.WorkflowRun, step *model.WorkflowStep, data interface{}) (bool, error) {
	return h.callback(run, step, data)
}

// ManualReviewManager 人工审核管理器
type ManualReviewManager struct {
	pending map[string]*ManualReviewRequest
}

// ManualReviewRequest 人工审核请求
type ManualReviewRequest struct {
	RunID     string
	StepID    string
	Run       *model.WorkflowRun
	Step      *model.WorkflowStep
	Data      interface{}
	Result    chan ManualReviewResult
}

// ManualReviewResult 人工审核结果
type ManualReviewResult struct {
	Approved bool
	Data     interface{}
	Error    error
}

// NewManualReviewManager 创建人工审核管理器
func NewManualReviewManager() *ManualReviewManager {
	return &ManualReviewManager{
		pending: make(map[string]*ManualReviewRequest),
	}
}

// SubmitReview 提交审核请求
func (m *ManualReviewManager) SubmitReview(req *ManualReviewRequest) {
	key := req.RunID + ":" + req.StepID
	m.pending[key] = req
}

// GetReview 获取审核请求
func (m *ManualReviewManager) GetReview(runID, stepID string) *ManualReviewRequest {
	key := runID + ":" + stepID
	return m.pending[key]
}

// CompleteReview 完成审核
func (m *ManualReviewManager) CompleteReview(runID, stepID string, result ManualReviewResult) {
	key := runID + ":" + stepID
	if req, ok := m.pending[key]; ok {
		req.Result <- result
		delete(m.pending, key)
	}
}

// ListPendingReviews 列出待审核的请求
func (m *ManualReviewManager) ListPendingReviews() []*ManualReviewRequest {
	var requests []*ManualReviewRequest
	for _, req := range m.pending {
		requests = append(requests, req)
	}
	return requests
}
