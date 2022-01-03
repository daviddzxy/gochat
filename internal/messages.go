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
	CreateRoomType string = "CREATE_ROOM"
	JoinRoomType          = "JOIN_ROOM"
	TextType              = "TEXT"
)

type JoinRoom struct {
	ChatRoomId int `json:"chatRoomId"`
}

type CreateRoom struct {
	ChatRoomName string `json:"chatRoomName"`
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
			return err, nil
		}
		parsedMsg = joinData
	case CreateRoomType:
		var createRoomData CreateRoom
		err := json.Unmarshal(msg, &createRoomData)
		if err != nil {
			return err, nil
		}
		parsedMsg = createRoomData
	case TextType:
		var textData Text
		err := json.Unmarshal(msg, &textData)
		if err != nil {
			return err, nil
		}
		parsedMsg = textData
	}
	return parsedMsg, nil
}

// Messages sent by server
const (
	SuccessCreateRoomType    string = "SUCCESS_CREATE_ROOM"
	SuccessJoinRoomType             = "SUCCESS_JOIN_ROOM"
	FailJoinRoomType                = "FAIL_JOIN_ROOM"
	UnableToParseMessageType        = "UNABLE_TO_PARSE"
)

type SuccessCreateRoom struct {
	ChatRoomId int `json:"chatRoomId"`
}

func NewSuccessCreateRoomMessage(chatRoomId int) []byte {
	env := &Envelope{Type: SuccessCreateRoomType}
	env.Data = &SuccessCreateRoom{chatRoomId}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

type SuccessJoinRoom struct {
	ChatRoomId int `json:"chatRoomId"`
}

func NewSuccessJoinRoomMessage(chatRoomId int) []byte {
	env := &Envelope{Type: SuccessJoinRoomType}
	env.Data = &SuccessJoinRoom{chatRoomId}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

type FailJoinRoom struct {
	ChatRoomId int `json:"chatRoomId"`
}

func NewFailJoinRoomMessage(chatRoomId int) []byte {
	env := &Envelope{Type: FailJoinRoomType}
	env.Data = &FailJoinRoom{chatRoomId}
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
