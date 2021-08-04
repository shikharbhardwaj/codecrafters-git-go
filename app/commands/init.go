package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	errors "github.com/shikharbhardwaj/codecrafters-git-go/app/errors"
	utils "github.com/shikharbhardwaj/codecrafters-git-go/app/utils"
)

func getDirsToMake() []string {
	return []string{".git", ".git/objects", ".git/refs", ".git/objects/pack", ".git/refs/heads", ".git/refs/tags"}
}

func checkEmptyRepoTarget(targetPath string) error {
	for _, dir := range getDirsToMake() {
		path := filepath.Join(targetPath, dir)
		if utils.PathExists(path) {
			return &errors.PathError{
				Op:   "AlreadyExists",
				Path: path,
				Err:  nil,
			}
		}
	}

	return nil
}

func getTargetDir(c *cli.Context) (string, error) {
	curDir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	targetDir := curDir

	if c.Args().Len() > 0 {
		targetDir, err = filepath.Abs(c.Args().Get(0))

		if err != nil {
			return "", err
		}
	}

	return targetDir, nil
}

func validateInitCommand(c *cli.Context) error {
	// If this is a init command without a target, the .git folder should not exist.
	// Else the .git folder within the target folder should not exist.
	targetDir, err := getTargetDir(c)

	if err != nil {
		return err
	}

	return checkEmptyRepoTarget(targetDir)
}

func initializeGitHead(c *cli.Context) error {
	headFileContents := []byte("ref: refs/heads/master\n")

	targetDir, err := getTargetDir(c)

	if err != nil {
		return err
	}

	headFilePath := filepath.Join(targetDir, ".git/HEAD")

	if err := ioutil.WriteFile(headFilePath, headFileContents, 0644); err != nil {
		return &errors.PathError{
			Op:   "Write",
			Path: headFilePath,
			Err:  err,
		}
	}

	return nil
}

func initializeGitDirs(c *cli.Context) error {
	targetDir, err := getTargetDir(c)

	if err != nil {
		return err
	}

	if c.Args().Len() > 0 {
		err = os.Mkdir(targetDir, 0755)
	}

	if err != nil {
		return err
	}

	for _, dir := range getDirsToMake() {
		dirToMake := filepath.Join(targetDir, dir)

		if err := os.Mkdir(dirToMake, 0755); err != nil {
			return &errors.PathError{
				Op:   "Mkdir",
				Path: dirToMake,
				Err:  err,
			}
		}
	}

	return nil
}

var InitCommand = &cli.Command{
	Name:     "init",
	HelpName: "init",
	Usage:    "Initialize a git repository",

	Before: func(c *cli.Context) error {
		utils.InfoLogger.Println("Validating preconditions for init command.")

		err := validateInitCommand(c)

		if err != nil {
			utils.ErrorLogger.Printf("Error when validating preconditions: %s\n", err.Error())

			utils.PrintError(err.Error(), c.Command.Name)
			cli.Exit(err.Error(), 1)
		}

		return err
	},

	Action: func(c *cli.Context) error {
		utils.InfoLogger.Println("Creating directories.")

		err := initializeGitDirs(c)

		if err != nil {
			utils.ErrorLogger.Printf("Error when creating directories: %s\n", err.Error())
			return err
		}

		utils.InfoLogger.Println("Creating HEAD marker.")

		err = initializeGitHead(c)

		if err != nil {
			utils.ErrorLogger.Printf("Error when creating HEAD marker: %s\n", err.Error())
			return err
		}

		fmt.Fprintln(c.App.Writer, "Initialized git directory")

		return nil
	},
}
