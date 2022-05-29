package internal

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type ClientMessage struct {
	clientId   int
	rawMessage []byte
}

type Client struct {
	id           int
	conn         *websocket.Conn
	roomSessions map[string]*RoomSession
}

func NewClient(id int, conn *websocket.Conn) *Client {
	c := &Client{id: id, conn: conn}
	c.roomSessions = make(map[string]*RoomSession)
	return c
}

type RoomSession struct {
	id     int
	handle string
	room   *room
}

type room struct {
	handle                 string
	roomSessions           map[int]*RoomSession
	roomSessionIdGenerator Generator
}

func NewChatRoom(handle string) *room {
	r := &room{handle: handle}
	r.roomSessions = make(map[int]*RoomSession)
	return r
}

func (chatRoom *room) createRoomSession(handle string) *RoomSession {
	return &RoomSession{id: chatRoom.roomSessionIdGenerator.generateId(), handle: handle, room: chatRoom}
}

type ChatServer struct {
	Address   string
	Pattern   string
	clients   map[int]*Client
	chatRooms map[string]*room
	onConnect chan *websocket.Conn
	onClose   chan *ClientMessage
	onMessage chan *ClientMessage
	upgrader  *websocket.Upgrader
}

func NewChatServer(address string, pattern string) *ChatServer {
	cs := &ChatServer{Address: address, Pattern: pattern}
	cs.chatRooms = make(map[string]*room)
	cs.clients = make(map[int]*Client)
	cs.onConnect = make(chan *websocket.Conn)
	cs.onClose = make(chan *ClientMessage)
	cs.onMessage = make(chan *ClientMessage)
	cs.upgrader = &websocket.Upgrader{}
	cs.upgrader.CheckOrigin = func(request *http.Request) bool { return true } // TODO: implement check origin function
	return cs
}

func (cs *ChatServer) connectionRequestHandler(responseWriter http.ResponseWriter, request *http.Request) {
	conn, err := cs.upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	cs.onConnect <- conn
}

func (cs *ChatServer) readFromConnection(c *Client) {
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

func (cs *ChatServer) openClient(id int, conn *websocket.Conn) {
	c := NewClient(id, conn)
	cs.clients[c.id] = c
	go cs.readFromConnection(c)
	log.Printf("New connection established: %d.\n", c.id)
}

func (cs *ChatServer) closeClient(id int) {
	err := cs.clients[id].conn.Close()
	if err != nil {
		log.Println(err)
	}
	delete(cs.clients, id)
}

func (cs *ChatServer) Run() {
	go func() {
		http.HandleFunc(cs.Pattern, cs.connectionRequestHandler)
		if err := http.ListenAndServe(cs.Address, nil); err != http.ErrServerClosed {
			log.Fatalf("Could not start web socket server: %s\n", err)
		}
	}()
	log.Printf("Chat server is listening on %s.\n", cs.Address)
	clientIdGenerator := Generator{}
	for {
		select {
		case conn := <-cs.onConnect:
			cs.openClient(clientIdGenerator.generateId(), conn)
		case clientMsg := <-cs.onClose:
			cs.closeClient(clientMsg.clientId)
		case clientMsg := <-cs.onMessage:
			log.Printf("New message %s received from client %d.\n", string(clientMsg.rawMessage), clientMsg.clientId)
			msg, err := ParseClientMessages(clientMsg.rawMessage)
			if err != nil {
				log.Printf("Unable to parse client message %s.\n", clientMsg.rawMessage)
			}
			c := cs.clients[clientMsg.clientId]
			switch msg.Type {
			case JoinType:
				joinData, _ := msg.Data.(Join)
				cs.handleJoinMessage(joinData, c)
			}
		}
	}
}

func (cs *ChatServer) handleJoinMessage(msg Join, c *Client) {
	r := cs.chatRooms[msg.RoomHandle]
	if r == nil {
		r = NewChatRoom(msg.RoomHandle)
	}
	rs := c.roomSessions[msg.RoomHandle]
	if rs == nil {
		rs := &RoomSession{id: r.roomSessionIdGenerator.generateId(), handle: msg.NewRoomSessionHandle, room: r}
		r.roomSessions[rs.id] = rs
		c.roomSessions[r.handle] = rs
		log.Printf("Client %d joined room %s, with room session id %d, handle %s",
			c.id,
			r.handle,
			rs.id,
			rs.handle,
		)
	} else {
		log.Printf("Client %d already in room %s with room session id %d, handle %s",
			c.id,
			r.handle,
			rs.id,
			rs.handle,
		)
	}
}
