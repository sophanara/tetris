package main

import (
	"fmt"
	"socorp/tetris/game"
)

func main() {
	tetrisView := game.NewTetrisView()
	ctxt := tetrisView.GameContext
	fmt.Printf("tetris view context, status:%s, score:%d \n", ctxt.Status, ctxt.Score)
	fmt.Println("Starting game....")
	tetrisView.StartWindow()
}
