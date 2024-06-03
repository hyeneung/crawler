package main

import (
	"context"
	"db/service"
	"db/utils"
	"io"
	"log"
	"log/slog"
	"net"
	"os"

	pb "db/service"

	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedResultInfoServer
}

func (s *server) InsertDomain(ctx context.Context, in *pb.UnaryRequest) (*pb.Response, error) {
	err := utils.InsertDomainDB(in.Id, in.Url)
	message := utils.DBLogMessage(in.Id, err)
	return &pb.Response{Id: in.Id, Message: message}, err
}

// Client-side Streaming RPC
func (s *server) InsertPosts(stream pb.ResultInfo_InsertPostsServer) error {
	var updatedCount uint32 = 0
	for {
		post, err := stream.Recv()
		if err == io.EOF {
			// Finished reading the order stream.
			return stream.SendAndClose(&wrapper.UInt32Value{Value: updatedCount})
		}
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		err = utils.InsertPostDB(post)
		message := utils.DBLogMessage(post.Id, err)
		if err != nil {
			updatedCount += 1
			slog.Debug(message)
		}
	}
}

// Server-side Streaming RPC
func (s *server) GetLogs(searchQuery *wrapper.UInt64Value, stream pb.ResultInfo_GetLogsServer) error {

	// for key, order := range orderMap {
	// 	log.Print(key, order)
	// 	for _, itemStr := range order.Items {
	// 		log.Print(itemStr)
	// 		if strings.Contains(itemStr, searchQuery.Value) {
	// 			// Send the matching orders in a stream
	// 			err := stream.Send(&order)
	// 			if err != nil {
	// 				return fmt.Errorf("error sending message to stream : %v", err)
	// 			}
	// 			log.Print("Matching Order Found : " + key)
	// 			break
	// 		}
	// 	}
	// }
	return nil
}

// Bi-directional Streaming RPC
func (s *server) InsertPosts_(stream pb.ResultInfo_InsertPosts_Server) error {
	for {
		post, err := stream.Recv()
		if err == io.EOF {
			slog.Debug("Client has sent all the messages")
			return nil
		}
		if err != nil {
			slog.Error(err.Error())
			return err
		}
		err = utils.InsertPostDB(post)
		message := utils.DBLogMessage(post.Id, err)
		if err != nil {
			return err
		}
		response := service.Response{Id: post.Id, Message: message}
		if err := stream.Send(&response); err != nil {
			slog.Error(err.Error())
			return err
		}
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterResultInfoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
