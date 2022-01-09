package main

import (
	"tmpchat/internal"
)

func main() {
	cs := internal.NewChatServer("localhost:8080", "/")
	cs.Run()
}
