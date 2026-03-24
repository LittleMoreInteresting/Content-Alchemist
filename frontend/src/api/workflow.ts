import { Workflow, WorkflowRun } from '../types';

// 创建工作流
export async function CreateWorkflow(workflow: Workflow): Promise<Workflow> {
  return await window.go.main.App.CreateWorkflow(workflow);
}

// 获取工作流
export async function GetWorkflow(id: string): Promise<Workflow> {
  return await window.go.main.App.GetWorkflow(id);
}

// 列出所有工作流
export async function ListWorkflows(): Promise<Workflow[]> {
  return await window.go.main.App.ListWorkflows();
}

// 更新工作流
export async function UpdateWorkflow(workflow: Workflow): Promise<void> {
  return await window.go.main.App.UpdateWorkflow(workflow);
}

// 删除工作流
export async function DeleteWorkflow(id: string): Promise<void> {
  return await window.go.main.App.DeleteWorkflow(id);
}

// 启动工作流
export async function StartWorkflow(
  workflowId: string,
  input: Record<string, any>
): Promise<WorkflowRun> {
  return await window.go.main.App.StartWorkflow(workflowId, input);
}

// 获取工作流运行状态
export async function GetWorkflowRun(runId: string): Promise<WorkflowRun> {
  return await window.go.main.App.GetWorkflowRun(runId);
}

// 列出工作流运行实例
export async function ListWorkflowRuns(workflowId: string): Promise<WorkflowRun[]> {
  return await window.go.main.App.ListWorkflowRuns(workflowId);
}

// 取消工作流运行
export async function CancelWorkflow(runId: string): Promise<void> {
  return await window.go.main.App.CancelWorkflow(runId);
}

// 获取Agent注册表
export async function GetAgentRegistry(): Promise<any[]> {
  return await window.go.main.App.GetAgentRegistry();
}
