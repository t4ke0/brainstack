package cmdtools

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Arg struct {
	Name  string
	Value string
}

type Arguments []Arg

var Arglist Arguments

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

// ADD other Arguments IF we need others
//var arguments = []string{"--project", "--todo"}
//
//func ArgParser(cmd string) (string, map[string]string, map[string]string) {
//	m := make(map[string]string)
//	tm := make(map[string]string)
//	f := strings.Fields(cmd)
//	main_cmd := f[0]
//	rest := f[1:]
//	for n, c := range rest {
//		for _, a := range arguments {
//			if c == a && c != "--todo" {
//				m[c] = rest[n+1]
//			} else if c == a && c == "--todo" {
//				var l []string
//				for _, w := range rest[n+1:] {
//					if w != "--project" {
//						l = append(l, w)
//					} else {
//						break
//					}
//				}
//				tm[c] = strings.Join(l, " ")
//			}
//		}
//	}
//	return main_cmd, m, tm
//}

func InitArg(argname string) Arguments {
	a := &Arg{Name: argname, Value: ""}
	Arglist = append(Arglist, *a)
	return Arglist
}

func ParseArg(cmd, more, exclude string) (string, Arguments) {
	scmd := strings.Split(cmd, " ")
	main_cmd := scmd[0]
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
	return main_cmd, Arglist
}

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
