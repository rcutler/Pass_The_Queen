import QtQuick 2.0

Component {
	Rectangle {
		width: 50
		height: 50
		color: "#D18B47"
		Image {
			id: image
			cache: false
			width: 50
			height: 50
			smooth: true
			source: capturedPieces.piece(index).image
			MouseArea {
				anchors.fill: parent
				// Put the function part of onClicked into its own file
				onClicked: {
					if (capturedPieces.piece(index).teamPiece == game.playerColor) {
						selectedPiece = index
					} else {
						selectedPiece = -1
					}
					console.log(selectedPiece)
				}
			}
		}
	}
}