package cmdtools

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

//Arg struct who hold the pairs of Name & value of the arguments
type Arg struct {
	Name  string
	Value string
}

//Arguments Array of Arg
type Arguments []Arg

//Arglist variable which it's type is the array of Arg
var Arglist Arguments

//ClearScreen clear terminal window for several OS systems
func ClearScreen() {
	switch OS := runtime.GOOS; OS {
	case "linux":
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
	case "windows":
		c := exec.Command("cmd", "/c", "cls")
		c.Stdout = os.Stdout
		c.Run()
	default:
		c := exec.Command("clear")
		c.Stdout = os.Stdout
		c.Run()
	}
}

//InitArg initialize argument put them in Arg struct it has Name & Value pair instances
func InitArg(argname string) Arguments {
	a := &Arg{Name: argname, Value: ""}
	Arglist = append(Arglist, *a)
	return Arglist
}

//ParseArg get the cmd , more and exclude as input
// for more is the argument that we need to get a big amount of string sequence as value
// and exclude is the argument that we are gonna skip when are we gonna get the value
// of `more` argument.
func ParseArg(cmd, more, exclude string) (string, Arguments) {
	scmd := strings.Split(cmd, " ")
	mainCmd := scmd[0]
	rest := scmd[1:]
	var narglist Arguments
	for j, i := range rest {
		for _, a := range Arglist {
			if strings.Trim(i, "--") == a.Name && more != a.Name {
				a.Value = rest[j+1]
				narglist = append(narglist, a)
				break
			} else if strings.Trim(i, "--") == a.Name && more == a.Name {
				var tempL []string
				for _, e := range rest[j+1:] {
					if strings.Trim(e, "--") != exclude {
						tempL = append(tempL, e)
					} else {
						break
					}
				}
				a.Value = strings.Join(tempL, " ")
				narglist = append(narglist, a)
				break
			} else {
				continue
			}
		}
	}
	Arglist = narglist
	return mainCmd, Arglist
}

//GetValue Get the argname as input and the argument list then returns a map[string]string
//which is the argument and it value as a map.
func GetValue(argname string, l Arguments) map[string]string {
	m := make(map[string]string)
	if len(l) != 0 {
		for _, a := range l {
			if a.Name == argname {
				m[argname] = a.Value
			}
		}
	}
	return m
}

//HelpMenu show help to the User
func HelpMenu() []string {
	var help []string
	help = []string{"init initialize your content from json file to the program", "show shows you your projects and todos in a table", "add  add projects and todos to your json file and save them there eg. add --project <project name> --todo <todo here> ", "save save changes you have made", "clear clears the screen", "done if you done with a todo"}
	return help
}
