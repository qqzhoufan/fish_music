package api

import (
	"fmt"
	"strings"
	"time"
)

// SongMetadata 歌曲元数据
type SongMetadata struct {
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	Year        int    `json:"year"`
	CoverURL    string `json:"cover_url"`
	Lyrics      string `json:"lyrics"`
	CountryCode string `json:"country_code"`
	CoverPath   string `json:"cover_path"`
}

// SearchResult 搜索结果（通用格式）
type SearchResult struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Duration int    `json:"duration"`
	URL      string `json:"url"`
	CoverURL string `json:"cover_url"`
	Year     int    `json:"year"`
}

// detectRegionFromArtistName 从歌手名称检测地区
func detectRegionFromArtistName(name string) string {
	// 简单的字符编码检测
	if containsCJK(name) {
		// 进一步区分中日韩
		if containsJapanese(name) {
			return "JP"
		}
		if containsKorean(name) {
			return "KR"
		}
		return "CN"
	}

	// 默认返回欧美
	return "US"
}

// containsCJK 检测是否包含中日韩文字
func containsCJK(s string) bool {
	for _, r := range s {
		if (r >= 0x4E00 && r <= 0x9FFF) || // CJK Unified Ideographs
			(r >= 0xAC00 && r <= 0xD7A3) || // Hangul Syllables
			(r >= 0x3040 && r <= 0x309F) { // Hiragana
			return true
		}
	}
	return false
}

// containsJapanese 检测是否包含日文
func containsJapanese(s string) bool {
	for _, r := range s {
		if (r >= 0x3040 && r <= 0x309F) || // Hiragana
			(r >= 0x30A0 && r <= 0x30FF) { // Katakana
			return true
		}
	}
	return false
}

// containsKorean 检测是否包含韩文
func containsKorean(s string) bool {
	for _, r := range s {
		if r >= 0xAC00 && r <= 0xD7A3 { // Hangul Syllables
			return true
		}
	}
	return false
}

// detectCountryCode 检测国家代码（替代之前的实现）
func detectCountryCode(artist string) string {
	code := detectRegionFromArtistName(artist)
	return code
}

// FormatLyrics 格式化歌词
func FormatLyrics(lyrics string, maxLength int) string {
	if len(lyrics) <= maxLength {
		return lyrics
	}

	// 截断歌词并添加提示
	return strings.TrimSpace(lyrics[:maxLength]) + "\n\n... (歌词过长，已截断)"
}

// GenerateUniqueHash 生成文件唯一哈希
func GenerateUniqueHash(url string, title string, artist string) string {
	return fmt.Sprintf("%d_%s_%s", time.Now().Unix(), title, artist)
}
