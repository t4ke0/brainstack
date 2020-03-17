package flaghandler

import (
	"../csvhandler"
	"fmt"
	"strings"
)

func HandleAddingIdeas(ideas string, file string) string {
	var msg string
	if ideas != "" {
		content := strings.Split(ideas, ",")
		csvhandler.WriteCsv(file, content)
		msg = "Success"
	} else {
		msg = "Failed"
	}
	return msg
}

func HandleDoneFlag(file string) ([]string, string) {
	var m string
	var f []string
	if file == "" {
		files := csvhandler.LookForCsvFiles("./")
		f = files
	} else {
		ideas := csvhandler.ReadCsv(file)
		fl, msg := csvhandler.RemoveLastElement(ideas)
		if msg != "" {
			m = msg
		}
		if rm := csvhandler.RemoveFile(file); rm {
			for _, line := range fl {
				csvhandler.WriteCsv(file, line)
				fmt.Println(line)
			}
		}
	}
	return f, m
}

func HandleResetFlag(file string) string {
	var msg string
	if rm := csvhandler.RemoveFile(file); rm {
		msg = ("Reset file successfully")
	}
	return msg
}
