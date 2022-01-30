package internal

import (
	"encoding/json"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Messages sent by client
const (
	JoinType    = "JOIN"
	PartType    = "PART"
	TextType    = "TEXT"
	ClientsType = "CLIENTS"
)

type Join struct {
	RoomName   string `json:"roomName"`
	ClientName string `json:"clientName"`
}

type Text struct {
	Text string `json:"text"`
}

func ParseClientMessages(rawMessage []byte) (*Message, error) {
	var msg json.RawMessage
	env := &Message{Data: &msg}
	err := json.Unmarshal(rawMessage, &env)
	if err != nil {
		return nil, err
	}

	switch env.Type {
	case JoinType:
		var joinData Join
		err := json.Unmarshal(msg, &joinData)
		if err != nil {
			return nil, err
		}
	case TextType:
		var textData Text
		err := json.Unmarshal(msg, &textData)
		if err != nil {
			return nil, err
		}
	}
	return env, nil
}

// Messages sent by server
const (
	SuccessJoinType    string = "SUCCESS_JOIN"
	SuccessPartType           = "SUCCESS_PART"
	ReceiveTextType           = "RECEIVE_TEXT"
	ReceiveClientsType        = "RECEIVE_CLIENTS"
)

type SuccessJoin struct {
	RoomName string `json:"roomName"`
}

type ReceiveClients struct {
	ClientNames []string `json:"clientNames"`
}

type ReceiveText struct {
	Text       string `json:"text"`
	ClientName string `json:"clientName"`
	Id         int    `json:"id"`
}

func NewSuccessJoinMessage(roomName string) []byte {
	env := &Message{Type: SuccessJoinType}
	env.Data = &SuccessJoin{roomName}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

func NewSuccessPartMessage() []byte {
	env := &Message{Type: SuccessPartType}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

func NewClientsMessage(clientNames []string) []byte {
	env := &Message{Type: ReceiveClientsType}
	env.Data = &ReceiveClients{clientNames}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

func NewReceiveTextMessage(text string, clientName string, id int) []byte {
	env := &Message{Type: ReceiveTextType}
	env.Data = &ReceiveText{text, clientName, id}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}
