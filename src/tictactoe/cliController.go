package tictactoe

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Command line interface controller
type cliController struct {
	game         *ticTacToe
	writer       io.Writer
	reader       io.Reader
	boardPrinter BoardWriter
	instructionBoardPrinter BoardWriter
}

// For simplicity this factory-method already has some defaults hardcoded
func NewCliController(g *ticTacToe, w io.Writer) *cliController {
	playersChars := [2]rune{'X', 'O'}

	return &cliController{game: g,
		writer:       w,
		reader:       os.Stdin,
		boardPrinter: &asciiBoardPrinter{w, playersChars},
		instructionBoardPrinter: &asciiInstructionBoardPrinter{w}}

}

func (c *cliController) Run() {
	fmt.Fprintln(c.writer, "Tic-tac-toe start")
	c.boardPrinter.WriteBoard(c.game.board)
	fmt.Fprint(c.writer, "\nWhen prompt please enter move by choosing one of the following fields:\n")
	c.instructionBoardPrinter.WriteBoard(c.game.board)
	fmt.Fprint(c.writer,
		"Type for example: 1\n"+
		"Or press 'q' to quit\n\n")

	moves := make(chan move)
	prompt := make(chan int, 1)
	prompt <- 1

	go c.moveLoop(moves, prompt)

	c.readInputLoop(moves, prompt)
}

func (c *cliController) readInputLoop(moves chan<- move, prompt <-chan int) {

	scanner := bufio.NewScanner(c.reader)

	for range prompt {
		fmt.Fprint(c.writer, "Enter your move: ")

		if false == scanner.Scan() {
			break
		}

		text := scanner.Text()

		if "q" == text {
			break
		}

		moves <- convertNumberStringToMove(text)
	}

	close(moves)
}

func (c *cliController) moveLoop(moves <-chan move, prompt chan<- int) {

	for aMove := range moves {
		c.game.makeMove(aMove)
		c.boardPrinter.WriteBoard(c.game.board)

		gameResult := c.game.getGameResult()

		if !gameResult.isFinished {
			prompt <- 1
			continue
		}

		if nil != gameResult.winningPlayer {
			fmt.Fprintf(c.writer, `The game won player "%s"`+"\n", *gameResult.winningPlayer)
		} else {
			fmt.Fprintf(c.writer, "Draw. Noone won")
		}

		break

	}
	close(prompt)
}

// Support method. It converts field to move on the board
func convertNumberStringToMove(s string) (m move) {
	i, _ := strconv.Atoi(s)

	m.x = uint((i - 1) / 3)
	m.y = uint((i - 1) % 3)

	return
}
