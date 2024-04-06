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

var boardTestPawn = [8][8]piece{
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {pawnBlack, pawnBlack, pawnBlack, pawnBlack, pawnBlack, pawnBlack, pawnBlack, pawnBlack},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, empty, pawnBlack, empty, empty, empty, empty, empty},
   {empty, empty, empty, pawnWhite, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {pawnWhite, pawnWhite, pawnWhite, pawnWhite, pawnWhite, pawnWhite, pawnWhite, pawnWhite},
   {empty, empty, empty, empty, empty, empty, empty, empty},
}

var boardTestRook = [8][8]piece{
   {empty, empty, pawnWhite, empty, empty, pawnWhite, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, pawnWhite, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, rookBlack, empty, empty, pawnWhite, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, pawnWhite, empty, empty, rookBlack, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
}

var boardTestBishop = [8][8]piece{
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, bishopWhite, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, pawnWhite, empty, bishopBlack, empty, empty, empty, empty},
   {empty, pawnWhite, empty, empty, empty, empty, empty, empty},
   {empty, pawnWhite, empty, empty, empty, empty, empty, empty},
   {empty, pawnWhite, empty, empty, empty, empty, empty, empty},
   {empty, pawnWhite, empty, empty, empty, empty, empty, empty},
}

var boardTestEmpty = [8][8]piece{
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
   {empty, empty, empty, empty, empty, empty, empty, empty},
}


