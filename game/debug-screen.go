package game

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type DebugScreen struct {
	messageChanel <-chan string
}

func NewDebugScreen(gs *GameService) *DebugScreen {
	messsageChanel := gs.messageSignal.Subscribe()
	return &DebugScreen{messageChanel: messsageChanel}
}

func (ds *DebugScreen) GetView() *tview.TextView {
	scoreView := tview.NewTextView()
	scoreView.SetBackgroundColor(tcell.ColorBlack).SetBorder(true).SetTitle("Debug console")

	go func() {
		for msg := range ds.messageChanel {
			fmt.Fprintln(scoreView, msg)
		}
	}()

	return scoreView
}
