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

func NewJoinMessage(chatRoomId int) []byte {
	env := &Envelope{Type: JoinType}
	env.Msg = &Join{ChatRoomId: chatRoomId}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
}

type CreateRoom struct {
	// TODO: meta room settings
}

func NewCreateRoomMessage() []byte {
	env := &Envelope{Type: JoinType}
	env.Msg = &CreateRoom{}
	jsonMsg, _ := json.Marshal(env)
	return jsonMsg
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

// Messages used by server
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

func ParseServerMessages(rawMessage []byte) (interface{}, error) {
	var msg json.RawMessage
	env := Envelope{Msg: &msg}
	err := json.Unmarshal(rawMessage, &env)
	if err != nil {
		return err, nil
	}

	var parsedMsg interface{}
	switch env.Type {
	case SuccessCreateRoomType:
		var msgSuccessCreateRoom SuccessCreateRoom
		err := json.Unmarshal(msg, &msgSuccessCreateRoom)
		if err != nil {
			return err, nil
		}
		parsedMsg = msgSuccessCreateRoom
	case SuccessJoinRoomType:
		var msgSuccessJoinRoom SuccessJoinRoom
		err := json.Unmarshal(msg, &msgSuccessJoinRoom)
		if err != nil {
			return err, nil
		}
		parsedMsg = msgSuccessJoinRoom
	case FailJoinRoomType:
		var msgFailJoinRoom FailJoinRoom
		err := json.Unmarshal(msg, &msgFailJoinRoom)
		if err != nil {
			return err, nil
		}
		parsedMsg = msgFailJoinRoom
	case UnableToParseMessageType:
		var msgUnableToParse UnableToParse
		err := json.Unmarshal(msg, &msgUnableToParse)
		if err != nil {
			return err, nil
		}
		parsedMsg = msgUnableToParse
	}
	return parsedMsg, nil
}
