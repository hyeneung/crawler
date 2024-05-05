package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

var (
	errNotFound   = errors.New("Not Found")
	errCantUpdate = errors.New("Cant update non-existing word")
	errWordExists = errors.New("That word already exists")
)

//	if xmlerr != nil {
//		panic(xmlerr)
//	}
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if res.StatusCode > 299 {
//		log.Fatalf("Response failed with status code: %d and\nbody:", res.StatusCode)
//	}
//
//	if err != nil {
//		log.Fatal(err)
//	}
func CheckUnmarshalErr(err error) {

}
func CheckGetXMLErr(err error) {

}

func CheckHttpResponse(resp *http.Response) {
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Status error: %v", resp.StatusCode)
	}
}

func CheckIOErr(err error) {

}
func CheckParseErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func CheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// String to print for the debug
// func (r RSSCrawler) String() string {
// 	var stringToReturn string
// 	date := r.lastUpdatedTime.Format(time.RFC822)
// 	switch r.status {
// 	case "NOT UPDATED":
// 		stringToReturn = fmt.Sprintf("[%s] No contents Updated", r.crawlerName)
// 	case "UPDATED":
// 		stringToReturn = fmt.Sprintf("[%s] updated database at %s", r.crawlerName, date)
// 	case "FAILED":
// 		stringToReturn = fmt.Sprintf("[%s] failed. \n\tERROR : %s", r.crawlerName, r.err)
// 	}
// 	return stringToReturn
// }
