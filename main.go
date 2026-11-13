package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	if err := newApp().Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newApp() *cli.Command {
	return &cli.Command{
		Name:      "hn",
		Version:   Version,
		Authors:   []any{"ytkhs"},
		Usage:     "A Hacker News CLI reader",
		UsageText: "hn [--json] [num]\nhn top [--json] [num]",
		Flags:     topFlags(),
		Action:    doTop,
		Commands:  Commands,
	}
}
