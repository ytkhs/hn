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

var Commands = []*cli.Command{
	commandTop,
}

var commandTop = &cli.Command{
	Name:      "top",
	Usage:     "list Hacker News top stories",
	ArgsUsage: "[num]",
	Action:    doTop,
}

func doTop(ctx context.Context, cmd *cli.Command) error {
	return getTopStories(parseNum(cmd.Args().First()))
}

func parseNum(arg string) int {
	n, err := strconv.Atoi(arg)
	if err != nil || n <= 0 || n > maxNum {
		return defaultNum
	}
	return n
}
