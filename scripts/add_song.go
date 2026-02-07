package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// 连接数据库
	connStr := "host=localhost port=5432 user=fish_music password=fish_music_pass dbname=fish_music sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatal("无法连接数据库:", err)
	}

	// 示例：添加一首歌
	addSong(db)
}

func addSong(db *sql.DB) {
	fmt.Println("=== 添加音乐到 Fish Music ===")
	fmt.Println()

	var title, artist, album, fileID string

	fmt.Print("歌曲标题: ")
	fmt.Scanln(&title)

	fmt.Print("歌手名称: ")
	fmt.Scanln(&artist)

	fmt.Print("专辑名称 (可选，直接回车跳过): ")
	fmt.Scanln(&album)

	fmt.Print("Telegram File ID: ")
	fmt.Scanln(&fileID)

	if title == "" || artist == "" || fileID == "" {
		fmt.Println("❌ 标题、歌手和 File ID 不能为空")
		return
	}

	// 生成简单的唯一哈希（可以用更复杂的方法）
	uniqueHash := fmt.Sprintf("%d_%s", time.Now().Unix(), title)

	// 插入数据库
	query := `
		INSERT INTO songs (unique_hash, file_id, source_url, title, artist, album, duration, file_size, country_code, year, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := db.Exec(query,
		uniqueHash,
		fileID,
		"", // source_url 可以为空
		title,
		artist,
		album,
		0,    // duration 未知
		0,    // file_size 未知
		"",   // country_code 未知
		0,    // year 未知
		"active",
	)

	if err != nil {
		fmt.Println("❌ 添加失败:", err)
		return
	}

	fmt.Println()
	fmt.Println("✅ 歌曲添加成功!")
	fmt.Printf("   %s - %s\n", artist, title)
}
