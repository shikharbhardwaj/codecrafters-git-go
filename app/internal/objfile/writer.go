package objfile

import (
	"bytes"
	"compress/zlib"
	"io"
	"strconv"

	errors "github.com/shikharbhardwaj/codecrafters-git-go/app/errors"
	"github.com/shikharbhardwaj/codecrafters-git-go/app/internal/plumbing"
)

// Writer is responsible for taking in a stream of bytes
// computing its hash and (optionally) storing it in the
// object store, if asked to do so at construction.
type Writer struct {
	raw io.Writer

	zlib   io.WriteCloser
	multi  io.Writer
	hasher plumbing.Hasher

	closed    bool
	remaining int64
}

func NewWriter(w io.Writer) (*Writer, error) {
	return &Writer{
		raw:  w,
		zlib: zlib.NewWriter(w),
	}, nil
}

func getHeaderBytes(t GitObjectType, size int64) []byte {
	buf := bytes.NewBufferString("")

	buf.Write(t.Bytes())
	buf.Write([]byte{' '})
	buf.Write([]byte(strconv.FormatInt(size, 10)))
	buf.Write([]byte{0})

	return buf.Bytes()
}

func (w *Writer) WriteHeader(t GitObjectType, size int64) error {
	if !t.Valid() {
		return errors.GitError{Message: "Invalid object type"}
	}

	if size < 0 {
		return errors.GitError{Message: "Invalid object size"}
	}

	defer w.prepareForWrite(t, size)

	_, err := w.zlib.Write(getHeaderBytes(t, size))

	return err
}

// Initialize the writer to write content and update the sha sum as it does.
func (w *Writer) prepareForWrite(t GitObjectType, size int64) {
	w.hasher = plumbing.NewHasher(getHeaderBytes(t, size))
	w.remaining = size

	w.multi = io.MultiWriter(w.hasher, w.zlib)
}

func (w *Writer) Hash() plumbing.Hash {
	return w.hasher.Sum()
}

func (w *Writer) Write(p []byte) (n int, err error) {
	if w.closed {
		return 0, errors.GitError{Message: "Attempting to write to a closed writer"}
	}

	// If an overwrite is attempted, write upto the last byte within size and raise an error.
	overwrite := false

	if int64(len(p)) > w.remaining {
		p = p[0:w.remaining]
		overwrite = true
	}

	n, err = w.multi.Write(p)
	w.remaining -= int64(n)

	if err == nil && overwrite {
		err = errors.GitError{Message: "Attempting to write beyond the size of the size of the object writer"}
	}

	return
}

func (w *Writer) Close() error {
	if err := w.zlib.Close(); err != nil {
		return err
	}

	w.closed = true
	return nil
}
