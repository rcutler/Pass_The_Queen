package main

import (
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
var is_supernode int
var norm_encoders map[string]*gob.Encoder
var super_encoders map[string]*gob.Encoder

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
	super_encoders[name].Encode(&mylib.Message{"", my_name, name, is_supernode})

	var msg mylib.Message
	for {
		dec.Decode(&msg)
		new_content := fmt.Sprintf("%v says: %v", msg.Source, msg.Content)
		fmt.Print(new_content)
		for dest, cur_enc := range norm_encoders {
			cur_enc.Encode(&mylib.Message{new_content, my_name, dest, is_supernode})
		}
	}

}

/* Connection to normal node or other supernode */
func serverSocketConnection(conn net.Conn) {

	defer conn.Close()

	var msg mylib.Message
	dec := gob.NewDecoder(conn)
	dec.Decode(&msg)

	if msg.Type == mylib.SUPER {
		super_encoders[msg.Source] = gob.NewEncoder(conn)
	} else {
		norm_encoders[msg.Source] = gob.NewEncoder(conn)
	}

	for {
		dec.Decode(&msg)
		new_content := fmt.Sprintf("%v says: %v", msg.Source, msg.Content)
		fmt.Print(new_content)
		for dest, cur_enc := range norm_encoders {
			if dest != msg.Source {
				cur_enc.Encode(&mylib.Message{new_content, my_name, dest, is_supernode})
			}
		}
		if msg.Type != mylib.SUPER {
			for dest, cur_enc := range super_encoders {
				cur_enc.Encode(&mylib.Message{new_content, msg.Source, dest, is_supernode})
			}
		}
	}
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

	my_name = os.Args[1]
	norm_encoders = make(map[string]*gob.Encoder)
	super_encoders = make(map[string]*gob.Encoder)
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
	port := rand.Int()%48127 + 1024 //1024 - 49151
	for {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
		if err != nil {
			port = rand.Int()%48127 + 1024 //1024 - 49151
		} else {
			go serverSocket(ln)
			break
		}
	}

	//Get list of nodes to connect to
	enc.Encode(&mylib.Message{fmt.Sprintf("%v:%v", my_name, port), my_name, "server", 0})
	var msg mylib.Message
	dec.Decode(&msg)
	decoded_message := strings.Split(msg.Content, " ")
	if decoded_message[0] == "true" {
		is_supernode = 1
	} else {
		is_supernode = 0
	}

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
		for dest, cur_enc := range super_encoders {
			cur_enc.Encode(&mylib.Message{in, my_name, dest, is_supernode})
		}
		for dest, cur_enc := range norm_encoders {
			cur_enc.Encode(&mylib.Message{in, my_name, dest, is_supernode})
		}
	}
}
