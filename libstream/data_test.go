package libstream

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var globalUUID_Stream []string

func InitializeDataInDB() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		"testu",
		"testup",
		"stream")
	DB, _ = sql.Open("postgres", dbinfo)

	for i := 0; i < 500; i++ {
		stream := NewStream()
		InsertToDB(&stream.S)
		//globalUUID_Stream = append(globalUUID_Stream, stream.ID)
	}
}

func TestStartStream(t *testing.T) {

	server := GetTestServer()

	req, _ := http.NewRequest("GET", "/test/run", nil)
	w := httptest.NewRecorder()
	//handler := http.HandlerFunc(server.StartNewStream)

	server.GetRouter().ServeHTTP(w, req)
	//assert.Equal(t, http.StatusFound, w.Code)
	//
	//
	//if rr.Code != http.StatusOK {
	//	t.Errorf("handler returned wrong status code: got %v want %v",
	//		rr.Code, http.StatusOK)
	//}
}

func TestInterruptStream(t *testing.T) {

	server := GetTestServer()
	InitializeDataInDB()

	for i := 20; i < 25; i++ {

		requestURL := fmt.Sprintf("/test/interrupt/%s", globalUUID_Stream[i])
		req, _ := http.NewRequest("GET", requestURL, nil)
		w := httptest.NewRecorder()

		server.GetRouter().ServeHTTP(w, req)
	}

}

func TestConcurrent(t *testing.T) {
	//wg := &sync.WaitGroup{}
	//for j := 0; j < 10; j++ {
	//	wg.Add(1)
	//	go func() {
	//		for i := 0; i < 500; i++ {
	//			TestPostData(t)
	//		}
	//		wg.Done()
	//	}()
	//	runtime.Gosched()
	//}
	//for j := 0; j < 10; j++ {
	//	wg.Add(1)
	//	go func() {
	//		for i := 0; i < 500; i++ {
	//			TestGetData(t)
	//		}
	//		wg.Done()
	//	}()
	//	runtime.Gosched()
	//}
	//wg.Wait()
}
