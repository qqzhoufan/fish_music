package database

import (
	"github.com/user/fish-music/internal/model"
	"gorm.io/gorm"
)

// SongRepository 歌曲数据访问层
type SongRepository struct {
	db *gorm.DB
}

// NewSongRepository 创建歌曲仓库
func NewSongRepository() *SongRepository {
	return &SongRepository{db: DB}
}

// FindByFileID 根据 FileID 查找歌曲
func (r *SongRepository) FindByFileID(fileID string) (*model.Song, error) {
	var song model.Song
	err := r.db.Where("file_id = ? AND status = ?", fileID, "active").First(&song).Error
	if err != nil {
		return nil, err
	}
	return &song, nil
}

// FindByUniqueHash 根据唯一哈希查找歌曲
func (r *SongRepository) FindByUniqueHash(hash string) (*model.Song, error) {
	var song model.Song
	err := r.db.Where("unique_hash = ?", hash).First(&song).Error
	if err != nil {
		return nil, err
	}
	return &song, nil
}

// Search 搜索歌曲（标题或歌手）
func (r *SongRepository) Search(keyword string, limit int) ([]*model.Song, error) {
	var songs []*model.Song
	query := r.db.Where("status = ?", "active")

	if keyword != "" {
		query = query.Where("title ILIKE ? OR artist ILIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Order("created_at DESC").Limit(limit).Find(&songs).Error
	return songs, err
}

// Create 创建歌曲记录
func (r *SongRepository) Create(song *model.Song) error {
	return r.db.Create(song).Error
}

// Update 更新歌曲记录
func (r *SongRepository) Update(song *model.Song) error {
	return r.db.Save(song).Error
}

// UpdateFileID 更新 FileID（用于补档）
func (r *SongRepository) UpdateFileID(id uint, fileID string) error {
	return r.db.Model(&model.Song{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"file_id":    fileID,
			"is_missing": false,
			"status":     "active",
		}).Error
}

// MarkMissing 标记为需要补档
func (r *SongRepository) MarkMissing(id uint) error {
	return r.db.Model(&model.Song{}).
		Where("id = ?", id).
		Update("is_missing", true).Error
}

// GetMissingSongs 获取所有需要补档的歌曲
func (r *SongRepository) GetMissingSongs() ([]*model.Song, error) {
	var songs []*model.Song
	err := r.db.Where("is_missing = ? OR status = ?", true, "missing").Find(&songs).Error
	return songs, err
}

// GetRandom 随机获取一首歌
func (r *SongRepository) GetRandom() (*model.Song, error) {
	var song model.Song
	err := r.db.Where("status = ?", "active").Order("RANDOM()").First(&song).Error
	if err != nil {
		return nil, err
	}
	return &song, nil
}

// GetRandomSongs 随机获取多首歌曲
func (r *SongRepository) GetRandomSongs(limit int) ([]*model.Song, error) {
	var songs []*model.Song
	err := r.db.Where("status = ?", "active").Order("RANDOM()").Limit(limit).Find(&songs).Error
	return songs, err
}

// GetStats 获取统计信息
func (r *SongRepository) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalCount int64
	if err := r.db.Model(&model.Song{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats["total_songs"] = totalCount

	var missingCount int64
	if err := r.db.Model(&model.Song{}).Where("is_missing = ?", true).Count(&missingCount).Error; err != nil {
		return nil, err
	}
	stats["missing_songs"] = missingCount

	var artistCount int64
	if err := r.db.Model(&model.Song{}).Select("COUNT(DISTINCT artist)").Scan(&artistCount).Error; err != nil {
		return nil, err
	}
	stats["total_artists"] = artistCount

	var todayCount int64
	if err := r.db.Model(&model.Song{}).
		Where("DATE(created_at) = CURRENT_DATE").
		Count(&todayCount).Error; err != nil {
		return nil, err
	}
	stats["today_added"] = todayCount

	return stats, nil
}

// ============================================
// UserRepository 用户数据访问层
// ============================================

// UserRepository 用户仓库
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository() *UserRepository {
	return &UserRepository{db: DB}
}

// FindByTelegramID 根据 Telegram ID 查找用户
func (r *UserRepository) FindByTelegramID(telegramID int64) (*model.User, error) {
	var user model.User
	err := r.db.Where("telegram_id = ?", telegramID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindOrCreate 查找或创建用户
func (r *UserRepository) FindOrCreate(telegramID int64, username, firstName, lastName string) (*model.User, error) {
	var user model.User
	err := r.db.Where("telegram_id = ?", telegramID).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		user = model.User{
			TelegramID: telegramID,
			Username:   username,
			FirstName:  firstName,
			LastName:   lastName,
		}
		if err := r.db.Create(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	return &user, err
}

// UpdateLastSeen 更新最后活跃时间
func (r *UserRepository) UpdateLastSeen(userID uint) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("last_seen", "NOW()").Error
}

// ============================================
// FavoriteRepository 收藏数据访问层
// ============================================

// FavoriteRepository 收藏仓库
type FavoriteRepository struct {
	db *gorm.DB
}

// NewFavoriteRepository 创建收藏仓库
func NewFavoriteRepository() *FavoriteRepository {
	return &FavoriteRepository{db: DB}
}

// Add 添加收藏
func (r *FavoriteRepository) Add(userID, songID uint) error {
	favorite := &model.Favorite{
		UserID: userID,
		SongID: songID,
	}
	return r.db.Create(favorite).Error
}

// Remove 取消收藏
func (r *FavoriteRepository) Remove(userID, songID uint) error {
	return r.db.Where("user_id = ? AND song_id = ?", userID, songID).Delete(&model.Favorite{}).Error
}

// IsFavorited 检查是否已收藏
func (r *FavoriteRepository) IsFavorited(userID, songID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Favorite{}).
		Where("user_id = ? AND song_id = ?", userID, songID).
		Count(&count).Error
	return count > 0, err
}

// GetByUser 获取用户收藏列表
func (r *FavoriteRepository) GetByUser(userID uint, limit int) ([]*model.Song, error) {
	var songs []*model.Song
	err := r.db.Table("songs").
		Select("songs.*").
		Joins("INNER JOIN favorites ON favorites.song_id = songs.id").
		Where("favorites.user_id = ?", userID).
		Order("favorites.created_at DESC").
		Limit(limit).
		Find(&songs).Error
	return songs, err
}

// ============================================
// HistoryRepository 历史记录数据访问层
// ============================================

// HistoryRepository 历史记录仓库
type HistoryRepository struct {
	db *gorm.DB
}

// NewHistoryRepository 创建历史记录仓库
func NewHistoryRepository() *HistoryRepository {
	return &HistoryRepository{db: DB}
}

// Add 添加历史记录
func (r *HistoryRepository) Add(userID, songID uint) error {
	history := &model.History{
		UserID: userID,
		SongID: songID,
	}
	return r.db.Create(history).Error
}

// GetByUser 获取用户历史记录
func (r *HistoryRepository) GetByUser(userID uint, limit int) ([]*model.Song, error) {
	var songs []*model.Song
	err := r.db.Table("songs").
		Select("DISTINCT ON (songs.id) songs.*").
		Joins("INNER JOIN history ON history.song_id = songs.id").
		Where("history.user_id = ?", userID).
		Order("songs.id, history.created_at DESC").
		Limit(limit).
		Find(&songs).Error
	return songs, err
}

// GetRecentHistory 获取最近的播放记录（返回完整的 History 记录）
func (r *HistoryRepository) GetRecentHistory(userID uint, limit int) ([]*model.History, error) {
	var histories []*model.History
	err := r.db.Preload("Song").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}
