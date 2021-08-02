package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	errors "github.com/shikharbhardwaj/codecrafters-git-go/app/errors"
	fs "github.com/shikharbhardwaj/codecrafters-git-go/app/internal/fs"
	"github.com/shikharbhardwaj/codecrafters-git-go/app/internal/objfile"
	"github.com/shikharbhardwaj/codecrafters-git-go/app/utils"
)

var HashObjectCommand = &cli.Command{
	Name:     "hash-object",
	HelpName: "hash-object",
	Usage:    "Compute object ID and optionally creates a blob from a file",

	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "w",
			Value: false,
			Usage: "Actually write the object into the object database",
		},
		&cli.StringFlag{
			Name:  "t",
			Value: "blob",
			Usage: "Specify the type",
		},
	},
	Action: func(c *cli.Context) error {
		utils.InfoLogger.Println("Validating preconditions for init command.")

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

		// Check if the file arg is supplied.
		if c.Args().Len() < 1 {
			err = errors.GitError{Message: "Need file-path to hash object"}

			return cli.Exit(err.Error(), 1)
		}

		path := c.Args().First()

		// Check if the object file path exists
		if !utils.PathExists(path) {
			err = errors.GitError{Message: fmt.Sprintf("Object not found at path %s", path)}

			return cli.Exit(err.Error(), 1)
		}

		f, err := os.Open(path)

		if err != nil {
			return err
		}

		defer f.Close()

		rawWriter := ioutil.Discard
		tempLocation := ""
		var tempFile *os.File

		if c.Bool("w") {
			tempFile, err = git.GetTempObjectFile()

			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			tempLocation = tempFile.Name()

			rawWriter = tempFile

			if err != nil {
				return cli.Exit(err.Error(), 1)
			}
		}

		objWriter, err := objfile.NewWriter(rawWriter)

		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		defer objWriter.Close()

		fi, err := f.Stat()

		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		objWriter.WriteHeader(objfile.Blob, fi.Size())

		io.Copy(objWriter, f)

		hash := objWriter.Hash().String()

		fmt.Fprintln(c.App.Writer, hash)

		if c.Bool("w") {
			err = objWriter.Close()

			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			err = tempFile.Close()

			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			objectPath := git.ComputeObjectPath(hash)
			objectDir := filepath.Dir(objectPath)

			err := os.MkdirAll(objectDir, os.ModePerm)

			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			err = os.Rename(tempLocation, objectPath)

			if err != nil {
				return cli.Exit(err.Error(), 1)
			}
		}

		return nil
	},
}
