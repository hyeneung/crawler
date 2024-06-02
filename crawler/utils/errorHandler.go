package utils

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/martinohmann/exit"
)

func CheckHttpResponse(resp *http.Response) {
	if resp.StatusCode != http.StatusOK {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
		slog.Error("HTTP response error", slog.Int("Status Code", resp.StatusCode))
		exit.Exit(errors.New("failed to parse xml"))
	}
}

func CheckErr(err error) {
	if err != nil {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
		slog.Error(err.Error())
		exit.Exit(err)
	}
}
