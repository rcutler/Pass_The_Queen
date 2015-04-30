/**
 * mylib.go: library module for common constants and structs
 * @author: Nicolas, Ryan, Xingchi
 */

package mylib

//Message types
const NONE int = 0
const CREATE_ROOM int = 1
const REQUEST_CONN_LIST int = 2
const CHAT_MESSAGE int = 3
const ACK int = 4
const NAK int = 5
const JOIN_ROOM int = 6
const START_GAME int = 7
const LEAVE_ROOM int = 8
const DELETE_ROOM int = 9
const SET_TEAM int = 10
const SET_COLOR int = 11
const MOVE int = 12
const LOCAL_INTRODUCTION int = 13
const GLOBAL_INTRODUCTION int = 14
const LEAVE_GLOBAL int = 15
const PLACE int = 16
const GAMEOVER int = 17

//Game constants
const TEAM1 int = 1
const TEAM2 int = 2
const EMPTY int = -1
const WHITE int = 1
const BLACK int = 2
const COLOR_WHITE string = "#D18B47"
const COLOR_BLACK string = "#FFCE9E"

//Game piece image locations
const BLACKROOK string = "../pieces/Rook_Black_60.png"
const BLACKKNIGHT string = "../pieces/Knight_Black_60.png"
const BLACKBISHOP string = "../pieces/Bishop_Black_60.png"
const BLACKQUEEN string = "../pieces/Queen_Black_60.png"
const BLACKKING string = "../pieces/King_Black_60.png"
const BLACKPAWN string = "../pieces/Pawn_Black_60.png"
const WHITEROOK string = "../pieces/Rook_White_60.png"
const WHITEKNIGHT string = "../pieces/Knight_White_60.png"
const WHITEBISHOP string = "../pieces/Bishop_White_60.png"
const WHITEQUEEN string = "../pieces/Queen_White_60.png"
const WHITEKING string = "../pieces/King_White_60.png"
const WHITEPAWN string = "../pieces/Pawn_White_60.png"

type Square struct {
	Color          string
	Type           string
	Image          string
	Empty          bool
	FromOtherBoard bool
	OrigPosition   bool
	Index          int
	TeamPiece      int
}

type CapturedPieces struct {
	Pieces []*CapturedPiece
	Len    int
}

func (cp *CapturedPieces) Empty() {
	cp.Pieces = nil
	cp.Len = 0
}

func (cp *CapturedPieces) Add(p CapturedPiece) {
	cp.Pieces = append(cp.Pieces, &p)
	cp.Len = len(cp.Pieces)
}

func (cp *CapturedPieces) Piece(index int) *CapturedPiece {
	return cp.Pieces[index]
}

type CapturedPiece struct {
	TeamPiece int
	Image     string
	Type      string
}

/* Network Message */
type Message struct {
	Content     string         //Content of message. Can be anything
	Orig_source string         //Original source of message
	Source      string         //Source node
	Dest        string         //Destination node
	Supernode   bool           //True if source node is a supernode
	Type        int            //Message type (see above)
	Timestamp   map[string]int //Vector timestamp
}

/* Vector clock */
type VectorClock struct {
	name         string         //Node name that owns this vector clock
	current_time map[string]int //current time in vector form
}

/* Default vector clock constructor */
func NewVectorClock(name string) VectorClock {
	var v VectorClock
	v.name = name
	v.current_time = make(map[string]int)
	v.current_time[name] = 1
	return v
}

/* Returns current vector clock time */
func (v *VectorClock) CurTime() map[string]int {
	ret := make(map[string]int)
	for k, v := range v.current_time {
		ret[k] = v
	}
	return ret
}

/* Compares current time with the input vector timestamp and
 * updates the current time to be the max of the two
 * If the input is nil, the time for vector clock owner node is
 * incremented by one */
func (v *VectorClock) Update(t map[string]int) {
	if t == nil {
		v.current_time[v.name]++
	} else {
		for name, time := range t {
			v.current_time[name] = max(time, v.current_time[name])
		}
	}
}

/* Removes a node from the timestamp vector */
func (v *VectorClock) Remove(name string) {
	delete(v.current_time, name)
}

/* Clears the vector clock */
func (v *VectorClock) Clear() {
	for k, _ := range v.current_time {
		if k != v.name {
			delete(v.current_time, k)
		}
	}
}

/* Returns the max of two ints */
func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
