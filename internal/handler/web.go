package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/user/fish-music/internal/database"
	"github.com/user/fish-music/internal/model"
)

// WebHandler Web 处理器
type WebHandler struct {
	username  string
	password  string
	songRepo  *database.SongRepository
}

// NewWebHandler 创建 Web 处理器
func NewWebHandler(
	username, password string,
	songRepo *database.SongRepository,
) *WebHandler {
	return &WebHandler{
		username: username,
		password: password,
		songRepo: songRepo,
	}
}

// RegisterRoutes 注册路由
func (h *WebHandler) RegisterRoutes(router *gin.Engine) {
	// 公开路由
	router.GET("/", h.basicAuth(h.handleIndex))
	router.GET("/api/stats", h.basicAuth(h.apiStats))

	// 歌曲管理
	router.GET("/api/songs", h.basicAuth(h.apiListSongs))
	router.GET("/api/songs/:id", h.basicAuth(h.apiGetSong))
	router.GET("/api/songs/missing", h.basicAuth(h.apiMissingSongs))
	router.POST("/api/songs/:id/reprocess", h.basicAuth(h.apiReprocessSong))
	router.PUT("/api/songs/:id", h.basicAuth(h.apiUpdateSong))
	router.DELETE("/api/songs/:id", h.basicAuth(h.apiDeleteSong))
}

// basicAuth Basic Auth 中间件
func (h *WebHandler) basicAuth(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok || username != h.username || password != h.password {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		handler(c)
	}
}

// handleIndex 首页
func (h *WebHandler) handleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Fish Music 管理后台",
	})
}

// apiStats 统计信息 API
func (h *WebHandler) apiStats(c *gin.Context) {
	stats, err := h.songRepo.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// apiListSongs 歌曲列表 API
func (h *WebHandler) apiListSongs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	keyword := c.Query("q")
	genre := c.Query("genre")
	language := c.Query("language")

	offset := (page - 1) * limit

	var songs []model.Song
	query := database.DB.Model(&model.Song{})

	if keyword != "" {
		query = query.Where("title ILIKE ? OR artist ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if genre != "" {
		query = query.Where("genre = ?", genre)
	}

	if language != "" {
		query = query.Where("language = ?", language)
	}

	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&songs)

	var total int64
	query.Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"songs": songs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// apiGetSong 获取单个歌曲信息
func (h *WebHandler) apiGetSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的歌曲ID"})
		return
	}

	var song model.Song
	if err := database.DB.Where("id = ?", id).First(&song).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "歌曲不存在"})
		return
	}

	c.JSON(http.StatusOK, song)
}

// apiMissingSongs 获取缺失歌曲列表
func (h *WebHandler) apiMissingSongs(c *gin.Context) {
	songs, err := h.songRepo.GetMissingSongs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"songs": songs,
		"count": len(songs),
	})
}

// apiReprocessSong 重新处理歌曲
func (h *WebHandler) apiReprocessSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的歌曲ID"})
		return
	}

	// TODO: 实现重新处理逻辑
	// 这里需要调用 MusicService.ReprocessMissingSong

	c.JSON(http.StatusOK, gin.H{
		"message": "重新处理任务已提交",
		"id":      id,
	})
}

// apiUpdateSong 更新歌曲信息
func (h *WebHandler) apiUpdateSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的歌曲ID"})
		return
	}

	var input struct {
		Title       string `json:"title"`
		Artist      string `json:"artist"`
		Album       string `json:"album"`
		Genre       string `json:"genre"`
		Language    string `json:"language"`
		CountryCode string `json:"country_code"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Artist != "" {
		updates["artist"] = input.Artist
	}
	if input.Album != "" {
		updates["album"] = input.Album
	}

	// genre、language 和 country_code 字段允许设置为空字符串（清空）
	updates["genre"] = input.Genre
	updates["language"] = input.Language
	updates["country_code"] = input.CountryCode

	if len(updates) > 0 {
		// 如果用户没有手动设置 country_code，但设置了 language，则自动根据语言设置国家代码
		if input.CountryCode == "" && input.Language != "" {
			var song model.Song
			if err := database.DB.Where("id = ?", id).First(&song).Error; err == nil {
				song.Language = input.Language
				song.UpdateCountryCodeByLanguage()
				updates["country_code"] = song.CountryCode
			}
		}

		if err := database.DB.Model(&model.Song{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// apiDeleteSong 删除歌曲
func (h *WebHandler) apiDeleteSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的歌曲ID"})
		return
	}

	if err := database.DB.Delete(&model.Song{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
