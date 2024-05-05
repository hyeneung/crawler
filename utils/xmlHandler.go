package utils

import (
	"encoding/xml"
	"io"
	"net/http"
)

type ParsedData struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

func getParsedData(url string) *ParsedData {
	res, err := http.Get(url)
	CheckGetXMLErr(err)
	CheckHttpResponse(res)
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	CheckIOErr(err)
	var posts ParsedData
	xmlerr := xml.Unmarshal(data, &posts)
	CheckUnmarshalErr(xmlerr)
	return &posts
}

func getLastIdxToUpdate(posts *ParsedData, updatedDate int64) int8 {
	lastUpdatedDate := UnixTime2Time(updatedDate)
	// fmt.Println("이전 업데이트한 날짜", lastUpdatedDate)
	var index int8 = 0
	for index < int8(len(posts.Items)) {
		pubDate := Str2time(posts.Items[index].PubDate)
		if pubDate.Compare(lastUpdatedDate) == 1 {
			// fmt.Println("업데이트 대상 날짜", pubDate)
			index++ // check next post when it needs to update
		} else {
			break
		}
	}
	return index - 1
}

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
