package internal

import (
	"testing"
)

func TestParseJoinMessage(t *testing.T) {
	rawMessage := []byte("{\"type\":\"JOIN_ROOM\",\"data\":{\"roomName\":\"text\", \"clientName\": \"text\"}}")
	msg, err := ParseClientMessages(rawMessage)
	if err != nil {
		t.Error("Unable to parse Join message.\n")
	}
	_, ok := msg.(JoinRoom)
	if !ok {
		t.Error("Message is not of type Join.\n")
	}
}

func TestParseTextMessage(t *testing.T) {
	rawMessage := []byte("{\"type\":\"TEXT\",\"data\":{\"chatRoomId\": 1, \"text\": \"text\"}}")
	msg, err := ParseClientMessages(rawMessage)
	if err != nil {
		t.Error("Unable to parse Text message.\n")
	}
	_, ok := msg.(Text)
	if !ok {
		t.Error("Message is not of type Text.\n")
	}
}

func TestNewSuccessJoinRoomMessage(t *testing.T) {
	expectedMessage := "{\"type\":\"SUCCESS_JOIN_ROOM\",\"data\":{\"roomName\":\"text\"}}"
	roomName := "text"
	joinMsg := string(NewSuccessJoinRoomMessage(roomName))
	if joinMsg != expectedMessage {
		t.Error("Unexpected SuccessJoinRoom message structure.\n")
	}
}

func TestNewNewClientNamesMessage(t *testing.T) {
	expectedMessage := "{\"type\":\"CLIENT_LIST\",\"data\":{\"clientNames\":[\"client1\",\"client2\"]}}"
	clientNames := []string{"client1", "client2"}
	joinMsg := string(NewClientNamesMessage(clientNames))
	if joinMsg != expectedMessage {
		t.Error("Unexpected SuccessJoinRoom message structure.\n")
	}
}
