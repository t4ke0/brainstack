package jsoncnt

import (
	"encoding/json"
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

//DoneTasks represent the tasks that has been marked as Done
type DoneTasks struct {
	ProjectName string   `json:"ProjectName"`
	Task        []string `json:"Task"`
}

// DoneTasksList type that we use to store all Done tasks as a JSON array
type DoneTasksList []DoneTasks

//JSONlist list of JSONcontent where we store all project & todos
//to have an easy way to show and store other JSONcontent instances
type JSONlist []JSONcontent

var (
	list  JSONlist
	dlist DoneTasksList
)

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

// OpenJSONfile accept filename as input it opens the file
func OpenJSONfile(filename string, readDoneT, readCnt bool) error {
	f, err := checkForFile(filename)
	if err != nil {
		return err
	}
	list = list[:0]
	dlist = dlist[:0]
	dec := json.NewDecoder(f)
	_, err = dec.Token()
	if err != nil {
		return err
	}
	var j JSONcontent
	var dt DoneTasks
	for dec.More() {
		if readCnt {
			err := dec.Decode(&j)
			if err != nil {
				return err
			}
			list = append(list, j)
		} else if readDoneT {
			err := dec.Decode(&dt)
			if err != nil {
				return err
			}
			dlist = append(dlist, dt)
		}
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

// ShowDoneTask show done tasks list
func ShowDoneTask() DoneTasksList {
	if len(dlist) != 0 {
		return dlist
	}
	return nil
}

// SearchList search if an element exist on JSONlist
func SearchList(elemnt string, cArr JSONlist) bool {
	if len(cArr) == 0 {
		return false
	}
	for _, e := range cArr {
		if e.Project == elemnt {
			return true
		}
	}
	return false
}

//SaveDoneTasks get DoneTasksList type as input then save
//	it on done tasks file.
func SaveDoneTasks(filename string, doneT DoneTasksList) (bool, error) {
	f, err := checkForFile(filename)
	if err != nil {
		return false, err
	}
	err = json.NewEncoder(f).Encode(doneT)
	if err != nil {
		return false, err
	}
	return true, nil
}

//SavePT saves JSONlist elements to a json file
func SavePT(filename string, data JSONlist) (bool, error) {
	f, err := checkForFile(filename)
	if err == nil {
		err = json.NewEncoder(f).Encode(data)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, err
}

// WriteJSONcnt Add To JSON file content
// add new Project and it todos into json file
func WriteJSONcnt(filename, project, todos string) (bool, error) {
	if project == "" && todos == "" {
		return false, nil
	}
	f, err := checkForFile(filename)
	if err != nil {
		return false, err
	}
	j := &JSONcontent{}
	if exist := SearchList(project, list); !exist {
		j.Project, j.Todos = project, todos
		list = append(list, *j)
		err := json.NewEncoder(f).Encode(list)
		if err != nil {
			return false, err
		}
		//"Saved into the File"
		return true, nil
	}
	//"Project Already exist"
	return false, nil
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
func LIFO(projectName string) (ok bool) {
	if len(list) == 0 || projectName == "" {
		ok = false
	}
	for n, i := range list {
		if i.Project == projectName {
			stodo := strings.Split(i.Todos, " ")
			if len(stodo) > 1 {
				ntodo := strings.Split(i.Todos, ",")
				ntodo = ntodo[:len(ntodo)-1]
				i.Todos = strings.Join(ntodo, "")
				list[n] = i
				ok = true
			}
		}
	}
	return
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
