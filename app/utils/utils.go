package utils

import (
	"fmt"
	"os"
)

func PrintError(message string, commandName string) {
	fmt.Printf("Error executing command %s: %s", commandName, message)
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}

	return true
}
