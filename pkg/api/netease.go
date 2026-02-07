package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// NeteaseAPI 网易云音乐第三方 API 客户端
// 使用开源项目: https://github.com/Binaryify/NeteaseCloudMusicApi
type NeteaseAPI struct {
	baseURL    string
	httpClient *http.Client
}

// NewNeteaseAPI 创建网易云 API 客户端
func NewNeteaseAPI(baseURL string) *NeteaseAPI {
	if baseURL == "" {
		// 使用多个可用的公开 API 服务
		baseURL = "https://music-api-xuanwuserver.vercel.app"
	}

	return &NeteaseAPI{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Artist 歌手信息
type Artist struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// SongInfo 歌曲信息
type SongInfo struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Artists  []Artist `json:"artists"`
	Album    struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		PicURL string `json:"picUrl"`
	} `json:"album"`
	Duration int    `json:"duration"`
	MusicID   string `json:"musicId"`
}

// SongURL 歌曲播放地址
type SongURL struct {
	ID   int64  `json:"id"`
	URL  string `json:"url"`
	Size int    `json:"size"`
}

// LyricResult 歌词结果
type LyricResult struct {
	Lyric string `json:"lyric"`
}

// NeteaseSearchResponse 搜索响应
type NeteaseSearchResponse struct {
	Result struct {
		Songs     []SongInfo `json:"songs"`
		SongCount int        `json:"songCount"`
	} `json:"result"`
	Code int `json:"code"`
}

// NeteaseSongURLResponse 歌曲地址响应
type NeteaseSongURLResponse struct {
	Data []SongURL `json:"data"`
	Code int       `json:"code"`
}

// NeteaseLyricResponse 歌词响应
type NeteaseLyricResponse struct {
	Lrc LyricResult `json:"lrc"`
	Code int         `json:"code"`
}

// Search 搜索音乐
func (n *NeteaseAPI) Search(keyword string, limit int) ([]SongInfo, error) {
	apiURL := fmt.Sprintf("%s/search?keywords=%s&limit=%d",
		n.baseURL, url.QueryEscape(keyword), limit)

	resp, err := n.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result NeteaseSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("API 返回错误: %d", result.Code)
	}

	return result.Result.Songs, nil
}

// GetSongURL 获取歌曲播放地址
func (n *NeteaseAPI) GetSongURL(songID int64) (string, error) {
	apiURL := fmt.Sprintf("%s/song/url?id=%d", n.baseURL, songID)

	resp, err := n.httpClient.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result NeteaseSongURLResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 200 || len(result.Data) == 0 {
		return "", fmt.Errorf("无法获取歌曲地址")
	}

	// 网易云的 URL 可能过期，返回网易云外链
	// 实际下载时可以尝试用这个 URL 或使用其他方法
	return result.Data[0].URL, nil
}

// GetLyric 获取歌词
func (n *NeteaseAPI) GetLyric(songID int64) (string, error) {
	apiURL := fmt.Sprintf("%s/lyric?id=%d", n.baseURL, songID)

	resp, err := n.httpClient.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result NeteaseLyricResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 200 {
		return "", fmt.Errorf("API 返回错误: %d", result.Code)
	}

	return result.Lrc.Lyric, nil
}

// GetSongDetail 获取歌曲详情
func (n *NeteaseAPI) GetSongDetail(songID int64) (*SongInfo, error) {
	apiURL := fmt.Sprintf("%s/song/detail?ids=%d", n.baseURL, songID)

	resp, err := n.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		Songs []SongInfo `json:"songs"`
		Code  int        `json:"code"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 200 || len(result.Songs) == 0 {
		return nil, fmt.Errorf("无法获取歌曲详情")
	}

	return &result.Songs[0], nil
}

// FormatArtist 格式化歌手列表
func FormatArtist(artists []Artist) string {
	names := make([]string, len(artists))
	for i, a := range artists {
		names[i] = a.Name
	}
	return strings.Join(names, ", ")
}

// ParseDuration 解析时长（毫秒转秒）
func ParseDuration(ms int) int {
	return ms / 1000
}
