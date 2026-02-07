package model

import (
	"fmt"
	"time"
)

// Song 歌曲模型
type Song struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UniqueHash  string    `gorm:"uniqueIndex;size:64;not null" json:"unique_hash"`      // 文件指纹，防止重复
	FileID      string    `gorm:"size:255;not null" json:"file_id"`                       // Telegram File ID
	SourceURL   string    `gorm:"size:512;not null" json:"source_url"`                    // 源链接，用于补档

	// 元数据
	Title       string    `gorm:"size:255;not null" json:"title"`                         // 歌曲标题
	Artist      string    `gorm:"size:255;not null" json:"artist"`                        // 歌手名称
	Album       string    `gorm:"size:255" json:"album"`                                  // 专辑名称
	Duration    int       `json:"duration"`                                               // 时长（秒）
	FileSize    int64     `json:"file_size"`                                              // 文件大小（字节）

	// 扩展元数据 (JSON)
	CountryCode string    `gorm:"size:10" json:"country_code"`                            // 国家代码 (CN, JP, US 等)
	Year        int       `json:"year"`                                                   // 发行年份
	CoverURL    string    `gorm:"size:512" json:"cover_url"`                              // 封面图片 URL
	Lyrics      string    `gorm:"type:text" json:"lyrics"`                                 // 歌词内容
	Genre       string    `gorm:"size:50" json:"genre"`                                    // 歌曲类型
	Language    string    `gorm:"size:50" json:"language"`                                  // 歌曲语言

	// 状态
	IsMissing   bool      `gorm:"default:false" json:"is_missing"`                        // 是否需要补档
	Status      string    `gorm:"size:20;default:active" json:"status"`                    // 状态: active, missing, processing

	// 时间戳
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Song) TableName() string {
	return "songs"
}

// SongMetadata 歌曲元数据结构（用于 JSON 序列化）
type SongMetadata struct {
	CountryCode string `json:"country_code"`
	Year        int    `json:"year"`
	CoverURL    string `json:"cover_url"`
	Lyrics      string `json:"lyrics"`
}

// GetCountryEmoji 获取国家 Emoji
func (s *Song) GetCountryEmoji() string {
	emojiMap := map[string]string{
		"CN": "🇨🇳",
		"JP": "🇯🇵",
		"US": "🇺🇸",
		"UK": "🇬🇧",
		"KR": "🇰🇷",
		"DE": "🇩🇪",
		"FR": "🇫🇷",
		"IT": "🇮🇹",
		"ES": "🇪🇸",
		"RU": "🇷🇺",
		"CA": "🇨🇦",
		"AU": "🇦🇺",
		"BR": "🇧🇷",
		"MX": "🇲🇽",
		"IN": "🇮🇳",
		"TW": "🇹🇼",
		"HK": "🇭🇰",
		"SG": "🇸🇬",
		"MY": "🇲🇾",
		"TH": "🇹🇭",
		"VN": "🇻🇳",
		"ID": "🇮🇩",
		"PH": "🇵🇭",
	}

	if emoji, ok := emojiMap[s.CountryCode]; ok {
		return emoji
	}
	return "🌍"
}

// GetYearText 获取年份文本（中文格式）
func (s *Song) GetYearText() string {
	if s.Year > 0 {
		return fmt.Sprintf("%d年", s.Year)
	}
	return "未知"
}

// GetGenreText 获取类型文本
func (s *Song) GetGenreText() string {
	if s.Genre != "" {
		return s.Genre
	}
	return "未分类"
}

// GetLanguageText 获取语言文本
func (s *Song) GetLanguageText() string {
	if s.Language != "" {
		return s.Language
	}
	return "未知"
}

// UpdateCountryCodeByLanguage 根据语言更新国家代码
func (s *Song) UpdateCountryCodeByLanguage() {
	// 语言到国家的映射
	langToCountry := map[string]string{
		"华语": "CN", // 中文 → 中国
		"英语": "US", // 英文 → 美国
		"日语": "JP", // 日文 → 日本
		"韩语": "KR", // 韩文 → 韩国
		"法语": "FR", // 法文 → 法国
		"德语": "DE", // 德文 → 德国
		"西班牙语": "ES", // 西班牙文 → 西班牙
		"俄语": "RU", // 俄文 → 俄罗斯
		"意大利语": "IT", // 意大利文 → 意大利
		"葡萄牙语": "BR", // 葡萄牙文 → 巴西
		"泰语": "TH", // 泰文 → 泰国
		"越南语": "VN", // 越南文 → 越南
		"印尼语": "ID", // 印尼文 → 印尼
		"马来语": "MY", // 马来文 → 马来西亚
		" Hindi": "IN", // 印地语 → 印度
		"Tagalog": "PH", // 菲律宾文 → 菲律宾
		"其他": "US", // 其他 → 默认美国
	}

	if countryCode, ok := langToCountry[s.Language]; ok {
		s.CountryCode = countryCode
	}
}
