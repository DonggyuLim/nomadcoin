package p2p

import (
	"fmt"
	"net/http"
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
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	peer := initPeer(conn, ip, openPort)
	time.Sleep(10 * time.Second)

	peer.inbox <- []byte("hello from 3000!")
}

func AddPeer(address, port, openPort string) {
	//port 4000 is requesting an upgrade from the port :3000

	url := fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[0:4])
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	utils.HandleErr(err)
	peer := initPeer(conn, address, port)
	time.Sleep(10 * time.Second)
	peer.inbox <- []byte("hello from 4000")
	conn.WriteMessage(websocket.TextMessage, []byte("Hello from port 4000!"))
}
