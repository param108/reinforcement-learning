package main

import (
	"math/rand"
)

type ProceduralPlayer interface {
	ChooseMove(board [9]int) Action
}

type Procedural struct {
	player      int
	otherPlayer int
}

func NewProcedural(player int) *Procedural {
	return &Procedural{
		player:      player,
		otherPlayer: 3 - player, // Assuming players are 1 and 2
	}
}

func (p *Procedural) isWinningMove(board [9]int, action Action, player int) bool {
	newBoard, _ := placeMove(board, player, action.X, action.Y)
	if hasWon(player, newBoard) {
		return true
	}
	return false
}

func (p *Procedural) getAvailableCentreOrCorners(board [9]int) []Action {
	availableActions := []Action{}
	centre := Action{X: 1, Y: 1, Player: p.player} // Center position
	corners := []Action{
		{X: 0, Y: 0, Player: p.player}, {X: 0, Y: 2, Player: p.player},
		{X: 2, Y: 0, Player: p.player}, {X: 2, Y: 2, Player: p.player},
	}

	if board[centre.Y*3+centre.X] == 0 {
		availableActions = append(availableActions, centre)
	}

	for _, corner := range corners {
		if board[corner.Y*3+corner.X] == 0 {
			availableActions = append(availableActions, corner)
		}
	}

	return availableActions
}

func (p *Procedural) ChooseMove(board [9]int) Action {
	generatedActions := generateActions(board, p.player)

	// if there is a winning move, take it
	for _, action := range generatedActions {
		if p.isWinningMove(board, action, p.player) {
			return action
		}
	}

	// if opponent has a winning move, block it
	for _, action := range generatedActions {
		if p.isWinningMove(board, action, p.otherPlayer) {
			return action
		}
	}

	// if the center or a corner is available, take it
	cornerAction := p.getAvailableCentreOrCorners(board)
	if len(cornerAction) > 0 {
		return cornerAction[0]
	}

	// choose a random available action
	return generatedActions[rand.Intn(len(generatedActions))]
}
