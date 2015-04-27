In order to resolve the qml dependencies in game.go, the following must be done.
Have the go path set to the parent directory of the src, bin, pkg directories.
In the src dirctory, pull down the go_qml with the following command:
	go get gopkg.in/qml.v1
This should get the gopkg.in/qml.v1 directory and files needed to use qml in go.
Make sure to have the 'go path' variable properly set. This will place the above
into the src directory where the go path points.
