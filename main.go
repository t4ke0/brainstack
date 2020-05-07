package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"./cmdhandler"
	"./cmdtools"
	"./csvhandler"
)

var c csvhandler.CSV

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage ./brainstack <csv/json-file-name>")
		os.Exit(1)
	}
	fileName := os.Args[1]
	if filepath.Ext(fileName) == ".csv" {
		f := csvhandler.OpenCsvFile(fileName)
		content := c.ReadCsv(f)
		CmdStream(content.Content, fileName)
	} else if filepath.Ext(fileName) == ".json" {
		cmdhandler.JSONcmdStream(fileName)
	} else {
		os.Exit(1)
	}
}

func addContentCsv(elements []string, filename string) [][]string {
	f := csvhandler.OpenCsvFile(filename)
	ncontent := c.AddContent(elements, f)
	return ncontent
}

//TODO : Remove CSV handling

// TODO: Rename This Function To csvStream
func CmdStream(content [][]string, filename string) {
	fmt.Printf("csv# ")
	cmd := bufio.NewScanner(os.Stdin)
	cmd.Scan()
	switch cmmd := strings.Split(cmd.Text(), " ")[0]; cmmd {
	case "show":
		// Show content of the csv file directly by opening the csv file using tablewriter
		c.PresentContent(filename)
		CmdStream(content, filename)
	case "done":
		if len(c.Content) == 0 {
			fmt.Println("YOU Don't have nothing on your todos")
			CmdStream(c.Content, filename)
		}
		tail := c.ReturnTail()
		ncontent := c.ReturnContent(tail)
		if len(ncontent) == 0 {
			fmt.Println("Empty")
			ncontent = [][]string{{""}}
		}
		c.Content = ncontent
		CmdStream(ncontent, filename)
	case "add":
		contentToAdd := strings.Split(cmd.Text(), " ")[1:]
		cnt := addContentCsv(contentToAdd, filename)
		content = cnt
		//content = append(content, contentToAdd)
		CmdStream(content, filename)
	case "clear":
		cmdtools.ClearScreen()
		CmdStream(content, filename)
	default:
		//fmt.Println("Wrong Command")
		CmdStream(content, filename)
	}
}
