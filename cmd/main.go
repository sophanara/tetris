package main

import (
	"fmt"
	"socorp/tetris/game"
)

func main() {
	tetrisView := game.NewTetrisView()
	fmt.Println("Starting game....")
	tetrisView.StartWindow()
	tetrisView.StartGame()
	if err := tetrisView.Application.Run(); err != nil {
		panic(err)
	}
}
