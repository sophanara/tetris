package game

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type GameView struct {
	TetrisScreen tview.Table
	GameService  *GameService
	BoardChanel  <-chan [][]int
}

func NewGameView(gameService *GameService) *GameView {
	boardChanel := gameService.boardSignal.Subscribe()
	return &GameView{GameService: gameService, BoardChanel: boardChanel}
}

func (gv *GameView) InitTable(maxRow int, maxCol int) {
	ts := gv.TetrisScreen
	ts.SetBorder(true)

	gv.TetrisScreen.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		// Calculate inner area (accounting for border)
		wspace := width / 4
		innerX := x + wspace
		innerY := y + 1
		// Draw the flex layout in the inner area
		return innerX, innerY, width, height
	})

	// Initialize the chanel handler that will update the view
	go func() {
		for board := range gv.BoardChanel {
			gv.UpdateView(board)
		}
	}()
}

func (gv *GameView) GetView() *tview.Box {
	gameBox := tview.NewBox().SetBackgroundColor(tcell.ColorRebeccaPurple).SetBorder(true).SetTitle("Go, Tetris!")
	tetrisScreen := tview.NewTable().SetBorders(true)
	gv.TetrisScreen = *tetrisScreen

	cmdPanel := tview.NewBox()

	gameBoxFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tetrisScreen, 0, 1, false)

	gameBoxRowFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(gameBoxFlex, 0, 100, false).
		AddItem(nil, 0, 2, false).
		AddItem(cmdPanel, 0, 10, true)

	gameBox.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		// Calculate inner area (accounting for border)
		innerX := x + 1
		innerY := y + 2
		innerWidth := width - 2
		innerHeight := height - 10
		// Draw the flex layout in the inner area
		gameBoxRowFlex.SetRect(innerX, innerY, innerWidth, innerHeight)
		gameBoxRowFlex.Draw(screen)
		return x, y, width, height
	})

	return gameBox
}

func (gv *GameView) UpdateView(board [][]int) {
	for r := 0; r < len(board); r++ {
		for c := 0; c < len(board[r]); c++ {
			gv.TetrisScreen.SetCell(r, c, tview.NewTableCell("   ").
				SetBackgroundColor(gv.GameService.GetColor(board[r][c])))
		}
	}
}
