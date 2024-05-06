package utils

import (
	"log/slog"
	"net/http"
	"os"
)

func CheckHttpResponse(resp *http.Response) {
	if resp.StatusCode != http.StatusOK {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
		slog.Error("HTTP response error", slog.Int("Status Code", resp.StatusCode))
	}
}

func CheckErr(err error) {
	if err != nil {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
		slog.Any("error", err)
	}
}
