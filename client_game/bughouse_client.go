package main

import (
	"Pass_The_Queen/messenger"
	"Pass_The_Queen/mylib"
	"bufio"
	"fmt"
	"gopkg.in/qml.v1"
	"os"
	"strconv"
	"strings"
	"time"
)

var chatting *ChatMsg
var globalchatting *ChatMsg
var recv_msg string
var game *Game
var capturedPieces *CapturedPieces
var turn int
var chessBoard *ChessBoard
var my_name string
var my_room string
var in_room bool
var in_game bool
var rooms map[string]string
var game_team int
var game_color int
var game_start int
var room_members []string

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

func main() {
	if err := qml.Run(run); err != nil {
		fmt.Fprint(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	return Run()
}

func StartGame(room string, player string, board int, color int, team int) {
	fmt.Println("I am in the start game function.... good news.")
	game.Name = room
	game.PlayerID = player
	game.Board = board
	game.PlayerColor = color
	game.TeamPlayer = team

	turn = 0

	game.InGame = true
	qml.Changed(game, &game.InGame)
}

func EndGame(board int, team int, reason string) {
	game.InGame = false
	in_room = false
	qml.Changed(game, &game.InGame)

	chessBoard.Empty()
	capturedPieces.Empty()
	qml.Changed(chessBoard, &chessBoard.Len)
	qml.Changed(capturedPieces, &capturedPieces.Len)
	chessBoard.initialize()
	qml.Changed(chessBoard, &chessBoard.Board)
	qml.Changed(chessBoard, &chessBoard.Len)
	fmt.Println(chessBoard.Len)

	messenger.Msnger.Leave_local()
	messenger.Msnger.Join_global()
	chessBoard.Time = 420
	qml.Changed(chessBoard, &chessBoard.Time)
}

func Run() error {
	engine := qml.NewEngine()

	temp := &ChessBoard{}

	chessBoard = temp

	temp2 := &CapturedPieces{}
	capturedPieces = temp2

	temp3 := &Game{}
	game = temp3

	tmp4 := &ChatMsg{}
	chatting = tmp4

	tmp5 := &ChatMsg{}
	globalchatting = tmp5

	chessBoard.initialize()

	game.InGame = false

	engine.Context().SetVar("game", game)
	engine.Context().SetVar("chessBoard", chessBoard)
	engine.Context().SetVar("capturedPieces", capturedPieces)
	engine.Context().SetVar("chatting", chatting)
	engine.Context().SetVar("globalchatting", globalchatting)

	component, err := engine.LoadFile("../src/Pass_The_Queen/qml/Application.qml")

	if err != nil {
		return err
	}

	window := component.CreateWindow(nil)

	window.Show()

	game_start = 1
	my_name = os.Args[1]
	in_room = false
	in_game = false
	rooms = make(map[string]string)

	messenger.Msnger = messenger.NewMessenger(my_name)

	messenger.Msnger.Login()

	go process_messages()

	start_network(engine)

	window.Wait()

	return nil
}

func start_network(engine1 *qml.Engine) {
	fmt.Printf("\nChat:\n")
	reader := bufio.NewReader(os.Stdin)
	for {
		in, _ := reader.ReadString('\n')
		in = strings.Split(in, "\n")[0]
		command := strings.Split(in, " ")
		//Print rooms
		if in == "rooms" {
			for room_name, owner := range rooms {
				fmt.Printf("%v (%v)\n", room_name, owner)
			}
			fmt.Println("")
			//Create a room
		} else if command[0] == "create" {
			if in_room {
				fmt.Println("Already in a room")
			} else {
				room_name := strings.Split(in, fmt.Sprintf("%v ", command[0]))[1]
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
			//Join an existing room
		} else if command[0] == "join" {
			if in_room {
				fmt.Println("Already in a room")
			} else {
				room_name := strings.Split(in, fmt.Sprintf("%v ", command[0]))[1]
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
		} else if in == "start" {
			if !in_room {
				fmt.Println("Not in a room")
			} else if rooms[my_room] != fmt.Sprintf("%v:%v", my_name, messenger.Msnger.Port) {
				fmt.Println("Not the room owner")
				messenger.Msnger.Leave_global()
				messenger.Msnger.Join_local(room_members)
				StartGame(my_room, my_name, game.Board, game_color, game_team)
			} else {
				in_game = true
				messenger.Msnger.Send_game_server(my_room, mylib.START_GAME)
				messenger.Msnger.Send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, messenger.Msnger.Port), mylib.START_GAME)
				delete(rooms, my_room)
				messenger.Msnger.Leave_global()
				messenger.Msnger.Join_local(room_members)
				StartGame(my_room, my_name, game.Board, game_color, game_team)
			}
		} else if in == "leave" {
			if in_game {
				messenger.Msnger.Leave_local()
				messenger.Msnger.Join_global()
			} else if in_room {
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
			} else {
				messenger.Msnger.Leave_global()
			}
			//Print list of room members
		} else if in == "members" {
			for i := range room_members {
				fmt.Println(room_members[i])
			}
		} else if command[0] == "set_team" { // Change the current players team.
			if !in_room {
				fmt.Println("Not in a room!")
			} else {
				if strings.Count(in, " ") == 0 {
					fmt.Println("Must provide an integer value of either 1 or 2 for set_team")
				} else {
					test := strings.Split(in, fmt.Sprintf("%v ", command[0]))[1]
					temp, err := strconv.Atoi(test)
					if err != nil || (temp != 1 && temp != 2) {
						fmt.Println("Must provide an integer value of either 1 or 2 for set_team")
					} else {
						game_team = temp
					}
				}
			}
		} else if command[0] == "set_color" { // Change the current players color
			if !in_room {
				fmt.Println("Not in a room!")
			} else {
				if strings.Count(in, " ") == 0 {
					fmt.Println("Must provide an integer value of either 1 or 2 for set_color")
				} else {
					test := strings.Split(in, fmt.Sprintf("%v ", command[0]))[1]
					temp, err := strconv.Atoi(test)
					if err != nil || (temp != 1 && temp != 2) {
						fmt.Println("Must provide an integer value of either 1 or 2 for set_color")
					} else {
						game_color = temp
					}
				}
			}
		} else {
			messenger.Msnger.Send_message(in, mylib.CHAT_MESSAGE)
		}
	}
}

func (chat ChatMsg) SendChatMsg(data string) {

	if strings.HasPrefix(data, "L ") {
		tmp1 := strings.TrimLeft(data, "L ")

		chatting.Msg += "Me: " + tmp1 + "\n"
		qml.Changed(chatting, &chatting.Msg)

	} else if strings.HasPrefix(data, "G ") {
		tmp2 := strings.TrimLeft(data, "G ")

		globalchatting.Msg += "Me: " + tmp2 + "\n"
		qml.Changed(globalchatting, &globalchatting.Msg)
	}

	fmt.Println("******************************************")

	fmt.Println("sending:" + data)

	fmt.Println("******************************************")
	messenger.Msnger.Send_message(data, mylib.CHAT_MESSAGE)
}

func process_messages() {
	var content string

	for {
		msg := messenger.Msnger.Receive_message()
		if msg.Type == mylib.NONE {
			time.Sleep(1000000000)
			continue
		}
		content = msg.Content
		if msg.Type == mylib.CHAT_MESSAGE {
			mesge := msg.Content

			if strings.HasPrefix(mesge, "L ") {
				content = fmt.Sprintf("%v says: %v", msg.Orig_source, strings.TrimLeft(mesge, "L "))
				fmt.Println("sadf     " + content)
				chatting.Msg += content + "\n"
				qml.Changed(chatting, &chatting.Msg)

			} else if strings.HasPrefix(mesge, "G ") {
				content = fmt.Sprintf("%v says: %v", msg.Orig_source, strings.TrimLeft(mesge, "G "))
				fmt.Println(content)
				globalchatting.Msg += content + "\n"
				qml.Changed(globalchatting, &globalchatting.Msg)
			}

		} else if msg.Type == mylib.CREATE_ROOM {
			decoded := strings.Split(content, ":")
			rooms[decoded[0]] = fmt.Sprintf("%v:%v:%v", decoded[1], decoded[2], decoded[3])
		} else if msg.Type == mylib.JOIN_ROOM {
			decoded := strings.Split(content, ":")
			if my_room == decoded[0] {
				room_members = append(room_members, fmt.Sprintf("%v:%v:%v", decoded[1], decoded[2], decoded[3]))
			}
		} else if msg.Type == mylib.START_GAME {
			decoded := strings.Split(content, ":")
			delete(rooms, decoded[0])
			if my_room == decoded[0] {
				in_game = true
			}
		} else if msg.Type == mylib.DELETE_ROOM {
			decoded := strings.Split(content, ":")
			delete(rooms, decoded[0])
			if my_room == decoded[0] {
				my_room = ""
				in_room = false
				room_members = make([]string, 0, 0)
			}
		} else if msg.Type == mylib.LEAVE_ROOM {
			content = msg.Content
			decoded := strings.Split(content, ":")
			if my_room == decoded[0] {
				for i := range room_members {
					if room_members[i] == fmt.Sprintf("%v:%v:%v", decoded[1], decoded[2], decoded[3]) {
						room_members = append(room_members[:i], room_members[i+1:]...)
						break
					}
				}
			}
		} else if msg.Type == mylib.MOVE {
			decoded := strings.Split(content, ":")
			// Should be board_num, player_team, player_color, turn, origLoc, newLoc, capture_pieceString
			board_num, _ := strconv.Atoi(decoded[0])
			team, _ := strconv.Atoi(decoded[1])
			color, _ := strconv.Atoi(decoded[2])
			turn, _ := strconv.Atoi(decoded[3])
			origLoc, _ := strconv.Atoi(decoded[4])
			newLoc, _ := strconv.Atoi(decoded[5])
			capturedI := decoded[6]
			capturedT := decoded[7]
			capturedU, _ := strconv.Atoi(decoded[8])
			UpdateFromOpponent(board_num, team, color, turn, origLoc, newLoc, capturedI, capturedT, capturedU)
		} else if msg.Type == mylib.PLACE {
			decoded := strings.Split(content, ":")
			// Should be board_num, player_team, player_color, turn, origLoc, newLoc, capture_pieceString
			board_num, _ := strconv.Atoi(decoded[0])
			team, _ := strconv.Atoi(decoded[1])
			color, _ := strconv.Atoi(decoded[2])
			turn, _ := strconv.Atoi(decoded[3])
			loc, _ := strconv.Atoi(decoded[4])
			piece, _ := strconv.Atoi(decoded[5])
			pieceImage := decoded[6]
			pieceType := decoded[7]
			pieceTeam, _ := strconv.Atoi(decoded[8])
			UpdatePlace(board_num, team, color, turn, loc, piece, pieceImage, pieceType, pieceTeam)
		} else if msg.Type == mylib.GAMEOVER {
			decoded := strings.Split(content, ":")
			board_num, _ := strconv.Atoi(decoded[0])
			team, _ := strconv.Atoi(decoded[1])
			cause := decoded[2]
			EndGame(board_num, team, cause)
			if cause == "King" {
				// The king was captured
			} else {
				// Somebody ran out of time
			}
		}
		msg.Type = mylib.NONE
	}
}

func (g *Game) ChangeTeam() string {
	if game_team == 1 {
		game_team = 2
		return "Team: 2"
	} else {
		game_team = 1
		return "Team: 1"
	}
}

func (g *Game) ChangeColor() string {
	if game_color == 1 {
		game_color = 2
		return "Color: Black"
	} else {
		game_color = 1
		return "Color: White"
	}
}

func (g *Game) ListGames() string {
	temp := ""
	for room_name, owner := range rooms {
		fmt.Printf("%v (%v)\n", room_name, owner)
		temp = temp + "\n" + room_name + " " + owner
	}
	fmt.Println("")
	return temp
}

func (g *Game) CheckGames(input string) bool {
	for room_name, _ := range rooms {
		if room_name == input {
			return true
		}
	}
	return false // No matches found
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
	if in_room {
		fmt.Println("Already in a room")
	} else {
		if rooms[room_name] == "" {
			fmt.Println("Error: Room with such name does not exist")
		} else {
			game_color = 2
			game_team = 2
			messenger.Msnger.Send_message(fmt.Sprintf("%v:%v:%v", room_name, my_name, messenger.Msnger.Address), mylib.JOIN_ROOM)
			in_room = true
			my_room = room_name
		}
	}
}

func (g *Game) LeaveRoom() {
	if in_game {
		messenger.Msnger.Leave_local()
		messenger.Msnger.Join_global()
	} else if in_room {
		//Delete room if room owner
		if rooms[my_room] == fmt.Sprintf("%v:%v", my_name, messenger.Msnger.Address) {
			messenger.Msnger.Send_game_server(my_room, mylib.DELETE_ROOM)
			messenger.Msnger.Send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, messenger.Msnger.Address), mylib.DELETE_ROOM)
			delete(rooms, my_room)
		} else {
			messenger.Msnger.Send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, messenger.Msnger.Address), mylib.LEAVE_ROOM)
		}
		in_room = false
		my_room = ""
		room_members = make([]string, 0, 0)
	}
}

func (g *Game) CreateRoom(room_name string) {
	if in_room {
		fmt.Println("Already in a room")
	} else {
		game_color = 1
		game_team = 1
		messenger.Msnger.Send_game_server(fmt.Sprintf("%v:%v:%v", room_name, my_name, messenger.Msnger.Address), mylib.CREATE_ROOM)
		msg := messenger.Msnger.Receive_game_server()
		if msg.Type == mylib.ACK {
			fmt.Printf("Creating room: %v\n", room_name)
			rooms[room_name] = fmt.Sprintf("%v:%v", my_name, messenger.Msnger.Address)
			messenger.Msnger.Send_message(fmt.Sprintf("%v:%v:%v", room_name, my_name, messenger.Msnger.Address), mylib.CREATE_ROOM)
			in_room = true
			my_room = room_name
		} else {
			fmt.Println("Error: Room name already taken")
		}
	}
}

func (g *Game) StartRoom(host int, boardNum int) {
	if !in_room {
		fmt.Println("Not in a room")
	} else if rooms[my_room] != fmt.Sprintf("%v:%v", my_name, messenger.Msnger.Address) {
		fmt.Println("Not the room owner")
		delete(rooms, my_room)
		messenger.Msnger.Leave_global()
		messenger.Msnger.Join_local(room_members)
		StartGame(my_room, my_name, boardNum, game_color, game_team)
	} else {
		messenger.Msnger.Send_game_server(my_room, mylib.START_GAME)
		messenger.Msnger.Send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, messenger.Msnger.Address), mylib.START_GAME)
		delete(rooms, my_room)
		messenger.Msnger.Leave_global()
		messenger.Msnger.Join_local(room_members)
		in_game = true
		StartGame(my_room, my_name, boardNum, game_color, game_team)
	}
}

func (c *ChessBoard) Timer() {
	if c.Time > 0 && game.InGame {
		if game.PlayerColor-1 == (turn % 2) {
			c.Time = c.Time - 1
			qml.Changed(c, &c.Time)
		}
	} else if game.InGame {
		temp := fmt.Sprintf("%v:%v:%v", game.Board, game.TeamPlayer, "Time")
		messenger.Msnger.Send_message(temp, mylib.GAMEOVER)
		EndGame(game.Board, game.TeamPlayer, "Time")
	}
}

func UpdateFromOpponent(board int, team int, color int, turnO int, origL int, newL int, capturedI string, capturedT string, capturedU int) {
	if board != game.Board && capturedI != "" && team == game.TeamPlayer { // Ignore this update then. Or could actually get the captured piece thingy.
		captured := new(CapturedPiece)
		captured.Image = capturedI
		captured.Type = capturedT
		captured.TeamPiece = capturedU
		capturedPieces.Add(*captured)
		qml.Changed(capturedPieces, &capturedPieces.Len)
	} else if board == game.Board && team != game.TeamPlayer {
		if turn != turnO {
			// Incompatible state. Different amounts of turns
			fmt.Println("Error with the number of turns.")
		} else {
			chessBoard.Update(origL, newL)
			turn++
		}
	} else { //if board != game.Board && team == game.TeamPlayer {
		// Can ignore, is a move not from your board
		// Can ignore if it
		//fmt.Println("DEBUG: An update from opponent or team mate that is not on this board has happend. Handled correctly!.")
		return
	}
}

func UpdatePlace(board int, team int, color int, turn0 int, loc int, p int, piecePImage string, piecePType string, piecePTeamPiece int) {
	if board != game.Board && team == game.TeamPlayer {
		// Remove the piece from captured pieces
		//capturedPieces.Pieces = append(capturedPieces.Pieces[:piece], capturedPieces.Pieces[piece+1:])
		tempP := &CapturedPieces{}
		tempPieces := tempP.Pieces
		for i := 0; i < capturedPieces.Len; i++ {
			if i != p {
				tempPieces = append(tempPieces, capturedPieces.Piece(i))
			}
		}

		capturedPieces.Pieces = tempPieces
		capturedPieces.Len = len(capturedPieces.Pieces)
		qml.Changed(capturedPieces, &capturedPieces.Len)
	} else if board == game.Board && team != game.TeamPlayer {
		temp := chessBoard.Square(loc)
		temp.Empty = false
		temp.FromOtherBoard = true
		temp.Image = piecePImage
		temp.OrigPosition = false
		temp.Type = piecePType
		temp.TeamPiece = piecePTeamPiece

		chessBoard.Board[loc].Empty = false
		chessBoard.Board[loc].FromOtherBoard = true
		chessBoard.Board[loc].Image = piecePImage
		chessBoard.Board[loc].OrigPosition = false
		chessBoard.Board[loc].Type = piecePType
		chessBoard.Board[loc].TeamPiece = piecePTeamPiece

		qml.Changed(chessBoard.Board[loc], &chessBoard.Board[loc].Image)
		turn++
	}
}

func (c *ChessBoard) Update(origLoc int, newLoc int) {
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

// In move piece, change the turn, and then send the state accross the network.
func (c *ChessBoard) MovePiece(origLoc int, newLoc int) {
	if game.PlayerColor == (turn%2)+1 {
		// Get the square values at the index locations
		origS := c.Square(origLoc)
		newS := c.Square(newLoc)

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
				temp := fmt.Sprintf("%v:%v:%v:%v:%v:%v:%v:%v:%v", game.Board, game.TeamPlayer, game.PlayerColor, turn, origLoc, newLoc, "", "", 0)
				//fmt.Println("DEBUG: state message to be sent accross network = ", temp)
				// client.Send_move(temp)
				// Find some way of having main send a message of temp to the other nodes.
				turn++

				// Add a send message to opponent about the move...
				messenger.Msnger.Send_message(temp, mylib.MOVE)
			} else {
				// For enemy piece, record the captured piece, move the piece there and end the turn.
				// Involves sending info over the network and changing the turn value to the other team.
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

				qml.Changed(capturedPieces, &capturedPieces.Len)
				// Send message to the opponent before incrementing turn
				// Send them the current player, the turn value, and the 2 index's.
				// In this case, send a caputed value equal to [Color][Piece] as a string or something.

				// Add a send captured pieces to partner
				temp := fmt.Sprintf("%v:%v:%v:%v:%v:%v:%v:%v:%v", game.Board, game.TeamPlayer, game.PlayerColor, turn, origLoc, newLoc, captured.Image, captured.Type, captured.TeamPiece)

				turn++

				messenger.Msnger.Send_message(temp, mylib.MOVE)

				if captured.Type == "King" {
					temp := fmt.Sprintf("%v:%v:%v", game.Board, game.TeamPlayer, "King")
					messenger.Msnger.Send_message(temp, mylib.GAMEOVER)
					EndGame(game.Board, game.TeamPlayer, "King")
				}
			}
		} else {
			fmt.Println("Cannot move pieces that are not your own!")
		}
	} else {
		fmt.Println("Can't move a piece when it is not your turn!")
	}
}

func (c *ChessBoard) PlacePiece(loc int, p int) {
	if game.PlayerColor == (turn%2)+1 {
		piece := capturedPieces.Piece(p)
		temp := c.Square(loc)
		temp.Empty = false
		temp.FromOtherBoard = true
		temp.Image = piece.Image
		temp.OrigPosition = false
		temp.Type = piece.Type
		temp.TeamPiece = piece.TeamPiece
		qml.Changed(c.Board[loc], &c.Board[loc].Image)

		temp2 := fmt.Sprintf("%v:%v:%v:%v:%v:%v:%v:%v:%v", game.Board, game.TeamPlayer, game.PlayerColor, turn, loc, p, piece.Image, piece.Type, piece.TeamPiece)

		turn++

		messenger.Msnger.Send_message(temp2, mylib.PLACE)

		tempP := &CapturedPieces{}
		tempPieces := tempP.Pieces
		for i := 0; i < capturedPieces.Len; i++ {
			if i != p {
				tempPieces = append(tempPieces, capturedPieces.Piece(i))
			}
		}
		capturedPieces.Pieces = tempPieces
		capturedPieces.Len = len(capturedPieces.Pieces)
		qml.Changed(capturedPieces, &capturedPieces.Len)

	}
}

func (cp *CapturedPieces) Add(p CapturedPiece) {
	cp.Pieces = append(cp.Pieces, &p)
	cp.Len = len(cp.Pieces)
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

func (c *ChessBoard) reinitialize() {

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

	for i := 16; i < 48; i++ {
		tmp := c.Square(i)

		tmp.Type = "EMPTY"
		tmp.Image = ""
		tmp.Empty = true
		tmp.TeamPiece = EMPTY
		tmp.FromOtherBoard = false
		tmp.OrigPosition = true

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

func (c *ChessBoard) Empty() {
	c.Board = nil
	c.Len = 0
}

func (cp *CapturedPieces) Empty() {
	cp.Pieces = nil
	cp.Len = 0
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
	c.reinitialize()
}

// Create a function to desconstruct the board.
