package cmdhandler

//TODO: JSONcmdParser Should be re-writed with a good way "gocyclo"
//TODO: JSONcmdParser should return somthing then print it in the main package
//TODO: introduce ncurses for your program
import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"../cmdtools"
	"../jsoncnt"
	"github.com/olekukonko/tablewriter"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initData(filename string) error {
	err := jsoncnt.OpenJSONfile(filename)
	if err == io.EOF {
		err = fmt.Errorf("%s", "empty")
	} else if err != nil {
		return err
	}
	return err
}

func dataPresentation() {
	l := jsoncnt.ShowJSONcnt()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Projects & Todos"})
	var data [][]string
	if len(l) != 0 {
		for _, p := range l {
			todo := strings.Split("=>"+p.Todos, "\n")
			projc := strings.Split(strings.ToUpper(p.Project+":"), "\n")
			data = append(data, projc, todo)
		}
		for _, v := range data {
			table.Append(v)
		}
		table.Render()
	} else {
		fmt.Println("You Have no Projects Or you Forgot to execute init cmd")
	}
}

func addNew(filename, project, todo string) (string, error) {
	isSaved, err := jsoncnt.WriteJSONcnt(filename, project, todo)
	if err != nil {
		return "", err
	}
	if isSaved {
		return "saved to File", nil
	}
	return "Wrong Args", nil
}

func todoAdder(project, todo string) string {
	todoAdded := jsoncnt.AddTodo(project, todo)
	if todoAdded {
		return fmt.Sprintf("added %s successfully\n", todo)
	}
	return "need project and todo values"
}

func saver(filename string) (string, error) {
	saved, err := jsoncnt.SaveCnt(filename)
	if err != nil {
		return "", err
	}
	if saved {
		return "Saved ...", nil
	}
	return "no changes has been maded", nil
}

//JSONcmdParser Parse the Commands executed by the user from the dashboard
func JSONcmdParser(cmd string, filename string) {
	cmdtools.InitArg("project")
	cmdtools.InitArg("todo")
	// Parse args
	mainCmd, a := cmdtools.ParseArg(cmd, "todo", "project")
	//Get args Values
	m := cmdtools.GetValue("project", a)
	m1 := cmdtools.GetValue("todo", a)

	switch mainCmd {
	case "init":
		initData(filename)
		JSONcmdStream(filename)
	case "show":
		dataPresentation()
		JSONcmdStream(filename)
	case "add":
		msg, err := addNew(filename, m["project"], m1["todo"])
		checkError(err)
		fmt.Println(msg)
		JSONcmdStream(filename)
	case "addTodo":
		msg := todoAdder(m["project"], m1["todo"])
		fmt.Println(msg)
		JSONcmdStream(filename)
	case "save":
		msg, err := saver(filename)
		checkError(err)
		fmt.Println(msg)
		JSONcmdStream(filename)
	case "clear":
		cmdtools.ClearScreen()
		JSONcmdStream(filename)
	case "LIFO":
		if ok := jsoncnt.LIFO(m["project"]); ok {
			fmt.Println("removed successFully")
		}
		JSONcmdStream(filename)
	case "FIFO":
		if ok := jsoncnt.FIFO(m["project"]); ok {
			fmt.Println("removed successFully")
		}
		JSONcmdStream(filename)
	case "removetodo":
		if ok := jsoncnt.RemoveTodo(m["project"], m1["todo"]); ok {
			fmt.Printf("removed %s SuccessFully\n", m1["todo"])
		}
		JSONcmdStream(filename)
	case "help":
		hl := cmdtools.HelpMenu()
		for _, h := range hl {
			fmt.Println(h)
		}
		JSONcmdStream(filename)
	default:
		fmt.Println("No Such Command Found")
		JSONcmdStream(filename)
	}
}

//JSONcmdStream the DashBoard of the User
func JSONcmdStream(filename string) {
	fmt.Printf("BRAINSTACK# ")
	cmd := bufio.NewScanner(os.Stdin)
	cmd.Scan()
	command := cmd.Text()
	JSONcmdParser(command, filename)
}
