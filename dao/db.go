package dao

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database interface {
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
}

var DB *gorm.DB

func InitDB(filename string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&Url{})
	log.Println("auto migration done")
}
