package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"content-alchemist/internal/model"
)

// ============ Workflow 操作 ============

// SaveWorkflow 保存工作流
func (d *DB) SaveWorkflow(workflow *model.Workflow) error {
	workflow.UpdatedAt = time.Now()
	if workflow.CreatedAt.IsZero() {
		workflow.CreatedAt = time.Now()
	}

	triggerJSON, _ := json.Marshal(workflow.Trigger)
	stepsJSON, _ := json.Marshal(workflow.Steps)

	query := `INSERT OR REPLACE INTO workflows 
		(id, name, description, trigger_config, steps, auto_publish, need_review, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := d.conn.Exec(query,
		workflow.ID,
		workflow.Name,
		workflow.Description,
		string(triggerJSON),
		string(stepsJSON),
		workflow.AutoPublish,
		workflow.NeedReview,
		workflow.Status,
		workflow.CreatedAt,
		workflow.UpdatedAt,
	)
	return err
}

// GetWorkflow 获取工作流
func (d *DB) GetWorkflow(id string) (*model.Workflow, error) {
	query := `SELECT id, name, description, trigger_config, steps, auto_publish, need_review, status, created_at, updated_at 
		FROM workflows WHERE id = ?`

	row := d.conn.QueryRow(query, id)

	var workflow model.Workflow
	var triggerStr, stepsStr string

	err := row.Scan(
		&workflow.ID,
		&workflow.Name,
		&workflow.Description,
		&triggerStr,
		&stepsStr,
		&workflow.AutoPublish,
		&workflow.NeedReview,
		&workflow.Status,
		&workflow.CreatedAt,
		&workflow.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("workflow not found")
	}
	if err != nil {
		return nil, err
	}

	if triggerStr != "" {
		json.Unmarshal([]byte(triggerStr), &workflow.Trigger)
	}
	if stepsStr != "" {
		json.Unmarshal([]byte(stepsStr), &workflow.Steps)
	}

	return &workflow, nil
}

// ListWorkflows 列出所有工作流
func (d *DB) ListWorkflows() ([]model.Workflow, error) {
	query := `SELECT id, name, description, trigger_config, steps, auto_publish, need_review, status, created_at, updated_at 
		FROM workflows ORDER BY updated_at DESC`

	rows, err := d.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []model.Workflow
	for rows.Next() {
		var workflow model.Workflow
		var triggerStr, stepsStr string

		err := rows.Scan(
			&workflow.ID,
			&workflow.Name,
			&workflow.Description,
			&triggerStr,
			&stepsStr,
			&workflow.AutoPublish,
			&workflow.NeedReview,
			&workflow.Status,
			&workflow.CreatedAt,
			&workflow.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if triggerStr != "" {
			json.Unmarshal([]byte(triggerStr), &workflow.Trigger)
		}
		if stepsStr != "" {
			json.Unmarshal([]byte(stepsStr), &workflow.Steps)
		}

		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

// DeleteWorkflow 删除工作流
func (d *DB) DeleteWorkflow(id string) error {
	_, err := d.conn.Exec("DELETE FROM workflows WHERE id = ?", id)
	return err
}

// ============ WorkflowRun 操作 ============

// SaveWorkflowRun 保存工作流运行实例
func (d *DB) SaveWorkflowRun(run *model.WorkflowRun) error {
	inputJSON, _ := json.Marshal(run.Input)
	outputJSON, _ := json.Marshal(run.Output)
	stepsJSON, _ := json.Marshal(run.Steps)

	query := `INSERT OR REPLACE INTO workflow_runs 
		(id, workflow_id, status, current_step, input, output, steps, started_at, completed_at, error)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := d.conn.Exec(query,
		run.ID,
		run.WorkflowID,
		run.Status,
		run.CurrentStep,
		string(inputJSON),
		string(outputJSON),
		string(stepsJSON),
		run.StartedAt,
		run.CompletedAt,
		run.Error,
	)
	return err
}

// GetWorkflowRun 获取工作流运行实例
func (d *DB) GetWorkflowRun(id string) (*model.WorkflowRun, error) {
	query := `SELECT id, workflow_id, status, current_step, input, output, steps, started_at, completed_at, error 
		FROM workflow_runs WHERE id = ?`

	row := d.conn.QueryRow(query, id)

	var run model.WorkflowRun
	var inputStr, outputStr, stepsStr string

	err := row.Scan(
		&run.ID,
		&run.WorkflowID,
		&run.Status,
		&run.CurrentStep,
		&inputStr,
		&outputStr,
		&stepsStr,
		&run.StartedAt,
		&run.CompletedAt,
		&run.Error,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("workflow run not found")
	}
	if err != nil {
		return nil, err
	}

	if inputStr != "" {
		json.Unmarshal([]byte(inputStr), &run.Input)
	}
	if outputStr != "" {
		json.Unmarshal([]byte(outputStr), &run.Output)
	}
	if stepsStr != "" {
		json.Unmarshal([]byte(stepsStr), &run.Steps)
	}

	return &run, nil
}

// ListWorkflowRuns 列出工作流运行实例
func (d *DB) ListWorkflowRuns(workflowID string) ([]model.WorkflowRun, error) {
	var query string
	var args []interface{}

	if workflowID == "" {
		query = `SELECT id, workflow_id, status, current_step, input, output, steps, started_at, completed_at, error 
			FROM workflow_runs ORDER BY started_at DESC`
	} else {
		query = `SELECT id, workflow_id, status, current_step, input, output, steps, started_at, completed_at, error 
			FROM workflow_runs WHERE workflow_id = ? ORDER BY started_at DESC`
		args = append(args, workflowID)
	}

	rows, err := d.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []model.WorkflowRun
	for rows.Next() {
		var run model.WorkflowRun
		var inputStr, outputStr, stepsStr string

		err := rows.Scan(
			&run.ID,
			&run.WorkflowID,
			&run.Status,
			&run.CurrentStep,
			&inputStr,
			&outputStr,
			&stepsStr,
			&run.StartedAt,
			&run.CompletedAt,
			&run.Error,
		)
		if err != nil {
			continue
		}

		if inputStr != "" {
			json.Unmarshal([]byte(inputStr), &run.Input)
		}
		if outputStr != "" {
			json.Unmarshal([]byte(outputStr), &run.Output)
		}
		if stepsStr != "" {
			json.Unmarshal([]byte(stepsStr), &run.Steps)
		}

		runs = append(runs, run)
	}

	return runs, nil
}

// ListActiveWorkflowRuns 列出活跃的工作流运行实例
func (d *DB) ListActiveWorkflowRuns() ([]model.WorkflowRun, error) {
	query := `SELECT id, workflow_id, status, current_step, input, output, steps, started_at, completed_at, error 
		FROM workflow_runs WHERE status IN ('pending', 'running') ORDER BY started_at ASC`

	rows, err := d.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []model.WorkflowRun
	for rows.Next() {
		var run model.WorkflowRun
		var inputStr, outputStr, stepsStr string

		err := rows.Scan(
			&run.ID,
			&run.WorkflowID,
			&run.Status,
			&run.CurrentStep,
			&inputStr,
			&outputStr,
			&stepsStr,
			&run.StartedAt,
			&run.CompletedAt,
			&run.Error,
		)
		if err != nil {
			continue
		}

		if inputStr != "" {
			json.Unmarshal([]byte(inputStr), &run.Input)
		}
		if outputStr != "" {
			json.Unmarshal([]byte(outputStr), &run.Output)
		}
		if stepsStr != "" {
			json.Unmarshal([]byte(stepsStr), &run.Steps)
		}

		runs = append(runs, run)
	}

	return runs, nil
}

// DeleteWorkflowRun 删除工作流运行实例
func (d *DB) DeleteWorkflowRun(id string) error {
	_, err := d.conn.Exec("DELETE FROM workflow_runs WHERE id = ?", id)
	return err
}
