//go:build darwin || windows || linux

package editor

import (
	"Content-Alchemist/backend/models"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// FileManager 文件管理器
type FileManager struct {
	// 可注入依赖，如logger等
}

// NewFileManager 创建文件管理器
func NewFileManager() *FileManager {
	return &FileManager{}
}

// ReadFileContent 读取文件内容
func (fm *FileManager) ReadFileContent(filePath string) (string, error) {
	// 清理路径
	filePath = filepath.Clean(filePath)

	// 检查文件是否存在
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", models.FileError{
				Code:    "FILE_NOT_FOUND",
				Message: fmt.Sprintf("文件不存在: %s", filePath),
			}
		}
		if os.IsPermission(err) {
			return "", models.FileError{
				Code:    "PERMISSION_DENIED",
				Message: fmt.Sprintf("无权限访问文件: %s", filePath),
			}
		}
		return "", fmt.Errorf("无法访问文件: %w", err)
	}

	// 检查是否为目录
	if info.IsDir() {
		return "", models.FileError{
			Code:    "IS_DIRECTORY",
			Message: fmt.Sprintf("路径是目录而非文件: %s", filePath),
		}
	}

	// 检查文件大小 (限制100MB)
	const maxSize = 100 * 1024 * 1024
	if info.Size() > maxSize {
		return "", models.FileError{
			Code:    "FILE_TOO_LARGE",
			Message: fmt.Sprintf("文件超过100MB限制: %s (%d MB)", filePath, info.Size()/(1024*1024)),
		}
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsPermission(err) {
			return "", models.FileError{
				Code:    "PERMISSION_DENIED",
				Message: fmt.Sprintf("无权限读取文件: %s", filePath),
			}
		}
		return "", fmt.Errorf("读取文件失败: %w", err)
	}

	return string(content), nil
}

// WriteFileContent 写入文件内容
func (fm *FileManager) WriteFileContent(filePath string, content string) error {
	// 清理路径
	filePath = filepath.Clean(filePath)

	// 检查文件是否存在
	info, err := os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("无法访问文件: %w", err)
	}

	// 如果存在，检查是否为目录
	if err == nil && info.IsDir() {
		return models.FileError{
			Code:    "IS_DIRECTORY",
			Message: fmt.Sprintf("路径是目录而非文件: %s", filePath),
		}
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		if os.IsPermission(err) {
			return models.FileError{
				Code:    "PERMISSION_DENIED",
				Message: fmt.Sprintf("无权限创建目录: %s", dir),
			}
		}
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		if os.IsPermission(err) {
			return models.FileError{
				Code:    "PERMISSION_DENIED",
				Message: fmt.Sprintf("无权限写入文件: %s", filePath),
			}
		}
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// CheckFileExists 检查文件是否存在
func (fm *FileManager) CheckFileExists(filePath string) (bool, error) {
	filePath = filepath.Clean(filePath)
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !info.IsDir(), nil
}

// CheckFileExistsSync 同步检查文件是否存在（无错误返回）
func (fm *FileManager) CheckFileExistsSync(filePath string) bool {
	exists, _ := fm.CheckFileExists(filePath)
	return exists
}

// CheckFileModified 检查文件是否被外部修改
// 返回：是否被修改, 修改时间, 错误
func (fm *FileManager) CheckFileModified(filePath string, lastModified time.Time) (bool, time.Time, error) {
	filePath = filepath.Clean(filePath)
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return true, time.Time{}, models.FileError{
				Code:    "FILE_NOT_FOUND",
				Message: "文件已被删除",
			}
		}
		return false, time.Time{}, err
	}

	currentModTime := info.ModTime()
	return currentModTime.After(lastModified), currentModTime, nil
}

// ExtractTitleFromContent 从内容中提取标题
func (fm *FileManager) ExtractTitleFromContent(content string, fallback string) string {
	// 匹配H1标题: # Title 或 # Title #
	re := regexp.MustCompile(`(?m)^#\s+(.+?)(?:\s+#)?$`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		title := strings.TrimSpace(matches[1])
		if title != "" {
			return title
		}
	}

	// 如果没有H1，使用前50个字符作为标题
	content = strings.TrimSpace(content)
	if content == "" {
		return fallback
	}

	// 移除Markdown标记
	content = regexp.MustCompile(`[#*_~\[\](){}|`+"`"+`]`).ReplaceAllString(content, "")
	content = strings.TrimSpace(content)

	if len(content) <= 50 {
		return content
	}
	return content[:50] + "..."
}

// ExtractSummaryFromContent 从内容中提取摘要
func (fm *FileManager) ExtractSummaryFromContent(content string) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return ""
	}

	// 移除代码块
	content = regexp.MustCompile("`"+`{3}[\s\S]*?`+"`"+`{3}`).ReplaceAllString(content, "")
	content = regexp.MustCompile("`"+`[^`+"`"+"`]*"+"`").ReplaceAllString(content, "")

	// 移除Markdown标记
	content = regexp.MustCompile(`[#*_~\[\](){}|`+"`"+`]`).ReplaceAllString(content, "")
	content = strings.TrimSpace(content)

	// 移除HTML标签
	content = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(content, "")

	// 提取前200个字符作为摘要
	if len(content) <= 200 {
		return content
	}
	return content[:200] + "..."
}

// CountWords 统计字数
func (fm *FileManager) CountWords(content string) int {
	// 移除代码块不计入字数
	content = regexp.MustCompile("`"+`{3}[\s\S]*?`+"`"+`{3}`).ReplaceAllString(content, "")
	content = regexp.MustCompile("`"+`[^`+"`"+"`]*"+"`").ReplaceAllString(content, "")

	// 中文：每个汉字算一个字
	// 英文：按空格分词
	chineseCount := 0
	englishWords := 0

	for _, r := range content {
		if r >= '\u4e00' && r <= '\u9fff' {
			chineseCount++
		}
	}

	// 移除中文后计算英文词数
	// Go regexp 不支持 \u 转义，使用原始 Unicode 字符范围
	nonChinese := regexp.MustCompile(`[一-龥]`).ReplaceAllString(content, "")
	words := strings.Fields(nonChinese)
	englishWords = len(words)

	return chineseCount + englishWords
}

// CreateNewArticleData 创建新文章数据
func (fm *FileManager) CreateNewArticleData(filePath string, content string) *models.Article {
	now := time.Now()
	fallbackTitle := filepath.Base(filePath)
	if ext := filepath.Ext(fallbackTitle); ext != "" {
		fallbackTitle = fallbackTitle[:len(fallbackTitle)-len(ext)]
	}

	title := fm.ExtractTitleFromContent(content, fallbackTitle)
	summary := fm.ExtractSummaryFromContent(content)
	wordCount := fm.CountWords(content)

	return &models.Article{
		UUID:         uuid.New().String(),
		FilePath:     filePath,
		Title:        title,
		Summary:      summary,
		Tags:         models.StringSlice{},
		WordCount:    wordCount,
		CreatedAt:    now,
		UpdatedAt:    now,
		LastOpenedAt: now,
	}
}

// UpdateArticleFromContent 根据内容更新文章元数据
func (fm *FileManager) UpdateArticleFromContent(article *models.Article, content string) {
	fallbackTitle := filepath.Base(article.FilePath)
	if ext := filepath.Ext(fallbackTitle); ext != "" {
		fallbackTitle = fallbackTitle[:len(fallbackTitle)-len(ext)]
	}

	article.Title = fm.ExtractTitleFromContent(content, fallbackTitle)
	article.Summary = fm.ExtractSummaryFromContent(content)
	article.WordCount = fm.CountWords(content)
	article.UpdatedAt = time.Now()
}

// GetDefaultSaveDirectory 获取默认保存目录
func (fm *FileManager) GetDefaultSaveDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("无法获取用户主目录: %w", err)
	}

	// 使用 Documents/ContentAlchemist
	docDir := filepath.Join(homeDir, "Documents", "ContentAlchemist")
	if err := os.MkdirAll(docDir, 0755); err != nil {
		return "", fmt.Errorf("无法创建文档目录: %w", err)
	}

	return docDir, nil
}

// GenerateUniqueFilename 生成唯一的文件名
func (fm *FileManager) GenerateUniqueFilename(dir string, baseName string) string {
	// 清理文件名
	baseName = fm.SanitizeFilename(baseName)
	if baseName == "" {
		baseName = "untitled"
	}

	filename := baseName + ".md"
	fullPath := filepath.Join(dir, filename)

	// 如果文件已存在，添加数字后缀
	counter := 1
	for {
		exists, _ := fm.CheckFileExists(fullPath)
		if !exists {
			break
		}
		filename = fmt.Sprintf("%s_%d.md", baseName, counter)
		fullPath = filepath.Join(dir, filename)
		counter++
	}

	return fullPath
}

// SanitizeFilename 清理文件名，移除非法字符
func (fm *FileManager) SanitizeFilename(name string) string {
	// 移除Windows和Unix的非法字符
	illegal := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
	name = illegal.ReplaceAllString(name, "")

	// 移除首尾空格和点
	name = strings.TrimSpace(name)
	name = strings.Trim(name, ".")

	// 限制长度
	if len(name) > 200 {
		name = name[:200]
	}

	return name
}
