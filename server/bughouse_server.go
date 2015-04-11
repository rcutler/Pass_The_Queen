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
//----Global p2p----
//Add chat GUI and broadcast chat messages to everyone in the network
//----Local p2p-----
//Add GUI option to join a specific room (local p2p)
//Add support for nodes leaving the global network (replace supernodes if needed, etc.)
//Move nodes from global to local p2p
//broadcast chat messages to everyone in the local group
//-----Future-------
//Add DNS lookup to find server address
//Add heartbeat messages in case nodes fail/leave
//Replace supernodes in case nodes fail/leave via election
//Make supernode network multi-tiered
//Integrate with game client

var supernodes []string
var member_count []int

func serverSocketConnection(conn net.Conn) {

	defer conn.Close()

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	cur_node := &mylib.Message{}
	dec.Decode(cur_node)
	fmt.Printf("Accepted connection from: %v\n", cur_node.Content)
	min_count := len(supernodes)
	min_index := 0
	for i := 0; i < len(member_count); i++ {
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
		supernodes = append(supernodes, cur_node.Content)
		member_count = append(member_count, 1)
		for i := 0; i < len(supernodes); i++ {
			if supernodes[i] != cur_node.Content {
				reply = fmt.Sprintf("%v %v", reply, supernodes[i])
			}
		}
	}
	enc.Encode(&mylib.Message{reply, "server", strings.Split(cur_node.Content, ":")[1], 0})
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
		switch {
		case in == "help\n" || in == "h\n":
			fmt.Println("Possible commands:")
			fmt.Println("supernodes")
		case in == "quit\n" || in == "exit\n" || in == "q\n":
			return
		case in == "supernodes\n":
			for i := 0; i < len(supernodes); i++ {
				fmt.Println(supernodes[i])
			}

		}

	}
}
