package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	h "github.com/ytkhs/hackernews-api-go"
)

type result struct {
	story h.Story
	err   error
}

type rankedStory struct {
	rank int
	h.Story
}

type jsonStory struct {
	Rank     int    `json:"rank"`
	ID       int    `json:"id"`
	Title    string `json:"title"`
	URL      string `json:"url,omitempty"`
	HNURL    string `json:"hn_url"`
	By       string `json:"by"`
	Score    int    `json:"score"`
	Comments int    `json:"comments"`
	Time     string `json:"time"`
}

func getTopStories(num int, asJSON bool) error {
	ids, err := h.GetStories("top")
	if err != nil {
		return err
	}

	if num > len(ids) {
		num = len(ids)
	}

	stories, errs := collect(fetchItems(ids[:num]))
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, "skip:", err)
	}

	if asJSON {
		return renderJSON(os.Stdout, stories)
	}
	renderText(os.Stdout, stories, time.Now(), useColor(stdoutIsTerminal(), os.Getenv("NO_COLOR")))
	return nil
}

// fetchItems fetches items concurrently and returns results in the same order as ids.
func fetchItems(ids []int) []result {
	results := make([]result, len(ids))
	done := make(chan struct{})
	for i, id := range ids {
		go func(i, id int) {
			story, err := h.GetItem(id)
			results[i] = result{story, err}
			done <- struct{}{}
		}(i, id)
	}
	for range ids {
		<-done
	}
	return results
}

// collect keeps each story's original HN rank even when earlier items are skipped.
func collect(results []result) ([]rankedStory, []error) {
	var stories []rankedStory
	var errs []error
	for i, r := range results {
		if r.err != nil {
			errs = append(errs, r.err)
			continue
		}
		if r.story.Deleted || r.story.Dead {
			continue
		}
		stories = append(stories, rankedStory{rank: i + 1, Story: r.story})
	}
	return stories, errs
}

const (
	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiDim    = "\x1b[2m"
	ansiOrange = "\x1b[38;5;208m"
)

func renderText(w io.Writer, stories []rankedStory, now time.Time, color bool) {
	paint := func(code, s string) string {
		if !color {
			return s
		}
		return code + s + ansiReset
	}

	width := 1
	for _, s := range stories {
		if n := len(strconv.Itoa(s.rank)); n > width {
			width = n
		}
	}
	indent := strings.Repeat(" ", width+2)

	for _, s := range stories {
		title := paint(ansiBold, fmt.Sprintf("%*d. %s", width, s.rank, s.Title))
		if d := domain(s.URL); d != "" {
			title += " " + paint(ansiDim, "("+d+")")
		}

		points := paint(ansiOrange, plural(s.Score, "point"))
		meta := paint(ansiDim, fmt.Sprintf(" by %s %s | %s",
			s.By, relativeTime(time.Unix(s.Time, 0), now), plural(s.Descendants, "comment")))

		link := s.URL
		if link == "" {
			link = hnURL(s.ID)
		}

		fmt.Fprintln(w, title)
		fmt.Fprintln(w, indent+points+meta)
		fmt.Fprintln(w, indent+link)
		fmt.Fprintln(w)
	}
}

func relativeTime(t, now time.Time) string {
	d := now.Sub(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return plural(int(d/time.Minute), "minute") + " ago"
	case d < 24*time.Hour:
		return plural(int(d/time.Hour), "hour") + " ago"
	default:
		return plural(int(d/(24*time.Hour)), "day") + " ago"
	}
}

func plural(n int, unit string) string {
	if n == 1 {
		return "1 " + unit
	}
	return strconv.Itoa(n) + " " + unit + "s"
}

func domain(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(u.Hostname(), "www.")
}

func hnURL(id int) string {
	return "https://news.ycombinator.com/item?id=" + strconv.Itoa(id)
}

// useColor follows the NO_COLOR convention (https://no-color.org/).
func useColor(isTerminal bool, noColor string) bool {
	return isTerminal && noColor == ""
}

func stdoutIsTerminal() bool {
	fi, err := os.Stdout.Stat()
	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}

func renderJSON(w io.Writer, stories []rankedStory) error {
	out := make([]jsonStory, 0, len(stories))
	for _, s := range stories {
		out = append(out, jsonStory{
			Rank:     s.rank,
			ID:       s.ID,
			Title:    s.Title,
			URL:      s.URL,
			HNURL:    hnURL(s.ID),
			By:       s.By,
			Score:    s.Score,
			Comments: s.Descendants,
			Time:     time.Unix(s.Time, 0).UTC().Format(time.RFC3339),
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
