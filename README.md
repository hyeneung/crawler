## 프로젝트 내용
* 기술 블로그들 크롤링해서 한 번에 검색할 수 있는 사이트 만들기
### 프로젝트 전체 구성
1. main package
* main.go : 크롤러 인스턴스들 생성. 실행
* crawler.go : 크롤러 동작 정의. 
    * run 호출 시 변동사항 여부 확인해서 db에 저장하는 함수 호출
    * [풀스택네트워킹 과제] 이 함수 호출을 gRPC로.
    * main 패키지와 utils 패키지가 다른 곳에서 돌아갈 수 있다는 전제.
    * [마이크로서비스프로그래밍 과제] 두 패키지를 별도의 컨테이너에서 실행.
2. utils package
* databaseHandler : db접근 이용
* errorHandler : 에러처리 , log남김
* timeHandler : time.Time이 heap쓰는거 피하려고 unixTime으로 저장하는 바람에 int64 string time 간에 변환할 일이 생겨서 만듦
* xmlHandler : http get 요청으로 RSS feed 파일 받아서 xml 파싱

### 해야할 것(우선순위)
1. 동작하는 코드
    완료
2. 학교 과제
    * docker compose로 mysql, line크롤러, 메루카리 크롤러 3개 서비스 구성.
    * main의 run에서 utils package 함수를 gRPC로 호출하도록
3. 성능 개선
    * 전체적으로 goroutine 적용해서 비동기처리
    * post 테이블의 url 도메인 부분/이외 부분 나누어서 post테이블, domain 테이블에 저장
4. 운영 측면 기능 추가
    * 하루 단위로 main.go의 RunCrawlers 호출되도록 (linux crontab 적용)
    * cmd로 xml path와 크롤러 이름, 크롤러 id 시작 번호 받아서 main 수행
5. 서비스 완전성(방학때 웹 공부용)
    * rss feed가 제공하지 않는 예전 post 정보 크롤링. db 저장. 
    * 크롤링한 정보를 웹에 게시판 형식으로 게시. 원하는 내용 검색해서 볼 수 있도록.

### DB(crawl_data)
 - insert할 데이터 : (title, url, pubDate)
    - title과 link는 필수적으로 포함. pubDate는 선택적. [참고] https://www.rssboard.org/rss-specification
 - 고민할 것 : (1)url 크기, url 중복도메인, (2)같은 게시물 중복 방지
 - (1)
    * 방법 1 : id를 크롤러별로 만들고, url 중복도메인 지우고, 크롤러-중복도메인 매핑 테이블 따로 만들기,
    * 방법 2 : url varchar(큰 값)해서 한 테이블에 넣기.
 - (2)
    * id 만들고 AUTO_INCREMENT로 하면 DB 수준에서 중복 처리 못함. 추가 코드 필요 -> goroutine 동기화 문제 걱정됨.
    * DB 기능 활용 위해 id 대신 크롤러id, pubDate(unix time)을 composite primary key로.
 - 최종 스키마 : post(crawler_id(PK), url(PK), title, pubDate(nullable) ), domain(crawler_id(PK,FK), domain_url)