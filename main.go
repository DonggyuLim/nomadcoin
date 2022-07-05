package main

import (
	"nomadcoin/explorer"
	"nomadcoin/rest"
)

func main() {
	rest.Start(4000)
	explorer.Start(3000)

}
