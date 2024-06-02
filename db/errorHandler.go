package db

import (
	"log/slog"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/martinohmann/exit"
)

// grpc
func checkDBInsertErr(err error) error {
	if err != nil {
		mysqlErr, _ := err.(*mysql.MySQLError)
		if mysqlErr != nil {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			slog.SetDefault(logger)
			// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
			switch mysqlErr.Number {
			case 1406: // pathURL : 500byte, title : 500byte 이거보다 큰 데이터 들어오면 에러 발생.
				// 코드 문제가 아니라 데이터의 문제라서 무시해도 될 에러들 처리. exit 호출 안함.
				slog.Error("Too long to insert. Skipping this data...")
				// case 1062:
				// 	// 같은 데이터를 db에 넣는 연산 반복하면 기본키 중복 입력해서 에러 발생.
				// 	slog.Error("Executed repeatedly. Change the \"lastUpdated\" value in crawler.go file or Delete utils/db/data")
			}
		} else {
			return checkErr(err)
		}
	}
	return nil
}

// grpc
func checkErr(err error) error {
	if err != nil {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
		slog.Error(err.Error())
		defer exit.Exit(err)
		return err
	}
	return nil
}
