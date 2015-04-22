import QtQuick 2.0
import QtQuick.Window 2.0
import QtQuick.Controls 1.1

Rectangle {
	width: Screen.width
	height: Screen.height

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
		}
	}

	// Chat box area


}