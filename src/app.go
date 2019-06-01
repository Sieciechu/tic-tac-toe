package main

import (
	"./tictactoe"
	"os"
)

func main() {

	game := tictactoe.NewTicTacToe()
	var controller = tictactoe.NewCliController(&game, os.Stdout)

	controller.Run()
}
