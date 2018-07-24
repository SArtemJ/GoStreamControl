package libstream

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestServerStart(t *testing.T) {
	server := GetTestServer()

	assert.Equal(t, "byroot", server.RootToken)
	assert.Equal(t, "/test/", server.APIPrefix)

	slice, _ := server.Storage.SelectAll(0, 0)
	assert.Equal(t, 0, len(slice))

	i, _ := server.Storage.Db.C(server.Storage.Collection).Count()
	assert.Equal(t, 0, i)
}

func TestRunStream(t *testing.T) {
	server := GetTestServer()

	//on local-machine localhost:8099
	//on docker 0.0.0.0:8099
	request := fmt.Sprintf("http://0.0.0.0:8099/test/run")
	req, _ := http.NewRequest("GET", request, nil)
	w := httptest.NewRecorder()
	server.GetRouter().ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestActivateStream(t *testing.T) {
	server := GetTestServer()

	request := fmt.Sprintf("http://0.0.0.0:8099/test/run")
	req, _ := http.NewRequest("GET", request, nil)
	w := httptest.NewRecorder()
	server.GetRouter().ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	count, _ := server.Storage.Db.C(server.Storage.Collection).Count()
	assert.Equal(t, 1, count)
}

// func TestConcurrent(t *testing.T) {
// 	wg := &sync.WaitGroup{}
// 	for j := 0; j < 10; j++ {
// 		wg.Add(1)
// 		go func() {
// 			for i := 0; i < 500; i++ {
// 				TestPostData(t)
// 			}
// 			wg.Done()
// 		}()
// 		runtime.Gosched()
// 	}
// 	for j := 0; j < 10; j++ {
// 		wg.Add(1)
// 		go func() {
// 			for i := 0; i < 500; i++ {
// 				TestGetData(t)
// 			}
// 			wg.Done()
// 		}()
// 		runtime.Gosched()
// 	}
// 	wg.Wait()
// }
