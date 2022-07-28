package p2p

import (
	"fmt"
	"net/http"
	"nomadcoin/blockchain"
	"nomadcoin/utils"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	//port :3000 will upgrade the request from 4000

	openPort := r.URL.Query().Get("openPort")
	ip := utils.IPSplitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" || ip != ""
	}
	fmt.Printf("%s,wants an upgrade\n", openPort)
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	peer := initPeer(conn, ip, openPort)
	time.Sleep(10 * time.Second)

	peer.inbox <- []byte("hello from 3000!")
}

func AddPeer(address, port, openPort string, broadcast bool) {
	//port 4000 is requesting an upgrade from the port :3000
	fmt.Printf("%s want to connect to port %s\n", openPort, port)
	url := fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	utils.HandleErr(err)
	peer := initPeer(conn, address, port)
	if !broadcast {
		broadcastNewPeer(peer)
		return
	} else {
		sendNewestBlock(peer)
	}

}

func BroadcastNewBlock(b *blockchain.Block) {
	for _, peer := range Peers.peerMap {
		notifyNewBlock(b, peer)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	for _, p := range Peers.peerMap {
		notifyNewTx(tx, p)
	}
}

func broadcastNewPeer(newPeer *peer) {
	for key, p := range Peers.peerMap {
		if key != newPeer.key {
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port)
			notifyNewPeer(payload, p)
		}
	}
}
