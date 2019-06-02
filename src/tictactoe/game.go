package tictactoe

import (
	"fmt"
)

type ticTacToe struct {
	players      [2]player
	board        Board
	movingPlayer uint
	movesLeft    uint
}

type player struct {
	name string
	char string
}

type Board [3][3]uint

type BoardWriter interface {
	WriteBoard(b Board) (n int, err error)
}

type move struct {
	x uint
	y uint
}

func NewTicTacToe() *ticTacToe {
	p1 := player{name: "player1", char: "X"}
	p2 := player{name: "player2", char: "O"}
	b := Board{}
	return &ticTacToe{players: [2]player{p1, p2},
		movingPlayer: 1,
		board:        b,
		movesLeft:    uint(len(b) * len(b[0]))}

}

func (g *ticTacToe) makeMove(m move) error {

	if 0 == g.movesLeft {
		return GameError("There are no more free moves")
	}

	if player := g.board[m.x][m.y]; player != 0 {
		return GameError("Field is already marked by player " + string(player))
	}

	g.board[m.x][m.y] = g.movingPlayer
	g.movesLeft--

	if 1 == g.movingPlayer {
		g.movingPlayer = 2
	} else {
		g.movingPlayer = 1
	}

	return nil
}

func (g *ticTacToe) checkForWinningSituation() *player {

	b := g.board

	var areSame = false
	var winningPlayer *player = nil

	for x, rowCount := 0, len(b); x < rowCount; x++ {
		areSame, winningPlayer = g.allFieldsTheSame(x, 0, 0, 1) // check rows
		if areSame {
			return winningPlayer
		}
	}
	areSame, winningPlayer = g.allFieldsTheSame(0, 0, 1, 1) // check first diagonal
	if areSame {
		return winningPlayer
	}

	for y, colCount := 0, len(b[0]); y < colCount; y++ {
		areSame, winningPlayer = g.allFieldsTheSame(0, y, 1, 0) // check cols
		if areSame {
			return winningPlayer
		}
	}
	areSame, winningPlayer = g.allFieldsTheSame(len(b)-1, 0, -1, 1) // check second diagonal
	if areSame {
		return winningPlayer
	}

	if 0 == g.movesLeft && !areSame {
		fmt.Println("Draw. Noone won")
		return nil
	}

	return nil
}

func (g *ticTacToe) allFieldsTheSame(x, y, dx, dy int) (areSame bool, winningPlayer *player) {
	b := &g.board

	winningPlayerNumber := b[x][y]

	if 0 == winningPlayerNumber {
		return false, nil
	}

	for step, maxSteps := 0, len(b)-1; step < maxSteps; x, y, step = x+dx, y+dy, step+1 {
		if winningPlayerNumber != b[x+dx][y+dy] {
			return false, nil
		}
	}

	return true, &g.players[winningPlayerNumber-1]
}
