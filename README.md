## 프로젝트 내용
* 기술 블로그들 검색 사이트 만들기
  * 회사 상관 없이 제목으로 검색 가능하도록 DB구성, 검색 웹사이트 구축
### 프로젝트 전체 구성
1. DB 저장 코드
* main package
  * main : 크롤러 인스턴스들 생성, 실행
  * crawler.go : 크롤러 동작 정의. 
* utils package
  * xmlHandler : http get 요청으로 RSS feed 파일 받아서 xml 파싱
  * errorHandler : 에러처리 , log남김
  * timeHandler : int64, string, time.Time 간 시간 표현 변환
    * time.Time이 heap쓰는거 피하기 위함
* utils/db package - 다른 컨테이너에서 실행될 것들
  * databaseHandler : DB handler 동작 정의
  * errorHandler : 에러처리 , log남김
  * mainDB : DB handler 객체 생성. gRPC 호출되는 함수 정의

1. 웹 사이트 운영
* later
### 해야할 것(우선순위)
1. 최소 기능 동작-rss feed를 DB에 저장 (완료)
2. 학교 과제
    * [x] docker compose로 mysql, 크롤러 2개 서비스 구성.
    * [ ] main package에서 db package의 함수를 gRPC로 호출
3. 성능 개선
    * [ ] 전체적으로 goroutine 적용해서 비동기처리
4. 운영 측면 기능 추가
    * [ ] 하루 단위로 프로그램 실행. 기술 블로그 업데이트 확인 및 반영 (linux crontab 적용)
    * [ ] cmd로 xml path와 크롤러 이름, 크롤러 id 시작 번호 받아서 프로그램 수행
5. 서비스 완전성(방학때 웹 공부용)
    * [ ] rss feed가 제공하지 않는 예전 post 정보 크롤링. db 저장. 
    * [ ] 크롤링한 정보를 웹에 게시판 형식으로 게시. 원하는 내용 검색해서 볼 수 있도록.
    * [ ] db 계정, 비번 정보 코드에서 빼서 따로 관리
###  docker로 실행
1. Dockerfile기반 go 이미지 빌드
```shell
docker build -t my-crawler:1.0 -f docker/Dockerfile .
```
2. docker compose 실행(해당 이미지 이용)
   * 한 번만 하면 volume mount가 돼서 컨테이너 지우고 다시 실행해도 남아 있음
   * utils/db/data 폴더를 지우거나 crawler.go 의 생성자에서 lastUpdated를 바꾸면 됨 
```shell
docker compose -f ./docker/docker-compose.yml up
```
3. 실행 결과
  * mysqladdmin에 성공적으로 연결되면 crawler 서비스 실행
![image](https://github.com/hyeneung/crawler/assets/71257602/7cb3f08c-a5c2-4947-8242-d713790da02b)
  * 당시 코드 : https://github.com/hyeneung/crawler/tree/a8cf288694e468338426b4d56386ad25eb273265
4. docker hub에 push
```shell
docker tag my-crawler:1.0 hyeneung/crawler:1.0
```
```shell
docker push hyeneung/crawler:1.0
```
### goroutine 적용기
1. 하나의 프로그램 내. 크롤러 동작 비동기 처리
     * [이전] main.go에서 한 회사의 블로그 정보를 다 저장해야 다른 회사꺼 DB에 저장
2. 하나의 크롤러 내. 다른 DB들 insert 동작 비동기 처리
     * [이전] domain DB에 insert한 후에야 post DB에 insert
3. 하나의 DB 작업 내. 다른 instance(record)들 DB insert 동작 비동기 처리
     * [이전] post DB에 이전 게시물 insert한 후에야 다음 게시물 insert
     * insert 작업 성공 횟수를 공유해서 수정 하기에 동기화 문제 발생 가능.
     * gRPC 쓰면서 goroutine 가능한지 확인할 것.
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
 * DB handler 객체로 만들어서 connection pool 관리
   * crawler의 확장성과 가용한 connection 수의 한계를 고려함