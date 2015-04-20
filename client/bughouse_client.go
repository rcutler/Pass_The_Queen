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
var game_start int

var msnger messenger.Messenger

//Local network
var room_members []string

func main() {

	game_start = 1
	my_name = os.Args[1]
	in_room = false
	in_game = false
	rooms = make(map[string]string)

	msnger = messenger.NewMessenger(my_name)

	fmt.Print("reached")
	msnger.Login()
	fmt.Print("reached")

	go process_messages()

	fmt.Printf("\nChat:\n")
	reader := bufio.NewReader(os.Stdin)
	for {
		in, _ := reader.ReadString('\n')
		in = strings.Split(in, "\n")[0]
		command := strings.Split(in, " ")
		fmt.Printf("Read: %q\n", in)
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
				// Need at stuff to set up teams and color
				fmt.Println("DEBUG start command from owner: ", my_room, " ", my_name, " ", 1, " ", game_color, " ", game_team)
				//start_local_chat(my_room, my_name, 1, game_color, game_team)
				start_local_chat()
				//				client_game.StartGame(my_room, my_name, 1, game_color, game_team)
			}
			//Leave a room
		} else if in == "start_guest" {
			if !in_room {
				fmt.Println("Not in a room")
			} else if rooms[my_room] == fmt.Sprintf("%v:%v", my_name, msnger.Port) {
				fmt.Println("The room owner, use 'start' instead")
			} else {
				in_game = true
				client_game.StartGame(my_room, my_name, 1, game_color, game_team)
			}
		} else if in == "leave" {
			if in_room {
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

func start_local_chat() {
	fmt.Print("STARTING A GAME")
	//TODO
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
		fmt.Println("Received a message")
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
				msnger.Send_message(fmt.Sprintf("%v:%v:%v", decoded[0], my_name, msnger.Port), msg.Type)
				// Need at stuff to set up teams and color
				/*fmt.Println("DEBUG: When am I here 2?")
				fmt.Println("DEBUG process_message: ", my_room, " ", my_name, " ", 1, " ", game_color, " ", game_team)
				start_local_chat(my_room, my_name, 1, game_color, game_team)*/
				start_local_chat()
				fmt.Println("DEBUG: game_color = ", game_color, " game_team = ", game_team, " player = ", my_name, " game_start = ", game_start)
				if game_start == 1 && game_team == 1 {
					game_start++
					fmt.Println("DEBUG: game_start = ", game_start)
					//client_game.StartGame(my_room, my_name, 1, game_color, game_team)
				}
				game_start++
				return
			}
			if msg.Source == decoded[1] {
				return
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
		}

		msg.Type = mylib.NONE
	}

}
