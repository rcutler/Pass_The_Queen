/*  This is the QML Main window for now.
	This contains the chess board logic
	Will need to separate out the main window, the game network
	aspect, and the global network aspect.

	This file will ultimately be the GameView

	TODO:
		GameView
		- Place the Game window stuff into its own file
		- Add in timer and countdown
		- Add in button for submitting a move
		- Add in button for undoing a move
		- Add in the team game chat area
		- Add area to view captured pieces from team mates
		- Be able to place a captured piece from team mate onto the board.

		GlobalView
		- Place the Global window stuff into its own file
		- Add the global chat area
		- Add the list of available games
		- Add search for game button
		- Add create a game button

		MainView
		- Add functionality to change the view depending on if
		playing a game or not.
		- Have it access both the global and game windows.
		- Give the window a title
		- Have it set up a view for the Game and the Global
			- Initalize it so the Global View is visable
			- When game created/started, make Global View not visable and Game View visable.
			- When game ends, revert back to Global View.

*/

import QtQuick 2.0
import QtQuick.Window 2.0
import QtQuick.Controls 1.1

Rectangle {
	id: gameView

	//width: Screen.width
	//height: Screen.height
	width: 600
	height: 600

	color: "#222222"

	// Add the ChessBoard File
	ChessBoard {
		id: board
		x: 20
		y: 20
	}

	ChatWin{
		id:chat
		x:800
		y:20
	}
	// Add a button for submitting a move.
	Button {
		x: 100
		y: 450
		width: 110
		height: 30
		visible: false
		id: submitMove
		text: "Submit Move"
		onClicked: {
			undoMove.visible = true
			submitMove.visible = false

		}
	}

	// Add a button for reverting a move.
	Button {
		x: 230
		y: 450
		width: 110
		height: 30
		id: undoMove
		text: "Undo Move"
		onClicked: {
			submitMove.visible = true
			undoMove.visible = false
		}
	}

	// Add a timer and display thing for it.


	// Add the CapturedPieces File
	/* CapturedPieces {
		id: captured
	}*/

	/* GameTime {
		id: time
	}
	*/

	// Add the GameChat File
	/* GameChat {
		id: chat
	}
	*/

}