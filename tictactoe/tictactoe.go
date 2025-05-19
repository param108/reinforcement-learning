package main

import (
	"bufio"
	"encoding/json"
	"errors"
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

func playRound(p1, p2 *RLPlayer) (*RLPlayer, []BoardAction, bool) {
	players := []*RLPlayer{p1, p2}
	history := []BoardAction{}
	board := [9]int{}

	curIdx := rand.Intn(2)
	otherIdx := 1 - curIdx

	for {
		player := players[curIdx]
		possibleActions := generateActions(board, player.Player)

		if len(possibleActions) == 0 {
			return nil, history, true // Draw
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

		if hasWon(player.Player, board) {
			return player, history, false
		}

		if !movesRemain(board) {
			return nil, history, true
		}

		curIdx, otherIdx = otherIdx, curIdx
	}
}

func learn(winner *RLPlayer, history []BoardAction) {
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

	// Step 2: Set the last board,action for the winner to probability 1
	last := playerMoves[len(playerMoves)-1]
	boardInt := boardToInt(last.Board)
	actionInt := actionToNumber(last.Action)
	if winner.ProbabilityTable[boardInt] == nil {
		winner.ProbabilityTable[boardInt] = make(map[int]float64)
	}
	winner.ProbabilityTable[boardInt][actionInt] = 1.0

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
		if nextProb == 0 {
			nextProb = 0.5
		}
		// Update
		winner.ProbabilityTable[curBoardInt][curActionInt] = curProb + 0.2*(nextProb-curProb)
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

func playRLPlayerVsHuman(p1file string, rlPlayerNum, humanPlayerNum int) error {
	table, err := readProbabilityTableFile(p1file)
	if err != nil {
		return err
	}
	ai := RLPlayer{
		ProbabilityTable: table,
		Player:           rlPlayerNum, // e.g. 1 (X)
	}

	human := humanPlayerNum // e.g. 2 (O)
	board := [9]int{}

	fmt.Printf("You are player %d\n", human)
	curPlayer := 1 // always start with 1 (X), change as desired

	reader := bufio.NewReader(os.Stdin)
	for {
		printBoard(board)
		var action Action
		if curPlayer == ai.Player {
			// AI's move
			actions := generateActions(board, ai.Player)
			if len(actions) == 0 {
				fmt.Println("Draw!")
				break
			}
			boardInt := boardToInt(board)
			bestProb := -1.0
			var bestActions []Action
			for _, a := range actions {
				actionInt := actionToNumber(a)
				prob, ok := ai.ProbabilityTable[boardInt][actionInt]
				if !ok {
					prob = 0.5
				}
				if prob > bestProb {
					bestProb = prob
					bestActions = []Action{a}
				} else if prob == bestProb {
					bestActions = append(bestActions, a)
				}
			}
			action = bestActions[rand.Intn(len(bestActions))]
			fmt.Printf("AI (Player %d) places at %d %d\n", ai.Player, action.X, action.Y)
			board, _ = placeMove(board, ai.Player, action.X, action.Y)
			if hasWon(ai.Player, board) {
				printBoard(board)
				fmt.Printf("AI (Player %d) wins!\n", ai.Player)
				break
			}
		} else if curPlayer == human {
			// Human's move
			for {
				fmt.Printf("Your move (enter x y, 0-based): ")
				line, err := reader.ReadString('\n')
				if err != nil {
					return err
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
				board, _ = placeMove(board, human, x, y)
				if hasWon(human, board) {
					printBoard(board)
					fmt.Println("You win!")
					return nil
				}
				break
			}
		} else {
			return errors.New("unexpected player turn")
		}

		if !movesRemain(board) {
			printBoard(board)
			fmt.Println("Game is a draw!")
			break
		}
		// Switch
		if curPlayer == ai.Player {
			curPlayer = human
		} else {
			curPlayer = ai.Player
		}
	}

	return nil
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
		totalGames := 10000000 // Or pick your number

		for i := 0; i < totalGames; i++ {
			winner, history, _ := playRound(p1, p2)
			if winner != nil {
				learn(winner, history)
			}
			// Optionally print progress, etc.
			if i%50 == 0 {
				fmt.Printf("Game %d\n", i)
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
		fmt.Printf("AI is player %d, Human is player %d\n", aiPlayerNum, humanPlayerNum)
		if err := playRLPlayerVsHuman("p1.json", aiPlayerNum, humanPlayerNum); err != nil {
			fmt.Println("Error:", err)
		}

	default:
		fmt.Println("Unknown mode. Use 'learn' or 'play'.")
	}
}
