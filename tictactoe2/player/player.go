package player

type Player interface {
	// MakeMove returns x, y, player
	MakeMove(board *board.Board) (int, int, int)
}
