package mylib

//Default type is 0 (no meaning)
const NONE int = 0
const CREATE_ROOM int = 1
const REQUEST_CONN_LIST int = 2
const CHAT_MESSAGE int = 3
const ACK = 4
const NAK = 5
const JOIN_ROOM int = 6
const START_GAME int = 7
const LEAVE_ROOM int = 8
const DELETE_ROOM int = 9
const SET_TEAM int = 10
const SET_COLOR int = 11

type Message struct {
	Content   string
	Source    string
	Dest      string
	Supernode bool
	Type      int
}
