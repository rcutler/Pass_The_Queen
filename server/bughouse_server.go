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
	"sync"
)

//Starting a game (host):
//1. (client.go) sendMessage(START_GAME)
//2. (client.go) msnger.Leave_global()
//-> (messenger.go) send LEAVE_GLOBAL (first to server -> ACK, then other nodes)
//-> (messenger.go) kill all node connections
//3. (client.go) msnger.connect(members[])
//-> (messenger.go) connect to members[]

//Starting a game (member):
//1. (client.go) receive START_GAME with room == my_room
//2. (client.go) msnger.Leave_global()
//-> (messenger.go) send LEAVE_GLOBAL (first to server -> ACK, then other nodes)
//-> (messenger.go) kill all node connections
//3. (client.go) msnger.connect(members[])
//-> (messenger.go) connect to members[]

//Supernode leave network:
//1. (client.go) receive LEAVE_GLOBAL with msg.is_supernode && (~is_supernode) && msg.origSource == msg.Source (|| heartbeat fails)
//-> (server.go) receive LEAVE_GLOBAL with msg.is_supernode (|| SUPER_MISSING from a child) -> remove from supernode table
//2. (client.go) msnger.Leave_global()
//-> (messenger.go) kill all connections
//3. (client.go) msnger.Join_global()
//-> (messenger.go) sendGameServer(REQUEST_CONN_LIST)

var supernodes map[string]int
var rooms map[string]string
var lock sync.Mutex
var test int

func serverSocketConnection(conn net.Conn) {

	test = 0
	defer conn.Close()

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	for {
		var msg mylib.Message
		dec.Decode(&msg)
		lock.Lock()
		if msg.Type == mylib.REQUEST_CONN_LIST {
			fmt.Printf("Accepted connection from: %v\n", msg.Content)
			min_count := len(supernodes)
			min_name := ""
			for node, num_members := range supernodes {
				if num_members < min_count {
					min_count = num_members
					min_name = node
				}
			}
			reply := ""
			if min_count < len(supernodes) {
				reply = fmt.Sprintf("false %v", min_name)
				supernodes[min_name]++
			} else {
				reply = "true"
				for node, _ := range supernodes {
					reply = fmt.Sprintf("%v %v", reply, node)
				}
				supernodes[msg.Content] = 0
			}
			enc.Encode(&mylib.Message{reply, "server", "server", strings.Split(msg.Content, ":")[1], false, mylib.REQUEST_CONN_LIST, nil})
		} else if msg.Type == mylib.CREATE_ROOM {
			name_available := true
			decoded := strings.Split(msg.Content, ":")
			for room_name := range rooms {
				if room_name == decoded[0] {
					name_available = false
				}
			}
			if name_available {
				enc.Encode(&mylib.Message{"", "server", "server", msg.Source, false, mylib.ACK, nil})
				rooms[decoded[0]] = fmt.Sprintf("%v:%v", decoded[1], decoded[2])
			} else {
				enc.Encode(&mylib.Message{"", "server", "server", msg.Source, false, mylib.NAK, nil})
			}
		} else if msg.Type == mylib.START_GAME {
			delete(rooms, msg.Content)
		} else if msg.Type == mylib.DELETE_ROOM {
			delete(rooms, msg.Content)
		} else if msg.Type == mylib.LEAVE_GLOBAL {
			if msg.Supernode {
				delete(supernodes, msg.Content)
			} else if supernodes[msg.Content] > 0 {
				supernodes[msg.Content]--
			}
			enc.Encode(&mylib.Message{"", "server", "server", msg.Source, false, mylib.ACK, nil})
		} else {
			lock.Unlock()
			return
		}
		lock.Unlock()
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
	supernodes = make(map[string]int)
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
			for node, num_members := range supernodes {
				fmt.Printf("%v: %v\n", node, num_members)
			}
		case in == "rooms":
			for room_name, owner := range rooms {
				fmt.Printf("%v (%v)\n", room_name, owner)
			}

		}

	}
}
