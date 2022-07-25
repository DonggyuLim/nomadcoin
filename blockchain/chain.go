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

func persistBlockchain(b *blockchain) {
	db.SaveCheckPoint(utils.ToBytes(b))
}

//block append in blockchain
func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1, difficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDiffculty = block.Difficulty
	persistBlockchain(b)
}
func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transactions...)

	}
	return txs
}

func FindTx(b *blockchain, targetID string) *Tx {
	for _, tx := range Txs(b) {
		if tx.Id == targetID {
			return tx
		}
	}
	return nil
}

//get allBlocks
func Blocks(b *blockchain) []*Block {
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
func recaculateDifficulty(b *blockchain) int {
	allBlocks := Blocks(b)
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
func difficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		//recaculate the difficulty
		return recaculateDifficulty(b)
	} else {
		return b.CurrentDiffculty
	}
}

//unspent txOut
//uhm...
func UTxOutsByAddress(address string, blockchain *blockchain) []*UtxOut {
	var uTxOuts []*UtxOut
	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(blockchain) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}
				if FindTx(b, input.TxId).TxOuts[input.Index].Address == address {
					creatorTxs[input.TxId] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Address == address {
					if _, ok := creatorTxs[tx.Id]; !ok {
						uTxOut := &UtxOut{tx.Id, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}

			}
		}
	}
	return uTxOuts
}

//owner amount
func BlanaceByAddress(address string, blockchain *blockchain) int {
	txOuts := UTxOutsByAddress(address, blockchain)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

func Blockchain() *blockchain {
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		fmt.Printf("NewestHash:%s\nHeight:%d\n", b.NewestHash, b.Height)
		//serch for checkpoint on the db
		checkpoint := db.Checkpoint()
		//db.Blockchain()이 nil 이라면 genesisBlock 생성
		if checkpoint == nil {
			b.AddBlock()
		} else {
			fmt.Println("Restoring..")
			b.restore(checkpoint)
		}

		//db가 있다면 restore b from bytes
	})
	fmt.Println(b.NewestHash)
	return b
}
