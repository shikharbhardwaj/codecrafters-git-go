package commands

import (
	"fmt"
	"io"

	errors "github.com/shikharbhardwaj/codecrafters-git-go/app/errors"
	"github.com/shikharbhardwaj/codecrafters-git-go/app/internal/fs"
	"github.com/shikharbhardwaj/codecrafters-git-go/app/internal/objfile"
	"github.com/shikharbhardwaj/codecrafters-git-go/app/internal/tree"
	"github.com/shikharbhardwaj/codecrafters-git-go/app/utils"
	"github.com/urfave/cli/v2"
)

var LsTreeCommand = &cli.Command{
	Name:     "ls-tree",
	HelpName: "ls-tree",
	Usage:    "Lists the contents of a given tree object",

	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "name-only",
			Usage: "List only filenames (instead of the \"long\" output), one per line.",
		},
	},

	Action: func(c *cli.Context) (err error) {
		utils.InfoLogger.Println("Validating preconditions for the init command.")

		workingDir := c.String("C")

		git, err := fs.FindGit(workingDir)

		if err != nil {
			utils.ErrorLogger.Println(err.Error())

			return cli.Exit(err.Error(), 1)
		}

		if c.Args().Len() < 1 {
			err = errors.GitError{Message: "Need tree_sha to list tree contents."}

			return cli.Exit(err.Error(), 1)
		}

		reader, err := git.GetObjectReader(c.Args().First())

		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		objreader, err := objfile.NewReader(reader)

		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		objtype, _, err := objreader.Header()

		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		if objtype != objfile.Tree {
			err = errors.GitError{Message: "Cannot ls-tree a non-tree object"}

			return cli.Exit(err.Error(), 1)
		}

		iterator := tree.TreeEntryIterator(objreader)

		for {
			entry, err := iterator()

			if err != nil {
				if err == io.EOF {
					return nil
				}

				return cli.Exit(err.Error(), 1)
			}

			fmt.Fprintf(c.App.Writer, "Entry with name: %s", entry.Name)
		}
	},
}
