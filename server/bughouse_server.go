/**
 * bughouse_server.go: Pass the Queen Bughouse Chess Server program
 * The server is responsible for accepting new nodes into the global
 * network and regulating which nodes are supernodes and which nodes
 * are regular nodes. The server is also responsible for replacing
 * supernodes when they crash or leave the network
 * @author: Nicolas
 */

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

var supernodes map[string]int
var rooms map[string]string
var lock sync.Mutex
var test int

/* Connection between server and a single Node */
func serverSocketConnection(conn net.Conn) {

	is_supernode := false
	supernode_name := ""
	node_name := ""
	test = 0
	defer conn.Close()

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	for {
		var msg mylib.Message
		dec.Decode(&msg)
		lock.Lock()

		//Node wishes to join the Global Network:
		//Responds with:
		//- supernode status [true/false]
		//- list of nodes to connect to
		//- list of rooms existing in the global network
		if msg.Type == mylib.REQUEST_CONN_LIST {
			fmt.Printf("Accepted connection from: %v\n", msg.Content)
			node_name = msg.Content
			min_count := len(supernodes)
			min_name := ""

			//Determine whether the new node will be a supernode or not:
			//Adds normal nodes to supernodes until all supernodes have
			//as many members as there are supernodes in the network
			//The next node after that will become a new supernode
			for node, num_members := range supernodes {
				if num_members < min_count {
					min_count = num_members
					min_name = node
				}
			}
			reply := ""
			if min_count < len(supernodes) {
				//Node will be a normal node
				is_supernode = false
				supernode_name = min_name
				reply = fmt.Sprintf("false %v", min_name)
				supernodes[min_name]++
			} else {
				//Node will be a super node
				reply = "true"
				is_supernode = true
				supernode_name = msg.Content
				for node, _ := range supernodes {
					reply = fmt.Sprintf("%v %v", reply, node)
				}
				supernodes[msg.Content] = 0
			}

			//Send a list of existing rooms
			for room_name, owner := range rooms {
				enc.Encode(&mylib.Message{fmt.Sprintf("%v:%v", room_name, owner), "server", "server", strings.Split(msg.Content, ":")[0], false, mylib.CREATE_ROOM, nil})
			}

			//Send a list of nodes to connect to
			enc.Encode(&mylib.Message{reply, "server", "server", strings.Split(msg.Content, ":")[0], false, mylib.REQUEST_CONN_LIST, nil})

		} else if msg.Type == mylib.CREATE_ROOM {
			//A node wants to create a new room
			name_available := true
			decoded := strings.Split(msg.Content, ":")
			//Check if the room name is available
			for room_name := range rooms {
				if room_name == decoded[0] {
					name_available = false
				}
			}
			if name_available {
				enc.Encode(&mylib.Message{"", "server", "server", msg.Source, false, mylib.ACK, nil})
				rooms[decoded[0]] = fmt.Sprintf("%v:%v:%v", decoded[1], decoded[2], decoded[3])
			} else {
				enc.Encode(&mylib.Message{"", "server", "server", msg.Source, false, mylib.NAK, nil})
			}
		} else if msg.Type == mylib.START_GAME {
			//A node started a game => remove the name from the rooms list
			delete(rooms, msg.Content)
		} else if msg.Type == mylib.DELETE_ROOM {
			//A room was deleted => remove the name from the rooms list
			delete(rooms, msg.Content)
		} else if msg.Type == mylib.LEAVE_GLOBAL {
			//A node left the global network
			if msg.Supernode {
				//Delete the supernode entry
				delete(supernodes, msg.Content)
			} else if supernodes[msg.Content] > 0 {
				//Decrement the member count of its supernode
				supernodes[msg.Content]--
			}
			//Delete any rooms that are in the name of the node
			for room_name, owner := range rooms {
				if owner == node_name {
					delete(rooms, room_name)
				}
			}
			//Send a success message
			enc.Encode(&mylib.Message{"", "server", "server", msg.Source, false, mylib.ACK, nil})
		} else {
			//Connected node crashed. Do same actions as mylib.LEAVE_GLOBAL
			lock.Unlock()
			if is_supernode {
				delete(supernodes, supernode_name)
			} else if supernodes[supernode_name] > 0 {
				supernodes[supernode_name]--
			}
			for room_name, owner := range rooms {
				if owner == node_name {
					delete(rooms, room_name)
				}
			}
			return
		}
		lock.Unlock()
	}

}

/* Listens for connecting nodes and spawns a connection thread */
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

/* Main Server routine. Starts the server socket and parses user commands */
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
	//Set up server socket
	go serverSocket(ln)

	//Parse user commands
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
