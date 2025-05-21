package main

import (
	"testing"
)

// Test hasWon function
func TestHasWon(t *testing.T) {
	tests := []struct {
		player int
		board  [9]int
		want   bool
	}{
		{1, [9]int{1, 1, 1, 0, 0, 0, 0, 0, 0}, true},  // Row win
		{2, [9]int{0, 0, 0, 2, 2, 2, 0, 0, 0}, true},  // Row win
		{1, [9]int{1, 0, 0, 1, 0, 0, 1, 0, 0}, true},  // Column win
		{2, [9]int{0, 0, 2, 0, 0, 2, 0, 0, 2}, true},  // Column win
		{1, [9]int{1, 0, 0, 0, 1, 0, 0, 0, 1}, true},  // Diagonal win
		{2, [9]int{0, 0, 2, 0, 2, 0, 2, 0, 0}, true},  // Diagonal win
		{1, [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}, false}, // No win
	}
	for _, tt := range tests {
		got := hasWon(tt.player, tt.board)
		if got != tt.want {
			t.Errorf("hasWon(%d, %v) = %v; want %v", tt.player, tt.board, got, tt.want)
		}
	}
}

// Test movesRemain function
func TestMovesRemain(t *testing.T) {
	tests := []struct {
		board [9]int
		want  bool
	}{
		{[9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}, true},  // All empty
		{[9]int{1, 2, 1, 2, 1, 2, 1, 2, 1}, false}, // Full board
		{[9]int{1, 0, 1, 2, 1, 2, 1, 2, 1}, true},  // One empty
	}
	for _, tt := range tests {
		got := movesRemain(tt.board)
		if got != tt.want {
			t.Errorf("movesRemain(%v) = %v; want %v", tt.board, got, tt.want)
		}
	}
}

// Test placeMove function
func TestPlaceMove(t *testing.T) {
	tests := []struct {
		board  [9]int
		player int
		x, y   int
		want   [9]int
		valid  bool
	}{
		{[9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}, 1, 0, 0, [9]int{1, 0, 0, 0, 0, 0, 0, 0, 0}, true},  // Valid move
		{[9]int{1, 0, 0, 0, 0, 0, 0, 0, 0}, 2, 0, 0, [9]int{1, 0, 0, 0, 0, 0, 0, 0, 0}, false}, // Occupied
		{[9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}, 1, 3, 3, [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}, false}, // Out of bounds
	}
	for _, tt := range tests {
		got, valid := placeMove(tt.board, tt.player, tt.x, tt.y)
		if got != tt.want || valid != tt.valid {
			t.Errorf("placeMove(%v, %d, %d, %d) = (%v, %v); want (%v, %v)", tt.board, tt.player, tt.x, tt.y, got, valid, tt.want, tt.valid)
		}
	}
}

// Test boardToInt function
func TestBoardToInt(t *testing.T) {
	tests := []struct {
		board [9]int
		want  int
	}{
		{[9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}, 0},     // Empty board
		{[9]int{1, 0, 0, 0, 0, 0, 0, 0, 0}, 1},     // Single move
		{[9]int{1, 2, 0, 0, 0, 0, 0, 0, 0}, 7},     // Mixed moves
		{[9]int{1, 2, 1, 2, 1, 2, 1, 2, 1}, 12301}, // Full board
	}
	for _, tt := range tests {
		got := boardToInt(tt.board)
		if got != tt.want {
			t.Errorf("boardToInt(%v) = %d; want %d", tt.board, got, tt.want)
		}
	}
}

// Test generateActions function
func TestGenerateActions(t *testing.T) {
	tests := []struct {
		board  [9]int
		player int
		want   int
	}{
		{[9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}, 1, 9}, // All empty
		{[9]int{1, 2, 1, 2, 1, 2, 1, 2, 1}, 2, 0}, // Full board
		{[9]int{1, 0, 1, 2, 1, 2, 1, 2, 1}, 2, 1}, // One empty
	}
	for _, tt := range tests {
		got := generateActions(tt.board, tt.player)
		if len(got) != tt.want {
			t.Errorf("generateActions(%v, %d) = %d actions; want %d actions", tt.board, tt.player, len(got), tt.want)
		}
	}
}
