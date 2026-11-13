# hn

## NAME

hn - A Hacker News CLI reader

## INSTALL

`go install github.com/ytkhs/hn@latest`

## COMMANDS

`hn [--json] [num]` (same as `hn top [--json] [num]`)

`hn --help` shows help.

### EXAMPLE

    hn // get 10 top stories(default)
    hn 30 // get 30 top stories(num option max = 500)
    hn --json 5 // output as JSON (for scripts and AI tools)

if you use iTerm2, clicking URL with command(⌘) key opens the news in your browser.

Output is colored when writing to a terminal. Set `NO_COLOR=1` to disable it.

### JSON OUTPUT

`--json` prints an array of stories in HN rank order. Errors go to stderr, so stdout is always valid JSON.

```json
[
  {
    "rank": 1,
    "id": 8863,
    "title": "My YC app: Dropbox - Throw away your USB drive",
    "url": "http://www.getdropbox.com/u/2/screencast.html",
    "hn_url": "https://news.ycombinator.com/item?id=8863",
    "by": "dhouston",
    "score": 111,
    "comments": 71,
    "time": "2007-04-04T19:16:40Z"
  }
]
```

`url` is omitted for text posts such as Ask HN. `rank` is the position on HN, so it may skip numbers when deleted stories are excluded.

## AUTHOR

ytkhs