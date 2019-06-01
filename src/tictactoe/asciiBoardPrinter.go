package tictactoe

import (
	"fmt"
	"io"
)

type asciiBoardPrinter struct {
	writer  io.Writer
	players [2]rune
}

func (a *asciiBoardPrinter) WriteBoard(b Board) (n int, err error) {
	board := ""

	for row, rowsWithWalls := 0, len(b)+2; row < rowsWithWalls; row++ {
		for col, colsWithWalls := 0, len(b[row])+2; col < colsWithWalls; col++ {
			if (0 == row || row == rowsWithWalls-1) || (0 == col || col == colsWithWalls-1) {
				board += fmt.Sprint("#")
			} else {

				char := " "
				if player := b[row-1][col-1]; player != 0 {
					char = string(a.players[player-1])
				}

				board += fmt.Sprint(char)
			}
		}
		board += fmt.Sprint("\n")
	}

	return a.writer.Write([]byte(board))
}
