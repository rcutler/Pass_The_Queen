import QtQuick 2.0
import QtQuick.Window 2.0
import QtQuick.Controls 1.1

/*
This is the main application window that gets spawned when the client gets spawned.
It will set up a global view. And have a local view for the game hidden until a game starts.

Each view will be a rectangle the size of the window. Will have a function as part of a button that will change from global to local when starting a game.
*/

Rectangle {
	id: applicationView

	//width: Screen.width
	//height: Screen.height
	width: 600
	height: 600

	color: "#222222"

	GlobalView {
		id: global
		visible: true
	}

	LocalView {
		id: local
		visible: false
	}

	// Try creating a timer here and having its text value read in the local view
	Item {
		id: timerItem
		Timer {
			id: timerGame
			interval: 1000
			running: false
			repeat: true
			onTriggered: {
				// Have this call a function in the go code that will decrement a count value and then do a qml.changed on it to have it updated. That way, can check to see if it is the players turn, if it is, decement. Otherwise just return. Then can have check in GO to see if game ended. And can resolve it that way.
				chessBoard.timer()
			}
		}
	}

}