package main

import (
	"shurl/dao"
	"shurl/router"
)

func main() {
	dao.InitDB("main.db")
	r := router.SetupRouter(dao.DB)
	r.Run()
}
