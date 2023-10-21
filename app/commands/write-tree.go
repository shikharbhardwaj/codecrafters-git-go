package commands

import (
	"os"
    "path/filepath"

	"github.com/shikharbhardwaj/codecrafters-git-go/app/internal/fs"
	"github.com/shikharbhardwaj/codecrafters-git-go/app/internal/objfile"
	"github.com/shikharbhardwaj/codecrafters-git-go/app/utils"
	"github.com/urfave/cli/v2"
)

var WriteTreeCommand = &cli.Command{
	Name:     "write-tree",
	HelpName: "write-tree",
	Usage:    "Write the contents of the current working directory as a tree object",

	Flags: []cli.Flag{},

	Action: func(c *cli.Context) (err error) {
		utils.InfoLogger.Println("Validating preconditions for the write-tree command.")

		workingDir := c.String("C")

		git, err := fs.FindGit(workingDir)

		if err != nil {
			utils.ErrorLogger.Println(err.Error())

			return cli.Exit(err.Error(), 1)
		}


        // 1. Iterate over the current directory to construct the tree object.
        //    The tree object is a file containing one entry for each blob in
        //    the working dir. For eg.
        //    
        //        100644 blob 4aab5f560862b45d7a9f1370b1c163b74484a24d    LICENSE.txt
        //        100644 blob 43ab992ed09fa756c56ff162d5fe303003b5ae0f    README.md
        //        100644 blob c10cb8bc2c114aba5a1cb20dea4c1597e5a3c193    pygit.py
        entries, err := git.GetTreeEntries()
        if err != nil {
            utils.ErrorLogger.Printf("Failed to get tree entries, err=%v\n", err)
            return cli.Exit(err.Error(), 1)
        }

        // 2. Write the tree object to a tempfile.
        tempfile, err := git.GetTempObjectFile()
        if err != nil {
			utils.ErrorLogger.Println(err.Error())

			return cli.Exit(err.Error(), 1)
        }

        objWriter, err := objfile.NewWriter(tempfile)
        if err != nil {
			utils.ErrorLogger.Println(err.Error())

			return cli.Exit(err.Error(), 1)
        }

        utils.InfoLogger.Printf("Got %d entries\n", len(entries))

        treeObjectSize := 0

        for _, entry := range(entries) {
            treeObjectSize += len([]byte(entry.String(false)))
        }

        err = objWriter.WriteHeader(objfile.Tree, int64(treeObjectSize))
        if err != nil {
            utils.ErrorLogger.Println(err.Error())

            return cli.Exit(err.Error(), 1)
        }

        for _, entry := range(entries) {
            entryBytes := []byte(entry.String(false))
            objWriter.Write(entryBytes)
        }

        hash := objWriter.Hash().String()

        err = objWriter.Close()
        if err != nil {
            return cli.Exit(err.Error(), 1)
        }

        err = tempfile.Close()
        if err != nil {
            return cli.Exit(err.Error(), 1)
        }

        // 3. Move the object to the store.
        objectPath := git.ComputeObjectPath(hash)
        objDir := filepath.Dir(objectPath)

        err = os.MkdirAll(objDir, os.ModePerm)
        if err != nil {
            return cli.Exit(err.Error(), 1)
        }

        err = os.Rename(tempfile.Name(), objectPath)
        if err != nil {
            return cli.Exit(err.Error(), 1)
        }

        return nil
	},
}
