package libstream

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
)

func TestServerStart(t *testing.T) {
	server := GetTestServer()

	req, _ := http.NewRequest("GET", "/test/status", nil)
	w := httptest.NewRecorder()
	//server.GetRouter().Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	//	t, err := route.GetPathTemplate()
	//	if err != nil {
	//		return err
	//	}
	//	fmt.Println(t)
	//	return nil
	//})
	server.GetRouter().ServeHTTP(w, req)
	assert.Equal(t, http.StatusFound, w.Code)
	//
	//resp := w.Result()
	//body, _ := ioutil.ReadAll(resp.Body)
	//
	//fmt.Println(resp.StatusCode)
	//fmt.Println(resp.Header.Get("Content-Type"))
	//fmt.Println(string(body))


}