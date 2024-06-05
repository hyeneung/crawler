package main

import (
	"crawler/service"
	"crawler/utils"
	"log/slog"
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

func getDomainURL(url string) string {
	// "https://techblog.lycorp.co.jp/ko/migrate-mysql-with-read-only-mode"
	// ["https:" "" "techblog.lycorp.co.jp" "ko" "migrate-mysql-with-read-only-mode"]
	return strings.Split(url, "/")[2]
}

// Init the crawler
func (r *rssCrawler) Init(stub *service.ResultInfoClient) {
	domainURL := getDomainURL(r.rss.url)
	// grpc unary
	slog.Info("starting gRPC unary")
	message := insertDomain(stub, r.id, domainURL)
	utils.LogInit(message, r.name)
}

// Run the crawler
func (r *rssCrawler) Run(stub *service.ResultInfoClient, currentTime int64) {
	var postNumToUpdate int32 = 0
	var postNumUpdated uint32 = 0
	// DB에 새로 넣어야 할 게시물 정보 가져옴
	var posts []utils.Post = utils.GetParsedData(r.rss.url)
	domainURL := getDomainURL(r.rss.url)
	var lastIdxToUpdate int32 = utils.CheckUpdatedPost(posts, r.id, domainURL, r.rss.lastUpdated)
	if lastIdxToUpdate < 0 {
		utils.LogRun(r.name, postNumToUpdate, postNumUpdated)
		return
	}
	postNumToUpdate = lastIdxToUpdate + 1

	// grpc client streaming
	if r.id == 0 {
		slog.Info("starting gRPC client streaming")
		postNumUpdated = insertPost_clientStreaming(stub, &posts, lastIdxToUpdate)
	} else {
		// grpc bidirectional streaming
		slog.Info("starting gRPC bidirectional streaming")
		postNumUpdated = insertPost_bidirectionalStreaming(stub, &posts, lastIdxToUpdate)
	}

	r.rss.lastUpdated = currentTime
	utils.LogRun(r.name, postNumToUpdate, postNumUpdated)
}
