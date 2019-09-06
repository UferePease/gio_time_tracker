package main

import (
	"context"
	"flag"
	"fmt"
	"gioui.org/ui"
	"gioui.org/ui/app"
	"gioui.org/ui/key"
	"gioui.org/ui/layout"
	"log"
	"net/http"
	"os"
)

type App struct {
	w *app.Window
	ui *TimerUI

	ctx           context.Context
	ctxCancel     context.CancelFunc
}

var (
	profile = flag.Bool("profile", false, "serve profiling data at http://localhost:6060")
)


func main() {
	flag.Parse()
	initProfiling()
	go func() {
		w := app.NewWindow(
			app.WithWidth(ui.Dp(400)),
			app.WithHeight(ui.Dp(600)),
			app.WithTitle("Time Tracka"),
		)

		if err := newApp(w).run(); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func initProfiling() {
	if !*profile {
		return
	}
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

func newApp(w *app.Window) *App {
	a := &App{
		w:	w,
	}

	a.ui = newUI()
	return a
}

func (a *App) run() error {
	//a.ui.profiling = *stats
	ops := new(ui.Ops)
	var cfg app.Config

	for {
		select {
		case e := <-a.w.Events():
			switch e := e.(type) {
			case key.Event:
				switch e.Name {
				case key.NameEscape:
					os.Exit(0)
				}
			case app.DestroyEvent:
				return e.Err
			case app.UpdateEvent:
				ops.Reset()
				cfg = e.Config
				cs := layout.RigidConstraints(e.Size)
				fmt.Println(a.ui.descLabel.Text)
				fmt.Println(a.ui.projectLabel.Text)
				a.ui.Layout(&cfg, a.w.Queue(), ops, cs)
				a.w.Update(ops)
			}
		}
	}

}



