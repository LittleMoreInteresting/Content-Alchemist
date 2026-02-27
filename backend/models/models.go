//go:build darwin || windows || linux

package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// StringSlice 用于在SQLite中存储JSON数组
type StringSlice []string

// Value 实现driver.Valuer接口
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

// Scan 实现sql.Scanner接口
func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = StringSlice{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan type %T into StringSlice", value)
	}

	return json.Unmarshal(bytes, s)
}

// Article 文章元数据结构
type Article struct {
	ID           int64       `json:"id"`
	UUID         string      `json:"uuid"`
	FilePath     string      `json:"filePath"`
	Title        string      `json:"title"`
	Summary      string      `json:"summary"`
	Tags         StringSlice `json:"tags"`
	WordCount    int         `json:"wordCount"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	LastOpenedAt time.Time   `json:"lastOpenedAt"`
}

// ArticleWithContent 包含文章内容的文章结构
type ArticleWithContent struct {
	Article Article `json:"article"`
	Content string  `json:"content"`
}

// EditorSettings 编辑器设置
type EditorSettings struct {
	DeepseekAPIKey string `json:"deepseekApiKey"`
	EditorTheme    string `json:"editorTheme"`
	FontSize       int    `json:"fontSize"`
	AIModel        string `json:"aiModel"`
}

// AIRequest AI请求结构
type AIRequest struct {
	SelectedText string `json:"selectedText"`
	Action       string `json:"action"` // "rewrite", "polish", "continue", "shorter", "code"
	Context      string `json:"context"` // 全文上下文（可选）
}

// AIResponse AI响应结构
type AIResponse struct {
	Result     string `json:"result"`
	Error      string `json:"error,omitempty"`
	TokensUsed int    `json:"tokensUsed,omitempty"`
}

// FileDialogOptions 文件对话框选项
type FileDialogOptions struct {
	DefaultDirectory string   `json:"defaultDirectory,omitempty"`
	DefaultFilename  string   `json:"defaultFilename,omitempty"`
	Title            string   `json:"title,omitempty"`
	Filters          []Filter `json:"filters,omitempty"`
}

// Filter 文件过滤器
type Filter struct {
	DisplayName string   `json:"displayName"`
	Pattern     string   `json:"pattern"` // 例如 "*.md;*.markdown"
}

// FileError 文件操作错误
type FileError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e FileError) Error() string {
	return e.Message
}

// Result 通用结果包装器
type Result[T any] struct {
	Data    T      `json:"data"`
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
}

// NewSuccessResult 创建成功结果
func NewSuccessResult[T any](data T) Result[T] {
	return Result[T]{
		Data:    data,
		Success: true,
	}
}

// NewErrorResult 创建错误结果
func NewErrorResult[T any](err string) Result[T] {
	var zero T
	return Result[T]{
		Data:    zero,
		Error:   err,
		Success: false,
	}
}

// ReadArticleResponse ReadArticle 方法的响应结构
type ReadArticleResponse struct {
	Article *Article `json:"article"`
	Content string   `json:"content"`
}
