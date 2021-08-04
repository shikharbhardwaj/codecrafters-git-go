package tree

import (
	"bufio"
	"strconv"

	errors "github.com/shikharbhardwaj/codecrafters-git-go/app/errors"
	"github.com/shikharbhardwaj/codecrafters-git-go/app/internal/objfile"
)

type Entry struct {
	Mode uint32
	Name string
	Sha  []byte
}

func TreeEntryIterator(r *objfile.Reader) func() (Entry, error) {
	scanner := bufio.NewScanner(r)

	tokenType := 0

	tokenScanner := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for _, c := range data {
			advance += 1

			shouldReturn := false

			switch tokenType {
			case 0:
				shouldReturn = (c == ' ')
			case 1:
				shouldReturn = (c == 0x0)
			case 2:
				shouldReturn = (len(token) == 20)
			}

			if shouldReturn {
				return
			}

			token = append(token, c)
		}

		return
	}

	scanner.Split(tokenScanner)

	return func() (entry Entry, err error) {
		tokenType = 0

		if scanner.Scan() {
			rawMode := scanner.Bytes()

			modeVal, err := strconv.ParseInt(string(rawMode), 10, 32)

			if err != nil {
				return entry, err
			}

			entry.Mode = uint32(modeVal)
		} else {
			err = errors.GitError{Message: "Could not find tree entry filemode field"}

			return
		}

		tokenType += 1

		if scanner.Scan() {
			rawName := scanner.Bytes()

			entry.Name = string(rawName)
		} else {
			err = errors.GitError{Message: "Could not find tree entry name field"}

			return
		}

		tokenType += 1

		if scanner.Scan() {
			entry.Sha = scanner.Bytes()
		} else {
			err = errors.GitError{Message: "Could not find tree entry sha field"}

			return
		}

		return
	}
}
