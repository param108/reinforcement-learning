package player

import (
	"github.com/param108/reinforcement-learning/tictactoe2/board"
)

type MinimaxPlayer struct {
	player int
}

func NewMinimaxPlayer(player int) *MinimaxPlayer {
	return &MinimaxPlayer{
		player: player,
	}
}

func (p *MinimaxPlayer) GetPlayer() int {
	return p.player
}

func (p *MinimaxPlayer) SetPlayer(player int) {
	p.player = player
}

func (p *MinimaxPlayer) minimax(board *board.Board, brdArray [9]int, maximizing bool) int {
	win := board.CalcWin(brdArray)
	if win == p.player {
		return 1
	} else if win == 3-p.player {
		return -1
	} else if win == 3 {
		return 0
	}

	if maximizing {
		maxEval := -2
		// more moves available
		actions := board.CalcPossibleMoves(brdArray)
		for _, action := range actions {
			newBoard := brdArray
			newBoard[action.X+3*action.Y] = p.player
			eval := p.minimax(board, newBoard, false)
			if eval > maxEval {
				maxEval = eval
			}
		}

		return maxEval
	} else {
		minEval := 2
		// more moves available
		actions := board.CalcPossibleMoves(brdArray)
		for _, action := range actions {
			newBoard := brdArray
			newBoard[action.X+3*action.Y] = 3 - p.player
			eval := p.minimax(board, newBoard, true)
			if eval < minEval {
				minEval = eval
			}
		}

		return minEval
	}
}

func (p *MinimaxPlayer) MakeMove(board *board.Board) (int, int, int) {
	brdArray := board.Get()
	actions := board.GetPossibleMoves()
	maxEval := -2
	maxIdx := 0
	for idx, action := range actions {
		newBoard := brdArray
		newBoard[action.X+3*action.Y] = action.Player
		eval := p.minimax(board, newBoard, false)
		if eval > maxEval {
			maxEval = eval
			maxIdx = idx
		}
	}

	return actions[maxIdx].X, actions[maxIdx].Y, p.player
}

func (p *MinimaxPlayer) Win() {
}

func (p *MinimaxPlayer) Lose() {
}
