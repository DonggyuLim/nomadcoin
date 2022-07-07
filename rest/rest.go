package rest

import (
	"encoding/json"

	"fmt"

	"log"

	"net/http"

	"nomadcoin/blockchain"

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

type errorRespons struct {
	errorMessage string `json:"errorMessage"`
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

			URL: url("/blocks/{hash}"),

			Method: "GET",

			Description: "See A Block",
		},
	}

	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":
		return
		// rw.Header().Add("Content-Type", "application/json")

		// json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())

	case "POST":
		return
		/* var addBlockBody addBlockBody

		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))

		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)

		rw.WriteHeader(http.StatusCreated)*/
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
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
