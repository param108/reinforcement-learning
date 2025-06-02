package main

import (
	"fmt"
	"os"

	"github.com/param108/reinforcement-learning/tictactoe2/game"
	"github.com/param108/reinforcement-learning/tictactoe2/player"
)

func main() {

	if os.Args[1] == "train" {
		// Create 2 human players
		player1 := player.NewLearnerPlayer(1, 0.2, 0.1, "learner")
		wins := 0
		lose := 0
		draw := 0
		for i := 0; i < 10000; i++ {
			fmt.Print("\r", "Playing as X", i)
			player2 := player.NewMinimaxPlayer(2)
			// Create a new game with player 1 starting
			g := game.NewGame(1, player1, player2, true)
			// Play the game
			result := g.Play()
			if result == player1.GetPlayer() {
				wins++
			} else if result == 3 {
				draw++
			} else {
				lose++
			}
		}
		fmt.Println("\nTraining as X finished. Wins:", wins, "Draws:", draw, "Losses:", lose)

		wins = 0
		lose = 0
		draw = 0

		player1.SetPlayer(2)
		for i := 0; i < 10000; i++ {
			fmt.Print("\r", "Playing as O", i)
			player2 := player.NewMinimaxPlayer(1)
			// Create a new game with player 1 starting
			g := game.NewGame(1, player2, player1, true)
			// Play the game
			result := g.Play()
			if result == player1.GetPlayer() {
				wins++
			} else if result == 3 {
				draw++
			} else {
				lose++
			}
		}
		fmt.Println("\nTraining as O finished. Wins:", wins, "Draws:", draw, "Losses:", lose)

		player1.SetPlayer(1)

		wins = 0
		lose = 0
		draw = 0
		for i := 0; i < 10000; i++ {
			fmt.Print("\r", "Playing as X", i)
			player2 := player.NewLearnerPlayer(2, 0.2, 0.1, "learner")
			// Create a new game with player 1 starting
			g := game.NewGame(1, player1, player2, true)
			// Play the game
			result := g.Play()
			if result == player1.GetPlayer() {
				wins++
			} else if result == 3 {
				draw++
			} else {
				lose++
			}
		}
		fmt.Println("\nTraining as X finished. Wins:", wins, "Draws:", draw, "Losses:", lose)

		wins = 0
		lose = 0
		draw = 0

		player1.SetPlayer(2)
		for i := 0; i < 10000; i++ {
			fmt.Print("\r", "Playing as O", i)
			player2 := player.NewLearnerPlayer(1, 0.2, 0.1, "learner")
			// Create a new game with player 1 starting
			g := game.NewGame(1, player2, player1, true)
			// Play the game
			result := g.Play()
			if result == player1.GetPlayer() {
				wins++
			} else if result == 3 {
				draw++
			} else {
				lose++
			}
		}
		fmt.Println("\nTraining as O finished. Wins:", wins, "Draws:", draw, "Losses:", lose)

		player1.SaveModel("learner_player.json")
		return
	}

	if os.Args[1] == "playX" {
		// Create a human player
		player1 := player.NewHumanPlayer(1)
		player2 := player.NewLearnerPlayer(2, 0.2, 0.01, "player")
		if err := player2.LoadModel("learner_player.json"); err != nil {
			fmt.Println("Error loading model:", err)
			return
		}

		g := game.NewGame(1, player1, player2, false)
		result := g.Play()

		g.GetBoard().Print()
		fmt.Println("Game result:", result)
		return
	}

	if os.Args[1] == "playO" {
		// Create a human player
		player2 := player.NewHumanPlayer(2)
		player1 := player.NewLearnerPlayer(1, 0.2, 0.01, "player")
		if err := player1.LoadModel("learner_player.json"); err != nil {
			fmt.Println("Error loading model:", err)
			return
		}

		g := game.NewGame(1, player1, player2, false)
		result := g.Play()

		g.GetBoard().Print()
		fmt.Println("Game result:", result)
		return
	}

	fmt.Println("Invalid command. Use 'train', 'playX', or 'playO'.")
}
