package main

import (
	"db-less-store/configs"
	"db-less-store/pkg/db"
)

func main() {
	conf := configs.LoadConfig()
	_ = db.NewDb(conf)
}
