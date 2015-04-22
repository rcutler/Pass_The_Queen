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
	name           string                  //Name of this node
	Port           int                     //Port number of this node
	is_supernode   bool                    //Whether this node is a supernode
	encoders       map[string]*gob.Encoder //List of encoders
	enc            *gob.Encoder            //Game Server encoder
	dec            *gob.Decoder            //Game Server decoder
	rcv_buffer     []*mylib.Message        //received message buffer (before ordering)
	deliver_buffer []*mylib.Message        //messages ready to for delivery (after ordering)
	v_clock        mylib.VectorClock       //Vector Clock
	received       []*mylib.Message        //All received messages
}

/* Messenger default constructor */
func NewMessenger(name string) Messenger {
	var m Messenger
	m.name = name
	m.is_supernode = false
	m.encoders = make(map[string]*gob.Encoder)
	m.v_clock = mylib.NewVectorClock(name)
	return m
}

/* Connect Messenger with game server and connect to the global chat */
func (m *Messenger) Login() {

	fmt.Println("starting client")
	fmt.Println("contacting team4.ece842.com for game server address (todo)")
	fmt.Println("contacting game server")

	//Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")

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
func (m *Messenger) serverSocket(ln net.Listener) {
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
func (m *Messenger) clientSocket(name string, port int) {
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
	defer delete(m.encoders, name)

	m.receive_messages(dec)
}

/* Connection to normal node or other supernodes */
func (m *Messenger) serverSocketConnection(conn net.Conn) {

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
func (m *Messenger) receive_messages(dec *gob.Decoder) {

	for {

		var msg mylib.Message

		//Receive message
		dec.Decode(&msg)
		content := msg.Content

		if msg.Type != mylib.NONE {

			//Compare Timestamps to see whether the message was already received
			already_received := false
			for t := range m.received {
				all_equal := true
				for name, time := range msg.Timestamp {
					if time != m.received[t].Timestamp[name] {
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
				content = fmt.Sprintf("%v says: %v", msg.Source, content)
				if msg.Orig_source != m.name {
					for dest, cur_enc := range m.encoders {
						//fmt.Printf("Forwarded: %q\n", mylib.Message{content, msg.Orig_source, m.name, dest, m.is_supernode, msg.Type, msg.Timestamp})
						cur_enc.Encode(&mylib.Message{content, msg.Orig_source, m.name, dest, m.is_supernode, msg.Type, msg.Timestamp})
					}
				}

				//Deliver any ready messages to the client
				m.deliver_ordered_messages()
			}
		}
	}
}

/* Send a message to the network */
func (m *Messenger) Send_message(content string, Type int) {
	m.received = append(m.received, &mylib.Message{content, m.name, m.name, "", m.is_supernode, Type, m.v_clock.CurTime()})
	//fmt.Printf("%q Sending: %q\n", m.v_clock.CurTime(), content)
	for dest, cur_enc := range m.encoders {
		//fmt.Printf("Sent: %q\n", mylib.Message{content, m.name, m.name, dest, m.is_supernode, Type, m.v_clock.CurTime()})
		cur_enc.Encode(&mylib.Message{content, m.name, m.name, dest, m.is_supernode, Type, m.v_clock.CurTime()})
	}
	m.v_clock.Update(nil)
}

/* Receive a message from the network */
func (m *Messenger) Receive_message() *mylib.Message {
	if len(m.deliver_buffer) != 0 {
		//fmt.Printf("Receiving: %q\n", m.deliver_buffer[0])
		msg := m.deliver_buffer[0]
		m.deliver_buffer = m.deliver_buffer[1:]
		return msg
	}
	return &mylib.Message{}
}

/* Send a message to the game server */
func (m *Messenger) Send_game_server(content string, Type int) {
	m.enc.Encode(&mylib.Message{content, m.name, m.name, "server", m.is_supernode, Type, m.v_clock.CurTime()})
}

/* Receive a message from the game server */
func (m *Messenger) Receive_game_server() *mylib.Message {
	var msg mylib.Message
	m.dec.Decode(&msg)
	return &msg
}

/* Check ordering rules and release any eligible messages for reception by client */
func (m *Messenger) deliver_ordered_messages() {
	for i := range m.rcv_buffer {
		cur_msg := m.rcv_buffer[i]
		ready_to_deliver := true

		for source, time := range cur_msg.Timestamp {
			if source == cur_msg.Orig_source {
				// Check Vj[j] == Vi[j] + 1
				if time != (m.v_clock.CurTime()[source]+1) && m.v_clock.CurTime()[source] != 0 {
					ready_to_deliver = false
					break
				}
			} else {
				// Check Vj[k] <= Vi[k] (k != j)
				if time > m.v_clock.CurTime()[source] && m.v_clock.CurTime()[source] != 0 {
					ready_to_deliver = false
				}
			}
		}

		if ready_to_deliver {
			m.v_clock.Update(cur_msg.Timestamp)
			m.deliver_buffer = append(m.deliver_buffer, cur_msg)
			m.rcv_buffer = append(m.rcv_buffer[:i], m.rcv_buffer[i+1:]...)
		}
	}
}
