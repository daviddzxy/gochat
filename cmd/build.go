package main

import (
	flag "github.com/spf13/pflag"
	"tmpchat/internal"
)

func main() {
	host := flag.StringP("host", "a", "localhost", "Server address")
	port := flag.StringP("port", "p", "8080", "Server port")
	flag.Parse()
	cs := internal.NewChatServer(*host, *port)
	cs.Run()
}
