package main

import (
	"./model"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type move struct {
	x uint
	y uint
}

func main() {

	game := model.NewGame()
	game.PrintWithWriter()

	ch := make(chan move)

	go moveLoop(&game, ch)

	fmt.Print("Tic-tac-toe start.\n" +
		"When prompt please enter move by typing cooridnates x,y={0,1,2}.\n" +
		"Type for example: 0 0\n" +
		"Or press 'q' to quit\n\n")

	readInputLoop(os.Stdin, ch)

}

func readInputLoop(r io.Reader, ch chan<- move) {
	scanner := bufio.NewScanner(r)

	defer close(ch)

	for {
		time.Sleep(500 * time.Millisecond)
		fmt.Print("Enter your move: ")


		if false == scanner.Scan() {
			break
		}

		text := scanner.Text()

		_ = text

		if "q" == text {
			break
		}

		ch <- convertTwoNumbersStringToMove(text)
	}
}

func moveLoop(g *model.Game, ch <-chan move) {

	for aMove := range ch {
		g.MakeMove(aMove.x, aMove.y)
		g.PrintWithWriter()
		g.CheckForWinningSituation()

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
