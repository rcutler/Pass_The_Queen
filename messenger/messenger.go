/**
 * messenger.go: Messenger object that sends/receives messages and interacts
 * with the game server. It also deals with transitioning between the
 * local and the global network.
 * @author: Nicolas
 */

package messenger

import (
	"Pass_The_Queen/mylib"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
)

/* Messenger object used to communicate with other nodes in the local or global chat */
type Messenger struct {
	name            string                  //Name of this node
	Port            int                     //Port of this node
	Address         string                  //Address of this node
	Is_supernode    bool                    //Whether this node is a supernode
	supernode       string                  //Name of this node's supernode
	local_encoders  map[string]*gob.Encoder //List of local encoders
	global_encoders map[string]*gob.Encoder //List of global encoders
	local_conns     []*net.Conn             //Array of local network connections
	global_conns    []*net.Conn             //Array of global network connections
	enc             *gob.Encoder            //Game Server encoder
	dec             *gob.Decoder            //Game Server decoder
	rcv_buffer      []*mylib.Message        //received message buffer (before ordering)
	deliver_buffer  []*mylib.Message        //messages ready to for delivery (after ordering)
	v_clock         mylib.VectorClock       //Vector Clock
	received        []*mylib.Message        //All received messages
	in_global       bool
}

var Msnger Messenger

/* Messenger default constructor */
func NewMessenger(name string) Messenger {
	var m Messenger
	m.name = name
	m.Is_supernode = false
	m.in_global = true
	m.local_encoders = make(map[string]*gob.Encoder)
	m.global_encoders = make(map[string]*gob.Encoder)
	m.v_clock = mylib.NewVectorClock(name)
	return m
}

/* Connect Messenger with game server and connect to the global chat */
func (m *Messenger) Login() {

	fmt.Println("starting client")
	fmt.Println("contacting game server")

	//Connect to server
	//conn, err := net.Dial("tcp", "localhost:8080") //Local server for debugging
	conn, err := net.Dial("tcp", "52.11.139.181:8080") //AWS Game server

	if err != nil {
		fmt.Println("Failed to connect to game server")
		log.Fatal("", err)
		return
	}

	m.enc = gob.NewEncoder(conn)
	m.dec = gob.NewDecoder(conn)

	//Get a random (unused) port number
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

	ip := ""

	//Get the network IP address
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error while searching ip address")
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println("Error while searching ip address")
			return
		}
		if ip != "" {
			break
		}
		for _, addr := range addrs {
			ip = addr.String()
			ip = strings.Split(ip, "/")[0]
			if len(ip) > 6 && strings.Split(ip, ".")[0] != "127" {
				break
			}
			ip = ""
		}
	}

	//Address = ip:port
	m.Address = fmt.Sprintf("%v:%v", ip, m.Port)

	//Join the global chat
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

/* Connection to this node's supernode(s) or fellow local nodes*/
func (m *Messenger) clientSocket(name string, address string, conn_type int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v", address))
	defer conn.Close()
	if err != nil {
		fmt.Printf("ERROR: Failed to connect to %v on address %v\n", name, address)
		return
	}
	fmt.Printf("Connected to %v on address %v\n", name, address)

	dec := gob.NewDecoder(conn)

	//Add new connection to either the global or local connections/encoders
	if conn_type == mylib.GLOBAL_INTRODUCTION {
		m.global_conns = append(m.global_conns, &conn)
		m.global_encoders[name] = gob.NewEncoder(conn)
		m.global_encoders[name].Encode(&mylib.Message{"", m.name, m.name, name, m.Is_supernode, conn_type, m.v_clock.CurTime()})
		defer delete(m.global_encoders, name)
	} else {
		m.local_conns = append(m.local_conns, &conn)
		m.local_encoders[name] = gob.NewEncoder(conn)
		m.local_encoders[name].Encode(&mylib.Message{"", m.name, m.name, name, m.Is_supernode, conn_type, m.v_clock.CurTime()})
		defer delete(m.local_encoders, name)
	}

	//Start message listening loop
	m.receive_messages(name, dec)
}

/* Connection to normal node, fellow supernodes, or fellow local nodes */
func (m *Messenger) serverSocketConnection(conn net.Conn) {

	defer conn.Close()

	var msg mylib.Message

	dec := gob.NewDecoder(conn)
	dec.Decode(&msg)
	name := msg.Source

	//Add incoming connection to either the global or local connections/encoders
	if msg.Type == mylib.GLOBAL_INTRODUCTION {
		fmt.Println("accepted global network connection")
		m.global_conns = append(m.global_conns, &conn)
		m.global_encoders[name] = gob.NewEncoder(conn)
		defer delete(m.global_encoders, name)
	} else {
		fmt.Println("accepted local network connection")
		m.local_conns = append(m.local_conns, &conn)
		m.local_encoders[name] = gob.NewEncoder(conn)
		defer delete(m.local_encoders, name)
	}

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

	//Start message listening loop
	m.receive_messages(name, dec)
}

/* Receives and handles incoming messages */
func (m *Messenger) receive_messages(name string, dec *gob.Decoder) {

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

			//Any messages originating from this node are discarded by default
			if msg.Orig_source == m.name {
				already_received = true
			}

			//Pass new messages to the client and forward the message to other nodes
			if !already_received {

				//Insert message in buffer
				m.rcv_buffer = append(m.rcv_buffer, &msg)
				m.received = append(m.received, &msg)

				if msg.Orig_source != m.name {
					if m.in_global {
						for dest, cur_enc := range m.global_encoders {
							cur_enc.Encode(&mylib.Message{content, msg.Orig_source, m.name, dest, m.Is_supernode, msg.Type, msg.Timestamp})
						}
					} else {
						for dest, cur_enc := range m.local_encoders {
							cur_enc.Encode(&mylib.Message{content, msg.Orig_source, m.name, dest, m.Is_supernode, msg.Type, msg.Timestamp})
						}
					}
				}

				//Deliver any ready messages to the client
				m.deliver_ordered_messages()
			}
		} else {
			//Other node crashed
			//Remove node from vector clock
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

			//If the other node was this node's supernode => reconnect to global chat
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
	//Send message
	if m.in_global {
		for dest, cur_enc := range m.global_encoders {
			cur_enc.Encode(&mylib.Message{content, m.name, m.name, dest, m.Is_supernode, Type, m.v_clock.CurTime()})
		}
	} else {
		for dest, cur_enc := range m.local_encoders {
			cur_enc.Encode(&mylib.Message{content, m.name, m.name, dest, m.Is_supernode, Type, m.v_clock.CurTime()})
		}
	}
	//Increment vector clock
	m.v_clock.Update(nil)
}

/* Deliver a message to the client */
func (m *Messenger) Receive_message() *mylib.Message {
	//Deliver any messages that are ready to be delivered
	if len(m.deliver_buffer) != 0 {
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
		log.Fatal("ERROR: failed to leave global net")
	}

	//Clear vector clock
	m.v_clock.Clear()

	for i := range m.global_conns {
		(*m.global_conns[i]).Close()
	}
	m.global_conns = nil
	m.in_global = false
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

	m.in_global = true

	//Get list of nodes to connect to
	m.enc.Encode(&mylib.Message{fmt.Sprintf("%v:%v", m.name, m.Address), m.name, m.name, "server", false, mylib.REQUEST_CONN_LIST, m.v_clock.CurTime()})
	m.received = nil

	for {
		var msg mylib.Message
		m.dec.Decode(&msg)
		//Message contains information about supernode status and which nodes to connect to
		if msg.Type == mylib.REQUEST_CONN_LIST {
			decoded_message := strings.Split(msg.Content, " ")

			if decoded_message[0] == "" {
				log.Fatal("ERROR: Game server failed to respond")
			}

			//Check supernode status
			m.Is_supernode = (decoded_message[0] == "true")
			fmt.Printf("Is supernode: %v\n", decoded_message[0])

			//Save supernode name
			if !m.Is_supernode {
				m.supernode = decoded_message[1]
				fmt.Printf("Supernode name: %v\n", decoded_message[1])
			} else {
				m.supernode = fmt.Sprintf("%v:%v", m.name, m.Address)
			}

			//Print list of nodes to connect to
			fmt.Printf("Connecting to:\n")
			for i := 1; i < len(decoded_message); i++ {
				fmt.Printf("%v\n", decoded_message[i])
			}

			//Connect to list of nodes to connect to
			for i := 1; i < len(decoded_message); i++ {
				cur_node := strings.Split(decoded_message[i], ":")
				cur_address := fmt.Sprintf("%v:%v", cur_node[1], cur_node[2])
				go m.clientSocket(cur_node[0], cur_address, mylib.GLOBAL_INTRODUCTION)
			}
			break
		}
		//Message contains existing rooms information => forward to client
		m.deliver_buffer = append(m.deliver_buffer, &msg)
	}

}

/* Join local chat */
func (m *Messenger) Join_local(members []string) {
	m.Is_supernode = true
	m.received = nil

	//Print list of nodes to connec to
	fmt.Printf("Connecting to:\n")
	for i := range members {
		fmt.Printf("%v\n", members[i])
	}

	//Connect to list of nodes to connect to
	for i := range members {
		cur_node := strings.Split(members[i], ":")
		cur_address := fmt.Sprintf("%v:%v", cur_node[1], cur_node[2])
		go m.clientSocket(cur_node[0], cur_address, mylib.LOCAL_INTRODUCTION)
	}
}
