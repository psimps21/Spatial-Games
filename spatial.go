package main

import (
	"bufio"
	"fmt"
	"gifhelper"
	"image"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
)

// The data stored in a single cell of a field
type Cell struct {
	strategy string  //represents "C" or "D" corresponding to the type of prisoner in the cell
	score    float64 //represents the score of the cell based on the prisoner's relationship with neighboring cells
}

// The game board is a 2D slice of Cell objects
type GameBoard [][]Cell

// InitializeGameBoard initializes a game board with the given number of rows and columns
func InitializeGameBoard(numRows, numCols int) GameBoard {
	var board GameBoard
	board = make(GameBoard, numRows)
	for i := range board {
		board[i] = make([]Cell, numCols)
	}
	return board
}

// CopyBoard produce a copy of a given board
func CopyBoard(board GameBoard) GameBoard {
	// If board is empty
	if len(board) == 0 {
		panic("Empty board given")
	}

	// var newBoard GameBoard
	newBoard := InitializeGameBoard(len(board), len(board[0]))
	for r := range board {
		for c := range board[0] {
			// Copy cell from old board
			newBoard[r][c] = board[r][c]
		}
	}
	return newBoard
}

// GenerateBoardFromFile generates a board from a given file
func GenerateBoardFromFile(filename string) GameBoard {
	// Open board file
	boardFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: Problem opening the given file")
	}
	defer boardFile.Close()

	// Loop over lines of board file
	scanner := bufio.NewScanner(boardFile)
	var board GameBoard
	var boardRow int
	firstLine := true
	for scanner.Scan() {
		// if line contains board coordinates
		if firstLine {
			items := strings.Split(scanner.Text(), " ")
			// Convert board coordinates to integers
			numRows, err1 := strconv.Atoi(items[0])
			numCols, err2 := strconv.Atoi(items[1])
			if err1 != nil {
				fmt.Println(err1)
			}
			if err2 != nil {
				fmt.Println(err1)
			}

			// Initialize game board of file specified size
			board = InitializeGameBoard(numRows, numCols)
			// initialize variable to tracking the board row in the file lines
			boardRow = 0
			firstLine = false
		} else { // Else line contains board row configs
			line := scanner.Text()
			// loop over characters in board row line
			for i := 0; i < len(line); i++ {
				// Initialize Cell for each column in the board row
				board[boardRow][i] = Cell{strategy: string(line[i])}
			}
			// Increment board row for next file line
			boardRow++
		}
	}
	err3 := scanner.Err()
	if err3 != nil {
		panic(err3)
	}
	return board
}

// PrintBoardStrategies print the strategy of cells in a given board
func PrintBoardStrategies(board GameBoard) {
	// Initialize 2d string slice for cell strategies
	var boardStrings [][]string
	boardStrings = make([][]string, len(board))
	for i := range board {
		boardStrings[i] = make([]string, len(board[i]))
	}

	// Populate slice with strategies of each cell
	for r := range board {
		for c := range board[r] {
			boardStrings[r][c] = board[r][c].strategy
		}
	}

	// Print the strategy for each row of the board
	for r := range board {
		fmt.Println(strings.Join(boardStrings[r], ""))
	}
}

// PrintBoardScores print the strategy of cells in a given board
func PrintBoardScores(board GameBoard) {
	// Initialize 2d string slice for cell strategies
	var boardStrings [][]string
	boardStrings = make([][]string, len(board))
	for i := range board {
		boardStrings[i] = make([]string, len(board[i]))
	}

	// Populate slice with strategies of each cell
	for r := range board {
		for c := range board[r] {
			boardStrings[r][c] = fmt.Sprintf("%f", board[r][c].score)
		}
	}

	// Print the strategy for each row of the board
	for r := range board {
		fmt.Println(strings.Join(boardStrings[r], " "))
	}
}

// TopNeighbor retruns the coordinates top neighbor Coords
func TopNeighbor(r, c int) []int {
	return []int{r - 1, c}
}

// BottomNeighbor returns coordinates of the Bottom neighbor
func BottomNeighbor(r, c int) []int {
	return []int{r + 1, c}
}

// LeftNeighbor returns coordinates of the Left neighbor
func LeftNeighbor(r, c int) []int {
	return []int{r, c - 1}
}

// RightNeighbor returns coordinates of the Right neighbor
func RightNeighbor(r, c int) []int {
	return []int{r, c + 1}
}

// TopLeftNeighbor returns coordinates of the Top Left neighbor
func TopLeftNeighbor(r, c int) []int {
	return []int{r - 1, c - 1}
}

// TopRightNeighbor returns coordinates of the Top Right neighbor
func TopRightNeighbor(r, c int) []int {
	return []int{r - 1, c + 1}
}

//  returns coordinates of the Bottom Left neighbor
func BottomLeftNeighbor(r, c int) []int {
	return []int{r + 1, c - 1}
}

// BottomLeftNeighbor returns coordinates of the Bottom Right neighbor
func BottomRightNeighbor(r, c int) []int {
	return []int{r + 1, c + 1}
}

// BottomNeighbors returns coordinates for all neighbors below a cell
func BottomNeighbors(r, c int) [][]int {
	return [][]int{BottomLeftNeighbor(r, c), BottomNeighbor(r, c), BottomRightNeighbor(r, c)}
}

// TopNeighbors returns coordinates for all neighbors above a cell
func TopNeighbors(r, c int) [][]int {
	return [][]int{TopLeftNeighbor(r, c), TopNeighbor(r, c), TopRightNeighbor(r, c)}
}

// LeftNeighbors returns coordinates for all neighbors to the left a cell
func LeftNeighbors(r, c int) [][]int {
	return [][]int{BottomLeftNeighbor(r, c), LeftNeighbor(r, c), TopLeftNeighbor(r, c)}
}

// RightNeighbors returns coordinates for all neighbors to the right a cell
func RightNeighbors(r, c int) [][]int {
	return [][]int{BottomRightNeighbor(r, c), RightNeighbor(r, c), TopRightNeighbor(r, c)}
}

// FindNeighborCoords finds the coordinates of neighbors for a given cell
func FindNeighborCoords(selfR, selfC, numRows, numCols int) [][]int {
	if numRows == 1 && numCols == 1 {
		fmt.Println("No neighbors to play against. Ending Game.")
		os.Exit(1)
	}

	// Coordinate slice to be returned
	var neighborCoords [][]int

	// Conditions for center or border cells
	if selfR > 0 && selfR < (numRows-1) && selfC > 0 && selfC < (numCols-1) { // If center cell
		neighborCoords = append(neighborCoords, TopNeighbor(selfR, selfC), BottomNeighbor(selfR, selfC))
		neighborCoords = append(neighborCoords, LeftNeighbors(selfR, selfC)...)
		neighborCoords = append(neighborCoords, RightNeighbors(selfR, selfC)...)
	} else if selfR == 0 && selfC > 0 && selfC < (numCols-1) { // top border
		neighborCoords = append(neighborCoords, LeftNeighbor(selfR, selfC), RightNeighbor(selfR, selfC))
		neighborCoords = append(neighborCoords, BottomNeighbors(selfR, selfC)...)
	} else if selfR == (numRows-1) && selfC > 0 && selfC < (numCols-1) { // bottom border
		neighborCoords = append(neighborCoords, LeftNeighbor(selfR, selfC), RightNeighbor(selfR, selfC))
		neighborCoords = append(neighborCoords, TopNeighbors(selfR, selfC)...)
	} else if selfR > 0 && selfR < (numRows-1) && selfC == 0 { // left border
		neighborCoords = append(neighborCoords, TopNeighbor(selfR, selfC), BottomNeighbor(selfR, selfC))
		neighborCoords = append(neighborCoords, RightNeighbors(selfR, selfC)...)
	} else if selfR > 0 && selfR < (numRows-1) && selfC == (numCols-1) { // right border
		neighborCoords = append(neighborCoords, TopNeighbor(selfR, selfC), BottomNeighbor(selfR, selfC))
		neighborCoords = append(neighborCoords, LeftNeighbors(selfR, selfC)...)
	} else if selfR == 0 && selfC == 0 { // top left corner
		neighborCoords = append(neighborCoords, RightNeighbor(selfR, selfC), BottomNeighbor(selfR, selfC), BottomRightNeighbor(selfR, selfC))
	} else if selfR == 0 && selfC == (numCols-1) { // top right corner
		neighborCoords = append(neighborCoords, LeftNeighbor(selfR, selfC), BottomNeighbor(selfR, selfC), BottomLeftNeighbor(selfR, selfC))
	} else if selfR == (numRows-1) && selfC == 0 { // bottom left corner
		neighborCoords = append(neighborCoords, RightNeighbor(selfR, selfC), TopNeighbor(selfR, selfC), TopRightNeighbor(selfR, selfC))
	} else if selfR == (numRows-1) && selfC == (numCols-1) { // bottom right corner
		neighborCoords = append(neighborCoords, LeftNeighbor(selfR, selfC), TopNeighbor(selfR, selfC), TopLeftNeighbor(selfR, selfC))
	} else {
		fmt.Printf("Coordinate out of range (%d,%d)", selfR, selfC)
		os.Exit(1)
	}
	return neighborCoords
}

// RunPrisonersDilemna runs prisoners dilemna between two given cells
func RunPrisonersDilemna(board GameBoard, b float64, selfR, selfC, oppR, oppC int) GameBoard {
	// Get self cell
	selfCell := board[selfR][selfC]
	// Get neighbor Cell
	oppCell := board[oppR][oppC]

	// If both self and opponent cooperate
	if selfCell.strategy == "C" && oppCell.strategy == "C" {
		board[selfR][selfC].score += 1.0
		board[oppR][oppC].score += 1.0
	} else if selfCell.strategy == "C" && oppCell.strategy == "D" {
		board[oppR][oppC].score += b
	} else if selfCell.strategy == "D" && oppCell.strategy == "C" {
		board[selfR][selfC].score += b
	} else if selfCell.strategy == "D" && oppCell.strategy == "D" {
		return board
	}
	return board
}

// PlayNeighbors plays the Prisoner's dilemna against all neighbors of a given cell
func PlayNeighbors(board GameBoard, neighborCoords [][]int, selfR, selfC int, b float64) GameBoard {
	// Loop through neighboring cell coords
	for nbrRow := range neighborCoords {
		// Set neighbor cell coordinates
		oppR := neighborCoords[nbrRow][0]
		oppC := neighborCoords[nbrRow][1]

		// Compete against neighbor Cell
		board = RunPrisonersDilemna(board, b, selfR, selfC, oppR, oppC)
	}
	return board
}

// EvolveBoardOnce evolves a board by one step
func EvolveBoardOnce(board GameBoard, b float64) GameBoard {
	// Define board size
	numRows := len(board)
	numCols := len(board[0])
	// fmt.Printf("Num Rows: %d\nNum cols: %d\n", numRows, numCols)

	// Update scores in each cell of the board
	for r := range board {
		for c := range board[0] {
			neighborCoords := FindNeighborCoords(r, c, numRows, numCols)
			board = PlayNeighbors(board, neighborCoords, r, c, b)
		}
	}

	// Create a new board with updated states from the current round and reset scores
	nextBoard := InitializeGameBoard(numRows, numCols)
	for r := range board {
		for c := range board[0] {
			maxCoords := []int{r, c}
			neighborCoords := FindNeighborCoords(r, c, numRows, numCols)

			for i := range neighborCoords {
				if board[maxCoords[0]][maxCoords[1]].score < board[neighborCoords[i][0]][neighborCoords[i][1]].score {
					maxCoords[0] = neighborCoords[i][0]
					maxCoords[1] = neighborCoords[i][1]
				}
				nextBoard[r][c] = Cell{board[maxCoords[0]][maxCoords[1]].strategy, 0}
			}
		}
	}
	return nextBoard
}

// DrawBoard creates an image from a GameBoard
func DrawBoard(board GameBoard) image.Image {
	cellWidth := 1
	height := len(board) * cellWidth
	width := len(board[0]) * cellWidth

	// declare colors
	blue := MakeColor(0, 0, 255)
	red := MakeColor(255, 0, 0)

	c := CreateNewPalettedCanvas(width, height, nil)
	for i := range board {
		for j := range board[0] {
			if board[i][j].strategy == "C" {
				c.SetFillColor(blue)
			} else if board[i][j].strategy == "D" {
				c.SetFillColor(red)
			} else {
				panic("Cell contains invalid strategy:" + board[i][j].strategy)
			}
			x := i * cellWidth
			y := j * cellWidth
			c.ClearRect(x, y, x+cellWidth, y+cellWidth)
			c.Fill()
		}
	}
	return GetImage(c)
}

// DrawBoards creates an image for a GameBoards in a given list
func DrawBoards(boards []GameBoard) []image.Image {
	numSteps := len(boards)
	imgList := make([]image.Image, numSteps)
	for i := range boards {
		imgList[i] = DrawBoard(boards[i])
	}
	return imgList
}

func main() {
	// Get arguments from command line
	filename := os.Args[1]
	b, err1 := strconv.ParseFloat(os.Args[2], 64)
	if err1 != nil {
		panic("Problem converting score b to a float")
	}
	steps, err2 := strconv.Atoi(os.Args[3])
	if err2 != nil {
		panic("Problem converting number of steps to an integer")
	}

	// Generate board from file
	initialBoard := GenerateBoardFromFile(filename)

	// Array for storing evolved boards
	var allBoards = []GameBoard{initialBoard}

	// Evolve the board
	nextBoard := CopyBoard(initialBoard)
	for step := 1; step <= steps; step++ {
		nextBoard = EvolveBoardOnce(nextBoard, b)
		allBoards = append(allBoards, CopyBoard(nextBoard))
	}

	imgList := DrawBoards(allBoards)
	finalImage := imgList[len(imgList)-1]

	f, err3 := os.Create("Prisoners.png")
	if err3 != nil {
		log.Fatal(err3)
	}

	err4 := png.Encode(f, finalImage)
	if err4 != nil {
		f.Close()
		log.Fatal(err4)
	}

	outputFile := "Prisoners"
	gifhelper.ImagesToGIF(imgList, outputFile)

	fmt.Println("Done")
}
