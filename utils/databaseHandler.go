package utils

import (
	// "database/sql"
	// "fmt"
	// "log"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// var (
//	errNotFound   = errors.New("Not Found")
// )

func insertDB(items []Item, crawlerID uint64) error {
	fmt.Println("업데이트할 아이템", items)
	fmt.Println("호출한 크롤러: ", crawlerID)
	return nil
}

func UpdateDB(crawlerID uint64, url string, updatedDate int64) (int8, int8) {
	posts := getParsedData(url)
	var updatedCount int8 = 0
	lastIndex := getLastIdxToUpdate(posts, updatedDate)
	var items []Item
	for idx := int8(0); idx <= lastIndex; idx++ {
		items = append(items, posts.Items[idx])
	}
	insertDB(items, crawlerID)
	return lastIndex + 1, updatedCount
}
