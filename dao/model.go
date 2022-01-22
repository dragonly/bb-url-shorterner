package dao

import "gorm.io/gorm"

type Dummy struct {
	gorm.Model
}

type Url struct {
	Short    string `gorm:"primarykey"`
	Original string
}
