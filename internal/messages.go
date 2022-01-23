package internal

import (
	"encoding/json"
)

type Envelope struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Messages sent by client
const (
	JoinRoomType        = "JOIN_ROOM"
	SendTextMessageType = "SEND_TEXT_MESSAGE"
)

type JoinRoom struct {
	RoomName   string `json:"roomName"`
	ClientName string `json:"clientName"`
}

type LeaveRoom struct {
	RoomName string `json:"roomName"`
}

type SendTextMessage struct {
	Text string `json:"text"`
}

func ParseClientMessages(rawMessage []byte) (interface{}, error) {
	var msg json.RawMessage
	env := Envelope{Data: &msg}
	err := json.Unmarshal(rawMessage, &env)
	if err != nil {
		return err, nil
	}

	var parsedMsg interface{}
	switch env.Type {
	case JoinRoomType:
		var joinRoomData JoinRoom
		err := json.Unmarshal(msg, &joinRoomData)
		if err != nil {
			return nil, err
		}
		parsedMsg = joinRoomData
	case SendTextMessageType:
		var textData SendTextMessage
		err := json.Unmarshal(msg, &textData)
		if err != nil {
			return nil, err
		}
		parsedMsg = textData
	}
	return parsedMsg, nil
}

// Messages sent by server
const (
	SuccessJoinRoomType    string = "SUCCESS_JOIN_ROOM"
	ClientListType                = "CLIENT_LIST"
	ReceiveTextMessageType        = "RECEIVE_TEXT_MESSAGE"
)

type SuccessJoinRoom struct {
	RoomName string `json:"roomName"`
}

type GetAllClientNames struct {
	ClientNames []string `json:"clientNames"`
}

type ReceiveTextMessage struct {
	Text       string `json:"text"`
	ClientName string `json:"clientName"`
	Id         int    `json:"id"`
}

func NewSuccessJoinRoomMessage(roomName string) []byte {
	env := &Envelope{Type: SuccessJoinRoomType}
	env.Data = &SuccessJoinRoom{roomName}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

func NewClientNamesMessage(clientNames []string) []byte {
	env := &Envelope{Type: ClientListType}
	env.Data = &GetAllClientNames{clientNames}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

func NewReceiveTextMessage(text string, clientName string, id int) []byte {
	env := &Envelope{Type: ReceiveTextMessageType}
	env.Data = &ReceiveTextMessage{text, clientName, id}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}
