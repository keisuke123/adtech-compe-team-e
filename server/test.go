package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	client := &http.Client{}

	// GETでもPOSTでも取得できた
	req, err := http.NewRequest("POST", "http://localhost:8080/people/1", nil)
	if err != nil {
		log.Printf("Faild to create request(AtCoder): %v", err)
	}

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))

}
