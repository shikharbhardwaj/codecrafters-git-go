package fs

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	errors "github.com/shikharbhardwaj/codecrafters-git-go/app/errors"
	utils "github.com/shikharbhardwaj/codecrafters-git-go/app/utils"
)

const (
	suffix     = ".git"
	objectPath = "objects"
	packPath   = "pack"
)

type Git struct {
	basedir string
}

// Find the directory containing the Git index folder from a given path.
// Currently, only finds the closest .git directory if it exists in the current tree.
// TODO: Make this support all different ways.
func FindGit(curDir string) (*Git, error) {
	for marker := curDir; marker != "/"; {
		utils.InfoLogger.Printf("Checking directory: %s\n", marker)
		if filepath.Base(marker) == suffix {
			return &Git{
				basedir: marker,
			}, nil
		}

		if utils.PathExists(filepath.Join(marker, suffix)) {
			return &Git{
				basedir: filepath.Join(marker, suffix),
			}, nil
		}

		marker = filepath.Dir(marker)
	}

	return nil, errors.GitError{
		Message: "Not a git repository (or any of the parent directories): .git",
	}
}

func (g Git) GetObjectReader(objectSha string) (io.Reader, error) {
	blobPath, err := g.GetObjectPath(objectSha)

	if err != nil {
		return nil, err
	}

	f, err := os.Open(blobPath)

	if err != nil {
		return nil, err
	}

	fileReader := bufio.NewReader(f)

	if err != nil {
		return nil, err
	}

	return fileReader, nil
}

func (g Git) GetTempObjectFile() (*os.File, error) {
	f, err := ioutil.TempFile(filepath.Join(g.basedir, objectPath, packPath), "tmp_obj_")

	return f, err
}
