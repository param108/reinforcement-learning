package main

import "fmt"

type TreeNode struct {
	board    [9]int
	toPlay   int         // player to play next
	children []*TreeNode // child nodes
	seen     bool        // whether this node has been seen in the tree
}

// generateTree generates a game tree using generateActions function
func generateTree(board [9]int, toPlay int) *TreeNode {
	root := &TreeNode{
		board:    board,
		toPlay:   toPlay,
		children: []*TreeNode{},
	}

	// Generate all possible actions for the current board state
	actions := generateActions(board, toPlay)

	for _, action := range actions {
		newBoard, _ := placeMove(board, toPlay, action.X, action.Y)
		childNode := generateTree(newBoard, 3-toPlay) // Switch player
		root.children = append(root.children, childNode)
	}

	return root
}

func treeToArray(node *TreeNode, output []*TreeNode) {
	if node == nil {
		return
	}
	output = append(output, node)
	for _, child := range node.children {
		treeToArray(child, output)
	}
	return
}

func clearSeen(node *TreeNode) {
	if node == nil {
		return
	}
	node.seen = false
	for _, child := range node.children {
		clearSeen(child)
	}
}

func clearAllSeen() {
	if X != nil {
		clearSeen(X)
	}
	if Y != nil {
		clearSeen(Y)
	}
}

func getNextUnseen(nodes []*TreeNode, lastIndex int, toPlay int) (int, *TreeNode) {
	if lastIndex >= len(nodes)-1 {
		return -1, nil
	}

	for i := lastIndex + 1; i < len(nodes); i++ {
		if !nodes[i].seen && nodes[i].toPlay == toPlay {
			nodes[i].seen = true
			return i, nodes[i]
		}
	}

	return -1, nil
}

func countNodes(node *TreeNode) int {
	if node == nil {
		return 0
	}
	count := 1 // Count this node
	for _, child := range node.children {
		count += countNodes(child)
	}
	return count
}

func iterateTree(start int, toPlay int) *TreeNode {
	if start == 1 {
		lastSeenX = getNextUnseen(X, toPlay)
		return lastSeenX
	} else if start == 2 {
		lastSeenY = getNextUnseen(Y, toPlay)
		return lastSeenY
	}
	return nil
}

var lastSeenX int
var lastSeenY int
var X, Y []*TreeNode

func init() {
	x := generateTree([9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}, 1)
	y := generateTree([9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}, 2)
	treeToArray(x, X)
	treeToArray(y, Y)
	fmt.Println("Game tree initialized.", countNodes(X), "nodes for player 1 and", countNodes(Y), "nodes for player 2.")
}
