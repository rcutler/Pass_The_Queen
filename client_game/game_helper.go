package client_game

import (
	"Pass_The_Queen/messenger"
	"Pass_The_Queen/mylib"
	//"bufio"
	"fmt"
	"gopkg.in/qml.v1"
	//"os"
	//"strconv"
	//"strings"
	//"time"
)

func (g *Game) ChangeTeam() {
	if game_team == 1 {
		game_team = 2
	} else {
		game_team = 1
	}
	fmt.Println(game_team)

}

func (g *Game) ChangeColor() {
	if game_color == 1 {
		game_color = 2
	} else {
		game_color = 1
	}
	fmt.Println(game_color)
}

func (g *Game) ListGames() string {
	temp := ""
	for room_name, owner := range rooms {
		fmt.Printf("%v (%v)\n", room_name, owner)
		//return room_name + " " + owner
		//return fmt.Sprint("%v (%v)\n", room_name, owner)
		temp = temp + "\n" + room_name + " " + owner
	}
	fmt.Println("")
	return temp
}

func (g *Game) Members() string {
	temp := ""
	for i := range room_members {
		fmt.Println(room_members[i])
		temp = temp + "\n" + room_members[i]
	}
	return temp
}

func (g *Game) JoinRoom(room_name string) {
	//"join" {
	if in_room {
		fmt.Println("Already in a room")
	} else {
		//room_name := strings.Split(in, fmt.Sprintf("%v ", command[0]))[1]
		if rooms[room_name] == "" {
			fmt.Println("Error: Room with such name does not exist")
		} else {
			game_color = 2
			game_team = 2
			messenger.Msnger.Send_message(fmt.Sprintf("%v:%v:%v", room_name, my_name, messenger.Msnger.Port), mylib.JOIN_ROOM)
			in_room = true
			my_room = room_name
		}
	}
}

func (g *Game) LeaveRoom() {
	//"leave" {
	if in_room {
		//Delete room if room owner
		if rooms[my_room] == fmt.Sprintf("%v:%v", my_name, messenger.Msnger.Port) {
			messenger.Msnger.Send_game_server(my_room, mylib.DELETE_ROOM)
			messenger.Msnger.Send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, messenger.Msnger.Port), mylib.DELETE_ROOM)
			delete(rooms, my_room)
		} else {
			messenger.Msnger.Send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, messenger.Msnger.Port), mylib.LEAVE_ROOM)
		}
		in_room = false
		my_room = ""
		room_members = make([]string, 0, 0)
	}
}

func (g *Game) CreateRoom(room_name string) {
	//"create" {
	if in_room {
		fmt.Println("Already in a room")
	} else {
		//room_name := strings.Split(in, fmt.Sprintf("%v ", command[0]))[1]
		fmt.Println("DEBUG: room_name = ", room_name)
		game_color = 1
		game_team = 1
		messenger.Msnger.Send_game_server(fmt.Sprintf("%v:%v:%v", room_name, my_name, messenger.Msnger.Port), mylib.CREATE_ROOM)
		msg := messenger.Msnger.Receive_game_server()
		if msg.Type == mylib.ACK {
			fmt.Printf("Creating room: %v\n", room_name)
			rooms[room_name] = fmt.Sprintf("%v:%v", my_name, messenger.Msnger.Port)
			messenger.Msnger.Send_message(fmt.Sprintf("%v:%v:%v", room_name, my_name, messenger.Msnger.Port), mylib.CREATE_ROOM)
			in_room = true
			my_room = room_name
		} else {
			fmt.Println("Error: Room name already taken")
		}
	}
}

func (g *Game) StartRoom(host int) {
	if !in_room {
		fmt.Println("Not in a room")
	} else {
		if host == 1 {
			messenger.Msnger.Send_game_server(my_room, mylib.START_GAME)
			messenger.Msnger.Send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, messenger.Msnger.Port), mylib.START_GAME)
			delete(rooms, my_room)
			fmt.Println("SHould only be here on A")
		}
		in_game = true
		// Need at stuff to set up teams and color
		//fmt.Println("DEBUG start command from owner: ", my_room, " ", my_name, " ", 1, " ", game_color, " ", game_team)
		//start_local_chat(my_room, my_name, 1, game_color, game_team)
		StartGame(my_room, my_name, 1, game_color, game_team)
	}
}

func (c *ChessBoard) Timer() {
	if game.TeamPlayer == (turn%2)+1 {
		c.Time = c.Time - 1
		qml.Changed(c, &c.Time)
	}
}

func UpdateFromOpponent(board int, team int, color int, turnO int, origL int, newL int, captured string) {
	// Do some error checking here to see if valid piece and what not....
	fmt.Println(board, " ", game.Board)
	fmt.Println(captured)
	fmt.Println(team, " ", game.TeamPlayer)
	if board != game.Board && captured != "" && team == game.TeamPlayer { // Ignore this update then. Or could actually get the captured piece thingy.
		// Update your captured pieces to add this new piece.
		fmt.Println("1")
	} else if board == game.Board && team != game.TeamPlayer {
		// The turn value should be the same
		fmt.Println("2")
		if turn != turnO {
			// Incompatible state. Different amounts of turns
			fmt.Println("Error with the number of turns.")
			fmt.Println("3")
		} else {
			chessBoard.Update(origL, newL)
			fmt.Println("4")
			turn++
		}
	} else { //if board != game.Board && team == game.TeamPlayer {
		// Can ignore, is a move not from your board
		// Can ignore if it
		fmt.Println("5")
		//fmt.Println("DEBUG: An update from opponent or team mate that is not on this board has happend. Handled correctly!.")
		return
	}
}

func (c *ChessBoard) Update(origLoc int, newLoc int) {
	//fmt.Println("DEBUG: origLoc chessboard = ", chessBoard.Board[origLoc])
	origS := c.Square(origLoc)
	newS := c.Square(newLoc)
	if origS.TeamPiece == newS.TeamPiece {
		fmt.Println("Cannot move a piece on top of a piece of the same color")
		fmt.Println("An exception would be castling, but that is not implemented yet.")
	} else {
		// For empty piece, move the piece there, and end the turn.
		// Involves sending info over the network and changing the turn value to the other team
		c.Board[newLoc].Empty = false
		c.Board[newLoc].TeamPiece = origS.TeamPiece
		c.Board[newLoc].Type = origS.Type
		c.Board[newLoc].OrigPosition = false
		c.Board[newLoc].FromOtherBoard = origS.FromOtherBoard
		c.Board[newLoc].Image = origS.Image

		newS.Empty = false
		newS.TeamPiece = origS.TeamPiece
		newS.Type = origS.Type
		newS.OrigPosition = false
		newS.FromOtherBoard = origS.FromOtherBoard
		newS.Image = origS.Image

		c.Board[origLoc].Image = ""
		c.Board[origLoc].Empty = true
		c.Board[origLoc].TeamPiece = EMPTY
		c.Board[origLoc].Type = "EMPTY"
		c.Board[origLoc].FromOtherBoard = false
		c.Board[origLoc].OrigPosition = false

		origS.Empty = true
		origS.TeamPiece = EMPTY
		origS.Type = "EMPTY"
		origS.OrigPosition = false
		origS.FromOtherBoard = false
		origS.Image = ""

		qml.Changed(c.Board[origLoc], &c.Board[origLoc].Image)
		qml.Changed(c.Board[newLoc], &c.Board[newLoc].Image)
	}
}

// Add a function for Move Piece
// In move piece, change the turn, and then send the state accross the network.
func (c *ChessBoard) MovePiece(origLoc int, newLoc int) {
	if game.PlayerColor == (turn%2)+1 {
		// Get the square values at the index locations
		origS := c.Square(origLoc)
		newS := c.Square(newLoc)

		// Add some debug checking here.
		//fmt.Println("DEBUG: MovePiece")
		//fmt.Println("DEBUG: Game turn is = ", turn)
		//fmt.Println("DEBUG: Current players color is (1 for white, 2 for black) = ", game.PlayerColor)
		//fmt.Println("DEBUG: The original Squares color is (1 for white, 2 for black) = ", origS.TeamPiece)

		// Check that the origS color is the same as the local player's color.
		if origS.TeamPiece == (turn%2)+1 {
			// Check to see the the newLoc square is empty, same color piece, other team piece
			// For same team piece, move is not valid, exit
			if origS.TeamPiece == newS.TeamPiece {
				fmt.Println("Cannot move a piece on top of a piece of the same color")
				fmt.Println("An exception would be castling, but that is not implemented yet.")
			} else if newS.Empty {
				// For empty piece, move the piece there, and end the turn.
				// Involves sending info over the network and changing the turn value to the other team
				c.Board[newLoc].Empty = false
				c.Board[newLoc].TeamPiece = origS.TeamPiece
				c.Board[newLoc].Type = origS.Type
				c.Board[newLoc].OrigPosition = false
				c.Board[newLoc].FromOtherBoard = origS.FromOtherBoard
				c.Board[newLoc].Image = origS.Image

				newS.Empty = false
				newS.TeamPiece = origS.TeamPiece
				newS.Type = origS.Type
				newS.OrigPosition = false
				newS.FromOtherBoard = origS.FromOtherBoard
				newS.Image = origS.Image

				c.Board[origLoc].Image = ""
				c.Board[origLoc].Empty = true
				c.Board[origLoc].TeamPiece = EMPTY
				c.Board[origLoc].Type = "EMPTY"
				c.Board[origLoc].FromOtherBoard = false
				c.Board[origLoc].OrigPosition = false

				origS.Empty = true
				origS.TeamPiece = EMPTY
				origS.Type = "EMPTY"
				origS.OrigPosition = false
				origS.FromOtherBoard = false
				origS.Image = ""

				qml.Changed(c.Board[origLoc], &c.Board[origLoc].Image)
				qml.Changed(c.Board[newLoc], &c.Board[newLoc].Image)

				// Send message to the opponent before incrementing turn
				// Send them the current player, the turn value, and the 2 index's.
				// In this case, send a caputed value to nil or 0.
				// Also need to send board value.
				// Message content as follows:
				// board_num:team:current_player:turn:origLoc:newLoc: :
				temp := fmt.Sprintf("%v:%v:%v:%v:%v:%v:%v", game.Board, game.TeamPlayer, game.PlayerColor, turn, origLoc, newLoc, "")
				fmt.Println("DEBUG: state message to be sent accross network = ", temp)
				// client.Send_move(temp)
				// Find some way of having main send a message of temp to the other nodes.
				turn++

				// Add a send message to opponent about the move...
				messenger.Msnger.Send_message(temp, mylib.MOVE)
			} else {
				// For enemy piece, record the captured piece, move the piece there and end the turn.
				// Involves sending info over the network and changing the turn value to the other team.
				fmt.Println("Need to finish this....")
				// Record the captured Piece and add it to the array of captured pieces...
				// Create a gridview to view all of the caputred pieces....
				captured := new(CapturedPiece)
				captured.Image = newS.Image
				captured.Type = newS.Type
				captured.TeamPiece = newS.TeamPiece
				capturedPieces.Add(*captured)

				c.Board[newLoc].Empty = false
				c.Board[newLoc].TeamPiece = origS.TeamPiece
				c.Board[newLoc].Type = origS.Type
				c.Board[newLoc].OrigPosition = false
				c.Board[newLoc].FromOtherBoard = origS.FromOtherBoard
				c.Board[newLoc].Image = origS.Image

				newS.Empty = false
				newS.TeamPiece = origS.TeamPiece
				newS.Type = origS.Type
				newS.OrigPosition = false
				newS.FromOtherBoard = origS.FromOtherBoard
				newS.Image = origS.Image

				c.Board[origLoc].Image = ""
				c.Board[origLoc].Empty = true
				c.Board[origLoc].TeamPiece = EMPTY
				c.Board[origLoc].Type = "EMPTY"
				c.Board[origLoc].FromOtherBoard = false
				c.Board[origLoc].OrigPosition = false

				origS.Empty = true
				origS.TeamPiece = EMPTY
				origS.Type = "EMPTY"
				origS.OrigPosition = false
				origS.FromOtherBoard = false
				origS.Image = ""

				qml.Changed(c.Board[origLoc], &c.Board[origLoc].Image)
				qml.Changed(c.Board[newLoc], &c.Board[newLoc].Image)

				//qml.Changed(capturedPieces.Pieces[capturedPieces.Len-1], &capturedPieces.Pieces[capturedPieces.Len-1].Image)

				// Send message to the opponent before incrementing turn
				// Send them the current player, the turn value, and the 2 index's.
				// In this case, send a caputed value equal to [Color][Piece] as a string or something.

				turn++
				// Add a send captured pieces to partner
				cp := fmt.Sprintf("%v-%v", captured.Type, captured.Image)
				temp := fmt.Sprintf("%v:%v:%v:%v:%v:%v:%v", game.Board, game.TeamPlayer, game.PlayerColor, turn, origLoc, newLoc, cp)
				fmt.Println("DEBUG: state message to be sent accross network = ", temp)

			}
		} else {
			fmt.Println("Cannot move pieces that are not your own!")
		}
	} else {
		fmt.Println("Can't move a piece when it is not your turn!")
		//fmt.Println("For testing, increment turn...")
		//turn++
	}
}

func (cp *CapturedPieces) Add(p CapturedPiece) {
	cp.Pieces = append(cp.Pieces, &p)
	cp.Len = len(cp.Pieces)
	fmt.Println("Length of captured pieces = ", cp.Len)
	// Add a qml changed part here for the viewing of captured pieces
	//qml.Changed(cp.Pieces[cp.Len-1], &cp.Pieces[cp.Len-1].Image)
	//temp := cp.Len - 1
	//fmt.Println(cp.Pieces[temp])
}

func (cp *CapturedPieces) Piece(index int) *CapturedPiece {
	return cp.Pieces[index]
}

func (c *ChessBoard) SetSquare(index int, square Square) {
	c.Board[index] = &square
}

func (c *ChessBoard) Square(index int) *Square {
	return c.Board[index]
}

func (c *ChessBoard) Add(square Square) {
	c.Board = append(c.Board, &square)
	c.Len = len(c.Board)
}

func (c *ChessBoard) initialize() {
	fmt.Println(game.Name)

	c.CurrentSquare = EMPTY
	c.NextSquare = EMPTY
	c.Turn = WHITE
	row := 0
	for i := 0; i < 64; i++ {
		s := new(Square)
		if (i % 8) == 0 {
			row++
		}
		if row%2 == 0 && i%2 == 0 {
			s.Color = COLOR_WHITE
		} else if row%2 == 0 && i%2 == 1 {
			s.Color = COLOR_BLACK
		} else if row%2 == 1 && i%2 == 0 {
			s.Color = COLOR_BLACK
		} else {
			s.Color = COLOR_WHITE
		}

		s.Index = i
		s.Type = "EMPTY"
		s.Image = ""
		s.Empty = true
		s.TeamPiece = EMPTY
		s.FromOtherBoard = false
		s.OrigPosition = true
		c.Add(*s)
	}

	// Add the black pieces
	for i := 0; i < 16; i++ {
		tmp := c.Square(i)
		tmp.Empty = false
		tmp.TeamPiece = BLACK
		tmp.FromOtherBoard = false
		tmp.OrigPosition = true
		if i == 0 || i == 7 {
			tmp.Type = "Rook"
			tmp.Image = BLACKROOK
		} else if i == 1 || i == 6 {
			tmp.Type = "Knight"
			tmp.Image = BLACKKNIGHT
		} else if i == 2 || i == 5 {
			tmp.Type = "Bishop"
			tmp.Image = BLACKBISHOP
		} else if i == 4 {
			tmp.Type = "King"
			tmp.Image = BLACKKING
		} else if i == 3 {
			tmp.Type = "Queen"
			tmp.Image = BLACKQUEEN
		} else {
			tmp.Type = "Pawn"
			tmp.Image = BLACKPAWN
		}
		c.SetSquare(i, *tmp)

	}
	// Add the white pieces
	for i := 48; i < 64; i++ {
		tmp := c.Square(i)
		tmp.Empty = false
		tmp.TeamPiece = WHITE
		tmp.FromOtherBoard = false
		tmp.OrigPosition = true
		if i == 56 || i == 63 {
			tmp.Type = "Rook"
			tmp.Image = WHITEROOK
		} else if i == 57 || i == 62 {
			tmp.Type = "Knight"
			tmp.Image = WHITEKNIGHT
		} else if i == 58 || i == 61 {
			tmp.Type = "Bishop"
			tmp.Image = WHITEBISHOP
		} else if i == 60 {
			tmp.Type = "King"
			tmp.Image = WHITEKING
		} else if i == 59 {
			tmp.Type = "Queen"
			tmp.Image = WHITEQUEEN
		} else {
			tmp.Type = "Pawn"
			tmp.Image = WHITEPAWN
		}
		c.SetSquare(i, *tmp)

	}
	c.Time = 420
}

// Create a function to desconstruct the board.
