package utils

import (
	"database/sql"
	"db/service"
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
	info := connectionInfo{username: "root", password: "1234", host: "127.0.0.1", port: 3306, database: "crawl_data"}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", info.username, info.password, info.host, info.port, info.database)
	conn, err := sql.Open("mysql", dsn)
	checkFatalErr(err)
	// conn.SetConnMaxLifetime(time.Minute * 3)
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10) // connection pool
	return conn
}

func InsertDomainDB(crawlerID uint64, domainURL string) error {
	conn := getConnection()
	defer conn.Close() // connection 반환(resource pool 이용)
	_, err := conn.Exec("INSERT INTO domain (id, url) VALUES (?, ?)", crawlerID, domainURL)
	return err
}

func InsertPostDB(post *service.Post) error {
	conn := getConnection()
	defer conn.Close() // connection 반환(resource pool 이용)
	// TODO - connection 하나 받을 때 post 하나씩 넣지 말고 한 번에 여러 개 넣을 것.
	stmt, err := conn.Prepare("INSERT INTO post (id, url, title, date) VALUES (?, ?, ?, ?)")
	checkFatalErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(post.Id, post.Link, post.Title, Str2UnixTime(post.PubDate))

	return err
}
