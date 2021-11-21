package internal

import (
	"encoding/json"
)

type Envelope struct {
	Type string      `json:"messageType"`
	Msg  interface{} `json:"message"`
}

// Messages sent by client
const (
	CreateRoomType string = "createRoom"
	JoinType              = "join"
)

type Join struct {
	ChatRoomId int `json:"chatRoomId"`
}

type CreateRoom struct {
	// TODO: meta room settings
}

func ParseClientMessages(rawMessage []byte) (interface{}, error) {
	var msg json.RawMessage
	env := Envelope{Msg: &msg}
	err := json.Unmarshal(rawMessage, &env)
	if err != nil {
		return err, nil
	}

	var parsedMsg interface{}
	switch env.Type {
	case JoinType:
		var joinMsg Join
		err := json.Unmarshal(msg, &joinMsg)
		if err != nil {
			return err, nil
		}
		parsedMsg = joinMsg
	case CreateRoomType:
		var createRoomMsg CreateRoom
		err := json.Unmarshal(msg, &createRoomMsg)
		if err != nil {
			return err, nil
		}
		parsedMsg = createRoomMsg
	}
	return parsedMsg, nil
}

// Messages sent by server
const (
	SuccessCreateRoomType    string = "successCreateRoom"
	SuccessJoinRoomType             = "successJoinRoom"
	FailJoinRoomType                = "failJoinRoom"
	UnableToParseMessageType        = "unableToParse"
)

type SuccessCreateRoom struct {
	ChatRoomId int `json:"chatRoomId"`
}

func NewSuccessCreateRoomMessage(chatRoomId int) []byte {
	env := &Envelope{Type: SuccessCreateRoomType}
	env.Msg = &SuccessCreateRoom{chatRoomId}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

type SuccessJoinRoom struct {
	ChatRoomId int `json:"chatRoomId"`
}

func NewSuccessJoinRoomMessage(chatRoomId int) []byte {
	env := &Envelope{Type: SuccessJoinRoomType}
	env.Msg = &SuccessJoinRoom{chatRoomId}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

type FailJoinRoom struct {
	ChatRoomId int `json:"chatRoomId"`
}

func NewFailJoinRoomMessage(chatRoomId int) []byte {
	env := &Envelope{Type: FailJoinRoomType}
	env.Msg = &FailJoinRoom{chatRoomId}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

type UnableToParse struct {
	// TODO: send info message
}

func NewUnableToParseMessage() []byte {
	env := &Envelope{Type: UnableToParseMessageType}
	env.Msg = &UnableToParse{}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}
