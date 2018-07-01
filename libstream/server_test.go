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
	server.GetRouter().ServeHTTP(w, req)
	assert.Equal(t, http.StatusFound, w.Code)
}