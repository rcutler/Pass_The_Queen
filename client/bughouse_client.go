package main

import (
	"Pass_The_Queen/client_game"
	"fmt"
	"gopkg.in/qml.v1"
	"os"
)

// Moving tthis into the messenger class so that I can use the existing messenger object from the game.go stuff.
//var msnger messenger.Messenger

// This may eventually be the only thing in this file.
// Will create a new file for the process function/helper.
// Move the run back into game.go
// For now, also move process_messages() to game.go
func main() {

	// Start the qml stuff here....
	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// Can eventually move this over to a GUI file/directory
// Will need to move this over to game.go...

func run() error {
	return client_game.Run()
}
