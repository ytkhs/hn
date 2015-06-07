package main

import (
	"fmt"
	h "github.com/qube81/hackernews-api-go"
	"time"
)

func getTopStories(num int) {
	ids, _ := h.GetStories("top")

	articleChan := make(chan h.Story)
	target := ids[0:num]

	for _, id := range target {
		go func(itemid int) {

			article, _ := h.GetItem(itemid)
			articleChan <- article

		}(id)
	}

	fmt.Println("=== Hacker News Top Stories ===")
	fmt.Println()

	for i := 0; i < len(target); i++ {
		if news := <-articleChan; !news.Deleted {
			fmt.Println(time.Unix(news.Time, 0), "by", news.By)
			fmt.Println("* ", news.Title)
			fmt.Println("> ", news.URL)
			fmt.Println()
		}
	}

}
