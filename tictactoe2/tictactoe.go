package main

import (
	"fmt"

	"github.com/param108/reinforcement-learning/tictactoe2/game"
	"github.com/param108/reinforcement-learning/tictactoe2/player"
)

func main() {
	// Create 2 human players
	player1 := player.NewHumanPlayer(1)
	player2 := player.NewMinimaxPlayer(2)
	// Create a new game with player 1 starting
	game := game.NewGame(1, player1, player2)
	// Play the game
	result := game.Play()
	// Print the final board state
	game.GetBoard().Print()
	// Print the result
	switch result {
	case 1:
		fmt.Println("Player 1 (X) wins!")
	case 2:
		fmt.Println("Player 2 (O) wins!")
	case 3:
		fmt.Println("It's a draw!")
	default:
		fmt.Println("Game continues or invalid state.")
	}
}
