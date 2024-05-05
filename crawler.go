package main

import (
	"crawler/utils"
	"fmt"
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
	name string
	rss  rss
}

// Constructor of RSSCrawler
func New(name, url string) *RSSCrawler {
	return &RSSCrawler{
		name: name,
		// rss:  rss{url: url, lastUpdated: time.Now().Unix()},
		rss: rss{url: url, lastUpdated: time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC).Unix()},
	}
}

// Run the crawler
func (r *RSSCrawler) Run(currentTime int64) {
	totalCount, successCount := utils.UpdateDB(r.rss.url, r.rss.lastUpdated)
	r.rss.lastUpdated = currentTime
	fmt.Printf("[%s] %d out of %d successfully updated at %s\n", r.name, totalCount, successCount, utils.UnixTime2Time(r.rss.lastUpdated))
}
