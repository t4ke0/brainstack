package auxiliary

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/TaKeO90/brainstack/jsoncnt"
	"github.com/jroimartin/gocui"
)

const doneTaskFile string = "done.json"

var (
	//lock this variable give the user the permission to execute command after initializing the json file
	lock = true

	//Editable global variable indicate wheter a certain view is in edit mode or not
	Editable = false

	//LHighP last selected project.
	LHighP string

	//CntList where we store projects and their todos.
	CntList jsoncnt.JSONlist

	//DoneList Where we store Done Task for a specific project.
	DoneList jsoncnt.DoneTasksList
)

// TaskDataSaver an interface that has a method that saves a new task
type TaskDataSaver interface {
	saveNewTasks() (bool, bool)
}

// ProjectDataSaver an interface that holds a method which saves new projects
type ProjectDataSaver interface {
	saveNewProjects() bool
}

// Tasks struct we store on it projects and their todos
type Tasks struct {
	//Selected project that the user save it todos
	Project string
	//New todos modified by the user
	Todos []string
}

// Projects has ProjectList where we store project view elements
type Projects struct {
	ProjectList []string
}

func removeView(views *[]*gocui.View, i int) {
	*views = append((*views)[:i], (*views)[i+1:]...)
}

func saveDoneTasks(pName, dTask string) (bool, error) {
	dT := &jsoncnt.DoneTasks{}
	if len(DoneList) != 0 {
		for i, n := range DoneList {
			// if project name exist in DoneList before
			if n.ProjectName == pName && !checkTasks(n.Task, dTask) {
				n.Task = append(n.Task, dTask)
				DoneList[i] = n
				removeTask(pName, dTask)
				// if project name doesn't exist before in Donelist
			} else if n.ProjectName != pName && !checkTasks(n.Task, dTask) {
				dT.ProjectName = pName
				dT.Task = append(dT.Task, dTask)
				DoneList = append(DoneList, *dT)
				removeTask(pName, dTask)
			}
		}
	} else {
		dT.ProjectName = pName
		dT.Task = append(dT.Task, dTask)
		DoneList = append(DoneList, *dT)
		removeTask(pName, dTask)
	}
	isOK, err := jsoncnt.SaveDoneTasks(doneTaskFile, DoneList)
	if err != nil {
		return false, err
	}
	return isOK, nil
}

//checkTasks checks if there is a duplicate task
func checkTasks(tasks []string, tName string) bool {
	for _, i := range tasks {
		if i == tName {
			return true
		}
		return false
	}
	return false
}

//removeTask removes the task that has been added as done one
func removeTask(pname, tname string) bool {
	for i, n := range CntList {
		if n.Project == pname {
			rmT := strings.Replace(n.Todos, tname, "", 1)
			nt := strings.Split(rmT, ",")
			n.Todos = strings.Join(nt, "")
			CntList[i] = n
		}
	}
	return false
}

func getDoneTask() error {
	err := jsoncnt.OpenJSONfile(doneTaskFile, true, false)
	if err != nil || err == io.EOF {
		return err
	}
	DoneList = jsoncnt.ShowDoneTask()
	return nil
}

func showDoneTasks(v *gocui.View, pname string) {
	//currentBuffer := v.BufferLines()
	for _, n := range DoneList {
		if n.ProjectName == pname {
			for _, t := range n.Task {
				fmt.Fprintf(v, "%s\n", t)
			}
		}
	}
}

// AUXILIARY FUNCTIONS ///
func formatInfo(v *gocui.View, info string, formatInf interface{}) {
	v.Clear()
	fmt.Fprintf(v, info, formatInf)
	err := v.SetCursor(0, 0)
	checkError(err)
}

func showInfo(v *gocui.View, info string) {
	v.Clear()
	fmt.Fprintf(v, info)
	err := v.SetCursor(0, 0)
	checkError(err)
}

func checkError(err error) {
	if err != nil && err != gocui.ErrUnknownView {
		log.Fatal(err)
	}
}

//function to get views Data just provide view Name and gui then specify if  you need line value or the whole bufferLine
func getVData(viewName string, g *gocui.Gui, line bool) ([]string, string) {
	// in here we gonna grab all view Data as []string
	if !line {
		v, err := g.View(viewName)
		checkError(err)
		return v.BufferLines(), ""
	}
	v, err := g.View(viewName)
	checkError(err)
	_, cy := v.Cursor()
	val, err := v.Line(cy)
	checkError(err)
	return []string{}, val
}

// getNewData get Data from views in function of Vctg
// Vtg is view category `project views or todo or done tasks`
func getNewData(g *gocui.Gui, Vctg string) (*Tasks, *Projects) {
	var data *Tasks
	var pdata *Projects
	if len(CntList) == 0 && lock {
		return nil, nil
	}
	if Vctg == "task" {
		todos, _ := getVData("todo", g, false)
		_, currentP := getVData("list", g, true)
		if currentP != "" && currentP != LHighP || len(todos) != 0 {
			data = &Tasks{LHighP, todos}
		} else {
			data = &Tasks{currentP, todos}
		}
		return data, nil
	} else if Vctg == "project" {
		projects, _ := getVData("list", g, false)
		pdata = &Projects{projects}
		return nil, pdata
	} else {
		return nil, nil
	}
}

func (vd *Tasks) saveNewTasks() (bool, bool) {
	if vd != nil && CntList != nil {
		for i, n := range CntList {
			if vd.Project == n.Project {
				if len(vd.Todos) != 0 {
					n.Todos = strings.Join(vd.Todos, ",")
					CntList[i] = n
				} else {
					return false, true
				}
			}
		}
		return true, true //NOTE: modified was (true, false)
	}
	return false, false
}

func (pjc *Projects) saveNewProjects() bool {
	//Check if CntList Contain any of those projects
	p := new(jsoncnt.JSONcontent)
	var ok bool
	if len(CntList) != 0 {
		for _, n := range pjc.ProjectList {
			isExist := jsoncnt.SearchList(n, CntList)
			if !isExist {
				p.Project = n
				CntList = append(CntList, *p)
				ok = true
			} else {
				ok = false
			}
		}
	} else {
		for _, n := range pjc.ProjectList {
			p.Project = n
			CntList = append(CntList, *p)
			ok = true
		}
	}
	return ok
}

func parseCmd(rawCmd string) map[string]string {
	mainCmd, subCmd := strings.Fields(rawCmd)[0], strings.Join(strings.Fields(rawCmd)[1:], " ")
	cmdMap := make(map[string]string)
	cmdMap[mainCmd] = subCmd
	return cmdMap
}

func initHandler(filename string, v *gocui.View) {
	// When initializing then we need to set the lock to false
	lock = false
	//Read file here
	err := jsoncnt.OpenJSONfile(filename, false, true)
	if err != nil || err == io.EOF {
		showInfo(v, "file is Empty you have no projects")
	} else {
		showInfo(v, `file initialized now you can interact with your data`)
		CntList = jsoncnt.ShowJSONcnt()
		err := getDoneTask()
		if err != nil {
			showInfo(v, "You have no Done tasks")
		}
	}
}

func saveHandler(filename string, v *gocui.View, g *gocui.Gui, vl string) {
	if !lock {
		var svD TaskDataSaver
		var pjD ProjectDataSaver
		if vl == "todos" {
			tsk, _ := getNewData(g, "task")
			svD = tsk
			isInit, isTodo := svD.saveNewTasks()
			if !isInit && !isTodo {
				showInfo(v, `You should run 'init' command first or you didn't modify todo section`)
			} else if !isInit && isTodo {
				showInfo(v, `Really ? wanna save even tough you didn't modify todos`)
			} else {
				ok, err := jsoncnt.SavePT(filename, CntList)
				if ok {
					showInfo(v, "saved")
				}
				if err != nil {
					showInfo(v, "failed to save")
				}
			}
		} else if vl == "projects" {
			_, pjct := getNewData(g, "project")
			pjD = pjct
			if ok := pjD.saveNewProjects(); ok {
				isSaved, err := jsoncnt.SavePT(filename, CntList)
				checkError(err)
				if isSaved {
					showInfo(v, "saved")
				} else {
					showInfo(v, "failed to save")
				}
			} else {
				showInfo(v, "failed maybe that project already exist")
			}
		} else if vl == "donetask" {
			isSaved, err := jsoncnt.SaveDoneTasks(doneTaskFile, DoneList)
			checkError(err)
			if isSaved {
				showInfo(v, "done tasks has been saved")
			}
		}
	} else {
		showInfo(v, "Should initialze the file first with `init` command")
	}
}

//TODO: handle Done Command.
func commandHandler(g *gocui.Gui, cmd string, v *gocui.View, filename string) {
	m := parseCmd(cmd)
	for k, vl := range m {
		switch k {
		case "init":
			initHandler(filename, v)
		case "save":
			saveHandler(filename, v, g, vl)
		case "project":
			if !lock {
				if LHighP != "" {
					formatInfo(v, "Current Selected Project : %s\n", LHighP)
				} else {
					showInfo(v, "You have not selected no project")
				}
			} else {
				showInfo(v, "Should initialze the file first with `init` command")
			}
		case "recover":
			if !lock {
				if recT := strings.TrimSpace(vl); recT != "" {
					recvT, isRec := recoverDone(g, recT)
					if isRec {
						formatInfo(v, "%s: task has been recovered", recvT)
						err := enterShowTodos(g, v)
						checkError(err)
					} else {
						showInfo(v, "Failed to recover the task")
					}
				}
			}
		case "set":
			if !lock {
				setValuesHandler(vl, v)
			}
		default:
			showInfo(v, "command not found")
		}
	}
}

// This Function helps to get the element and it value to set
func setValuesHandler(vl string, notifV *gocui.View) {
	if vl != "" {
		element, value := strings.Fields(vl)[0], strings.Join(strings.Fields(vl)[1:], " ")
		switch element {
		case "project":
			if value != "" {
				if ok := setCurrentProject(value); ok {
					formatInfo(notifV, "Current Project : %s", value)
				} else {
					showInfo(notifV, "Failed to set Project, or project you are trying to set doesn't exist")
				}
			}
		}
	}
}

func printProjects(data jsoncnt.JSONlist, v *gocui.View) {
	if len(data) != 0 {
		v.Clear()
		for _, n := range data {
			fmt.Fprintf(v, "%s\n", n.Project)
		}
	}
}

func orderViews(views *[]*gocui.View) {
	for i, v := range *views {
		switch v.Name() {
		case "cmd":
			if i != 0 {
				swap(i, 0, views)
			}
		case "list":
			if i != 1 {
				swap(i, 1, views)
			}
		case "todo":
			if i != 2 {
				swap(i, 2, views)
			}
		case "done":
			if i != 3 {
				swap(i, 3, views)
			}
		case "notif":
			if i != 4 {
				swap(i, 4, views)
			}
		}
	}
}

func swap(i, j int, views *[]*gocui.View) {
	(*views)[i], (*views)[j] = (*views)[j], (*views)[i]
}

// THIS IS THE FUNCTIONS WE GONNA USE WITH KEYBINDING

func editToNormal(g *gocui.Gui, v *gocui.View) error {
	nv, err := g.View("notif")
	if err != nil {
		return err
	}
	views := g.Views()
	if Editable {
		for _, v := range views {
			if v.Name() != "notif" {
				Editable = false
				v.Editable = Editable
				showInfo(nv, "Normal Mode")
			}
		}
	}
	return nil
}

func normalToEdit(g *gocui.Gui, v *gocui.View) error {
	nv, err := g.View("notif")
	if err != nil {
		return err
	}
	view := g.Views()
	if !Editable {
		Editable = true
		for _, v := range view {
			if v.Name() != "notif" {
				v.Editable = Editable
				showInfo(nv, "edit mode")
			}
		}
	}
	return nil
}

func enterShowTodos(g *gocui.Gui, v *gocui.View) error {
	lv, err := g.View("list")
	if err != nil {
		return err
	}
	tv, err := g.View("todo")
	if err != nil {
		return err
	}
	dv, err := g.View("done")
	if err != nil {
		return err
	}
	nv, err := g.View("notif")
	if err != nil {
		return err
	}
	_, cy := lv.Cursor()
	vBuff := lv.BufferLines()
	if len(vBuff) != 0 {
		for i, j := range CntList {
			for _, x := range vBuff {
				if i == cy && j.Project == x {
					LHighP = x
					if LHighP != "" {
						formatInfo(nv, "Current Project: %s", x)
					}
					dv.Clear()
					showDoneTasks(dv, x)
					tv.Clear()
					for _, T := range strings.Split(j.Todos, ",") {
						fmt.Fprintf(tv, "%s\n", strings.TrimSpace(T))
					}
				}
			}
		}
	}
	return nil
}

func arrowUp(g *gocui.Gui, v *gocui.View) error {
	vN := g.CurrentView()
	Lines := vN.BufferLines()
	if cx, cy := vN.Cursor(); cy != 0 {
		cy--
		err := vN.SetCursor(cx, cy)
		if err != nil {
			return err
		}
	} else if cy == 0 {
		err := vN.SetCursor(cx, len(Lines)-1)
		if err != nil {
			return err
		}
	}
	return nil
}

func arrowDown(g *gocui.Gui, v *gocui.View) error {
	currentV := g.CurrentView()
	cx, cy := currentV.Cursor()
	cy++
	Lines := currentV.BufferLines()
	err := currentV.SetCursor(cx, cy)
	if err != nil {
		return err
	}
	if cy == len(Lines) {
		err := currentV.SetCursor(cx, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func tabHandler(g *gocui.Gui, v *gocui.View) error {
	views := g.Views()
	orderViews(&views)
	removeView(&views, 4)
	currentV := g.CurrentView()
	for i, n := range views {
		if currentV == n {
			if i != len(views)-1 {
				i++
			} else {
				i = 0
			}
			_, err := g.SetCurrentView(views[i].Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func endKeyHandler(g *gocui.Gui, v *gocui.View) error {
	currentV := g.CurrentView()
	cBuff := currentV.Buffer()
	cx, cy := currentV.Cursor()
	if cBuff != "" {
		word, err := currentV.Line(cy)
		if err != nil {
			return err
		}
		if cx != len(word)-1 && word != "" {
			err := currentV.SetCursor(len(word)-1, cy)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func homeKeyHandler(g *gocui.Gui, v *gocui.View) error {
	currentV := g.CurrentView()
	cBuff := currentV.Buffer()
	cx, cy := currentV.Cursor()
	if cBuff != "" {
		word, err := currentV.Line(cy)
		if err != nil {
			return err
		}
		if word != "" && cx != 0 {
			if err := currentV.SetCursor(0, cy); err != nil {
				return err
			}
		}
	}
	return nil
}

func setCurrentProject(projectName string) bool {
	//TODO: check if that project exist or not cause we cant set a project that is not exist
	if LHighP != "" {
		for _, n := range CntList {
			if n.Project == projectName {
				LHighP = projectName
				return true
			}
		}
	}
	return false
}

func recoverDone(g *gocui.Gui, taskName string) (string, bool) {
	var recoverTask string
	var ok bool
	if taskName != "" {
		recoverTask = taskName
	} else {
		_, recP := getVData("done", g, true)
		recoverTask = recP
	}

	for i, n := range DoneList {
		if n.ProjectName == LHighP {
			nDone := strings.Replace(strings.Join(n.Task, ""), recoverTask, "", 1)
			n.Task = strings.Split(nDone, "")
			DoneList[i] = n
			ok = true
		}
	}
	for i, j := range CntList {
		if j.Project == LHighP {
			nTask := strings.Split(j.Todos, "")
			nTask = append(nTask, recoverTask)
			j.Todos = strings.Join(nTask, "")
			CntList[i] = j
			ok = true
		}
	}
	return recoverTask, ok
}

/////////////////////////

// KeyBindingHandler Handles Keybinding
func KeyBindingHandler(g *gocui.Gui, filename string) {
	lv, err := g.View("list")
	checkError(err)
	tv, err := g.View("todo")
	checkError(err)
	nv, err := g.View("notif")
	checkError(err)
	cv, err := g.View("cmd")
	checkError(err)
	dv, err := g.View("done")
	checkError(err)

	// HANDLING CTRL+C FOR QUITING THE PROGRAM
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})

	// HANDLING CTRL+D FOR SWITCHING FROM EDIT MODE TO NORMAL MODE
	g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, editToNormal)

	// HANDLING CTRL+O FOR SWITCHING FROM NORMAL MODE TO EDIT MODE
	g.SetKeybinding("", gocui.KeyCtrlO, gocui.ModNone, normalToEdit)

	// HANDLING ENTER KEY FOR PROJECT VIEW TO SHOW TASKS IN `TODOS` VIEW
	g.SetKeybinding(lv.Name(), gocui.KeyEnter, gocui.ModNone, enterShowTodos)

	// HANDLING UP ARROW KEY
	g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, arrowUp)

	// HANDLING DOWN ARROW FOR PROJECT VIEW
	g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, arrowDown)

	// HANDLING TAB KEY FOR SWITCHING BETWEEN VIEWS EXCEPT NOTIF VIEW
	g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, tabHandler)

	// HANDLING END KEY
	g.SetKeybinding("", gocui.KeyEnd, gocui.ModNone, endKeyHandler)

	// HANDLING HOME KEY
	g.SetKeybinding("", gocui.KeyHome, gocui.ModNone, homeKeyHandler)

	// HANDLING DELETE KEY FOR COMMAND VIEW
	g.SetKeybinding(cv.Name(), gocui.KeyDelete, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		cv.Clear()
		err := cv.SetCursor(0, 0)
		if err != nil {
			return err
		}
		return nil
	})

	//HANDLING CTRL+R IN DONE TASK VIEW TO RECOVER A DONE TASK
	g.SetKeybinding(dv.Name(), gocui.KeyCtrlR, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if recT, isRecovered := recoverDone(g, ""); isRecovered {
			formatInfo(nv, "%s:  task recovered\n", recT)
			err := enterShowTodos(g, v)
			if err != nil {
				return err
			}
		} else {
			showInfo(nv, "Failed To restore the task")
		}
		return nil
	})

	// HANDLING ENTER KEY FOR COMMAND VIEW FOR EXECUTING COMMANDS
	g.SetKeybinding(cv.Name(), gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if currentCmd := cv.Buffer(); currentCmd != "" {
			commandHandler(g, strings.TrimSpace(currentCmd), nv, filename)
			printProjects(CntList, lv)
			cv.Clear()
			err := cv.SetCursor(0, 0)
			if err != nil {
				return err
			}
		}
		return nil
	})
	// HANDLING ENTER KEY FOR TODOS VIEW WHICH MARK AS THE HIGHLIGHTED
	//	TASK AS DONE ONE AND SHOW THAT TASK ON DONE VIEW
	g.SetKeybinding(tv.Name(), gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		//GET THE CURRENT PROJECT
		_, cProj := getVData(lv.Name(), g, true)
		//GET CURRENT HIGHLIGHTED TASK
		_, cTask := getVData(tv.Name(), g, true)
		ok, err := saveDoneTasks(cProj, cTask)
		if err != nil {
			return err
		}
		if ok {
			//			dv.Clear()
			//			showDoneTasks(dv, cProj)
			if err := enterShowTodos(g, v); err != nil {
				return err
			}
		} else {
			showInfo(nv, "Failed to add that task to the Done tasks list")
		}
		return nil
	})
}
