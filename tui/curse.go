package main

import (
	"log"
	"os"
	"path/filepath"

	"./auxilary"
	"github.com/jroimartin/gocui"
)

//   TODO : we don't need command argument i Guess we just need to save the views at their current status when the user uses the command `save`
//#2 TODO:get keybinding from config file

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//InitGui initialize the GUI
func InitGui() (*gocui.Gui, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}
	g.Cursor = true
	return g, nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	_, err := g.SetView("list", 0, (maxY/2)-17, (maxY/2)+13, (maxY/2)+13)
	if err != nil {
		return err
	}
	_, err = g.SetView("todo", (maxX/2)-60, (maxY/2)-17, maxX-52, (maxY/2)+13)
	if err != nil {
		return err
	}
	_, err = g.SetView("cmd", (maxX/5)-1, (maxY/2)+14, (maxY/2)+80, (maxY/2)+18)
	if err != nil {
		return err
	}

	_, err = g.SetView("done", (maxX/2)-4, (maxY/2)-17, maxX-39, (maxY/2)+13)
	if err != nil {
		return err
	}

	_, err = g.SetView("notif", 0, 0, (maxY/2)+49, (maxY/2)-18)
	if err != nil {
		return err
	}
	return nil
}

//TODO: return []*gocui.View instead of tuple
func createViews(g *gocui.Gui) (*gocui.View, *gocui.View, *gocui.View, *gocui.View, *gocui.View) {
	maxX, maxY := g.Size()
	lv, err := g.SetView("list", 0, (maxY/2)-17, (maxY/2)+13, (maxY/2)+13)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Cannot Create List view", err)
	}
	tv, err := g.SetView("todo", (maxX/2)-60, (maxY/2)-17, (maxY/2)+52, (maxY/2)+13)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Cannot Create Todo View", err)
	}
	cv, err := g.SetView("cmd", (maxX/5)-1, (maxY/2)+14, (maxY/2)+80, (maxY/2)+18)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Cannot Create Cmd view", err)
	}

	dv, err := g.SetView("done", (maxX/2)-4, (maxY/2)-17, maxX-39, (maxY/2)+13)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Cannot Create DoneTasks view", err)
	}

	nv, err := g.SetView("notif", 0, 0, (maxY/2)+49, (maxY/2)-18)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Cannot Create notif view")
	}
	return lv, tv, cv, nv, dv
}

func configViews(v *gocui.View, title string, x, y int, highlight bool, auto bool) {
	if auto {
		v.Editable = auxilary.Editable
	} else {
		v.Editable = false
	}
	v.SetCursor(x, y)
	v.SelFgColor = gocui.ColorGreen
	v.Highlight = highlight
	v.Title = title
}

func handleViews(g *gocui.Gui, filename string) {

	lv, tv, cv, nv, dv := createViews(g)

	configViews(lv, "Projects", 0, 0, true, false)
	configViews(tv, "Todos", 0, 0, true, true)
	configViews(cv, "Commands", 0, 0, true, true)
	configViews(nv, "Info", 0, 0, false, false)
	configViews(dv, "Done Tasks", 0, 0, true, true)

	_, err := g.SetCurrentView(cv.Name())
	checkError(err)
	auxilary.KeyBindingHandler(g, filename)
}

func main() {
	g, err := InitGui()
	checkError(err)

	g.SetManagerFunc(layout)

	defer g.Close()

	var filename string

	if args := os.Args; len(args) >= 2 {
		if file := args[1]; filepath.Ext(file) == ".json" {
			filename = file
		}
	} else {
		g.Close()
		os.Exit(2)
	}

	handleViews(g, filename)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatal(err)
	}

}
