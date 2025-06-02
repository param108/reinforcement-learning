package player

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"

	"github.com/param108/reinforcement-learning/tictactoe2/board"
)

type LearnerPlayer struct {
	player       int               // 1 or 2
	model        map[int64]float64 // Model to store Q-values
	epsilon      float64           // Exploration rate
	history      []int64           // History of moves for training
	learningRate float64           // Learning rate for Q-learning
	mode         string            // Mode of the player (e.g., "learner", "explorer")
}

func NewLearnerPlayer(player int, epsilon float64, learningRate float64, mode string) *LearnerPlayer {
	return &LearnerPlayer{
		player:       player,
		epsilon:      epsilon,
		model:        make(map[int64]float64),
		history:      []int64{},
		learningRate: learningRate,
		mode:         mode,
	}
}

func (lp *LearnerPlayer) GetPlayer() int {
	return lp.player
}

func (lp *LearnerPlayer) SetPlayer(player int) {
	lp.player = player
}

func (lp *LearnerPlayer) LoadModel(path string) error {
	m := map[string]float64{}

	// Load the model from the specified path
	fp, err := os.Open(path)
	if err != nil {
		return err
	}

	defer fp.Close()
	data, err := io.ReadAll(fp)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	for k, v := range m {
		id, err := strconv.Atoi(k)
		if err != nil {
			return errors.New("invalid key in model: " + k)
		}

		lp.model[int64(id)] = v
	}
	return nil
}

// SaveModel saves the model to the specified path
func (lp *LearnerPlayer) SaveModel(path string) error {
	m := make(map[string]float64)

	// Convert the model to a map with string keys
	for k, v := range lp.model {
		m[strconv.FormatInt(k, 10)] = v
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	if _, err := fp.Write(data); err != nil {
		return err
	}

	return nil
}

// AddHistoryEntry adds a new entry to the player's history
func (lp *LearnerPlayer) AddHistoryEntry(b *board.Board, x, y int) {
	// Create a new history entry with the current board state
	newBoard, _ := b.TryMove(b.Get(), x, y, lp.player)

	id := b.CalcID(newBoard, b.GetStart(), lp.player)
	lp.history = append(lp.history, id)
}

// MakeMove is a placeholder for the learner player logic
func (lp *LearnerPlayer) MakeMove(b *board.Board) (int, int, int) {
	maxValue := float64(-1) // Initialize to a very low value
	maxActions := []board.Action{}

	actions := b.GetPossibleMoves()
	if lp.mode == "learner" {
		if randFloat := rand.Float64(); randFloat < lp.epsilon {
			// Explore: choose a random action
			action := actions[rand.Intn(len(actions))]
			lp.AddHistoryEntry(b, action.X, action.Y)
			return action.X, action.Y, lp.player
		}
	}

	startBoard := b.Get()
	for _, action := range actions {
		newBoard, _ := b.TryMove(startBoard, action.X, action.Y, lp.player)
		id := b.CalcID(newBoard, b.GetStart(), lp.player)
		if _, exists := lp.model[id]; !exists {
			// if the board state is a winning state, assign value 1.0,
			// if it's a losing state, assign value 0,
			// otherwise initialize to 0.5
			if b.CalcWin(newBoard) == lp.player {
				lp.model[id] = 1 // Initialize Q-value if not present
			} else {
				lp.model[id] = 0.5 // Initialize Q-value if not present
			}
		}

		if lp.mode != "learner" {
			fmt.Println("Action:", action.X, action.Y, "ID:", id, "Value:", lp.model[id])
		}

		value := lp.model[id]
		if value > maxValue {
			maxValue = value
			maxActions = []board.Action{action}
		} else if value >= maxValue-0.0001 && value <= maxValue+0.0001 {
			maxActions = append(maxActions, action)
		}
	}

	action := maxActions[rand.Intn(len(maxActions))]
	lp.AddHistoryEntry(b, action.X, action.Y)

	return action.X, action.Y, lp.player
}

func (lp *LearnerPlayer) Win() {
	if lp.mode == "learner" {
		// Update the model based on the history of moves
		for i := len(lp.history) - 2; i >= 0; i-- {
			id := lp.history[i]
			nextID := lp.history[i+1]
			if _, exists := lp.model[id]; !exists {
				lp.model[id] = 0.5 // Increase Q-value for winning moves
			}

			// old := lp.model[id]
			lp.model[id] += lp.learningRate * (lp.model[nextID] - lp.model[id]) // Update Q-value
			// fmt.Println("Updating ID:", id, "Old:", old, "Value:", lp.model[id])

			if lp.model[id] < 0 {
				lp.model[id] = 0 // Ensure Q-value does not go below 0
			}

			if lp.model[id] > 1 {
				lp.model[id] = 1 // Ensure Q-value does not exceed 1
			}
		}
	}
	lp.history = []int64{} // Clear history after updating
}

func (lp *LearnerPlayer) Lose() {
	if lp.mode == "learner" {
		nextValue := float64(0.0) // Losing state has a value of 0

		// Update the model based on the history of moves
		for i := len(lp.history) - 1; i >= 0; i-- {
			id := lp.history[i]

			if i != len(lp.history)-1 {
				nextID := lp.history[i+1]
				nextValue = lp.model[nextID]
			}

			if _, exists := lp.model[id]; !exists {
				lp.model[id] = 0.5 // Decrease Q-value for losing moves
			}

			// old := lp.model[id]
			lp.model[id] += lp.learningRate * (nextValue - lp.model[id]) // Update Q-value
			//fmt.Println("Updating ID:", id, "Old:", old, "Value:", lp.model[id])
			if lp.model[id] < 0 {
				lp.model[id] = 0 // Ensure Q-value does not go below 0
			}

			if lp.model[id] > 1 {
				lp.model[id] = 1 // Ensure Q-value does not exceed 1
			}
		}
	}
	lp.history = []int64{} // Clear history after updating
}
