package main

import (
	"crawler/utils"
	"log/slog"
	"os"
	"strings"
	"time"
)

// rss type
// rss 2.0 standards : https://www.rssboard.org/rss-specification
type rss struct {
	url         string
	lastUpdated int64
}

// RSSCrawler type
type RSSCrawler struct {
	id   uint64
	name string
	rss  rss
}

// Constructor of RSSCrawler
func New(_id uint64, _name string, _url string) *RSSCrawler {
	return &RSSCrawler{
		id:   _id,
		name: _name,
		rss: rss{
			url:         _url,
			lastUpdated: time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC).Unix(),
			// lastUpdated: time.Now().Unix()
		},
	}
}
func logInit(crawlerName string) {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}).WithAttrs([]slog.Attr{
		slog.String("crawler", crawlerName),
	})
	logger := slog.New(logHandler)
	logger.Info("Init Done")
}

// Init the crawler
func (r *RSSCrawler) Init() {
	// "https://techblog.lycorp.co.jp/ko/migrate-mysql-with-read-only-mode"
	// ["https:" "" "techblog.lycorp.co.jp" "ko" "migrate-mysql-with-read-only-mode"]
	// domainURL := strings.Split(r.rss.url, "/")[2]

	// // TODO - grpc unary
	// err := db.InsertDomain(r.id, domainURL)

	// utils.CheckDBInsertErr(err)
	// logInit(r.name)
}

// log the result
func logRun(crawlerName string, totalCount int8, successCount uint8) {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}).WithAttrs([]slog.Attr{
		slog.String("crawler", crawlerName),
		slog.Group("results",
			"newly posted", totalCount,
			"updated", successCount,
		),
	})
	logger := slog.New(logHandler)
	logger.Info("Done")
}

// Run the crawler
func (r *RSSCrawler) Run(currentTime int64) {
	var postNumToUpdate int8 = 0
	var postNumUpdated uint8 = 0
	posts := utils.GetParsedData(r.rss.url)
	domainURL := strings.Split(r.rss.url, "/")[2]

	lastIdxToUpdate := utils.CheckUpdatedPost(posts, domainURL, r.rss.lastUpdated)
	if lastIdxToUpdate < 0 {
		logRun(r.name, postNumToUpdate, postNumUpdated)
		return
	}
	// var postsToUpdate []utils.Post = posts[:lastIdxToUpdate+1]

	// // TODO - grpc client streaming
	// var successCount uint8 = db.InsertPosts_(postsToUpdate)
	// // TODO - grpc server streaming
	// var logs []string = db.GetLogs(r.id)
	// // TODO - grpc biddirectional streaming
	// var logs []string = db.InsertPosts(postsToUpdate)

	r.rss.lastUpdated = currentTime
	logRun(r.name, postNumToUpdate, postNumUpdated)
}
