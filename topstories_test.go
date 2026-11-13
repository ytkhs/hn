package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"errors"

	h "github.com/ytkhs/hackernews-api-go"
)

var sampleStories = []rankedStory{
	{rank: 1, Story: h.Story{ID: 1, Title: "First", URL: "https://example.com/1", By: "alice", Score: 100, Descendants: 5, Time: 0}},
	{rank: 3, Story: h.Story{ID: 2, Title: "Ask HN: Second", By: "bob", Score: 50, Descendants: 7, Time: 60}},
}

func TestRenderJSONKeepsRankOrderAndFields(t *testing.T) {
	var buf bytes.Buffer
	if err := renderJSON(&buf, sampleStories); err != nil {
		t.Fatal(err)
	}

	var got []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, buf.String())
	}
	if len(got) != 2 {
		t.Fatalf("got %d items", len(got))
	}

	first := got[0]
	want := map[string]any{
		"rank":     float64(1),
		"id":       float64(1),
		"title":    "First",
		"url":      "https://example.com/1",
		"hn_url":   "https://news.ycombinator.com/item?id=1",
		"by":       "alice",
		"score":    float64(100),
		"comments": float64(5),
		"time":     "1970-01-01T00:00:00Z",
	}
	for k, v := range want {
		if first[k] != v {
			t.Errorf("first[%q] = %v, want %v", k, first[k], v)
		}
	}
	if got[1]["rank"] != float64(3) || got[1]["title"] != "Ask HN: Second" {
		t.Errorf("second item = %v", got[1])
	}
}

func TestRenderJSONOmitsEmptyURL(t *testing.T) {
	var buf bytes.Buffer
	renderJSON(&buf, sampleStories[1:])
	if strings.Contains(buf.String(), `"url"`) {
		t.Errorf("url should be omitted for text posts: %s", buf.String())
	}
}

func TestRenderJSONEmptyIsArray(t *testing.T) {
	var buf bytes.Buffer
	renderJSON(&buf, nil)
	if strings.TrimSpace(buf.String()) != "[]" {
		t.Errorf("got %q, want []", buf.String())
	}
}

func TestRenderTextLayout(t *testing.T) {
	now := time.Unix(2*3600+60, 0)
	var buf bytes.Buffer
	renderText(&buf, []rankedStory{
		{rank: 1, Story: h.Story{ID: 1, Title: "First", URL: "https://www.example.com/a", By: "alice", Score: 100, Descendants: 5, Time: 0}},
		{rank: 12, Story: h.Story{ID: 2, Title: "Ask HN: Second", By: "bob", Score: 1, Descendants: 1, Time: 2 * 3600}},
	}, now, false)

	want := ` 1. First (example.com)
    100 points by alice 2 hours ago | 5 comments
    https://www.example.com/a

12. Ask HN: Second
    1 point by bob 1 minute ago | 1 comment
    https://news.ycombinator.com/item?id=2

`
	if got := buf.String(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderTextWithoutColorHasNoEscapes(t *testing.T) {
	var buf bytes.Buffer
	renderText(&buf, sampleStories, time.Unix(0, 0), false)
	if strings.Contains(buf.String(), "\x1b[") {
		t.Errorf("unexpected escape sequence: %q", buf.String())
	}
}

func TestRenderTextWithColorHighlightsTitle(t *testing.T) {
	var buf bytes.Buffer
	renderText(&buf, sampleStories, time.Unix(0, 0), true)
	if !strings.Contains(buf.String(), "\x1b[1m1. First\x1b[0m") {
		t.Errorf("title should be bold: %q", buf.String())
	}
}

func TestRelativeTime(t *testing.T) {
	now := time.Unix(10*86400, 0)
	cases := []struct {
		ago  time.Duration
		want string
	}{
		{10 * time.Second, "just now"},
		{time.Minute, "1 minute ago"},
		{59 * time.Minute, "59 minutes ago"},
		{time.Hour, "1 hour ago"},
		{23 * time.Hour, "23 hours ago"},
		{24 * time.Hour, "1 day ago"},
		{72 * time.Hour, "3 days ago"},
		{-time.Hour, "just now"},
	}
	for _, c := range cases {
		if got := relativeTime(now.Add(-c.ago), now); got != c.want {
			t.Errorf("relativeTime(-%v) = %q, want %q", c.ago, got, c.want)
		}
	}
}

func TestDomain(t *testing.T) {
	cases := map[string]string{
		"https://www.example.com/a?b=c": "example.com",
		"https://ice-wm.org/":           "ice-wm.org",
		"":                              "",
		"::not a url":                   "",
	}
	for in, want := range cases {
		if got := domain(in); got != want {
			t.Errorf("domain(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestUseColor(t *testing.T) {
	if useColor(false, "") {
		t.Error("non-terminal should not use color")
	}
	if useColor(true, "1") {
		t.Error("NO_COLOR should disable color")
	}
	if !useColor(true, "") {
		t.Error("terminal without NO_COLOR should use color")
	}
}

func TestCollectSkipsDeletedAndErrorsKeepingRank(t *testing.T) {
	results := []result{
		{story: h.Story{ID: 1, Title: "ok"}},
		{story: h.Story{ID: 2, Deleted: true}},
		{err: errTest},
		{story: h.Story{ID: 4, Title: "ok2"}},
	}
	stories, errs := collect(results)
	if len(stories) != 2 || stories[0].ID != 1 || stories[1].ID != 4 {
		t.Errorf("stories = %+v", stories)
	}
	if stories[1].rank != 4 {
		t.Errorf("rank should keep original HN position, got %d", stories[1].rank)
	}
	if len(errs) != 1 {
		t.Errorf("errs = %v", errs)
	}
}

var errTest = errors.New("test")
