package player

import "github.com/param108/reinforcement-learning/tictactoe2/board"

type Player interface {
	// MakeMove returns x, y, player
	MakeMove(board *board.Board) (int, int, int)
	Win()
	Lose()
	GetPlayer() int
	SetPlayer(player int)
}
