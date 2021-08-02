package main

import (
	"os"

	"github.com/urfave/cli/v2"

	commands "github.com/shikharbhardwaj/codecrafters-git-go/app/commands"
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
		commands.InitCommand,
		commands.CatFileCommand,
		commands.HashObjectCommand,
	}

	app.Run(os.Args)
}
