package main

import (
	"Pass_The_Queen/client_game"
	"Pass_The_Queen/mylib"
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

var my_name string
var my_room string
var my_port int
var is_supernode bool
var in_room bool
var in_game bool
var rooms map[string]string
var game_team int
var game_color int
var game_start int

//Global network
var norm_encoders map[string]*gob.Encoder
var super_encoders map[string]*gob.Encoder

//Local network
var room_members []string

/* Connection to this node's supernode(s) */
func clientSocket(name string, port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%v", port))
	defer conn.Close()
	if err != nil {
		fmt.Printf("Failed to connect to %v on port %v\n", name, port)
		log.Fatal("", err)
		return
	}

	super_encoders[name] = gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	//Introduce this node
	super_encoders[name].Encode(&mylib.Message{"", my_name, name, is_supernode, 0})
	defer delete(super_encoders, name)

	process_messages(dec)
}

/* Connection to normal node or other supernode */
func serverSocketConnection(conn net.Conn) {

	defer conn.Close()

	var msg mylib.Message

	dec := gob.NewDecoder(conn)
	dec.Decode(&msg)
	name := msg.Source

	if msg.Supernode {
		super_encoders[name] = gob.NewEncoder(conn)
		defer delete(super_encoders, name)
	} else {
		norm_encoders[name] = gob.NewEncoder(conn)
		defer delete(norm_encoders, name)
	}

	process_messages(dec)

}

/* Receives incoming connections from normal nodes/other supernodes */
func serverSocket(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("client failed to accept connection")
			continue
		}
		go serverSocketConnection(conn)
	}
}

func main() {

	game_start = 1
	my_name = os.Args[1]
	in_room = false
	in_game = false
	norm_encoders = make(map[string]*gob.Encoder)
	super_encoders = make(map[string]*gob.Encoder)
	rooms = make(map[string]string)

	fmt.Println("starting client")
	fmt.Println("contacting team4.ece842.com for game server address (todo)")
	fmt.Println("contacting game server")

	//Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	defer conn.Close()

	if err != nil {
		fmt.Println("Failed to connect to server on port 8080")
		log.Fatal("", err)
		return
	}
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	//Set up server for incoming connections
	my_port := rand.Int()%48127 + 1024 //1024 - 49151
	for {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%v", my_port))
		if err != nil {
			my_port = rand.Int()%48127 + 1024 //1024 - 49151
		} else {
			go serverSocket(ln)
			break
		}
	}

	//Get list of nodes to connect to
	enc.Encode(&mylib.Message{fmt.Sprintf("%v:%v", my_name, my_port), my_name, "server", false, mylib.REQUEST_CONN_LIST})
	var msg mylib.Message
	dec.Decode(&msg)
	decoded_message := strings.Split(msg.Content, " ")
	is_supernode = (decoded_message[0] == "true")

	//Print list of nodes to connect to
	fmt.Printf("Is supernode: %v\n", decoded_message[0])
	fmt.Printf("Connecting to:\n")
	for i := 1; i < len(decoded_message); i++ {
		fmt.Printf("%v\n", decoded_message[i])
	}

	//Connect to list of nodes to connect to
	for i := 1; i < len(decoded_message); i++ {
		cur_node := strings.Split(decoded_message[i], ":")
		cur_port, _ := strconv.Atoi(cur_node[1])
		go clientSocket(cur_node[0], cur_port)
	}

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
				enc.Encode(&mylib.Message{fmt.Sprintf("%v:%v:%v", room_name, my_name, my_port), my_name, "server", is_supernode, mylib.CREATE_ROOM})
				dec.Decode(&msg)
				if msg.Type == mylib.ACK {
					fmt.Printf("Creating room: %v\n", room_name)
					rooms[room_name] = fmt.Sprintf("%v:%v", my_name, my_port)
					send_message(fmt.Sprintf("%v:%v:%v", room_name, my_name, my_port), mylib.CREATE_ROOM)
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
					send_message(fmt.Sprintf("%v:%v:%v", room_name, my_name, my_port), mylib.JOIN_ROOM)
					in_room = true
					my_room = room_name
				}
			}
		} else if in == "start" {
			if !in_room {
				fmt.Println("Not in a room")
			} else if rooms[my_room] != fmt.Sprintf("%v:%v", my_name, my_port) {
				fmt.Println("Not the room owner")
			} else {
				in_game = true
				enc.Encode(&mylib.Message{my_room, my_name, "server", is_supernode, mylib.START_GAME})
				send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, my_port), mylib.START_GAME)
				delete(rooms, my_room)
				// Need at stuff to set up teams and color
				fmt.Println("DEBUG start command from owner: ", my_room, " ", my_name, " ", 1, " ", game_color, " ", game_team)
				//start_local_chat(my_room, my_name, 1, game_color, game_team)
				start_local_chat()
				client_game.StartGame(my_room, my_name, 1, game_color, game_team)
			}
			//Leave a room
		} else if in == "start_guest" {
			if !in_room {
				fmt.Println("Not in a room")
			} else if rooms[my_room] == fmt.Sprintf("%v:%v", my_name, my_port) {
				fmt.Println("The room owner, use 'start' instead")
			} else {
				in_game = true
				client_game.StartGame(my_room, my_name, 1, game_color, game_team)
			}
		} else if in == "leave" {
			if in_room {
				//Delete room if room owner
				if rooms[my_room] == fmt.Sprintf("%v:%v", my_name, my_port) {
					enc.Encode(&mylib.Message{my_room, my_name, "server", is_supernode, mylib.DELETE_ROOM})
					send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, my_port), mylib.DELETE_ROOM)
					delete(rooms, my_room)
				} else {
					send_message(fmt.Sprintf("%v:%v:%v", my_room, my_name, my_port), mylib.LEAVE_ROOM)
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
			send_message(in, mylib.CHAT_MESSAGE)
		}
	}
}

// Send message to all normal nodes (and potentially super nodes as well)
func send_message(content string, Type int) {
	for dest, cur_enc := range super_encoders {
		cur_enc.Encode(&mylib.Message{content, my_name, dest, is_supernode, Type})
	}
	for dest, cur_enc := range norm_encoders {
		cur_enc.Encode(&mylib.Message{content, my_name, dest, is_supernode, Type})
	}
}

func forward_message(content string, source string, Type int, broadcast bool) {
	for dest, cur_enc := range norm_encoders {
		if dest != source {
			cur_enc.Encode(&mylib.Message{content, my_name, dest, is_supernode, Type})
		}
	}
	if broadcast {
		for dest, cur_enc := range super_encoders {
			cur_enc.Encode(&mylib.Message{content, my_name, dest, is_supernode, Type})
		}
	}
}

// Need at stuff to set up teams and color
/*
func start_local_chat(room string, player string, board int, color int, team int) {
	is_supernode = true
	for i := range room_members {
		decoded := strings.Split(room_members[i], ":")
		port, _ := strconv.Atoi(decoded[1])
		go clientSocket(decoded[0], port)
	}
	// Set up the game state with the initialize function here.
	// Put in a new file/package called client_game
	// Need at stuff to set up teams and color
	fmt.Println("DEBUG start_local_chat: ", room, " ", player, " ", board, " ", color, " ", team)
	client_game.StartGame(room, player, board, color, team)
}*/
func start_local_chat() {
	is_supernode = true
	for i := range room_members {
		decoded := strings.Split(room_members[i], ":")
		port, _ := strconv.Atoi(decoded[1])
		go clientSocket(decoded[0], port)
	}
}

func process_messages(dec *gob.Decoder) {
	var msg mylib.Message
	var content string

	for {
		dec.Decode(&msg)
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
				forward_message(content, msg.Source, msg.Type, !msg.Supernode)
				send_message(fmt.Sprintf("%v:%v:%v", decoded[0], my_name, my_port), msg.Type)
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
				forward_message(content, msg.Source, msg.Type, !msg.Supernode)
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
		} else {
			return
		}

		forward_message(content, msg.Source, msg.Type, !msg.Supernode)

		msg.Type = mylib.NONE
	}

}
