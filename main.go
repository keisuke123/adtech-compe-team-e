package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type BID struct {
	id           string  `json:"id,omitempty"`
	bidPrice     float64 `json:"bidPrice,omitempty"`
	advertiserId string  `json:"advertiserId,omitempty"`
	nurl         string  `json:"nurl,omitempty"`
}

var bid BID

func main() {
	bid = BID{
		id:           "asdf",
		bidPrice:     3.14,
		advertiserId: "zxcv",
		nurl:         "aiueo",
	}
	router := mux.NewRouter()
	router.HandleFunc("/people/{id}", BidHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func BidHandler(w http.ResponseWriter, r *http.Request) {
	// peopleをjsonにエンコードしてwに書き込む？
	json.NewEncoder(w).Encode(bid)
}
