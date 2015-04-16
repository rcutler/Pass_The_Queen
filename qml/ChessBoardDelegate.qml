import QtQuick 2.0

Component {
	Rectangle {
		width: 50
		height: 50
		color: chessBoard.square(index).color
		Image {
			id: image
			cache: false
			width: 50
			height: 50
			smooth: true
			source: chessBoard.square(index).image
			MouseArea {
				anchors.fill: parent
				// Put the function part of onClicked into its own file
				onClicked: {
					if(board.selectedSquare != -1){
						if(board.selectedSquare == index){
							board.selectedSquare = -1
							//console.log("Reset the selected square: " + board.selectedSquare)
						}
						else if (board.nextSquare == index) {
							board.nextSquare = -1
							//console.log("Reset the nextSquare: " + board.nextSquare)
						}
						else {
							board.nextSquare = index
							//console.log("Set the next Square to index: " + board.nextSquare + " and the index is: " + index)
							chessBoard.movePiece(board.selectedSquare, board.nextSquare)
							board.nextSquare = -1
							board.selectedSquare = -1
							//console.log("Am I here?")
						}
					}
					else {
						board.selectedSquare = index
						//console.log("Set selected Square to index: " + board.selectedSquare + " and the index is: " + index)
					}
					//console.log("DONE!")
				}
			}
		}
	}
}