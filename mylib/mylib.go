package mylib

//Default type is 0 (no meaning)
const NONE int = 0
const SUPER int = 1

type Message struct {
	Content string
	Source  string
	Dest    string
	Type    int
}
