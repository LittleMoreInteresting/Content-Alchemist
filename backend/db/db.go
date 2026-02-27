//go:build darwin || windows || linux

package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"Content-Alchemist/backend/models"

	_ "modernc.org/sqlite"
)

// DB 数据库连接包装器
type DB struct {
	conn *sql.DB
}

// sql.NullString 辅助函数
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// nullInt64ToInt 辅助函数
func nullInt64ToInt(ni sql.NullInt64) int {
	if ni.Valid {
		return int(ni.Int64)
	}
	return 0
}

// nullInt64ToTime 辅助函数
func nullInt64ToTime(ni sql.NullInt64) time.Time {
	if ni.Valid {
		return time.Unix(ni.Int64, 0)
	}
	return time.Time{}
}

// New 创建数据库连接
func New(dataDir string) (*DB, error) {
	dbPath := filepath.Join(dataDir, "content-alchemist.db")

	// 确保目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	conn, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 设置连接池
	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)

	db := &DB{conn: conn}

	// 初始化Schema
	if err := db.initSchema(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to init schema: %w", err)
	}

	return db, nil
}

// Close 关闭数据库连接
func (d *DB) Close() error {
	return d.conn.Close()
}

// initSchema 初始化数据库表结构
func (d *DB) initSchema() error {
	schema := `
-- 文章元数据表
CREATE TABLE IF NOT EXISTS articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT UNIQUE NOT NULL,
    file_path TEXT UNIQUE NOT NULL,
    title TEXT,
    summary TEXT,
    tags TEXT DEFAULT '[]',
    word_count INTEGER DEFAULT 0,
    created_at INTEGER,
    updated_at INTEGER,
    last_opened_at INTEGER
);

-- 文章索引
CREATE INDEX IF NOT EXISTS idx_articles_uuid ON articles(uuid);
CREATE INDEX IF NOT EXISTS idx_articles_last_opened ON articles(last_opened_at DESC);

-- 应用配置表
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT
);

-- 初始化默认配置
INSERT OR IGNORE INTO settings (key, value) VALUES
('deepseek_api_key', ''),
('editor_theme', 'github-dark'),
('font_size', '16'),
('ai_model', 'deepseek-chat'),
('window_width', '1400'),
('window_height', '900');
`

	_, err := d.conn.Exec(schema)
	return err
}

// GetArticleByPath 根据文件路径获取文章
func (d *DB) GetArticleByPath(filePath string) (*models.Article, error) {
	var article models.Article
	var tagsJSON string
	var createdAt, updatedAt, lastOpenedAt sql.NullInt64

	row := d.conn.QueryRow(
		`SELECT id, uuid, file_path, title, summary, tags, word_count, created_at, updated_at, last_opened_at
		 FROM articles WHERE file_path = ?`,
		filePath,
	)

	err := row.Scan(
		&article.ID,
		&article.UUID,
		&article.FilePath,
		&article.Title,
		&article.Summary,
		&tagsJSON,
		&article.WordCount,
		&createdAt,
		&updatedAt,
		&lastOpenedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan article: %w", err)
	}

	// 解析时间戳
	article.CreatedAt = nullInt64ToTime(createdAt)
	article.UpdatedAt = nullInt64ToTime(updatedAt)
	article.LastOpenedAt = nullInt64ToTime(lastOpenedAt)

	// 解析tags JSON
	if err := article.Tags.Scan(tagsJSON); err != nil {
		return nil, fmt.Errorf("failed to scan tags: %w", err)
	}

	return &article, nil
}

// GetArticleByUUID 根据UUID获取文章
func (d *DB) GetArticleByUUID(uuid string) (*models.Article, error) {
	var article models.Article
	var tagsJSON string
	var createdAt, updatedAt, lastOpenedAt sql.NullInt64

	row := d.conn.QueryRow(
		`SELECT id, uuid, file_path, title, summary, tags, word_count, created_at, updated_at, last_opened_at
		 FROM articles WHERE uuid = ?`,
		uuid,
	)

	err := row.Scan(
		&article.ID,
		&article.UUID,
		&article.FilePath,
		&article.Title,
		&article.Summary,
		&tagsJSON,
		&article.WordCount,
		&createdAt,
		&updatedAt,
		&lastOpenedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan article: %w", err)
	}

	// 解析时间戳
	article.CreatedAt = nullInt64ToTime(createdAt)
	article.UpdatedAt = nullInt64ToTime(updatedAt)
	article.LastOpenedAt = nullInt64ToTime(lastOpenedAt)

	// 解析tags JSON
	if err := article.Tags.Scan(tagsJSON); err != nil {
		return nil, fmt.Errorf("failed to scan tags: %w", err)
	}

	return &article, nil
}

// CreateArticle 创建新文章记录
func (d *DB) CreateArticle(article *models.Article) error {
	now := time.Now().Unix()

	result, err := d.conn.Exec(
		`INSERT INTO articles (uuid, file_path, title, summary, tags, word_count, created_at, updated_at, last_opened_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		article.UUID,
		article.FilePath,
		article.Title,
		article.Summary,
		article.Tags,
		article.WordCount,
		now,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to insert article: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	article.ID = id
	article.CreatedAt = time.Unix(now, 0)
	article.UpdatedAt = time.Unix(now, 0)
	article.LastOpenedAt = time.Unix(now, 0)

	return nil
}

// UpdateArticle 更新文章记录
func (d *DB) UpdateArticle(article *models.Article) error {
	now := time.Now().Unix()

	_, err := d.conn.Exec(
		`UPDATE articles
		 SET title = ?, summary = ?, tags = ?, word_count = ?, updated_at = ?
		 WHERE uuid = ?`,
		article.Title,
		article.Summary,
		article.Tags,
		article.WordCount,
		now,
		article.UUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update article: %w", err)
	}

	article.UpdatedAt = time.Unix(now, 0)
	return nil
}

// UpdateLastOpened 更新最后打开时间
func (d *DB) UpdateLastOpened(uuid string) error {
	now := time.Now().Unix()
	_, err := d.conn.Exec(
		`UPDATE articles SET last_opened_at = ? WHERE uuid = ?`,
		now,
		uuid,
	)
	return err
}

// GetRecentArticles 获取最近打开的文章列表
func (d *DB) GetRecentArticles(limit int) ([]*models.Article, error) {
	rows, err := d.conn.Query(
		`SELECT id, uuid, file_path, title, summary, tags, word_count, created_at, updated_at, last_opened_at
		 FROM articles
		 WHERE last_opened_at IS NOT NULL
		 ORDER BY last_opened_at DESC
		 LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent articles: %w", err)
	}
	defer rows.Close()

	return d.scanArticles(rows)
}

// scanArticles 扫描文章列表
func (d *DB) scanArticles(rows *sql.Rows) ([]*models.Article, error) {
	var articles []*models.Article

	for rows.Next() {
		var article models.Article
		var tagsJSON string
		var createdAt, updatedAt, lastOpenedAt sql.NullInt64

		err := rows.Scan(
			&article.ID,
			&article.UUID,
			&article.FilePath,
			&article.Title,
			&article.Summary,
			&tagsJSON,
			&article.WordCount,
			&createdAt,
			&updatedAt,
			&lastOpenedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan article row: %w", err)
		}

		// 解析时间戳
		article.CreatedAt = nullInt64ToTime(createdAt)
		article.UpdatedAt = nullInt64ToTime(updatedAt)
		article.LastOpenedAt = nullInt64ToTime(lastOpenedAt)

		// 解析tags JSON
		if err := article.Tags.Scan(tagsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan tags: %w", err)
		}

		articles = append(articles, &article)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return articles, nil
}

// DeleteArticle 删除文章记录
func (d *DB) DeleteArticle(uuid string) error {
	_, err := d.conn.Exec(`DELETE FROM articles WHERE uuid = ?`, uuid)
	return err
}

// GetSetting 获取设置项
func (d *DB) GetSetting(key string) (string, error) {
	var value string
	row := d.conn.QueryRow(`SELECT value FROM settings WHERE key = ?`, key)
	err := row.Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

// SetSetting 设置配置项
func (d *DB) SetSetting(key, value string) error {
	_, err := d.conn.Exec(
		`INSERT INTO settings (key, value) VALUES (?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key,
		value,
	)
	return err
}

// GetAllSettings 获取所有设置
func (d *DB) GetAllSettings() (map[string]string, error) {
	rows, err := d.conn.Query(`SELECT key, value FROM settings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		settings[key] = value
	}

	return settings, rows.Err()
}
