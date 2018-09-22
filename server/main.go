package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	. "github.com/aerospike/aerospike-client-go"
	"github.com/gorilla/mux"
)

import (
	as "github.com/aerospike/aerospike-client-go"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type BID struct {
	Id           string  `json:"Id,omitempty"`
	BidPrice     float64 `json:"BidPrice,omitempty"`
	AdvertiserId string  `json:"AdvertiserId,omitempty"`
	Nurl         string  `json:"Nurl,omitempty"`
}

type BidParam struct {
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

type MlParams struct {
	FloorPrice     int
	MediaId        string
	Timestamp      int64
	OsType         string
	BannerSize     int
	BannerPosition int
	DeviceType     int
	Gender         string
	Age            float64
	Income         float64
	HasChild       string
	IsMarried      string
}

type WinNotice struct {
	Id      string
	Price   float64
	IsClick int
}

var bid BID
var keys []*Key
var aerospikeClient *Client
var err error
var userDemographics map[string]UserDemographics

func main() {
	// Aerospike config
	aerospikeClient, err = as.NewClient("35.221.100.18", 3000)
	panicOnError(err)

	//defer closeAerospile()

	fmt.Println("==== Aerospike connected ====")

	// generate keys
	for i := 0; i < 20; i++ {
		tmpKey, tmpErr := NewKey("test", "aerospike", i)
		panicOnError(tmpErr)
		keys = append(keys, tmpKey)
	}

	for i := 0; i < 20; i++ {
		rec, err := aerospikeClient.Get(nil, keys[i])
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		fmt.Printf("%#v\n", *rec)
	}

	bid = BID{
		Id:           "asdf",
		BidPrice:     3.14,
		AdvertiserId: "zxcv",
		Nurl:         "aiueo",
	}

	loadGob(&userDemographics)

	router := mux.NewRouter()
	router.HandleFunc("/bid_req", BidHandler).Methods("POST")
	router.HandleFunc("/win_notice", WinHandler).Methods("POST")
	router.HandleFunc("/hoge", HogeHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))

	fmt.Println("==== Listen ====")
}

func HogeHandler(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var bidParams BidParam
	err := decoder.Decode(&bidParams)
	if err != nil {
		writer.Write([]byte("json decode error" + err.Error() + "\n"))
	}

	fmt.Println(bidParams)

	json.NewEncoder(writer).Encode(bid)
}

func BidHandler(w http.ResponseWriter, r *http.Request) {
	// get POST parameters
	decoder := json.NewDecoder(r.Body)
	var bidParams BidParam
	err := decoder.Decode(&bidParams)
	if err != nil {
		w.Write([]byte("json decode error" + err.Error() + "\n"))
	}

	targetUserDemographics := userDemographics[bidParams.DeviceId]

	mlParams := MlParams{
		FloorPrice:     bidParams.FloorPrice,
		MediaId:        bidParams.MediaId,
		Timestamp:      bidParams.Timestamp,
		OsType:         bidParams.OsType,
		BannerSize:     bidParams.BannerSize,
		BannerPosition: bidParams.BannerPosition,
		DeviceType:     bidParams.DeviceType,
		Gender:         targetUserDemographics.Gender,
		Age:            targetUserDemographics.Age,
		Income:         targetUserDemographics.Income,
		HasChild:       targetUserDemographics.HasChild,
		IsMarried:      targetUserDemographics.IsMarried,
	}

	jsonBytes, err := json.Marshal(mlParams)
	panicOnError(err)

	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8080/hoge", strings.NewReader(string(jsonBytes)))
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := httpClient.Do(req)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	defer response.Body.Close()

	// convert to Integer for CTR
	byteArray, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(byteArray))

	// HTTP response
	json.NewEncoder(w).Encode(bid)
}

/**
 *  GET    /win_notice?advId=xxxx
 *	 params advId
 */
func WinHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var winNoticeParams WinNotice
	err := decoder.Decode(&winNoticeParams)
	if err != nil {
		w.Write([]byte("json decode error" + err.Error() + "\n"))
	}

	// get a advId from GET parameter
	params := mux.Vars(r)
	advId, err := strconv.Atoi(params["advId"])
	panicOnError(err)

	// decrease a budget
	decreaseBudget(advId, winNoticeParams.Price)

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

func closeAerospile() {
	for i := 0; i < len(keys); i++ {
		_, err := aerospikeClient.Delete(nil, keys[i])
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
	}
}
