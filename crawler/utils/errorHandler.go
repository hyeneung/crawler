package utils

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/martinohmann/exit"
)

func CheckHttpResponse(resp *http.Response) {
	if resp.StatusCode != http.StatusOK {
		slog.SetDefault(SlogLogger)
		slog.Error("HTTP response error", slog.Int("Status Code", resp.StatusCode))
		exit.Exit(errors.New("failed to parse xml"))
	}
}

func CheckErr(err error) {
	if err != nil {
		slog.SetDefault(SlogLogger)
		slog.Error(err.Error())
		// exit.Exit(err)
	}
}
