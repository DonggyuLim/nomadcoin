package rest

import (
	"encoding/json"
	"strconv"

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

type addBlockBody struct {
	Message string
}

func documentation(rw http.ResponseWriter, r *http.Request) {

	data := []urlDescription{

		{

			URL: url("/"),

			Method: "GET",

			Description: "See Documentation",
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

			URL: url("/blocks/{height}"),

			Method: "GET",

			Description: "See A Block",
		},
	}

	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		rw.Header().Add("Content-Type", "application/json")

		json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())

	case "POST":

		var addBlockBody addBlockBody

		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))

		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)

		rw.WriteHeader(http.StatusCreated)
	}
}

type errorRespons struct {
	errorMessage string `json:"errorMessage"`
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["height"])
	//strconv = 타입변환해주는 라이브러리
	//https://pkg.go.dev/strconv
	utils.HandleErr(err)
	block, err := blockchain.GetBlockchain().GetBlock(id)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorRespons{fmt.Sprintf("%v", err)})
	} else {
		encoder.Encode(block)
	}

}
func jsonContentTypeMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}
func Start(aPort int) {
	router := mux.NewRouter()
	//모든 요청에 	rw.Header().Add("Content-Type", "application/json") 해줌.

	port = fmt.Sprintf(":%d", aPort)
	router.Use(jsonContentTypeMiddleWare)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks/{height:[0-9]+}", block).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
