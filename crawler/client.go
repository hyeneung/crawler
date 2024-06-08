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

// grpc server streaming(not using)
func saveServerLogs(stub *pb.ResultInfoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	stream, err := (*stub).GetLogs(ctx, &emptypb.Empty{})
	utils.CheckErr(err, utils.SlogLogger)

	filePath := "./log/serverLog.log"
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	utils.CheckErr(err, utils.SlogLogger)
	defer file.Close()

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

// grpc biddirectional streaming(not using)
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
		res, err := streamPost.Recv()
		if err == io.EOF {
			break
		}
		if res.Message == "Succeed" {
			successCount++
		}
	}
	c <- successCount
}

func main() {
	address := os.Getenv("GRPC_SERVER_ADDRESS")

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	stub := pb.NewResultInfoClient(conn)

	slog.SetDefault(utils.SlogLogger)
	// rss 2.0 standards : https://www.rssboard.org/rss-specification
	configFilePath := "./config-crawler.yaml"
	config := getConfig(configFilePath)

	var wg sync.WaitGroup
	for i, c := range config.Crawlers {
		wg.Add(1)
		go func(id uint64, crawler Crawler) {
			defer wg.Done()
			crawler.Init(&stub, &config) // domain table에 crawler id, domain url 저장
			crawler.Run(&stub, &config)  // post table에 게시물 정보 저장
		}(uint64(i), c)
	}
	wg.Wait()
	writeConfig(configFilePath, config)
	slog.Info("Successfully Finished")
}
