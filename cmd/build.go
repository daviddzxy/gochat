package main

import (
	"tmpchat/internal"
)

func main() {
	s := internal.ChatServer{Address:"localhost:8080", Pattern:"/"}
	s.Run()
}
