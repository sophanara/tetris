package game

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PreviewView struct {
	GameService     *GameService
	Preview         tview.Table
	NextShapeChanel <-chan [][]int
}

func NewPreviewView(gameService *GameService) *PreviewView {
	nextShapeChanel := gameService.nextShapeSignal.Subscribe()
	preView := &PreviewView{GameService: gameService, NextShapeChanel: nextShapeChanel}
	return preView
}

func (p *PreviewView) UpdateView(shape [][]int) {
	preview := p.Preview
	rLength := len(shape)
	cLength := len(shape[0])

	for r := 0; r < rLength; r++ {
		for c := 0; c < cLength; c++ {
			color := p.GameService.GetColor(shape[r][c])
			preview.SetCell(r, c, tview.NewTableCell("   ").
				SetBackgroundColor(color))
		}
	}
}

func (p *PreviewView) GetView() *tview.Box {
	previewBox := tview.NewBox().SetBackgroundColor(tcell.ColorBlack).SetBorder(true).SetTitle("Preview")

	preview := tview.NewTable().SetBorders(true)
	p.Preview = *preview
	previewViewFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(preview, 0, 1, false)
	previewBox.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		// Calculate inner area (accounting for border)
		innerX := x + 10
		innerY := y + 5
		innerWidth := width - 20
		innerHeight := height - 10
		// Draw the flex layout in the inner area
		previewViewFlex.SetRect(innerX, innerY, innerWidth, innerHeight)
		previewViewFlex.Draw(screen)
		return x, y, width, height
	})

	//initialize the chanel handler that will update the view
	go func() {
		for chanel := range p.NextShapeChanel {
			p.UpdateView(chanel)
		}
	}()
	return previewBox
}
