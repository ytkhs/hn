package main

import (
	"fmt"
	"os"
	"time"

	h "github.com/ytkhs/hackernews-api-go"
)

type result struct {
	story h.Story
	err   error
}

func getTopStories(num int) error {
	ids, err := h.GetStories("top")
	if err != nil {
		return err
	}

	if num > len(ids) {
		num = len(ids)
	}
	target := ids[:num]

	results := make(chan result)
	for _, id := range target {
		go func(itemid int) {
			story, err := h.GetItem(itemid)
			results <- result{story, err}
		}(id)
	}

	fmt.Println("=== Hacker News Top Stories ===")
	fmt.Println()

	for range target {
		r := <-results
		if r.err != nil {
			fmt.Fprintln(os.Stderr, "skip:", r.err)
			continue
		}
		if r.story.Deleted {
			continue
		}
		fmt.Println(time.Unix(r.story.Time, 0), "by", r.story.By)
		fmt.Println("* ", r.story.Title)
		fmt.Println("> ", r.story.URL)
		fmt.Println()
	}

	return nil
}
