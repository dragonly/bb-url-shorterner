package dao

type Url struct {
	Short    string `gorm:"primarykey"`
	Original string
}
