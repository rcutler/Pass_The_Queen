package messenger

import (
	"Pass_The_Queen/mylib"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

/* Messenger object used to communicate with other nodes in the local or global chat */
type Messenger struct {
	name         string                  //Name of this node
	Port         int                     //Port number of this node
	is_supernode bool                    //Whether this node is a supernode
	encoders     map[string]*gob.Encoder //List of encoders
	enc          *gob.Encoder            //Game Server encoder
	dec          *gob.Decoder            //Game Server decoder
	rcv_buffer   []*mylib.Message        //Received message buffer
	v_clock      mylib.VectorClock       //Vector Clock
	received     []*mylib.Message        //All received messages
}

/* Messenger default constructor */
func NewMessenger(name string) Messenger {
	var m Messenger
	m.name = name
	m.is_supernode = false
	m.encoders = make(map[string]*gob.Encoder)
	m.v_clock = mylib.NewVectorClock(name)
	//m.rcv_buffer = make([]*mylib.Message,0)
	return m
}

/* Connect Messenger with game server and connect to the global chat */
func (m Messenger) Login() {

	fmt.Println("starting client")
	fmt.Println("contacting team4.ece842.com for game server address (todo)")
	fmt.Println("contacting game server")

	//Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	defer conn.Close()

	if err != nil {
		fmt.Println("Failed to connect to server on Port 8080")
		log.Fatal("", err)
		return
	}

	m.enc = gob.NewEncoder(conn)
	m.dec = gob.NewDecoder(conn)

	//Set up server for incoming connections
	m.Port = rand.Int()%48127 + 1024 //1024 - 49151
	for {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%v", m.Port))
		if err != nil {
			m.Port = rand.Int()%48127 + 1024 //1024 - 49151
		} else {
			go m.serverSocket(ln)
			break
		}
	}

	//Get list of nodes to connect to
	m.enc.Encode(&mylib.Message{fmt.Sprintf("%v:%v", m.name, m.Port), m.name, m.name, "server", false, mylib.REQUEST_CONN_LIST, m.v_clock.CurTime()})
	m.v_clock.Update(nil)
	var msg mylib.Message
	m.dec.Decode(&msg)
	decoded_message := strings.Split(msg.Content, " ")
	m.is_supernode = (decoded_message[0] == "true")

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
		go m.clientSocket(cur_node[0], cur_port)
	}
}

/* Receives incoming connections from normal nodes/other supernodes */
func (m Messenger) serverSocket(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("client failed to accept connection")
			continue
		}
		go m.serverSocketConnection(conn)
	}
}

/* Connection to this node's supernode(s) */
func (m Messenger) clientSocket(name string, port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%v", port))
	defer conn.Close()
	if err != nil {
		fmt.Printf("Failed to connect to %v on port %v\n", name, port)
		log.Fatal("", err)
		return
	}

	m.encoders[name] = gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	//Introduce this node
	m.encoders[name].Encode(&mylib.Message{"", m.name, m.name, name, m.is_supernode, mylib.NONE, m.v_clock.CurTime()})
	m.v_clock.Update(nil)
	defer delete(m.encoders, name)

	m.receive_messages(dec)
}

/* Connection to normal node or other supernodes */
func (m Messenger) serverSocketConnection(conn net.Conn) {

	defer conn.Close()

	var msg mylib.Message

	dec := gob.NewDecoder(conn)
	dec.Decode(&msg)
	name := msg.Source

	m.encoders[name] = gob.NewEncoder(conn)
	defer delete(m.encoders, name)

	m.receive_messages(dec)
}

/* Receives and handles incoming messages */
func (m Messenger) receive_messages(dec *gob.Decoder) {
	var msg mylib.Message
	var content string

	for {
		//Receive message
		dec.Decode(&msg)
		content = msg.Content

		if msg.Type != mylib.NONE {
			fmt.Printf("Received: %q from %q (orig) and %q (imm) to %q (super?: %q) of type %q at time %q\n",
				msg.Content, msg.Orig_source, msg.Source, msg.Dest, m.is_supernode, msg.Type, msg.Timestamp)

			//Compare Timestamps to see whether the message was already received
			already_received := false
			for t := range m.received {
				all_equal := true
				for name, time := range m.received[t].Timestamp {
					if time != msg.Timestamp[name] {
						all_equal = false
						break
					}
				}
				if all_equal {
					already_received = true
					break
				}
			}

			//Pass new messages to the client and forward the message to other nodes
			if !already_received {

				//Insert message in buffer
				m.rcv_buffer = append(m.rcv_buffer, &msg)
				m.received = append(m.received, &msg)

				//TODO: break connections on certain condition such as starting a game etc.
				//TODO: set is_supernode to true on starting a game
				//TODO: provide functions to connect to other members on starting a game

				//Forward message to other nodes if it did not originate at this node
				if msg.Orig_source != m.name {
					for dest, cur_enc := range m.encoders {
						cur_enc.Encode(&mylib.Message{content, msg.Orig_source, m.name, dest, m.is_supernode, msg.Type, msg.Timestamp})
					}
				}
			}
		}
		msg.Type = mylib.NONE
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
}
func start_local_chat() {
	is_supernode = true
	for i := range room_members {
		decoded := strings.Split(room_members[i], ":")
		port, _ := strconv.Atoi(decoded[1])
		go clientSocket(decoded[0], port)
	}
}*/

// Send message to all connected nodes
func (m Messenger) Send_message(content string, Type int) {
	fmt.Printf("Inside send_message\n")
	for dest, cur_enc := range m.encoders {
		fmt.Printf("Sending: %q from %q and %q to %q (super?: %q) of type %q at time %q\n", content, m.name, m.name, dest, m.is_supernode, Type, m.v_clock.CurTime())
		cur_enc.Encode(&mylib.Message{content, m.name, m.name, dest, m.is_supernode, Type, m.v_clock.CurTime()})
	}
	m.v_clock.Update(nil)
}

// Send message to game server
func (m Messenger) Send_game_server(content string, Type int) {
	m.enc.Encode(&mylib.Message{content, m.name, m.name, "server", m.is_supernode, Type, m.v_clock.CurTime()})
	m.v_clock.Update(nil)
}

func (m Messenger) Receive_game_server() *mylib.Message {
	var msg mylib.Message
	m.dec.Decode(&msg)
	return &msg
}

func (m Messenger) Receive_message() *mylib.Message {
	if len(m.rcv_buffer) != 0 {
		msg := m.rcv_buffer[0]
		m.rcv_buffer = m.rcv_buffer[1:]
		return msg
	}
	return &mylib.Message{}
}
