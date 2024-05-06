package utils

func insertDB(item Item) error {
	// fmt.Println("업데이트할 아이템", item)
	return nil
}

func UpdateDB(url string, updatedDate int64) (int8, int8) {
	posts := getParsedData(url)
	var updatedCount int8 = 0
	lastIndex := getLastIdxToUpdate(posts, updatedDate)
	for idx := int8(0); idx <= lastIndex; idx++ {
		var item Item = posts.Items[idx]
		err := insertDB(item)
		if err == nil {
			updatedCount += 1
		}
	}
	return lastIndex + 1, updatedCount
}
