package main

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    cursor coordinate
    selected coordinate
    board [8][8]string
}

type coordinate struct {
    x int
    y int
}

type chessPosition struct {
}

const rowsAndColums = 8

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

var player1Name = "player 1"
var player2Name = "player 2"


func main(){
    logToFile("**** new start ****")
    p := tea.NewProgram(initialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }

}

func initialModel() model {
    return model {
        cursor: coordinate{0, 0},
        selected: coordinate{-1, -1},
        board: [8][8]string{
           {"♖", "♘", "♗", "♕", "♔", "♘", "♗", "♖"},
           {"♗", "♗", "♗", "♗", "♗", "♗", "♗", "♗"},
           {" ", " ", " ", " ", " ", " ", " ", " "},
           {" ", " ", " ", " ", " ", " ", " ", " "},
           {" ", " ", " ", " ", " ", " ", " ", " "},
           {" ", " ", " ", " ", " ", " ", " ", " "},
           {"♟︎", "♟︎", "♟︎", "♟︎", "♟︎", "♟︎", "♟︎", "♟︎"},
           {"♜", "♞", "♝", "♛", "♚", "♝", "♞", "♜"},
        },
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
        case "ctrl+c", "q":
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


    s += player1Name + ": []\n"

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

            s += color + "| " + piece +" |" + White
        }

        s += "\n"

        // draw borders
        for j := 0; j < rowsAndColums; j++ {
            color := White


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


    s += player2Name + ": []\n"

    logToFile("cursor: " + strconv.Itoa(m.cursor.x) + " " + strconv.Itoa(m.cursor.y) +
       "selected: " + strconv.Itoa(m.selected.x) + strconv.Itoa(m.selected.y))

    return s
}

func (m  *model) selectSquare(){
    logToFile("selected: " + strconv.Itoa(m.cursor.x) + " " + strconv.Itoa(m.cursor.y))
    if m.selected == m.cursor{
        m.selected = coordinate{-1, -1}
    } else {
        m.selected = coordinate{m.cursor.x, m.cursor.y}
    }
}

func logToFile(msg string){
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
