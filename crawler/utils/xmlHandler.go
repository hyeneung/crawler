package utils

import (
	"encoding/xml"
	"io"
	"net/http"
	"sync"
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
	var wg sync.WaitGroup
	numPosts := int32(len(posts))
	lastUpdatedDate := UnixTime2Time(updatedDate)
	pathStartIdx := len("https://") + len(domainURL)

	resultCh := make(chan int32)
	worker := func(start, end int32) {
		defer wg.Done()
		var lastIdx int32 = -1
		for i := start; i < end; i++ {
			post := posts[i]
			pubDate := Str2time(post.PubDate)
			if pubDate.Compare(lastUpdatedDate) == 1 {
				// URL parsing
				posts[i].Link = post.Link[pathStartIdx:]
				posts[i].Id = id
				lastIdx = i
			} else {
				break
			}
		}
		resultCh <- lastIdx
	}

	const workerCount int32 = 4
	// segmentSize = (numPosts//workerCount) + 1
	segmentSize := (numPosts + workerCount - 1) / workerCount

	for i := int32(0); i < workerCount; i++ {
		start := i * segmentSize
		end := start + segmentSize
		if end > numPosts {
			end = numPosts
		}
		wg.Add(1)
		go worker(start, end)
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var maxIdx int32 = -1
	for idx := range resultCh {
		if idx > maxIdx {
			maxIdx = idx
		}
	}

	return maxIdx
}
