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
)

var chatting *ChatMsg
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

//Local network
var room_members []string

func StartGame(room string, player string, board int, color int, team int) {
	fmt.Println("I am in the start game function.... good news.")
	game.Name = room
	game.PlayerID = player
	game.Board = board
	game.PlayerColor = color
	game.TeamPlayer = team

	fmt.Println("DEBUG: Board Values. game.Board = ", game.Board, " board = ", board)

	turn = 0
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

	chessBoard.initialize()

	engine.Context().SetVar("game", game)
	engine.Context().SetVar("chessBoard", chessBoard)
	//chat room variable
	engine.Context().SetVar("chatting", chatting)

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
				messenger.Msnger.Leave_global()
				messenger.Msnger.Join_local(room_members)
				StartGame(my_room, my_name, 1, game_color, game_team)
			} else {
				in_game = true
				messenger.Msnger.Send_game_server(my_room, mylib.START_GAME)
				messenger.Msnger.Send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, messenger.Msnger.Port), mylib.START_GAME)
				delete(rooms, my_room)
				messenger.Msnger.Leave_global()
				messenger.Msnger.Join_local(room_members)
				//start_game() // Calls StartGame
				StartGame(my_room, my_name, 1, game_color, game_team)
			}
			//Leave a room
			/*		} else if in == "start_guest" {
					if !in_room {
						fmt.Println("Not in a room")
					} else if rooms[my_room] == fmt.Sprintf("%v:%v", my_name, messenger.Msnger.Port) {
						fmt.Println("The room owner, use 'start' instead")
					} else {
						in_game = true
						StartGame(my_room, my_name, 1, game_color, game_team)
					}*/
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
		} else {
			messenger.Msnger.Send_message(in, mylib.CHAT_MESSAGE)
		}
	}
}

func (chat ChatMsg) SendChatMsg(data string) {
	fmt.Println("******************************************")
	fmt.Println("sending: " + data)
	chatting.Msg += "Me: " + data + "\n"
	qml.Changed(chatting, &chatting.Msg)
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
			content = fmt.Sprintf("%v says: %v", msg.Source, msg.Content)
			fmt.Println(content)
			chatting.Msg += content + "\n"
			qml.Changed(chatting, &chatting.Msg)

		} else if msg.Type == mylib.CREATE_ROOM {
			decoded := strings.Split(content, ":")
			rooms[decoded[0]] = fmt.Sprintf("%v:%v", decoded[1], decoded[2])
		} else if msg.Type == mylib.JOIN_ROOM {
			decoded := strings.Split(content, ":")
			if my_room == decoded[0] {
				room_members = append(room_members, fmt.Sprintf("%v:%v", decoded[1], decoded[2]))
			}
		} else if msg.Type == mylib.START_GAME {
			decoded := strings.Split(content, ":")
			delete(rooms, decoded[0])
			if my_room == decoded[0] {
				in_game = true
				messenger.Msnger.Leave_global()
				messenger.Msnger.Join_local(room_members)
				//start_game_guest()
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
			board_num, _ := strconv.Atoi(decoded[0])
			team, _ := strconv.Atoi(decoded[1])
			color, _ := strconv.Atoi(decoded[2])
			turn, _ := strconv.Atoi(decoded[3])
			origLoc, _ := strconv.Atoi(decoded[4])
			newLoc, _ := strconv.Atoi(decoded[5])
			captured := decoded[6]
			fmt.Println("DEBUG: decoded = ", decoded)
			UpdateFromOpponent(board_num, team, color, turn, origLoc, newLoc, captured)
		}
		msg.Type = mylib.NONE
	}
}
