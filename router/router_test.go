package router_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"shurl/router"

	"github.com/stretchr/testify/assert"
)

func TestStaticRoutes(t *testing.T) {
	assert := assert.New(t)
	// static routes do not use database, so pass nil here
	router := router.SetupRouter(nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusMovedPermanently, w.Code)
	assert.True(reflect.DeepEqual([]string{"/web"}, w.Header()["Location"]))
}
