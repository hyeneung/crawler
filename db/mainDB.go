package db

import "errors"

// grpc unary
func InsertDomain(crawlerID uint64, domainURL string) error {
	// db 객체 메서드 호출
	return errors.New("InsertDomain error")
}

// grpc client streaming
func InsertPosts_(posts []Post) uint8 {
	return 0 // successCount
}

// grpc biddirectional streaming
func InsertPosts(posts []Post) []string {
	return []string{"some Log"}
}

// grpc server streaming
func GetLogs(id uint64) []string {
	return []string{"some Log"}
}

func MainDB() {
	// TODO - gRPC 서버 로직
	// TODO - connection pool 이용
}
