import QtQuick 2.0
import QtQuick.Window 2.0
import QtQuick.Controls 1.1

Rectangle {
	width: applicationView.width
	height: applicationView.height
	//width: 600
	//height: 600

	color: "#222222"

	// Add the ChessBoard File
	ChessBoard {
		id: board
		x: 20
		y: 20
	}

	ChatWin {
		id: chat
		x: 450
		y: 20
	}
	CapturedPieces {
		id: cp
		x: 20
		y: 500
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
			local.visible = false
			global.visible = true
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
	Item {
		x: 100
		y: 450
		id: test

		Text {
			id: timerText
			font.pointSize: 24
			color: "white"
			text: chessBoard.time
		}
	}

	// Add the CapturedPieces File


	// Add the local chat box area.

}