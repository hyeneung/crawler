package main

import (
	"context"
	"db/service"
	"db/utils"
	"io"
	"log"
	"net"

	pb "db/service"

	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedResultInfoServer
}

func InsertDomain(crawlerID uint64, domainURL string) *service.Response {
	err := utils.InsertDomainDB(crawlerID, domainURL)
	message, _ := utils.GetResponseMessage(err)
	return &service.Response{Id: crawlerID, Message: message}
}
func (s *server) InsertDomain(ctx context.Context, in *pb.UnaryRequest) (*pb.Response, error) {
	err := utils.InsertDomainDB(in.Id, in.Url)
	message, _ := utils.GetResponseMessage(err)
	log.Println(message)
	return &pb.Response{Id: in.Id, Message: message}, status.New(codes.OK, "").Err()
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
			return err
		}
		err = utils.InsertPostDB(post)
		_, is_success := utils.GetResponseMessage(err)
		if is_success {
			updatedCount += 1
			log.Printf("Updated")
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

	responses := []service.Response{} // cannot use &response (value of type *Response) as *service.Response value
	for {
		post, err := stream.Recv()
		if err == io.EOF {
			// Client has sent all the messages
			// Send remaining shipments
			for _, response := range responses {
				if err := stream.Send(&response); err != nil {
					return err
				}
			}
			return nil
		}
		if err != nil {
			log.Println(err)
			return err
		}
		err = utils.InsertPostDB(post)
		message, _ := utils.GetResponseMessage(err)
		responses = append(responses, service.Response{Id: post.Id, Message: message})
	}
}

func main() {
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
