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
	Is_supernode   bool                    //Whether this node is a supernode
	supernode      string                  //Name of this node's supernode
	encoders       map[string]*gob.Encoder //List of encoders
	local_conns    []*net.Conn             //Array of local network connections
	global_conns   []*net.Conn             //Array of global network connections
	enc            *gob.Encoder            //Game Server encoder
	dec            *gob.Decoder            //Game Server decoder
	rcv_buffer     []*mylib.Message        //received message buffer (before ordering)
	deliver_buffer []*mylib.Message        //messages ready to for delivery (after ordering)
	v_clock        mylib.VectorClock       //Vector Clock
	received       []*mylib.Message        //All received messages
}

var Msnger Messenger

/* Messenger default constructor */
func NewMessenger(name string) Messenger {
	var m Messenger
	m.name = name
	m.Is_supernode = false
	m.encoders = make(map[string]*gob.Encoder)
	m.v_clock = mylib.NewVectorClock(name)
	return m
}

/* Connect Messenger with game server and connect to the global chat */
func (m *Messenger) Login() {

	fmt.Println("starting client")
	fmt.Println("contacting team4.ece842.com")

	//Connect to server
	//conn, err := net.Dial("tcp", "localhost:8080")
	conn, err := net.Dial("tcp", "52.11.139.181:8080")

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

	m.Join_global()

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
func (m *Messenger) clientSocket(name string, port int, conn_type int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%v", port))
	defer conn.Close()
	if err != nil {
		fmt.Printf("Failed to connect to %v on port %v\n", name, port)
		log.Fatal("", err)
		return
	}
	fmt.Printf("Connected to %v on port %v\n", name, port)

	m.encoders[name] = gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	if conn_type == mylib.GLOBAL_INTRODUCTION {
		m.global_conns = append(m.global_conns, &conn)
	} else {
		m.local_conns = append(m.local_conns, &conn)
	}
	//Introduce this node
	m.encoders[name].Encode(&mylib.Message{"", m.name, m.name, name, m.Is_supernode, conn_type, m.v_clock.CurTime()})
	defer delete(m.encoders, name)

	m.receive_messages(name, dec)
}

/* Connection to normal node or other supernodes */
func (m *Messenger) serverSocketConnection(conn net.Conn) {

	defer conn.Close()

	var msg mylib.Message

	dec := gob.NewDecoder(conn)
	dec.Decode(&msg)
	name := msg.Source
	if msg.Type == mylib.GLOBAL_INTRODUCTION {
		m.global_conns = append(m.global_conns, &conn)
	} else {
		m.local_conns = append(m.local_conns, &conn)
	}

	m.encoders[name] = gob.NewEncoder(conn)
	defer delete(m.encoders, name)

	//Remove new connection's previous history from vector clock and received array
	m.v_clock.Remove(name)
	for {
		success := true
		for i := range m.received {
			if m.received[i].Orig_source == name {
				m.received = append(m.received[:i], m.received[i+1:]...)
				success = false
				break
			}
		}
		if success {
			break
		}
	}

	m.receive_messages(name, dec)
}

/* Receives and handles incoming messages */
func (m *Messenger) receive_messages(name string, dec *gob.Decoder) {

	for {

		var msg mylib.Message

		//Receive message
		dec.Decode(&msg)
		content := msg.Content

		//fmt.Printf("Received: %q at %q\n", msg, m.v_clock.CurTime())

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

			if msg.Orig_source == m.name {
				already_received = true
			}

			//Pass new messages to the client and forward the message to other nodes
			if !already_received {

				//Insert message in buffer
				m.rcv_buffer = append(m.rcv_buffer, &msg)
				m.received = append(m.received, &msg)

				//Forward message to other nodes if it did not originate at this node
				if msg.Type == mylib.CHAT_MESSAGE {
					content = fmt.Sprintf("%v says: %v", msg.Source, content)
				}
				if msg.Orig_source != m.name {
					for dest, cur_enc := range m.encoders {
						//fmt.Printf("Forwarded: %q\n", mylib.Message{content, msg.Orig_source, m.name, dest, m.Is_supernode, msg.Type, msg.Timestamp})
						cur_enc.Encode(&mylib.Message{content, msg.Orig_source, m.name, dest, m.Is_supernode, msg.Type, msg.Timestamp})
					}
				}

				//Deliver any ready messages to the client
				m.deliver_ordered_messages()
			}
		} else {
			m.v_clock.Remove(name)
			for {
				success := true
				for i := range m.received {
					if m.received[i].Orig_source == name {
						m.received = append(m.received[:i], m.received[i+1:]...)
						success = false
						break
					}
				}
				if success {
					break
				}
			}
			supernode_name := strings.Split(m.supernode, ":")[0]
			if !m.Is_supernode && name == supernode_name {
				m.Leave_global()
				m.Join_global()
			}
			return
		}
	}
}

/* Send a message to the network */
func (m *Messenger) Send_message(content string, Type int) {
	m.received = append(m.received, &mylib.Message{content, m.name, m.name, "", m.Is_supernode, Type, m.v_clock.CurTime()})
	//fmt.Printf("%q Sending: %q\n", m.v_clock.CurTime(), content)
	for dest, cur_enc := range m.encoders {
		//fmt.Printf("Sent: %q\n", mylib.Message{content, m.name, m.name, dest, m.Is_supernode, Type, m.v_clock.CurTime()})
		cur_enc.Encode(&mylib.Message{content, m.name, m.name, dest, m.Is_supernode, Type, m.v_clock.CurTime()})
	}
	m.v_clock.Update(nil)
}

/* Receive a message from the network */
func (m *Messenger) Receive_message() *mylib.Message {
	if len(m.deliver_buffer) != 0 {
		//	fmt.Printf("Receiving: %q\n", m.deliver_buffer[0])
		msg := m.deliver_buffer[0]
		m.deliver_buffer = m.deliver_buffer[1:]
		return msg
	}
	return &mylib.Message{}
}

/* Send a message to the game server */
func (m *Messenger) Send_game_server(content string, Type int) {
	m.enc.Encode(&mylib.Message{content, m.name, m.name, "server", m.Is_supernode, Type, m.v_clock.CurTime()})
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

/* Leave global chat */
func (m *Messenger) Leave_global() {
	m.Send_game_server(m.supernode, mylib.LEAVE_GLOBAL)

	var msg mylib.Message
	m.dec.Decode(&msg)
	if msg.Type != mylib.ACK {
		fmt.Println("FATAL ERROR: failed to leave global net")
	}

	m.v_clock.Clear()
	//m.Send_message(m.supernode, mylib.LEAVE_GLOBAL)

	for i := range m.global_conns {
		(*m.global_conns[i]).Close()
	}
	m.global_conns = nil

}

/* Leave local chat */
func (m *Messenger) Leave_local() {
	for i := range m.local_conns {
		(*m.local_conns[i]).Close()
	}
	m.local_conns = nil

}

/* Join global chat */
func (m *Messenger) Join_global() {

	//Get list of nodes to connect to
	m.enc.Encode(&mylib.Message{fmt.Sprintf("%v:%v", m.name, m.Port), m.name, m.name, "server", false, mylib.REQUEST_CONN_LIST, m.v_clock.CurTime()})
	m.received = nil
	var msg mylib.Message
	m.dec.Decode(&msg)
	decoded_message := strings.Split(msg.Content, " ")
	m.Is_supernode = (decoded_message[0] == "true")
	fmt.Printf("Is supernode: %v\n", decoded_message[0])
	if !m.Is_supernode {
		m.supernode = decoded_message[1]
		fmt.Printf("Supernode name: %v\n", decoded_message[1])
	} else {
		m.supernode = fmt.Sprintf("%v:%v", m.name, m.Port)
	}

	//Print list of nodes to connect to
	fmt.Printf("Connecting to:\n")
	for i := 1; i < len(decoded_message); i++ {
		fmt.Printf("%v\n", decoded_message[i])
	}

	//Connect to list of nodes to connect to
	for i := 1; i < len(decoded_message); i++ {
		cur_node := strings.Split(decoded_message[i], ":")
		cur_port, _ := strconv.Atoi(cur_node[1])
		go m.clientSocket(cur_node[0], cur_port, mylib.GLOBAL_INTRODUCTION)
	}
}

/* Join local chat */
func (m *Messenger) Join_local(members []string) {
	m.Is_supernode = true
	m.received = nil

	fmt.Printf("Connecting to:\n")
	for i := range members {
		fmt.Printf("%v\n", members[i])
	}

	//Connect to list of nodes to connect to
	for i := range members {
		cur_node := strings.Split(members[i], ":")
		cur_port, _ := strconv.Atoi(cur_node[1])
		go m.clientSocket(cur_node[0], cur_port, mylib.LOCAL_INTRODUCTION)
	}
}
