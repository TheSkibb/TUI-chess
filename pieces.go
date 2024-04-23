package main

/* pieces */

var pawnBlack = piece{
    unicode : "♙",
    pieceColor: pieceColorBlack,
    possibleMoves: []coordinate{{1, 2}, {0, 0}},
}

var rookBlack = piece{
    unicode: "♖",
    pieceColor: pieceColorBlack,
}

var knightBlack = piece{
    unicode: "♘",
    pieceColor: pieceColorBlack,
}

var bishopBlack = piece{
    unicode: "♗",
    pieceColor: pieceColorBlack,
}

var queenBlack = piece{
    unicode: "♕",
    pieceColor: pieceColorBlack,
}

var kingBlack = piece{
    unicode: "♔",
    pieceColor: pieceColorBlack,
}

var pawnWhite = piece{
    unicode : "♟︎",
    pieceColor: pieceColorWhite,
    possibleMoves: []coordinate{{1, 2}, {0, 0}},
}

var rookWhite = piece{
    unicode: "♜",
    pieceColor: pieceColorWhite,
}

var knightWhite = piece{
    unicode: "♞",
    pieceColor: pieceColorWhite,
}

var bishopWhite = piece{
    unicode: "♝",
    pieceColor: pieceColorWhite,
}

var queenWhite = piece{
    unicode: "♛",
    pieceColor: pieceColorWhite,
}

var kingWhite = piece{
    unicode: "♚",
    pieceColor: pieceColorWhite,
}

var empty = piece{
    unicode: " ",
    pieceColor: "none",
}

