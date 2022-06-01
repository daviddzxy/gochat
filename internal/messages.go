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
	JoinType = "JOIN"
	PartType = "PART"
	TextType = "TEXT"
)

type Join struct {
	RoomHandle           string `json:"roomHandle"`
	NewRoomSessionHandle string `json:"roomSessionHandle"`
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
	ReceiveTextType string = "RECEIVE_TEXT"
)

type ReceiveText struct {
	Content       string `json:"content"`
	RoomSessionId int    `json:"roomSessionId"`
}

func NewReceiveTextMessage(content string, RoomSessionId int) []byte {
	env := &Message{Type: ReceiveTextType}
	env.Data = &ReceiveText{content, RoomSessionId}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}
