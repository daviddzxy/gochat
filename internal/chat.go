package internal

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

type ChatServer struct {
	Address string
	Pattern string
	conns []*websocket.Conn
	newConn chan *websocket.Conn
}

var upgrader = websocket.Upgrader{}

func (chatServer *ChatServer) Run() {
	chatServer.newConn = make(chan *websocket.Conn)

	go func() {
		http.HandleFunc(chatServer.Pattern, chatServer.connectionRequestHandler)
		if err := http.ListenAndServe(chatServer.Address, nil); err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()
	log.Println("Chat server is listening.")
	for {
		select {
		case conn := <- chatServer.newConn:
			chatServer.conns = append(chatServer.conns, conn)
			log.Println("New connection")
		}
	}
}

func (chatServer *ChatServer) connectionRequestHandler(responseWriter http.ResponseWriter, request *http.Request) {
	upgrader.CheckOrigin = func(request *http.Request) bool { return  true } // TODO: implement check origin function
	conn, err := upgrader.Upgrade(responseWriter, request, nil)

	if err != nil {
		log.Println(err)
	}
	chatServer.newConn <- conn
	// TODO: read from connection
}


// chatClient is used for testing purposes only
type chatClient struct {
	Address string
	Pattern string
	conn *websocket.Conn
}

func (chatClient *chatClient) Connect() {
	u := url.URL{Scheme: "ws", Host: chatClient.Address, Path: chatClient.Pattern}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Connected to the chat server.")
	chatClient.conn = c
}

func (chatClient *chatClient) SendMessage(message string){
	err := chatClient.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println(err)
	}
}

func (chatClient *chatClient) Close() {
	err := chatClient.conn.Close()
	if err != nil {
		log.Println(err)
	}
}