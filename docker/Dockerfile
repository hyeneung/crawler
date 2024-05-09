FROM golang:1.22

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum *.go ./
RUN go mod download && go mod verify

COPY utils/db/init_db.sql utils/*.go ./utils/
RUN go build -v -o /usr/local/bin/execute_crawler ./
# CMD ["빌드 파일 이름 - execute_crawler"]