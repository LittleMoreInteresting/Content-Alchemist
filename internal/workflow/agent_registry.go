package workflow

import (
	"fmt"
	"sync"
)

// AgentRegistry Agent注册表
type AgentRegistry struct {
	agents map[string]AgentInfo
	mu     sync.RWMutex
}

// AgentInfo Agent信息
type AgentInfo struct {
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Icon        string            `json:"icon"`
	Category    string            `json:"category"`
	InputSchema map[string]any    `json:"inputSchema"`
	OutputSchema map[string]any   `json:"outputSchema"`
	ConfigSchema map[string]any   `json:"configSchema"`
}

// NewAgentRegistry 创建Agent注册表
func NewAgentRegistry() *AgentRegistry {
	registry := &AgentRegistry{
		agents: make(map[string]AgentInfo),
	}
	
	// 注册内置Agent
	registry.registerBuiltinAgents()
	
	return registry
}

// Register 注册Agent
func (r *AgentRegistry) Register(info AgentInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if info.Type == "" {
		return fmt.Errorf("agent type is required")
	}
	
	r.agents[info.Type] = info
	return nil
}

// Get 获取Agent信息
func (r *AgentRegistry) Get(agentType string) (AgentInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	info, ok := r.agents[agentType]
	return info, ok
}

// List 列出所有Agent
func (r *AgentRegistry) List() []AgentInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var list []AgentInfo
	for _, info := range r.agents {
		list = append(list, info)
	}
	
	return list
}

// ListByCategory 按分类列出Agent
func (r *AgentRegistry) ListByCategory(category string) []AgentInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var list []AgentInfo
	for _, info := range r.agents {
		if info.Category == category {
			list = append(list, info)
		}
	}
	
	return list
}

// GetCategories 获取所有分类
func (r *AgentRegistry) GetCategories() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	categories := make(map[string]bool)
	for _, info := range r.agents {
		categories[info.Category] = true
	}
	
	var list []string
	for cat := range categories {
		list = append(list, cat)
	}
	
	return list
}

// Unregister 注销Agent
func (r *AgentRegistry) Unregister(agentType string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	delete(r.agents, agentType)
}

// registerBuiltinAgents 注册内置Agent
func (r *AgentRegistry) registerBuiltinAgents() {
	// 选题相关
	r.Register(AgentInfo{
		Type:        "hot_fetch",
		Name:        "获取热点",
		Description: "从各大平台获取实时热点数据",
		Icon:        "TrendCharts",
		Category:    "topic",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"platforms": map[string]any{
					"type":        "array",
					"description": "要获取热点的平台列表",
					"items": map[string]any{
						"type": "string",
						"enum": []string{"weibo", "zhihu", "baidu", "toutiao"},
					},
				},
				"limit": map[string]any{
					"type":        "integer",
					"description": "获取热点数量",
					"default":     20,
				},
			},
		},
		OutputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"trends": map[string]any{
					"type":        "array",
					"description": "热点列表",
				},
			},
		},
	})
	
	r.Register(AgentInfo{
		Type:        "topic_generate",
		Name:        "生成选题",
		Description: "基于热点数据AI生成优质选题",
		Icon:        "Lightbulb",
		Category:    "topic",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"trends": map[string]any{
					"type":        "array",
					"description": "热点数据",
				},
				"limit": map[string]any{
					"type":        "integer",
					"description": "生成选题数量",
					"default":     5,
				},
				"min_score": map[string]any{
					"type":        "number",
					"description": "最低评分",
					"default":     70,
				},
			},
		},
		OutputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"topics": map[string]any{
					"type":        "array",
					"description": "选题列表",
				},
			},
		},
	})
	
	r.Register(AgentInfo{
		Type:        "topic_select",
		Name:        "选择选题",
		Description: "从候选选题中选择最佳选题",
		Icon:        "Select",
		Category:    "topic",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"topics": map[string]any{
					"type":        "array",
					"description": "候选选题列表",
				},
				"strategy": map[string]any{
					"type":        "string",
					"description": "选择策略",
					"enum":        []string{"highest_score", "random", "manual"},
					"default":     "highest_score",
				},
			},
		},
		OutputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"selected_topic": map[string]any{
					"type":        "object",
					"description": "选中的选题",
				},
			},
		},
	})
	
	// 创作相关
	r.Register(AgentInfo{
		Type:        "outline_generate",
		Name:        "生成大纲",
		Description: "基于选题生成文章大纲",
		Icon:        "List",
		Category:    "content",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"title": map[string]any{
					"type":        "string",
					"description": "文章标题",
				},
				"style": map[string]any{
					"type":        "string",
					"description": "写作风格",
					"default":     "干货专业",
				},
				"sections": map[string]any{
					"type":        "integer",
					"description": "章节数量",
					"default":     5,
				},
			},
		},
		OutputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"outline": map[string]any{
					"type":        "array",
					"description": "大纲节点列表",
				},
			},
		},
	})
	
	r.Register(AgentInfo{
		Type:        "content_write",
		Name:        "创作内容",
		Description: "根据大纲创作文章内容",
		Icon:        "Edit",
		Category:    "content",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"outline": map[string]any{
					"type":        "array",
					"description": "大纲",
				},
				"title": map[string]any{
					"type":        "string",
					"description": "标题",
				},
				"word_count": map[string]any{
					"type":        "integer",
					"description": "目标字数",
					"default":     2000,
				},
			},
		},
		OutputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"article": map[string]any{
					"type":        "object",
					"description": "文章对象",
				},
			},
		},
	})
	
	r.Register(AgentInfo{
		Type:        "content_review",
		Name:        "质量审核",
		Description: "审核文章质量",
		Icon:        "View",
		Category:    "content",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"article": map[string]any{
					"type":        "object",
					"description": "文章对象",
				},
				"min_score": map[string]any{
					"type":        "number",
					"description": "最低通过分数",
					"default":     80,
				},
			},
		},
		OutputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"review": map[string]any{
					"type":        "object",
					"description": "审核结果",
				},
			},
		},
	})
	
	r.Register(AgentInfo{
		Type:        "content_rewrite",
		Name:        "重写优化",
		Description: "根据审核意见重写文章",
		Icon:        "Refresh",
		Category:    "content",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"article": map[string]any{
					"type":        "object",
					"description": "文章对象",
				},
				"issues": map[string]any{
					"type":        "array",
					"description": "需要修改的问题",
				},
			},
		},
		OutputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"article": map[string]any{
					"type":        "object",
					"description": "重写后的文章",
				},
			},
		},
	})
	
	// 图片相关
	r.Register(AgentInfo{
		Type:        "image_generate",
		Name:        "生成配图",
		Description: "AI生成文章封面图或插图",
		Icon:        "Picture",
		Category:    "image",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"title": map[string]any{
					"type":        "string",
					"description": "文章标题",
				},
				"keywords": map[string]any{
					"type":        "array",
					"description": "关键词",
				},
				"image_type": map[string]any{
					"type":        "string",
					"description": "图片类型",
					"enum":        []string{"cover", "illustration"},
					"default":     "cover",
				},
			},
		},
		OutputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"image_url": map[string]any{
					"type":        "string",
					"description": "图片URL",
				},
			},
		},
	})
	
	// 排版相关
	r.Register(AgentInfo{
		Type:        "layout_apply",
		Name:        "公众号排版",
		Description: "应用公众号排版样式",
		Icon:        "Document",
		Category:    "layout",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"article": map[string]any{
					"type":        "object",
					"description": "文章对象",
				},
				"theme": map[string]any{
					"type":        "string",
					"description": "主题样式",
					"enum":        []string{"科技蓝", "简约白", "活力橙"},
					"default":     "科技蓝",
				},
			},
		},
		OutputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"html": map[string]any{
					"type":        "string",
					"description": "排版后的HTML",
				},
			},
		},
	})
	
	// 发布相关
	r.Register(AgentInfo{
		Type:        "publish_draft",
		Name:        "发布草稿",
		Description: "发布到微信公众号草稿箱",
		Icon:        "Upload",
		Category:    "publish",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"article": map[string]any{
					"type":        "object",
					"description": "文章对象",
				},
				"html": map[string]any{
					"type":        "string",
					"description": "HTML内容",
				},
				"account": map[string]any{
					"type":        "string",
					"description": "发布账号",
					"default":     "default",
				},
			},
		},
		OutputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"media_id": map[string]any{
					"type":        "string",
					"description": "微信返回的media_id",
				},
				"published": map[string]any{
					"type":        "boolean",
					"description": "是否发布成功",
				},
			},
		},
	})
	
	// 控制流相关
	r.Register(AgentInfo{
		Type:        "condition",
		Name:        "条件判断",
		Description: "根据条件选择执行路径",
		Icon:        "Share",
		Category:    "control",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"condition": map[string]any{
					"type":        "string",
					"description": "条件表达式",
				},
				"true_next": map[string]any{
					"type":        "string",
					"description": "条件为真时的下一步",
				},
				"false_next": map[string]any{
					"type":        "string",
					"description": "条件为假时的下一步",
				},
			},
		},
		OutputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"condition": map[string]any{
					"type":        "boolean",
					"description": "条件结果",
				},
			},
		},
	})
	
	r.Register(AgentInfo{
		Type:        "delay",
		Name:        "延迟等待",
		Description: "延迟指定时间后继续",
		Icon:        "Timer",
		Category:    "control",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"duration": map[string]any{
					"type":        "integer",
					"description": "延迟时间（秒）",
					"default":     5,
				},
			},
		},
		OutputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"delayed": map[string]any{
					"type":        "boolean",
					"description": "是否延迟成功",
				},
			},
		},
	})
	
	r.Register(AgentInfo{
		Type:        "manual_review",
		Name:        "人工审核",
		Description: "暂停等待人工审核",
		Icon:        "User",
		Category:    "control",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"message": map[string]any{
					"type":        "string",
					"description": "审核提示信息",
				},
			},
		},
		OutputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"approved": map[string]any{
					"type":        "boolean",
					"description": "是否审核通过",
				},
			},
		},
	})
}
