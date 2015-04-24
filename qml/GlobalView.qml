import QtQuick 2.0
import QtQuick.Window 2.0
import QtQuick.Controls 1.1

Rectangle {
	width: applicationView.width
	height: applicationView.height

	property int host: 0

	// Button For joining a room
	Button {
		x: 200
		y: 200
		width: 150
		height: 40
		id: joinRoom
		text: "Join a room"
		onClicked: {
			startGame.visible = true
			changeTeam.visible = true
			changeColor.visible = true
			joinRoom.visible = false
			createRoom.visible = false
			leaveRoom.visible = true
			listRooms.y = 400
			game.joinRoom("a") // Replace the string value with value from a text field
		}
	}

	// Button for starting a room
	Button {
		x: 200
		y: 100
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
			leaveRoom.visible = false
			listRooms.y = 300
			// Create a function for like starting a timer when the game starts if it is your turn or something. And then have another thing to call 
			game.startRoom(host)
			host = 0
			timerGame.start()
		}
	}

	// Button for creating a room
	Button {
		x: 200
		y: 100
		width: 150
		height: 40
		id: createRoom
		text: "Create a room"
		onClicked: {
			startGame.visible = true
			changeTeam.visible = true
			changeColor.visible = true
			leaveRoom.visible = true
			joinRoom.visible = false
			createRoom.visible = false
			listRooms.y = 400
			host = 1
			game.createRoom("a") // Replace the string value with value from a text field
		}
	}

	// Button for changing team you are on
	Button {
		x: 200
		y: 200
		width: 150
		height: 40
		id: changeTeam
		visible: false
		text: "Change Team"
		onClicked: {
			// Change the team value in the back end and where it is displayed.
			game.changeTeam()
		}
	}

	// Button for changing the color you are playing as
	Button {
		x: 200
		y: 300
		width: 150
		height: 40
		visible: false
		id: changeColor
		text: "Change Color"
		onClicked: {
			// Change the color value in the back end and where it is displayed
			game.changeColor()
		}
	}	

	// Button to leave the room currently in
	Button {
		x: 200
		y: 500
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
			leaveRoom.visible = false
			listRooms.y = 300
			game.leaveRoom()
		}
	}
	
	// List of rooms that are available


	// Button to update the list of rooms available
	Button {
		x: 200
		y: 300
		width: 150
		height: 40
		id: listRooms
		visible: true
		text: "List Rooms"
		onClicked: {
			// Send a message and update the gui
			game.listGames()
			// Set some text area equal to the text returned
		}
	}

	// Create an item. The item is a text box area with a list of games


	// Create an item. The item is a text box area with a list of game members


	// Chat box area

}