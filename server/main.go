package main

import (
	"encoding/json"
	"fmt"
	as "github.com/aerospike/aerospike-client-go"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type BidParam struct {
	Id             string `json:"id"`
	FloorPrice     int    `json:"floorPrice"`
	DeviceId       string `json:"deviceId"`
	MediaId        int    `json:"mediaId"`
	Timestamp      int64  `json:"timeStamp"`
	OsType         string `json:"osType"`
	BannerSize     int    `json:"bannerSize"`
	BannerPosition int    `json:"bannerPosition"`
	DeviceType     int    `json:"deviceType"`
}

type MlParams struct {
	FloorPrice     int     `json:"floorPrice"`
	MediaId        int     `json:"mediaId"`
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
	AdvId int
	Score float64
	Ctr   float64
}

type CTRs struct {
	Adv01 float64 `json:"adv01"`
	Adv02 float64 `json:"adv02"`
	Adv03 float64 `json:"adv03"`
	Adv04 float64 `json:"adv04"`
	Adv05 float64 `json:"adv05"`
	Adv06 float64 `json:"adv06"`
	Adv07 float64 `json:"adv07"`
	Adv08 float64 `json:"adv08"`
	Adv09 float64 `json:"adv09"`
	Adv10 float64 `json:"adv10"`
	Adv11 float64 `json:"adv11"`
	Adv12 float64 `json:"adv12"`
	Adv13 float64 `json:"adv13"`
	Adv14 float64 `json:"adv14"`
	Adv15 float64 `json:"adv15"`
	Adv16 float64 `json:"adv16"`
	Adv17 float64 `json:"adv17"`
	Adv18 float64 `json:"adv18"`
	Adv19 float64 `json:"adv19"`
	Adv20 float64 `json:"adv20"`
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
	runtime.GOMAXPROCS(8)

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
		price, err := strconv.ParseFloat(winNoticeParams.Price, 64)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		decreaseBudget(advId, price)
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
	}

	jsonBytes, err := json.Marshal(mlParams)
	panicOnError(err)

	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", "http://10.146.0.13:3000/predict/ctr", strings.NewReader(string(jsonBytes)))
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
	var tmpCtrs CTRs
	json.Unmarshal(byteArray, &tmpCtrs)

	fmt.Println("%+v\n", tmpCtrs)
	// TODO: ここでbodyを読んで"CTR"を取り出す
	var ctrs [20]float64
	ctrs[0] = tmpCtrs.Adv01
	ctrs[1] = tmpCtrs.Adv02
	ctrs[2] = tmpCtrs.Adv03
	ctrs[3] = tmpCtrs.Adv04
	ctrs[4] = tmpCtrs.Adv05
	ctrs[5] = tmpCtrs.Adv06
	ctrs[6] = tmpCtrs.Adv07
	ctrs[7] = tmpCtrs.Adv08
	ctrs[8] = tmpCtrs.Adv09
	ctrs[9] = tmpCtrs.Adv10
	ctrs[10] = tmpCtrs.Adv11
	ctrs[11] = tmpCtrs.Adv12
	ctrs[12] = tmpCtrs.Adv12
	ctrs[13] = tmpCtrs.Adv13
	ctrs[14] = tmpCtrs.Adv14
	ctrs[15] = tmpCtrs.Adv15
	ctrs[16] = tmpCtrs.Adv16
	ctrs[17] = tmpCtrs.Adv17
	ctrs[18] = tmpCtrs.Adv18
	ctrs[19] = tmpCtrs.Adv19

	// get the best advId

	var budgetsPercentage [20]float64
	rand.Seed(time.Now().UnixNano())
	currentBalance := getCurrentBalance()
	for i := 0; i < 20; i++ {
		budgetsPercentage[i] = originalBudgets[i] / currentBalance[i]
	}

	// 一番いいスコアと会社の情報を得る
	bestScoreInfo := scoring(ctrs, budgetsPercentage)

	if bestScoreInfo.AdvId == -1 {
		ctx.SetStatusCode(204)
	} else {
		// 広告を出す会社のID
		advCompanyId := bestScoreInfo.AdvId

		// 会社のCPCとりだし
		advInfo, err := aerospikeClient.Get(nil, keys[advCompanyId])
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

		ctr := bestScoreInfo.Ctr
		cpc := advInfo.Bins["cpc"].(float64)
		nurl := "http://35.186.252.136/win_notice?advId=" + strconv.Itoa(advCompanyId)

		bidResponse := BidResponse{
			Id:           bidParams.Id,
			BidPrice:     cpc * ctr * 1000.0,
			AdvertisedId: advIds[advCompanyId],
			Nurl:         nurl,
		}

		// HTTP response
		encode, err := json.Marshal(bidResponse)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

		ctx.SetBody(encode)
	}
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
	var scores []Score

	for i := 0; i < len(ctrs); i++ {
		if balancePercentage[i] <= 0.1 {
			continue
		}
		scores = append(scores, Score{AdvId: i, Score: 0.0, Ctr: ctrs[i]})
	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Ctr > scores[j].Ctr
	})

	if len(scores) == 0 {
		return Score{AdvId: -1, Score: -1, Ctr: 204}
	}
	return scores[0]
}

func getCurrentBalance() [20]float64 {
	var ret [20]float64
	for i := 0; i < 20; i++ {
		rec, err := aerospikeClient.Get(nil, keys[i])
		ret[i] = rec.Bins["budget"].(float64)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		fmt.Printf("%#v\n", *rec)
	}
	return ret
}
