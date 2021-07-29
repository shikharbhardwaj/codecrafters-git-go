package fs

import (
	"path/filepath"

	errors "github.com/shikharbhardwaj/codecrafters-git-go/app/errors"
	utils "github.com/shikharbhardwaj/codecrafters-git-go/app/utils"
)

// Get the path of an object stored in Git, also checking if it exists.
func (g Git) GetObjectPath(objectSha string) (string, error) {
	objectPath := filepath.Join(g.basedir, "objects", objectSha[:2], objectSha[2:])

	if !utils.PathExists(objectPath) {
		return "", &errors.PathError{
			Op:   "stat",
			Path: objectPath,
			Err:  nil,
		}
	}

	return objectPath, nil
}
