package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}
type Blockchain struct {
	Blocks []*Block
}

var b *Blockchain
var once sync.Once

func (b *Block) caculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlocks := len(GetBlockchain().Blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockchain().Blocks[totalBlocks-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(GetBlockchain().Blocks) + 1}
	newBlock.caculateHash()
	return &newBlock
}
func (b *Blockchain) AddBlock(data string) {
	b.Blocks = append(b.Blocks, createBlock(data))
}

func GetBlockchain() *Blockchain {
	if b == nil {
		once.Do(func() {
			b = &Blockchain{}
			b.AddBlock("Genesis")
		})

	}
	return b
}

func (b *Blockchain) getLastHash() string {
	if len(b.Blocks) > 0 {
		return b.Blocks[len(b.Blocks)-1].Hash
	}
	return ""
}

func (b *Blockchain) addBlock(data string) {
	newBlock := Block{data, "", b.getLastHash(), len(GetBlockchain().Blocks) + 1}
	hash := sha256.Sum256([]byte(newBlock.Data + newBlock.PrevHash))
	newBlock.Hash = fmt.Sprintf("%x", hash)
	b.Blocks = append(b.Blocks, &newBlock)
}

func (b *Blockchain) listBlocks() {
	for _, block := range b.Blocks {
		fmt.Printf("Data : %s\n ,", block.Data)
		fmt.Printf("Hash : %s\n ,", block.Hash)
		fmt.Printf("Prev Hash : %s\n ,", block.PrevHash)
	}
}

func (b *Blockchain) AllBlocks() []*Block {
	return b.Blocks
}

var ErrNotFound = errors.New("Block not found")

func (b *Blockchain) GetBlock(height int) (*Block, error) {
	if height > len(b.Blocks) {
		return nil, ErrNotFound
	}
	return b.Blocks[height-1], nil
}
