package main

import (
	"crawler/utils"
	"log/slog"
	"os"
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
func New(_id uint64, _name, _url string) *RSSCrawler {
	return &RSSCrawler{
		id:   _id,
		name: _name,
		// rss:  rss{url: _url, lastUpdated: time.Now().Unix()},
		rss: rss{url: _url, lastUpdated: time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC).Unix()},
	}
}

// Run the crawler
func (r *RSSCrawler) Run(currentTime int64) {
	totalCount, successCount := utils.UpdateDB(r.id, r.rss.url, r.rss.lastUpdated)
	r.rss.lastUpdated = currentTime

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}).WithAttrs([]slog.Attr{
		slog.String("crawler", r.name),
		slog.Group("results",
			"newly posted", totalCount,
			"updated", successCount,
		),
	})
	logger := slog.New(logHandler)
	logger.Debug("executed")
}
