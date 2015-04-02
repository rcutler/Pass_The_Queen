package main

import (
	"Pass_The_Queen/mylib"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	var connected_nodes []string
	name := os.Args[1]
	fmt.Println("starting client")
	fmt.Println("contacting team4.ece842.com for game server address (todo)")
	fmt.Println("contacting game server")

	//Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Connection failed", err)
	}
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	enc.Encode(&mylib.Message{name})

	//Get list of nodes to connect to
	cur_supernode := &mylib.Message{}
	dec.Decode(cur_supernode)
	is_supernode := cur_supernode.Content

	for {
		dec.Decode(cur_supernode)
		//todo: make "stop" be part of the message struct rather than its own message
		if cur_supernode.Content == "STOP" {
			break
		} else {
			connected_nodes = append(connected_nodes, cur_supernode.Content)
		}
	}

	//Print list of nodes to connect to
	fmt.Printf("Is supernode: %v\n", is_supernode)
	fmt.Printf("Connecting to:\n")
	for i := 0; i < len(connected_nodes); i++ {
		fmt.Printf("%v\n", connected_nodes[i])
	}

	//Connect to other nodes
	//todo

	conn.Close()
	fmt.Println("closing client")
}
