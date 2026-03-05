//go:build darwin || windows || linux

package backend

import (
	"Content-Alchemist/backend/ai"
	"Content-Alchemist/backend/db"
	"Content-Alchemist/backend/editor"
	"Content-Alchemist/backend/models"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App 应用主结构体
type App struct {
	ctx       context.Context
	db        *db.DB
	fileMgr   *editor.FileManager
	aiService *ai.Service
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
		fmt.Printf("无法获取配置目录: %v\n", err)
	}

	a.configDir = filepath.Join(configDir, "ContentAlchemist")

	// 初始化数据库
	database, err := db.New(a.configDir)
	if err != nil {
		fmt.Printf("无法初始化数据库: %v\n", err)
	}
	a.db = database

	// 初始化 AI 服务
	a.initAIService()
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

// RenameArticleByTitle 根据新标题重命名文章文件
// 返回新的文件路径
func (a *App) RenameArticleByTitle(uuid string, newTitle string, content string) (string, error) {
	if a.db == nil {
		return "", fmt.Errorf("数据库未初始化")
	}

	// 获取文章记录
	article, err := a.db.GetArticleByUUID(uuid)
	if err != nil {
		return "", fmt.Errorf("查询文章失败: %w", err)
	}

	if article == nil {
		return "", models.FileError{
			Code:    "ARTICLE_NOT_FOUND",
			Message: fmt.Sprintf("未找到UUID对应的文章: %s", uuid),
		}
	}

	// 清理标题，作为文件名
	newFilename := a.fileMgr.SanitizeFilename(newTitle)
	if newFilename == "" {
		newFilename = "untitled"
	}

	// 获取原文件所在目录
	oldDir := filepath.Dir(article.FilePath)

	// 生成新的文件路径
	newPath := a.fileMgr.GenerateUniqueFilename(oldDir, newFilename)

	// 如果是不同路径，执行文件重命名操作
	if filepath.Clean(newPath) != filepath.Clean(article.FilePath) {
		// 写入新文件
		if err := a.fileMgr.WriteFileContent(newPath, content); err != nil {
			return "", err
		}

		// 删除旧文件
		if err := os.Remove(article.FilePath); err != nil && !os.IsNotExist(err) {
			// 删除新文件，回滚
			os.Remove(newPath)
			return "", fmt.Errorf("删除旧文件失败: %w", err)
		}

		// 更新文件路径
		article.FilePath = newPath
	}

	// 更新标题和元数据
	article.Title = newTitle
	a.fileMgr.UpdateArticleFromContent(article, content)
	if err := a.db.UpdateArticle(article); err != nil {
		return "", fmt.Errorf("更新文章记录失败: %w", err)
	}

	return article.FilePath, nil
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

// ============================================
// Wails绑定方法 - AI配置
// ============================================

// GetAIConfig 获取AI配置
func (a *App) GetAIConfig() (*models.AIConfig, error) {
	if a.db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	config := &models.AIConfig{
		BaseURL:     "https://api.deepseek.com/v1",
		Token:       "",
		Temperature: 0.7,
		Model:       "deepseek-chat",
	}

	if baseURL, err := a.db.GetSetting("ai_base_url"); err == nil && baseURL != "" {
		config.BaseURL = baseURL
	}

	if token, err := a.db.GetSetting("ai_token"); err == nil {
		config.Token = token
	}

	if temp, err := a.db.GetSetting("ai_temperature"); err == nil && temp != "" {
		if t, err := strconv.ParseFloat(temp, 64); err == nil {
			config.Temperature = t
		}
	}

	if model, err := a.db.GetSetting("ai_model"); err == nil && model != "" {
		config.Model = model
	}

	return config, nil
}

// SaveAIConfig 保存AI配置
func (a *App) SaveAIConfig(config *models.AIConfig) error {
	if a.db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	if err := a.db.SetSetting("ai_base_url", config.BaseURL); err != nil {
		return fmt.Errorf("保存base_url失败: %w", err)
	}

	if err := a.db.SetSetting("ai_token", config.Token); err != nil {
		return fmt.Errorf("保存token失败: %w", err)
	}

	tempStr := strconv.FormatFloat(config.Temperature, 'f', 2, 64)
	if err := a.db.SetSetting("ai_temperature", tempStr); err != nil {
		return fmt.Errorf("保存temperature失败: %w", err)
	}

	if err := a.db.SetSetting("ai_model", config.Model); err != nil {
		return fmt.Errorf("保存model失败: %w", err)
	}

	// 更新 AI 服务配置
	if a.aiService != nil {
		a.aiService.UpdateConfig(config)
	}

	return nil
}

// GetPositioning 获取公众号定位配置
func (a *App) GetPositioning() (string, error) {
	if a.db == nil {
		return "", fmt.Errorf("数据库未初始化")
	}

	positioning, err := a.db.GetSetting("positioning")
	if err != nil {
		return "", fmt.Errorf("获取公众号定位失败: %w", err)
	}

	return positioning, nil
}

// SavePositioning 保存公众号定位配置
func (a *App) SavePositioning(positioning string) error {
	if a.db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	if err := a.db.SetSetting("positioning", positioning); err != nil {
		return fmt.Errorf("保存公众号定位失败: %w", err)
	}

	return nil
}

// ============================================
// AI 服务相关
// ============================================

// initAIService 初始化 AI 服务
func (a *App) initAIService() {
	config, err := a.GetAIConfig()
	if err != nil {
		fmt.Printf("获取AI配置失败: %v\n", err)
		return
	}
	a.aiService = ai.NewService(config)
}

// GenerateOutline 根据标题、写作要求和公众号定位生成大纲和候选标题
// title: 文章标题
// requirements: 写作要求
// positioning: 公众号定位
// 返回: 包含3个候选标题和大纲的结果
func (a *App) GenerateOutline(title string, requirements string, positioning string) (*models.GenerateOutlineResult, error) {
	if a.aiService == nil {
		// 重新初始化 AI 服务
		a.initAIService()
	}

	if title == "" {
		return nil, fmt.Errorf("标题不能为空")
	}

	result, err := a.aiService.GenerateOutlineWithTitles(title, requirements, positioning)
	if err != nil {
		return nil, err
	}

	return &models.GenerateOutlineResult{
		Titles:  result.Titles,
		Outline: result.Outline,
	}, nil
}

// GenerateArticle 根据大纲生成文章
// title: 文章标题
// outline: 文章大纲
// requirements: 写作要求
// 返回: 生成的文章内容
func (a *App) GenerateArticle(title string, outline string, requirements string) (string, error) {
	if a.aiService == nil {
		a.initAIService()
	}

	if title == "" {
		return "", fmt.Errorf("标题不能为空")
	}

	if outline == "" {
		return "", fmt.Errorf("大纲不能为空，请先生成大纲")
	}

	article, err := a.aiService.GenerateArticleFromOutline(title, outline, requirements)
	if err != nil {
		return "", err
	}

	return article, nil
}

// SaveArticleWithSmartNaming 智能保存文章
// 如果文章未保存到本地文件，则根据标题自动生成文件名并弹出保存对话框
// 如果已保存，则直接保存
func (a *App) SaveArticleWithSmartNaming(uuid string, title string, content string) (*models.Article, error) {
	if a.db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	// 获取文章记录
	article, err := a.db.GetArticleByUUID(uuid)
	if err != nil {
		return nil, fmt.Errorf("查询文章失败: %w", err)
	}

	// 如果文章不存在，或者文件路径为空，或者是临时文件，需要选择保存位置
	needsNewFile := article == nil || article.FilePath == "" || !a.fileMgr.CheckFileExistsSync(article.FilePath)

	if needsNewFile {
		// 根据标题生成文件名
		filename := a.fileMgr.SanitizeFilename(title)
		if filename == "" {
			filename = "untitled"
		}
		filename += ".md"

		// 打开保存对话框
		filePath, err := a.SaveFileDialog(filename)
		if err != nil {
			return nil, err
		}
		if filePath == "" {
			// 用户取消
			return nil, nil
		}

		// 保存到新文件
		if article == nil {
			// 创建新文章
			article = a.fileMgr.CreateNewArticleData(filePath, content)
			article.Title = title
			if err := a.fileMgr.WriteFileContent(filePath, content); err != nil {
				return nil, err
			}
			if err := a.db.CreateArticle(article); err != nil {
				return nil, fmt.Errorf("创建文章记录失败: %w", err)
			}
		} else {
			// 更新现有文章
			article.FilePath = filePath
			article.Title = title
			a.fileMgr.UpdateArticleFromContent(article, content)
			if err := a.fileMgr.WriteFileContent(filePath, content); err != nil {
				return nil, err
			}
			if err := a.db.UpdateArticle(article); err != nil {
				return nil, fmt.Errorf("更新文章记录失败: %w", err)
			}
		}

		return article, nil
	}

	// 文章已存在，直接保存
	if err := a.SaveArticle(uuid, content); err != nil {
		return nil, err
	}

	// 更新标题
	if article.Title != title {
		article.Title = title
		a.fileMgr.UpdateArticleFromContent(article, content)
		if err := a.db.UpdateArticle(article); err != nil {
			return nil, fmt.Errorf("更新文章标题失败: %w", err)
		}
	}

	return article, nil
}
