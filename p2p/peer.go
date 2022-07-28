package p2p

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type peer struct {
	key     string
	address string
	port    string
	conn    *websocket.Conn
	inbox   chan []byte
}

func (p *peer) close() {
	Peers.mutex.Lock()
	defer Peers.mutex.Unlock()
	p.conn.Close()
	delete(Peers.peerMap, p.key)
}

func (p *peer) read() {
	//delete peer in case of error
	defer p.close()
	for {
		m := Message{}
		err := p.conn.ReadJSON(&m)
		if err != nil {

			break
		}
		handleMsg(&m, p)
	}
}

func (p *peer) write() {
	defer p.close()
	for {
		m, ok := <-p.inbox
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}

type peers struct {
	peerMap map[string]*peer
	mutex   sync.Mutex
}

func AllPeers(p *peers) []string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	keys := []string{}
	for key := range p.peerMap {
		keys = append(keys, key)
	}
	return keys
}

var Peers peers = peers{
	peerMap: make(map[string]*peer),
	mutex:   sync.Mutex{},
}

func initPeer(conn *websocket.Conn, address, port string) *peer {
	Peers.mutex.Lock()
	defer Peers.mutex.Unlock()
	key := fmt.Sprintf("%s:%s", address, port)
	p := &peer{

		conn:    conn,
		inbox:   make(chan []byte),
		address: address,
		key:     key,
		port:    port,
	}
	Peers.peerMap[key] = p
	go p.read()
	go p.write()
	return p
}
