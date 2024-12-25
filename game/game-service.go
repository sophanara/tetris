package game

import (
	"fmt"
	"math/rand"
	signal "socorp/tetris/libs"

	"github.com/gdamore/tcell/v2"
)

var shapes = [][][]int{
	{ // I shape (long bar)
		{0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
	},
	{ // O shape (square)
		{0, 0, 0, 0, 0, 0},
		{0, 0, 2, 2, 0, 0},
		{0, 0, 2, 2, 0, 0},
		{0, 0, 0, 0, 0, 0},
	},
	{ // T shape
		{0, 0, 0, 0, 0, 0},
		{0, 0, 3, 0, 0, 0},
		{0, 3, 3, 3, 0, 0},
		{0, 0, 0, 0, 0, 0},
	},
	{ // S shape
		{0, 0, 0, 0, 0, 0},
		{0, 0, 4, 4, 0, 0},
		{0, 4, 4, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
	},
	{ // Z shape
		{0, 0, 0, 0, 0, 0},
		{0, 5, 5, 0, 0, 0},
		{0, 0, 5, 5, 0, 0},
		{0, 0, 0, 0, 0, 0},
	},
	{ // J shape
		{0, 0, 0, 0, 0, 0},
		{0, 6, 0, 0, 0, 0},
		{0, 6, 6, 6, 0, 0},
		{0, 0, 0, 0, 0, 0},
	},
	{ // L shape
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 7, 0, 0},
		{0, 7, 7, 7, 0, 0},
		{0, 0, 0, 0, 0, 0},
	},
}

const (
	BOARD_WIDTH  = 10
	BOARD_HEIGHT = 15
)

type GameContext struct {
	CurrentShape       *int
	NextShape          *int
	Board              [][]int
	Status             GameStatus
	Score              int
	MaxRow             int
	MaxCol             int
	CurrentRowPosition int
	CurrentColPosition int
}

type GameService struct {
	GameContext     *GameContext
	nextShapeSignal *signal.Signal[[][]int]
	boardSignal     *signal.Signal[[][]int]
	messageSignal   *signal.Signal[string]
	shapes          [][][]int
}

func NewGameService() *GameService {
	// generate the random number
	randNb := rand.Intn(7)

	gameContext := &GameContext{
		Status:    END,
		Score:     0,
		MaxRow:    BOARD_HEIGHT,
		MaxCol:    BOARD_WIDTH,
		NextShape: &randNb,
	}

	// get the shape table
	shape := shapes[randNb]
	nextShapeSignal := signal.NewSignal(shape)
	messageSignal := signal.NewSignal("")

	gameService := &GameService{GameContext: gameContext, shapes: shapes, nextShapeSignal: nextShapeSignal}
	board := gameService.initBoard()
	gameContext.Board = board
	boardSignal := signal.NewSignal(board)

	gameService.boardSignal = boardSignal
	gameService.messageSignal = messageSignal
	return gameService
}

func (gs *GameService) InitNextShape() {
	ctxt := gs.GameContext

	// generate the random number
	randNb := rand.Intn(7)

	ctxt.NextShape = &randNb
	// get the shape table
	shape := gs.shapes[randNb]

	gs.nextShapeSignal.Set(shape)
}

func (gs *GameService) updateCurrentShape() {
	ctxt := gs.GameContext
	if ctxt.CurrentShape == nil {
		ctxt.CurrentShape = ctxt.NextShape
	}
}

func (gs *GameService) InitGame() {
	gs.updateCurrentShape()
}

func (gs *GameService) GetColor(color int) tcell.Color {
	switch color {
	case 1:
		return tcell.ColorYellow
	case 2:
		return tcell.ColorRed
	case 3:
		return tcell.ColorGreen
	case 4:
		return tcell.ColorBlue
	case 5:
		return tcell.ColorViolet
	case 6:
		return tcell.ColorPink
	case 7:
		return tcell.ColorOrange
	default:
		return tcell.ColorBlack
	}
}

func (gs *GameService) getCurrentShape() [][]int {
	curShape := gs.GameContext.CurrentShape
	return gs.shapes[*curShape]
}

// func (gs *GameService) getNextShape() [][]int {
// 	nextShape := gs.GameContext.NextShape
// 	return gs.shapes[*nextShape]
// }

func (gs *GameService) initBoard() [][]int {
	board := make([][]int, BOARD_HEIGHT)
	for i := 0; i < BOARD_HEIGHT; i++ {
		board[i] = make([]int, BOARD_WIDTH)
	}
	return board
}

func (gv *GameService) DropBlock() {
	shape := gv.getCurrentShape()
	maxCol := gv.GameContext.MaxCol
	xpos := (maxCol - len(shape[0])) / 2
	ypos := 0

	gv.GameContext.CurrentColPosition = xpos
	gv.GameContext.CurrentRowPosition = ypos

	gv.ApplyPiecePosition()
}

func (gv *GameService) ApplyPiecePosition() {
	debugText := fmt.Sprintf("updatePiecePosition %d %d\n",
		gv.GameContext.CurrentRowPosition, gv.GameContext.CurrentColPosition)
	gv.messageSignal.Set(debugText)

	shape := gv.getCurrentShape()
	xpos := gv.GameContext.CurrentColPosition
	ypos := gv.GameContext.CurrentRowPosition

	// Copy the previous board add the new shape
	board := make([][]int, BOARD_HEIGHT)
	for i := range board {
		board[i] = make([]int, BOARD_WIDTH)
		copy(board[i], gv.GameContext.Board[i])
	}
	// draw the shape
	for r := 0; r < len(shape); r++ {
		for c := 0; c < len(shape[r]); c++ {
			value := shape[r][c]
			if value != 0 {
				board[r+ypos][c+xpos] = value
			}
		}
	}

	gv.boardSignal.Set(board)
}

// Check if all the cell of the current piece can be moved
func (gv *GameService) canMovePiece(nbRow int, nbCol int) bool {
	ctxt := gv.GameContext
	curRowPosition := ctxt.CurrentRowPosition
	curColPosition := ctxt.CurrentColPosition
	curShape := gv.getCurrentShape()
	board := ctxt.Board

	gv.messageSignal.Set(fmt.Sprintf("canMovePiece %d, %d \n", curRowPosition, curColPosition))
	for r := 0; r < len(curShape); r++ {
		for c := 0; c < len(curShape[r]); c++ {
			// ignore blank cell
			if curShape[r][c] == 0 {
				continue
			}
			newRowPosition := curRowPosition + r + nbRow
			newColPosition := curColPosition + c + nbCol
			if newRowPosition < 0 ||
				newColPosition < 0 ||
				newRowPosition >= BOARD_HEIGHT ||
				newColPosition >= BOARD_WIDTH ||
				board[newRowPosition][newColPosition] != 0 {
				return false
			}
		}
	}
	return true
}

func (gv *GameService) movePiece(bnRow int, nbCol int) {
	// check if the piece can be moved
	if !gv.canMovePiece(bnRow, nbCol) {
		return
	}
	// move the piece
	gv.GameContext.CurrentRowPosition += bnRow
	gv.GameContext.CurrentColPosition += nbCol
	gv.ApplyPiecePosition()
}

// Move the current piece down
func (gv *GameService) MoveDown() {
	gv.movePiece(1, 0)
}

// Move the current piece left
func (gv *GameService) MoveLeft() {
	gv.movePiece(0, -1)
}

// move the current piece right
func (gv *GameService) MoveRight() {
	gv.movePiece(0, 1)
}
