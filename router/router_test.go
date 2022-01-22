package router

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"shurl/dao"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RouterSuite struct {
	suite.Suite
	router *gin.Engine
}

const dbName = "router.db"

func setupDB() {
	dao.InitDB(dbName)
	// create dummy entities in database
	if err := dao.DB.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < 10; i++ {
			dummy := dao.Dummy{}
			if err := tx.Create(&dummy).Error; err != nil {
				log.Println("fail")
				return err
			}
			log.Println("success")
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

func tearDownDB() {
	if err := os.Remove(dbName); err != nil {
		panic(err)
	}
}

func (s *RouterSuite) SetupSuite() {
	log.Println("SetupSuite")
	setupDB()
	s.router = SetupRouter()
}

func (s *RouterSuite) TearDownSuite() {
	log.Println("TearDownSuite")
	tearDownDB()
}

func (s *RouterSuite) Test1() {
	log.Println("test1")
	assert := assert.New(s.T())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	s.router.ServeHTTP(w, req)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("pong", w.Body.String())
}
func (s *RouterSuite) Test2() {
	log.Println("test2")
}
func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(RouterSuite))
}
