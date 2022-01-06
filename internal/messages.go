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
	JoinRoomType = "JOIN_ROOM"
	TextType     = "TEXT"
)

type JoinRoom struct {
	RoomName string `json:"roomName"`
	UserName string `json:"userName"`
}

type Text struct {
	ChatRoomId int    `json:"chatRoomId"`
	Text       string `json:"text"`
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
		var joinData JoinRoom
		err := json.Unmarshal(msg, &joinData)
		if err != nil {
			return nil, err
		}
		parsedMsg = joinData
	case TextType:
		var textData Text
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
	SuccessJoinRoomType      string = "SUCCESS_JOIN_ROOM"
	UnableToParseMessageType        = "UNABLE_TO_PARSE"
)

type SuccessCreateRoom struct {
	ChatRoomId int `json:"chatRoomId"`
}

type SuccessJoinRoom struct {
	RoomName  string   `json:"roomName"`
	UserNames []string `json:"userNames"`
}

func NewSuccessJoinRoomMessage(roomName string, userNames []string) []byte {
	env := &Envelope{Type: SuccessJoinRoomType}
	env.Data = &SuccessJoinRoom{roomName, userNames}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

type UnableToParse struct {
	// TODO: send info message
}

func NewUnableToParseMessage() []byte {
	env := &Envelope{Type: UnableToParseMessageType}
	env.Data = &UnableToParse{}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}
