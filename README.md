## 프로젝트 내용
* 기술 블로그들 검색 사이트 만들기
  * 회사 상관 없이 제목으로 검색 가능하도록 DB구성, 검색 웹사이트 구축
### 프로젝트 전체 구성
1. crawler
   - RSS feed 파일 바탕으로 게시물 크롤링
   - gRPC client. gRPC로 db(handler)의 insert 함수 호출
2. db(handler)
   - gRPC server. 게시물 정보 받아 db에 insert
3. 웹 사이트 운영
   * later
### 해야할 것(우선순위)
1. 최소 기능 동작-rss feed를 DB에 저장 (완료)
2. 학교 과제
    * [x] docker compose로 mysql, 크롤러 2개 서비스 구성.
    * [x] gRPC로 crawler에서 db(handler)의 함수 호출. 4가지 방식 이용 
    * [ ] kubernetes 이용해서 crawler pod과 db pod 만들어 통신
         * crawler/db/mysql 컨테이너 만듦
         * db, mysql은 같은 pod에. mysql volume 추가
         * grpc_health_probe와 mysql health check 넣을 것
3. 성능 개선
    * [x] 전체적으로 goroutine 적용해서 비동기처리
4. 운영 측면 기능 추가
    * [ ] 하루 단위로 프로그램 실행. 기술 블로그 업데이트 확인 및 반영 (linux crontab 적용 or k8s cronjob)
    * [x] 설정 파일로 크롤러 정보 받아 프로그램 수행
5. 서비스 완전성(방학때 웹 공부용)
    * [ ] rss feed가 제공하지 않는 예전 post 정보 크롤링. db 저장. 
    * [ ] 크롤링한 정보를 웹에 게시판 형식으로 게시. 원하는 내용 검색해서 볼 수 있도록.
    * [ ] db 계정, 비번 정보 환경변수로 따로 관리
###  docker로 크롤링 정보 저장 코드 실행
1. git pull 없이 docker hub 이미지와 최소 파일 만으로 실행
     - distribute 폴더 다운, config-crawler.yaml 에서 crawler 정보 수정, 실행
```shell
docker compose up
```
1. git pull 후 직접 이미지 빌드하여 실행
```shell
docker build -t db-handler:1.0 -f docker/Dockerfile_dbhandler . 
```
  - ERROR: failed to solve: Canceled: context canceled 시 db/data 폴더 삭제
```shell
docker build -t crawler:1.0 -f docker/Dockerfile_crawler .
```
```shell
docker compose -f ./docker/docker-compose-on-hand.yml up
```
1. 결과
  - mysql 컨테이너, db-handler(grpc server), crawler(grpc client) 순으로 실행됨. 시간 좀 걸림.
### goroutine 적용
1. 각 블로그 별 크롤러 동작 비동기 처리
     * crawler/client.go/main
2. DB에 새로 반영해야 할 게시물 확인. 확인해야 할 영역 나누어 비동기 처리
     * crawler/utils/xmlHandler.go/CheckUpdatedPost
3. gRPC streaming 함수들. 여러 개 보내는 동작, 받는 동작 비동기 처리. 
     * crawler/client.go/  db/server.go/
## 고민한 것
### DB
1. 스키마
 * url 중복도메인 - 한 회사의 기술 블로그 게시물 많음
    * 방법 1 : url varchar(큰 값) 혹은 text 타입으로 한 테이블에 넣기
      * 메모리 낭비 생기지만 링크 보여줄 때 조인 안해도 됨.
    * 방법 2 : id를 크롤러별로 만들고, url 중복도메인 지우고, 크롤러-중복도메인 매핑 테이블 따로 만들기
      * 메모리 효율적. 
      * 게시물 특성 상 한 번 넣으면 수정도 변경도 안되는 상황 -> 캐싱
 * primary key 무엇으로 할 것인가
    * 방법 1 : id 만들고 AUTO_INCREMENT
      * DB 수준에서 중복 처리 못함. 코드 레벨에서 처리 -> goroutine 동기화 처리해야. 병목 예상됨
    * 방법 2 : 크롤러id, pubDate(unix time)을 composite primary key로.
      * DB 기능 활용 가능
 * 최종 스키마
   *  domain(crawler_id(PK), domain_url)
   *  post(crawler_id(PK, FK), url(PK), title, pubDate(nullable) )
      *  pubDate nullable인 이유 : RSS 파일에서 optional임 [참고]https://www.rssboard.org/rss-specification
2. 트랜잭션 적용
 * 하나씩 insert만 하니까 트랜잭션 적용 시 db 성능 저하 우려되어 적용 안함