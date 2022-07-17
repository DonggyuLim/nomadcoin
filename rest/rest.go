package rest

import (
	"encoding/json"

	"fmt"

	"log"

	"net/http"

	"nomadcoin/blockchain"
	"nomadcoin/utils"

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
type errorRespons struct {
	ErrorMessage string `json:"errorMessage"`
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
	}

	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		// rw.Header().Add("Content-Type", "application/json")

		json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())

	case "POST":

		blockchain.Blockchain().AddBlock()

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
		errMessage := errorRespons{fmt.Sprintf("%v", err)}
		encoder.Encode(errMessage)

	} else {
		encoder.Encode(block)
	}
}

//response = latest block of blockchain
func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Blockchain())
}

func jsonContentTypeMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}
func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.Blockchain().BlanaceByAddress(address)
		json.NewEncoder(rw).Encode(balanaceResponse{address, amount})
	default:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blockchain().TxOutsByAddress(address)))
	}

}

func Start(aPort int) {
	router := mux.NewRouter()
	//모든 요청에 	rw.Header().Add("Content-Type", "application/json") 해줌.
	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleWare)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status)
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
