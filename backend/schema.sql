-- Content Alchemist 数据库 Schema
-- SQLite 数据库初始化脚本

-- ============================================
-- 文章元数据表
-- ============================================
CREATE TABLE IF NOT EXISTS articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT UNIQUE NOT NULL,           -- 业务ID，前端识别用 (UUID v4)
    file_path TEXT UNIQUE NOT NULL,      -- 本地绝对路径 (规范化后)
    title TEXT,                          -- 缓存标题（H1或文件名）
    summary TEXT,                        -- 文章摘要（前200字）
    tags TEXT DEFAULT '[]',              -- JSON数组 ["Go", "架构"]
    word_count INTEGER DEFAULT 0,        -- 字数统计
    created_at INTEGER,                  -- Unix timestamp (秒)
    updated_at INTEGER,                  -- Unix timestamp (秒)
    last_opened_at INTEGER               -- Unix timestamp，用于"最近打开"列表
);

-- ============================================
-- 索引
-- ============================================
-- UUID 查询
CREATE INDEX IF NOT EXISTS idx_articles_uuid ON articles(uuid);

-- 文件路径查询
CREATE INDEX IF NOT EXISTS idx_articles_path ON articles(file_path);

-- 最近打开列表排序
CREATE INDEX IF NOT EXISTS idx_articles_last_opened ON articles(last_opened_at DESC);

-- 标题搜索（可选，如果需要简单搜索功能）
CREATE INDEX IF NOT EXISTS idx_articles_title ON articles(title);

-- ============================================
-- 应用配置表
-- ============================================
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- ============================================
-- 默认配置
-- ============================================
INSERT OR IGNORE INTO settings (key, value) VALUES
('deepseek_api_key', ''),            -- DeepSeek API密钥
('editor_theme', 'github-dark'),     -- 编辑器主题
('font_size', '16'),                 -- 字体大小
('ai_model', 'deepseek-chat'),       -- AI模型
('window_width', '1400'),            -- 窗口宽度
('window_height', '900'),            -- 窗口高度
('sidebar_visible', 'true'),         -- 侧边栏是否可见
('auto_save_interval', '2000');      -- 自动保存间隔(毫秒)

-- ============================================
-- 可选：AI使用记录表（用于统计和限制）
-- ============================================
CREATE TABLE IF NOT EXISTS ai_usage (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    article_uuid TEXT NOT NULL,          -- 关联的文章
    action TEXT NOT NULL,                -- 'rewrite', 'polish', 'continue', etc.
    prompt_tokens INTEGER DEFAULT 0,     -- 输入token数
    completion_tokens INTEGER DEFAULT 0, -- 输出token数
    created_at INTEGER DEFAULT (strftime('%s', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_ai_usage_article ON ai_usage(article_uuid);
CREATE INDEX IF NOT EXISTS idx_ai_usage_created ON ai_usage(created_at DESC);

-- ============================================
-- 视图：简化查询
-- ============================================
-- 最近文章视图
CREATE VIEW IF NOT EXISTS v_recent_articles AS
SELECT
    uuid,
    file_path,
    title,
    summary,
    word_count,
    datetime(last_opened_at, 'unixepoch', 'localtime') as last_opened_formatted,
    last_opened_at
FROM articles
WHERE last_opened_at IS NOT NULL
ORDER BY last_opened_at DESC;

-- ============================================
-- 触发器：自动更新时间戳
-- ============================================
-- 更新文章时自动设置 updated_at
CREATE TRIGGER IF NOT EXISTS tr_articles_update_timestamp
AFTER UPDATE ON articles
FOR EACH ROW
BEGIN
    UPDATE articles SET updated_at = strftime('%s', 'now')
    WHERE id = NEW.id AND OLD.updated_at = NEW.updated_at;
END;

-- ============================================
-- 迁移记录表（用于未来版本升级）
-- ============================================
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- 标记初始版本
INSERT OR IGNORE INTO schema_migrations (version) VALUES (1);
