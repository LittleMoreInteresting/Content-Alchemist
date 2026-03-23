package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"content-alchemist/internal/model"
)

// ============ Topic 操作 ============

// SaveTopic 保存选题
func (d *DB) SaveTopic(topic *model.Topic) error {
	topic.UpdatedAt = time.Now()
	if topic.CreatedAt.IsZero() {
		topic.CreatedAt = time.Now()
	}

	keywordsJSON, _ := json.Marshal(topic.Keywords)
	referencesJSON, _ := json.Marshal(topic.References)
	anglesJSON, _ := json.Marshal(topic.Angles)

	query := `INSERT OR REPLACE INTO topics 
		(id, title, category, source, source_url, score, hot_score, comp_score, fit_score, 
		 keywords, summary, reference_urls, angles, status, workflow_run_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := d.conn.Exec(query,
		topic.ID,
		topic.Title,
		topic.Category,
		topic.Source,
		topic.SourceURL,
		topic.Score,
		topic.HotScore,
		topic.CompScore,
		topic.FitScore,
		string(keywordsJSON),
		topic.Summary,
		string(referencesJSON),
		string(anglesJSON),
		topic.Status,
		topic.WorkflowRunID,
		topic.CreatedAt,
		topic.UpdatedAt,
	)
	return err
}

// GetTopic 获取选题
func (d *DB) GetTopic(id string) (*model.Topic, error) {
	query := `SELECT id, title, category, source, source_url, score, hot_score, comp_score, fit_score, 
		keywords, summary, reference_urls, angles, status, workflow_run_id, created_at, updated_at 
		FROM topics WHERE id = ?`

	row := d.conn.QueryRow(query, id)

	var topic model.Topic
	var keywordsStr, referencesStr, anglesStr string

	err := row.Scan(
		&topic.ID,
		&topic.Title,
		&topic.Category,
		&topic.Source,
		&topic.SourceURL,
		&topic.Score,
		&topic.HotScore,
		&topic.CompScore,
		&topic.FitScore,
		&keywordsStr,
		&topic.Summary,
		&referencesStr,
		&anglesStr,
		&topic.Status,
		&topic.WorkflowRunID,
		&topic.CreatedAt,
		&topic.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("topic not found")
	}
	if err != nil {
		return nil, err
	}

	if keywordsStr != "" {
		json.Unmarshal([]byte(keywordsStr), &topic.Keywords)
	}
	if referencesStr != "" {
		json.Unmarshal([]byte(referencesStr), &topic.References)
	}
	if anglesStr != "" {
		json.Unmarshal([]byte(anglesStr), &topic.Angles)
	}

	return &topic, nil
}

// ListTopics 列出选题
func (d *DB) ListTopics(status string, limit int) ([]model.Topic, error) {
	var query string
	var args []interface{}

	if status == "" {
		query = `SELECT id, title, category, source, source_url, score, hot_score, comp_score, fit_score, 
			keywords, summary, reference_urls, angles, status, workflow_run_id, created_at, updated_at 
			FROM topics ORDER BY score DESC, created_at DESC`
	} else {
		query = `SELECT id, title, category, source, source_url, score, hot_score, comp_score, fit_score, 
			keywords, summary, reference_urls, angles, status, workflow_run_id, created_at, updated_at 
			FROM topics WHERE status = ? ORDER BY score DESC, created_at DESC`
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

	var topics []model.Topic
	for rows.Next() {
		var topic model.Topic
		var keywordsStr, referencesStr, anglesStr string

		err := rows.Scan(
			&topic.ID,
			&topic.Title,
			&topic.Category,
			&topic.Source,
			&topic.SourceURL,
			&topic.Score,
			&topic.HotScore,
			&topic.CompScore,
			&topic.FitScore,
			&keywordsStr,
			&topic.Summary,
			&referencesStr,
			&anglesStr,
			&topic.Status,
			&topic.WorkflowRunID,
			&topic.CreatedAt,
			&topic.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if keywordsStr != "" {
			json.Unmarshal([]byte(keywordsStr), &topic.Keywords)
		}
		if referencesStr != "" {
			json.Unmarshal([]byte(referencesStr), &topic.References)
		}
		if anglesStr != "" {
			json.Unmarshal([]byte(anglesStr), &topic.Angles)
		}

		topics = append(topics, topic)
	}

	return topics, nil
}

// ListTopicsBySource 根据来源列出选题
func (d *DB) ListTopicsBySource(source string, limit int) ([]model.Topic, error) {
	query := `SELECT id, title, category, source, source_url, score, hot_score, comp_score, fit_score, 
		keywords, summary, reference_urls, angles, status, workflow_run_id, created_at, updated_at 
		FROM topics WHERE source = ? ORDER BY score DESC, created_at DESC`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := d.conn.Query(query, source)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []model.Topic
	for rows.Next() {
		var topic model.Topic
		var keywordsStr, referencesStr, anglesStr string

		err := rows.Scan(
			&topic.ID,
			&topic.Title,
			&topic.Category,
			&topic.Source,
			&topic.SourceURL,
			&topic.Score,
			&topic.HotScore,
			&topic.CompScore,
			&topic.FitScore,
			&keywordsStr,
			&topic.Summary,
			&referencesStr,
			&anglesStr,
			&topic.Status,
			&topic.WorkflowRunID,
			&topic.CreatedAt,
			&topic.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if keywordsStr != "" {
			json.Unmarshal([]byte(keywordsStr), &topic.Keywords)
		}
		if referencesStr != "" {
			json.Unmarshal([]byte(referencesStr), &topic.References)
		}
		if anglesStr != "" {
			json.Unmarshal([]byte(anglesStr), &topic.Angles)
		}

		topics = append(topics, topic)
	}

	return topics, nil
}

// DeleteTopic 删除选题
func (d *DB) DeleteTopic(id string) error {
	_, err := d.conn.Exec("DELETE FROM topics WHERE id = ?", id)
	return err
}

// SearchTopics 搜索选题
func (d *DB) SearchTopics(keyword string) ([]model.Topic, error) {
	query := `SELECT id, title, category, source, source_url, score, hot_score, comp_score, fit_score, 
		keywords, summary, reference_urls, angles, status, workflow_run_id, created_at, updated_at 
		FROM topics WHERE title LIKE ? OR summary LIKE ? OR keywords LIKE ?
		ORDER BY score DESC, created_at DESC`

	likeKeyword := "%" + keyword + "%"
	rows, err := d.conn.Query(query, likeKeyword, likeKeyword, likeKeyword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []model.Topic
	for rows.Next() {
		var topic model.Topic
		var keywordsStr, referencesStr, anglesStr string

		err := rows.Scan(
			&topic.ID,
			&topic.Title,
			&topic.Category,
			&topic.Source,
			&topic.SourceURL,
			&topic.Score,
			&topic.HotScore,
			&topic.CompScore,
			&topic.FitScore,
			&keywordsStr,
			&topic.Summary,
			&referencesStr,
			&anglesStr,
			&topic.Status,
			&topic.WorkflowRunID,
			&topic.CreatedAt,
			&topic.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if keywordsStr != "" {
			json.Unmarshal([]byte(keywordsStr), &topic.Keywords)
		}
		if referencesStr != "" {
			json.Unmarshal([]byte(referencesStr), &topic.References)
		}
		if anglesStr != "" {
			json.Unmarshal([]byte(anglesStr), &topic.Angles)
		}

		topics = append(topics, topic)
	}

	return topics, nil
}

// ============ HotTrend 操作 ============

// SaveHotTrend 保存热点趋势
func (d *DB) SaveHotTrend(trend *model.HotTrend) error {
	if trend.CreatedAt.IsZero() {
		trend.CreatedAt = time.Now()
	}

	query := `INSERT OR REPLACE INTO hot_trends 
		(id, platform, title, url, hot_rank, hot_value, category, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := d.conn.Exec(query,
		trend.ID,
		trend.Platform,
		trend.Title,
		trend.URL,
		trend.HotRank,
		trend.HotValue,
		trend.Category,
		trend.CreatedAt,
	)
	return err
}

// GetHotTrends 获取热点趋势列表
func (d *DB) GetHotTrends(platform string, limit int) ([]model.HotTrend, error) {
	var query string
	var args []interface{}

	if platform == "" {
		query = `SELECT id, platform, title, url, hot_rank, hot_value, category, created_at 
			FROM hot_trends ORDER BY hot_value DESC, created_at DESC`
	} else {
		query = `SELECT id, platform, title, url, hot_rank, hot_value, category, created_at 
			FROM hot_trends WHERE platform = ? ORDER BY hot_value DESC, created_at DESC`
		args = append(args, platform)
	}

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := d.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []model.HotTrend
	for rows.Next() {
		var trend model.HotTrend

		err := rows.Scan(
			&trend.ID,
			&trend.Platform,
			&trend.Title,
			&trend.URL,
			&trend.HotRank,
			&trend.HotValue,
			&trend.Category,
			&trend.CreatedAt,
		)
		if err != nil {
			continue
		}

		trends = append(trends, trend)
	}

	return trends, nil
}

// DeleteOldHotTrends 删除旧的热点趋势数据
func (d *DB) DeleteOldHotTrends(before time.Time) error {
	_, err := d.conn.Exec("DELETE FROM hot_trends WHERE created_at < ?", before)
	return err
}

// ClearHotTrendsByPlatform 清空指定平台的热点数据
func (d *DB) ClearHotTrendsByPlatform(platform string) error {
	_, err := d.conn.Exec("DELETE FROM hot_trends WHERE platform = ?", platform)
	return err
}

// BatchSaveHotTrends 批量保存热点趋势
func (d *DB) BatchSaveHotTrends(trends []model.HotTrend) error {
	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT OR REPLACE INTO hot_trends 
		(id, platform, title, url, hot_rank, hot_value, category, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()
	for _, trend := range trends {
		if trend.CreatedAt.IsZero() {
			trend.CreatedAt = now
		}
		if trend.ID == "" {
			trend.ID = fmt.Sprintf("%s_%d_%s", trend.Platform, trend.HotRank, now.Format("20060102"))
		}

		_, err := stmt.Exec(
			trend.ID,
			trend.Platform,
			trend.Title,
			trend.URL,
			trend.HotRank,
			trend.HotValue,
			trend.Category,
			trend.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetHotTrendPlatforms 获取所有热点平台
func (d *DB) GetHotTrendPlatforms() ([]string, error) {
	query := `SELECT DISTINCT platform FROM hot_trends ORDER BY platform`

	rows, err := d.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var platforms []string
	for rows.Next() {
		var platform string
		if err := rows.Scan(&platform); err != nil {
			continue
		}
		platforms = append(platforms, platform)
	}

	return platforms, nil
}

// SearchHotTrends 搜索热点趋势
func (d *DB) SearchHotTrends(keyword string) ([]model.HotTrend, error) {
	query := `SELECT id, platform, title, url, hot_rank, hot_value, category, created_at 
		FROM hot_trends WHERE title LIKE ?
		ORDER BY hot_value DESC, created_at DESC`

	likeKeyword := "%" + keyword + "%"
	rows, err := d.conn.Query(query, likeKeyword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []model.HotTrend
	for rows.Next() {
		var trend model.HotTrend

		err := rows.Scan(
			&trend.ID,
			&trend.Platform,
			&trend.Title,
			&trend.URL,
			&trend.HotRank,
			&trend.HotValue,
			&trend.Category,
			&trend.CreatedAt,
		)
		if err != nil {
			continue
		}

		trends = append(trends, trend)
	}

	return trends, nil
}

// IsDuplicateTrend 检查是否是重复的热点
func (d *DB) IsDuplicateTrend(platform string, title string) (bool, error) {
	// 规范化标题进行比较
	normalizedTitle := strings.TrimSpace(title)
	
	var count int
	query := `SELECT COUNT(*) FROM hot_trends WHERE platform = ? AND title = ? AND created_at > datetime('now', '-1 day')`
	err := d.conn.QueryRow(query, platform, normalizedTitle).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
