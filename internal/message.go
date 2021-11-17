package internal

import "encoding/json"

const (
	JoinType string = "join"
)

type Envelope struct {
	Type string `json:"messageType"`
	Msg interface{} `json:"message"`
}

type JoinMessage struct {
	ChatRoomId int `json:"chatRoomId"`
}

func NewJoinMessage(chatRoomId int) []byte {
	msg := &Envelope{Type: JoinType}
	msg.Msg = &JoinMessage{ChatRoomId: chatRoomId}
	jsonMsg, _ := json.Marshal(msg)
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
		var joinMsg JoinMessage
		err := json.Unmarshal(msg, &joinMsg)
		if err != nil {
			return err, nil
		}
		parsedMsg = joinMsg
	}
	return parsedMsg, nil
}


