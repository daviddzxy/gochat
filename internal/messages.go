package internal

import (
	"encoding/json"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Messages sent by Client
const (
	JoinType = "JOIN"
	PartType = "PART"
	TextType = "TEXT"
)

type Join struct {
	RoomHandle    string `json:"roomHandle"`
	SessionHandle string `json:"roomSessionHandle"`
}

type Part struct {
	RoomHandle string `json:"roomHandle"`
}

type Text struct {
	RoomHandle string `json:"roomHandle"`
	Content    string `json:"content"`
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
	case PartType:
		var partData Part
		err := json.Unmarshal(jsonData, &partData)
		if err != nil {
			return nil, err
		}
		message.Data = partData
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

const (
	SuccessJoinType string = "SUCCESS_JOIN"
	SuccessPartType        = "SUCCESS_PART"
	ReceiveTextType        = "RECEIVE_TEXT"
)

type ReceiveText struct {
	Content       string `json:"content"`
	RoomHandle    string `json:"roomHandle"`
	RoomSessionId int    `json:"roomSessionId"`
}

func NewReceiveTextMessage(content string, roomHandle string, RoomSessionId int) []byte {
	env := &Message{Type: ReceiveTextType}
	env.Data = &ReceiveText{content, roomHandle, RoomSessionId}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

type SucessJoin struct {
	RoomHandle    string         `json:"roomHandle"`
	RoomSessionId int            `json:"roomSessionId"`
	RoomSessions  []*roomSession `json:"roomSessions"`
}

func NewSuccessJoin(roomHandle string, RoomSessionId int, roomSessions []*roomSession) []byte {
	env := &Message{Type: SuccessJoinType}
	env.Data = &SucessJoin{roomHandle, RoomSessionId, roomSessions}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

type SuccessPart struct {
	RoomHandle string `json:"roomHandle"`
}

func NewSuccessPart(roomHandle string) []byte {
	env := &Message{Type: SuccessPartType}
	env.Data = &SuccessPart{roomHandle}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}
