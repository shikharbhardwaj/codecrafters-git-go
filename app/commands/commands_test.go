package commands_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shikharbhardwaj/codecrafters-git-go/app/commands"
	"github.com/shikharbhardwaj/codecrafters-git-go/app/utils"
	"github.com/urfave/cli/v2"
)

var gitDir string
var app cli.App
var buf *bytes.Buffer

func TestMain(m *testing.M) {
	// Create a temp dir
	utils.InfoLogger.Println("Running main.")

	var err error

	baseDir, err := ioutil.TempDir(os.TempDir(), "git_ditto_test_repo_")

	gitDir = filepath.Join(baseDir, "test")

	defer os.RemoveAll(baseDir)

	buf = bytes.NewBufferString("")

	if err != nil {
		fmt.Println("Could not create temp dir for testing.")
		os.Exit(1)
	}

	workingDir, err := os.Getwd()

	if err != nil {
		fmt.Println("Could not get working dir.")
		os.Exit(1)
	}

	app = cli.App{
		Writer: buf,
	}

	app.Commands = []*cli.Command{
		commands.InitCommand,
		commands.CatFileCommand,
		commands.HashObjectCommand,
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "C",
			Value:       workingDir,
			DefaultText: "",
			Usage:       "Run as if git was started in <path> instead of the current working directory.",
		},
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestInit(t *testing.T) {
	cases := []struct {
		testArgs []string
	}{
		{testArgs: []string{"foo", "init", gitDir}},
	}

	for _, c := range cases {
		app.Run(c.testArgs)

		// check if the required dir structure exists.
		utils.Expect(t, utils.PathExists(filepath.Join(gitDir, ".git")), true)
		utils.Expect(t, utils.PathExists(filepath.Join(gitDir, ".git/HEAD")), true)
		utils.Expect(t, utils.PathExists(filepath.Join(gitDir, ".git/objects")), true)
		utils.Expect(t, utils.PathExists(filepath.Join(gitDir, ".git/refs")), true)
		utils.Expect(t, utils.PathExists(filepath.Join(gitDir, ".git/refs/heads")), true)
		utils.Expect(t, utils.PathExists(filepath.Join(gitDir, ".git/refs/tags")), true)
	}

	t.Cleanup(func() {
		err := os.RemoveAll(gitDir)

		if err != nil {
			fmt.Printf("Could not cleanup after init: %s\n", err.Error())
		}
	})
}

func TestHashObject(t *testing.T) {

	cases := []struct {
		testArgs   []string
		testStdout bool
		stdoutFile string
	}{
		{testArgs: []string{"foo", "init", gitDir}, testStdout: false},
		{testArgs: []string{"foo", "-C", gitDir, "hash-object", "testdata/test.dat"}, testStdout: true, stdoutFile: "testdata/test.sum"},
		{testArgs: []string{"foo", "-C", gitDir, "hash-object", "-w", "testdata/test.dat"}, testStdout: true, stdoutFile: "testdata/test.sum"},
	}

	for _, c := range cases {
		err := app.Run(c.testArgs)

		utils.Expect(t, err, nil)

		if c.testStdout {
			utils.ExpectFileContent(t, c.stdoutFile, buf.String())
		}

		buf.Reset()
	}

	t.Cleanup(func() {
		err := os.RemoveAll(gitDir)

		if err != nil {
			fmt.Printf("Could not cleanup after init: %s\n", err.Error())
		}
	})
}

func TestCatFile(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/test.sum")

	if err != nil {
		t.Fatal("Could not load sha1 sum for test file.")
	}

	data = bytes.Replace(data, []byte("\r\n"), []byte("\n"), -1)
	data = bytes.Replace(data, []byte("\r\n"), []byte("\n"), -1)

	sha := strings.TrimSpace(string(data))

	cases := []struct {
		testArgs   []string
		testStdout bool
		stdoutFile string
	}{
		{testArgs: []string{"foo", "init", gitDir}},
		{testArgs: []string{"foo", "-C", gitDir, "hash-object", "-w", "testdata/test.dat"}},
		{testArgs: []string{"foo", "-C", gitDir, "cat-file", "-p", sha}, testStdout: true, stdoutFile: "testdata/test.dat"},
	}

	for _, c := range cases {
		err := app.Run(c.testArgs)

		utils.Expect(t, err, nil)

		if c.testStdout {
			utils.ExpectFileContent(t, c.stdoutFile, buf.String())
		}

		buf.Reset()
	}

	t.Cleanup(func() {
		err := os.RemoveAll(gitDir)

		if err != nil {
			fmt.Printf("Could not cleanup after init: %s\n", err.Error())
		}
	})
}
