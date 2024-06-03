package main

import (
	"context"
	"crawler/service"
	"crawler/utils"
	"io"
	"log"
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
func logInit(res *service.Response, crawlerName string) {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}).WithAttrs([]slog.Attr{
		slog.String("crawler", crawlerName),
		slog.String("db_log", res.Message),
	})
	logger := slog.New(logHandler)
	logger.Info("Init Done")
}

// Init the crawler
func (r *RSSCrawler) Init(stub *service.ResultInfoClient) {
	// "https://techblog.lycorp.co.jp/ko/migrate-mysql-with-read-only-mode"
	// ["https:" "" "techblog.lycorp.co.jp" "ko" "migrate-mysql-with-read-only-mode"]
	domainURL := strings.Split(r.rss.url, "/")[2]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// grpc unary
	res, err := (*stub).InsertDomain(ctx, &service.UnaryRequest{Id: r.id, Url: domainURL})
	if err != nil {
		slog.Error(err.Error())
	}
	logInit(res, r.name)
}

// log the result
func logRun(crawlerName string, totalCount int32, successCount uint32) {
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
func (r *RSSCrawler) Run(stub *service.ResultInfoClient, currentTime int64) {
	var postNumToUpdate int32 = 0
	var postNumUpdated uint32 = 0
	var posts []utils.Post = utils.GetParsedData(r.rss.url)
	domainURL := strings.Split(r.rss.url, "/")[2]
	var lastIdxToUpdate int32 = utils.CheckUpdatedPost(posts, r.id, domainURL, r.rss.lastUpdated)
	if lastIdxToUpdate < 0 {
		logRun(r.name, postNumToUpdate, postNumUpdated)
		return
	}
	postNumToUpdate = lastIdxToUpdate + 1

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// // grpc client streaming
	// insertStream, err := (*stub).InsertPosts(ctx)
	// utils.CheckErr(err)
	// for i := 0; i < int(lastIdxToUpdate+1); i++ {
	// 	post := posts[i]
	// 	data := service.Post{Id: post.Id, Title: post.Title,
	// 		Link: post.Link, PubDate: post.PubDate}
	// 	err := insertStream.Send(&data)
	// 	utils.CheckErr(err)
	// }
	// res, err := insertStream.CloseAndRecv()
	// postNumUpdated = res.Value
	// if err != nil {
	// 	slog.Error(err.Error())
	// }

	// // TODO - grpc server streaming
	// var logs []string = db.GetLogs(r.id)

	// grpc biddirectional streaming
	insertStream, err := (*stub).InsertPosts_(ctx)
	utils.CheckErr(err)
	channel := make(chan struct{})
	go asncClientBidirectionalRPC(insertStream, channel)
	for i := 0; i < int(lastIdxToUpdate+1); i++ {
		post := posts[i]
		data := service.Post{Id: post.Id, Title: post.Title,
			Link: post.Link, PubDate: post.PubDate}
		err := insertStream.Send(&data)
		utils.CheckErr(err)
	}
	if err := insertStream.CloseSend(); err != nil {
		log.Fatal(err)
	}
	channel <- struct{}{}
	r.rss.lastUpdated = currentTime
	logRun(r.name, postNumToUpdate, postNumUpdated)
}

func asncClientBidirectionalRPC(streamPost service.ResultInfo_InsertPosts_Client, c chan struct{}) {
	for {
		res, err := streamPost.Recv()
		if err == io.EOF {
			break
		}
		// TODO - insert 성공 횟수 처리
		slog.Debug(res.Message)
	}
	<-c
}
