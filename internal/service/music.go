package service

import (
	"fmt"

	"github.com/user/fish-music/internal/database"
	"github.com/user/fish-music/internal/model"
	"github.com/user/fish-music/pkg/api"
)

// MusicService 音乐服务
type MusicService struct {
	searchClient *api.NeteaseAPI
	songRepo     *database.SongRepository
}

// NewMusicService 创建音乐服务
func NewMusicService(
	searchClient *api.NeteaseAPI,
	songRepo *database.SongRepository,
) *MusicService {
	return &MusicService{
		searchClient: searchClient,
		songRepo:     songRepo,
	}
}

// ProcessDownload 处理下载任务（完整流程）
func (s *MusicService) ProcessDownload(searchResult api.SongInfo) (*model.Song, error) {
	// 1. 检查是否已存在
	uniqueHash := api.GenerateUniqueHash("", searchResult.Name, api.FormatArtist(searchResult.Artists))
	existingSong, err := s.songRepo.FindByUniqueHash(uniqueHash)
	if err == nil && existingSong != nil {
		return existingSong, nil // 已存在，直接返回
	}

	// TODO: 实现下载流程
	return nil, fmt.Errorf("自动下载功能暂不可用，请手动添加")
}

// SearchMusic 搜索音乐
func (s *MusicService) SearchMusic(keyword string) ([]*model.Song, []api.SongInfo, error) {
	// 1. 先从数据库搜索
	dbSongs, err := s.songRepo.Search(keyword, 10)
	if err == nil && len(dbSongs) > 0 {
		return dbSongs, nil, nil
	}

	// 2. 数据库无结果，API 也不稳定
	return nil, nil, fmt.Errorf("数据库无结果，请先添加音乐")
}

// ReprocessMissingSong 重新处理缺失的歌曲
func (s *MusicService) ReprocessMissingSong(songID uint) error {
	// TODO: 实现重新处理逻辑
	return fmt.Errorf("补档功能开发中")
}

// GetSongByID 根据 ID 获取歌曲
func (s *MusicService) GetSongByID(id uint) (*model.Song, error) {
	var song model.Song
	err := database.DB.Where("id = ?", id).First(&song).Error
	if err != nil {
		return nil, err
	}
	return &song, nil
}

// SearchAndDownload 搜索并下载（异步任务）
type SearchAndDownloadTask struct {
	musicSvc *MusicService
	keyword  string
	chatID   int64
}

func (t *SearchAndDownloadTask) Execute() error {
	_, apiResults, err := t.musicSvc.SearchMusic(t.keyword)
	if err != nil {
		return err
	}

	if len(apiResults) > 0 {
		// 下载第一个结果
		_, err := t.musicSvc.ProcessDownload(apiResults[0])
		return err
	}

	return fmt.Errorf("未找到相关歌曲")
}
