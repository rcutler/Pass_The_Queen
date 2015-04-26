package client_game

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

	//"encoding/json"
)
//chatting message
type ChatMsg struct{ 
	Msg string
	//ChatTextChanged bool
}
var chatting *ChatMsg
var recv_msg string

var game *Game
var capturedPieces CapturedPieces
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

//Local network
var room_members []string

func Run() error {

	//tmp := &qml.Engine{}
	//engine = tmp
	engine := qml.NewEngine()

	temp := &ChessBoard{}

	chessBoard = temp

	temp2 := CapturedPieces{}
	capturedPieces = temp2

	temp3 := &Game{}
	game = temp3

	tmp4:=&ChatMsg{}
	chatting=tmp4

	chessBoard.initialize()

	engine.Context().SetVar("game", game)
	engine.Context().SetVar("chessBoard", chessBoard)
	//chat room variable
	engine.Context().SetVar("chatting",chatting)
	
	
	component, err := engine.LoadFile("../src/Pass_The_Queen/qml/Application.qml")

	if err != nil {
		return err
	}

	window := component.CreateWindow(nil)

	window.Show()


	// chatContent:=window.Root().ObjectByName("chatContent")

	// type Jason struct{
	// 	content string
	// }
	// c:=Jason{"lalala"}
	// fmt.Println(c)
	// chatContent.Call("append","content: lalala")


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

func (chat ChatMsg)SendChatMsg(data string){
		fmt.Println("******************************************");
		fmt.Println("sending: "+data)
		// chatting.Msg= "is it success?"
		// qml.Changed(chatting,&chatting.Msg)
		fmt.Println("******************************************");
		messenger.Msnger.Send_message(data, mylib.CHAT_MESSAGE)
}



func StartGame(room string, player string, board int, color int, team int) { //, engine1 *qml.Engine) {
	fmt.Println("I am in the start game function.... good news.")
	game.Name = room
	game.PlayerID = player
	game.Board = board
	game.PlayerColor = color
	game.TeamPlayer = team

	fmt.Println("IN StartGame: ", game.Board)
	fmt.Println("IN StartGame: ", game.TeamPlayer)

	turn = 0

	/*	temp := &ChessBoard{}

		chessBoard = temp

		temp2 := CapturedPieces{}
		capturedPieces = temp2*/

	//capturedPieces := &CapturedPieces{}
	//	chessBoard.initialize()

	t3 := new(CapturedPiece)
	t3.Image = BLACKROOK
	t3.TeamPiece = 1
	t3.Type = "ROOK"

	capturedPieces.Add(*t3)

	// Check that everything in the chessboard is initialized properly.
	for i := 0; i < 64; i++ {
		s := chessBoard.Square(i)
		fmt.Println("Index: ", s.Index, "  Color: ", s.Color, " Piece Type: ", s.Type, " Piece image loc: ", s.Image, " Piece color: ", s.TeamPiece)
	}

	//fmt.Println("DEBUG: HERE???? 1")

	//engine1.Context().SetVar("game", game)
	//engine1.Context().SetVar("chessBoard", chessBoard)
	//engine1.Context().SetVar("capturedPieces", capturedPieces)

	//fmt.Println("DEBUG: HERE???? 2")
}

// Create a function for ending the game...

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
			} else {
				in_game = true
				messenger.Msnger.Send_game_server(my_room, mylib.START_GAME)
				messenger.Msnger.Send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, messenger.Msnger.Port), mylib.START_GAME)
				delete(rooms, my_room)
				// Need at stuff to set up teams and color
				//fmt.Println("DEBUG start command from owner: ", my_room, " ", my_name, " ", 1, " ", game_color, " ", game_team)
				//start_local_chat(my_room, my_name, 1, game_color, game_team)
				StartGame(my_room, my_name, 1, game_color, game_team)
			}
			//Leave a room
		} else if in == "start_guest" {
			if !in_room {
				fmt.Println("Not in a room")
			} else if rooms[my_room] == fmt.Sprintf("%v:%v", my_name, messenger.Msnger.Port) {
				fmt.Println("The room owner, use 'start' instead")
			} else {
				in_game = true
				StartGame(my_room, my_name, 1, game_color, game_team)
			}
		} else if in == "leave" {
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
			//Print list of room members
		} else if in == "members" {
			for i := range room_members {
				fmt.Println(room_members[i])
			}
		} else if command[0] == "set_team" { // Change the current players team.
			if !in_room {
				fmt.Println("Not in a room!")
				//client_game.StartGame("a", "b", 1, 2, 2)
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
		} else if command[0] == "test_move" {
			messenger.Msnger.Send_message("1:2:2:0:51:3:Queen-../pieces/Queen_Black_60.png", mylib.MOVE)
		} else {
			messenger.Msnger.Send_message(in, mylib.CHAT_MESSAGE)
		}
	}
}

func process_messages() {
	var content string

	for {
		//	fmt.Print("Looping")
		msg := messenger.Msnger.Receive_message()
		if msg.Type == mylib.NONE {
			time.Sleep(1000000000)
			continue
		}
		content = msg.Content
		if msg.Type == mylib.CHAT_MESSAGE {
			content = fmt.Sprintf("%v says: %v", msg.Source, msg.Content)
			fmt.Println(content)
		} else if msg.Type == mylib.CREATE_ROOM {
			decoded := strings.Split(content, ":")
			rooms[decoded[0]] = fmt.Sprintf("%v:%v", decoded[1], decoded[2])
		} else if msg.Type == mylib.JOIN_ROOM {
			decoded := strings.Split(content, ":")
			if my_room == decoded[0] {
				room_members = append(room_members, fmt.Sprintf("%v:%v", decoded[1], decoded[2]))
			}
		} else if msg.Type == mylib.START_GAME { // Shouldn't be called any more...
			// Or if received and not the game host, do the start_guest functionality.
			decoded := strings.Split(content, ":")
			delete(rooms, decoded[0])
			if my_room == decoded[0] {
				in_game = true
				messenger.Msnger.Send_message(fmt.Sprintf("%v:%v:%v", decoded[0], my_name, messenger.Msnger.Port), msg.Type)
				// Need at stuff to set up teams and color
				if game_start == 1 && game_team == 1 {
					fmt.Println("DEBUG: game_color = ", game_color, " game_team = ", game_team, " player = ", my_name, " game_start = ", game_start)
					game_start++
					fmt.Println("DEBUG: game_start = ", game_start)
					//client_game.StartGame(my_room, my_name, 1, game_color, game_team)
				}
				game_start++
				//return
			}
			if msg.Source == decoded[1] {
				//return
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
					if room_members[i] == fmt.Sprintf("%v:%v", decoded[1], decoded[2]) {
						room_members = append(room_members[:i], room_members[i+1:]...)
						break
					}
				}
			}
		} else if msg.Type == mylib.MOVE {
			fmt.Println("Got a move message")
			decoded := strings.Split(content, ":")
			// Should be board_num, player_team, player_color, turn, origLoc, newLoc, capture_pieceString
			fmt.Println(decoded)
			board_num, _ := strconv.Atoi(decoded[0])
			team, _ := strconv.Atoi(decoded[1])
			color, _ := strconv.Atoi(decoded[2])
			turn, _ := strconv.Atoi(decoded[3])
			origLoc, _ := strconv.Atoi(decoded[4])
			newLoc, _ := strconv.Atoi(decoded[5])
			captured := decoded[6]
			UpdateFromOpponent(board_num, team, color, turn, origLoc, newLoc, captured)
		}else if msg.Type== mylib.CHAT_MESSAGE {
			chatting.Msg= (msg.Orig_source+":" + msg.Content)
			qml.Changed(chatting,&chatting.Msg)
		}

		msg.Type = mylib.NONE
	}

}
