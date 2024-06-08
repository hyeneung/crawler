package utils

import (
	"database/sql"
	"db/service"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type connectionInfo struct {
	username string
	password string
	host     string
	port     string
	database string
}

func getConnection() *sql.DB {
	info := connectionInfo{
		username: os.Getenv("DB_USER"),
		password: os.Getenv("DB_PASS"),
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		database: os.Getenv("DB_NAME"),
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", info.username, info.password, info.host, info.port, info.database)
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
	stmt, err := conn.Prepare("INSERT INTO domain (id, url) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(crawlerID, domainURL)
	return err
}

func InsertPostDB(post *service.Post) error {
	conn := getConnection()
	defer conn.Close() // connection 반환(resource pool 이용)
	stmt, err := conn.Prepare("INSERT INTO post (id, url, title, date) VALUES (?, ?, ?, ?)")

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.Id, post.Link, post.Title, Str2UnixTime(post.PubDate))
	return err
}
