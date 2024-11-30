package game

import (
	"math/rand"
	"sync"

	"github.com/gdamore/tcell/v2"
)

var shapes = [][][]int{
	{ // I shape (long bar)
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	},
	{ // O shape (square)
		{0, 0, 0, 0},
		{0, 2, 2, 0},
		{0, 2, 2, 0},
		{0, 0, 0, 0},
	},
	{ // T shape
		{0, 0, 0, 0},
		{0, 3, 0, 0},
		{3, 3, 3, 0},
		{0, 0, 0, 0},
	},
	{ // S shape
		{0, 0, 0, 0},
		{0, 4, 4, 0},
		{4, 4, 0, 0},
		{0, 0, 0, 0},
	},
	{ // Z shape
		{0, 0, 0, 0},
		{5, 5, 0, 0},
		{0, 5, 5, 0},
		{0, 0, 0, 0},
	},
	{ // J shape
		{0, 0, 0, 0},
		{6, 0, 0, 0},
		{6, 6, 6, 0},
		{0, 0, 0, 0},
	},
	{ // L shape
		{0, 0, 0, 0},
		{0, 0, 7, 0},
		{7, 7, 7, 0},
		{0, 0, 0, 0},
	},
}

const (
	BOARD_WIDTH  = 10
	BOARD_HEIGHT = 20
)

type NextShapeObserver interface {
	UpdateNext(nextShape [][]int)
}

type GameService struct {
	CurrentShape      *int
	NextShape         *int
	shapes            [][][]int
	nextShapeObserver []NextShapeObserver
	mu                sync.Mutex
}

func NewGameService() *GameService {
	return &GameService{shapes: shapes}
}

func (gs *GameService) InitNextShape() {
	// generate the random number
	randNb := rand.Intn(7)

	gs.NextShape = &randNb
	if gs.CurrentShape == nil {
		gs.CurrentShape = gs.NextShape
	}
	// get the shape table
	shape := gs.shapes[randNb]

	// triger observer next shape
	for i := 0; i < len(gs.nextShapeObserver); i++ {
		observer := gs.nextShapeObserver[i]
		observer.UpdateNext(shape)
	}
}

func (gs *GameService) InitGame() {
	gs.InitNextShape()
}

func (gs *GameService) AddObserver(observer NextShapeObserver) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.nextShapeObserver = append(gs.nextShapeObserver, observer)
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
