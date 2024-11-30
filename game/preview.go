package game

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PreviewView struct {
	GameService *GameService
	preview     tview.Table
}

// UpdateNext implements NextShapeObserver.
func (p *PreviewView) UpdateNext(nextShape [][]int) {
	// update the preview
	p.updateView(nextShape)
}

func NewPreviewView(gameService *GameService) *PreviewView {
	preView := &PreviewView{GameService: gameService}
	gameService.AddObserver(preView)
	return preView
}

func (p *PreviewView) updateView(shape [][]int) {
	preview := p.preview
	rLength := len(shape)
	cLength := len(shape[0])

	for r := 0; r < rLength; r++ {
		for c := 0; c < cLength; c++ {
			color := p.GameService.GetColor(shape[r][c])
			preview.SetCell(r, c, tview.NewTableCell("   ").
				SetMaxWidth(20).
				SetBackgroundColor(color))
		}
	}
	preview.SetTitle("block")
}

func (p *PreviewView) GetView() *tview.Box {
	previewBox := tview.NewBox().SetBackgroundColor(tcell.ColorBlack).SetBorder(true).SetTitle("Preview")

	preview := tview.NewTable().SetBorders(true)
	p.preview = *preview
	p.GameService.InitGame()

	previewViewFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(preview, 0, 1, true)
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

	return previewBox
}
