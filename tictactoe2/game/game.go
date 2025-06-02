package game

import (
	"github.com/param108/reinforcement-learning/tictactoe2/board"
	"github.com/param108/reinforcement-learning/tictactoe2/player"
)

type Game struct {
	brd *board.Board
	// 1 - x 2 -o
	players [3]player.Player
	silent  bool // If true, no output is printed
}

func NewGame(player int, xplayer, oplayer player.Player, silent bool) *Game {
	g := &Game{
		brd:    board.NewBoard(player),
		silent: silent,
	}

	g.players[1] = xplayer
	g.players[2] = oplayer

	return g
}

func (g *Game) Play() int {
	for g.brd.CheckWin() == 0 {
		if !g.silent {
			// Print the board
			g.brd.Print()
		}
		x, y, player := g.players[g.brd.NextPlayer()].MakeMove(g.brd)
		g.brd.MakeMove(x, y, player)
	}

	if g.brd.CheckWin() == 1 {
		g.players[1].Win()
		g.players[2].Lose()
	} else if g.brd.CheckWin() == 2 {
		g.players[2].Win()
		g.players[1].Lose()
	}

	return g.brd.CheckWin()
}

func (g *Game) GetBoard() *board.Board {
	return g.brd
}
