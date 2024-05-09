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

func getConnection() *sql.DB {
	// docker exec -it docker-crawler-1 /bin/bash   // host : "db"
	info := connectionInfo{username: "root", password: "1234", host: "db", port: 3306, database: "crawl_data"}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", info.username, info.password, info.host, info.port, info.database)
	conn, err := sql.Open("mysql", dsn)
	CheckErr(err)
	// conn.SetConnMaxLifetime(time.Minute * 3)
	// conn.SetMaxOpenConns(10)
	// conn.SetMaxIdleConns(10)
	return conn
}

func InsertDomainDB(crawlerID uint64, domainURL string) {
	conn := getConnection()
	defer conn.Close()
	_, err := conn.Exec("INSERT INTO domain (id, url) VALUES (?, ?)", crawlerID, domainURL)
	CheckDBInsertErr(err) // 임시. 실제로는 참조 무결성 위반 가능성 배제 위해 아래꺼 exit되는거 써야 함.
	// CheckErr(err)
}

func insertPostDB(posts []Post, crawlerID uint64) uint8 {
	conn := getConnection()
	defer conn.Close()
	stmt, err := conn.Prepare("INSERT INTO post (id, url, title, date) VALUES (?, ?, ?, ?)")
	CheckErr(err)
	defer stmt.Close()
	var updatedCount uint8 = 0
	for _, post := range posts {
		// goroutine 적용 - updatedCount 동기화 문제 발생가능. 처리방법 확인
		_, err := stmt.Exec(crawlerID, post.Link, post.Title, Str2UnixTime(post.PubDate))
		CheckDBInsertErr(err)
		if err == nil {
			updatedCount += 1
		}
	}
	return updatedCount
}

func InsertPostDB(crawlerID uint64, rssURL string, domainURL string, updatedDate int64) (int8, uint8) {
	posts := getParsedData(rssURL)
	var updatedCount uint8 = 0
	lastIndex := getLastIdxToUpdate(posts, domainURL, updatedDate)
	if lastIndex < 0 {
		return 0, updatedCount
	}
	updatedCount = insertPostDB(posts[:lastIndex+1], crawlerID)
	return lastIndex + 1, updatedCount
}
