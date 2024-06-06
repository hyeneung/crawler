package test

import (
	"sync"
	"testing"
	"time"
)

type crawlerInfo struct {
	name   string
	rssURL string
}
type rss struct {
	url         string
	lastUpdated int64
}

// RSSCrawler type
type rssCrawler struct {
	id   uint64
	name string
	rss  rss
}

// Constructor of RSSCrawler
func New(_id uint64, _name string, _url string) *rssCrawler {
	return &rssCrawler{
		id:   _id,
		name: _name,
		rss: rss{
			url:         _url,
			lastUpdated: time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC).Unix(),
			// lastUpdated: time.Now().Unix()
		},
	}
}
func (r *rssCrawler) Run(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Logf("%d-%d", r.id, i)
	}
}
func TestGoroutine(t *testing.T) {

	const lineURL = "https://techblog.lycorp.co.jp/ko/feed/index.xml"
	const mercariURL = "https://engineering.mercari.com/en/blog/feed.xml"
	const kurlyURL = "https://helloworld.kurly.com/feed.xml"
	var crawlerInfos = []crawlerInfo{
		{name: "lineCrawler", rssURL: lineURL},
		{name: "mercariCrawler", rssURL: mercariURL},
		{name: "kurleyCrawler", rssURL: kurlyURL},
	}

	var wg sync.WaitGroup
	for i := 0; i < len(crawlerInfos); i++ {
		wg.Add(1)
		go func(id uint64) {
			defer wg.Done()
			crawler := crawlerInfos[id]
			rssCrawler := New(id, crawler.name, crawler.rssURL)
			t.Logf("%d", rssCrawler.id)
			rssCrawler.Run(t)
		}(uint64(i))
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
