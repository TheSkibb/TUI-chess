package main

var boardDefault = [8][8]piece{
   {rookBlack, knightBlack, bishopBlack, queenBlack, kingBlack, knightBlack, bishopBlack, rookBlack},
   {pawnBlack, pawnBlack, pawnBlack, pawnBlack, pawnBlack, pawnBlack, pawnBlack, pawnBlack},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {pawnWhite, pawnWhite, pawnWhite, pawnWhite, pawnWhite, pawnWhite, pawnWhite, pawnWhite},
   {rookWhite, knightWhite, bishopWhite, queenWhite, kingWhite, bishopWhite, knightWhite, rookWhite},
}

var boardTestRook = [8][8]piece{
   {rookBlack, knightBlack, bishopBlack, queenBlack, kingBlack, knightBlack, bishopBlack, rookBlack},
   {pawnBlack, pawnBlack, pawnBlack, pawnBlack, pawnBlack, pawnBlack, pawnBlack, pawnBlack},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, rookBlack, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, rookBlack, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
}
