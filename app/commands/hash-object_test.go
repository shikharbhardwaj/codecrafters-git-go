package commands

import (
	"bytes"
	"testing"

	"github.com/urfave/cli/v2"

	utils "github.com/shikharbhardwaj/codecrafters-git-go/app/utils"
)

func Test_hash_object_without_write(t *testing.T) {
	cases := []struct {
		testArgs []string
	}{
		{testArgs: []string{"foo", "hash-object", "test_data/test.dat"}},
	}

	for _, c := range cases {
		buf := bytes.NewBufferString("")
		app := &cli.App{
			Writer: buf,
		}

		app.Commands = []*cli.Command{
			HashObjectCommand,
		}

		err := app.Run(c.testArgs)

		utils.Expect(t, err, nil)
		utils.ExpectFileContent(t, "test_data/test.sum", buf.String())
	}
}
