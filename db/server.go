package main

import (
	"context"
	"db/utils"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net"
	"os"
	"sync"

	pb "db/service"

	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/martinohmann/exit"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedResultInfoServer
}

// Unary
func (s *server) InsertDomain(ctx context.Context, in *pb.UnaryRequest) (*pb.Response, error) {
	slog.Info("[Received] "+in.Url, "Id : ", in.Id)
	err := utils.InsertDomainDB(in.Id, in.Url)
	message := utils.DBLogMessage(in.Id, err) // log 남김
	return &pb.Response{Id: in.Id, Message: message}, err
}

// Client-side Streaming RPC
func (s *server) InsertPosts(stream pb.ResultInfo_InsertPostsServer) error {
	var updatedCount uint32 = 0
	for {
		post, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&wrapper.UInt32Value{Value: updatedCount})
		}
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		slog.Info("[Received] "+post.Title, "Id : ", post.Id)
		err = utils.InsertPostDB(post)
		utils.DBLogMessage(post.Id, err)
		if err == nil {
			updatedCount += 1
		} else {
			slog.Error(err.Error())
			return err
		}
	}
}

// Server-side Streaming RPC
func (s *server) GetLogs(empty *emptypb.Empty, stream pb.ResultInfo_GetLogsServer) error {
	dir := "./log"
	files, err := os.ReadDir(dir)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	if files == nil {
		slog.Error("directory is empty")
		return errors.New("directory is empty")
	}
	slog.Info("log file fetched")
	// goroutine
	var wg sync.WaitGroup
	errChan := make(chan error, len(files))
	for _, file := range files {
		wg.Add(1)
		go func(file fs.DirEntry) {
			defer wg.Done()
			logFile, err := os.ReadFile(dir + "/" + file.Name())
			if err != nil {
				errChan <- err
				return
			}
			if err := stream.Send(&pb.LogData{Message: logFile}); err != nil {
				errChan <- err
				return
			}
		}(file)
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()
	if err := <-errChan; err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

// Bi-directional Streaming RPC
func (s *server) InsertPosts_(stream pb.ResultInfo_InsertPosts_Server) error {
	var wg sync.WaitGroup
	postChan := make(chan *pb.Post, 10)
	errChan := make(chan error, 1)
	// goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			post, err := stream.Recv()
			if err == io.EOF {
				close(postChan)
				return
			}
			if err != nil {
				slog.Error(err.Error())
				errChan <- err
				return
			}
			slog.Info("[Received] "+post.Title, "Id : ", post.Id)
			postChan <- post
		}
	}()

	for post := range postChan {
		err := utils.InsertPostDB(post)
		message := utils.DBLogMessage(post.Id, err)
		if err != nil {
			close(errChan)
			return err
		}
		response := pb.Response{Id: post.Id, Message: message}
		if err := stream.Send(&response); err != nil {
			slog.Error(err.Error())
			close(errChan)
			return err
		}
	}
	wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

const (
	port = ":50051"
)

func main() {
	slog.SetDefault(utils.SlogLogger)
	slog.Info("server starting")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		slog.Error(err.Error())
		exit.Exit(err)
	}
	s := grpc.NewServer()
	pb.RegisterResultInfoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		slog.Error(err.Error())
		exit.Exit(err)
	}
}
