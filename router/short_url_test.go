package router_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"shurl/common"
	"shurl/dao"
	"shurl/router"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RouterSuite struct {
	suite.Suite
	router *gin.Engine
}

type Response struct {
	Message string `json:"message"`
}

const (
	dbName = "router.db"
)

func (s *RouterSuite) setupDB() dao.Database {
	dao.InitDB(dbName)
	return dao.DB
}

func tearDownDB() {
	if err := os.Remove(dbName); err != nil {
		panic(err)
	}
}

func (s *RouterSuite) SetupSuite() {
	log.Println("SetupSuite")
	db := s.setupDB()
	s.router = router.SetupRouter(db)
}

func (s *RouterSuite) TearDownSuite() {
	log.Println("TearDownSuite")
	tearDownDB()
}

// test the normal cases:
// 1. generate a batch of random long url, and call the shorten api to get the corresponding short links
// 2. call the lookup api to check that all the short - original url mappings are correct
func (s *RouterSuite) TestGenShortUrlAndLookUpOriginalUrl() {
	type Res1 struct {
		Link string `json:"link"`
	}
	type Res2 struct {
		Url string `json:"url"`
	}
	log.Println("generating short urls")
	const count = 10000
	urls := make([]dao.Url, 0, count)
	for i := 0; i < count; i++ {
		url := dao.Url{
			Original: "http://test.com/" + common.RandString(100),
		}
		urls = append(urls, url)
	}
	for i, url := range urls {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader(fmt.Sprintf("{\"url\":\"%s\"}", url.Original)))
		s.router.ServeHTTP(w, req)
		s.Equal(http.StatusOK, w.Code)
		res := Res1{}
		s.NoError(json.Unmarshal(w.Body.Bytes(), &res))
		s.Equal(common.ShortUrlLen, len(res.Link))
		// store short link for later look up test
		urls[i].Short = res.Link
	}

	log.Println("look up original urls")
	for _, url := range urls {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/url/%s", url.Short), nil)
		s.router.ServeHTTP(w, req)

		s.Equal(http.StatusOK, w.Code)
		res := Res2{}
		s.NoError(json.Unmarshal(w.Body.Bytes(), &res))
		s.Equal(url.Original, res.Url)
	}
}

// test the shorten api expecting bad json error
func (s *RouterSuite) TestGenShortUrlBadJson() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader("bad json"))
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusBadRequest, w.Code)
	res := Response{}
	s.NoError(json.Unmarshal(w.Body.Bytes(), &res))
	s.Equal("invalid character 'b' looking for beginning of value", res.Message)
}

// test the shorten api expecting invalid url error
func (s *RouterSuite) TestGenShortUrlInvalidUrl() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader("{\"url\":\"bad url\"}"))
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusBadRequest, w.Code)
	res := Response{}
	s.NoError(json.Unmarshal(w.Body.Bytes(), &res))
	s.Equal("invalid url", res.Message)
}

// test the lookup api expecting invalid url error (can not past the regex test)
func (s *RouterSuite) TestLookUpOriginalUrlInvalidShortUrl() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/url/invalid_short_url", nil)
	s.router.ServeHTTP(w, req)
	s.Equal(http.StatusBadRequest, w.Code)
	res := Response{}
	s.NoError(json.Unmarshal(w.Body.Bytes(), &res))
	s.Equal("invalid short url", res.Message)
}

// test the shorten api's code paths handling database errors
// use mock to trigger database errors
func (s *RouterSuite) TestGenShortUrlErrors() {
	ctrl := gomock.NewController(s.T())
	db := dao.NewMockDatabase(ctrl)
	defer ctrl.Finish()

	testHTTP := func() {
		router := router.SetupRouter(db)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader("{\"url\":\"www.test.com\"}"))
		router.ServeHTTP(w, req)
		s.Equal(http.StatusInternalServerError, w.Code)
		res := Response{}
		s.NoError(json.Unmarshal(w.Body.Bytes(), &res))
		s.Equal("error has occurred", res.Message)
	}

	s.Run("database error", func() {
		// trigger sqlite3 error handling, which retries short url generation
		c1 := db.EXPECT().Create(gomock.Any()).Return(&gorm.DB{Error: sqlite3.Error{Code: sqlite3.ErrConstraint}})
		c2 := db.EXPECT().Create(gomock.Any()).Return(&gorm.DB{Error: sqlite3.Error{Code: sqlite3.ErrInternal}})
		gomock.InOrder(c1, c2)
		testHTTP()
	})

	s.Run("gorm error", func() {
		db.EXPECT().Create(gomock.Any()).Return(&gorm.DB{Error: errors.New("fake gorm errors")})
		testHTTP()
	})
}

// test the lookup api's code paths handling database errors
// use mock to trigger database errors
func (s *RouterSuite) TestLookUpOriginalUrlErrors() {
	ctrl := gomock.NewController(s.T())
	db := dao.NewMockDatabase(ctrl)
	defer ctrl.Finish()

	router := router.SetupRouter(db)
	req, _ := http.NewRequest("GET", "/api/url/abcdefg", nil)

	s.Run("database errors", func() {
		w := httptest.NewRecorder()
		db.EXPECT().First(gomock.Any()).Return(&gorm.DB{Error: sqlite3.Error{}})
		router.ServeHTTP(w, req)
		s.Equal(http.StatusInternalServerError, w.Code)
		res := Response{}
		s.NoError(json.Unmarshal(w.Body.Bytes(), &res))
		s.Equal("error has occurred", res.Message)
	})

	s.Run("not found", func() {
		w := httptest.NewRecorder()
		db.EXPECT().First(gomock.Any()).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
		router.ServeHTTP(w, req)
		s.Equal(http.StatusNotFound, w.Code)
		res := Response{}
		s.NoError(json.Unmarshal(w.Body.Bytes(), &res))
		s.Equal("url not found", res.Message)
	})

	s.Run("gorm errors", func() {
		w := httptest.NewRecorder()
		db.EXPECT().First(gomock.Any()).Return(&gorm.DB{Error: gorm.ErrNotImplemented})
		router.ServeHTTP(w, req)
		s.Equal(http.StatusInternalServerError, w.Code)
		res := Response{}
		s.NoError(json.Unmarshal(w.Body.Bytes(), &res))
		s.Equal("error has occurred", res.Message)
	})
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(RouterSuite))
}
