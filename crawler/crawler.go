package main

import (
	"crawler/service"
	"crawler/utils"
	"log/slog"
	"strings"
)

func getDomainURL(url string) string {
	// "https://techblog.lycorp.co.jp/ko/migrate-mysql-with-read-only-mode"
	// ["https:" "" "techblog.lycorp.co.jp" "ko" "migrate-mysql-with-read-only-mode"]
	return strings.Split(url, "/")[2]
}

// Init the crawler
func (r *Crawler) Init(stub *service.ResultInfoClient, config *Config) {
	logger := utils.SlogLogger.With(
		slog.Uint64("crawlerId", r.ID),
	)
	logger.Info("starting Init")

	i := getCrawlerIdx(*config, r.ID)
	// 이미 domain table에 있는 경우
	if config.Crawlers[i].Initialized {
		logger.Info("Init Done")
		return
	}
	domainURL := getDomainURL(r.RSS.URL)
	// grpc unary
	message := insertDomain(stub, r.ID, domainURL)
	utils.LogInit(message, r.Name)
	config.Crawlers[i].Initialized = true
}

// Run the crawler
func (r *Crawler) Run(stub *service.ResultInfoClient, config *Config) {
	logger := utils.SlogLogger.With(
		slog.Uint64("crawlerId", r.ID),
	)
	var postNumToUpdate int32 = 0
	var postNumUpdated uint32 = 0
	var posts []utils.Post = utils.GetParsedData(r.RSS.URL)
	domainURL := getDomainURL(r.RSS.URL)
	var lastIdxToUpdate int32 = utils.CheckUpdatedPost(posts, r.ID, domainURL, r.RSS.LastUpdated)
	if lastIdxToUpdate < 0 {
		utils.LogRun(r.Name, postNumToUpdate, postNumUpdated)
		return
	}
	postNumToUpdate = lastIdxToUpdate + 1

	// grpc client streaming
	logger.Info("starting Update")
	postNumUpdated = insertPost_clientStreaming(stub, &posts, lastIdxToUpdate)

	updateConfig(config, r.ID)
	utils.LogRun(r.Name, postNumToUpdate, postNumUpdated)
}
