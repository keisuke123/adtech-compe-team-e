package main

import (
	"encoding/json"
	"fmt"
	as "github.com/aerospike/aerospike-client-go"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

type BidParam struct {
	Id             string `json:"id"`
	FloorPrice     int    `json:"floorPrice"`
	DeviceId       string `json:"deviceId"`
	MediaId        string `json:"mediaId"`
	Timestamp      int64  `json:"timeStamp"`
	OsType         string `json:"osType"`
	BannerSize     int    `json:"bannerSize"`
	BannerPosition int    `json:"bannerPosition"`
	DeviceType     int    `json:"deviceType"`
}

type MlParams struct {
	FloorPrice     int     `json:"floorPrice"`
	MediaId        string  `json:"mediaId"`
	Timestamp      int64   `json:"timestamp"`
	OsType         string  `json:"osType"`
	BannerSize     int     `json:"bannerSize"`
	BannerPosition int     `json:"bannerPosition"`
	DeviceType     int     `json:"deviceType"`
	Gender         string  `json:"gender"`
	Age            float64 `json:"age"`
	Income         float64 `json:"income"`
	HasChild       string  `json:"hasChild"`
	IsMarried      string  `json:"isMarried"`
	DeviceId       string  `json:"deviceId"`
	Id             string  `json:"id"`
	AdvId          int     `json:"advId"`
}

type BidResponse struct {
	Id           string  `json:"id"`
	BidPrice     float64 `json:"bidPrice"`
	AdvertisedId string  `json:"advertisedId"`
	Nurl         string  `json:"nurl"`
}

type WinNotice struct {
	Id      string `json:"id"`
	Price   string `json:"price"`
	IsClick string `json:"isClick"`
}

type Score struct {
	advId int
	Score float64
	Ctr   float64
}

var keys []*as.Key
var aerospikeClient *as.Client
var err error
var userDemographics map[string]UserDemographics
var advIds [20]string
var originalBudgets [20]float64

func main() {
	// Aerospike config
	aerospikeClient, err = as.NewClient("35.221.100.18", 3000)
	panicOnError(err)

	defer aerospikeClient.Close()

	fmt.Println("==== Aerospike connected ====")

	// generate keys
	for i := 0; i < 20; i++ {
		tmpKey, tmpErr := as.NewKey("test", "aerospike", i)
		panicOnError(tmpErr)
		keys = append(keys, tmpKey)
	}

	for i := 0; i < 20; i++ {
		rec, err := aerospikeClient.Get(nil, keys[i])
		originalBudgets[i] = rec.Bins["budget"].(float64)
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

	// get the best advId
	var ctrs [20]float64
	var budgetsPercentage [20]float64
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20; i++ {
		ctrs[i] = rand.Float64()
		budgetsPercentage[i] = 100000.0 / originalBudgets[i]
	}

	// 一番いいスコアと会社の情報を得る
	bestScoreInfo := scoring(ctrs, budgetsPercentage)

	// 広告を出す会社のID
	advCompanyId := bestScoreInfo.advId

	// 会社のCPCとりだし
	advInfo, err := aerospikeClient.Get(nil, keys[advCompanyId])
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	ctr := bestScoreInfo.Ctr
	cpc := advInfo.Bins["cpc"].(float64)
	nurl := "http://35.186.252.136:8080/win_notice?advId=" + strconv.Itoa(advCompanyId)

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

// decrease company's budget
func decreaseBudget(advId int, price float64) {
	aerospikeClient.Operate(as.NewWritePolicy(0, 0), keys[advId], as.AddOp(as.NewBin("budget", -price)), as.GetOp())
}

// scoring
func scoring(ctrs [20]float64, balancePercentage [20]float64) Score {
	scores := make([]Score, 20)
	for i := 0; i < len(ctrs); i++ {
		//	TODO: implment
		scores[i].advId = i
		scores[i].Score = ctrs[i] * balancePercentage[i]
		scores[i].Ctr = ctrs[i]
	}
	sort.Slice(scores, func(i, j int) bool {
		if scores[i].Score == scores[j].Score {
			return scores[i].Ctr > scores[i].Ctr
		} else {
			return scores[i].Score > scores[j].Score
		}
	})

	return scores[0]
}
