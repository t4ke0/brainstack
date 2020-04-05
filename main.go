package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"./csvhandler"
)

func main() {
	start := time.Now()
	f := csvhandler.OpenCsvFile("file.csv")
	content := csvhandler.ReadCsv(f)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Read File in %d\n", elapsed)
	CmdStream(content)
}

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

func CmdStream(content [][]string) {
	fmt.Printf("# ")
	cmd := bufio.NewScanner(os.Stdin)
	cmd.Scan()
	switch cmd.Text() {
	case "show":
		csvhandler.PresentContent(content)
		CmdStream(content)
	case "done":
		tail := csvhandler.ReturnTail(content)
		ncontent := csvhandler.ReturnContent(tail, content)
		if len(ncontent) == 0 {
			fmt.Println("Empty")
			ncontent = [][]string{{""}}
		}
		CmdStream(ncontent)
	case "clear":
		ClearScreen()
		CmdStream(content)
	default:
		fmt.Println("Wrong Command")
		CmdStream(content)
	}
}
