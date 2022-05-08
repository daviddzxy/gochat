package internal

import (
	"encoding/json"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ClientDetails struct {
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
	SuccessJoinType  string = "SUCCESS_JOIN"
	SuccessPartType         = "SUCCESS_PART"
	ReceiveTextType         = "RECEIVE_TEXT"
	AddClientType           = "ADD_CLIENT"
	RemoveClientType        = "REMOVE_CLIENT"
)

type SuccessJoin struct {
	RoomName      string          `json:"roomName"`
	ClientDetails []ClientDetails `json:"clients"`
}

type ReceiveText struct {
	Text     string `json:"text"`
	ClientId int    `json:"clientId"`
	TextId   int    `json:"textId"`
}

type AddClient struct {
	ClientDetail ClientDetails `json:"client"`
}

type RemoveClient struct {
	Id int `json:"id"`
}

func NewSuccessJoinMessage(roomName string, users []ClientDetails) []byte {
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

func NewReceiveTextMessage(text string, clientId int, TextId int) []byte {
	env := &Message{Type: ReceiveTextType}
	env.Data = &ReceiveText{text, clientId, TextId}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

func NewAddClientMessage(clientName string, clientId int) []byte {
	env := &Message{Type: AddClientType}
	env.Data = &AddClient{ClientDetails{clientName, clientId}}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

func NewRemoveClientMessage(clientId int) []byte {
	env := &Message{Type: RemoveClientType}
	env.Data = &RemoveClient{clientId}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}
