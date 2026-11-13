package main

import (
	"context"
	"strconv"

	"github.com/urfave/cli/v3"
)

const (
	defaultNum = 10
	maxNum     = 500
)

// runTop is a variable so tests can replace the network call.
var runTop = getTopStories

var Commands = []*cli.Command{
	{
		Name:      "top",
		Usage:     "list Hacker News top stories (default command)",
		ArgsUsage: "[num]",
		Flags:     topFlags(),
		Action:    doTop,
	},
}

func topFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{Name: "json", Usage: "output stories as JSON"},
	}
}

func doTop(ctx context.Context, cmd *cli.Command) error {
	return runTop(parseNum(cmd.Args().First()), cmd.Bool("json"))
}

func parseNum(arg string) int {
	n, err := strconv.Atoi(arg)
	if err != nil || n <= 0 || n > maxNum {
		return defaultNum
	}
	return n
}
