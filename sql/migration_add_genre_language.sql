-- Fish Music Database Migration
-- 添加类型和语言字段
-- 版本: v1.1
-- 创建日期: 2026-02-08

-- 添加类型和语言字段到 songs 表
ALTER TABLE songs ADD COLUMN IF NOT EXISTS genre VARCHAR(50);
ALTER TABLE songs ADD COLUMN IF NOT EXISTS language VARCHAR(50);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_songs_genre ON songs(genre);
CREATE INDEX IF NOT EXISTS idx_songs_language ON songs(language);

-- 添加注释
COMMENT ON COLUMN songs.genre IS '歌曲类型: 流行, 摇滚, 说唱, 民谣, 电子, 古典, 爵士, 其他';
COMMENT ON COLUMN songs.language IS '歌曲语言: 华语, 英语, 日语, 韩语, 其他';

-- 更新现有数据（根据国家代码推断语言）
UPDATE songs SET language = '华语' WHERE country_code = 'CN' AND language IS NULL;
UPDATE songs SET language = '英语' WHERE country_code = 'US' AND language IS NULL;
UPDATE songs SET language = '日语' WHERE country_code = 'JP' AND language IS NULL;
UPDATE songs SET language = '韩语' WHERE country_code = 'KR' AND language IS NULL;

-- 如果没有匹配的语言，设置为其他
UPDATE songs SET language = '其他' WHERE language IS NULL;
