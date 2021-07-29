package commands

import (
	"os"

	"github.com/urfave/cli/v2"

	errors "github.com/shikharbhardwaj/codecrafters-git-go/app/errors"
	fs "github.com/shikharbhardwaj/codecrafters-git-go/app/internal/fs"
	"github.com/shikharbhardwaj/codecrafters-git-go/app/utils"
)

var CatFileCommand = &cli.Command{
	Name:     "cat-file",
	HelpName: "cat-file",
	Usage:    "Provide content or type and size information for repository objects",

	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "pretty-print",
			Aliases: []string{"p"},
			Value:   false,
			Usage:   "Pretty-print the contents of <object> based on its type.",
		},
	},

	Before: func(c *cli.Context) error {
		utils.InfoLogger.Println("Validating preconditions for init command.")

		if c.Args().Len() < 1 {
			utils.ErrorLogger.Println("Need blob SHA to do a cat-file.")

			cli.Exit("Invalid argument", 1)
		}

		if !c.Bool("pretty-print") {
			err := &errors.GitError{
				Message: "cat-file called without -p flag.",
			}

			utils.ErrorLogger.Println(err.Error())

			cli.Exit(err.Error(), 1)

			return err
		}

		return nil
	},

	Action: func(c *cli.Context) error {
		blobSha := c.Args().Get(0)

		utils.InfoLogger.Printf("cat-file for blob SHA: %s\n", blobSha)

		workingDir, err := os.Getwd()

		if err != nil {
			utils.ErrorLogger.Println(err.Error())

			cli.Exit(err.Error(), 1)
		}

		git, err := fs.FindGit(workingDir)

		if err != nil {
			utils.ErrorLogger.Println(err.Error())

			cli.Exit(err.Error(), 1)

			return err
		}

		err = git.WritePrettyObject(blobSha, os.Stdout)

		if err != nil {
			utils.ErrorLogger.Println(err.Error())

			cli.Exit(err.Error(), 1)

			return err
		}

		return nil
	},
}
