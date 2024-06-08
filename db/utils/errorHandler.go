package utils

import (
	"log/slog"

	"github.com/go-sql-driver/mysql"
	"github.com/martinohmann/exit"
)

func DBLogMessage(title string, crawlerId uint64, err error) string {
	var message string
	var logger *slog.Logger
	if title == "" {
		logger = SlogLogger.With(
			slog.Uint64("crawler_id", crawlerId),
		)
	} else {
		logger = SlogLogger.With(
			slog.Uint64("crawler_id", crawlerId),
			slog.String("title", title),
		)
	}
	if err == nil {
		message = "Succeed"
		logger.Info(message)
		return message
	}
	mysqlErr, _ := err.(*mysql.MySQLError)
	if mysqlErr == nil {
		message = "[Failed] Not using MySQL"
		logger.Error(message)
		logger.Error(err.Error())
		return message
	}
	switch mysqlErr.Number {
	case 1406:
		message = "[Failed] URL or title exceeded 500 bytes"
		logger.Error(message)
	case 1062:
		message = "[Failed] Dupulicated data insertion. Change the \"lastUpdated\" value in crawler.go file or Delete utils/db/data"
		logger.Error(message)
	default:
		message = err.Error()
		logger.Error(message)
	}
	return message
}

func checkFatalErr(err error) {
	if err != nil {
		slog.SetDefault(SlogLogger)
		slog.Error(err.Error())
		exit.Exit(err)
	}
}
