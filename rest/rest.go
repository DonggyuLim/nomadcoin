package rest

import (
	"encoding/json"

	"fmt"

	"log"

	"net/http"

	"nomadcoin/blockchain"
	"nomadcoin/p2p"
	"nomadcoin/utils"
	"nomadcoin/wallet"

	"github.com/gorilla/mux"
)

var port string

type url string

func (u url) MarshalText() ([]byte, error) {

	url := fmt.Sprintf("http://localhost%s%s", port, u)

	return []byte(url), nil

}

type urlDescription struct {
	URL url `json:"url"`

	Method string `json:"method"`

	Description string `json:"description"`

	Payload string `json:"payload,omitempty"`
}
type balanaceResponse struct {
	Address string `json:"address"`
	Blanace int    `json:"blanace"`
}

type addTxPayload struct {
	To     string
	Amount int
}
type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type myWalletResponse struct {
	Address string `json:"address"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {

	data := []urlDescription{

		{

			URL: url("/"),

			Method: "GET",

			Description: "See Documentation",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the blockchain",
		},
		{

			URL: url("/blocks"),

			Method: "GET",

			Description: "See All Blocks",
		},

		{

			URL: url("/blocks"),

			Method: "POST",

			Description: "Add A Block",

			Payload: "data:string",
		},

		{

			URL: url("/blocks/{hash}"),

			Method: "GET",

			Description: "See A Block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
		},
		{
			URL:         url("/mempool"),
			Method:      "GET",
			Description: "mempool view",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to Web Sockets",
		},
	}

	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		// rw.Header().Add("Content-Type", "application/json")
		block := blockchain.Blockchain()
		json.NewEncoder(rw).Encode(blockchain.Blocks(block))

	case "POST":

		newBlock := blockchain.Blockchain().AddBlock()
		p2p.BroadcastNewBlock(newBlock)
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	//mux.Vars = 라우트 변수를 맵핑해줌.

	vars := mux.Vars(r)

	hash := vars["hash"]
	//strconv = 타입변환해주는 라이브러리
	//https://pkg.go.dev/strconv
	block, err := blockchain.FindBlock(hash)

	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		errMessage := errorResponse{fmt.Sprintf("%v", err)}
		encoder.Encode(errMessage)

	} else {
		encoder.Encode(block)
	}
}

//response = latest block of blockchain
func status(rw http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.Blockchain(), rw)
}

func jsonContentTypeMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}

// /balance/{address} -> balance
func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		chain := blockchain.Blockchain()
		amount := blockchain.BlanaceByAddress(address, chain)
		json.NewEncoder(rw).Encode(balanaceResponse{address, amount})
	default:
		chain := blockchain.Blockchain()
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, chain)))
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool().Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload

	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{"not enough funds"})
		return
	}
	go p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)

}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(myWalletResponse{Address: address})
}

type addPeerPayload struct {
	Address, Port string
}

func addPeers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, port[1:], true)
		rw.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	}

}

func Start(aPort int) {
	router := mux.NewRouter()
	//모든 요청에 	rw.Header().Add("Content-Type", "application/json") 해줌.

	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleWare, loggerMiddleware)

	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool)
	router.HandleFunc("/transactions", transactions).Methods("POST")
	router.HandleFunc("/wallet", myWallet)
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/peers", addPeers).Methods("GET", "POST")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
