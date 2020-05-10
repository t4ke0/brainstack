package csvhandler

import (
	"encoding/csv"
	//"fmt"
	"io"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
)

// OpenCsvFile open the csv file
func OpenCsvFile(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// CSV struct contains a content instance that we will use \
// to interact with csv content
type CSV struct {
	Content [][]string
}

//// SwapId function swaps the Id to the begining of the array
//func SwapId(record []string) []string {
//	for i := range record {
//		record[len(record)-1], record[i] = record[i], record[len(record)-1]
//	}
//	return record
//}

// ReadCsv reads the csv file
func (c *CSV) ReadCsv(rc *os.File) *CSV {
	r := csv.NewReader(rc)
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if err, ok := err.(*csv.ParseError); ok && err.Err == csv.ErrFieldCount {
				log.Fatal(err)
			}
		}
		c.Content = append(c.Content, rec)
	}
	//fmt.Println(c.Content)
	return c
}

// AddContent function handle adding new Stuff to the Csv File
func (c *CSV) AddContent(content []string, rc *os.File) [][]string {
	// add Content to our struct then add it to the file
	c.Content = append(c.Content, content)
	w := csv.NewWriter(rc)
	if err := w.Write(content); err != nil {
		log.Fatal(err)
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	return c.Content
}

// PresentContent using terminal table to show content of the csv file
func (c *CSV) PresentContent(filename string) {
	content := c.Content
	table := tablewriter.NewWriter(os.Stdout)
	for _, cnt := range content {
		table.Append(cnt)
	}
	table.Render()
}

// ReturnTail check for the size of the tail of the content
// If tail > 1 => we remove the last elem of that tail then \
// return the new tail to append it to the content in an other func
// If tail < 1 => we return an empty slice to inform the other func \
// to update the content
func (c *CSV) ReturnTail() []string {
	content := c.Content
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
func (c *CSV) ReturnContent(tail []string) [][]string {
	content := c.Content
	if len(tail) >= 1 {
		content = append(content[:len(content)-1], tail)
		return content
	} else if len(tail) < 1 {
		content = content[:len(content)-1]
		return content
	}
	return content

}
