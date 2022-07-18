package blockchain

import (
	"errors"
	"nomadcoin/utils"
	"time"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOut"`
}

type TxIn struct {
	TxId  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}
type UtxOut struct {
	TxID   string
	Index  int
	Amount int
}
type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

func isOnMempool(uTxOut *UtxOut) bool {

	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxId == uTxOut.TxID && input.Index == uTxOut.Index {
				return true
			}

		}
	}
	return false
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "CoinBase"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	//check balance
	blockchain := Blockchain()
	if BlanaceByAddress(from, blockchain) < amount {
		return nil, errors.New("Not enought funds")
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutsByAddress(from, blockchain)
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("gyu", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

//miner will confirm transaction in mempool
//mempool must empty
func (m *mempool) TxToConfirm() []*Tx {
	//address changed necessary
	coinbase := makeCoinbaseTx("gyu")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
