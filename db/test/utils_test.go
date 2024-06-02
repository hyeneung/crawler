package test

import (
	"db/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertDomain(t *testing.T) {
	res := service.InsertDomain(1, "test url")
	assert.Equal(t, res.Id, uint64(1))
	assert.Equal(t, res.Message, "Succeed")

	res = service.InsertDomain(1, "test url")
	assert.Equal(t, res.Message, "[Failed] Dupulicated data insertion. Change the \"lastUpdated\" value in crawler.go file or Delete utils/db/data")

	largeBytes := make([]byte, 501)
	longString := string(largeBytes)
	res = service.InsertDomain(1, longString)
	assert.Equal(t, res.Message, "[Failed] URL or title exceeded 500 bytes")
}
