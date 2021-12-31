package internal

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var idRoomGenerator = Generator{}

type Room struct {
	id      int
	clients map[int]*Client
}

func NewRoom() *Room {
	r := &Room{id: idRoomGenerator.generateId()}
	r.clients = make(map[int]*Client)
	return r
}

type ClientMessage struct {
	clientId   int
	rawMessage []byte
}

type Client struct {
	id   int
	conn *websocket.Conn
}

var idClientGenerator = Generator{}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{id: idClientGenerator.generateId(), conn: conn}
}

type ChatServer struct {
	Address   string
	Pattern   string
	clients   map[int]*Client
	onConnect chan *websocket.Conn
	onClose   chan *ClientMessage
	onMessage chan *ClientMessage
	chatRooms map[int]*Room
	upgrader  *websocket.Upgrader
}

func (cs *ChatServer) Run() {
	cs.clients = make(map[int]*Client)
	cs.chatRooms = make(map[int]*Room)
	cs.onConnect = make(chan *websocket.Conn)
	cs.onClose = make(chan *ClientMessage)
	cs.upgrader = &websocket.Upgrader{}
	cs.upgrader.CheckOrigin = func(request *http.Request) bool { return true } // TODO: implement check origin function

	go func() {
		http.HandleFunc(cs.Pattern, cs.connectionRequestHandler)
		if err := http.ListenAndServe(cs.Address, nil); err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()
	log.Println("Chat server is listening.")
	for {
		select {
		case conn := <-cs.onConnect:
			c := NewClient(conn)
			cs.clients[c.id] = c
			log.Printf("New connection established: %d.\n", c.id)
			go cs.readFromClient(c)
		case clientMsg := <-cs.onClose:
			cs.closeClient(clientMsg.clientId)
			log.Printf("Connection %d closed.\n", clientMsg.clientId)
		case clientMsg := <-cs.onMessage:
			log.Printf("New message received from client %d.\n", clientMsg.clientId)
			msg, err := ParseClientMessages(clientMsg.rawMessage)
			client := cs.clients[clientMsg.clientId]
			if err != nil {
				cs.writeToClient(client, NewUnableToParseMessage())
				log.Printf("Unable to parse client message %s.\n", clientMsg.rawMessage)
			}
			switch m := msg.(type) {
			case Text:
				//check if client is in room
				//broadcast message to other clients
			case JoinRoom:
				chatRoomId := m.ChatRoomId
				if cs.chatRooms[chatRoomId] != nil {
					cs.chatRooms[chatRoomId].clients[client.id] = client
					cs.writeToClient(client, NewSuccessJoinRoomMessage(chatRoomId))
					log.Printf("Client %d successfully joined room %d.\n", client.id, chatRoomId)
				} else {
					cs.writeToClient(client, NewFailJoinRoomMessage(chatRoomId))
					log.Printf("Client %d tried to join not exisiting room %d.\n", client.id, chatRoomId)
				}
				log.Printf("Client %d joined room %d.\n", client.id, m.ChatRoomId)
			case CreateRoom:
				r := NewRoom()
				cs.chatRooms[r.id] = r
				cs.writeToClient(cs.clients[clientMsg.clientId], NewSuccessCreateRoomMessage(r.id))
				log.Printf("New room has been created.\n")
			}
		}
	}
}

func (cs *ChatServer) connectionRequestHandler(responseWriter http.ResponseWriter, request *http.Request) {
	conn, err := cs.upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	cs.onConnect <- conn
}

func (cs *ChatServer) readFromClient(c *Client) {
	for {
		_, p, err := c.conn.ReadMessage()
		msg := &ClientMessage{clientId: c.id, rawMessage: p}
		if err != nil {
			log.Println(err)
			cs.onClose <- msg
			return
		}
		cs.onMessage <- msg
	}
}

func (cs *ChatServer) writeToClient(c *Client, rawMessage []byte) {
	err := c.conn.WriteMessage(websocket.TextMessage, rawMessage)
	if err != nil {
		log.Printf("Unable to send message %s to client %d.\n", string(rawMessage), c.id)
	}
}

func (cs *ChatServer) closeClient(id int) {
	err := cs.clients[id].conn.Close()
	if err != nil {
		log.Println(err)
	}
	delete(cs.clients, id)
}
