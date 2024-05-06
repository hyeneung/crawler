package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type connectionInfo struct {
	username string
	password string
	host     string
	port     int
	database string
}

func insertDB(items []Item, crawlerID uint64) uint8 {
	info := connectionInfo{username: "root", password: "1234", host: "127.0.0.1", port: 3306, database: "crawl_data"}
	// build the DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", info.username, info.password, info.host, info.port, info.database)
	// Open the connection
	db, err := sql.Open("mysql", dsn)
	CheckErr(err)
	defer db.Close()
	// db.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)

	stmtInsertPost, err := db.Prepare("INSERT INTO post (id, url, title, date) VALUES (?, ?, ?, ?)")
	CheckErr(err)
	defer stmtInsertPost.Close()
	var updatedCount uint8 = 0
	for _, item := range items {
		// goroutine 적용 - updatedCount 동기화 문제 발생가능. 처리방법 확인
		_, err := stmtInsertPost.Exec(crawlerID, item.Link, item.Title, Str2UnixTime(item.PubDate))
		CheckErr(err)
		if err == nil {
			updatedCount += 1
		}
	}
	return updatedCount
}

func UpdateDB(crawlerID uint64, url string, updatedDate int64) (int8, uint8) {
	posts := getParsedData(url)
	var updatedCount uint8 = 0
	lastIndex := getLastIdxToUpdate(posts, updatedDate)
	var items []Item
	for idx := int8(0); idx <= lastIndex; idx++ {
		items = append(items, posts.Items[idx])
	}
	updatedCount = insertDB(items, crawlerID)
	return lastIndex + 1, updatedCount
}
