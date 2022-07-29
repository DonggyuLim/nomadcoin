package main

import (
	"nomadcoin/cli"
	"nomadcoin/db"
)

func main() {
	defer db.Close()
	db.InitDB()
	cli.Start()
}
