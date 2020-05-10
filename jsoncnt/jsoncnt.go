package jsoncnt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type JSONcontent struct {
	Project string `json:"Project"`
	Todos   string `json:"Todos"`
}

type JSONlist []JSONcontent

var list JSONlist

func checkForFile(filename string) (*os.File, error) {
	var f *os.File
	dir, err := ioutil.ReadDir("./")
	if err != nil {
		return nil, err
	}
	for _, file := range dir {
		if file.Name() == filename {
			f, err = os.Open(filename)
			if err != nil {
				return nil, err
			}
		} else {
			f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
			if err != nil {
				return nil, err
			}
		}
	}
	return f, nil
}

// OpenJSONfile accept filename as input it opens the file & return *os.File
func OpenJSONfile(filename string) error {
	f, err := checkForFile(filename)
	if err != nil {
		return err
	}
	list = list[:0]
	dec := json.NewDecoder(f)
	_, err = dec.Token()
	if err != nil {
		return err
	}
	var j JSONcontent
	for dec.More() {
		err := dec.Decode(&j)
		if err != nil {
			return err
		}
		list = append(list, j)
	}
	_, err = dec.Token()
	if err != nil {
		return err
	}
	return nil
}

// ShowJSONcnt
func ShowJSONcnt() JSONlist {
	if len(list) != 0 {
		return list
	} else {
		return nil
	}
}

// WriteJSONcnt
func WriteJSONcnt(filename, project, todos string) error {
	f, err := checkForFile(filename)
	if err != nil {
		return err
	}
	var j = &JSONcontent{}
	for _, p := range list {
		if p.Project != project {
			j.Project, j.Todos = project, todos
			list = append(list, *j)
			err := json.NewEncoder(f).Encode(list)
			if err != nil {
				return err
			}
			break
		} else {
			fmt.Println("Project Already Exist")
			break
		}
	}
	return nil
}

//LIFO last in first out
func LIFO(projectName string) JSONlist {
	nlist := JSONlist{}
	for _, i := range list {
		if i.Project == projectName {
			stodo := strings.Split(i.Todos, "")
			if len(stodo) > 1 {
				ntodo := strings.Split(i.Todos, ",")
				ntodo = ntodo[:len(ntodo)-1]
				i.Todos = strings.Join(ntodo, "")
				nlist = append(nlist, i)
			}
			break
		} else {
			nlist = append(nlist, i)
			continue
		}
	}
	list = nlist
	return list
}
