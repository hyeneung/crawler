package main

import (
	"time"
)

type crawlerInfo struct {
	name   string
	rssURL string
}

func main() {
	const lineURL = "https://techblog.lycorp.co.jp/ko/feed/index.xml"
	const mercariURL = "https://engineering.mercari.com/en/blog/feed.xml"
	var crawlerInfos = []crawlerInfo{
		{name: "lineCrawler", rssURL: lineURL},
		{name: "mercariCrawler", rssURL: mercariURL},
	}

	var id uint64 = 0
	for _, crawlerInfo := range crawlerInfos {
		// TODO - goroutine 적용
		rssCrawler := New(id, crawlerInfo.name, crawlerInfo.rssURL)
		rssCrawler.Init()                 // domain DB에 crawler id, domain url 저장
		rssCrawler.Run(time.Now().Unix()) // post DB에 게시물 정보 저장
		id += 1
	}
}
