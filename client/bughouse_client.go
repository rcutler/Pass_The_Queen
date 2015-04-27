package main

import (
	"Pass_The_Queen/client_game"
	"Pass_The_Queen/messenger"
	"Pass_The_Queen/mylib"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var my_name string
var my_room string
var in_room bool
var in_game bool
var rooms map[string]string
var game_team int
var game_color int

var msnger messenger.Messenger

//Local network
var room_members []string

func main() {

	my_name = os.Args[1]
	in_room = false
	in_game = false
	rooms = make(map[string]string)

	msnger = messenger.NewMessenger(my_name)

	msnger.Login()

	go process_messages()

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
				msnger.Send_game_server(fmt.Sprintf("%v:%v:%v", room_name, my_name, msnger.Port), mylib.CREATE_ROOM)
				msg := msnger.Receive_game_server()
				if msg.Type == mylib.ACK {
					fmt.Printf("Creating room: %v\n", room_name)
					rooms[room_name] = fmt.Sprintf("%v:%v", my_name, msnger.Port)
					msnger.Send_message(fmt.Sprintf("%v:%v:%v", room_name, my_name, msnger.Port), mylib.CREATE_ROOM)
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
					msnger.Send_message(fmt.Sprintf("%v:%v:%v", room_name, my_name, msnger.Port), mylib.JOIN_ROOM)
					in_room = true
					my_room = room_name
				}
			}
		} else if in == "start" {
			if !in_room {
				fmt.Println("Not in a room")
			} else if rooms[my_room] != fmt.Sprintf("%v:%v", my_name, msnger.Port) {
				fmt.Println("Not the room owner")
			} else {
				in_game = true
				msnger.Send_game_server(my_room, mylib.START_GAME)
				msnger.Send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, msnger.Port), mylib.START_GAME)
				delete(rooms, my_room)
				msnger.Leave_global()
				msnger.Join_local(room_members)
				start_game()
			}
			//Leave a room
		} else if in == "leave" {
			if in_game {
				msnger.Leave_local()
				msnger.Join_global()
			} else if in_room {
				//Delete room if room owner
				if rooms[my_room] == fmt.Sprintf("%v:%v", my_name, msnger.Port) {
					msnger.Send_game_server(my_room, mylib.DELETE_ROOM)
					msnger.Send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, msnger.Port), mylib.DELETE_ROOM)
					delete(rooms, my_room)
				} else {
					msnger.Send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, msnger.Port), mylib.LEAVE_ROOM)
				}
				in_room = false
				my_room = ""
				room_members = make([]string, 0, 0)
			} else {
				msnger.Leave_global()
				return
			}
			//Print list of room members
		} else if in == "members" {
			for i := range room_members {
				fmt.Println(room_members[i])
			}
		} else if command[0] == "set_team" { // Change the current players team.
			if !in_room {
				fmt.Println("Not in a room!")
				client_game.StartGame("a", "b", 1, 2, 2)
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
			msnger.Send_message(in, mylib.CHAT_MESSAGE)
		}
	}
}

func start_game() {
	fmt.Println("Inside start_game")
	//client_game.StartGame(my_room, my_name, 1, game_color, game_team)
}

func start_game_guest() {
	fmt.Println("Inside start_game_guest")
	//client_game.StartGame(my_room, my_name, 1, game_color, game_team)
}
func process_messages() {
	var content string

	for {
		//	fmt.Print("Looping")
		msg := msnger.Receive_message()
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
		} else if msg.Type == mylib.START_GAME {
			decoded := strings.Split(content, ":")
			delete(rooms, decoded[0])
			if my_room == decoded[0] {
				in_game = true
				msnger.Leave_global()
				msnger.Join_local(room_members)
				start_game_guest()
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
		} else if msg.Type == mylib.LEAVE_GLOBAL {
			if msg.Orig_source == msg.Source && msg.Supernode && !msnger.Is_supernode {
				msnger.Leave_global()
				msnger.Join_global()
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
			client_game.UpdateFromOpponent(board_num, team, color, turn, origLoc, newLoc, captured)
		}

		msg.Type = mylib.NONE
	}

}
