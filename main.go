package main

import (
	"fmt"
	"sort"
	"strconv"
)

type Player struct {
	Symbol string
	AI     bool
}

type Game struct {
	Turn    int
	Cells   []string
	Players []Player
	Active  bool
}

var SymbolIDs map[int]string

func NewGame(HumanSymbol string) *Game {
	players := []Player{
		Player{
			Symbol: "X",
			AI:     true,
		},
		Player{
			Symbol: "O",
			AI:     true,
		},
	}
	SymbolIDs = make(map[int]string)
	SymbolIDs[0] = players[0].Symbol
	SymbolIDs[1] = players[1].Symbol

	if HumanSymbol == "X" {
		players[0].AI = false
	} else if HumanSymbol == "O" {
		players[1].AI = false
	}
	return &Game{
		Turn:    0,
		Cells:   []string{"-", "-", "-", "-", "-", "-", "-", "-", "-"},
		Players: players,
		Active:  true,
	}
}

func (g *Game) takePosition(position int, symbol string) {
	g.Cells[position] = symbol
}

func (g *Game) printBoard() {
	for i, pos := range g.Cells {
		switch i % 3 {
		case 0:
			fmt.Printf("%v |", pos)
		case 1:
			fmt.Printf(" %v |", pos)
		case 2:
			fmt.Printf(" %v\n", pos)
		}
	}
	fmt.Println()
}

func (g *Game) isFreeSpace(space int) bool {
	return g.Cells[space] == "-"
}

func symbolWon(board []string, symbol string) bool {
	return ((board[0] == symbol && board[1] == symbol && board[2] == symbol) || //First Horizontal
		(board[3] == symbol && board[4] == symbol && board[5] == symbol) || //Second Horizontal
		(board[6] == symbol && board[7] == symbol && board[8] == symbol) || //Third Horizontal
		(board[0] == symbol && board[3] == symbol && board[6] == symbol) || //First Vertical
		(board[1] == symbol && board[4] == symbol && board[7] == symbol) || //Second Vertical
		(board[2] == symbol && board[5] == symbol && board[8] == symbol) || //Third Vertical
		(board[0] == symbol && board[4] == symbol && board[8] == symbol) || //First Diagonal
		(board[2] == symbol && board[4] == symbol && board[6] == symbol)) //Second Diagoonal
}

//Returns an array of available moves
func (g *Game) getAvailableMoves() []int {
	var arr []int
	for i, cell := range g.Cells {
		if cell == "-" {
			arr = append(arr, i)
		}
	}
	return arr
}

//Returns an array of available Moves
func getAvailableMoves(board []string) []int {
	var arr []int
	for i, cell := range board {
		if cell == "-" {
			arr = append(arr, i)
		}
	}
	return arr
}

//Returns an array of taken spots
func (g *Game) getUnavailableMoves() []int {
	var arr []int
	for i, cell := range g.Cells {
		if cell != "-" {
			arr = append(arr, i)
		}
	}
	return arr
}

func (g *Game) getPlayerMove() int {
	var input string
	//TODO: Print only available slots
	fmt.Printf("Enter move number\n0 | 1 | 2\n3 | 4 | 5\n6 | 7 | 8\n")
	fmt.Scanln(&input)
	i, err := strconv.Atoi(input)
	//Check Move Number
	if err != nil {
		fmt.Println("Entry must be a number")
		return g.getPlayerMove()
	}
	//TODO: Check Number in range
	if i < 0 || i > 8 {
		fmt.Println("Entry must be in range")
		return g.getPlayerMove()
	}

	//Check Available Moves
	availableMoves := g.getAvailableMoves()
	if sort.SearchInts(availableMoves, i) == len(availableMoves) {
		fmt.Println("Entry must be available")
		return g.getPlayerMove()
	}
	return i
}

type move struct {
	index int
	score int
}

//MINIMAX
func miniMax(board []string, symbol, playerSymbol, opponantSymbol string) move {
	i := 0
	availableMoves := getAvailableMoves(board)
	var m move
	var moves []move

	if symbolWon(board, symbol) {
		return move{score: 10}
	} else if len(availableMoves) == 0 {
		return move{score: 0}
	} else {
		if symbol == playerSymbol {
			if symbolWon(board, opponantSymbol) {
				return move{score: -10}
			}
		} else {
			if symbolWon(board, playerSymbol) {
				return move{score: -10}
			}
		}
	}

	for i = 0; i < len(availableMoves); i++ {
		m = move{
			index: availableMoves[i],
		}
		board[m.index] = playerSymbol
		result := miniMax(board, symbol, opponantSymbol, playerSymbol)
		m.score = result.score
		board[m.index] = "-"
		moves = append(moves, m)
	}

	var bestScore int
	var bestIndex int
	if symbol == playerSymbol {
		bestScore = -10000
		for i = 0; i < len(moves); i++ {
			if moves[i].score > bestScore {
				bestScore = moves[i].score
				bestIndex = i
			}
		}
	} else {
		bestScore = 10000
		for i = 0; i < len(moves); i++ {
			if moves[i].score < bestScore {
				bestScore = moves[i].score
				bestIndex = i
			}
		}
	}

	return moves[bestIndex]
}

//END MINIMAX

func (g *Game) getAIMove(playerSymbol, opponantSymbol string) int {
	m := miniMax(g.Cells, playerSymbol, playerSymbol, opponantSymbol)
	return m.index
}

func (g *Game) playerTurn() {
	var i int
	p := g.Players[g.Turn]
	if p.AI {
		i = g.getAIMove(p.Symbol, g.Players[1-g.Turn].Symbol)
	} else {
		i = g.getPlayerMove()
	}
	g.takePosition(i, p.Symbol)
	if symbolWon(g.Cells, p.Symbol) {
		fmt.Printf("Game Over! %v is the winner!\n", p.Symbol)
		g.Active = false
	}
	g.Turn = 1 - g.Turn //Swap between 1 and 0 and 0 and 1
}

func (g *Game) Loop() {
	for g.Active {
		g.printBoard()
		g.playerTurn()
		if len(g.getAvailableMoves()) == 0 {
			g.Active = false
			fmt.Printf("Cat got it!\n")
		}
	}
}

func main() {
	g := NewGame("X")
	g.Loop()
}
