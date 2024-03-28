package main

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    cursorPosX int
    cursorPosY int
}

type chessPosition struct {
}

const rowsAndColums = 8
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
        cursorPosX: 0,
        cursorPosY: 0,
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
            if m.cursorPosY < rowsAndColums - 1 {
                m.cursorPosY ++
            }

        case "k", "up":
            if m.cursorPosY != 0 {
                m.cursorPosY --
            }

        case "l", "right":
            if m.cursorPosX < rowsAndColums - 1 {
                m.cursorPosX ++
            }
        case "h", "left":
            if m.cursorPosX != 0{
                m.cursorPosX --
            }

        }

    }
    return m, nil
}

func (m model) View() string{

    s := ""


    s += player1Name + ": []\n"

    s += "|---||---||---||---||---||---||---||---|\n"

    for i := 0; i < rowsAndColums; i++ {

        // draw cells
        for j := 0; j < rowsAndColums; j++ {
            if i == m.cursorPosY && j == m.cursorPosX {
                s += highlightColor + "| X |" + White
            } else {
                s += "|   |"
            }
        }

        s += "\n"

        // draw borders
        for j := 0; j < rowsAndColums; j++ {
            if i == m.cursorPosY && j == m.cursorPosX || i == m.cursorPosY - 1 && j == m.cursorPosX {
                s += highlightColor + "|---|" + White
            } else {
                s += "|---|"
            }
        }

        s += "\n"
    }


    s += player2Name + ": []\n"

    logToFile(strconv.Itoa(m.cursorPosX) + " " + strconv.Itoa(m.cursorPosY))

    return s
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
