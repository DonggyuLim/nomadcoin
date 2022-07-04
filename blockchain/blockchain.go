package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct {
	Data     string
	hash     string
	prevHash string
}
type blockchain struct {
	blocks []*block
}

var b *blockchain
var once sync.Once

func (b *block) caculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.prevHash))
	b.hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlocks := len(GetBlockchain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().blocks[totalBlocks-1].hash
}

func createBlock(data string) *block {
	newBlock := block{data, "", getLastHash()}
	newBlock.caculateHash()
	return &newBlock
}
func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis")
		})

	}
	return b
}

func (b *blockchain) getLastHash() string {
	if len(b.blocks) > 0 {
		return b.blocks[len(b.blocks)-1].hash
	}
	return ""
}

func (b *blockchain) addBlock(data string) {
	newBlock := block{data, "", b.getLastHash()}
	hash := sha256.Sum256([]byte(newBlock.Data + newBlock.prevHash))
	newBlock.hash = fmt.Sprintf("%x", hash)
	b.blocks = append(b.blocks, &newBlock)
}

func (b *blockchain) listBlocks() {
	for _, block := range b.blocks {
		fmt.Printf("Data : %s\n ,", block.Data)
		fmt.Printf("Hash : %s\n ,", block.hash)
		fmt.Printf("Prev Hash : %s\n ,", block.prevHash)
	}
}

func (b *blockchain) AllBlocks() []*block {
	return b.blocks
}
