package main

import (
	pb "crawler/service"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// https://github.com/grpc/grpc/blob/master/doc/naming.md
	address = "dns:localhost:50051"
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
}
