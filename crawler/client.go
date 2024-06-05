package main

import (
	"context"
	pb "crawler/service"
	"crawler/utils"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Unary
func insertDomain(stub *pb.ResultInfoClient, id uint64, url string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := (*stub).InsertDomain(ctx, &pb.UnaryRequest{Id: id, Url: url})
	if err != nil {
		slog.Error(err.Error())
	}
	return res.Message
}

// grpc client streaming
func insertPost_clientStreaming(stub *pb.ResultInfoClient, posts *[]utils.Post, lastIdxToUpdate int32) uint32 {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	insertStream, err := (*stub).InsertPosts(ctx)
	utils.CheckErr(err)
	for i := 0; i < int(lastIdxToUpdate+1); i++ {
		post := (*posts)[i]
		data := pb.Post{Id: post.Id, Title: post.Title,
			Link: post.Link, PubDate: post.PubDate}
		err := insertStream.Send(&data)
		utils.CheckErr(err)
	}
	res, err := insertStream.CloseAndRecv()
	if err != nil {
		slog.Error(err.Error())
	}
	return res.Value
}

// grpc server streaming
func printServerLogs(stub *pb.ResultInfoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	stream, err := (*stub).GetLogs(ctx, &emptypb.Empty{})
	utils.CheckErr(err)
	for {
		logs, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error(err.Error())
		}
		dbLog := strings.Replace(string(logs.Message), `\"`, `"`, -1)
		fmt.Println("server streaming", dbLog)
	}
}

// grpc biddirectional streaming
func insertPost_bidirectionalStreaming(stub *pb.ResultInfoClient, posts *[]utils.Post, lastIdxToUpdate int32) uint32 {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	insertStream, err := (*stub).InsertPosts_(ctx)
	utils.CheckErr(err)

	c := make(chan uint32)
	go asncClientBidirectionalRPC(insertStream, c)

	for i := 0; i < int(lastIdxToUpdate+1); i++ {
		post := (*posts)[i]
		data := pb.Post{Id: post.Id, Title: post.Title,
			Link: post.Link, PubDate: post.PubDate}
		err := insertStream.Send(&data)
		utils.CheckErr(err)
	}
	if err := insertStream.CloseSend(); err != nil {
		log.Fatal(err)
	}
	return <-c
}
func asncClientBidirectionalRPC(streamPost pb.ResultInfo_InsertPosts_Client, c chan<- uint32) {
	var successCount uint32 = 0
	for {
		res, err := streamPost.Recv()
		if err == io.EOF {
			break
		}
		successCount++
		slog.Info(res.Message)
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

	var id uint64 = 0
	for _, crawlerInfo := range crawlerInfos {
		// TODO - goroutine 적용
		rssCrawler := New(id, crawlerInfo.name, crawlerInfo.rssURL)
		rssCrawler.Init(&stub)                   // domain DB에 crawler id, domain url 저장
		rssCrawler.Run(&stub, time.Now().Unix()) // post DB에 게시물 정보 저장
		id += 1
	}
	// grpc server streaming
	slog.Info("starting gRPC server streaming")
	printServerLogs(&stub)
}
