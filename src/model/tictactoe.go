package model

import (
	"fmt"
	"io"
	"os"
)

type GameError string

func (e GameError) Error() string {
	return string(e)
}

type TicTacToe struct {
	players      [2]player
	board        Board
	movingPlayer uint
	movesLeft    uint
	writer 		 io.Writer
}

func NewTicTacToe() TicTacToe {
	p1 := player{name: "player1", char: "X"}
	p2 := player{name: "player2", char: "O"}
	b := Board{}
	return TicTacToe{players: [2]player{p1, p2},
		movingPlayer: 1,
		board:        b,
		movesLeft:    uint(len(b) * len(b[0])),
		writer: os.Stdout}

}

func (g *TicTacToe) MakeMove(x, y uint) error {

	if 0 == g.movesLeft {
		return GameError("There are no more free moves")
	}

	if player := g.board[x][y]; player != 0 {
		return GameError("Field is already marked by player " + string(player))
	}

	g.board[x][y] = g.movingPlayer
	g.movesLeft--

	if 1 == g.movingPlayer {
		g.movingPlayer = 2
	} else {
		g.movingPlayer = 1
	}

	return nil
}

func (g *TicTacToe) CheckForWinningSituation() {

	b := &g.board

	var areSame = false
	var winningPlayer *player = nil

	for x, rowCount := 0, len(b); x < rowCount; x++ {
		areSame, winningPlayer = g.allFieldsTheSame(x, 0, 0, 1) // check rows
		if areSame {
			fmt.Printf("The game won player %s\n", winningPlayer.name)
			return
		}
	}
	areSame, winningPlayer = g.allFieldsTheSame(0, 0, 1, 1) // check first diagonal
	if areSame {
		fmt.Printf("The game won player %s\n", winningPlayer.name)
		return
	}

	for y, colCount := 0, len(b[0]); y < colCount; y++ {
		areSame, winningPlayer = g.allFieldsTheSame(0, y, 1, 0) // check cols
		if areSame {
			fmt.Printf("The game won player %s\n", winningPlayer.name)
			return
		}
	}
	areSame, winningPlayer = g.allFieldsTheSame(len(b)-1, 0, -1, 1) // check second diagonal
	if areSame {
		fmt.Printf("The game won player %s\n", winningPlayer.name)
		return
	}

	if 0 == g.movesLeft && !areSame {
		fmt.Println("Draw. Noone won")
		return
	}
}

func (g *TicTacToe) allFieldsTheSame(x, y, dx, dy int) (areSame bool, winningPlayer *player) {
	b := &g.board

	winningPlayerNumber := b[x][y]

	if 0 == winningPlayerNumber {
		return false, nil
	}

	for step, maxSteps := 0, len(b)-1; step < maxSteps; x,y, step = x+dx, y+dy, step+1 {
		if winningPlayerNumber != b[x+dx][y+dy] {
			return false, nil
		}
	}

	return true, &g.players[winningPlayerNumber-1]
}

func (g *TicTacToe) PrintWithWriter() {
	g.writer.Write([]byte(g.getBoard()))
}

func (g *TicTacToe) getBoard() string {
	b := &g.board

	board := ""

	for row, rowsWithWalls := 0, len(b)+2; row < rowsWithWalls; row++ {
		for col, colsWithWalls := 0, len(b[row])+2; col < colsWithWalls; col++ {
			if (0 == row || row == rowsWithWalls-1) || (0 == col || col == colsWithWalls-1) {
				board += fmt.Sprint("#")
			} else {

				char := " "
				if player := b[row-1][col-1]; player != 0 {
					char = g.players[player-1].char
				}

				board += fmt.Sprint(char)
			}
		}
		board += fmt.Sprint("\n")
	}

	return board
}

func (g *TicTacToe) PrintBoard() {
	b := &g.board

	for i, height := 0, len(b)+2; i < height; i++ {
		for k, width := 0, len(b[i])+2; k < width; k++ {
			if (0 == i || i == height-1) || (0 == k || k == width-1) {
				fmt.Print("#")
			} else {

				char := " "
				if player := b[i-1][k-1]; player != 0 {
					char = g.players[player-1].char
				}

				fmt.Print(char)
			}
		}
		fmt.Print("\n")
	}
}

func (g *TicTacToe) SetPlayer(playerNumber int, playerName string, char string) error {
	if playerNumber > 2 || playerNumber < 1 {
		return GameError("Player number must be 1 or 2")
	}

	g.players[playerNumber].char = char
	g.players[playerNumber].name = playerName

	return nil
}

type player struct {
	name string
	char string
}

type Board [3][3]uint
