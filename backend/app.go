//go:build darwin || windows || linux

package backend

import (
	"Content-Alchemist/backend/db"
	"Content-Alchemist/backend/editor"
	"Content-Alchemist/backend/models"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App 应用主结构体
type App struct {
	ctx       context.Context
	db        *db.DB
	fileMgr   *editor.FileManager
	configDir string
}

// NewApp 创建新应用实例
func NewApp() *App {
	return &App{
		fileMgr: editor.NewFileManager(),
	}
}

// Startup 在应用启动时调用（Wails生命周期）
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// 获取配置目录
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("无法获取配置目录: %w", err)
	}

	a.configDir = filepath.Join(configDir, "ContentAlchemist")

	// 初始化数据库
	database, err := db.New(a.configDir)
	if err != nil {
		fmt.Printf("无法初始化数据库: %w", err)
	}
	a.db = database
}

// Shutdown 在应用关闭时调用（Wails生命周期）
func (a *App) Shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// ============================================
// Wails绑定方法 - 文件对话框
// ============================================

// OpenFileDialog 打开文件选择对话框
// 返回: 选中的文件路径，如果用户取消则返回空字符串
func (a *App) OpenFileDialog() (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("应用未初始化")
	}

	// 获取默认目录
	defaultDir, err := a.fileMgr.GetDefaultSaveDirectory()
	if err != nil {
		defaultDir = "."
	}

	// 打开文件对话框
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "打开 Markdown 文件",
		DefaultDirectory: defaultDir,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Markdown 文件",
				Pattern:     "*.md;*.markdown;*.mdown",
			},
			{
				DisplayName: "所有文件",
				Pattern:     "*",
			},
		},
		ShowHiddenFiles: false,
	})

	if err != nil {
		return "", fmt.Errorf("打开对话框失败: %w", err)
	}

	// 用户取消
	if selection == "" {
		return "", nil
	}

	// 转换为绝对路径
	absPath, err := filepath.Abs(selection)
	if err != nil {
		return "", fmt.Errorf("无法获取绝对路径: %w", err)
	}

	return absPath, nil
}

// SaveFileDialog 打开保存文件对话框
// 返回: 选中的文件路径，如果用户取消则返回空字符串
func (a *App) SaveFileDialog(defaultFilename string) (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("应用未初始化")
	}

	// 获取默认目录
	defaultDir, err := a.fileMgr.GetDefaultSaveDirectory()
	if err != nil {
		defaultDir = "."
	}

	// 如果没有提供默认文件名，生成一个
	if defaultFilename == "" {
		defaultFilename = "untitled.md"
	}

	// 确保有.md扩展名
	if filepath.Ext(defaultFilename) == "" {
		defaultFilename += ".md"
	}

	// 打开保存对话框
	selection, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            "保存 Markdown 文件",
		DefaultDirectory: defaultDir,
		DefaultFilename:  defaultFilename,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Markdown 文件",
				Pattern:     "*.md",
			},
			{
				DisplayName: "所有文件",
				Pattern:     "*",
			},
		},
		ShowHiddenFiles: false,
	})

	if err != nil {
		return "", fmt.Errorf("打开保存对话框失败: %w", err)
	}

	// 用户取消
	if selection == "" {
		return "", nil
	}

	// 转换为绝对路径并确保.md扩展名
	absPath, err := filepath.Abs(selection)
	if err != nil {
		return "", fmt.Errorf("无法获取绝对路径: %w", err)
	}

	if filepath.Ext(absPath) == "" {
		absPath += ".md"
	}

	return absPath, nil
}

// OpenDirectoryDialog 打开目录选择对话框
// 返回: 选中的目录路径，如果用户取消则返回空字符串
func (a *App) OpenDirectoryDialog() (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("应用未初始化")
	}

	defaultDir, err := a.fileMgr.GetDefaultSaveDirectory()
	if err != nil {
		defaultDir = "."
	}

	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "选择文件夹",
		DefaultDirectory: defaultDir,
	})

	if err != nil {
		return "", fmt.Errorf("打开目录对话框失败: %w", err)
	}

	if selection == "" {
		return "", nil
	}

	absPath, err := filepath.Abs(selection)
	if err != nil {
		return "", fmt.Errorf("无法获取绝对路径: %w", err)
	}

	return absPath, nil
}

// ============================================
// Wails绑定方法 - 文件操作
// ============================================

// ReadArticle 读取文章
// 如果文章不在数据库中，会创建新记录
// 返回: 文章元数据和内容
func (a *App) ReadArticle(filePath string) (*models.ReadArticleResponse, error) {
	if a.db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	// 清理路径
	filePath = filepath.Clean(filePath)

	// 读取文件内容
	content, err := a.fileMgr.ReadFileContent(filePath)
	if err != nil {
		return nil, err
	}

	// 检查数据库中是否已有此文章
	article, err := a.db.GetArticleByPath(filePath)
	if err != nil {
		return nil, fmt.Errorf("查询数据库失败: %w", err)
	}

	if article == nil {
		// 创建新文章记录
		article = a.fileMgr.CreateNewArticleData(filePath, content)
		if err := a.db.CreateArticle(article); err != nil {
			return nil, fmt.Errorf("创建文章记录失败: %w", err)
		}
	} else {
		// 更新元数据（标题、摘要、字数可能已变化）
		a.fileMgr.UpdateArticleFromContent(article, content)
		if err := a.db.UpdateArticle(article); err != nil {
			return nil, fmt.Errorf("更新文章记录失败: %w", err)
		}
	}

	// 更新最后打开时间
	if err := a.db.UpdateLastOpened(article.UUID); err != nil {
		// 非致命错误，记录但继续
		runtime.LogWarningf(a.ctx, "更新最后打开时间失败: %v", err)
	}

	return &models.ReadArticleResponse{
		Article: article,
		Content: content,
	}, nil
}

// SaveArticle 保存文章
// 保存内容到文件，并更新数据库元数据
func (a *App) SaveArticle(uuid string, content string) error {
	if a.db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 获取文章记录
	article, err := a.db.GetArticleByUUID(uuid)
	if err != nil {
		return fmt.Errorf("查询文章失败: %w", err)
	}

	if article == nil {
		return models.FileError{
			Code:    "ARTICLE_NOT_FOUND",
			Message: fmt.Sprintf("未找到UUID对应的文章: %s", uuid),
		}
	}

	// 检查文件是否被外部修改
	exists, err := a.fileMgr.CheckFileExists(article.FilePath)
	if err != nil {
		return fmt.Errorf("检查文件状态失败: %w", err)
	}

	if exists {
		// 文件存在，检查是否被外部修改
		modified, _, err := a.fileMgr.CheckFileModified(article.FilePath, article.UpdatedAt)
		if err != nil {
			// 如果是文件不存在错误，忽略（可能被删除）
			if _, ok := err.(models.FileError); !ok {
				return fmt.Errorf("检查文件修改状态失败: %w", err)
			}
		} else if modified {
			return models.FileError{
				Code:    "FILE_MODIFIED_EXTERNALLY",
				Message: "文件已被外部程序修改，请先刷新或另存为",
			}
		}
	}

	// 写入文件
	if err := a.fileMgr.WriteFileContent(article.FilePath, content); err != nil {
		return err
	}

	// 更新元数据
	a.fileMgr.UpdateArticleFromContent(article, content)
	if err := a.db.UpdateArticle(article); err != nil {
		return fmt.Errorf("更新数据库记录失败: %w", err)
	}

	return nil
}

// SaveArticleAs 另存为
// 将文章保存到新路径，创建新的数据库记录
func (a *App) SaveArticleAs(uuid string, newPath string, content string) (*models.Article, error) {
	if a.db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	newPath = filepath.Clean(newPath)

	// 检查新路径是否已存在其他文章
	existingArticle, err := a.db.GetArticleByPath(newPath)
	if err != nil {
		return nil, fmt.Errorf("查询数据库失败: %w", err)
	}

	if existingArticle != nil && existingArticle.UUID != uuid {
		return nil, models.FileError{
			Code:    "FILE_EXISTS",
			Message: "目标文件已存在，且关联其他文章",
		}
	}

	// 写入新文件
	if err := a.fileMgr.WriteFileContent(newPath, content); err != nil {
		return nil, err
	}

	// 如果已存在相同路径的文章，更新它
	if existingArticle != nil {
		a.fileMgr.UpdateArticleFromContent(existingArticle, content)
		if err := a.db.UpdateArticle(existingArticle); err != nil {
			return nil, fmt.Errorf("更新文章记录失败: %w", err)
		}
		return existingArticle, nil
	}

	// 创建新文章记录
	article := a.fileMgr.CreateNewArticleData(newPath, content)
	if err := a.db.CreateArticle(article); err != nil {
		return nil, fmt.Errorf("创建文章记录失败: %w", err)
	}

	return article, nil
}

// CreateNewArticle 创建新文章
// 打开保存对话框，创建空文件
func (a *App) CreateNewArticle() (*models.Article, error) {
	if a.db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	// 打开保存对话框
	filePath, err := a.SaveFileDialog("untitled.md")
	if err != nil {
		return nil, err
	}

	// 用户取消
	if filePath == "" {
		return nil, nil
	}

	// 检查文件是否已存在
	exists, err := a.fileMgr.CheckFileExists(filePath)
	if err != nil {
		return nil, fmt.Errorf("检查文件失败: %w", err)
	}

	if exists {
		return nil, models.FileError{
			Code:    "FILE_EXISTS",
			Message: "文件已存在，请使用打开功能",
		}
	}

	// 创建空文件
	content := ""
	if err := a.fileMgr.WriteFileContent(filePath, content); err != nil {
		return nil, err
	}

	// 创建数据库记录
	article := a.fileMgr.CreateNewArticleData(filePath, content)
	if err := a.db.CreateArticle(article); err != nil {
		return nil, fmt.Errorf("创建文章记录失败: %w", err)
	}

	return article, nil
}

// CreateNewArticleWithContent 创建新文章并指定初始内容
func (a *App) CreateNewArticleWithContent(defaultFilename string, content string) (*models.Article, error) {
	if a.db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	// 打开保存对话框
	filePath, err := a.SaveFileDialog(defaultFilename)
	if err != nil {
		return nil, err
	}

	if filePath == "" {
		return nil, nil
	}

	// 如果文件已存在，检查是否关联同一文章
	exists, err := a.fileMgr.CheckFileExists(filePath)
	if err != nil {
		return nil, fmt.Errorf("检查文件失败: %w", err)
	}

	if exists {
		existingArticle, err := a.db.GetArticleByPath(filePath)
		if err != nil {
			return nil, fmt.Errorf("查询数据库失败: %w", err)
		}

		if existingArticle != nil {
			// 更新现有文章
			a.fileMgr.UpdateArticleFromContent(existingArticle, content)
			if err := a.fileMgr.WriteFileContent(filePath, content); err != nil {
				return nil, err
			}
			if err := a.db.UpdateArticle(existingArticle); err != nil {
				return nil, fmt.Errorf("更新文章记录失败: %w", err)
			}
			return existingArticle, nil
		}
	}

	// 创建新文件和记录
	if err := a.fileMgr.WriteFileContent(filePath, content); err != nil {
		return nil, err
	}

	article := a.fileMgr.CreateNewArticleData(filePath, content)
	if err := a.db.CreateArticle(article); err != nil {
		return nil, fmt.Errorf("创建文章记录失败: %w", err)
	}

	return article, nil
}

// ============================================
// Wails绑定方法 - 文章元数据
// ============================================

// GetRecentArticles 获取最近打开的文章列表
func (a *App) GetRecentArticles(limit int) ([]*models.Article, error) {
	if a.db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return a.db.GetRecentArticles(limit)
}

// UpdateArticleMeta 更新文章元数据（不触碰文件内容）
func (a *App) UpdateArticleMeta(uuid string, tags []string) error {
	if a.db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	article, err := a.db.GetArticleByUUID(uuid)
	if err != nil {
		return fmt.Errorf("查询文章失败: %w", err)
	}

	if article == nil {
		return models.FileError{
			Code:    "ARTICLE_NOT_FOUND",
			Message: fmt.Sprintf("未找到UUID对应的文章: %s", uuid),
		}
	}

	article.Tags = tags
	return a.db.UpdateArticle(article)
}

// GetArticleByUUID 根据UUID获取文章元数据
func (a *App) GetArticleByUUID(uuid string) (*models.Article, error) {
	if a.db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	article, err := a.db.GetArticleByUUID(uuid)
	if err != nil {
		return nil, fmt.Errorf("查询文章失败: %w", err)
	}

	if article == nil {
		return nil, models.FileError{
			Code:    "ARTICLE_NOT_FOUND",
			Message: fmt.Sprintf("未找到UUID对应的文章: %s", uuid),
		}
	}

	return article, nil
}

// DeleteArticle 删除文章记录（不删除文件）
func (a *App) DeleteArticle(uuid string) error {
	if a.db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	return a.db.DeleteArticle(uuid)
}

// DeleteArticleAndFile 删除文章记录和文件
func (a *App) DeleteArticleAndFile(uuid string) error {
	if a.db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	article, err := a.db.GetArticleByUUID(uuid)
	if err != nil {
		return fmt.Errorf("查询文章失败: %w", err)
	}

	if article == nil {
		return models.FileError{
			Code:    "ARTICLE_NOT_FOUND",
			Message: fmt.Sprintf("未找到UUID对应的文章: %s", uuid),
		}
	}

	// 删除文件
	if err := os.Remove(article.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	// 删除数据库记录
	return a.db.DeleteArticle(uuid)
}

// CheckFileExists 检查文件是否存在
func (a *App) CheckFileExists(filePath string) (bool, error) {
	return a.fileMgr.CheckFileExists(filePath)
}

// GetFileInfo 获取文件信息
func (a *App) GetFileInfo(filePath string) (map[string]interface{}, error) {
	filePath = filepath.Clean(filePath)

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]interface{}{
				"exists": false,
			}, nil
		}
		return nil, err
	}

	return map[string]interface{}{
		"exists":  true,
		"isDir":   info.IsDir(),
		"size":    info.Size(),
		"modTime": info.ModTime().Unix(),
		"path":    filePath,
	}, nil
}
