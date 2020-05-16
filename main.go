package main

import (
	"fmt"
	"os"
	"path/filepath"

	"./cmdhandler"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage ./brainstack <json-file-name>")
		os.Exit(1)
	}
	fileName := os.Args[1]
	if ext := filepath.Ext(fileName); ext == ".json" {
		cmdhandler.JSONcmdStream(fileName)
	} else {
		fmt.Printf("You Need A json File to Start -> want (json) got %s\n", ext)
		os.Exit(1)
	}
}
