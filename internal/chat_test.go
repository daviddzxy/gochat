package internal

import (
	"testing"
	"time"
)

func TestChatServerConnectionHandler(t *testing.T) {
	s := ChatServer{Address: "localhost:8080", Pattern:"/chat"}
	c := chatClient{Address: "localhost:8080", Pattern:"/chat"}
	go s.Run()
	maxAttempts := 3
	currAttempt := 0
	for c.conn == nil && currAttempt <= maxAttempts{
		c.Connect()
		time.Sleep(100 * time.Millisecond)
		currAttempt += 1
	}
	if len(s.conns) != 1 {
		t.Error("Server does not contain a connection.")
	}
}