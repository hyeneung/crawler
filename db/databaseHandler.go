package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	Title   string
	Link    string
	PubDate string
}

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
	checkErr(err)
	// conn.SetConnMaxLifetime(time.Minute * 3)
	// conn.SetMaxOpenConns(10)
	// conn.SetMaxIdleConns(10)
	return conn
}

func insertDomainDB(crawlerID uint64, domainURL string) {
	conn := getConnection()
	defer conn.Close()
	_, err := conn.Exec("INSERT INTO domain (id, url) VALUES (?, ?)", crawlerID, domainURL)
	checkErr(err)
}

func insertPostDB(calledCnt uint8, crawlerID uint64, url string, title string, date uint64) uint8 {
	conn := getConnection()
	// defer conn.Close() // resource pool 이용.반환
	stmt, err := conn.Prepare("INSERT INTO post (id, url, title, date) VALUES (?, ?, ?, ?)")
	checkErr(err)
	defer stmt.Close()
	var updatedCount uint8 = 0
	for _, post := range posts {
		// goroutine 적용 - updatedCount 동기화 문제 발생가능. 처리방법 확인
		_, err := stmt.Exec(crawlerID, post.Link, post.Title, Str2UnixTime(post.PubDate))
		checkDBInsertErr(err)
		if err == nil {
			updatedCount += 1
		}
	}
	return updatedCount
}
