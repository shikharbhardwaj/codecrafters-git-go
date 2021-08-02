package objfile

import (
	"fmt"

	errors "github.com/shikharbhardwaj/codecrafters-git-go/app/errors"
)

type GitObjectType int

const (
	Blob GitObjectType = iota
	Tree
	Commit
)

func ObjectTypeToNameMapping() map[GitObjectType]string {
	return map[GitObjectType]string{
		Blob:   "blob",
		Tree:   "tree",
		Commit: "commit",
	}
}

func ObjectNameToTypeMapping() map[string]GitObjectType {
	typeToName := ObjectTypeToNameMapping()

	nameToType := make(map[string]GitObjectType, len(typeToName))

	for k, v := range typeToName {
		nameToType[v] = k
	}

	return nameToType
}

func DetectObjectType(typeString string) (GitObjectType, error) {
	objType, prs := ObjectNameToTypeMapping()[typeString]

	if prs {
		return objType, nil
	} else {
		return 0, errors.GitError{Message: fmt.Sprintf("Unkown object type: '%s'", typeString)}
	}
}

func (t GitObjectType) Valid() bool {
	_, prs := ObjectTypeToNameMapping()[t]

	return prs
}

func (t GitObjectType) String() string {
	return ObjectTypeToNameMapping()[t]
}

func (t GitObjectType) Bytes() []byte {
	return []byte(t.String())
}
