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
)

type Join struct {
	RoomHandle           string `json:"roomHandle"`
	NewRoomSessionHandle string `json:"roomSessionHandle"`
}

type Part struct {
	RoomHandle string `json:"roomHandle"`
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
	}
	return message, nil
}
