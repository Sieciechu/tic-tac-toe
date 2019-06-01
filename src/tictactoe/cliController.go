package tictactoe

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
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

	ch := make(chan move)
	defer close(ch)

	go c.moveLoop(ch)

	c.readInputLoop(ch)
}

func (c *cliController) readInputLoop(ch chan<- move) {

	scanner := bufio.NewScanner(c.reader)

	for {
		time.Sleep(500 * time.Millisecond)
		fmt.Fprint(c.writer, "Enter your move: ")

		if false == scanner.Scan() {
			break
		}

		text := scanner.Text()

		if "q" == text {
			break
		}

		ch <- convertTwoNumbersStringToMove(text)
	}
}

func (c *cliController) moveLoop(ch <-chan move) {

	for aMove := range ch {
		c.game.makeMove(aMove)
		c.boardPrinter.WriteBoard(c.game.board)
		c.game.checkForWinningSituation()

	}
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
