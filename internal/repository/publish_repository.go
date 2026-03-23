package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"content-alchemist/internal/model"
)

// ============ PublishAccount 操作 ============

// SavePublishAccount 保存发布账号
func (d *DB) SavePublishAccount(account *model.PublishAccount) error {
	// 加密敏感信息
	encryptedSecret, _ := d.crypto.Encrypt(account.AppSecret)
	encryptedToken, _ := d.crypto.Encrypt(account.AccessToken)
	encryptedRefresh, _ := d.crypto.Encrypt(account.RefreshToken)

	configJSON, _ := json.Marshal(account.Config)

	query := `INSERT OR REPLACE INTO publish_accounts 
		(id, platform, name, account_id, account_type, app_id, app_secret, access_token, 
		 token_expiry, refresh_token, config, status, last_used_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := d.conn.Exec(query,
		account.ID,
		account.Platform,
		account.Name,
		account.AccountID,
		account.AccountType,
		account.AppID,
		encryptedSecret,
		encryptedToken,
		account.TokenExpiry,
		encryptedRefresh,
		string(configJSON),
		account.Status,
		account.LastUsedAt,
		account.CreatedAt,
	)
	return err
}

// GetPublishAccount 获取发布账号
func (d *DB) GetPublishAccount(id string) (*model.PublishAccount, error) {
	query := `SELECT id, platform, name, account_id, account_type, app_id, app_secret, access_token, 
		token_expiry, refresh_token, config, status, last_used_at, created_at 
		FROM publish_accounts WHERE id = ?`

	row := d.conn.QueryRow(query, id)

	var account model.PublishAccount
	var secretStr, tokenStr, refreshStr, configStr string

	err := row.Scan(
		&account.ID,
		&account.Platform,
		&account.Name,
		&account.AccountID,
		&account.AccountType,
		&account.AppID,
		&secretStr,
		&tokenStr,
		&account.TokenExpiry,
		&refreshStr,
		&configStr,
		&account.Status,
		&account.LastUsedAt,
		&account.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("publish account not found")
	}
	if err != nil {
		return nil, err
	}

	// 解密敏感信息
	if secretStr != "" {
		account.AppSecret, _ = d.crypto.Decrypt(secretStr)
	}
	if tokenStr != "" {
		account.AccessToken, _ = d.crypto.Decrypt(tokenStr)
	}
	if refreshStr != "" {
		account.RefreshToken, _ = d.crypto.Decrypt(refreshStr)
	}
	if configStr != "" {
		json.Unmarshal([]byte(configStr), &account.Config)
	}

	return &account, nil
}

// GetPublishAccountByPlatform 根据平台获取发布账号
func (d *DB) GetPublishAccountByPlatform(platform string) (*model.PublishAccount, error) {
	query := `SELECT id, platform, name, account_id, account_type, app_id, app_secret, access_token, 
		token_expiry, refresh_token, config, status, last_used_at, created_at 
		FROM publish_accounts WHERE platform = ? AND status = 'active' LIMIT 1`

	row := d.conn.QueryRow(query, platform)

	var account model.PublishAccount
	var secretStr, tokenStr, refreshStr, configStr string

	err := row.Scan(
		&account.ID,
		&account.Platform,
		&account.Name,
		&account.AccountID,
		&account.AccountType,
		&account.AppID,
		&secretStr,
		&tokenStr,
		&account.TokenExpiry,
		&refreshStr,
		&configStr,
		&account.Status,
		&account.LastUsedAt,
		&account.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no active account found for platform: %s", platform)
	}
	if err != nil {
		return nil, err
	}

	// 解密敏感信息
	if secretStr != "" {
		account.AppSecret, _ = d.crypto.Decrypt(secretStr)
	}
	if tokenStr != "" {
		account.AccessToken, _ = d.crypto.Decrypt(tokenStr)
	}
	if refreshStr != "" {
		account.RefreshToken, _ = d.crypto.Decrypt(refreshStr)
	}
	if configStr != "" {
		json.Unmarshal([]byte(configStr), &account.Config)
	}

	return &account, nil
}

// ListPublishAccounts 列出所有发布账号
func (d *DB) ListPublishAccounts() ([]model.PublishAccount, error) {
	query := `SELECT id, platform, name, account_id, account_type, app_id, 
		token_expiry, config, status, last_used_at, created_at 
		FROM publish_accounts ORDER BY created_at DESC`

	rows, err := d.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []model.PublishAccount
	for rows.Next() {
		var account model.PublishAccount
		var configStr string

		err := rows.Scan(
			&account.ID,
			&account.Platform,
			&account.Name,
			&account.AccountID,
			&account.AccountType,
			&account.AppID,
			&account.TokenExpiry,
			&configStr,
			&account.Status,
			&account.LastUsedAt,
			&account.CreatedAt,
		)
		if err != nil {
			continue
		}

		if configStr != "" {
			json.Unmarshal([]byte(configStr), &account.Config)
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

// DeletePublishAccount 删除发布账号
func (d *DB) DeletePublishAccount(id string) error {
	_, err := d.conn.Exec("DELETE FROM publish_accounts WHERE id = ?", id)
	return err
}

// UpdateAccountLastUsed 更新账号最后使用时间
func (d *DB) UpdateAccountLastUsed(id string) error {
	_, err := d.conn.Exec(
		"UPDATE publish_accounts SET last_used_at = ? WHERE id = ?",
		time.Now(), id,
	)
	return err
}

// ============ PublishTask 操作 ============

// SavePublishTask 保存发布任务
func (d *DB) SavePublishTask(task *model.PublishTask) error {
	task.UpdatedAt = time.Now()
	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}

	tagsJSON, _ := json.Marshal(task.Tags)

	query := `INSERT OR REPLACE INTO publish_tasks 
		(id, article_id, account_id, title, content, cover_image, summary, author, tags, category,
		 original_url, schedule_type, schedule_at, status, platform_id, platform_url, error, 
		 published_at, retry_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := d.conn.Exec(query,
		task.ID,
		task.ArticleID,
		task.AccountID,
		task.Title,
		task.Content,
		task.CoverImage,
		task.Summary,
		task.Author,
		string(tagsJSON),
		task.Category,
		task.OriginalURL,
		task.ScheduleType,
		task.ScheduleAt,
		task.Status,
		task.PlatformID,
		task.PlatformURL,
		task.Error,
		task.PublishedAt,
		task.RetryCount,
		task.CreatedAt,
		task.UpdatedAt,
	)
	return err
}

// GetPublishTask 获取发布任务
func (d *DB) GetPublishTask(id string) (*model.PublishTask, error) {
	query := `SELECT id, article_id, account_id, title, content, cover_image, summary, author, tags, category,
		original_url, schedule_type, schedule_at, status, platform_id, platform_url, error, 
		published_at, retry_count, created_at, updated_at 
		FROM publish_tasks WHERE id = ?`

	row := d.conn.QueryRow(query, id)

	var task model.PublishTask
	var tagsStr string

	err := row.Scan(
		&task.ID,
		&task.ArticleID,
		&task.AccountID,
		&task.Title,
		&task.Content,
		&task.CoverImage,
		&task.Summary,
		&task.Author,
		&tagsStr,
		&task.Category,
		&task.OriginalURL,
		&task.ScheduleType,
		&task.ScheduleAt,
		&task.Status,
		&task.PlatformID,
		&task.PlatformURL,
		&task.Error,
		&task.PublishedAt,
		&task.RetryCount,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("publish task not found")
	}
	if err != nil {
		return nil, err
	}

	if tagsStr != "" {
		json.Unmarshal([]byte(tagsStr), &task.Tags)
	}

	return &task, nil
}

// ListPublishTasks 列出发布任务
func (d *DB) ListPublishTasks(status string, limit int) ([]model.PublishTask, error) {
	var query string
	var args []interface{}

	if status == "" {
		query = `SELECT id, article_id, account_id, title, content, cover_image, summary, author, tags, category,
			original_url, schedule_type, schedule_at, status, platform_id, platform_url, error, 
			published_at, retry_count, created_at, updated_at 
			FROM publish_tasks ORDER BY created_at DESC`
	} else {
		query = `SELECT id, article_id, account_id, title, content, cover_image, summary, author, tags, category,
			original_url, schedule_type, schedule_at, status, platform_id, platform_url, error, 
			published_at, retry_count, created_at, updated_at 
			FROM publish_tasks WHERE status = ? ORDER BY created_at DESC`
		args = append(args, status)
	}

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := d.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.PublishTask
	for rows.Next() {
		var task model.PublishTask
		var tagsStr string

		err := rows.Scan(
			&task.ID,
			&task.ArticleID,
			&task.AccountID,
			&task.Title,
			&task.Content,
			&task.CoverImage,
			&task.Summary,
			&task.Author,
			&tagsStr,
			&task.Category,
			&task.OriginalURL,
			&task.ScheduleType,
			&task.ScheduleAt,
			&task.Status,
			&task.PlatformID,
			&task.PlatformURL,
			&task.Error,
			&task.PublishedAt,
			&task.RetryCount,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if tagsStr != "" {
			json.Unmarshal([]byte(tagsStr), &task.Tags)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetScheduledTasks 获取待调度的任务
func (d *DB) GetScheduledTasks(before time.Time) ([]model.PublishTask, error) {
	query := `SELECT id, article_id, account_id, title, content, cover_image, summary, author, tags, category,
		original_url, schedule_type, schedule_at, status, platform_id, platform_url, error, 
		published_at, retry_count, created_at, updated_at 
		FROM publish_tasks 
		WHERE status = 'scheduled' AND schedule_at <= ?
		ORDER BY schedule_at ASC`

	rows, err := d.conn.Query(query, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.PublishTask
	for rows.Next() {
		var task model.PublishTask
		var tagsStr string

		err := rows.Scan(
			&task.ID,
			&task.ArticleID,
			&task.AccountID,
			&task.Title,
			&task.Content,
			&task.CoverImage,
			&task.Summary,
			&task.Author,
			&tagsStr,
			&task.Category,
			&task.OriginalURL,
			&task.ScheduleType,
			&task.ScheduleAt,
			&task.Status,
			&task.PlatformID,
			&task.PlatformURL,
			&task.Error,
			&task.PublishedAt,
			&task.RetryCount,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if tagsStr != "" {
			json.Unmarshal([]byte(tagsStr), &task.Tags)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetPendingTasks 获取待处理的任务
func (d *DB) GetPendingTasks(limit int) ([]model.PublishTask, error) {
	query := `SELECT id, article_id, account_id, title, content, cover_image, summary, author, tags, category,
		original_url, schedule_type, schedule_at, status, platform_id, platform_url, error, 
		published_at, retry_count, created_at, updated_at 
		FROM publish_tasks 
		WHERE status = 'pending'
		ORDER BY created_at ASC`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := d.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.PublishTask
	for rows.Next() {
		var task model.PublishTask
		var tagsStr string

		err := rows.Scan(
			&task.ID,
			&task.ArticleID,
			&task.AccountID,
			&task.Title,
			&task.Content,
			&task.CoverImage,
			&task.Summary,
			&task.Author,
			&tagsStr,
			&task.Category,
			&task.OriginalURL,
			&task.ScheduleType,
			&task.ScheduleAt,
			&task.Status,
			&task.PlatformID,
			&task.PlatformURL,
			&task.Error,
			&task.PublishedAt,
			&task.RetryCount,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if tagsStr != "" {
			json.Unmarshal([]byte(tagsStr), &task.Tags)
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// DeletePublishTask 删除发布任务
func (d *DB) DeletePublishTask(id string) error {
	_, err := d.conn.Exec("DELETE FROM publish_tasks WHERE id = ?", id)
	return err
}

// GetPublishStats 获取发布统计
func (d *DB) GetPublishStats() (*model.PublishStats, error) {
	stats := &model.PublishStats{}

	// 总数
	d.conn.QueryRow("SELECT COUNT(*) FROM publish_tasks").Scan(&stats.TotalCount)

	// 成功数
	d.conn.QueryRow("SELECT COUNT(*) FROM publish_tasks WHERE status = 'published'").Scan(&stats.SuccessCount)

	// 失败数
	d.conn.QueryRow("SELECT COUNT(*) FROM publish_tasks WHERE status = 'failed'").Scan(&stats.FailedCount)

	// 待处理数
	d.conn.QueryRow("SELECT COUNT(*) FROM publish_tasks WHERE status = 'pending'").Scan(&stats.PendingCount)

	// 定时数
	d.conn.QueryRow("SELECT COUNT(*) FROM publish_tasks WHERE status = 'scheduled'").Scan(&stats.ScheduleCount)

	return stats, nil
}

// CancelScheduledTask 取消定时任务
func (d *DB) CancelScheduledTask(id string) error {
	_, err := d.conn.Exec(
		"UPDATE publish_tasks SET status = 'cancelled', updated_at = ? WHERE id = ? AND status = 'scheduled'",
		time.Now(), id,
	)
	return err
}
