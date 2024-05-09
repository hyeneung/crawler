## 프로젝트 내용
* 기술 블로그들 검색 사이트 만들기
  * 회사 상관 없이 제목으로 검색 가능하도록 DB구성, 검색 웹사이트 구축
### 프로젝트 전체 구성
1. DB 저장 코드
* main package
  * main.go : 크롤러 인스턴스들 생성, 실행
  * crawler.go : 크롤러 동작 정의. 
* utils package
  * databaseHandler : db접근, 이용
  * errorHandler : 에러처리 , log남김
  * timeHandler : int64 string time 간 시간 표현 변환
    * time.Time이 heap쓰는거 피하기 위함
  * xmlHandler : http get 요청으로 RSS feed 파일 받아서 xml 파싱
2. 웹 사이트 운영
* later
### 해야할 것(우선순위)
1. 최소 기능 동작-rss feed를 DB에 저장 (완료)
2. 학교 과제
    * [x] docker compose로 mysql, 크롤러 2개 서비스 구성.
    * [ ] main의 run에서 utils package 함수를 gRPC로 호출
3. 성능 개선
    * [ ] 전체적으로 goroutine 적용해서 비동기처리
    * [ ] post 테이블의 url 도메인 부분/이외 부분 나누어서 post테이블, domain 테이블에 저장
    * [ ] db 계정, 비번 관리
4. 운영 측면 기능 추가
    * [ ] 하루 단위로 기술 블로그 업데이트 확인 및 반영 (linux crontab 적용)
    * [ ] cmd로 xml path와 크롤러 이름, 크롤러 id 시작 번호 받아서 main 수행
5. 서비스 완전성(방학때 웹 공부용)
    * [ ] rss feed가 제공하지 않는 예전 post 정보 크롤링. db 저장. 
    * [ ] 크롤링한 정보를 웹에 게시판 형식으로 게시. 원하는 내용 검색해서 볼 수 있도록.

###  docker로 실행
1. Dockerfile기반 go 이미지 빌드
```shell
docker build -t my-crawler:1.0 -f docker/Dockerfile .
```
2. docker compose 실행(해당 이미지 이용)
   * 한 번만 하면 volume mount가 돼서 컨테이너 지우고 다시 실행해도 남아 있음
   * utils/db/data 폴더를 지우거나 crawler.go 의 생성자에서 lastUpdated를 바꾸면 됨 
```shell
docker compose -f ./docker/docker-compose.yml up -d
```
3. 실행 결과
  * mysqladdmin에 성공적으로 연결되면 crawler 서비스 실행
![alt text](doc\imgs\compose-result.png)
4. docker hub에 push
```shell
docker tag my-crawler:1.0 hyeneung/crawler:1.0
```
```shell
docker push hyeneung/crawler:1.0
```
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
      * DB 수준에서 중복 처리 못함. 추가 코드 필요 -> goroutine 동기화 문제 걱정됨.
    * 방법 2 : 크롤러id, pubDate(unix time)을 composite primary key로.
      * DB 기능 활용 가능
 * 최종 스키마
   *  domain(crawler_id(PK), domain_url)
   *  post(crawler_id(PK, FK), url(PK), title, pubDate(nullable) )
      *  pubDate nullable인 이유 : RSS 파일에서 optional임 [참고]https://www.rssboard.org/rss-specification
2. 트랜잭션 적용
 * 필요 없음

### 코드 설계
 * DB handler에서 db 연결 코드 중복(crawler init, run)
   * 객체로 만들어서 connection pool 관리
     * gRPC로 public 메서드 호출 가능한가? - 일단 과제용부터 구현하고 고민하는 걸로.