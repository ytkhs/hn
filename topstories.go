package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/k0kubun/pp"
)

type TopStoryIds []int

func getTopStories() {
	getTopIds()
}

func getTopIds() []int {

	println(APIEndpoint + "topstories.json")

	response, err := http.Get(APIEndpoint + "topstories.json")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var data TopStoryIds
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return data
}
