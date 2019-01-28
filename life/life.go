package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type boardSymbol uint8

const alive boardSymbol = '*'
const dead boardSymbol = ' '

type life struct {
	startConfig  [][]boardSymbol
	currentState [][]boardSymbol
	tempState    [][]boardSymbol
	dimX, dimY   int
}

func deep2DCopy(x, y int, board [][]boardSymbol) [][]boardSymbol {
	res := make([][]boardSymbol, x)
	for i := range board {
		res[i] = make([]boardSymbol, y)
		copy(res[i], board[i])
	}
	return res
}

// lines longer than 65536 characters are not supported
func NewLife(path string) *life {
	l := new(life)

	configFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer configFile.Close()

	scanner := bufio.NewScanner(configFile)
	scanner.Scan()
	dimensions := strings.Split(scanner.Text(), " ")
	l.dimX, _ = strconv.Atoi(dimensions[0])
	l.dimY, _ = strconv.Atoi(dimensions[1])

	l.startConfig = make([][]boardSymbol, l.dimX)
	curRow := 0
	for scanner.Scan() {
		row := scanner.Text()
		l.startConfig[curRow] = make([]boardSymbol, l.dimY)
		for curCol, char := range row {
			var write boardSymbol
			if char == '-' {
				write = dead
			} else if char == '*' {
				write = alive
			} else {
				fmt.Println("invalid config")
				return nil
			}
			l.startConfig[curRow][curCol] = write
		}
		curRow++
	}
	if curRow != l.dimX {
		fmt.Println("invalid config")
		return nil
	}

	l.currentState = deep2DCopy(l.dimX, l.dimY, l.startConfig)
	l.tempState = deep2DCopy(l.dimX, l.dimY, l.startConfig)
	return l
}

func (l *life) printBoard() {
	for _, row := range l.currentState {
		for _, symbol := range row {
			fmt.Printf(" %c ", symbol)
		}
		fmt.Println()
	}
}

func (l *life) getAliveNeighboursCnt(x, y int) uint8 {
	var count uint8 = 0
	var xUpperBoarder bool = x == 0
	var xLowerBoarder bool = x == l.dimX-1
	var yLeftBoarder bool = y == 0
	var yRightBoarder bool = y == l.dimY-1

	if !xUpperBoarder {
		if !yLeftBoarder && l.currentState[x-1][y-1] == alive {
			count++
		}
		if l.currentState[x-1][y] == alive {
			count++
		}
		if !yRightBoarder && l.currentState[x-1][y+1] == alive {
			count++
		}
	}
	if !xLowerBoarder {
		if !yLeftBoarder && l.currentState[x+1][y-1] == alive {
			count++
		}
		if l.currentState[x+1][y] == alive {
			count++
		}
		if !yRightBoarder && l.currentState[x+1][y+1] == alive {
			count++
		}
	}
	if !yLeftBoarder && l.currentState[x][y-1] == alive {
		count++
	}
	if !yRightBoarder && l.currentState[x][y+1] == alive {
		count++
	}

	return count
}

func (l *life) livesOn(x, y int) bool {
	aliveNeighboursCnt := l.getAliveNeighboursCnt(x, y)
	if l.currentState[x][y] == alive {
		return aliveNeighboursCnt == 2 || aliveNeighboursCnt == 3
	} else {
		return aliveNeighboursCnt == 3
	}
}

func (l *life) nextGeneration() {
	for x, row := range l.currentState {
		for y, _ := range row {
			if l.livesOn(x, y) {
				l.tempState[x][y] = alive
			} else {
				l.tempState[x][y] = dead
			}
		}
	}
	tmp := l.currentState
	l.currentState = l.tempState
	l.tempState = tmp
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	l := NewLife("predefined_configs/pulsar")
	for {
		clearScreen()
		l.printBoard()
		l.nextGeneration()
		time.Sleep(time.Second)
	}
}
