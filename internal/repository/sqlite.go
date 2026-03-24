package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"content-alchemist/internal/model"
	"content-alchemist/internal/utils"

	_ "modernc.org/sqlite"
)

// DB SQLite数据库管理器
type DB struct {
	conn   *sql.DB
	crypto *utils.Crypto
}

// ErrKeyNotFound 密钥未找到错误
var ErrKeyNotFound = utils.ErrKeyNotFound

// NewDB 创建数据库连接
func NewDB() (*DB, error) {
	// 首先检查是否有加密密钥
	if !utils.HasKey() {
		return nil, ErrKeyNotFound
	}

	dataDir, err := dbDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dataDir, "data.db")
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	crypto, err := utils.NewCrypto()
	if err != nil {
		return nil, err
	}

	db := &DB{conn: conn, crypto: crypto}

	// 初始化表结构
	if err := db.initTables(); err != nil {
		return nil, err
	}

	return db, nil
}
func dbDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dataDir := filepath.Join(homeDir, ".content-alchemist")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", err
	}
	return dataDir, nil
}

// InitWithNewKey 生成新密钥并初始化数据库
func InitWithNewKey() error {
	// 生成并存储新密钥
	if err := utils.GenerateAndStoreKey(); err != nil {
		return fmt.Errorf("generate and store key failed: %w", err)
	}
	// 清理数据库
	dataDir, err := dbDir()
	if err != nil {
		return err
	}
	return os.RemoveAll(dataDir)
}

// HasKey 检查系统中是否存在加密密钥
func HasKey() bool {
	return utils.HasKey()
}

// Close 关闭数据库连接
func (d *DB) Close() error {
	return d.conn.Close()
}

// initTables 初始化数据库表
func (d *DB) initTables() error {
	queries := []string{
		// 原有表
		`CREATE TABLE IF NOT EXISTS config (
			id TEXT PRIMARY KEY,
			api_base_url TEXT DEFAULT 'https://api.deepseek.com',
			api_key TEXT,
			model TEXT DEFAULT 'deepseek-chat',
			temperature REAL DEFAULT 0.7,
			style_tags TEXT,
			audience TEXT,
			persona TEXT DEFAULT '我',
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS articles (
			id TEXT PRIMARY KEY,
			title TEXT,
			content TEXT,
			outline TEXT,
			status TEXT DEFAULT 'draft',
			source_type TEXT DEFAULT 'manual',
			workflow_run_id TEXT,
			topic_id TEXT,
			quality_score REAL DEFAULT 0,
			read_time INTEGER DEFAULT 0,
			word_count INTEGER DEFAULT 0,
			publish_task_id TEXT,
			published_at DATETIME,
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS materials (
			id TEXT PRIMARY KEY,
			material_type TEXT,
			title TEXT,
			content TEXT,
			tags TEXT,
			source TEXT,
			created_at DATETIME,
			usage_count INTEGER DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS versions (
			id TEXT PRIMARY KEY,
			article_id TEXT,
			content TEXT,
			snapshot TEXT,
			created_at DATETIME,
			version_type TEXT
		)`,
		// 工作流相关表
		`CREATE TABLE IF NOT EXISTS workflows (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			trigger_config TEXT,
			steps TEXT,
			auto_publish BOOLEAN DEFAULT 0,
			need_review BOOLEAN DEFAULT 1,
			status TEXT DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS workflow_runs (
			id TEXT PRIMARY KEY,
			workflow_id TEXT NOT NULL,
			status TEXT DEFAULT 'pending',
			current_step TEXT,
			input TEXT,
			output TEXT,
			steps TEXT,
			started_at DATETIME,
			completed_at DATETIME,
			error TEXT,
			FOREIGN KEY (workflow_id) REFERENCES workflows(id)
		)`,
		// 选题相关表
		`CREATE TABLE IF NOT EXISTS topics (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			category TEXT,
			source TEXT,
			source_url TEXT,
			score REAL DEFAULT 0,
			hot_score REAL DEFAULT 0,
			comp_score REAL DEFAULT 0,
			fit_score REAL DEFAULT 0,
			keywords TEXT,
			summary TEXT,
			reference_urls TEXT,
			angles TEXT,
			reason TEXT,
			status TEXT DEFAULT 'pending',
			workflow_run_id TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS hot_trends (
			id TEXT PRIMARY KEY,
			platform TEXT NOT NULL,
			title TEXT NOT NULL,
			url TEXT,
			hot_rank INTEGER,
			hot_value REAL,
			category TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		// 发布相关表
		`CREATE TABLE IF NOT EXISTS publish_accounts (
			id TEXT PRIMARY KEY,
			platform TEXT NOT NULL,
			name TEXT NOT NULL,
			account_id TEXT,
			account_type TEXT,
			app_id TEXT,
			app_secret TEXT,
			access_token TEXT,
			token_expiry DATETIME,
			refresh_token TEXT,
			config TEXT,
			status TEXT DEFAULT 'active',
			last_used_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS publish_tasks (
			id TEXT PRIMARY KEY,
			article_id TEXT NOT NULL,
			account_id TEXT NOT NULL,
			title TEXT,
			content TEXT,
			cover_image TEXT,
			summary TEXT,
			author TEXT,
			tags TEXT,
			category TEXT,
			original_url TEXT,
			schedule_type TEXT DEFAULT 'immediate',
			schedule_at DATETIME,
			status TEXT DEFAULT 'pending',
			platform_id TEXT,
			platform_url TEXT,
			error TEXT,
			published_at DATETIME,
			retry_count INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (account_id) REFERENCES publish_accounts(id)
		)`,
		// 创建索引
		`CREATE INDEX IF NOT EXISTS idx_workflow_runs_workflow_id ON workflow_runs(workflow_id)`,
		`CREATE INDEX IF NOT EXISTS idx_workflow_runs_status ON workflow_runs(status)`,
		`CREATE INDEX IF NOT EXISTS idx_topics_status ON topics(status)`,
		`CREATE INDEX IF NOT EXISTS idx_topics_source ON topics(source)`,
		`CREATE INDEX IF NOT EXISTS idx_hot_trends_platform ON hot_trends(platform)`,
		`CREATE INDEX IF NOT EXISTS idx_publish_tasks_status ON publish_tasks(status)`,
		`CREATE INDEX IF NOT EXISTS idx_publish_tasks_schedule ON publish_tasks(schedule_at)`,
	}

	for _, query := range queries {
		if _, err := d.conn.Exec(query); err != nil {
			return fmt.Errorf("init table failed: %w", err)
		}
	}
	return nil
}

// ============ Config 操作 ============

// SaveConfig 保存配置（加密API Key）
func (d *DB) SaveConfig(config *model.Config) error {
	// 加密 API Key
	encryptedKey, err := d.crypto.Encrypt(config.APIKey)
	if err != nil {
		return fmt.Errorf("encrypt api key failed: %w", err)
	}

	config.UpdatedAt = time.Now()
	if config.CreatedAt.IsZero() {
		config.CreatedAt = time.Now()
	}

	query := `INSERT OR REPLACE INTO config 
		(id, api_base_url, api_key, model, temperature, style_tags, audience, persona, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = d.conn.Exec(query,
		config.ID,
		config.APIBaseURL,
		encryptedKey,
		config.Model,
		config.Temperature,
		strings.Join(config.StyleTags, ","),
		config.Audience,
		config.Persona,
		config.CreatedAt,
		config.UpdatedAt,
	)
	return err
}

// GetConfig 获取配置（解密API Key）
func (d *DB) GetConfig() (*model.Config, error) {
	query := `SELECT id, api_base_url, api_key, model, temperature, style_tags, audience, persona, created_at, updated_at 
		FROM config LIMIT 1`

	row := d.conn.QueryRow(query)

	var config model.Config
	var encryptedKey, tagsStr string

	err := row.Scan(
		&config.ID,
		&config.APIBaseURL,
		&encryptedKey,
		&config.Model,
		&config.Temperature,
		&tagsStr,
		&config.Audience,
		&config.Persona,
		&config.CreatedAt,
		&config.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("config not found")
	}
	if err != nil {
		return nil, err
	}
	if encryptedKey != "" {
		// 解密 API Key
		decryptedKey, err := d.crypto.Decrypt(encryptedKey)
		if err != nil {
			return nil, fmt.Errorf("decrypt api key failed: %w", err)
		}
		config.APIKey = decryptedKey
	}

	if tagsStr != "" {
		config.StyleTags = strings.Split(tagsStr, ",")
	} else {
		config.StyleTags = []string{}
	}

	return &config, nil
}

// HasConfig 检查是否有配置
func (d *DB) HasConfig() bool {
	var count int
	err := d.conn.QueryRow("SELECT COUNT(*) FROM config").Scan(&count)
	return err == nil && count > 0
}

// ============ Article 操作 ============

// SaveArticle 保存文章
func (d *DB) SaveArticle(article *model.Article) error {
	article.UpdatedAt = time.Now()

	query := `INSERT OR REPLACE INTO articles 
		(id, title, content, outline, status, source_type, workflow_run_id, topic_id,
		 quality_score, read_time, word_count, publish_task_id, published_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	outlineJSON, _ := json.Marshal(article.Outline)

	_, err := d.conn.Exec(query,
		article.ID,
		article.Title,
		article.Content,
		string(outlineJSON),
		article.Status,
		article.SourceType,
		article.WorkflowRunID,
		article.TopicID,
		article.QualityScore,
		article.ReadTime,
		article.WordCount,
		article.PublishTaskID,
		article.PublishedAt,
		article.CreatedAt,
		article.UpdatedAt,
	)
	return err
}

// GetArticle 获取文章
func (d *DB) GetArticle(id string) (*model.Article, error) {
	query := `SELECT id, title, content, outline, status, source_type, workflow_run_id, topic_id,
		quality_score, read_time, word_count, publish_task_id, published_at, created_at, updated_at 
		FROM articles WHERE id = ?`

	row := d.conn.QueryRow(query, id)

	var article model.Article
	var outlineStr string

	err := row.Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&outlineStr,
		&article.Status,
		&article.SourceType,
		&article.WorkflowRunID,
		&article.TopicID,
		&article.QualityScore,
		&article.ReadTime,
		&article.WordCount,
		&article.PublishTaskID,
		&article.PublishedAt,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if outlineStr != "" {
		json.Unmarshal([]byte(outlineStr), &article.Outline)
	}

	return &article, nil
}

// ListArticles 列出所有文章
func (d *DB) ListArticles() ([]model.Article, error) {
	query := `SELECT id, title, content, outline, status, source_type, workflow_run_id, topic_id,
		quality_score, read_time, word_count, publish_task_id, published_at, created_at, updated_at 
		FROM articles ORDER BY updated_at DESC`

	rows, err := d.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		var outlineStr string

		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&outlineStr,
			&article.Status,
			&article.SourceType,
			&article.WorkflowRunID,
			&article.TopicID,
			&article.QualityScore,
			&article.ReadTime,
			&article.WordCount,
			&article.PublishTaskID,
			&article.PublishedAt,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if outlineStr != "" {
			json.Unmarshal([]byte(outlineStr), &article.Outline)
		}

		articles = append(articles, article)
	}

	return articles, nil
}

// DeleteArticle 删除文章
func (d *DB) DeleteArticle(id string) error {
	_, err := d.conn.Exec("DELETE FROM articles WHERE id = ?", id)
	return err
}

// ============ Material 操作 ============

// SaveMaterial 保存素材
func (d *DB) SaveMaterial(material *model.Material) error {
	query := `INSERT OR REPLACE INTO materials 
		(id, material_type, title, content, tags, source, created_at, usage_count)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := d.conn.Exec(query,
		material.ID,
		material.Type,
		material.Title,
		material.Content,
		strings.Join(material.Tags, ","),
		material.Source,
		material.CreatedAt,
		material.UsageCount,
	)
	return err
}

// GetMaterials 获取素材列表
func (d *DB) GetMaterials(materialType string) ([]model.Material, error) {
	var query string
	var args []interface{}

	if materialType == "" {
		query = `SELECT id, material_type, title, content, tags, source, created_at, usage_count 
			FROM materials ORDER BY usage_count DESC, created_at DESC`
	} else {
		query = `SELECT id, material_type, title, content, tags, source, created_at, usage_count 
			FROM materials WHERE material_type = ? ORDER BY usage_count DESC, created_at DESC`
		args = append(args, materialType)
	}

	rows, err := d.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []model.Material
	for rows.Next() {
		var m model.Material
		var tagsStr string

		err := rows.Scan(
			&m.ID,
			&m.Type,
			&m.Title,
			&m.Content,
			&tagsStr,
			&m.Source,
			&m.CreatedAt,
			&m.UsageCount,
		)
		if err != nil {
			continue
		}

		if tagsStr != "" {
			m.Tags = strings.Split(tagsStr, ",")
		}
		materials = append(materials, m)
	}

	return materials, nil
}

// DeleteMaterial 删除素材
func (d *DB) DeleteMaterial(id string) error {
	_, err := d.conn.Exec("DELETE FROM materials WHERE id = ?", id)
	return err
}
