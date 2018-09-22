package main

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"

	"github.com/gorilla/mux"
	. "github.com/aerospike/aerospike-client-go"
)

import (
	as "github.com/aerospike/aerospike-client-go"
	"strconv"
	"os"
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

var (
	host      string = "35.221.183.111"
	port      int    = 8081
	namespace string = "hoge"
	set       string = "users"
)

var keys []*Key
var aerospikeClient *Client
var err error

func main() {
	// Aerospike config
	aerospikeClient, err = as.NewClient("35.221.100.18", 3000)
	panicOnError(err)

	// generate keys
	for i := 0 ; i < 20 ; i++ {
		tmpKey, tmpErr := NewKey("test", "aerospike", i)
		panicOnError(tmpErr)
		keys = append(keys, tmpKey)
	}

	for i:=0; i < 20 ; i++ {
		rec, _ := aerospikeClient.Get(nil, keys[i])
		fmt.Printf("%#v\n", *rec)
	}

	bid = BID {
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

func printError(format string, a ...interface{}) {
	fmt.Printf("error: "+format+"\n", a...)
	os.Exit(1)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func decreaseBudget(advId int, price float64) {
	aerospikeClient.Operate(NewWritePolicy(0, 0), keys[advId], as.AddOp(as.NewBin("budget", price)), as.GetOp())
}
