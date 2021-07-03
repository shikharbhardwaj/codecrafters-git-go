package main

import (
	"fmt"
	"os"
)

func printError(message string, commandName string) {
	fmt.Printf("Error executing command %s: %s", commandName, message)
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}

	return false
}
