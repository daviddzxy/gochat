package internal

import (
	"encoding/json"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type User struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

// Messages sent by client
const (
	JoinType = "JOIN"
	PartType = "PART"
	TextType = "TEXT"
)

type Join struct {
	RoomName   string `json:"roomName"`
	ClientName string `json:"clientName"`
}

type Text struct {
	Text string `json:"text"`
}

func ParseClientMessages(rawMessage []byte) (*Message, error) {
	var jsonData json.RawMessage
	message := &Message{Data: &jsonData}
	err := json.Unmarshal(rawMessage, &message)
	if err != nil {
		return nil, err
	}

	switch message.Type {
	case JoinType:
		var joinData Join
		err := json.Unmarshal(jsonData, &joinData)
		if err != nil {
			return nil, err
		}
		message.Data = joinData
	case TextType:
		var textData Text
		err := json.Unmarshal(jsonData, &textData)
		if err != nil {
			return nil, err
		}
		message.Data = textData
	}
	return message, nil
}

// Messages sent by server
const (
	SuccessJoinType string = "SUCCESS_JOIN"
	SuccessPartType        = "SUCCESS_PART"
	ReceiveTextType        = "RECEIVE_TEXT"
)

type SuccessJoin struct {
	RoomName string `json:"roomName"`
	Users    []User `json:"users"`
}

type ReceiveText struct {
	Text       string `json:"text"`
	ClientName string `json:"clientName"`
	Id         int    `json:"id"`
}

func NewSuccessJoinMessage(roomName string, users []User) []byte {
	env := &Message{Type: SuccessJoinType}
	env.Data = &SuccessJoin{roomName, users}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

func NewSuccessPartMessage() []byte {
	env := &Message{Type: SuccessPartType}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

func NewReceiveTextMessage(text string, clientName string, id int) []byte {
	env := &Message{Type: ReceiveTextType}
	env.Data = &ReceiveText{text, clientName, id}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}
