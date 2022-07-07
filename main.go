package main

import (
	"nomadcoin/blockchain"
	"nomadcoin/cli"
)

func main() {
	blockchain.Blockchain()
	cli.Start()
}
