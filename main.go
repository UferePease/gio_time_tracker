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
	"os"
)

type App struct {
	w *app.Window
	ui *TimerUI

	ctx           context.Context
	ctxCancel     context.CancelFunc
}

func main() {
	flag.Parse()

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



