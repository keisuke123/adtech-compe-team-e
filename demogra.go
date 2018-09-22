package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type USER struct {
	Gender    string
	Age       float64
	Income    float64
	HasChild  string
	IsMarried string
}

// エラー処理
func failOnError(err error) {
	if err != nil {
		log.Print("Error: ", err)
	}
}

// gobファイルを保存
func storeGob(data interface{}) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		log.Print(err)
	}
	err = ioutil.WriteFile("demogra", buffer.Bytes(), 0600)
	if err != nil {
		log.Print(err)
	}
}

// gobを読み出す
func loadGob(data interface{}) {
	raw, err := ioutil.ReadFile("demogra")
	if err != nil {
		log.Print(err)
	}
	buffer := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(data)
	if err != nil {
		log.Print(err)
	}
}

func main() {
	// csvを読む
	file, err := os.Open("demogra_cleaned.csv")
	if err != nil {
		failOnError(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// mapにする
	users := make(map[string]USER)
	for {
		record, err := reader.Read() // 1行読み出す
		if err == io.EOF {
			break
		} else {
			failOnError(err)
		}

		// 番号がないデータは飛ばす
		if record[0] == "" {
			continue
		}

		/*
			age := record[2]
			income := record[3]
			female := record[4]
			male := record[5]
			notMarried := record[6]
			married := record[7]
			no := record[8]
			yes := record[9]
		*/

		id := record[1]
		ageFloat, _ := strconv.ParseFloat(record[2], 32)
		incomeFloat, _ := strconv.ParseFloat(record[3], 32)
		// ageInt := int(ageFloat)
		gender := "male"
		if record[4] == "1" {
			gender = "female"
		}
		hasChild := "yes"
		if record[8] == "1" {
			hasChild = "no"
		}
		isMarried := "yes"
		if record[6] == "1" {
			isMarried = "no"
		}

		// bool値、両方0の場合はデフォルトでYesにしてるけどいいのか？
		users[id] = USER{
			Gender:    gender,
			Age:       ageFloat,
			Income:    incomeFloat,
			HasChild:  hasChild,
			IsMarried: isMarried,
		}
	}

	// gobファイルに保存する
	storeGob(users)

	// Gobを読み込んでtmpにいれる
	var tmp map[string]USER
	loadGob(&tmp)
	// fmt.Println(tmp) // gobファイルから変換されたmapが表示される

	// 文字列型でIDを指定すればその他の値を取り出せる
	fmt.Println(tmp["b95cbd52-a9c6-489e-a7e6-102212058bb6"])

}
