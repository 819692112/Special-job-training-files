package main

import (
	"tgserver/db"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db.TgCreate()
}
