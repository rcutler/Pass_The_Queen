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

	width: Screen.width
	height: Screen.height

	color: "#222222"

	GlobalView {
		id: global
		visible: true
	}

	LocalView {
		id: local
		visible: false
	}

}