package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

const (
	appName = "git-ditto"
)

func main() {
	var app = &cli.App{
		Name:  appName,
		Usage: "Implementation of basic git functionality in Go.",
	}

	app.Commands = []*cli.Command{
		initCommand,
	}

	app.Run(os.Args)
}
