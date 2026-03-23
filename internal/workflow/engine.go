package workflow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"content-alchemist/internal/model"
	"content-alchemist/internal/repository"
	"content-alchemist/internal/workflow/agents"
)

// Agent 工作流Agent接口
type Agent interface {
	Execute(ctx context.Context, input any) (any, error)
}

// AgentFactory Agent工厂函数
type AgentFactory func() Agent

// Engine 工作流引擎
type Engine struct {
	db         *repository.DB
	agents     map[string]AgentFactory
	registry   *AgentRegistry
	events     *EventBus
	hooks      *Hooks
	runs       map[string]*RunContext
	mu         sync.RWMutex
	stopCh     chan struct{}
}

// RunContext 运行上下文
type RunContext struct {
	Run      *model.WorkflowRun
	Workflow *model.Workflow
	Cancel   context.CancelFunc
}

// NewEngine 创建工作流引擎
func NewEngine(db *repository.DB) *Engine {
	engine := &Engine{
		db:       db,
		agents:   make(map[string]AgentFactory),
		registry: NewAgentRegistry(),
		events:   NewEventBus(),
		hooks:    NewHooks(),
		runs:     make(map[string]*RunContext),
		stopCh:   make(chan struct{}),
	}

	// 注册内置Agent
	engine.registerBuiltinAgents()

	return engine
}

// RegisterAgent 注册Agent
func (e *Engine) RegisterAgent(agentType string, factory AgentFactory) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.agents[agentType] = factory
}

// Start 启动工作流
func (e *Engine) Start(ctx context.Context, workflowID string, input map[string]any) (*model.WorkflowRun, error) {
	// 1. 加载工作流定义
	workflow, err := e.db.GetWorkflow(workflowID)
	if err != nil {
		return nil, fmt.Errorf("get workflow failed: %w", err)
	}

	// 2. 检查工作流状态
	if workflow.Status != model.WorkflowStatusActive {
		return nil, fmt.Errorf("workflow is not active: %s", workflow.Status)
	}

	// 3. 创建运行实例
	run := &model.WorkflowRun{
		ID:         generateID(),
		WorkflowID: workflowID,
		Status:     model.StepStatusRunning,
		Input:      input,
		Steps:      make([]model.WorkflowRunStep, len(workflow.Steps)),
		StartedAt:  time.Now(),
	}

	// 4. 初始化步骤状态
	for i, step := range workflow.Steps {
		run.Steps[i] = model.WorkflowRunStep{
			StepID:    step.ID,
			Status:    model.StepStatusPending,
			RetryCount: 0,
		}
	}

	// 5. 保存运行状态
	if err := e.db.SaveWorkflowRun(run); err != nil {
		return nil, fmt.Errorf("save workflow run failed: %w", err)
	}

	// 6. 启动执行协程
	runCtx, cancel := context.WithCancel(ctx)
	e.mu.Lock()
	e.runs[run.ID] = &RunContext{
		Run:      run,
		Workflow: workflow,
		Cancel:   cancel,
	}
	e.mu.Unlock()

	go e.execute(runCtx, run, workflow)

	return run, nil
}

// Pause 暂停工作流运行
func (e *Engine) Pause(runID string) error {
	e.mu.RLock()
	ctx, ok := e.runs[runID]
	e.mu.RUnlock()

	if !ok {
		return fmt.Errorf("workflow run not found: %s", runID)
	}

	ctx.Run.Status = model.WorkflowStatusPaused
	return e.db.SaveWorkflowRun(ctx.Run)
}

// Resume 恢复工作流运行
func (e *Engine) Resume(ctx context.Context, runID string) error {
	// 加载运行实例
	run, err := e.db.GetWorkflowRun(runID)
	if err != nil {
		return err
	}

	if run.Status != model.WorkflowStatusPaused {
		return fmt.Errorf("workflow run is not paused: %s", run.Status)
	}

	// 加载工作流定义
	workflow, err := e.db.GetWorkflow(run.WorkflowID)
	if err != nil {
		return err
	}

	// 恢复执行
	run.Status = model.StepStatusRunning
	e.db.SaveWorkflowRun(run)

	runCtx, cancel := context.WithCancel(ctx)
	e.mu.Lock()
	e.runs[runID] = &RunContext{
		Run:      run,
		Workflow: workflow,
		Cancel:   cancel,
	}
	e.mu.Unlock()

	go e.execute(runCtx, run, workflow)

	return nil
}

// Cancel 取消工作流运行
func (e *Engine) Cancel(runID string) error {
	e.mu.Lock()
	ctx, ok := e.runs[runID]
	if ok {
		ctx.Cancel()
		delete(e.runs, runID)
	}
	e.mu.Unlock()

	if !ok {
		// 从数据库加载并更新状态
		run, err := e.db.GetWorkflowRun(runID)
		if err != nil {
			return err
		}
		run.Status = "cancelled"
		return e.db.SaveWorkflowRun(run)
	}

	ctx.Run.Status = "cancelled"
	return e.db.SaveWorkflowRun(ctx.Run)
}

// GetRun 获取工作流运行状态
func (e *Engine) GetRun(runID string) (*model.WorkflowRun, error) {
	return e.db.GetWorkflowRun(runID)
}

// ListRuns 列出工作流运行实例
func (e *Engine) ListRuns(workflowID string) ([]model.WorkflowRun, error) {
	return e.db.ListWorkflowRuns(workflowID)
}

// execute 执行工作流
func (e *Engine) execute(ctx context.Context, run *model.WorkflowRun, workflow *model.Workflow) {
	defer func() {
		e.mu.Lock()
		delete(e.runs, run.ID)
		e.mu.Unlock()
	}()

	// 触发前置钩子
	e.hooks.Trigger(BeforeWorkflow, run, nil, nil)

	// 构建步骤映射
	stepMap := make(map[string]model.WorkflowStep)
	for _, step := range workflow.Steps {
		stepMap[step.ID] = step
	}

	// 找到当前应该执行的步骤
	currentStepID := e.getCurrentStepID(run)
	if currentStepID == "" {
		currentStepID = workflow.Steps[0].ID
	}

	for currentStepID != "" {
		select {
		case <-ctx.Done():
			run.Status = "cancelled"
			e.db.SaveWorkflowRun(run)
			return
		default:
		}

		// 检查是否被暂停
		if run.Status == model.WorkflowStatusPaused {
			return
		}

		step, ok := stepMap[currentStepID]
		if !ok {
			run.Status = model.WorkflowStatusFailed
			run.Error = fmt.Sprintf("step not found: %s", currentStepID)
			e.db.SaveWorkflowRun(run)
			return
		}

		// 执行步骤
		result, err := e.executeStep(ctx, run, step)

		if err != nil {
			// 错误处理
			handleAction := e.getErrorAction(step, run)
			switch handleAction {
			case model.ErrorActionRetry:
				// 重试逻辑在executeStep中处理
				continue
			case model.ErrorActionSkip:
				currentStepID = step.NextStep
				continue
			case model.ErrorActionAbort:
				run.Status = model.WorkflowStatusFailed
				run.Error = err.Error()
				e.db.SaveWorkflowRun(run)
				e.hooks.Trigger(AfterWorkflow, run, nil, err)
				return
			case model.ErrorActionManual:
				run.Status = model.WorkflowStatusPaused
				e.db.SaveWorkflowRun(run)
				return
			}
		}

		// 确定下一步
		currentStepID = e.getNextStep(step, result, run)

		// 触发步骤后置钩子
		e.hooks.Trigger(AfterStep, run, &step, result)
	}

	// 完成
	run.Status = model.WorkflowStatusCompleted
	now := time.Now()
	run.CompletedAt = &now
	e.db.SaveWorkflowRun(run)

	// 触发工作流后置钩子
	e.hooks.Trigger(AfterWorkflow, run, nil, nil)
}

// executeStep 执行单个步骤
func (e *Engine) executeStep(ctx context.Context, run *model.WorkflowRun, step model.WorkflowStep) (any, error) {
	// 查找步骤索引
	stepIndex := -1
	for i, s := range run.Steps {
		if s.StepID == step.ID {
			stepIndex = i
			break
		}
	}
	if stepIndex == -1 {
		return nil, fmt.Errorf("step not found in run: %s", step.ID)
	}

	// 更新步骤状态为运行中
	run.Steps[stepIndex].Status = model.StepStatusRunning
	run.Steps[stepIndex].StartedAt = time.Now()
	e.db.SaveWorkflowRun(run)

	// 触发步骤前置钩子
	e.hooks.Trigger(BeforeStep, run, &step, nil)

	// 获取Agent
	factory, ok := e.agents[step.Type]
	if !ok {
		err := fmt.Errorf("agent not registered: %s", step.Type)
		e.updateStepFailed(run, stepIndex, err)
		return nil, err
	}

	// 构建输入
	input := e.buildStepInput(run, step)
	run.Steps[stepIndex].Input = input
	e.db.SaveWorkflowRun(run)

	// 执行（带重试）
	var output any
	var err error
	maxRetries := 3
	for attempt := 0; attempt <= maxRetries; attempt++ {
		agent := factory()
		output, err = agent.Execute(ctx, input)
		if err == nil {
			break
		}
		if attempt < maxRetries {
			time.Sleep(time.Duration(attempt+1) * time.Second)
		}
	}

	if err != nil {
		e.updateStepFailed(run, stepIndex, err)
		return nil, err
	}

	// 更新成功状态
	run.Steps[stepIndex].Status = model.StepStatusCompleted
	run.Steps[stepIndex].Output = output
	now := time.Now()
	run.Steps[stepIndex].CompletedAt = &now

	// 更新运行状态
	run.CurrentStep = step.ID
	run.Output = mergeOutput(run.Output, output)
	e.db.SaveWorkflowRun(run)

	return output, nil
}

// updateStepFailed 更新步骤失败状态
func (e *Engine) updateStepFailed(run *model.WorkflowRun, stepIndex int, err error) {
	run.Steps[stepIndex].Status = model.StepStatusFailed
	run.Steps[stepIndex].Error = err.Error()
	run.Steps[stepIndex].RetryCount++
	now := time.Now()
	run.Steps[stepIndex].CompletedAt = &now
	e.db.SaveWorkflowRun(run)
}

// getCurrentStepID 获取当前应该执行的步骤ID
func (e *Engine) getCurrentStepID(run *model.WorkflowRun) string {
	// 找到第一个未完成的步骤
	for _, step := range run.Steps {
		if step.Status == model.StepStatusPending {
			return step.StepID
		}
		if step.Status == model.StepStatusRunning {
			return step.StepID
		}
	}
	return ""
}

// getNextStep 根据当前步骤获取下一步
func (e *Engine) getNextStep(currentStep model.WorkflowStep, result any, run *model.WorkflowRun) string {
	// 如果是条件步骤，根据结果选择分支
	if currentStep.Type == model.StepTypeCondition && currentStep.Config != nil {
		conditionResult := false
		if result != nil {
			if r, ok := result.(map[string]any); ok {
				if v, ok := r["condition"].(bool); ok {
					conditionResult = v
				}
			}
		}

		if conditionResult {
			if nextID, ok := currentStep.Config["true_next"].(string); ok && nextID != "" {
				return nextID
			}
		} else {
			if nextID, ok := currentStep.Config["false_next"].(string); ok && nextID != "" {
				return nextID
			}
		}
	}

	return currentStep.NextStep
}

// getErrorAction 获取错误处理动作
func (e *Engine) getErrorAction(step model.WorkflowStep, run *model.WorkflowRun) model.ErrorAction {
	if step.OnError != "" {
		return model.ErrorAction(step.OnError)
	}
	return model.ErrorActionRetry
}

// buildStepInput 构建步骤输入
func (e *Engine) buildStepInput(run *model.WorkflowRun, step model.WorkflowStep) any {
	input := make(map[string]any)

	// 合并工作流输入
	for k, v := range run.Input {
		input[k] = v
	}

	// 合并之前步骤的输出
	for k, v := range run.Output {
		input[k] = v
	}

	// 合并步骤配置
	for k, v := range step.Config {
		input[k] = v
	}

	return input
}

// mergeOutput 合并输出
func mergeOutput(existing map[string]any, newOutput any) map[string]any {
	if existing == nil {
		existing = make(map[string]any)
	}

	if newOutput == nil {
		return existing
	}

	if output, ok := newOutput.(map[string]any); ok {
		for k, v := range output {
			existing[k] = v
		}
	}

	return existing
}

// generateID 生成唯一ID
func generateID() string {
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), time.Now().Unix())
}

// Stop 停止工作流引擎
func (e *Engine) Stop() {
	close(e.stopCh)

	// 取消所有正在运行的工作流
	e.mu.Lock()
	for _, ctx := range e.runs {
		ctx.Cancel()
	}
	e.mu.Unlock()
}

// GetEventBus 获取事件总线
func (e *Engine) GetEventBus() *EventBus {
	return e.events
}

// GetHooks 获取钩子系统
func (e *Engine) GetHooks() *Hooks {
	return e.hooks
}

// registerBuiltinAgents 注册内置Agent
func (e *Engine) registerBuiltinAgents() {
	// 简单Agent - 不需要外部服务
	e.RegisterAgent("delay", func() Agent {
		return &agents.DelayAgent{}
	})
	e.RegisterAgent("condition", func() Agent {
		return &agents.ConditionAgent{}
	})
	
	// 这些Agent需要在Start时传入对应的Service
	// 这里先注册占位，实际使用时再初始化
	e.RegisterAgent("hot_fetch", func() Agent {
		return &agents.HotFetchAgent{}
	})
	e.RegisterAgent("topic_generate", func() Agent {
		return &agents.TopicGenerateAgent{}
	})
	e.RegisterAgent("topic_select", func() Agent {
		return &agents.TopicSelectAgent{}
	})
	e.RegisterAgent("outline_generate", func() Agent {
		return &agents.OutlineGenerateAgent{}
	})
	e.RegisterAgent("content_write", func() Agent {
		return &agents.ContentWriteAgent{}
	})
	e.RegisterAgent("content_review", func() Agent {
		return &agents.ContentReviewAgent{}
	})
	e.RegisterAgent("layout_apply", func() Agent {
		return &agents.LayoutApplyAgent{}
	})
}

// GetRegistry 获取Agent注册表
func (e *Engine) GetRegistry() *AgentRegistry {
	return e.registry
}
