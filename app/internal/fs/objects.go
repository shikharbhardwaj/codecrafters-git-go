package fs

import (
	"path/filepath"

	errors "github.com/shikharbhardwaj/codecrafters-git-go/app/errors"
	utils "github.com/shikharbhardwaj/codecrafters-git-go/app/utils"
)

func (g Git) ComputeObjectPath(objectSha string) string {
	return filepath.Join(g.basedir, objectPath, objectSha[:2], objectSha[2:])
}

// Get the path of an object stored in Git, also checking if it exists.
func (g Git) GetObjectPath(objectSha string) (string, error) {
	objectPath := g.ComputeObjectPath(objectSha)

	if !utils.PathExists(objectPath) {
		return "", &errors.PathError{
			Op:   "stat",
			Path: objectPath,
			Err:  nil,
		}
	}

	return objectPath, nil
}
