package service

import (
	"db/utils"
)

// RSSCrawler type
type Response struct {
	Id      uint64
	Message string
}

// grpc unary
func InsertDomain(crawlerID uint64, domainURL string) Response {
	err := utils.InsertDomainDB(crawlerID, domainURL)
	message, _ := utils.GetResponseMessage(err)
	res := Response{Id: crawlerID, Message: message}
	return res
}

// grpc client streaming
func InsertPosts_(crawlerID uint64, posts []utils.Post) uint8 {
	var updatedCount uint8 = 0
	for _, post := range posts {
		err := utils.InsertPostDB(crawlerID, post)
		_, is_success := utils.GetResponseMessage(err)
		if is_success {
			updatedCount += 1
		}
	}
	return updatedCount // successCount
}

// grpc biddirectional streaming
func InsertPosts(crawlerID uint64, posts []utils.Post) []Response {
	responses := []Response{}
	for _, post := range posts {
		err := utils.InsertPostDB(crawlerID, post)
		message, _ := utils.GetResponseMessage(err)
		responses = append(responses, Response{Id: crawlerID, Message: message})
	}
	return responses
}

// grpc server streaming
func GetLogs(id uint64) []Response {
	responses := []Response{}
	responses = append(responses, Response{Id: 1, Message: "test"})
	return responses
}
