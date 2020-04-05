package csvhandler

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
)

// CSV struct contains a content instance that we will use \
// to interact with csv content
type CSV struct {
	Content [][]string
}

var c CSV

// OpenCsvFile open the csv file
func OpenCsvFile(filename string) *os.File {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// ReadCsv reads the csv file
func ReadCsv(rc *os.File) [][]string {
	r := csv.NewReader(rc)
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		c.Content = append(c.Content, rec)
	}
	return c.Content
}

// PresentContent using terminal table to show content of the csv file
func PresentContent(content [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	for _, v := range content {
		table.Append(v)
	}
	table.Render()
}

// ReturnTail check for the size of the tail of the content
// If tail > 1 => we remove the last elem of that tail then \
// return the new tail to append it to the content in an other func
// If tail < 1 => we return an empty slice to inform the other func \
// to update the content
func ReturnTail(content [][]string) []string {

	if tail := content[len(content)-1]; len(tail) > 1 {
		ntail := tail[:len(tail)-1]
		return ntail
	} else if len(tail) <= 1 {
		ntail := tail[:len(tail)-1]
		return ntail
	} else if len(tail) == 0 {
		return []string{}
	}
	return []string{}

}

// ReturnContent function takes the tail who is a []string \
// Returns new content after processing it
// if tail is not empty we append that tail to our content \
// and return a new content [][]string
// Otherwise "if tail is empty" we cut it from the content \
// & return the new content
func ReturnContent(tail []string, content [][]string) [][]string {

	if len(tail) >= 1 {
		content = append(content[:len(content)-1], tail)
		return content
	} else if len(tail) < 1 {
		content = content[:len(content)-1]
		return content
	}
	return content

}
