package main

import (
	"context"
	pb "crawler/service"
	"crawler/utils"
	"io"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/martinohmann/exit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Unary
func insertDomain(stub *pb.ResultInfoClient, id uint64, url string) string {
	logger := utils.SlogLogger.With(
		slog.Uint64("crawlerId", id),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := (*stub).InsertDomain(ctx, &pb.UnaryRequest{Id: id, Url: url})
	if err != nil {
		logger.Error(err.Error())
	}
	return res.Message
}

// grpc client streaming
func insertPost_clientStreaming(stub *pb.ResultInfoClient, posts *[]utils.Post, lastIdxToUpdate int32) uint32 {
	logger := utils.SlogLogger.With(
		slog.Uint64("crawlerId", (*posts)[0].Id),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	insertStream, err := (*stub).InsertPosts(ctx)
	utils.CheckErr(err, logger)

	// goroutine
	var wg sync.WaitGroup
	wg.Add(int(lastIdxToUpdate) + 1)

	worker := func(post utils.Post) {
		defer wg.Done()
		data := pb.Post{Id: post.Id, Title: post.Title, Link: post.Link, PubDate: post.PubDate}
		err := insertStream.Send(&data)
		utils.CheckErr(err, logger)
	}
	for i := 0; i < int(lastIdxToUpdate)+1; i++ {
		go worker((*posts)[i])
	}
	wg.Wait()

	res, err := insertStream.CloseAndRecv()
	if err != nil {
		logger.Error(err.Error())
	}
	return res.Value
}

// grpc server streaming
func saveServerLogs(stub *pb.ResultInfoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	stream, err := (*stub).GetLogs(ctx, &emptypb.Empty{})
	utils.CheckErr(err, utils.SlogLogger)

	filePath := "./log/serverLog.log"
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	utils.CheckErr(err, utils.SlogLogger)

	for {
		logs, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error(err.Error())
		}
		_, err = file.Write(logs.Message)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

// grpc biddirectional streaming
func insertPost_bidirectionalStreaming(stub *pb.ResultInfoClient, posts *[]utils.Post, lastIdxToUpdate int32) uint32 {
	logger := utils.SlogLogger.With(
		slog.Uint64("crawlerId", (*posts)[0].Id),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	insertStream, err := (*stub).InsertPosts_(ctx)
	utils.CheckErr(err, logger)

	// goroutine(receiver)
	c := make(chan uint32, 10)
	go receiveWorker(insertStream, c)

	// goroutine(sender)
	var wg sync.WaitGroup
	wg.Add(int(lastIdxToUpdate) + 1)
	sendWorker := func(post utils.Post) {
		defer wg.Done()
		data := pb.Post{Id: post.Id, Title: post.Title, Link: post.Link, PubDate: post.PubDate}
		err := insertStream.Send(&data)
		utils.CheckErr(err, logger)
	}
	for i := 0; i < int(lastIdxToUpdate)+1; i++ {
		go sendWorker((*posts)[i])
	}
	wg.Wait()

	// 다 보냈으면 close
	if err := insertStream.CloseSend(); err != nil {
		logger.Error(err.Error())
		exit.Exit(err)
	}
	return <-c // receiver로부터 결과 취합
}
func receiveWorker(streamPost pb.ResultInfo_InsertPosts_Client, c chan<- uint32) {
	var successCount uint32 = 0
	for {
		_, err := streamPost.Recv()
		if err == io.EOF {
			break
		}
		successCount++
	}
	c <- successCount
}

type crawlerInfo struct {
	name   string
	rssURL string
}

const (
	// https://github.com/grpc/grpc/blob/master/doc/naming.md
	address = "dns:localhost:50051"
)

func main() {
	const lineURL = "https://techblog.lycorp.co.jp/ko/feed/index.xml"
	const mercariURL = "https://engineering.mercari.com/en/blog/feed.xml"
	const kurlyURL = "https://helloworld.kurly.com/feed.xml"
	var crawlerInfos = []crawlerInfo{
		{name: "lineCrawler", rssURL: lineURL},
		{name: "mercariCrawler", rssURL: mercariURL},
		{name: "kurleyCrawler", rssURL: kurlyURL},
	}

	slog.SetDefault(utils.SlogLogger)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	stub := pb.NewResultInfoClient(conn)

	var wg sync.WaitGroup
	for i, c := range crawlerInfos {
		wg.Add(1)
		go func(id uint64, c crawlerInfo) {
			defer wg.Done()
			rssCrawler := New(id, c.name, c.rssURL)
			rssCrawler.Init(&stub)                   // domain DB에 crawler id, domain url 저장
			rssCrawler.Run(&stub, time.Now().Unix()) // post DB에 게시물 정보 저장
		}(uint64(i), c)
	}
	wg.Wait()

	// grpc server streaming
	slog.Info("starting gRPC server streaming")
	saveServerLogs(&stub)
	slog.Info("finished gRPC server streaming")
}
