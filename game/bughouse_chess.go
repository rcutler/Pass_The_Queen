/*
	This is the main go file. Will instantiate the windows.
	Will first start the global window at start up.
	Handle the networking code.
	When a game is created/started, will switch over to the
	game windows/view.
	Will need to add state information for the game such as game id,
	player ids, which player is on which team, which player is on which
	board, board id, which players are associated with what board.

	TODO: Chessboard
	- Add state to only allow one move. Then switch to other teams turn
	when a button is pressed.
	- Allow for state to be reverted at press of a button instead of
	switching to other teams turn.
	- Add in state for clocks relating to the game/board. (add a game struct?)
	- Create an array of games based off of how many players there are. Use
	the networking code/logic to determine this.

	TODO: CapturePieces
	- Add state functionality for captured pieces.
	- Something along the lines of Type, Color, and image.
	- Be able to add the piece to the board.

	TODO: Break out the game view, global view into their own files.
	TODO: Have the main windows set up the networking connections and then
	set up the window
	TODO: Add support for changing the view when entering a game or finishing a game

*/

package main

import (
	"fmt"
	"gopkg.in/qml.v1"
	"os"
	"time"
)

func main() {
	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	engine := qml.NewEngine()

	// Create the custom type here for either pieces or squares
	chessBoard := &ChessBoard{}
	// Set the variable in the context for the custom type above.
	engine.Context().SetVar("chessBoard", chessBoard)

	// Set context variables for the squares that are selected. Along with rest of the board
	chessBoard.initialize()

	// Test to make sure everything in chessBoard is initialized properly.
	for i := 0; i < 64; i++ {
		s := chessBoard.Square(i)
		fmt.Println("Index: ", s.Index, "  Color: ", s.Color, " Piece Type: ", s.PieceType, " Piece image loc: ", s.PieceImage, " Piece color: ", s.PieceTeam)
	}

	// Load the main qml file

	component, err := engine.LoadFile("GameView.qml")
	if err != nil {
		return err
	}

	// Create the new window
	window := component.CreateWindow(nil)
	window.Show()

	time.Sleep(2 * time.Second)
	chessBoard.MovePiece(1, 16)

	window.Wait()
	return nil
}

type ChessBoard struct {
	Board   []*Square
	Len     int
	SSquare int
	NSquare int
	CTurn   bool // If Whites turn, will be true. Will be false if blacks turn.
}

type Square struct {
	Index               int
	Color               string
	PieceType           string
	PieceImage          string
	PieceEmpty          bool
	PieceTeam           bool
	PieceFromOtherBoard bool
	PieceOrigPosition   bool
}

func (c *ChessBoard) initialize() {
	// Initialize the empty board
	c.SSquare = -1
	c.NSquare = -1
	c.CTurn = true
	row := 0
	for i := 0; i < 64; i++ {
		s := new(Square)
		if (i % 8) == 0 {
			row++
		}
		if row%2 == 0 && i%2 == 0 {
			s.Color = "#D18B47"
		} else if row%2 == 0 && i%2 == 1 {
			s.Color = "#FFCE9E"
		} else if row%2 == 1 && i%2 == 0 {
			s.Color = "#FFCE9E"
		} else {
			s.Color = "#D18B47"
		}
		s.Index = i
		s.PieceType = "Empty"
		s.PieceImage = ""
		s.PieceEmpty = true
		s.PieceTeam = false
		s.PieceFromOtherBoard = false
		s.PieceOrigPosition = false
		c.Add(*s)
	}

	// Add the black pieces
	for i := 0; i < 16; i++ {
		tmp := c.Square(i)

		tmp.PieceEmpty = false
		tmp.PieceTeam = false
		tmp.PieceFromOtherBoard = false
		tmp.PieceOrigPosition = true

		if i == 0 || i == 7 {
			tmp.PieceType = "Rook"
			tmp.PieceImage = "./Pieces/Rook_Black_60.png"
		} else if i == 1 || i == 6 {
			tmp.PieceType = "Knight"
			tmp.PieceImage = "./Pieces/Knight_Black_60.png"
		} else if i == 2 || i == 5 {
			tmp.PieceType = "Bishop"
			tmp.PieceImage = "./Pieces/Bishop_Black_60.png"
		} else if i == 4 {
			tmp.PieceType = "King"
			tmp.PieceImage = "./Pieces/King_Black_60.png"
		} else if i == 3 {
			tmp.PieceType = "Queen"
			tmp.PieceImage = "./Pieces/Queen_Black_60.png"
		} else {
			tmp.PieceType = "Pawn"
			tmp.PieceImage = "./Pieces/Pawn_Black_60.png"
		}
		c.SetSquare(i, *tmp)
	}
	// Add the white pieces
	for i := 48; i < 64; i++ {
		tmp := c.Square(i)

		tmp.PieceEmpty = false
		tmp.PieceTeam = true
		tmp.PieceFromOtherBoard = false
		tmp.PieceOrigPosition = true

		if i == 56 || i == 63 {
			tmp.PieceType = "Rook"
			tmp.PieceImage = "./Pieces/Rook_White_60.png"
		} else if i == 57 || i == 62 {
			tmp.PieceType = "Knight"
			tmp.PieceImage = "./Pieces/Knight_White_60.png"
		} else if i == 58 || i == 61 {
			tmp.PieceType = "Bishop"
			tmp.PieceImage = "./Pieces/Bishop_White_60.png"
		} else if i == 60 {
			tmp.PieceType = "King"
			tmp.PieceImage = "./Pieces/King_White_60.png"
		} else if i == 59 {
			tmp.PieceType = "Queen"
			tmp.PieceImage = "./Pieces/Queen_White_60.png"
		} else {
			tmp.PieceType = "Pawn"
			tmp.PieceImage = "./Pieces/Pawn_White_60.png"
		}
		c.SetSquare(i, *tmp)
	}
}

func (c *ChessBoard) SetSelectedSquare(index int) {
	c.SSquare = index
}

func (c *ChessBoard) SelectedSquare(index int) int {
	return c.SSquare
}

func (c *ChessBoard) SetNextSquare(index int) {
	c.NSquare = index
}

func (c *ChessBoard) NextSquare(index int) int {
	return c.NSquare
}

func (c *ChessBoard) Square(index int) *Square {
	return c.Board[index]
}

func (c *ChessBoard) SetSquare(index int, square Square) {
	c.Board[index] = &square
}

func (c *ChessBoard) Add(square Square) {
	c.Board = append(c.Board, &square)
	c.Len = len(c.Board)
}

func (c *ChessBoard) MovePiece(oldLoc int, newLoc int) {
	// TODO: Add a check for castling...
	oldl := c.Square(oldLoc)
	newl := c.Square(newLoc)

	fmt.Println("DEBUG: MovePiece")
	fmt.Println("DEBUG: TURN = ", c.CTurn, " old square -> INDEX = ", oldl.Index, " COLOR = ", oldl.Color, " PieceType = ", oldl.PieceType, " Image = ", oldl.PieceImage, " Is Empty = ", oldl.PieceEmpty, " Team = ", oldl.PieceTeam, " Other Board = ", oldl.PieceFromOtherBoard, " Orig Pos = ", oldl.PieceOrigPosition)
	fmt.Println("DEBUG: TURN = ", c.CTurn, " new square -> INDEX = ", newl.Index, " COLOR = ", newl.Color, " PieceType = ", newl.PieceType, " Image = ", newl.PieceImage, " Is Empty = ", newl.PieceEmpty, " Team = ", newl.PieceTeam, " Other Board = ", newl.PieceFromOtherBoard, " Orig Pos = ", newl.PieceOrigPosition)
	fmt.Println("DEBUG: FIN")

	if oldl.PieceTeam == c.CTurn && !oldl.PieceEmpty {

		// Check to see if capturing a piece.
		if !newl.PieceEmpty && newl.PieceTeam != c.CTurn {
			// Add the piece at the captured location to the other boards.
			// TODO: How to save these... An array/slice of square objects?
			fmt.Println("Captured a piece! A ", newl.PieceType, " of color ", newl.PieceTeam, " on this teams turn ", c.CTurn)
		}

		newl.PieceEmpty = false
		newl.PieceImage = oldl.PieceImage
		newl.PieceType = oldl.PieceType
		newl.PieceOrigPosition = false
		newl.PieceFromOtherBoard = oldl.PieceFromOtherBoard
		newl.PieceTeam = oldl.PieceTeam

		oldl.PieceEmpty = true
		oldl.PieceImage = ""
		oldl.PieceType = "Empty"
		oldl.PieceOrigPosition = false
		oldl.PieceFromOtherBoard = false

		c.SetSquare(oldLoc, *oldl)
		c.SetSquare(newLoc, *newl)

		qml.Changed(newl, &newl.PieceImage)
		qml.Changed(oldl, &oldl.PieceImage)
		c.CTurn = !c.CTurn
	} else {
		fmt.Println("Wrong turn for that piece.")
	}

	// If the s square is empty, do nothing.
	// If the s square is of the other teams color, do nothing.
	// If the new square is the same as the old square, do nothing.

	// If the new square is empty, just move the piece from s to tmp and update c.

	// If the new square has a piece of the same color, do nothing.

	// If the new square has a piece of the other color, remove that piece and place the piece from s there.

	// When that happens, needs to keep track of the capture piece to pass it to team mate.

}
