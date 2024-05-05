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
* utils errorHandler 코드 중복 없애고 로그 남기는 로직 추가
* databaseHandler 구현
2. 학교 과제
* docker compose로 mysql, line크롤러, 메루카리 크롤러 3개 서비스 구성.
* main의 run에서 utils package 함수를 gRPC로 호출하도록
3. 성능 개선
* 전체적으로 goroutine 적용해서 비동기처리
4. 운영 측면 기능 추가
* 하루 단위로 main.go의 RunCrawlers 호출되도록 (linux crontab 적용)
5. 서비스 완전성(방학때 웹 공부용)
* rss feed가 제공하지 않는 예전 post 정보 크롤링. db 저장. 
* 크롤링한 정보를 웹에 게시판 형식으로 게시. 원하는 내용 검색해서 볼 수 있도록.
