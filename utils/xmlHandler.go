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
	CheckErr(err)
	CheckHttpResponse(res)
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	CheckErr(err)
	var posts ParsedData
	xmlerr := xml.Unmarshal(data, &posts)
	CheckErr(xmlerr)
	return &posts
}

func getLastIdxToUpdate(posts *ParsedData, updatedDate int64) int8 {
	lastUpdatedDate := UnixTime2Time(updatedDate)
	var index int8 = 0
	for index < int8(len(posts.Items)) {
		pubDate := Str2time(posts.Items[index].PubDate)
		if pubDate.Compare(lastUpdatedDate) == 1 {
			index++ // check next post when it needs to update
		} else {
			break
		}
	}
	return index - 1
}
