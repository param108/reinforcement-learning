package board

import "fmt"

type Action struct {
	X      int // X coordinate (0-2)
	Y      int // Y coordinate (0-2)
	Player int // Player number (1 or 2
}

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
	for i := 0; i < 9; i++ {
		newBoard[i] = b.board[i]
	}
	return newBoard
}

func (b *Board) GetStart() int {
	return b.start
}

func (b *Board) GetPossibleMoves() []Action {
	var actions []Action
	for i := 0; i < 9; i++ {
		if b.board[i] == 0 { // Cell is empty
			x := i % 3
			y := i / 3
			actions = append(actions, Action{X: x, Y: y, Player: b.player})
		}
	}
	return actions
}

func (b *Board) CalcPossibleMoves(brd [9]int) []Action {
	var actions []Action
	for i := 0; i < 9; i++ {
		if brd[i] == 0 { // Cell is empty
			x := i % 3
			y := i / 3
			actions = append(actions, Action{X: x, Y: y})
		}
	}
	return actions
}

func (b *Board) NextPlayer() int {
	return b.player
}

func (b *Board) CalcID(brd [9]int, start int, player int) int64 {
	// Calculate a unique ID for the board state
	// This is a base-3 encoding of the board, with additional bits for start and player
	var id int64 = 0

	// Encode board (base-3)
	for i := 0; i < 9; i++ {
		id = id*3 + int64(brd[i])
	}

	// Map start and player from (1,2) to (0,1)
	startBit := int64(start - 1)
	playerBit := int64(player - 1)

	// Add start (base-2)
	id = id*2 + startBit
	// Add player (base-2)
	id = id*2 + playerBit

	return id
}

func (b *Board) ID() int64 {
	return b.CalcID(b.board, b.start, b.player)
}

// TryMove attempts to place player's mark at (x, y).
// Returns true if move succeeded, false if the cell was not empty.
func (b *Board) TryMove(brd [9]int, x, y, player int) ([9]int, bool) {
	ret := [9]int{}
	idx := x + 3*y
	if idx < 0 || idx >= 9 {
		return ret, false // Out of bounds
	}
	if brd[idx] != 0 {
		return ret, false // Cell already taken
	}

	var newBoard [9]int
	for i := 0; i < 9; i++ {
		newBoard[i] = brd[i]
	}

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
	return b.CalcWin(b.board)
}

func (b *Board) CalcWin(board [9]int) int {
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

	for _, pattern := range winPatterns {
		a, bIdx, c := pattern[0], pattern[1], pattern[2]
		if board[a] != 0 && board[a] == board[bIdx] && board[bIdx] == board[c] {
			return board[a]
		}
	}

	// Check for draw (no zeros left)
	for i := 0; i < 9; i++ {
		if board[i] == 0 {
			return 0 // Game continues, empty spots left
		}
	}

	return 3 // Draw
}

// Print prints the board in a human-readable format.
func (b *Board) Print() {
	brd := b.Get()
	fmt.Println("Current Board:")
	b.PrintBoard(brd)
	fmt.Println("Next Player:", b.player)
}

func (b *Board) PrintBoard(brd [9]int) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			idx := i*3 + j
			switch brd[idx] {
			case 0:
				fmt.Print(" . ")
			case 1:
				fmt.Print(" X ")
			case 2:
				fmt.Print(" O ")
			}
			if j < 2 {
				fmt.Print("|")
			}
		}
		fmt.Println()
		if i < 2 {
			fmt.Println("-----------")
		}
	}
}
