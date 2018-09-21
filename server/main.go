package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type BID struct {
	Id           string  `json:"Id,omitempty"`
	BidPrice     float64 `json:"BidPrice,omitempty"`
	AdvertiserId string  `json:"AdvertiserId,omitempty"`
	Nurl         string  `json:"Nurl,omitempty"`
}

type BID_PARAM struct {
	Id             string
	FloorPrice     int
	DeviceId       string
	MediaId        string
	Timestamp      int64
	OsType         string
	BannerSize     int
	BannerPosition int
	DeviceType     int
}

type WIN_NOTICE struct {
	Id      string
	Price   float64
	IsClick int
}

var bid BID

func main() {
	bid = BID{
		Id:           "asdf",
		BidPrice:     3.14,
		AdvertiserId: "zxcv",
		Nurl:         "aiueo",
	}
	router := mux.NewRouter()
	router.HandleFunc("/bid_req", BidHandler).Methods("POST")
	router.HandleFunc("/win_notice", WinHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func BidHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params BID_PARAM
	err := decoder.Decode(&params)
	if err != nil {
		w.Write([]byte("json decode error" + err.Error() + "\n"))
	}
	// peopleをjsonにエンコードしてwに書き込む？
	json.NewEncoder(w).Encode(bid)
}

func WinHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params WIN_NOTICE
	err := decoder.Decode(&params)
	if err != nil {
		w.Write([]byte("json decode error" + err.Error() + "\n"))
	}
	w.WriteHeader(204)
}
