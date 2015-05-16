package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

type topStoryIds []int

type news struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	CommentIds  []int  `json:"kids"`
	Score       int    `json:"score"`
	Text        string `json:"text"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

func display(item news) {
	fmt.Printf("* %s\n", item.Title)
	fmt.Printf("→ %s\n", item.URL)
	fmt.Println()
}

func getTopStories(num int) {

	fmt.Println("=== Hacker News Top Stories ===")
	fmt.Println()

	newsChan := make(chan news)

	for _, id := range getTopIds()[0:num] {
		go func(id int) {
			item := getNews(id)
			newsChan <- item
		}(id)
	}

	for i := 0; i < num; i++ {
		display(<-newsChan)
	}

}

func getTopIds() []int {

	body := getJSON(APIEndpoint + "topstories.json")

	var data topStoryIds
	err := json.NewDecoder(body).Decode(&data)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return data
}

func getNews(id int) news {

	body := getJSON(APIEndpoint + "item/" + strconv.Itoa(id) + ".json")

	var data news
	err := json.NewDecoder(body).Decode(&data)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return data

}
