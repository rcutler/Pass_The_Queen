package main

import (
	"Pass_The_Queen/mylib"
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

//Todo's:
//----Current goals-----
//Add support for nodes leaving the global network (replace supernodes if needed, etc.)
//-----Future-------
//Add GUI
//Add DNS lookup to find server address
//Add heartbeat messages in case nodes fail/leave
//Allow nodes to send global messages until the game actually starto
//Replace supernodes in case nodes fail/leave via election
//Make supernode network multi-tiered
//Integrate with game client

var supernodes []string
var member_count []int
var rooms map[string]string

func serverSocketConnection(conn net.Conn) {

	defer conn.Close()

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	var msg mylib.Message
	for {
		dec.Decode(&msg)
		if msg.Type == mylib.REQUEST_CONN_LIST {
			fmt.Printf("Accepted connection from: %v\n", msg.Content)
			min_count := len(supernodes)
			min_index := 0
			for i := range member_count {
				if member_count[i] < min_count {
					min_count = member_count[i]
					min_index = i
				}
			}
			reply := ""
			if min_count < len(supernodes) {
				reply = fmt.Sprintf("false %v", supernodes[min_index])
				member_count[min_index]++
			} else {
				reply = "true"
				supernodes = append(supernodes, msg.Content)
				member_count = append(member_count, 1)
				for i := range supernodes {
					if supernodes[i] != msg.Content {
						reply = fmt.Sprintf("%v %v", reply, supernodes[i])
					}
				}
			}
			enc.Encode(&mylib.Message{reply, "server", strings.Split(msg.Content, ":")[1], false, mylib.REQUEST_CONN_LIST})
		} else if msg.Type == mylib.CREATE_ROOM {
			name_available := true
			decoded := strings.Split(msg.Content, ":")
			for room_name := range rooms {
				if room_name == decoded[0] {
					name_available = false
				}
			}
			if name_available {
				enc.Encode(&mylib.Message{"", "server", msg.Source, false, mylib.ACK})
				rooms[decoded[0]] = fmt.Sprintf("%v:%v", decoded[1], decoded[2])
			} else {
				enc.Encode(&mylib.Message{"", "server", msg.Source, false, mylib.NAK})
			}
		} else if msg.Type == mylib.START_GAME {
			delete(rooms, msg.Content)
			return
		} else if msg.Type == mylib.DELETE_ROOM {
			delete(rooms, msg.Content)
		} else {
			return
		}
		msg.Type = mylib.NONE
	}

}

func serverSocket(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("server failed to accept connection")
			continue
		}
		go serverSocketConnection(conn)
	}
}

/* Main Server routine. Accepts client connections. */
func main() {
	fmt.Println("starting server")
	rooms = make(map[string]string)
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Failed to set up server on port 8080")
		log.Fatal("", err)
		return
	}
	go serverSocket(ln)
	reader := bufio.NewReader(os.Stdin)
	for {
		in, _ := reader.ReadString('\n')
		in = strings.Split(in, "\n")[0]
		switch {
		case in == "help" || in == "h":
			fmt.Println("Possible commands:")
			fmt.Println("supernodes")
			fmt.Println("rooms")
		case in == "quit" || in == "exit" || in == "q":
			return
		case in == "supernodes":
			for i := range supernodes {
				fmt.Println(supernodes[i])
			}
		case in == "rooms":
			for room_name, owner := range rooms {
				fmt.Printf("%v (%v)\n", room_name, owner)
			}

		}

	}
}
