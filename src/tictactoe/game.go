package tictactoe

// ticTacToe struct containing needed information about the tic-tac-toe game: players, board, which turn it is
//	(and how much moves left)
type ticTacToe struct {
	players      [2]player
	board        Board
	movingPlayer uint
	movesLeft    uint
}

// simple type for player
type player string

// Board for tic tac toe, matrix 3*3, each field has one of 3 values: 0=empty, 1=player1, 2=player2
// This type contains just data about tic-tac-toe board (it's not a view)
type Board [3][3]uint

// BoardWriter is the interface for writing the board to any output
// BoardWriter writes b Board to an output defined in concrete implementation.
// It returns the number of bytes written and any error encountered that caused the write to stop early.
// Write must not modify the board data, even temporarily.
type BoardWriter interface {
	WriteBoard(b Board) (n int, err error)
}

// move struct contains information about move coordinates
type move struct {
	x uint
	y uint
}

// contains information about the gameResult: if game is finished (still playing, win, draw)
//	and if someone wins the winner
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

// makes move for the player which the current turn belongs to
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

// returns the game result (see gameResult struct)
func (g *ticTacToe) getGameResult() (result gameResult) {

	b := g.board

	const ROW_COUNT = len(b)
	const COL_COUNT = len(b[0])
	const NUMBER_OF_WORKERS = ROW_COUNT + COL_COUNT + 2 // 2 is for two diagonals

	var gameResults = make(chan gameResult, NUMBER_OF_WORKERS)

	for x := 0; x < ROW_COUNT; x++ {
		go g.allFieldsTheSame(gameResults, x, 0, 0, 1) // check rows

	}
	go g.allFieldsTheSame(gameResults, 0, 0, 1, 1) // check first diagonal

	for y := 0; y < COL_COUNT; y++ {
		go g.allFieldsTheSame(gameResults, 0, y, 1, 0) // check cols
	}
	go g.allFieldsTheSame(gameResults, ROW_COUNT-1, 0, -1, 1) // check second diagonal

	for i := 0; i < NUMBER_OF_WORKERS; i++ {
		result = <-gameResults
		if result.isFinished {
			return result
		}
	}

	if 0 == g.movesLeft {
		result.isFinished = true
	}

	return result
}

// This is a support method for getGameResult. It checks if on the board for the given line there is a winner and game is finished.
// How? It traverses n=len(board) fields starting from x,y and increasing coordinates by dx, dy each time.
// And pushes the result to the channel.
func (g *ticTacToe) allFieldsTheSame(res chan<- gameResult, x, y, dx, dy int) {
	b := &g.board

	winningPlayerNumber := b[x][y]

	if 0 == winningPlayerNumber {
		res <- gameResult{false, nil}
		return
	}

	for step, maxSteps := 0, len(b)-1; step < maxSteps; x, y, step = x+dx, y+dy, step+1 {
		if winningPlayerNumber != b[x+dx][y+dy] {
			res <- gameResult{false, nil}
			return
		}
	}

	res <- gameResult{true, &g.players[winningPlayerNumber-1]}
	return
}
