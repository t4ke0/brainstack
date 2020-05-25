package jsoncnt

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//JSONcontent struct it has the instance that we are using to store
//project and todos as a json format
type JSONcontent struct {
	Project string `json:"Project"`
	Todos   string `json:"Todos"`
}

//JSONlist list of JSONcontent where we store all project & todos
//to have an easy way to show and store other JSONcontent instances
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

// ShowJSONcnt show content of JSONlist
func ShowJSONcnt() JSONlist {
	if len(list) != 0 {
		return list
	}
	return nil
}

func searchList(elemnt string) bool {
	if len(list) == 0 {
		return false
	}
	for _, e := range list {
		if e.Project == elemnt {
			return true
		}
	}
	return false
}

// WriteJSONcnt Add To JSON file content
func WriteJSONcnt(filename, project, todos string) (error, bool) {
	if project == "" && todos == "" {
		return nil, false
	}
	f, err := checkForFile(filename)
	if err != nil {
		return err, false
	}
	j := &JSONcontent{}
	if exist := searchList(project); !exist {
		j.Project, j.Todos = project, todos
		list = append(list, *j)
		err := json.NewEncoder(f).Encode(list)
		if err != nil {
			return err, false
		}
		//fmt.Println("Saved into the File")
		return nil, true
	}
	//fmt.Println("Project Already exist")
	return nil, false
}

//SaveCnt save the list in it's current situation to the json file
func SaveCnt(filename string) (bool, error) {
	// try to remove the file first
	if err := os.Remove(filename); err != nil {
		return false, err
	}

	f, err := checkForFile(filename)
	if err != nil {
		return false, err
	}
	if len(list) != 0 {
		err := json.NewEncoder(f).Encode(list)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

//LIFO last in first out
func LIFO(projectName string) bool {
	nlist := JSONlist{}
	if len(list) == 0 || projectName == "" {
		return false
	}
	for _, i := range list {
		if i.Project == projectName {
			stodo := strings.Split(i.Todos, " ")
			if len(stodo) > 1 {
				ntodo := strings.Split(i.Todos, ",")
				ntodo = ntodo[:len(ntodo)-1]
				i.Todos = strings.Join(ntodo, "")
				nlist = append(nlist, i)
			}
		} else {
			nlist = append(nlist, i)
			continue
		}
	}
	list = nlist
	return true
}

//AddTodo add Todo to a particular project
func AddTodo(projectName, newTodo string) bool {
	var nlist JSONlist
	if projectName == "" && newTodo == "" {
		return false
	}
	if len(list) == 0 {
		return false
	}
	for _, n := range list {
		if n.Project == projectName {
			todoS := strings.Split(n.Todos, ",")
			todoS = append(todoS, newTodo)
			n.Todos = strings.Join(todoS, ",")
			nlist = append(nlist, n)
		} else {
			nlist = append(nlist, n)
		}
	}
	list = nlist
	return true
}

//FIFO remove the first element of the todo string
func FIFO(projectName string) bool {
	var nlist JSONlist
	if len(list) == 0 || projectName == "" {
		return false
	}
	for _, n := range list {
		if n.Project == projectName {
			stodo := strings.Split(n.Todos, ",")
			if len(stodo) > 1 {
				ntodo := stodo[len(stodo)-1:]
				n.Todos = strings.Join(ntodo, ",")
				nlist = append(nlist, n)
			}
		} else {
			nlist = append(nlist, n)
		}
	}
	list = nlist
	return true
}

//RemoveTodo remove todo for a particular project
func RemoveTodo(projectName, todo string) bool {
	var nlist JSONlist
	if len(list) == 0 || projectName == "" {
		return false
	}
	for _, n := range list {
		if n.Project == projectName {
			stodo := strings.Replace(n.Todos, todo, "", 1)
			ntodo := strings.Split(stodo, ",")
			n.Todos = strings.Join(ntodo, "")
			nlist = append(nlist, n)
		} else {
			nlist = append(nlist, n)
		}
	}
	list = nlist
	return true
}
