package main

import (
	"./tictactoe"
	"os"
)

func main() {

	game := tictactoe.NewTicTacToe()
	controller := tictactoe.NewCliController(game, os.Stdout)

	controller.Run()
}
