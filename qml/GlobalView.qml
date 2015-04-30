import QtQuick 2.0
import QtQuick.Window 2.0
import QtQuick.Controls 1.1

Rectangle {
	width: applicationView.width
	height: applicationView.height

	property int host: 0
	property int boardNum: 1

	property bool butt: false	

	GlobalChatWin{
		x:370
		y:100
	}
	// Button For joining a room
	Button {
		x: 50
		y: 200
		width: 150
		height: 40
		id: joinRoom
		text: "Join a room"
		onClicked: {
			if (gameName.text != "" && game.checkGames(gameName.text)) {
				startGame.visible = true
				changeTeam.visible = true
				changeColor.visible = true
				leaveRoom.visible = true
				joinRoom.visible = false
				createRoom.visible = false
				host = 1
				gameName.visible = false
				gameColor.visible = true
				gameTeam.visible = true
				gameBoardInput.visible = true
				changeBoard.visible = true
				gameBoard.visible = true
				gameColor.text = "Color: Black"
				gameTeam.text = "Team: 2"
				game.joinRoom(gameName.text) // Replace the string value with value from a text field
				console.log(gameName.text)
			}
		}
	}

	// Button for starting a room
	Button {
		x: 50
		y: 150
		width: 150
		height: 40
		id: startGame
		visible: false
		text: "Start a game"
		onClicked: {
			local.visible = true
			global.visible = false
			// Reset global view to default
			joinRoom.visible = true
			createRoom.visible = true
			startGame.visible = false
			changeTeam.visible = false
			changeColor.visible = false
			gameColor.visible = false
			gameTeam.visible = false
			changeBoard.visible = false
			gameBoard.visible = false
			leaveRoom.visible = false
			gameBoardInput.visible = false
			gameName.visible = true
			game.startRoom(host, boardNum)
			host = 0
			timerGame.start()
		}
	}

	// Button for creating a room
	Button {
		x: 50
		y: 150
		width: 150
		height: 40
		id: createRoom
		text: "Create a room"
		onClicked: {
			console.log(game.checkGames(gameName.text))
			//game.listGames()
			if (gameName.text != "" && !game.checkGames(gameName.text)) {
				startGame.visible = true
				changeTeam.visible = true
				changeColor.visible = true
				leaveRoom.visible = true
				joinRoom.visible = false
				createRoom.visible = false
				host = 1
				gameName.visible = false
				gameColor.visible = true
				gameTeam.visible = true
				gameBoardInput.visible = true
				changeBoard.visible = true
				gameBoard.visible = true
				gameColor.text = "Color: White"
				gameTeam.text = "Team: 1"
				game.createRoom(gameName.text) // Replace the string value with value from a text field
				console.log(gameName.text)
			}
		}
	}

	// Button for changing team you are on
	Button {
		x: 50
		y: 270
		width: 150
		height: 40
		id: changeTeam
		visible: false
		text: "Change Team"
		onClicked: {
			// Change the team value in the back end and where it is displayed.
			gameTeam.text = game.changeTeam()
		}
	}

	Text {
		id: gameTeam
		visible: false
		x: 50
		y: 320
		width: 150
		height: 40

	}

	// Button for changing the color you are playing as
	Button {
		x: 50
		y: 200
		width: 150
		height: 40
		visible: false
		id: changeColor
		text: "Change Color"
		onClicked: {
			// Change the color value in the back end and where it is displayed
			gameColor.text = game.changeColor()
		}
	}	

	Text {
		id: gameColor
		visible: false
		x: 50
		y: 250
		width: 150
		height: 40

	}

	Button {
		x: 50
		y: 350
		width: 150
		height: 40
		visible: false
		id: changeBoard
		text: "Change Board"
		onClicked: {
			// Change the color value in the back end and where it is displayed
			boardNum = gameBoardInput.text
		}
	}	

	Text {
		id: gameBoard
		visible: false
		x: 50
		y: 400
		text: boardNum
		width: 150
		height: 40
	}

	TextField {
		id: gameBoardInput
		placeholderText: "Enter board number"
		validator: IntValidator {bottom: 1; top: 50;}
		focus: true
		x: 50
		y: 430
		visible: false
		width: 150
		height: 40
	}

	// Button to leave the room currently in
	Button {
		x: 50
		y: 480
		width: 150
		height: 40
		id: leaveRoom
		visible: false
		text: "Leave Room"
		onClicked: {
			joinRoom.visible = true
			createRoom.visible = true
			startGame.visible = false
			changeTeam.visible = false
			changeColor.visible = false
			changeBoard.visible = false
			gameBoard.visible = false
			leaveRoom.visible = false
			gameName.visible = true
			gameBoardInput.visible = false
			gameColor.visible = false
			gameTeam.visible = false 
			game.leaveRoom()
		}
	}
	
	// List of rooms that are available


	// Button to update the list of rooms available
	Button {
		x: 50
		y: 100
		width: 150
		height: 40
		id: listRooms
		visible: true
		text: "List Rooms"
		onClicked: {
			// Send a message and update the gui
			gameList.text = game.listGames()
			// Set some text area equal to the text returned
		}
	}

	// Create an item. The item is a text box area with a list of games
	Text {
		id: gameListHeader
		text: "List of available games"
		width: 300
		height:100
		x: 210
		y: 100
	}
	Text {
		id: gameList
		text: ""
		x: 210
		y: 100
		width: 300
		height: 110
	}

	// Input text box for game name thingy.
	TextField {
		id: gameName
		placeholderText: "Enter a room name"
		x: 50
		y: 250
		width: 150
		height: 40
	}

	// Chat box area

}