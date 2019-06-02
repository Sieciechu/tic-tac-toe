package tictactoe

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type cliController struct {
	game         *ticTacToe
	writer       io.Writer
	reader       io.Reader
	boardPrinter BoardWriter
}

func NewCliController(g *ticTacToe, w io.Writer) *cliController {
	playersChars := [2]rune{'X', 'O'}
	var a BoardWriter = &asciiBoardPrinter{w, playersChars}

	return &cliController{game: g,
		writer:       w,
		reader:       os.Stdin,
		boardPrinter: a}

}

func (c *cliController) Run() {
	fmt.Fprintln(c.writer, "Tic-tac-toe start")
	c.boardPrinter.WriteBoard(c.game.board)
	fmt.Fprint(c.writer, "When prompt please enter move by typing cooridnates x,y={0,1,2}.\n"+
		"Type for example: 0 0\n"+
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

		moves <- convertTwoNumbersStringToMove(text)
	}

	close(moves)
}

func (c *cliController) moveLoop(moves <-chan move, prompt chan<- int) {

	for aMove := range moves {
		c.game.makeMove(aMove)
		c.boardPrinter.WriteBoard(c.game.board)

		if player := c.game.checkForWinningSituation(); player != nil {
			fmt.Fprintf(c.writer, `The game won player "%s"`+"\n", player.name)
			break
		}

		prompt <- 1
	}
	close(prompt)
}

func convertTwoNumbersStringToMove(s string) move {
	var n []uint
	for _, f := range strings.Fields(s) {
		i, err := strconv.Atoi(f)
		if err == nil {
			n = append(n, uint(i))
		}
	}

	return move{n[0], n[1]}
}
