package player

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/param108/reinforcement-learning/tictactoe2/board"
)

type HumanPlayer struct {
	player int // 1 or 2
}

func NewHumanPlayer(player int) *HumanPlayer {
	return &HumanPlayer{
		player: player,
	}
}

// MakeMove asks the user for coordinates and validates the move
func (hp *HumanPlayer) MakeMove(b *board.Board) (int, int, int) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Player %d, enter your move as \"x y\": ", hp.player)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		parts := strings.Split(text, " ")
		if len(parts) != 2 {
			fmt.Println("Please enter two numbers separated by a space.")
			continue
		}
		x, err1 := strconv.Atoi(parts[0])
		y, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil || x < 0 || x > 2 || y < 0 || y > 2 {
			fmt.Println("Invalid coordinates. Enter x and y between 0 and 2.")
			continue
		}
		_, ok := b.TryMove(b.Get(), x, y, hp.player)
		if !ok {
			fmt.Println("Invalid move. Cell already taken. Try again.")
			continue
		}
		return x, y, hp.player
	}
}

func (hp *HumanPlayer) Win() {
}

func (hp *HumanPlayer) Lose() {
}

func (hp *HumanPlayer) GetPlayer() int {
	return hp.player
}

func (hp *HumanPlayer) SetPlayer(player int) {
	hp.player = player
}
