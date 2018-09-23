package main

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

type UserDemographics struct {
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
	raw, err := ioutil.ReadFile(os.Getenv("DEMOGRA_PATH"))
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
