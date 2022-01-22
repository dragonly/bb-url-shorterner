package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"shurl/common"
	"shurl/dao"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type RouterSuite struct {
	suite.Suite
	router *gin.Engine
	urls   []dao.Url
}

const (
	dbName = "router.db"
	count  = 10000
)

func (s *RouterSuite) setupDB() {
	dao.InitDB(dbName)
	s.urls = make([]dao.Url, 0, count)
	for i := 0; i < count; i++ {
		url := dao.Url{
			Original: "http://test.com/" + common.RandString(100),
		}
		s.urls = append(s.urls, url)
		// log.Printf("generated url: %v\n", s.urls)
	}
}

func tearDownDB() {
	if err := os.Remove(dbName); err != nil {
		panic(err)
	}
}

func (s *RouterSuite) SetupSuite() {
	log.Println("SetupSuite")
	s.setupDB()
	s.router = SetupRouter()
}

func (s *RouterSuite) TearDownSuite() {
	log.Println("TearDownSuite")
	tearDownDB()
}

// test success path
func (s *RouterSuite) TestGenShortUrlAndLookUpOriginalUrl() {
	type Res1 struct {
		Link string `json:"link"`
	}
	type Res2 struct {
		Url string `json:"url"`
	}
	log.Println("generating short urls")
	for i, url := range s.urls {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader(fmt.Sprintf("{\"url\":\"%s\"}", url.Original)))
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)
		res := Res1{}
		s.NoError(json.Unmarshal(w.Body.Bytes(), &res))
		s.Equal(common.ShortUrlLen, len(res.Link))
		// store short link for later look up test
		s.urls[i].Short = res.Link
	}

	log.Println("look up original urls")
	for _, url := range s.urls {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/url/%s", url.Short), nil)
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)
		res := Res2{}
		s.NoError(json.Unmarshal(w.Body.Bytes(), &res))
		s.Equal(url.Original, res.Url)
	}
}

func (s *RouterSuite) TestGenShortUrlBadJson() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader("bad json"))
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	type Res struct {
		Message string `json:"message"`
	}
	res := Res{}
	s.NoError(json.Unmarshal(w.Body.Bytes(), &res))
	s.Equal("invalid character 'b' looking for beginning of value", res.Message)
}

func (s *RouterSuite) TestGenShortUrlInvalidUrl() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader("{\"url\":\"bad url\"}"))
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	type Res struct {
		Message string `json:"message"`
	}
	res := Res{}
	s.NoError(json.Unmarshal(w.Body.Bytes(), &res))
	s.Equal("invalid url", res.Message)
}

func (s *RouterSuite) TestLookUpOriginalUrlInvalidShortUrl() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/url/invalid_short_url", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	type Res struct {
		Message string `json:"message"`
	}
	res := Res{}
	s.NoError(json.Unmarshal(w.Body.Bytes(), &res))
	s.Equal("invalid short url", res.Message)
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(RouterSuite))
}
