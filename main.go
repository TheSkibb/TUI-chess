package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
    "strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    cursor coordinate
    selected coordinate
    board [8][8]piece
    capturedP1 []piece
    capturedP2 []piece
    moveLog []string
    player1 player
    player2 player

    // player turn can be 1 for player1, 2 for player2 or 0 if freemoving
    playerTurn int
}

type coordinate struct {
    x int
    y int
}

type piece struct {
    unicode string
    pieceColor pieceColor
    possibleMoves []coordinate
}

type player struct {
    name string
    possibleMoves []coordinate
    checked bool
}

type pieceColor string

const (
    pieceColorBlack pieceColor = "black"
    pieceColorWhite pieceColor = "white"
)

const rowsAndColums = 8
const debugging = 1

/* colors */
const Reset  = "\033[0m"
const Red    = "\033[31m"
const Green  = "\033[32m"
const Yellow = "\033[33m"
const Blue   = "\033[34m"
const Purple = "\033[35m"
const Cyan   = "\033[36m"
const Gray   = "\033[37m"
const White  = "\033[97m"

/* text color defaults */
var highlightColor = Blue
var selectedColor = Red
var possibleMoveColor = Yellow
var boardColor = White
var pieceMarkupColor = White

func main(){

    logToFile("**** new start ****")

    // read config file
    err := setColors()

    if err != nil {
        fmt.Println("error in config file")
        return 
    }

    var model model

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
    m := model {
        cursor: coordinate{4, 4},
        selected: coordinate{-1, -1},
        board: boardDefault,
        player1: player{
            name: "player 1",
            checked: false,
        },
        player2: player{
            name: "player 2",
            checked: false,
        },
    }
    switch mode{
        case "default":
            m.cursor = coordinate{4, 7}
            m.playerTurn = 1
            m.calculateMoves()
            return nil, m
        case "freeplay":
            m.cursor = coordinate{4, 7}
            m.calculateMoves()
            return nil, m

        case "testRook":
            m.board = boardTestRook
            m.calculateMoves()
            return nil, m

        case "testPawn":
            m.board = boardTestPawn
            m.calculateMoves()
            return nil, m

        case "testBishop":
            m.board = boardTestBishop
            m.calculateMoves()
            return nil, m

        case "testKnight":
            m.board = boardTestKnight
            m.calculateMoves()
            return nil, m
        
        case "testQueen":
            m.board = boardTestQueen
            m.calculateMoves()
        return nil, m

        case "testKing":
            m.board = boardTestKing
            m.calculateMoves()
        return nil, m

        case "testEmpty":
            m.board = boardTestEmpty
            return nil, m

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

        /* quit program */
        case "ctrl+c", "ctrl+d", "q":
            fmt.Println("Thanks for playing!")
            return m, tea.Quit

        /* move cursor down */
        case "j", "down":
            if m.cursor.y < rowsAndColums - 1 {
                m.cursor.y ++
            }

        /* move cursor up */
        case "k", "up":
            if m.cursor.y != 0 {
                m.cursor.y --
            }

        /* move cursor right */
        case "l", "right":
            if m.cursor.x < rowsAndColums - 1 {
                m.cursor.x ++
            }

        /* move cursor left */
        case "h", "left":
            if m.cursor.x != 0{
                m.cursor.x --
            }

        /* select piece */
        case "enter", " ":
            m.selectSquare()
        }

    }
    return m, nil
}

func (m model) View() string{

    s := ""

    s += m.player2.name + ": [" +pieceArrToString(m.capturedP2) + "]\n"

    s += boardColor + "|---||---||---||---||---||---||---||---|\n"

    var selectedPiece piece

    if m.selected.x != -1 {
        selectedPiece = m.board[m.selected.y][m.selected.x]
        logToFile("the selected piece has " + strconv.Itoa(len(selectedPiece.possibleMoves)) + "options")
    }

    for i := 0; i < rowsAndColums; i++ {

        // draw cells
        for j := 0; j < rowsAndColums; j++ {

            piece := m.board[i][j]

            color := boardColor

            if i == m.cursor.y && j == m.cursor.x {
                color = highlightColor
            }

            if i == m.selected.y && j == m.selected.x {
                color = selectedColor
            }

            if m.selected.x != -1 {
                for _, possibleMove := range selectedPiece.possibleMoves{
                    if i == possibleMove.y && j == possibleMove.x {
                        color = possibleMoveColor
                        break
                    }
                }
            }

            s += color + "| " + pieceMarkupColor + piece.unicode + color + " |" + boardColor
        }

        s += "\n"

        // draw borders
        for j := 0; j < rowsAndColums; j++ {
            color := boardColor

            if m.selected.x != -1 {
                for _, possibleMove := range selectedPiece.possibleMoves{
                    if i == possibleMove.y && j == possibleMove.x || i == possibleMove.y - 1 && j == possibleMove.x {
                        color = possibleMoveColor
                        break
                    }
                }
            }

            if i == m.selected.y && j == m.selected.x || i == m.selected.y - 1 && j == m.selected.x {
                color = selectedColor
            }

            if i == m.cursor.y && j == m.cursor.x || i == m.cursor.y - 1 && j == m.cursor.x {
                color = highlightColor
            }

            s += color + "|---|" + boardColor
        }

        s += "\n"
    }


    s += m.player1.name + ": [" + pieceArrToString(m.capturedP1) + "]\n"

    return s
}

func (m *model) calculateMoves(){
    //iterate over all pieces on the board and calculate the possible moves of each piece
    logToFile("calculating all possible moves")

    for i := 0; i < rowsAndColums; i++ {
        for j := 0; j < rowsAndColums; j++ {
            m.calculatePossibleMoves(coordinate{i, j})
        }
    }
    logToFile("should have " + strconv.Itoa(len(m.board[6][0].possibleMoves)) + "options")
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
    if m.board[c.y][c.x].unicode == empty.unicode {
        return true
    } 
    return false
}

func (m model) checkIfChecking(c coordinate) bool {
    logToFile("king is checked")
    piece := m.board[c.y][c.x]
    if piece.unicode == kingBlack.unicode || piece.unicode == kingWhite.unicode {
        return true
    }

    return false
}

func (m model) checkIfSameColor(c1 coordinate, c2 coordinate) bool {
    return m.board[c1.y][c1.x].pieceColor == m.board[c2.y][c2.x].pieceColor
}

// checks if the move is a valid move and whether or not the move puts the king in check
// a move that puts the king in check will be invalid
func (m model) checkValidMove(piecePos, movePos coordinate) (valid bool, checking bool) {

    valid = false


    checking = m.checkIfChecking(movePos)

    return false, false
}

func (m  *model) selectSquare(){

    logToFile("square selected")

    cursorpiece := m.board[m.cursor.y][m.cursor.x]
    var selectedPiece piece

    if m.selected.x != -1 {
        selectedPiece = m.board[m.selected.y][m.selected.x]
    }
    
    //check player is deselecting
    if m.selected == m.cursor{
        m.selected = coordinate{-1, -1}
        // TODO fix deselection to stop displaying the currently selected piece
        //m.possibleMoves = []coordinate{}

    } else {

        //check if selection is a possible move
        /* TODO fix moving of pieces by selecting a possible move */
        if m.selected.x != -1 {
            for _, pos := range selectedPiece.possibleMoves{
                logToFile(strconv.Itoa(pos.x) + strconv.Itoa(pos.y))
                logToFile(strconv.Itoa(m.cursor.x) + strconv.Itoa(m.cursor.y))
                if pos == m.cursor{
                    m.movePiece(pos, m.cursor)
                    logToFile("moved piece")
                    return
                }
            }
        }

        //check if player is selecting the right piece
        if (m.playerTurn == 1 && cursorpiece.pieceColor == pieceColorBlack) || 
            (m.playerTurn == 2 && cursorpiece.pieceColor == pieceColorWhite) {
            return
        }

        //check if selection is an empty selectSquare
        if m.board[m.cursor.y][m.cursor.x].unicode == empty.unicode {
            return
        }

        m.selected = coordinate{m.cursor.x, m.cursor.y}
    }
}

func (m *model) movePiece(pos coordinate, piecePos coordinate){

    piece := m.board[m.selected.y][m.selected.x]

    //capturing
    if m.board[pos.y][pos.x].unicode != empty.unicode {
        if m.board[pos.y][pos.x].pieceColor == pieceColorWhite {
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

    // TODO fix to recalculate all possible moves
    //m.possibleMoves = []coordinate{}

    // switch turn
    if m.playerTurn == 1 {
        m.playerTurn = 2
    } else if m.playerTurn == 2 {
        m.playerTurn = 1
    }

    //recalculate possible moves

    m.calculateMoves()
}

func (m *model) calculatePossibleMoves(c coordinate) (checking bool){
    piece := m.board[c.y][c.x]
    m.player1.checked = false
    m.player2.checked = false

    switch piece.unicode{
        /* pawn movement */
        case "♙", "♟︎":
            logToFile("calculating pawn moves")
            return m.calculatePossibleMovesPawn(c)

        /* rook movement */
        case "♖", "♜":
            logToFile("calculating rook moves")
            return m.calculatePossibleMovesRook(c)

        /* bishop movement */
        case "♗", "♝":
            return m.calculatePossibleMovesBishop(c)

        /* knight movement */
        case "♞", "♘":
            return m.calculatePossibleMovesKnight(c)

        /* queen movement */
        case "♕", "♛":
            return m.calculatePossibleMovesQueen(c)

        /* king movement */
        case "♔", "♚":
            return m.calculatePossibleMovesKing(c)
    }
    return false
}

func (m *model) calculatePossibleMovesPawn(pos coordinate) (checking bool) {
    piece := &m.board[pos.y][pos.x]

    direction := 1

    if piece.pieceColor == pieceColorWhite {
        direction = -1
    }

    piece.possibleMoves = []coordinate{
    }

    //check forward
    if pos.y != 0 && pos.y != rowsAndColums -1 && m.board[pos.y + 1 * direction][pos.x].unicode == empty.unicode {
        piece.possibleMoves = append(piece.possibleMoves, coordinate{pos.x, pos.y + 1 * direction})
    }

    //check if in initial position
    if (piece.pieceColor == "black" && pos.y == 1 ||
    piece.pieceColor == "white" && pos.y == 6 ) && 
    m.board[pos.y + 2 * direction][pos.x].unicode == empty.unicode {
        piece.possibleMoves = append(piece.possibleMoves, coordinate{pos.x, pos.y + 2 * direction})
    }

    //check for diagonals
    if pos.x != 0 && pos.y != 0 && pos.y != rowsAndColums - 1 &&
    m.board[pos.y + 1 * direction][pos.x - 1].unicode != empty.unicode && 
    !m.checkIfSameColor(coordinate{pos.x - 1, pos.y + 1 * direction}, coordinate{pos.x, m.cursor.y}) {
        if m.checkIfChecking(coordinate{pos.x - 1, pos.y + 1 * direction}) {
            checking = true
        }
        piece.possibleMoves = append(piece.possibleMoves, coordinate{pos.x - 1, pos.y + 1 * direction})
    }
    if pos.x != 7 && pos.y != 0 && pos.y != rowsAndColums - 1 &&
    m.board[pos.y + 1 * direction][pos.x + 1].unicode != empty.unicode &&
    !m.checkIfSameColor(coordinate{pos.x + 1, pos.y + 1 * direction}, coordinate{pos.x, m.cursor.y}) {
        if m.checkIfChecking(coordinate{pos.x - 1, pos.y + 1 * direction}) {
            checking = true
        }
        piece.possibleMoves = append(piece.possibleMoves, coordinate{pos.x + 1, pos.y + 1 * direction})
    }
    return checking
}

func (m *model) calculatePossibleMovesRook(pos coordinate) (checking bool) {

    piece := &m.board[pos.y][pos.x]

    piece.possibleMoves = []coordinate{
    }

    // down
    for i := 1; pos.y + i < rowsAndColums; i ++ {
        moveCoor := coordinate{pos.x, pos.y + i}

        if m.checkIfEmpty(moveCoor){
            // if empty no need to check if checking
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }

    // up
    for i := 1; pos.y - i >= 0; i ++ {
        moveCoor := coordinate{pos.x, pos.y - i}

        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }

    // right
    for i := 1; pos.x + i < rowsAndColums; i++ {
        moveCoor := coordinate{pos.x + i, pos.y}

        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }

    // left
    for i := 1; pos.x - i >= 0; i++ {
        
        moveCoor := coordinate{pos.x - i, pos.y}

        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }

    return checking
}

func (m *model) calculatePossibleMovesBishop(pos coordinate) (checking bool) {

    piece := &m.board[pos.y][pos.x]

    piece.possibleMoves = []coordinate{}
            
    // right down
    for i := 1; pos.y + i < rowsAndColums && pos.x + i < rowsAndColums; i++ {
        moveCoor := coordinate{pos.x + i, pos.y + i}
        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }

    // right up
    for i := 1; pos.y - i >= 0 && pos.x + i < rowsAndColums; i++ {
        moveCoor := coordinate{pos.x + i, pos.y - i}
        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }

    // left down
    for i := 1; pos.y + i < rowsAndColums && pos.x - i >= 0; i++ {
        moveCoor := coordinate{pos.x - i, pos.y + i}
        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }

    // left up
    for i := 1; pos.y - i >= 0 && pos.x - i >= 0; i++ {
        moveCoor := coordinate{pos.x - i, pos.y - i}
        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }
    return checking
}

func (m *model) calculatePossibleMovesKnight(pos coordinate) (checking bool) {
    piece := &m.board[pos.y][pos.x]
    logToFile("calculating knights")

    piece.possibleMoves = []coordinate{}

    checkMoves := []coordinate{
        {pos.x + 1, pos.y + 2},
        {pos.x + 1, pos.y - 2},
        {pos.x - 1, pos.y + 2},
        {pos.x - 1, pos.y - 2},
        {pos.x + 2, pos.y + 1},
        {pos.x + 2, pos.y - 1},
        {pos.x - 2, pos.y + 1},
        {pos.x - 2, pos.y - 1},
    }

    for _, move := range checkMoves {
        if move.x >= rowsAndColums || move.x < 0 || move.y >= rowsAndColums || move.y < 0 {
            continue
         }
        if m.checkIfEmpty(move){
            piece.possibleMoves = append(piece.possibleMoves, move)
        } else if !m.checkIfSameColor(move, pos) {
            if m.checkIfChecking(move) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, move)
            continue
        } else {
            continue
        }
    }
    return checking
}

func (m *model) calculatePossibleMovesQueen(pos coordinate) (checking bool) {

    piece := &m.board[pos.y][pos.x]

    piece.possibleMoves = []coordinate{}
    // queen logic is just copy paste rook and bishop

    // down
    for i := 1; pos.y + i < rowsAndColums; i ++ {
        moveCoor := coordinate{pos.x, pos.y + i}

        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }

    // up
    for i := 1; pos.y - i >= 0; i ++ {
        moveCoor := coordinate{pos.x, pos.y - i}

        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }

    //left
    for i := 1; pos.x + i < rowsAndColums; i++ {
        moveCoor := coordinate{pos.x + i, pos.y}

        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }

    // left
    for i := 1; pos.x - i >= 0; i++ {
        
        moveCoor := coordinate{pos.x - i, pos.y}

        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }
    
    // right down
    for i := 1; pos.y + i < rowsAndColums && pos.x + i < rowsAndColums; i++ {
        moveCoor := coordinate{pos.x + i, pos.y + i}
        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }

    // right up
    for i := 1; pos.y - i >= 0 && pos.x + i < rowsAndColums; i++ {
        moveCoor := coordinate{pos.x + i, pos.y - i}
        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }

    // left down
    for i := 1; pos.y + i < rowsAndColums && pos.x - i >= 0; i++ {
        moveCoor := coordinate{pos.x - i, pos.y + i}
        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }

    // left up
    for i := 1; pos.y - i >= 0 && pos.x - i >= 0; i++ {
        moveCoor := coordinate{pos.x - i, pos.y - i}
        if m.checkIfEmpty(moveCoor){
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
        } else if !m.checkIfSameColor(moveCoor, pos) {
            if m.checkIfChecking(moveCoor) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, moveCoor)
            break
        } else {
            break
        }
    }

    return checking
}

func (m *model) calculatePossibleMovesKing(pos coordinate) (checking bool) {

    piece := &m.board[pos.y][pos.x]
    piece.possibleMoves = []coordinate{} 
    checkMoves := []coordinate{
        {pos.x + 1, pos.y},
        {pos.x - 1, pos.y},
        {pos.x, pos.y + 1},
        {pos.x, pos.y - 1},
        {pos.x + 1, pos.y + 1},
        {pos.x + 1, pos.y - 1},
        {pos.x - 1, pos.y + 1},
        {pos.x - 1, pos.y - 1},
    }
    for _, move := range checkMoves {
        if move.x >= rowsAndColums || move.x < 0 || move.y >= rowsAndColums || move.y < 0 {
            continue
         }
        if m.checkIfEmpty(move){
            piece.possibleMoves = append(piece.possibleMoves, move)
        } else if !m.checkIfSameColor(move, pos) {
            if m.checkIfChecking(move) {
                checking = true
            }
            piece.possibleMoves = append(piece.possibleMoves, move)
            continue
        } else {
            continue
        }
    }
    return checking
}

// adds the move to the list of moves formatted as a chess move
// https://www.chessstrategyonline.com/content/tutorials/basic-chess-concepts-chess-notation
// https://en.wikipedia.org/wiki/Portable_Game_Notation
func (m model) logMove(original_pos coordinate, new_pos coordinate){
}

// check if the king is in check
func (m model) checkForCheck() bool{
    return false
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

func setColors() error {

    readFile, err := os.Open("./conf.txt")

    if err != nil {
        return nil
    }
    fileScanner := bufio.NewScanner(readFile)

    fileScanner.Split(bufio.ScanLines)

    for fileScanner.Scan() {
        err = setColorConfig(fileScanner.Text())
    }

    readFile.Close()

    return err
}

func setColorConfig(config string) error {
    line_split := strings.Split(config, " ")
    var err error

    if len(line_split) != 2 {
        return errors.New("malformed config line")
    }

    switch line_split[0]{
        case "board-color": 
            err, boardColor = getColor(line_split[1])
        case "piece-color":
            err, pieceMarkupColor = getColor(line_split[1])
        case "select-color":
            err, selectedColor = getColor(line_split[1])
        case "highlight-color":
            err, highlightColor = getColor(line_split[1])
        case "possible-color":
            err, possibleMoveColor = getColor(line_split[1])
    }

    return err
}

func getColor(c string) (error, string) {
    switch c{
        case "red", "Red":
            return nil, Red
        case "green", "Green":
            return nil, Green
        case "yellow", "Yellow":
            return nil, Yellow
        case "blue", "Blue":
            return nil, Blue
        case "purple", "Purple":
            return nil, Purple
        case "cyan", "Cyan":
            return nil, Cyan
        case "gray", "Gray":
            return nil, Gray
        case "white", "White":
            return nil, White
    }
    return nil, ""
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
