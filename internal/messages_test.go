package internal

import (
	"testing"
)


func TestParseJoinMessage(t *testing.T) {
	rawMessage := []byte("{\"messageType\":\"join\",\"message\":{\"chatRoomId\":1}}")
	msg, err := ParseClientMessages(rawMessage)
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
	msg, err := ParseClientMessages(rawMessage)
	if err != nil{
		t.Error("Unable to parse CreateRoom message.")
	}
	_, ok := msg.(CreateRoom)
	if !ok {
		t.Error("Message is not of type CreateRoom.")
	}
}

func TestNewSuccessCreateRoomMessage(t *testing.T) {
	expectedMessage := "{\"messageType\":\"successCreateRoom\",\"message\":{\"chatRoomId\":1}}"
	chatRoomId := 1
	joinMsg := string(NewSuccessCreateRoomMessage(chatRoomId))
	if joinMsg != expectedMessage{
		t.Error("Unexpected SuccessCreateRoom message structure.")
	}
}

func TestNewSuccessJoinRoomMessage(t *testing.T){
	expectedMessage := "{\"messageType\":\"successJoinRoom\",\"message\":{\"chatRoomId\":1}}"
	chatRoomId := 1
	joinMsg := string(NewSuccessJoinRoomMessage(chatRoomId))
	if joinMsg != expectedMessage{
		t.Error("Unexpected SuccessJoinRoom message structure.")
	}
}

func TestNewFailJoinRoomTypeMessage(t *testing.T) {
	expectedMessage := "{\"messageType\":\"failJoinRoom\",\"message\":{\"chatRoomId\":1}}"
	chatRoomId := 1
	joinMsg := string(NewFailJoinRoomMessage(chatRoomId))
	if joinMsg != expectedMessage{
		t.Error("Unexpected FailJoinRoom message structure.")
	}
}

func TestNewUnableToParseMessage(t *testing.T) {
	expectedMessage := "{\"messageType\":\"unableToParse\",\"message\":{}}"
	joinMsg := string(NewUnableToParseMessage())
	if joinMsg != expectedMessage{
		t.Error("Unexpected UnableToParse message structure.")
	}
}
