package cmdhandler

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"../cmdtools"
	"../jsoncnt"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func JSONcmdParser(cmd string, filename string) {
	//main_cmd, mp, mt := cmdtools.ArgParser(cmd)
	//Init cmd args
	cmdtools.InitArg("project")
	cmdtools.InitArg("todo")
	// Parse args
	main_cmd, a := cmdtools.ParseArg(cmd, "todo", "project")

	err := jsoncnt.OpenJSONfile(filename)
	if err == io.EOF {
		fmt.Println("empty")
	} else if err != nil {
		log.Fatal(err)
	}
	//Get args Values
	m := cmdtools.GetValue("project", a)
	m1 := cmdtools.GetValue("todo", a)

	switch main_cmd {
	case "show":
		l := jsoncnt.ShowJSONcnt()
		for _, p := range l {
			fmt.Printf("Project:%s\t Todos:%s\n", p.Project, p.Todos)
		}
		JSONcmdStream(filename)
	case "add":
		if m["project"] != "" && m1["todo"] != "" {
			err := jsoncnt.WriteJSONcnt(filename, m["project"], m1["todo"])
			if err != nil {
				log.Fatal(err)
			}
			JSONcmdStream(filename)
		} else {
			fmt.Println("Wrong Args")
			JSONcmdStream(filename)
		}
	case "clear":
		cmdtools.ClearScreen()
		JSONcmdStream(filename)
	default:
		fmt.Println("No such Command")
		JSONcmdStream(filename)
	}
}

func JSONcmdStream(filename string) {
	fmt.Printf("json# ")
	cmd := bufio.NewScanner(os.Stdin)
	cmd.Scan()
	command := cmd.Text()
	JSONcmdParser(command, filename)
}
