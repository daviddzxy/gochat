package internal

import (
	"testing"
	"time"
)

func SetUpServerClient() (*ChatServer, *chatClient){
	s := &ChatServer{Address: "localhost:8080", Pattern:"/chat"}
	c := &chatClient{Address: "localhost:8080", Pattern:"/chat"}
	go s.Run()
	maxAttempts := 3
	currAttempt := 0
	for c.conn == nil && currAttempt <= maxAttempts{
		c.Connect()
		time.Sleep(100 * time.Millisecond)
		currAttempt += 1
	}
	return s, c
}

func TestEstablishAndCloseConnection(t *testing.T) {
	s, c := SetUpServerClient()
	if len(s.clients) != 1 {
		t.Error("Server does not contain a connection.")
	}
	c.Close()

	maxAttempts := 3
	currAttempt := 0
	for len(s.clients) > 0 && currAttempt <= maxAttempts {
		time.Sleep(100 * time.Millisecond)
		currAttempt += 1
	}
	if len(s.clients) != 0 {
		t.Error("Connection was not successfully closed.")
	}
}