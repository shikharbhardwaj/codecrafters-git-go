package tree

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
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
		advance = 0

		for _, c := range data {
			shouldReturn := false

			switch tokenType {
			case 0:
				shouldReturn = (c == ' ')
				advance = 1
			case 1:
				shouldReturn = (c == 0x0)
				advance = 1
			case 2:
				shouldReturn = (len(token) == 20)
			}

			if shouldReturn {
				break
			}

			token = append(token, c)
		}

		advance += len(token)

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
			err = scanner.Err()

			if err == nil {
				err = io.EOF
			}

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

func (e *Entry) String(nameOnly bool) string {
	buf := bytes.NewBufferString("")

	if !nameOnly {
		fmt.Fprintf(buf, "%06d\t", e.Mode)
	}

	objectType := objfile.Tree

	if e.Mode >= 100_000 {
		objectType = objfile.Blob
	}

	if !nameOnly {
		fmt.Fprintf(buf, "%s\t", objectType.String())
	}

	if !nameOnly {
		fmt.Fprintf(buf, "%s\t", hex.EncodeToString(e.Sha))
	}

	fmt.Fprint(buf, e.Name)

	return buf.String()
}

