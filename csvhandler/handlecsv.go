package csvhandler

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//CheckError Checks For Errors then it logs them
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//CreateFile creates New csv file or open an existing one
func CreateFile(file string) *os.File {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	CheckError(err)
	return f
}

//RemoveFile removes an existing csv file
func RemoveFile(file string) bool {
	err := os.Remove(file)
	CheckError(err)
	return true
}

//ReadCsv reads content from a csv file
func ReadCsv(filename string) [][]string {
	f, err := os.Open(filename)
	CheckError(err)
	defer f.Close()
	csvreader := csv.NewReader(f)
	rec, errno := csvreader.ReadAll()
	CheckError(errno)
	return rec
}

//WriteCsv write content on the csv file
func WriteCsv(filename string, content []string) {
	f := CreateFile(filename)
	defer f.Close()
	csvwriter := csv.NewWriter(f)
	csvwriter.Write(content)
	csvwriter.Flush()
}

//LookForCsvFiles Search csv files in the current directory
func LookForCsvFiles(file string) []string {
	var existF []string
	dir, err := ioutil.ReadDir(file)
	CheckError(err)
	for _, i := range dir {
		if filepath.Ext(i.Name()) == ".csv" {
			existF = append(existF, i.Name())
		}
	}
	return existF
}

//RemoveLastElement removes the last element of the stack when done flag is true
func RemoveLastElement(ideas [][]string) ([][]string, string) {
	var message string
	index := 1
	tail := ideas[len(ideas)-index]
	//If We Reach the End of the Ideas we show a message for the User
	if ideas[0][0] == "." {
		message = "You Are Done reset the file now"
		return ideas, message
	}
	//if first element of the tail is a 'DOT'
	if tail[0] == "." {
		ideas = append(ideas[:len(ideas)-(index+1)], ideas[len(ideas)-(index+1)])
		tail = ideas[len(ideas)-(index)]
	}
	//if length of the tail is bigger than one
	if len(tail) > 1 {
		//tot == TailOfTail
		tot := tail[len(tail)-index]
		if tot != "." {
			rep := strings.ReplaceAll(strings.Join(tail, ","), tot, ".")
			//Replacing the last element of the tail with a DOT
			ideas = append(ideas[:len(ideas)-1], strings.Split(rep, ","))
		} else {
			for _, i := range tail {
				if i == "." {
					index++
				}
			}
			tot := tail[len(tail)-index]
			rep := strings.ReplaceAll(strings.Join(tail, ","), tot, ".")
			ideas = append(ideas[:len(ideas)-1], strings.Split(rep, ","))
		}
		// else if the length of the tail is less than or equal ONE
	} else if len(tail) <= 1 {
		t := tail[index-1]
		//Replacing the old tail with a DOT
		nt := strings.ReplaceAll(strings.Join(tail, ","), t, ".")
		ideas = append(ideas[:len(ideas)-index], strings.Split(nt, ""))
	}
	return ideas, message
}
