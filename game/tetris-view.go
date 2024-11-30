package game

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type GameStatus int

const (
	END GameStatus = iota
	START
)

func (status GameStatus) String() string {
	return [...]string{"Start", "end"}[status]
}

type GameContext struct {
	Status GameStatus
	Score  int
}

type TetrisView struct {
	GameService *GameService
	GameContext GameContext
}

func NewTetrisView() *TetrisView {
	gameContext := &GameContext{
		Status: END,
		Score:  0,
	}
	tetrisService := NewGameService()
	tetrisView := &TetrisView{
		GameContext: *gameContext,
		GameService: tetrisService,
	}
	return tetrisView
}

func (tv *TetrisView) StartWindow() {
	outerBox := tview.NewBox()

	gameBox := tview.NewBox().SetBackgroundColor(tcell.ColorBlack).SetBorder(true).SetTitle("Go, Tetris!")
	previewBox := NewPreviewView(tv.GameService).GetView()
	scoreView := tview.NewBox().SetBackgroundColor(tcell.ColorBlack).SetBorder(true)

	// Add the game box to the new flex element
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(previewBox, 0, 1, true).
		AddItem(gameBox, 0, 2, true).
		AddItem(scoreView, 0, 1, true)

	outerBox.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		// Calculate inner area (accounting for border)
		innerX := x + 1
		innerY := y + 1
		innerWidth := width - 2
		innerHeight := height - 2

		// Draw the flex layout in the inner area
		flex.SetRect(innerX, innerY, innerWidth, innerHeight)
		flex.Draw(screen)

		return x, y, width, height
	})

	tetrisScreen := tview.NewBox().SetBackgroundColor(tcell.ColorDarkSlateGray)
	gameBoxFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tetrisScreen, 0, 1, true)

	gameBox.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		// Calculate inner area (accounting for border)
		innerX := x + 10
		innerY := y + 5
		innerWidth := width - 20
		innerHeight := height - 10
		// Draw the flex layout in the inner area
		gameBoxFlex.SetRect(innerX, innerY, innerWidth, innerHeight)
		gameBoxFlex.Draw(screen)
		return x, y, width, height
	})

	// Draw the outerBox full screen
	if err := tview.NewApplication().SetRoot(outerBox, true).Run(); err != nil {
		panic(err)
	}
}
