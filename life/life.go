// Package life contains the logic running the game of life.
package life

import (
	"bufio"
	"fmt"
	"errors"
	"os"
	"path"
	"strconv"
	"strings"
)

type boardSymbol uint8

const alive boardSymbol = '*'
const dead boardSymbol = ' '
const configFolder = "predefined_configs"

// The game is represented by its current state.
// Additionally, the dimensions of the board are recorded as well as a
// temporary state used for computing the next generation and the beginning
// state is recorded so it can be restored.
type Life struct {
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

// Constructs a new game by reading the provided `config`.
// Configurations containing lines longer than 65536 characters are not supported.
func NewLife(config string) (*Life, error) {
	l := new(Life)

	configFile, err := os.Open(path.Join(configFolder, config))
	if err != nil {
		return nil, errors.New("the configuration you specified does not exist")
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
				return nil, errors.New("infalid config")
			}
			l.startConfig[curRow][curCol] = write
		}
		curRow++
	}
	if curRow != l.dimX {
		return nil, errors.New("infalid config")
	}

	l.currentState = deep2DCopy(l.dimX, l.dimY, l.startConfig)
	l.tempState = deep2DCopy(l.dimX, l.dimY, l.startConfig)
	return l, nil
}

// Returns a string representing the current state of the game.
func (l *Life) Printable() string {
	board := ""
	for _, row := range l.currentState {
		for _, symbol := range row {
			board += fmt.Sprintf(" %c ", symbol)
		}
		board += "\n"
	}
	return board
}

func (l *Life) getAliveNeighboursCnt(x, y int) uint8 {
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

func (l *Life) livesOn(x, y int) bool {
	aliveNeighboursCnt := l.getAliveNeighboursCnt(x, y)
	if l.currentState[x][y] == alive {
		return aliveNeighboursCnt == 2 || aliveNeighboursCnt == 3
	} else {
		return aliveNeighboursCnt == 3
	}
}

// Computes the next generation of the game and updates the state.
func (l *Life) NextGeneration() {
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

func isNum(char rune) bool {
	return char >= '0' && char <= '9'
}

func decodeInt(line string) (int, int) {
	read := 0
	for _, char := range line {
		if !isNum(char) {
			break
		}
		read++
	}
	res, _ := strconv.Atoi(line[:read])
	if res == 0 {
		res++
	}
	return res, read
}

func decodeConfig(path string) *Life {
	l := new(Life)

	config, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer config.Close()

	scanner := bufio.NewScanner(config)
	scanner.Scan()
	dimensions := strings.Split(scanner.Text(), " ")
	l.dimX, _ = strconv.Atoi(dimensions[0])
	l.dimY, _ = strconv.Atoi(dimensions[1])
	l.startConfig = make([][]boardSymbol, l.dimX)
	for i := 0; i < l.dimX; i++ {
		l.startConfig[i] = make([]boardSymbol, l.dimY)
	}

	currRow := 0
	currCol := 0
	var currSymbol byte
	for scanner.Scan() {
		line := scanner.Text()
		lineLen := len(line)
		read := 0
		var write boardSymbol
		for read < lineLen {
			num, nextSymbol := decodeInt(line)
			fmt.Println(num, nextSymbol)
			currSymbol = line[nextSymbol]
			read += nextSymbol + 1
			line = line[nextSymbol+1:]
			fmt.Println(line)
			if currSymbol == 'b' {
				write = dead
			} else if currSymbol == 'o' {
				write = alive
			} else if currSymbol == '$' {
				currRow++
				currCol = 0
				continue
			} else {
				break
			}
			for i := 0; i < num; i++ {
				l.startConfig[currRow][currCol] = write
				currCol++
			}
		}
	}
	l.currentState = deep2DCopy(l.dimX, l.dimY, l.startConfig)
	l.tempState = deep2DCopy(l.dimX, l.dimY, l.startConfig)
	return l
}
