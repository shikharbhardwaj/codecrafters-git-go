package utils

import (
	"path/filepath"

	errors "github.com/shikharbhardwaj/codecrafters-git-go/app/errors"
)

const (
	OBJECT_SHA_ID_PREFIX = 2
	GIT_FOLDER_NAME      = ".git"
)

// Find the directory containing the Git folder from a given path.
// Currently, only finds the closest .git directory if it exists in the current tree.
// TODO: Make this support all different ways.
func FindGit(curDir string) (*Git, error) {
	for marker := curDir; marker != "/"; {
		InfoLogger.Printf("Checking directory: %s\n", marker)
		if filepath.Base(marker) == GIT_FOLDER_NAME {
			return &Git{
				basedir: marker,
			}, nil
		}

		if PathExists(filepath.Join(marker, GIT_FOLDER_NAME)) {
			return &Git{
				basedir: filepath.Join(marker, GIT_FOLDER_NAME),
			}, nil
		}

		marker = filepath.Dir(marker)
	}

	return nil, errors.GitError{
		Message: "Not a git repository (or any of the parent directories): .git",
	}
}

// Get the path of an object stored in Git, also checking if it exists.
func (g Git) GetObjectPath(objectSha string) (string, error) {
	objectPath := filepath.Join(g.basedir, "objects", objectSha[:2], objectSha[2:])

	if !PathExists(objectPath) {
		return "", &errors.PathError{
			Op:   "stat",
			Path: objectPath,
			Err:  nil,
		}
	}

	return objectPath, nil
}
