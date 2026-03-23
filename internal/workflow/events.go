package workflow

import (
	"sync"
	"time"

	"content-alchemist/internal/model"
)

// EventType 事件类型
type EventType string

const (
	// 工作流事件
	EventWorkflowStarted   EventType = "workflow:started"
	EventWorkflowCompleted EventType = "workflow:completed"
	EventWorkflowFailed    EventType = "workflow:failed"
	EventWorkflowPaused    EventType = "workflow:paused"
	EventWorkflowCancelled EventType = "workflow:cancelled"

	// 步骤事件
	EventStepStarted   EventType = "step:started"
	EventStepCompleted EventType = "step:completed"
	EventStepFailed    EventType = "step:failed"
	EventStepRetry     EventType = "step:retry"
	EventStepSkipped   EventType = "step:skipped"
)

// Event 事件
type Event struct {
	Type      EventType              `json:"type"`
	RunID     string                 `json:"runId"`
	StepID    string                 `json:"stepId,omitempty"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// EventHandler 事件处理器
type EventHandler func(event Event)

// EventBus 事件总线
type EventBus struct {
	handlers map[EventType][]EventHandler
	mu       sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType][]EventHandler),
	}
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// Unsubscribe 取消订阅事件
func (eb *EventBus) Unsubscribe(eventType EventType, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	handlers := eb.handlers[eventType]
	for i, h := range handlers {
		if &h == &handler {
			eb.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// Publish 发布事件
func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	handlers := eb.handlers[event.Type]
	eb.mu.RUnlock()

	for _, handler := range handlers {
		go handler(event)
	}
}

// PublishWorkflowEvent 发布工作流事件
func (eb *EventBus) PublishWorkflowEvent(eventType EventType, run *model.WorkflowRun, data map[string]interface{}) {
	event := Event{
		Type:      eventType,
		RunID:     run.ID,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}
	eb.Publish(event)
}

// PublishStepEvent 发布步骤事件
func (eb *EventBus) PublishStepEvent(eventType EventType, run *model.WorkflowRun, step *model.WorkflowStep, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	event := Event{
		Type:      eventType,
		RunID:     run.ID,
		StepID:    step.ID,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}
	eb.Publish(event)
}

// WorkflowEventEmitter 工作流事件发射器
type WorkflowEventEmitter struct {
	bus *EventBus
}

// NewWorkflowEventEmitter 创建工作流事件发射器
func NewWorkflowEventEmitter(bus *EventBus) *WorkflowEventEmitter {
	return &WorkflowEventEmitter{bus: bus}
}

// EmitWorkflowStarted 发射工作流开始事件
func (e *WorkflowEventEmitter) EmitWorkflowStarted(run *model.WorkflowRun) {
	e.bus.PublishWorkflowEvent(EventWorkflowStarted, run, map[string]interface{}{
		"workflowId": run.WorkflowID,
		"input":      run.Input,
	})
}

// EmitWorkflowCompleted 发射工作流完成事件
func (e *WorkflowEventEmitter) EmitWorkflowCompleted(run *model.WorkflowRun) {
	e.bus.PublishWorkflowEvent(EventWorkflowCompleted, run, map[string]interface{}{
		"workflowId":  run.WorkflowID,
		"output":      run.Output,
		"duration":    run.CompletedAt.Sub(run.StartedAt).Seconds(),
	})
}

// EmitWorkflowFailed 发射工作流失败事件
func (e *WorkflowEventEmitter) EmitWorkflowFailed(run *model.WorkflowRun, err error) {
	e.bus.PublishWorkflowEvent(EventWorkflowFailed, run, map[string]interface{}{
		"workflowId": run.WorkflowID,
		"error":      err.Error(),
		"stepId":     run.CurrentStep,
	})
}

// EmitStepStarted 发射步骤开始事件
func (e *WorkflowEventEmitter) EmitStepStarted(run *model.WorkflowRun, step *model.WorkflowStep) {
	e.bus.PublishStepEvent(EventStepStarted, run, step, map[string]interface{}{
		"stepType": step.Type,
		"stepName": step.Name,
	})
}

// EmitStepCompleted 发射步骤完成事件
func (e *WorkflowEventEmitter) EmitStepCompleted(run *model.WorkflowRun, step *model.WorkflowStep, output interface{}) {
	e.bus.PublishStepEvent(EventStepCompleted, run, step, map[string]interface{}{
		"stepType": step.Type,
		"stepName": step.Name,
		"output":   output,
	})
}

// EmitStepFailed 发射步骤失败事件
func (e *WorkflowEventEmitter) EmitStepFailed(run *model.WorkflowRun, step *model.WorkflowStep, err error) {
	e.bus.PublishStepEvent(EventStepFailed, run, step, map[string]interface{}{
		"stepType": step.Type,
		"stepName": step.Name,
		"error":    err.Error(),
	})
}

// EmitStepRetry 发射步骤重试事件
func (e *WorkflowEventEmitter) EmitStepRetry(run *model.WorkflowRun, step *model.WorkflowStep, retryCount int) {
	e.bus.PublishStepEvent(EventStepRetry, run, step, map[string]interface{}{
		"stepType":   step.Type,
		"stepName":   step.Name,
		"retryCount": retryCount,
	})
}
