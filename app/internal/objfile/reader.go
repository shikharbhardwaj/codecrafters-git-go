package objfile

import (
	"bufio"
	"compress/zlib"
	"io"
	"strconv"
	"strings"

	errors "github.com/shikharbhardwaj/codecrafters-git-go/app/errors"
)

type Reader struct {
	multi io.Reader
	zlib  io.ReadCloser
}

func (r *Reader) Read(p []byte) (n int, err error) {
	return r.multi.Read(p)
}

func (r *Reader) ReadUntil(delim byte) ([]byte, error) {
	bufReader := bufio.NewReader(r)

	return bufReader.ReadBytes(delim)
}

func NewReader(r io.Reader) (*Reader, error) {
	zlib, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &Reader{
		zlib: zlib,
	}, nil
}

func (r *Reader) prepareForRead(t GitObjectType, size int64) {
	r.multi = r.zlib
}

func (r *Reader) Header() (t GitObjectType, size int64, err error) {
	bufReader := bufio.NewReader(r)

	bytes, err := bufReader.ReadBytes(' ')

	if err != nil {
		return
	}

	t, err = DetectObjectType(string(bytes))

	if err != nil {
		return
	}

	bytes, err = bufReader.ReadBytes(0)

	if err != nil {
		return
	}

	size, err = strconv.ParseInt(string(bytes), 10, 64)

	if err != nil {
		err = errors.GitError{Message: "Could not detect object header."}
		return
	}

	defer r.prepareForRead(t, size)

	return
}

func getObjectHeader(r io.Reader) ([]byte, error) {
	bufReader := bufio.NewReader(r)

	bytes, err := bufReader.ReadBytes(0x0)

	if err != nil {
		return "", err
	}

	return bytes, err
}

func getTypeStringFromHeader(header []byte) (string, error) {
	header_string := string(header)

	parts := strings.Split(header_string, " ")

	if len(parts) < 2 {
		return "", errors.GitError{Message: "Could not detect object type from object header."}
	}

	return parts[0], nil
}

func GetObjectType(objectReader io.Reader) (GitObjectType, error) {
	header, err := getObjectHeader(objectReader)

	if err != nil {
		return 0, err
	}

	type_string, err := getTypeStringFromHeader(header)

	return DetectObjectType(type_string)
}