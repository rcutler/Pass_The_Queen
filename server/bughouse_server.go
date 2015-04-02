package main

import (
	"Pass_The_Queen/mylib"
	"encoding/gob"
	"fmt"
	"net"
)

//Todo's:
//----Global p2p----
//Get server to return clients addresses rather than names when connecting
//Get client nodes to connect to each other after talking to the server
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

func handleConnection(conn net.Conn) {

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
	if min_count < len(supernodes) {
		enc.Encode(&mylib.Message{"false"})
		enc.Encode(&mylib.Message{supernodes[min_index]})
		member_count[min_index]++
	} else {
		enc.Encode(&mylib.Message{"true"})
		supernodes = append(supernodes, cur_node.Content)
		member_count = append(member_count, 1)
		for i := 0; i < len(supernodes); i++ {
			//fmt.Printf("Sending: %v\n", supernodes[i])
			if supernodes[i] != cur_node.Content {
				enc.Encode(&mylib.Message{supernodes[i]})
			}
		}
	}
	//todo: make stop part of the last message (different field), rather than its own
	enc.Encode(&mylib.Message{"STOP"})
}

func main() {
	fmt.Println("starting server")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("server failed to listen to port 8080")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("server failed to accept connection")
			continue
		}
		go handleConnection(conn)
	}
}
