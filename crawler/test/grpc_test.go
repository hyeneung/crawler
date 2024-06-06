package test

import (
	"crawler/utils"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInsertDomain(t *testing.T) {
	posts := utils.GetParsedData("https://techblog.lycorp.co.jp/ko/feed/index.xml")

	domainURL := strings.Split("https://techblog.lycorp.co.jp/ko/feed/index.xml", "/")[2]

	lastIdxToUpdate := utils.CheckUpdatedPost(posts, uint64(13), domainURL, time.Now().Unix())

	assert.Equal(t, posts[0].Id, uint64(13))
	assert.Equal(t, lastIdxToUpdate, int8(-1))
}
