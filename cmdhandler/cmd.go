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
		err := jsoncnt.OpenJSONfile(filename)
		if err == io.EOF {
			fmt.Println("empty")
		} else if err != nil {
			log.Fatal(err)
		}
		JSONcmdStream(filename)
	case "show":
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
		JSONcmdStream(filename)
	case "add":
		err, isSaved := jsoncnt.WriteJSONcnt(filename, m["project"], m1["todo"])
		checkError(err)
		if isSaved {
			fmt.Println("saved to File")
			JSONcmdStream(filename)
		} else {
			fmt.Println("Wrong Args")
			JSONcmdStream(filename)
		}
	case "addTodo":
		todoAdded := jsoncnt.AddTodo(m["project"], m1["todo"])
		if todoAdded {
			fmt.Printf("added %s successfully\n", m1["todo"])
		} else {
			fmt.Println("need project and todo values")
		}
		JSONcmdStream(filename)
	case "save":
		saved, err := jsoncnt.SaveCnt(filename)
		checkError(err)
		if saved {
			fmt.Println("Saved ...")
		} else {
			fmt.Println("no changes has been maded")
		}
		JSONcmdStream(filename)
	case "clear":
		cmdtools.ClearScreen()
		JSONcmdStream(filename)
	case "FIFO":
		if ok := jsoncnt.LIFO(m["project"]); ok {
			fmt.Println("removed successFully")
		}
		JSONcmdStream(filename)
	case "LIFO":
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
