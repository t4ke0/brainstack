package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/TaKeO90/brainstack/cmdhandler"
)

var (
	runtui bool
	file   string
	runcmd bool
)

func main() {
	flag.BoolVar(&runtui, "runtui", false, "use it to run the program in tui mode")
	flag.BoolVar(&runcmd, "runcmd", false, "use this flag to run cli version of the program")
	flag.StringVar(&file, "file", "", "use it to specify the json file that we will use to store your project and their tasks on it")
	flag.Parse()
	if runcmd && file != "" {
		if ext := filepath.Ext(file); ext == ".json" {
			cmdhandler.JSONcmdStream(file)
		} else {
			fmt.Printf("You Need A json File to Start -> want (json) got %s\n", ext)
			os.Exit(1)
		}
	} else if runtui && file != "" {
		tuiPath := "./tui/curse.go"
		cmd := exec.Command("go", "run", tuiPath, file)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
