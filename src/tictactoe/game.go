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

// This is a support method. It checks if on the board for the given line there is a winner and game is finished
// res : channel for the result
// x : the row of the board to check
// y : the column of the board to check
// dx : delta-x, the amount by which x will be increased each step
// dy : delta-y, the amount by which y will be increased each step
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
