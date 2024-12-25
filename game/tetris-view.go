package game

import (
	"time"

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

type TetrisView struct {
	ScoreBox        ScoreView
	PreviewBox      PreviewView
	GameBox         GameView
	DebugBox        DebugScreen
	OuterBox        *tview.Box
	GameService     *GameService
	NextShapeChanel <-chan [][]int
	BordChanel      <-chan [][]int
	Application     tview.Application
}

func NewTetrisView() *TetrisView {
	outerBox := tview.NewBox()
	tetrisService := NewGameService()
	gameBox := NewGameView(tetrisService)
	previewBox := NewPreviewView(tetrisService)
	scoreBox := NewScoreView()
	debugbox := NewDebugScreen(tetrisService)

	tetrisView := &TetrisView{
		Application: *tview.NewApplication(),
		OuterBox:    outerBox,
		GameService: tetrisService,
		GameBox:     *gameBox,
		PreviewBox:  *previewBox,
		ScoreBox:    *scoreBox,
		DebugBox:    *debugbox,
	}

	return tetrisView
}

// UpdateNext implements NextShapeObserver.
func (tv *TetrisView) UpdateNext(nextShape [][]int) {
	// update the preview
	tv.PreviewBox.UpdateView(nextShape)
}

func (tv *TetrisView) StartGame() {
	tv.GameService.InitGame()
	tv.GameService.DropBlock()
	// Set up key bindings
	tv.Application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()

		switch key {
		case tcell.KeyEscape:
			tv.Application.Stop()
		case tcell.KeyDown:
			tv.GameService.MoveDown()
		case tcell.KeyLeft:
			tv.GameService.MoveLeft()
		case tcell.KeyRight:
			tv.GameService.MoveRight()
		}
		return event
	})

	// Start the game loop
	go tv.startPeriodicUpdate()
}

func (tv *TetrisView) StartWindow() {
	gameBox := tv.GameBox.GetView()
	previewBox := tv.PreviewBox.GetView()
	scoreView := tv.ScoreBox.GetView()
	debugScreen := tv.DebugBox.GetView()

	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(scoreView, 0, 1, false).
		AddItem(debugScreen, 0, 1, false)
	// Add the game box to the new flex element
	flex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(previewBox, 0, 1, false).
		AddItem(gameBox, 0, 2, false).
		AddItem(rightFlex, 0, 1, false)

	tv.OuterBox.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
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

	ctxt := tv.GameService.GameContext
	tv.GameBox.InitTable(ctxt.MaxRow, ctxt.MaxCol)
	tv.Application.SetRoot(tv.OuterBox, true)
}

func (tv *TetrisView) startPeriodicUpdate() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		tv.GameService.MoveDown()
		go tv.Application.Draw()
	}
}
