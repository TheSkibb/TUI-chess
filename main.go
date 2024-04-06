package main

import (
	"fmt"
    "errors"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    cursor coordinate
    selected coordinate
    board [8][8]piece
    possibleMoves []coordinate
}

type coordinate struct {
    x int
    y int
}

type piece struct {
    unicode string
    pieceColor pieceColor
}

type pieceColor string

const (
    black pieceColor = "black"
    white pieceColor = "white"
)

const rowsAndColums = 8
const debugging = 1

//colors 

var Reset  = "\033[0m"
var Red    = "\033[31m"
var Green  = "\033[32m"
var Yellow = "\033[33m"
var Blue   = "\033[34m"
var Purple = "\033[35m"
var Cyan   = "\033[36m"
var Gray   = "\033[37m"
var White  = "\033[97m"

var highlightColor = Blue
var selectedColor = Red
var possibleMoveColor = Yellow

var player1Name = "player 1"
var player2Name = "player 2"


func main(){

    logToFile("**** new start ****")

    var model model
    var err error

    args := os.Args[1:]
    if len(args) == 0 {
        err, model = initialModel("default")
    } else {
        err, model = initialModel(args[0])
    }


    if err != nil {
        fmt.Printf("%v", err)
        os.Exit(1)
    }

    p := tea.NewProgram(model)

    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }

}

func initialModel(mode string) (error, model) {
    switch mode{
        case "default":
            return nil, model {
                cursor: coordinate{0, 0},
                selected: coordinate{-1, -1},
                board: boardDefault,
            }

        case "testRook":
            return nil, model {
                cursor: coordinate{0, 0},
                selected: coordinate{-1, -1},
                board: boardTestRook,
            }

        case "testPawn":
            return nil, model {
                cursor: coordinate{0, 0},
                selected: coordinate{-1, -1},
                board: boardTestPawn,
            }

        case "testEmpty":
            return nil, model {
                cursor: coordinate{0, 0},
                selected: coordinate{-1, -1},
                board: boardTestEmpty,
            }

    default:
        return errors.New("unrecognized board"), model{}
    }
}

func (m model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return nil
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    switch msg := msg.(type) {

    // Is it a key press?
    case tea.KeyMsg:

        // Cool, what was the actual key pressed?
        switch msg.String() {

        // These keys should exit the program.
        case "ctrl+c", "ctrl+d", "q":
            return m, tea.Quit

        case "j", "down":
            if m.cursor.y < rowsAndColums - 1 {
                m.cursor.y ++
            }

        case "k", "up":
            if m.cursor.y != 0 {
                m.cursor.y --
            }

        case "l", "right":
            if m.cursor.x < rowsAndColums - 1 {
                m.cursor.x ++
            }
        case "h", "left":
            if m.cursor.x != 0{
                m.cursor.x --
            }
        case "enter", " ":
            m.selectSquare()
        }

    }
    return m, nil
}

func (m model) View() string{

    // every cell is 6x3
    // |---||---|
    // |   ||   |
    // |---||---|

    s := ""


    s += player2Name + ": []\n"

    s += "|---||---||---||---||---||---||---||---|\n"

    for i := 0; i < rowsAndColums; i++ {

        // draw cells
        for j := 0; j < rowsAndColums; j++ {

            piece := m.board[i][j]

            color := White

            if i == m.cursor.y && j == m.cursor.x {
                color = highlightColor
            }

            if i == m.selected.y && j == m.selected.x {
                color = selectedColor
            }

            for _, possibleMove := range m.possibleMoves{
                if i == possibleMove.y && j == possibleMove.x {
                    color = possibleMoveColor
                    break
                }
            }

            s += color + "| " + piece.unicode +" |" + White
        }

        s += "\n"

        // draw borders
        for j := 0; j < rowsAndColums; j++ {
            color := White

            for _, possibleMove := range m.possibleMoves{
                if i == possibleMove.y && j == possibleMove.x || i == possibleMove.y - 1 && j == possibleMove.x {
                    color = possibleMoveColor
                    break
                }
            }
            if i == m.selected.y && j == m.selected.x || i == m.selected.y - 1 && j == m.selected.x {
                color = selectedColor
            }

            if i == m.cursor.y && j == m.cursor.x || i == m.cursor.y - 1 && j == m.cursor.x {
                color = highlightColor
            }

            s += color + "|---|" + White
        }

        s += "\n"
    }


    s += player1Name + ": []\n"

    logToFile("cursor: " + strconv.Itoa(m.cursor.x) + " " + strconv.Itoa(m.cursor.y) +
       "selected: " + strconv.Itoa(m.selected.x) + strconv.Itoa(m.selected.y))

    return s
}

func (m *model) calculateMoves(){
    //logToFile(strconv.Itoa(m.selected.x) + strconv.Itoa(m.selected.y))
    piece := m.board[m.selected.y][m.selected.x]
    logToFile("calculating" + piece.unicode)

    //direction is either 1 if piece is black, or -1 if it is white
    direction := 1

    if piece.unicode == " " {
        m.possibleMoves = []coordinate{}
    }else{
        switch piece.unicode{

        /* pawn movement */
        case "♙", "♟︎":
            if piece.pieceColor == white {
                direction = -1
                logToFile("is white piece")
            }

            m.possibleMoves = []coordinate{
            }

            //check forward
            if m.board[m.selected.y + 1 * direction][m.selected.x] == empty {
                logToFile("next square is empty")
                m.possibleMoves = append(m.possibleMoves, coordinate{m.selected.x, m.selected.y + 1 * direction})
            }

            //check if in initial position
            if (piece.pieceColor == "black" && m.selected.y == 1 ||
            piece.pieceColor == "white" && m.selected.y == 6 ) && 
            m.board[m.selected.y + 2 * direction][m.selected.x] == empty {
                m.possibleMoves = append(m.possibleMoves, coordinate{m.selected.x, m.selected.y + 2 * direction})
            }

            //check for diagonals
            if m.selected.x != 0 && 
            m.board[m.selected.y + 1 * direction][m.selected.x - 1] != empty && 
            !m.checkIfSameColor(coordinate{m.selected.x - 1, m.selected.y + 1 * direction}, coordinate{m.selected.x, m.cursor.y}) {
                m.possibleMoves = append(m.possibleMoves, coordinate{m.selected.x - 1, m.selected.y + 1 * direction})
            }
            if m.selected.x != 7 && 
            m.board[m.selected.y + 1 * direction][m.selected.x + 1] != empty &&
            !m.checkIfSameColor(coordinate{m.selected.x + 1, m.selected.y + 1 * direction}, coordinate{m.selected.x, m.cursor.y}) {
                m.possibleMoves = append(m.possibleMoves, coordinate{m.selected.x + 1, m.selected.y + 1 * direction})
            }

        /* rook movement */
        case "♖", "♜":
            m.possibleMoves = []coordinate{
            }

            // down
            for i := 1; m.selected.y + i < rowsAndColums; i ++ {
                moveCoor := coordinate{m.selected.x, m.cursor.y + i}

                if m.checkIfEmpty(moveCoor){
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                } else if !m.checkIfSameColor(moveCoor, m.selected) {
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                    break
                } else {
                    break
                }
            }

            // up
            for i := 1; m.selected.y - i >= 0; i ++ {
                moveCoor := coordinate{m.selected.x, m.cursor.y - i}

                if m.checkIfEmpty(moveCoor){
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                } else if !m.checkIfSameColor(moveCoor, m.selected) {
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                    break
                } else {
                    break
                }
            }

            //left
            for i := 1; m.selected.x + i < rowsAndColums; i++ {
                moveCoor := coordinate{m.selected.x + i, m.cursor.y}

                if m.checkIfEmpty(moveCoor){
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                } else if !m.checkIfSameColor(moveCoor, m.selected) {
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                    break
                } else {
                    break
                }
            }

            // left
            for i := 1; m.selected.x - i >= 0; i++ {
                
                moveCoor := coordinate{m.selected.x - i, m.cursor.y}

                if m.checkIfEmpty(moveCoor){
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                } else if !m.checkIfSameColor(moveCoor, m.selected) {
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                    break
                } else {
                    break
                }
            }


        default:
            m.possibleMoves = []coordinate{}
        }

    }

    logToFile("length of possible moves: " + strconv.Itoa(len(m.possibleMoves)))
}

//creates an array of coordinates {-1, -1}, which will not be 
func createEmptyCoordinateArray(length int) []coordinate{
    emptyArray := make([]coordinate, length)
    
    for i := 0; i < length; i++ {
        emptyArray[i] = coordinate{-1, -1}
    }

    logToFile(strconv.Itoa(emptyArray[10].x))

    return emptyArray
}

func (m model) checkIfEmpty(c coordinate) bool {
    if m.board[c.y][c.x] == empty {
        return true
    } 
    return false
}

func (m model) checkIfSameColor(c1 coordinate, c2 coordinate) bool {
    logToFile("***** colors***")
    logToFile(string(m.board[c1.y][c1.x].pieceColor))
    logToFile(string(m.board[c2.y][c2.x].pieceColor))
    return m.board[c1.y][c1.x].pieceColor == m.board[c2.y][c2.x].pieceColor
}

func (m  *model) selectSquare(){
    logToFile("selected: " + strconv.Itoa(m.cursor.x) + " " + strconv.Itoa(m.cursor.y))
    
    //check player is deselecting
    if m.selected == m.cursor{
        m.selected = coordinate{-1, -1}
        m.possibleMoves = []coordinate{}
    } else {
        /* selecting */

        //check if selection is a possible move
        for _, pos := range m.possibleMoves{
            if pos == m.cursor{
                m.movePiece(pos, m.cursor)
                return
            }
        }

        m.selected = coordinate{m.cursor.x, m.cursor.y}

        m.calculateMoves()
    }
}

func (m *model) movePiece(pos coordinate, piecePos coordinate){

    piece := m.board[m.selected.y][m.selected.x]

    logToFile("moving piece" + piece.unicode)

    //move piece
    m.board[pos.y][pos.x] = piece
    m.board[m.selected.y][m.selected.x] = empty

    //reset selected and possibleMoves
    m.selected = coordinate{-1, -1}
    m.possibleMoves = []coordinate{}
}

func logToFile(msg string){
    if debugging == 0 {
        return
    }

    f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}

    _, err = f.WriteString(msg + "\n")

    if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
    }

	defer f.Close()
}
