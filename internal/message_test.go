package internal

import (
	"testing"
)

func TestNewJoinMessage(t *testing.T) {
	expectedMessage := "{\"messageType\":\"join\",\"message\":{\"chatRoomId\":1}}"
	chatRoomId := 1
	joinMsg := string(NewJoinMessage(chatRoomId))
	if joinMsg != expectedMessage{
		t.Error("Unexpected Join message structure.")
	}
}

func TestParseJoinMessage(t *testing.T) {
	rawMessage := []byte("{\"messageType\":\"join\",\"message\":{\"chatRoomId\":1}}")
	msg, err := ParseMessage(rawMessage)
	if err != nil{
		t.Error("Unable to parse Join message.")
	}
	_, ok := msg.(Join)
	if !ok {
		t.Error("Message is not of type Join.")
	}
}

func TestParseCreateRoomMessage(t *testing.T) {
	rawMessage := []byte("{\"messageType\":\"createRoom\",\"message\":{}}")
	msg, err := ParseMessage(rawMessage)
	if err != nil{
		t.Error("Unable to parse CreateRoom message.")
	}
	_, ok := msg.(CreateRoom)
	if !ok {
		t.Error("Message is not of type CreateRoom.")
	}
}