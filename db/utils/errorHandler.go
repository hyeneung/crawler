package utils

import (
	"log/slog"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/martinohmann/exit"
)

func DBLogMessage(crawlerId uint64, err error) string {
	var message string
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}).WithAttrs([]slog.Attr{
		slog.Uint64("crawlerId", crawlerId),
	})
	logger := slog.New(logHandler)

	slog.SetDefault(logger)
	if err == nil {
		message = "Succeed"
		slog.Debug(message)
		return message
	}
	mysqlErr, _ := err.(*mysql.MySQLError)
	if mysqlErr == nil {
		message = "[Failed] Not using MySQL"
		slog.Error(message)
		return message
	}
	switch mysqlErr.Number {
	case 1406:
		message = "[Failed] URL or title exceeded 500 bytes"
		slog.Error(message)
	case 1062:
		message = "[Failed] Dupulicated data insertion. Change the \"lastUpdated\" value in crawler.go file or Delete utils/db/data"
		slog.Error(message)
	default:
		message = err.Error()
		slog.Error(message)
	}
	return message
}

func checkFatalErr(err error) {
	if err != nil {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
		slog.Error(err.Error())
		exit.Exit(err)
	}
}
