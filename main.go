package main

import "time"

type crawler struct {
	name   string
	rssURL string
}

func RunCrawlers(crawlers []*RSSCrawler) {
	// goroutine 적용
	for _, crawlerInstance := range crawlers {
		crawlerInstance.Run(time.Now().Unix())
	}
}

func main() {
	const lineURL = "https://techblog.lycorp.co.jp/ko/feed/index.xml"
	const mercariURL = "https://engineering.mercari.com/en/blog/feed.xml"
	var crawlerInfos = []crawler{
		{"lineCrawler", lineURL},
		{"mercariCrawler", mercariURL},
	}

	var crawlers []*RSSCrawler
	var id uint64 = 0
	for _, crawlerInfo := range crawlerInfos {
		crawlerInstance := New(id, crawlerInfo.name, crawlerInfo.rssURL)
		crawlers = append(crawlers, crawlerInstance)
		id += 1
	}
	RunCrawlers(crawlers)
}
