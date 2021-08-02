package utils

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"
)

func Expect(t *testing.T, a interface{}, b interface{}) {
	t.Helper()

	if !reflect.DeepEqual(a, b) {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func ExpectFileContent(t *testing.T, file, expected string) {
	data, err := ioutil.ReadFile(file)
	// Ignore windows line endings
	// TODO: Replace with bytes.ReplaceAll when support for Go 1.11 is dropped
	data = bytes.Replace(data, []byte("\r\n"), []byte("\n"), -1)
	Expect(t, err, nil)
	Expect(t, string(data), expected)
}
