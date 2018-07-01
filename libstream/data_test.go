package libstream

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

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


	req, _ := http.NewRequest("GET", "/test/interrupt/", nil)
	w := httptest.NewRecorder()
	//handler := http.HandlerFunc(server.StartNewStream)

	server.GetRouter().ServeHTTP(w, req)
	//testData := CreateToInsert(9978)
	//requestBody := bytes.NewBuffer(nil)
	//jsonapi.MarshalOnePayloadEmbedded(requestBody, testData)
	//
	//req, err := http.NewRequest("POST", "/api/v1/values", requestBody)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//req.Header.Set("Content-Type", jsonapi.MediaType)
	//
	//rr := httptest.NewRecorder()
	//handler := http.HandlerFunc(PostData)
	//handler.ServeHTTP(rr, req)
	//
	//if rr.Code != http.StatusOK {
	//	t.Errorf("handler returned wrong status code: got %v want %v",
	//		rr.Code, http.StatusOK)
	//}

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

//func CreateToInsert(id int) *NestedData {
	//var k NestedData
	//k.ID = id
	//k.Name = "test" + strconv.Itoa(id)
	//k.Value = "test" + strconv.Itoa(id)
	//return &k
//}


func GetDataFromDB(server * StreamServer) {
	SelectAll()
}