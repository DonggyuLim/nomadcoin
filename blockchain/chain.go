package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"nomadcoin/db"
	"nomadcoin/utils"
	"sync"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once

//byte 값을 decoding 해주는 함수
func (b *blockchain) restore(data []byte) {
	utils.HandleErr(gob.NewDecoder(bytes.NewReader(data)).Decode(b))

}

func (b *blockchain) persist() {
	db.SaveBlockChain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			fmt.Printf("NewestHash:%s\nHeight:%d\n", b.NewestHash, b.Height)
			//serch for checkpoint on the db
			checkpoint := db.Checkpoint()
			//db.Blockchain()이 nil 이라면 genesisBlock 생성
			if checkpoint == nil {
				b.AddBlock("Genesis")
			} else {
				fmt.Println("Restoring..")
				b.restore(checkpoint)
			}

			//db가 있다면 restore b from bytes

		})
	}
	fmt.Println(b.NewestHash)
	return b
}
