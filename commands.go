package main

import (
	"github.com/urfave/cli"
	"strconv"
)

var Commands = []cli.Command{
	commandTop,
}

var commandTop = cli.Command{
	Name:   "top",
	Usage:  "list Hacker News top stories",
	Action: doTop,
}

func doTop(c *cli.Context) error {

	num := 10

	param, err := strconv.Atoi(c.Args().First())
	if err == nil && param > 0 && param <= 500 {
		num = param
	}

	return getTopStories(num)
}
