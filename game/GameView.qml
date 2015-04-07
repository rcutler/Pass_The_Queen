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

Rectangle {
	id: gameView

	width: Screen.width
	height: Screen.height

	color: "#222222"

	// Add the ChessBoard File
	ChessBoard {
		id: board
		x: 20
		y: 20
	}

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