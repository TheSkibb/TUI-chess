package main

import (
	"errors"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    cursor coordinate
    selected coordinate
    board [8][8]piece
    possibleMoves []coordinate
    capturedP1 []piece
    capturedP2 []piece
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

        case "testBishop":
        return nil, model {
                cursor: coordinate{0, 0},
                selected: coordinate{-1, -1},
                board: boardTestBishop,
            }

        case "testKnight":
        return nil, model {
                cursor: coordinate{0, 0},
                selected: coordinate{-1, -1},
                board: boardTestKnight,
            }
        
        case "testQueen":
        return nil, model {
                cursor: coordinate{0, 0},
                selected: coordinate{-1, -1},
                board: boardTestQueen,
            }

        case "testKing":
        return nil, model {
                cursor: coordinate{0, 0},
                selected: coordinate{-1, -1},
                board: boardTestKing,
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

    s := ""


    s += player2Name + ": [" +pieceArrToString(m.capturedP2) + "]\n"

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


    s += player1Name + ": [" + pieceArrToString(m.capturedP1) + "]\n"

    return s
}

func (m *model) calculateMoves(){

    piece := m.board[m.selected.y][m.selected.x]

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
            }

            m.possibleMoves = []coordinate{
            }

            //check forward
            if m.board[m.selected.y + 1 * direction][m.selected.x] == empty {
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

        /* bishop movement */
        case "♗", "♝":
            m.possibleMoves = []coordinate{}
            
            // right down
            for i := 1; m.selected.y + i < rowsAndColums && m.selected.x + i < rowsAndColums; i++ {
                moveCoor := coordinate{m.selected.x + i, m.selected.y + i}
                m.possibleMoves = append(m.possibleMoves, moveCoor)
            }

            // right up
            for i := 1; m.selected.y - i >= 0 && m.selected.x + i < rowsAndColums; i++ {
                moveCoor := coordinate{m.selected.x + i, m.selected.y - i}
                if m.checkIfEmpty(moveCoor){
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                } else if !m.checkIfSameColor(moveCoor, m.selected) {
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                    break
                } else {
                    break
                }
            }

            // left down
            for i := 1; m.selected.y + i < rowsAndColums && m.selected.x - i >= 0; i++ {
                moveCoor := coordinate{m.selected.x - i, m.selected.y + i}
                if m.checkIfEmpty(moveCoor){
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                } else if !m.checkIfSameColor(moveCoor, m.selected) {
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                    break
                } else {
                    break
                }
            }

            // left up
            for i := 1; m.selected.y - i >= 0 && m.selected.x - i >= 0; i++ {
                moveCoor := coordinate{m.selected.x - i, m.selected.y - i}
                if m.checkIfEmpty(moveCoor){
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                } else if !m.checkIfSameColor(moveCoor, m.selected) {
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                    break
                } else {
                    break
                }
            }

        /* knight movement */
        case "♞", "♘":
            m.possibleMoves = []coordinate{}

            checkMoves := []coordinate{
                {m.selected.x + 1, m.cursor.y + 2},
                {m.selected.x + 1, m.cursor.y - 2},
                {m.selected.x - 1, m.cursor.y + 2},
                {m.selected.x - 1, m.cursor.y - 2},
                {m.selected.x + 2, m.cursor.y + 1},
                {m.selected.x + 2, m.cursor.y - 1},
                {m.selected.x - 2, m.cursor.y + 1},
                {m.selected.x - 2, m.cursor.y - 1},
            }

            for _, move := range checkMoves {
                if move.x >= rowsAndColums || move.x < 0 || move.y >= rowsAndColums || move.y < 0 {
                    continue
                 }
                if m.checkIfEmpty(move){
                    m.possibleMoves = append(m.possibleMoves, move)
                } else if !m.checkIfSameColor(move, m.selected) {
                    m.possibleMoves = append(m.possibleMoves, move)
                    continue
                } else {
                    continue
                }
            }

        /* queen movement */
        case "♕", "♛":
            m.possibleMoves = []coordinate{}
            // queen logic is just copy paste rook and bishop

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
            
            // right down
            for i := 1; m.selected.y + i < rowsAndColums && m.selected.x + i < rowsAndColums; i++ {
                moveCoor := coordinate{m.selected.x + i, m.selected.y + i}
                m.possibleMoves = append(m.possibleMoves, moveCoor)
            }

            // right up
            for i := 1; m.selected.y - i >= 0 && m.selected.x + i < rowsAndColums; i++ {
                moveCoor := coordinate{m.selected.x + i, m.selected.y - i}
                if m.checkIfEmpty(moveCoor){
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                } else if !m.checkIfSameColor(moveCoor, m.selected) {
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                    break
                } else {
                    break
                }
            }

            // left down
            for i := 1; m.selected.y + i < rowsAndColums && m.selected.x - i >= 0; i++ {
                moveCoor := coordinate{m.selected.x - i, m.selected.y + i}
                if m.checkIfEmpty(moveCoor){
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                } else if !m.checkIfSameColor(moveCoor, m.selected) {
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                    break
                } else {
                    break
                }
            }

            // left up
            for i := 1; m.selected.y - i >= 0 && m.selected.x - i >= 0; i++ {
                moveCoor := coordinate{m.selected.x - i, m.selected.y - i}
                if m.checkIfEmpty(moveCoor){
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                } else if !m.checkIfSameColor(moveCoor, m.selected) {
                    m.possibleMoves = append(m.possibleMoves, moveCoor)
                    break
                } else {
                    break
                }
            }



        /* king movement */
        case "♔", "♚":
            m.possibleMoves = []coordinate{}
            checkMoves := []coordinate{
                {m.selected.x + 1, m.selected.y},
                {m.selected.x - 1, m.selected.y},
                {m.selected.x, m.selected.y + 1},
                {m.selected.x, m.selected.y - 1},
                {m.selected.x + 1, m.selected.y + 1},
                {m.selected.x + 1, m.selected.y - 1},
                {m.selected.x - 1, m.selected.y + 1},
                {m.selected.x - 1, m.selected.y - 1},
            }
            for _, move := range checkMoves {
                if move.x >= rowsAndColums || move.x < 0 || move.y >= rowsAndColums || move.y < 0 {
                    continue
                 }
                if m.checkIfEmpty(move){
                    m.possibleMoves = append(m.possibleMoves, move)
                } else if !m.checkIfSameColor(move, m.selected) {
                    m.possibleMoves = append(m.possibleMoves, move)
                    continue
                } else {
                    continue
                }
            }

        default:
            m.possibleMoves = []coordinate{}
        }

    }
}

//creates an array of coordinates {-1, -1}, which will not be 
func createEmptyCoordinateArray(length int) []coordinate{
    emptyArray := make([]coordinate, length)
    
    for i := 0; i < length; i++ {
        emptyArray[i] = coordinate{-1, -1}
    }

    return emptyArray
}

func (m model) checkIfEmpty(c coordinate) bool {
    if m.board[c.y][c.x] == empty {
        return true
    } 
    return false
}

func (m model) checkIfSameColor(c1 coordinate, c2 coordinate) bool {
    return m.board[c1.y][c1.x].pieceColor == m.board[c2.y][c2.x].pieceColor
}

func (m  *model) selectSquare(){
    
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

        //check if selection is an empty selectSquare
        if m.board[m.cursor.y][m.cursor.x] == empty {
            return
        }

        m.selected = coordinate{m.cursor.x, m.cursor.y}

        m.calculateMoves()
    }
}

func (m *model) movePiece(pos coordinate, piecePos coordinate){

    piece := m.board[m.selected.y][m.selected.x]

    //capturing
    if m.board[pos.y][pos.x] != empty {
        if m.board[pos.y][pos.x].pieceColor == white {
            m.capturedP2 = append(m.capturedP2, m.board[pos.y][pos.x])
        } else {
            m.capturedP1 = append(m.capturedP1, m.board[pos.y][pos.x])
        }
    }

    //move piece
    m.board[pos.y][pos.x] = piece
    m.board[m.selected.y][m.selected.x] = empty

    //reset selected and possibleMoves
    m.selected = coordinate{-1, -1}
    m.possibleMoves = []coordinate{}
}

//create a string of the unicode characters for an array of pieces
func pieceArrToString(a []piece) string {
    str := ""

    for i, piece := range a {
        str += piece.unicode
        if i != len(a) - 1 {
            str += " "
        }
    }

    return str
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
