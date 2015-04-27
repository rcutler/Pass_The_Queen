// View and user logic for the chess board

import QtQuick 2.0

GridView {
		
		property int selectedPiece: -1

		cellWidth: 50
		cellHeight: 50
		width: 50
		height: 400

		model: capturedPieces.len

		// Can put the delegate into its own file
		delegate: CapturedPiecesDelegate {
			id: cpDelegate
		}
	}