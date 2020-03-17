package main

import (
	"./flaghandler"
	"flag"
	"fmt"
)

func main() {
	var (
		file  string
		ideas string
		add   bool
		done  bool
		reset bool
	)
	//flags for the program
	flag.StringVar(&file, "f", "", "Specify the File Name Where To put your ideas")
	flag.StringVar(&ideas, "i", "", "add ideas to the file")
	flag.BoolVar(&add, "add", false, "true if you wanna add the ideas")
	flag.BoolVar(&done, "done", false, "true if you finish executing an idea")
	flag.BoolVar(&reset, "reset", false, "reset the file")
	flag.Parse()

	if add && file != "" {
		msg := flaghandler.HandleAddingIdeas(ideas, file)
		fmt.Println(msg)
	} else if done {
		f, m := flaghandler.HandleDoneFlag(file)
		if len(f) != 0 {
			fmt.Printf("Should Specify a file with 'f' flag %v\n", f)
		} else if m != "" {
			fmt.Println(m)
		}
	} else if reset && file != "" {
		m := flaghandler.HandleResetFlag(file)
		fmt.Println(m)
	}
}
