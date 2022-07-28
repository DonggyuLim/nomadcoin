package p2p

import (
	"encoding/json"
	"fmt"
	"nomadcoin/blockchain"
	"nomadcoin/utils"
	"strings"
)

type MessageKind int

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlockRequest
	MessageAllBlockResponse
	MessageNewBlockNotify
	MessageTransaction
	MessageNewPeer
)

type Message struct {
	Kind    MessageKind
	Payload []byte
}

//Message -> json
func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:    kind,
		Payload: utils.ToJSON(payload),
	}

	return utils.ToJSON(m)
}

func sendNewestBlock(p *peer) {
	fmt.Printf("Sending newest block to %s\n", p.key)
	block, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
	utils.HandleErr(err)
	m := makeMessage(MessageNewestBlock, block)
	p.inbox <- m
}

func requestAllBlocks(p *peer) {
	//요청하는 거니까 payload = nil
	m := makeMessage(MessageAllBlockRequest, nil)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlockResponse, blockchain.Blocks(blockchain.Blockchain()))
	p.inbox <- m
}

func handleMsg(m *Message, p *peer) {
	fmt.Printf("Received the newest block from %s\n", p.key)
	switch m.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		err := json.Unmarshal(m.Payload, &payload)
		utils.HandleErr(err)
		b, err := blockchain.FindBlock(blockchain.Blockchain().NewestHash)
		utils.HandleErr(err)
		if payload.Height >= b.Height {
			fmt.Printf("Requesting all blocks from %s\n", p.key)
			requestAllBlocks(p)
		} else {
			fmt.Printf("Sending newest block to from %s\n", p.key)
			sendNewestBlock(p)
		}
	case MessageAllBlockRequest:
		fmt.Printf("%s wants all the blocks\n", p.key)
		sendAllBlocks(p)
	case MessageAllBlockResponse:
		fmt.Printf("Received all the blocks from %s\n", p.key)
		var payload []*blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().Replace(payload)
	case MessageNewBlockNotify:
		var payload *blockchain.Block
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Blockchain().AddPeerBlock(payload)
	case MessageTransaction:
		var payload *blockchain.Tx
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		blockchain.Mempool().AddPeerTx(payload)
	case MessageNewPeer:
		var payload string
		utils.HandleErr(json.Unmarshal(m.Payload, &payload))
		fmt.Printf("I will now /ws upgrade %s", payload)
		parts := strings.Split(payload, ":")
		AddPeer(parts[0], parts[1], parts[2], false)
	}
}

func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNewBlockNotify, b)
	p.inbox <- m
}

func notifyNewTx(tx *blockchain.Tx, p *peer) {
	m := makeMessage(MessageTransaction, tx)
	p.inbox <- m
}

func notifyNewPeer(address string, p *peer) {
	m := makeMessage(MessageNewPeer, address)
	p.inbox <- m
}
