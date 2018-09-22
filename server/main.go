package main

import (
	"encoding/json"
	"fmt"
	as "github.com/aerospike/aerospike-client-go"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"strconv"
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
	DeviceId       string
	Id             string
	AdvId          int
}

type BidResponse struct {
	Id           string
	BidPrice     float64
	AdvertisedId string
	Nurl         string
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
var advIds [20]string

func main() {
	// Aerospike config
	aerospikeClient, err = as.NewClient("35.221.100.18", 3000)
	panicOnError(err)

	defer aerospikeClient.Close()

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

	advIds[0] = "adv01"
	advIds[1] = "adv02"
	advIds[2] = "adv03"
	advIds[3] = "adv04"
	advIds[4] = "adv05"
	advIds[5] = "adv06"
	advIds[6] = "adv07"
	advIds[7] = "adv08"
	advIds[8] = "adv09"
	advIds[9] = "adv10"
	advIds[10] = "adv11"
	advIds[11] = "adv12"
	advIds[12] = "adv13"
	advIds[13] = "adv14"
	advIds[14] = "adv15"
	advIds[15] = "adv16"
	advIds[16] = "adv17"
	advIds[17] = "adv18"
	advIds[18] = "adv19"
	advIds[19] = "adv20"

	loadGob(&userDemographics)

	router := fasthttprouter.New()
	router.POST("/bid_req", bidRequestHandler)
	router.POST("/win_notice", winNoticeHandler)
	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))

	fmt.Println("==== Listen ====")
}

/**
 *    POST /win_notice
 */
func winNoticeHandler(ctx *fasthttp.RequestCtx) {
	var winNoticeParams WinNotice
	fmt.Println(string(ctx.PostBody()))
	err := json.Unmarshal(ctx.PostBody(), &winNoticeParams)
	if err != nil {
		ctx.Write([]byte("json decode error" + err.Error() + "\n"))
	}

	// get a advId from GET parameter
	params := ctx.QueryArgs()
	advId, err := params.GetUint("advId")
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	panicOnError(err)

	fmt.Printf("%+v\n", winNoticeParams)
	fmt.Printf("advId : %d\n", advId)

	// decrease a budget if an Ad was clicked
	if winNoticeParams.IsClick == "1" { // click
		fmt.Println("reduce")
		price, err := strconv.ParseFloat(winNoticeParams.Price, 64)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		decreaseBudget(advId, price)
	} else { // not click
		fmt.Println("not reduce")
	}

	ctx.SetStatusCode(204)
}

/**
 *   POST /bid_req
 */
func bidRequestHandler(ctx *fasthttp.RequestCtx) {
	// get POST parameters
	var bidParams BidParam
	if err := json.Unmarshal(ctx.PostBody(), &bidParams); err != nil {
		ctx.Write([]byte("json decode error" + err.Error() + "\n"))
	}

	targetUserDemographics := userDemographics[bidParams.DeviceId]

	fmt.Printf("%+v\n", targetUserDemographics)

		var mlParams MlParams
		mlParams = MlParams{
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
			DeviceId:       bidParams.DeviceId,
			Id:             bidParams.Id,
			AdvId:          11,
		}

		jsonBytes, err := json.Marshal(mlParams)
		panicOnError(err)
		fmt.Println(mlParams)
		fmt.Println(jsonBytes)
		/*
		httpClient := &http.Client{}
		// TODO: URL変える
		req, err := http.NewRequest("POST", "http://35.231.37.137:3000/predict/ctr", strings.NewReader(string(jsonBytes)))
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
		fmt.Println("CTR")
		fmt.Println(string(byteArray))
		// TODO: ここでbodyを読んで"CTR"を取り出す
	*/
	// scoring(ctr, current balance)
	advCompanyId := 0

	// 会社のCPCとりだし
	advInfo, err := aerospikeClient.Get(nil, keys[advCompanyId])
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	ctr := 0.05
	cpc := advInfo.Bins["cpc"].(float64)
	nurl := "http://localhost:8080/win_notice?advId=" + strconv.Itoa(advCompanyId)

	bidResponse := BidResponse{
		Id:           bidParams.Id,
		BidPrice:     cpc * ctr * 1000.0,
		AdvertisedId: advIds[advCompanyId],
		Nurl:         nurl,
	}

	fmt.Println("bitResponse")

	fmt.Printf("%+v\n\n", bidResponse)

	// HTTP response
	encode, err := json.Marshal(bidResponse)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	ctx.SetBody(encode)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func decreaseBudget(advId int, price float64) {
	aerospikeClient.Operate(NewWritePolicy(0, 0), keys[advId], as.AddOp(as.NewBin("budget", -price)), as.GetOp())
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
