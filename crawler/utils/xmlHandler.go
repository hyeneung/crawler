package utils

import (
	"encoding/xml"
	"io"
	"net/http"
)

type ParsedData struct {
	Data []Post `xml:"channel>item"`
}

type Post struct {
	Id      uint64
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

func GetParsedData(url string) []Post {
	res, err := http.Get(url)
	CheckErr(err, SlogLogger)
	CheckHttpResponse(res)
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	CheckErr(err, SlogLogger)
	var posts ParsedData
	xmlerr := xml.Unmarshal(data, &posts)
	CheckErr(xmlerr, SlogLogger)
	return posts.Data
}

func CheckUpdatedPost(posts []Post, id uint64, domainURL string, updatedDate int64) int32 {
	lastUpdatedDate := UnixTime2Time(updatedDate)
	var index int32 = 0
	pathStartIdx := len("https://") + len(domainURL)
	for index < int32(len(posts)) {
		post := posts[index]
		// id 할당
		posts[index].Id = id
		pubDate := Str2time(post.PubDate)
		if pubDate.Compare(lastUpdatedDate) == 1 {
			// URL parsing
			posts[index].Link = post.Link[pathStartIdx:]
			index++ // check next post when it needs to update
		} else {
			break
		}
	}
	return index - 1
}
