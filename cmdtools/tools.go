package cmdtools

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

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

var arguments = []string{"--project", "--todo"}

func ArgParser(cmd string) (string, map[string]string, map[string]string) {
	m := make(map[string]string)
	tm := make(map[string]string)
	f := strings.Fields(cmd)
	main_cmd := f[0]
	rest := f[1:]
	for n, c := range rest {
		for _, a := range arguments {
			if c == a && c != "--todo" {
				m[c] = rest[n+1]
			} else if c == a && c == "--todo" {
				var l []string
				for _, w := range rest[n+1:] {
					if w != "--project" {
						l = append(l, w)
					} else {
						break
					}
				}
				tm[c] = strings.Join(l, " ")
			}
		}
	}
	return main_cmd, m, tm
}
