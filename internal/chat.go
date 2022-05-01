package internal

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Room struct {
	name    string
	clients map[int]*Client
}

func NewRoom(name string) *Room {
	r := &Room{name: name}
	r.clients = make(map[int]*Client)
	return r
}

func (r *Room) getClientNames() []string {
	clientNames := make([]string, len(r.clients))
	index := 0
	for _, c := range r.clients {
		clientNames[index] = c.clientName
		index += 1
	}
	return clientNames
}

type ClientMessage struct {
	clientId   int
	rawMessage []byte
}

type Client struct {
	id         int
	conn       *websocket.Conn
	clientName string
}

var idClientGenerator = Generator{}
var idMessageGenerator = Generator{}

func NewClient(conn *websocket.Conn) *Client {
	c := &Client{id: idClientGenerator.generateId(), conn: conn, clientName: ""}
	return c
}

type ChatServer struct {
	Address       string
	Pattern       string
	clients       map[int]*Client
	clientsByRoom map[int]*Room
	onConnect     chan *websocket.Conn
	onClose       chan *ClientMessage
	onMessage     chan *ClientMessage
	chatRooms     map[string]*Room
	upgrader      *websocket.Upgrader
}

func NewChatServer(address string, pattern string) *ChatServer {
	cs := &ChatServer{Address: address, Pattern: pattern}
	cs.clients = make(map[int]*Client)
	cs.clientsByRoom = make(map[int]*Room)
	cs.chatRooms = make(map[string]*Room)
	cs.onConnect = make(chan *websocket.Conn)
	cs.onClose = make(chan *ClientMessage)
	cs.onMessage = make(chan *ClientMessage)
	cs.upgrader = &websocket.Upgrader{}
	cs.upgrader.CheckOrigin = func(request *http.Request) bool { return true } // TODO: implement check origin function
	return cs
}

func (cs *ChatServer) Run() {
	go func() {
		http.HandleFunc(cs.Pattern, cs.connectionRequestHandler)
		if err := http.ListenAndServe(cs.Address, nil); err != http.ErrServerClosed {
			log.Fatalf("Could not start web socket server: %s\n", err)
		}
	}()
	log.Printf("Chat server is listening on %s.\n", cs.Address)
	for {
		select {
		case conn := <-cs.onConnect:
			c := NewClient(conn)
			cs.clients[c.id] = c
			log.Printf("New connection established: %d.\n", c.id)
			go cs.readFromClient(c)
		case clientMsg := <-cs.onClose:
			_ = cs.removeClientFromRoom(clientMsg.clientId)
			cs.closeClient(clientMsg.clientId)
			log.Printf("Connection %d closed.\n", clientMsg.clientId)
		case clientMsg := <-cs.onMessage:
			log.Printf("New message %s received from client %d.\n", string(clientMsg.rawMessage), clientMsg.clientId)
			message, err := ParseClientMessages(clientMsg.rawMessage)
			if err != nil {
				log.Printf("Unable to parse client message %s.\n", clientMsg.rawMessage)
			}
			client := cs.clients[clientMsg.clientId]
			switch message.Type {
			case TextType:
				textData, _ := message.Data.(Text)
				cs.handleTextMessage(textData, client)
			case JoinType:
				joinData, _ := message.Data.(Join)
				cs.handleJoinMessage(joinData, client)
			case PartType:
				cs.handlePartMessage(client)
			}
		}
	}
}

func (cs *ChatServer) handleJoinMessage(data Join, c *Client) {
	err := cs.addClientToRoom(c.id, data.RoomName)
	if err != nil && errors.Is(err, ClientAlreadyInRoomError) {
		_ = cs.removeClientFromRoom(c.id)
		_ = cs.addClientToRoom(c.id, data.RoomName)
	}
	c.clientName = data.ClientName
	cs.writeToClient(c, NewSuccessJoinMessage(data.RoomName, cs.chatRooms[data.RoomName].getClientNames()))
	log.Printf("Client %d joined room %s with name %s.\n", c.id, data.RoomName, data.ClientName)
}

func (cs *ChatServer) handlePartMessage(c *Client) {
	room := cs.clientsByRoom[c.id].name
	if err := cs.removeClientFromRoom(c.id); err == nil {
		cs.writeToClient(c, NewSuccessPartMessage())
		log.Printf("Client %d left room %s.\n", c.id, room)
	} else {
		log.Println(err)
	}
}

func (cs *ChatServer) removeClientFromRoom(clientId int) error {
	r := cs.clientsByRoom[clientId]
	c := cs.clients[clientId]
	if r != nil {
		delete(r.clients, c.id)
		delete(cs.clientsByRoom, c.id)
		log.Printf("Client %d left room %s.\n", c.id, r.name)
		if len(r.clients) == 0 {
			delete(cs.chatRooms, r.name)
			log.Printf("Room %s was deleted.\n", r.name)
		}
		return nil
	} else {
		return ClientDoesNotBelongToAnyRoomError
	}
}

func (cs *ChatServer) addClientToRoom(clientId int, roomName string) error {
	c := cs.clients[clientId]
	r := cs.clientsByRoom[clientId]
	if r == nil {
		if cs.chatRooms[roomName] == nil {
			r = NewRoom(roomName)
			cs.chatRooms[roomName] = r
			log.Printf("New room %s has been created.\n", roomName)
		} else {
			r = cs.chatRooms[roomName]
		}
		r.clients[c.id] = c
		cs.clientsByRoom[c.id] = r
		log.Printf("Client %d joined room %s.\n", c.id, roomName)
		return nil
	} else {
		return errors.New("client is already in a room")
	}
}

func (cs *ChatServer) handleTextMessage(data Text, c *Client) {
	room := cs.clientsByRoom[c.id]
	if room != nil {
		receiveTextMessage := NewReceiveTextMessage(data.Text, c.clientName, idMessageGenerator.generateId())
		cs.broadcastMessage(room, receiveTextMessage)
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
	log.Printf("Message %s sent to client %d\n", string(rawMessage), c.id)
}

func (cs *ChatServer) broadcastMessage(r *Room, rawMessage []byte) {
	wg := sync.WaitGroup{}
	wg.Add(len(r.clients))
	for _, c := range r.clients {
		go func(c *Client) {
			defer wg.Done()
			cs.writeToClient(c, rawMessage)
		}(c)
	}
	wg.Wait()
}

func (cs *ChatServer) closeClient(id int) {
	err := cs.clients[id].conn.Close()
	if err != nil {
		log.Println(err)
	}
	delete(cs.clients, id)
}
