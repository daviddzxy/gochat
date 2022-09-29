package main

import (
	flag "github.com/spf13/pflag"
	"tmpchat/internal"
)

func main() {
	port := flag.StringP("port", "p", "8080", "Server port")
	flag.Parse()
	cs := internal.NewChatServer(*port)
	cs.Run()
}
