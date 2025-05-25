package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func hasWon(player int, board [9]int) bool {
	winningCombos := [8][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}
	for _, combo := range winningCombos {
		if board[combo[0]] == player && board[combo[1]] == player && board[combo[2]] == player {
			return true
		}
	}
	return false
}

func movesRemain(board [9]int) bool {
	for _, v := range board {
		if v == 0 {
			return true
		}
	}
	return false
}

func placeMove(board [9]int, player int, x int, y int) ([9]int, bool) {
	// Check that x and y are in the valid range
	if x < 0 || x > 2 || y < 0 || y > 2 {
		return board, false // Out of bounds
	}
	pos := y*3 + x
	if board[pos] != 0 {
		return board, false // Position already occupied
	}
	newBoard := board
	newBoard[pos] = player
	return newBoard, true
}

func boardToInt(board [9]int) int {
	total := 0
	multiplier := 1
	for i := 0; i < 9; i++ {
		total += board[i] * multiplier
		multiplier *= 3
	}
	return total
}

func printBoard(board [9]int) {
	symbols := []string{".", "X", "O"}
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			idx := y*3 + x
			fmt.Print(symbols[board[idx]])
			if x < 2 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

type Action struct {
	X      int
	Y      int
	Player int
}

func generateActions(board [9]int, player int) []Action {
	actions := []Action{}
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			idx := y*3 + x
			if board[idx] == 0 {
				actions = append(actions, Action{X: x, Y: y, Player: player})
			}
		}
	}
	return actions
}

func actionToNumber(a Action) int {
	return a.X*100 + a.Y*10
}

type RLPlayer struct {
	ProbabilityTable map[int]map[int]float64 // map[boardAsInt][actionAsInt]probability
	Player           int                     // 1 for one player, 2 for the other
}

func NewRLPlayer(player int) *RLPlayer {
	return &RLPlayer{
		ProbabilityTable: make(map[int]map[int]float64),
		Player:           player,
	}
}

type BoardAction struct {
	Board  [9]int
	Action Action
}

func learn(winner *RLPlayer, history []BoardAction, win bool) {
	// Step 1: Filter out BoardActions not belonging to the winning player
	playerMoves := []BoardAction{}
	for _, ba := range history {
		if ba.Action.Player == winner.Player {
			playerMoves = append(playerMoves, ba)
		}
	}
	if len(playerMoves) == 0 {
		return
	}

	// Step 2: Handle the last board,action based on the win parameter
	last := playerMoves[len(playerMoves)-1]
	boardInt := boardToInt(last.Board)
	actionInt := actionToNumber(last.Action)
	if winner.ProbabilityTable[boardInt] == nil {
		winner.ProbabilityTable[boardInt] = make(map[int]float64)
	}

	if win {
		// If the player wins, set the last probability to 1
		winner.ProbabilityTable[boardInt][actionInt] = 1.0
	} else {
		// If the player loses, set the last probability to 0
		winner.ProbabilityTable[boardInt][actionInt] = 0.0
	}

	// Step 3: Back-propagate probabilities
	// Move backward from penultimate to first
	for i := len(playerMoves) - 2; i >= 0; i-- {
		cur := playerMoves[i]
		next := playerMoves[i+1]
		curBoardInt := boardToInt(cur.Board)
		curActionInt := actionToNumber(cur.Action)
		nextBoardInt := boardToInt(next.Board)
		nextActionInt := actionToNumber(next.Action)

		if winner.ProbabilityTable[curBoardInt] == nil {
			winner.ProbabilityTable[curBoardInt] = make(map[int]float64)
		}
		// Get current and next probability
		curProb := winner.ProbabilityTable[curBoardInt][curActionInt]
		nextProb := winner.ProbabilityTable[nextBoardInt][nextActionInt]
		// Default to 0.5 if never seen
		if curProb == 0 {
			curProb = 0.5
		}

		prob := 0.2 * (nextProb - curProb)
		winner.ProbabilityTable[curBoardInt][curActionInt] += prob
	}
}

func encodeProbabilityTable(table map[int]map[int]float64) map[string]map[string]float64 {
	out := make(map[string]map[string]float64)
	for boardInt, actions := range table {
		boardKey := fmt.Sprintf("%d", boardInt)
		out[boardKey] = make(map[string]float64)
		for actionInt, prob := range actions {
			actionKey := fmt.Sprintf("%d", actionInt)
			out[boardKey][actionKey] = prob
		}
	}
	return out
}

func writeProbabilityTableFile(table map[int]map[int]float64, filename string) error {
	converted := encodeProbabilityTable(table)
	data, err := json.MarshalIndent(converted, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func decodeProbabilityTable(data map[string]map[string]float64) (map[int]map[int]float64, error) {
	result := make(map[int]map[int]float64)
	for boardStr, actionMap := range data {
		boardInt, err := strconv.Atoi(boardStr)
		if err != nil {
			return nil, fmt.Errorf("invalid board key: %s", boardStr)
		}
		result[boardInt] = make(map[int]float64)
		for actionStr, prob := range actionMap {
			actionInt, err := strconv.Atoi(actionStr)
			if err != nil {
				return nil, fmt.Errorf("invalid action key: %s", actionStr)
			}
			result[boardInt][actionInt] = prob
		}
	}
	return result, nil
}

func readProbabilityTableFile(filename string) (map[int]map[int]float64, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var intermediate map[string]map[string]float64
	if err := json.Unmarshal(data, &intermediate); err != nil {
		return nil, err
	}

	return decodeProbabilityTable(intermediate)
}

func playRound(
	p1, p2 *RLPlayer,
	humanPlayerNum int, isHuman bool) (*RLPlayer, []BoardAction, bool, error) {
	players := []*RLPlayer{p1, p2}
	history := []BoardAction{}
	board := [9]int{}

	curIdx := rand.Intn(2)
	otherIdx := 1 - curIdx

	reader := bufio.NewReader(os.Stdin)

	for {
		player := players[curIdx]
		if isHuman && player.Player == humanPlayerNum {
			// Human's move
			printBoard(board)
			for {
				fmt.Printf("Your move (enter x y, 0-based): ")
			AGAIN:
				line, err := reader.ReadString('\n')
				if err != nil {
					fmt.Printf("Error reading input: %v\n", err)
					goto AGAIN
				}
				line = strings.TrimSpace(line)
				parts := strings.Fields(line)
				if len(parts) != 2 {
					fmt.Println("Please enter two integers separated by space.")
					continue
				}
				x, err1 := strconv.Atoi(parts[0])
				y, err2 := strconv.Atoi(parts[1])
				if err1 != nil || err2 != nil || x < 0 || x > 2 || y < 0 || y > 2 {
					fmt.Println("Invalid input. x and y must be integers in [0,2].")
					continue
				}
				pos := y*3 + x
				if board[pos] != 0 {
					fmt.Println("Position already occupied. Try again.")
					continue
				}
				board, _ = placeMove(board, humanPlayerNum, x, y)
				history = append(history, BoardAction{
					Board:  board,
					Action: Action{X: x, Y: y, Player: humanPlayerNum},
				})
				break
			}
		} else {
			// AI's move
			possibleActions := generateActions(board, player.Player)
			if len(possibleActions) == 0 {
				return nil, history, true, nil // Draw
			}

			boardInt := boardToInt(board)
			if player.ProbabilityTable[boardInt] == nil {
				player.ProbabilityTable[boardInt] = make(map[int]float64)
			}

			bestProb := -1.0
			var bestActions []Action

			for _, a := range possibleActions {
				actionInt := actionToNumber(a)
				prob, ok := player.ProbabilityTable[boardInt][actionInt]
				if !ok {
					player.ProbabilityTable[boardInt][actionInt] = 0.5
					prob = 0.5
				}
				if prob > bestProb {
					bestProb = prob
					bestActions = []Action{a}
				} else if prob == bestProb {
					bestActions = append(bestActions, a)
				}
			}

			var chosen Action
			if rand.Intn(10) < 8 && len(bestActions) > 0 { // 80% exploit
				chosen = bestActions[rand.Intn(len(bestActions))]
			} else {
				chosen = possibleActions[rand.Intn(len(possibleActions))]
			}

			// Clone board before move
			boardBefore := board
			newBoard, _ := placeMove(board, player.Player, chosen.X, chosen.Y)
			history = append(history, BoardAction{
				Board:  boardBefore,
				Action: chosen,
			})
			board = newBoard
		}

		// Check for win or draw
		if hasWon(player.Player, board) {
			return player, history, false, nil
		}

		if !movesRemain(board) {
			return nil, history, true, nil
		}

		// Switch players
		curIdx, otherIdx = otherIdx, curIdx
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ttt <mode>")
		fmt.Println("Modes: learn, play")
		return
	}
	mode := os.Args[1]

	switch mode {
	case "learn":
		// Two RL players: X=1, O=2
		p1 := NewRLPlayer(1)
		p2 := NewRLPlayer(2)
		totalGames := 20000000 // Or pick your number
		wins := 0
		losses := 0
		draws := 0

		for i := 0; i < totalGames; i++ {
			winner, history, _, _ := playRound(p1, p2, 0, false)
			if winner != nil {
				if winner == p1 {
					learn(p1, history, true)
					learn(p2, history, false)
					wins++
				} else {
					learn(p1, history, false)
					learn(p2, history, true)
					losses++
				}
			} else {
				draws++
			}
			// Optionally print progress, etc.
			if i%50 == 0 {
				fmt.Printf("Game %d wins %d losses %d draws %d\n", i, wins, losses, draws)
			}
		}

		// Write p1's table as p1.json, overwrite
		err := writeProbabilityTableFile(p1.ProbabilityTable, "p1.json")
		if err != nil {
			fmt.Println("Error writing p1.json:", err)
		} else {
			fmt.Println("Learning complete. p1.json written.")
		}

	case "play":
		// Randomly assign AI & Human to X (1) and O (2)
		aiPlayerNum := 1 + rand.Intn(2) // 1 or 2
		humanPlayerNum := 1
		if aiPlayerNum == 1 {
			humanPlayerNum = 2
		}

		p1 := NewRLPlayer(aiPlayerNum)
		p2 := NewRLPlayer(humanPlayerNum)
		p1.ProbabilityTable, _ = readProbabilityTableFile("p1.json")

		fmt.Printf("AI is player %d, Human is player %d\n", aiPlayerNum, humanPlayerNum)
		player, history, draw, err := playRound(p1, p2, humanPlayerNum, true)
		if err != nil {
			fmt.Println("Error during game:", err)
		}

		printBoard(history[len(history)-1].Board)
		if draw {
			fmt.Println("It's a draw!")
			return
		}

		if player == p1 {
			fmt.Printf("AI Player %d wins!\n", player.Player)
		} else {
			fmt.Printf("Human Player %d wins!\n", humanPlayerNum)
		}
	default:
		fmt.Println("Unknown mode. Use 'learn' or 'play'.")
	}
}
