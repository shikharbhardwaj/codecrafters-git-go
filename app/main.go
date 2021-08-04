package main

import (
	"fmt"
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

	workingDir, err := os.Getwd()

	if err != nil {
		fmt.Println("Could not get working dir.")
		os.Exit(1)
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "C",
			Value:       workingDir,
			DefaultText: "<current working dir>",
			Usage:       "Run as if git was started in <path>",
		},
	}

	app.Commands = []*cli.Command{
		commands.InitCommand,
		commands.CatFileCommand,
		commands.HashObjectCommand,
		commands.LsTreeCommand,
	}

	app.Run(os.Args)
}
