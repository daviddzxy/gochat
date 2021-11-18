package internal

import "encoding/json"

const (
	CreatRoomType string = "createRoom"
	JoinType = "join"
)

type Envelope struct {
	Type string      `json:"messageType"`
	Msg  interface{} `json:"message"`
}

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

func ParseMessage(input []byte) (interface{}, error) {
	var msg json.RawMessage
	env := Envelope{Msg: &msg}
	err := json.Unmarshal(input, &env)
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
	case CreatRoomType:
		var createRoomMsg CreateRoom
		err := json.Unmarshal(msg, &createRoomMsg)
		if err != nil {
			return err, nil
		}
		parsedMsg = createRoomMsg
	}
	return parsedMsg, nil
}
