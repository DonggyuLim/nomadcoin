package blockchain

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"nomadcoin/db"
	"nomadcoin/utils"
	"sync"
)

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
)

type blockchain struct {
	NewestHash       string `json:"newestHash"`
	Height           int    `json:"height"`
	CurrentDiffculty int    `json:"currentDifficulty"`
	mutex            sync.Mutex
}
type Storage interface {
	FindBlock(hash string) []byte
	SaveBlock(hash string, data []byte)
	SaveChain(data []byte)
	LoadChain() []byte
	DeleteAllBlocks()
}

var b *blockchain
var once sync.Once
var dbStorage Storage = db.DB{}

//byte 값을 decoding 해주는 함수
func (b *blockchain) restore(data []byte) {
	utils.HandleErr(gob.NewDecoder(bytes.NewReader(data)).Decode(b))
}

func persistBlockchain(b *blockchain) {
	dbStorage.SaveChain(utils.ToBytes(b))
}

//block append in blockchain
func (b *blockchain) AddBlock() *Block {
	block := createBlock(b.NewestHash, b.Height+1, difficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDiffculty = block.Difficulty
	persistBlockchain(b)
	return block
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
	b.mutex.Lock()
	defer b.mutex.Unlock()
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
		checkpoint := dbStorage.LoadChain()
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

func Status(b *blockchain, rw http.ResponseWriter) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	err := json.NewEncoder(rw).Encode(b)
	utils.HandleErr(err)
}

func (b *blockchain) Replace(newBlocks []*Block) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.CurrentDiffculty = newBlocks[0].Difficulty
	b.Height = len(newBlocks)
	b.NewestHash = newBlocks[0].Hash
	persistBlockchain(b)
	dbStorage.DeleteAllBlocks()
	for _, block := range newBlocks {
		block.persistBlock()
	}
}

func (b *blockchain) AddPeerBlock(block *Block) {
	b.mutex.Lock()
	m.mutex.Lock()
	defer b.mutex.Unlock()
	defer m.mutex.Unlock()
	b.Height += 1
	b.CurrentDiffculty = block.Difficulty
	b.NewestHash = block.Hash
	persistBlockchain(b)
	block.persistBlock()

	// 블록을 검증해야함 그냥 블록에 담겨진 트랜잭션을 멤풀에서 지우는게 아닌
	//검증을 하는 방법을 고민해볼 필요가 있음.
	//그리고 거짓 블록이라면 받으면 안됨.
	for _, tx := range block.Transactions {
		_, ok := m.Txs[tx.Id]
		if ok {
			delete(m.Txs, tx.Id)
		}
	}
}
