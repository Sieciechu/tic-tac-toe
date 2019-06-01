package tictactoe

type GameError string

func (e GameError) Error() string {
	return string(e)
}

