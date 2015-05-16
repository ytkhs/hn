package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
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

func doTop(c *cli.Context) {

	num := 10

	param, err := strconv.Atoi(c.Args().First())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else if param <= 500 {
		num = param
	}

	getTopStories(num)
}
