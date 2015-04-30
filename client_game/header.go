package client_game

// File containing all of the game structures and constant values

type Game struct {
	Name        string
	PlayerID    string
	Board       int
	PlayerColor int
	TeamPlayer  int
	InGame      bool
}

type ChessBoard struct {
	Board         []*Square
	Len           int
	CurrentSquare int
	NextSquare    int
	Turn          int
	Time          int
}

type Square struct {
	Color          string
	Type           string
	Image          string
	Empty          bool
	FromOtherBoard bool
	OrigPosition   bool
	Index          int
	TeamPiece      int
}

type CapturedPieces struct {
	Pieces []*CapturedPiece
	Len    int
}

type CapturedPiece struct {
	TeamPiece int
	Image     string
	Type      string
}

//chatting message
type ChatMsg struct {
	Msg string
	//ChatTextChanged bool
}

const TEAM1 int = 1
const TEAM2 int = 2
const EMPTY int = -1
const WHITE int = 1
const BLACK int = 2
const COLOR_WHITE string = "#D18B47"
const COLOR_BLACK string = "#FFCE9E"
const BLACKROOK string = "../pieces/Rook_Black_60.png"
const BLACKKNIGHT string = "../pieces/Knight_Black_60.png"
const BLACKBISHOP string = "../pieces/Bishop_Black_60.png"
const BLACKQUEEN string = "../pieces/Queen_Black_60.png"
const BLACKKING string = "../pieces/King_Black_60.png"
const BLACKPAWN string = "../pieces/Pawn_Black_60.png"
const WHITEROOK string = "../pieces/Rook_White_60.png"
const WHITEKNIGHT string = "../pieces/Knight_White_60.png"
const WHITEBISHOP string = "../pieces/Bishop_White_60.png"
const WHITEQUEEN string = "../pieces/Queen_White_60.png"
const WHITEKING string = "../pieces/King_White_60.png"
const WHITEPAWN string = "../pieces/Pawn_White_60.png"
