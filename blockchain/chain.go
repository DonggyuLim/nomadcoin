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
	NewestHash       string `json:"newestHash"`
	Height           int    `json:"height"`
	CurrentDiffculty int    `json:"currentDifficulty"`
}

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
)

var b *blockchain
var once sync.Once

//byte 값을 decoding 해주는 함수
func (b *blockchain) restore(data []byte) {
	utils.HandleErr(gob.NewDecoder(bytes.NewReader(data)).Decode(b))
}

func (b *blockchain) persist() {
	db.SaveCheckPoint(utils.ToBytes(b))
}

//block append in blockchain
func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDiffculty = block.Difficulty
	b.persist()
}

//get allBlocks
func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

//Difficulty are controled by block generation time.
func (b *blockchain) recaculateDifficulty() int {
	allBlocks := b.Blocks()
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval
	if actualTime < expectedTime {
		return b.CurrentDiffculty + 1
	} else if actualTime > expectedTime {
		return b.CurrentDiffculty - 1
	} else {
		return b.CurrentDiffculty
	}

}

//difficulty control
func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		//recaculate the difficulty
		return b.recaculateDifficulty()
	} else {
		return b.CurrentDiffculty
	}
}

func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
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
