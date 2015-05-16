package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var Commands = []cli.Command{
	commandHoge,
}

var commandHoge = cli.Command{
	Name:  "hoge",
	Usage: "function hoge",
	Description: `
    The quick brown fox jumps over the lazy dog
    The quick brown fox jumps over the lazy dog
    The quick brown fox jumps over the lazy dog
`,
	Action: doHoge,
}

func main() {
	app := cli.NewApp()
	app.Name = "hn"
	app.Usage = "sample cli app"
	app.Commands = Commands
	app.Run(os.Args)
}

func doHoge(c *cli.Context) {
	println("Hello, hn!")
}
