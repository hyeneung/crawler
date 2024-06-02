package utils

import (
	"log/slog"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/martinohmann/exit"
)

func GetResponseMessage(err error) (string, bool) {
	var message string
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	if err == nil {
		return "Succeed", true
	}
	mysqlErr, _ := err.(*mysql.MySQLError)
	if mysqlErr == nil {
		message = "[Failed] Not using MySQL"
		slog.Error(message)
		return message, false
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
	}
	return message, false
}

func checkFatalErr(err error) {
	if err != nil {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
		slog.Error(err.Error())
		exit.Exit(err)
	}
}
