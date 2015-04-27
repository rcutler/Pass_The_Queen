// View and user logic for the chess board

import QtQuick 2.0

GridView {
		//id: board
		
		property int selectedSquare: -1
		property int nextSquare: -1

		cellWidth: 50
		cellHeight: 50
		width: 400
		height: 400

		model: chessBoard.len

		// Can put the delegate into its own file
		delegate: ChessBoardDelegate {
			id: boardDelegate
		}
	}