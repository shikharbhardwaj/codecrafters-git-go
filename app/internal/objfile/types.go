package objfile

import errors "github.com/shikharbhardwaj/codecrafters-git-go/app/errors"

type GitObjectType int

const (
	Blob GitObjectType = iota
	Tree
	Commit
)

func DetectObjectType(type_string string) (GitObjectType, error) {
	switch type_string {
	case "blob":
		return Blob, nil
	case "commit":
		return Commit, nil
	case "tree":
		return Tree, nil
	}

	return 0, errors.GitError{Message: "Unkown object type."}
}
