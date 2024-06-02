package main

import (
	"crawler/utils"
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

// Init the crawler
func (r *RSSCrawler) Init() {
	// "https://techblog.lycorp.co.jp/ko/migrate-mysql-with-read-only-mode"
	// ["https:" "" "techblog.lycorp.co.jp" "ko" "migrate-mysql-with-read-only-mode"]
	domainURL := strings.Split(r.rss.url, "/")[2]
	utils.InsertDomainDB(r.id, domainURL)
}

// Run the crawler
func (r *RSSCrawler) Run(currentTime int64) (int8, uint8) {
	domainURL := strings.Split(r.rss.url, "/")[2]
	totalCount, successCount := utils.InsertPostDB(r.id, r.rss.url, domainURL, r.rss.lastUpdated)
	r.rss.lastUpdated = currentTime
	return totalCount, successCount
}
