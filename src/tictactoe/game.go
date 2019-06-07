package tictactoe

type ticTacToe struct {
	players      [2]player
	board        Board
	movingPlayer uint
	movesLeft    uint
}

type player string

// Board for tic tac toe, each field has one of 3 values: 0=empty, 1=player1, 2=player2
type Board [3][3]uint

// BoardWriter is the interface for writing the board to any output
// BoardWriter writes b Board to an output defined in concrete implementation.
// It returns the number of bytes written and any error encountered that caused the write to stop early.
// Write must not modify the slice data, even temporarily.
type BoardWriter interface {
	WriteBoard(b Board) (n int, err error)
}

type move struct {
	x uint
	y uint
}

type gameResult struct {
	isFinished    bool
	winningPlayer *player
}

// For simplicity this factory-method already has some defaults hardcoded
func NewTicTacToe() *ticTacToe {
	p1 := player("player1")
	p2 := player("player2")
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

func (g *ticTacToe) getGameResult() gameResult {

	b := g.board

	var areSame = false
	var winningPlayer *player = nil

	for x, rowCount := 0, len(b); x < rowCount; x++ {
		areSame, winningPlayer = g.allFieldsTheSame(x, 0, 0, 1) // check rows
		if areSame {
			return gameResult{true, winningPlayer}
		}
	}
	areSame, winningPlayer = g.allFieldsTheSame(0, 0, 1, 1) // check first diagonal
	if areSame {
		return gameResult{true, winningPlayer}
	}

	for y, colCount := 0, len(b[0]); y < colCount; y++ {
		areSame, winningPlayer = g.allFieldsTheSame(0, y, 1, 0) // check cols
		if areSame {
			return gameResult{true, winningPlayer}
		}
	}
	areSame, winningPlayer = g.allFieldsTheSame(len(b)-1, 0, -1, 1) // check second diagonal
	if areSame {
		return gameResult{true, winningPlayer}
	}

	if 0 == g.movesLeft && !areSame {
		return gameResult{true, nil}
	}

	return gameResult{false, nil}
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
