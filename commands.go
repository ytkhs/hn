package main

import (
	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandTop,
}

var commandTop = cli.Command{
	Name:   "top",
	Usage:  "list Hacker News top stories",
	Action: doTop,
}

func doTop(c *cli.Context) {
	getTopStories()
}
