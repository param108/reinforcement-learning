package board

type Board struct {
	board  [9]int
	start  int
	player int
}

func NewBoard(start int) *Board {
	return &Board{
		start:  start,
		player: start,
	}
}

func (b *Board) Get() [9]int {
	var newBoard [9]int
	copy(newBoard, b.board)
	return newBoard
}

func (b *Board) NextPlayer() {
	return b.player
}

func (b *Board) ID() int64 {
	var id int64 = 0

	// Encode board (base-3)
	for i := 0; i < 9; i++ {
		id = id*3 + int64(b.board[i])
	}

	// Map start and player from (1,2) to (0,1)
	startBit := int64(b.start - 1)
	playerBit := int64(b.player - 1)

	// Add start (base-2)
	id = id*2 + startBit
	// Add player (base-2)
	id = id*2 + playerBit

	return id
}

// TryMove attempts to place player's mark at (x, y).
// Returns true if move succeeded, false if the cell was not empty.
func (b *Board) TryMove(x, y, player int) ([9]int, bool) {
	idx := x + 3*y
	if idx < 0 || idx >= 9 {
		return nil, false // Out of bounds
	}
	if b.board[idx] != 0 {
		return nil, false // Cell already taken
	}

	var newBoard [9]int
	copy(newBoard, b.board)

	newBoard[idx] = player
	return newBoard, true
}

// TryMove attempts to place player's mark at (x, y).
// Returns true if move succeeded, false if the cell was not empty.
func (b *Board) MakeMove(x, y, player int) bool {

	if b.player != player {
		return false
	}

	idx := x + 3*y
	if idx < 0 || idx >= 9 {
		return false // Out of bounds
	}
	if b.board[idx] != 0 {
		return false // Cell already taken
	}

	b.board[idx] = player

	// toggle the next player
	b.player = 3 - b.player
	return true
}

// CheckWin returns:
//
//	0 if the game continues (no win, no draw),
//	1 if player 1 wins,
//	2 if player 2 wins,
//	3 if the game is a draw
func (b *Board) CheckWin() int {
	winPatterns := [8][3]int{
		{0, 1, 2}, // rows
		{3, 4, 5},
		{6, 7, 8},
		{0, 3, 6}, // columns
		{1, 4, 7},
		{2, 5, 8},
		{0, 4, 8}, // diagonals
		{2, 4, 6},
	}

	// Check for winning pattern
	for _, pattern := range winPatterns {
		a, bIdx, c := pattern[0], pattern[1], pattern[2]
		if b.board[a] != 0 && b.board[a] == b.board[bIdx] && b.board[bIdx] == b.board[c] {
			return b.board[a]
		}
	}

	// Check for draw (no zeros left)
	for i := 0; i < 9; i++ {
		if b.board[i] == 0 {
			return 0 // Game continues, empty spots left
		}
	}

	return 3 // Draw
}
