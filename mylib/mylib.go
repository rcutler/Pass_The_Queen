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
	Content     string
	Orig_source string
	Source      string
	Dest        string
	Supernode   bool
	Type        int
	Timestamp   map[string]int
}

type VectorClock struct {
	name         string
	current_time map[string]int
}

func NewVectorClock(name string) VectorClock {
	var v VectorClock
	v.name = name
	v.current_time = make(map[string]int)
	v.current_time[name] = 0
	return v
}

func (v *VectorClock) CurTime() map[string]int {
	ret := make(map[string]int)
	for k, v := range v.current_time {
		ret[k] = v
	}
	return ret
}

func (v *VectorClock) Update(t map[string]int) {
	if t == nil {
		v.current_time[v.name]++
	} else {
		for name, time := range t {
			v.current_time[name] = max(time, v.current_time[name])
		}
	}
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
