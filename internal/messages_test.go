package internal

import (
	"testing"
)

func TestParseJoinMessage(t *testing.T) {
	rawMessage := []byte("{\"messageType\":\"joinRoom\",\"message\":{\"chatRoomId\":1}}")
	msg, err := ParseClientMessages(rawMessage)
	if err != nil {
		t.Error("Unable to parse Join message.\n")
	}
	_, ok := msg.(JoinRoom)
	if !ok {
		t.Error("Message is not of type Join.\n")
	}
}

func TestParseCreateRoomMessage(t *testing.T) {
	rawMessage := []byte("{\"messageType\":\"createRoom\",\"message\":{}}")
	msg, err := ParseClientMessages(rawMessage)
	if err != nil {
		t.Error("Unable to parse CreateRoom message.\n")
	}
	_, ok := msg.(CreateRoom)
	if !ok {
		t.Error("Message is not of type CreateRoom.\n")
	}
}

func TestParseTextMessage(t *testing.T) {
	rawMessage := []byte("{\"messageType\":\"text\",\"message\":{\"chatRoomId\": 1, \"text\": \"Sample text\"}}")
	msg, err := ParseClientMessages(rawMessage)
	if err != nil {
		t.Error("Unable to parse Text message.\n")
	}
	_, ok := msg.(Text)
	if !ok {
		t.Error("Message is not of type Text.\n")
	}
}

func TestNewSuccessCreateRoomMessage(t *testing.T) {
	expectedMessage := "{\"messageType\":\"successCreateRoom\",\"message\":{\"chatRoomId\":1}}"
	chatRoomId := 1
	joinMsg := string(NewSuccessCreateRoomMessage(chatRoomId))
	if joinMsg != expectedMessage {
		t.Error("Unexpected SuccessCreateRoom message structure.\n")
	}
}

func TestNewSuccessJoinRoomMessage(t *testing.T) {
	expectedMessage := "{\"messageType\":\"successJoinRoom\",\"message\":{\"chatRoomId\":1}}"
	chatRoomId := 1
	joinMsg := string(NewSuccessJoinRoomMessage(chatRoomId))
	if joinMsg != expectedMessage {
		t.Error("Unexpected SuccessJoinRoom message structure.\n")
	}
}

func TestNewFailJoinRoomTypeMessage(t *testing.T) {
	expectedMessage := "{\"messageType\":\"failJoinRoom\",\"message\":{\"chatRoomId\":1}}"
	chatRoomId := 1
	joinMsg := string(NewFailJoinRoomMessage(chatRoomId))
	if joinMsg != expectedMessage {
		t.Error("Unexpected FailJoinRoom message structure.\n")
	}
}

func TestNewUnableToParseMessage(t *testing.T) {
	expectedMessage := "{\"messageType\":\"unableToParse\",\"message\":{}}"
	joinMsg := string(NewUnableToParseMessage())
	if joinMsg != expectedMessage {
		t.Error("Unexpected UnableToParse message structure.\n")
	}
}
