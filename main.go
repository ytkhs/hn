package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := newApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newApp() *cli.App {

	app := cli.NewApp()
	app.Name = "hn"
	app.Version = Version
	app.Authors = []cli.Author{{Name: "ytkhs"}}
	app.Usage = "A Hacker News CLI reader"
	app.Commands = Commands

	return app
}
