
import QtQuick 2.0
import QtQuick.Window 2.0
import QtQuick.Controls 1.1

Item {
    id: top

    width: 400
    height: 600

    ListModel {
        id: chatContent
        ListElement {
            content: "chatting.msg"
        }
    }


    Rectangle {
        id: background
        z: 0
        anchors.fill: parent
        color: "#5d5b59"
        //color: "#100359"
    }


    Rectangle {
        id: chatBox
        opacity: 1
        anchors.centerIn: top

        color: "#5d5b59"
        border.color: "black"
        border.width: 1
        radius: 5
        anchors.fill: parent

        function sendMessage(){

            var hasFocus = input.focus;
            input.focus = false;

            var data = input.text
            input.clear()
            //chatContent.append({content: "Me: " + data})
            
            //send message out
            chatting.sendChatMsg(data)
            

            chatView.positionViewAtEnd()

            input.focus = hasFocus; 
        }
        Item {
            anchors.fill: parent
            anchors.margins: 10

            Rectangle {
                height: parent.height - input.height - 15
                width: parent.width
                color: "#d7d6d5"
                anchors.top: parent.top
                border.color: "black"
                border.width: 1
                radius: 5

                ListView {
                    id: chatView
                    width: parent.width-5
                    height: parent.height-5
                    anchors.centerIn: parent
                    model: chatContent
                    clip: true
                    delegate: Component {
                        Text {
                            font.pointSize: 14
                            text: chatting.msg
                        }
                    }
                }
            }

            InputBox {
                id: input
                Keys.onReturnPressed: chatBox.sendMessage()
                height: sendButton.height
                width: parent.width - sendButton.width - 15
                anchors.left: parent.left
                anchors.bottom: parent.bottom
            }

            Button {
                id: sendButton
                anchors.right: parent.right
                height: parent.height*0.15
                text: "Send"
                onClicked: chatBox.sendMessage()
                anchors.bottom: parent.bottom
            }

        }
    }


}
