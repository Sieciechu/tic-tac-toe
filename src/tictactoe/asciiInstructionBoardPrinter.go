package tictactoe

import (
	"fmt"
	"io"
)

// asciiInstructionBoardPrinter implements WriteBoard.
// It prints the board for instruction purposes of cliController
type asciiInstructionBoardPrinter struct {
	writer  io.Writer
}

func (a *asciiInstructionBoardPrinter) WriteBoard(b Board) (n int, err error) {
	board := ""

	field := 1
	for row, rowsWithWalls := 0, len(b)+2; row < rowsWithWalls; row++ {
		for col, colsWithWalls := 0, len(b[row])+2; col < colsWithWalls; col++ {
			if (0 == row || row == rowsWithWalls-1) || (0 == col || col == colsWithWalls-1) {
				board += fmt.Sprint("#")
			} else {
				board += fmt.Sprint(field)
				field++
			}
		}
		board += fmt.Sprint("\n")
	}

	return a.writer.Write([]byte(board))
}
