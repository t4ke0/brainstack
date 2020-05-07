package cmdhandler

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	//"strings"

	"../cmdtools"
	"../jsoncnt"
	//	"../jsonhandler"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func JSONcmdParser(cmd string, filename string) {
	main_cmd, mp, mt := cmdtools.ArgParser(cmd)
	err := jsoncnt.OpenJSONfile(filename)
	if err == io.EOF {
		fmt.Println("empty")
	} else if err != nil {
		log.Fatal(err)
	}
	switch main_cmd {
	case "show":
		l := jsoncnt.ShowJSONcnt()
		for _, p := range l {
			fmt.Printf("Project:%s\t Todos:%s\n", p.Project, p.Todos)
		}
		JSONcmdStream(filename)
	case "add":
		err := jsoncnt.WriteJSONcnt(filename, mp["--project"], mt["--todo"])
		if err != nil {
			log.Fatal(err)
		}
		JSONcmdStream(filename)
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
