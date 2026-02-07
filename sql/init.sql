-- Fish Music Database Schema
-- 版本: v1.0
-- 创建日期: 2026-02-07

-- 删除已存在的表（开发环境使用）
-- DROP TABLE IF EXISTS history CASCADE;
-- DROP TABLE IF EXISTS favorites CASCADE;
-- DROP TABLE IF EXISTS songs CASCADE;
-- DROP TABLE IF EXISTS users CASCADE;

-- ============================================
-- 用户表
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    language VARCHAR(10) DEFAULT 'zh',
    is_admin BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id);
CREATE INDEX IF NOT EXISTS idx_users_is_admin ON users(is_admin);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

-- ============================================
-- 歌曲表
-- ============================================
CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    unique_hash VARCHAR(64) UNIQUE NOT NULL,
    file_id VARCHAR(255) NOT NULL,
    source_url VARCHAR(512) NOT NULL,

    -- 元数据
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    album VARCHAR(255),
    duration INTEGER DEFAULT 0,
    file_size BIGINT DEFAULT 0,

    -- 扩展元数据
    country_code VARCHAR(10),
    year INTEGER,
    cover_url VARCHAR(512),
    lyrics TEXT,

    -- 状态
    is_missing BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'active',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_songs_unique_hash ON songs(unique_hash);
CREATE INDEX IF NOT EXISTS idx_songs_file_id ON songs(file_id);
CREATE INDEX IF NOT EXISTS idx_songs_title ON songs(title);
CREATE INDEX IF NOT EXISTS idx_songs_artist ON songs(artist);
CREATE INDEX IF NOT EXISTS idx_songs_status ON songs(status);
CREATE INDEX IF NOT EXISTS idx_songs_is_missing ON songs(is_missing);
CREATE INDEX IF NOT EXISTS idx_songs_created_at ON songs(created_at);

-- ============================================
-- 收藏表
-- ============================================
CREATE TABLE IF NOT EXISTS favorites (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    song_id INTEGER NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, song_id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_song_id ON favorites(song_id);
CREATE INDEX IF NOT EXISTS idx_favorites_created_at ON favorites(created_at);

-- ============================================
-- 历史表
-- ============================================
CREATE TABLE IF NOT EXISTS history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    song_id INTEGER NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_history_user_id ON history(user_id);
CREATE INDEX IF NOT EXISTS idx_history_song_id ON history(song_id);
CREATE INDEX IF NOT EXISTS idx_history_user_created ON history(user_id, created_at DESC);

-- ============================================
-- 触发器: 自动更新 updated_at
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 应用到各表
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_songs_updated_at ON songs;
CREATE TRIGGER update_songs_updated_at
    BEFORE UPDATE ON songs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 视图: 统计信息
-- ============================================
CREATE OR REPLACE VIEW song_stats AS
SELECT
    COUNT(*) as total_songs,
    COUNT(DISTINCT artist) as total_artists,
    COUNT(CASE WHEN is_missing = TRUE THEN 1 END) as missing_songs,
    SUM(file_size) as total_storage_bytes,
    AVG(duration) as avg_duration
FROM songs;

-- ============================================
-- 视图: 热门歌曲
-- ============================================
CREATE OR REPLACE VIEW popular_songs AS
SELECT
    s.id,
    s.title,
    s.artist,
    s.album,
    COUNT(h.id) as play_count
FROM songs s
LEFT JOIN history h ON s.id = h.song_id
GROUP BY s.id, s.title, s.artist, s.album
ORDER BY play_count DESC
LIMIT 100;

-- ============================================
-- 插入默认管理员（可选）
-- ============================================
-- 注意: 需要手动设置正确的 telegram_id
-- INSERT INTO users (telegram_id, username, is_admin, is_active)
-- VALUES (123456789, 'admin', TRUE, TRUE)
-- ON CONFLICT (telegram_id) DO NOTHING;

-- ============================================
-- 数据清理函数（可选）
-- ============================================
-- 清理 90 天前的历史记录
CREATE OR REPLACE FUNCTION clean_old_history()
RETURNS void AS $$
BEGIN
    DELETE FROM history
    WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '90 days';
END;
$$ LANGUAGE plpgsql;
