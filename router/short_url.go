package router

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"

	"shurl/common"
	"shurl/dao"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

var (
	validUrl = regexp.MustCompile(`[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&\/=]*)?`)
)

// genShortUrl generates short url according to the original url
func genShortUrl(url string) string {
	h := md5.New()
	io.WriteString(h, url)
	sum := h.Sum(nil)
	// the md5 sum contains 16 bytes, and base64 encoding translates 3 bytes to 4 characters.
	// we only want 7 characters as bitly, so trimming the input can save some cpu cycles.
	return base64.URLEncoding.EncodeToString(sum[:6])[:common.ShortUrlLen]
}

func registerShortUrlAPIs(r *gin.RouterGroup, db dao.Database) {
	// generate short url
	r.POST("/shorten", func(c *gin.Context) {
		// parse request body
		type Req struct {
			Url string `json:"url"`
		}
		req := Req{}
		if err := c.BindJSON(&req); err != nil {
			log.Printf("failed to bind json: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		if !validUrl.MatchString(req.Url) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid url",
			})
			return
		}

		for {
			randPrefix := common.RandString(6)
			shortUrl := genShortUrl(randPrefix + req.Url)
			urlModel := dao.Url{
				Short:    shortUrl,
				Original: req.Url,
			}
			if err := db.Create(&urlModel).Error; err != nil {
				if sqliteErr, ok := err.(sqlite3.Error); ok {
					if sqliteErr.Code == sqlite3.ErrConstraint {
						// duplicated short url (violates primary key constraint)
						// need to regenerate short url with a new random prefix
						log.Println("duplicated primary key")
						// TODO: limit retry
						continue
					} else {
						// database errors
						log.Printf("Error: database error, %v\n", err)
						c.JSON(http.StatusInternalServerError, gin.H{
							"message": "error has occurred",
						})
					}
				} else {
					// gorm errors
					log.Printf("Error: gorm error, %v\n", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "error has occurred",
					})
				}
				return
			}

			// normal case
			c.JSON(http.StatusOK, gin.H{
				"link": shortUrl,
			})
			return
		}
	})

	// get original url
	r.GET("/url/:short_url", func(c *gin.Context) {
		shortUrl := c.Param("short_url")
		if len(shortUrl) != common.ShortUrlLen {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid short url",
			})
			return
		}

		urlModel := dao.Url{
			Short: shortUrl,
		}
		if err := db.First(&urlModel).Error; err != nil {
			if _, ok := err.(sqlite3.Error); ok {
				// database errors
				log.Printf("Error: database error, %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "error has occurred",
				})
			} else {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, gin.H{
						"message": "url not found",
					})
				} else {
					// other gorm errors
					log.Printf("Error: gorm error, %v\n", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "error has occurred",
					})
				}
			}
			return
		}

		// normal case
		c.JSON(http.StatusOK, gin.H{
			"url": urlModel.Original,
		})
	})
}
