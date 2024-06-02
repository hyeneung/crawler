package main

import (
	"log/slog"
	"os"
	"time"
)

type crawlerInfo struct {
	name   string
	rssURL string
}

func logDomain(crawlerName string) {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}).WithAttrs([]slog.Attr{
		slog.String("crawler", crawlerName),
	})
	logger := slog.New(logHandler)
	logger.Info("Init Done")
}

// log the result
func logPost(crawlerName string, totalCount int8, successCount uint8) {
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

func main() {
	const lineURL = "https://techblog.lycorp.co.jp/ko/feed/index.xml"
	const mercariURL = "https://engineering.mercari.com/en/blog/feed.xml"
	var crawlerInfos = []crawlerInfo{
		{name: "lineCrawler", rssURL: lineURL},
		{name: "mercariCrawler", rssURL: mercariURL},
	}

	var id uint64 = 0
	for _, crawlerInfo := range crawlerInfos {
		// goroutine 적용
		rssCrawler := New(id, crawlerInfo.name, crawlerInfo.rssURL)
		rssCrawler.Init() // domain DB에 crawler id, domain url 저장
		logDomain(rssCrawler.name)
		totalCount, successCount := rssCrawler.Run(time.Now().Unix()) // post DB에 게시물 정보 저장
		logPost(rssCrawler.name, totalCount, successCount)
		id += 1
	}
}
