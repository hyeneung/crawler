package utils

import (
	"log"
	"log/slog"
	"os"
	"time"

	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

var lumberjackLogger = &lumberjack.Logger{
	Filename:   getLogFilePath(),
	MaxSize:    20,   // A file can be up to 20M.
	MaxBackups: 5,    // Save up to 5 files at the same time
	MaxAge:     10,   // A file can be saved for up to 10 days.
	Compress:   true, // Compress with gzip.
}

var SlogLogger = slog.New(
	slogmulti.Fanout(
		slog.NewJSONHandler(os.Stdout, nil),
		slog.NewJSONHandler(lumberjackLogger, nil),
	),
)

func getLogFilePath() string {
	logPath := "./log/"
	logFileName := time.Now().Format("2006-01-01_15h04m05s") + ".log"
	fileName := logPath + logFileName
	if _, err := os.Stat(fileName); err != nil {
		log.Println("creating file", fileName)
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
		}
	}
	return fileName
}

func LogInit(message string, crawlerName string) {
	initLogger := SlogLogger.With(
		slog.String("crawler", crawlerName),
		slog.String("db_log", message),
	)
	slog.SetDefault(initLogger)
	initLogger.Info("Init Done")
}

// log the result
func LogRun(crawlerName string, totalCount int32, successCount uint32) {
	resultLogger := SlogLogger.With(
		slog.String("crawler", crawlerName),
		slog.Group("results",
			"newly posted", totalCount,
			"updated", successCount,
		),
	)
	slog.SetDefault(resultLogger)

	resultLogger.Info("Done")
}
