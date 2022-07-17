package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"nomadcoin/db"
	"nomadcoin/utils"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

//DB 에서 받은 byte slice를 Decode 함.
func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

var ErrNotFound = errors.New("block not found")

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

//마이닝 함수
func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)

	for {
		hash := utils.Hash(b)
		b.Timestamp = int(time.Now().Unix())

		fmt.Printf("Hash:%s\nNonce:%d\nTarget:%s\n\n", b.Hash, b.Nonce, target)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash

			break
		} else {
			b.Nonce++
		}
	}
}

func createBlock(prevHash string, height int) *Block {
	block := &Block{

		Hash:         "",
		PrevHash:     prevHash,
		Height:       height,
		Difficulty:   Blockchain().difficulty(),
		Nonce:        0,
		Transactions: []*Tx{makeCoinbaseTx("gyu")},
		//address changed necessary
	}
	block.mine()
	block.persist()
	return block
}
