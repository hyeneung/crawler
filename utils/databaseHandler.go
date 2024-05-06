package utils

// var (
//	errNotFound   = errors.New("Not Found")
// )

func insertDB(item Item) error {
	// fmt.Println("업데이트할 아이템", item)
	return nil
}

func UpdateDB(url string, updatedDate int64) (int8, int8) {
	posts := getParsedData(url)
	var updatedCount int8 = 0
	lastIndex := getLastIdxToUpdate(posts, updatedDate)
	// goroutine 적용 - updateCount 동기화 문제 괜찮은지 확인.
	for idx := int8(0); idx <= lastIndex; idx++ {
		err := insertDB(posts.Items[idx])
		if err == nil {
			updatedCount += 1
		}
	}
	return lastIndex + 1, updatedCount
}
