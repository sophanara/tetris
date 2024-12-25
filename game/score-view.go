package game

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ScoreView struct{}

func NewScoreView() *ScoreView {
	return &ScoreView{}
}

func (sv *ScoreView) GetView() *tview.Box {
	scoreView := tview.NewBox().SetBackgroundColor(tcell.ColorBlack).SetBorder(true)
	return scoreView
}
