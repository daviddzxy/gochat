package internal

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type clientMessage struct {
	clientId   int
	rawMessage []byte
}

type client struct {
	id           int
	conn         *websocket.Conn
	roomSessions map[string]*roomSession
}

func NewClient(id int, conn *websocket.Conn) *client {
	c := &client{id: id, conn: conn}
	c.roomSessions = make(map[string]*roomSession)
	return c
}

type roomSession struct {
	id     int
	handle string
	room   *room
	client *client
}

func (rs *roomSession) writeMessage(m []byte) {
	err := rs.client.conn.WriteMessage(websocket.TextMessage, m)
	if err != nil {
		log.Printf("Write message to room session failed - {sessionId: %d, clientId %d, message: %s}",
			rs.id,
			rs.client.id,
			string(m),
		)
		return
	}
	log.Printf("Write message to room session - {sessionId: %d, clientId %d, message: %s}",
		rs.id,
		rs.client.id,
		string(m),
	)
}

type room struct {
	handle                 string
	roomSessions           map[int]*roomSession
	roomSessionIdGenerator Generator
}

func NewChatRoom(handle string) *room {
	r := &room{handle: handle}
	r.roomSessions = make(map[int]*roomSession)
	return r
}

func (r *room) isEmpty() bool {
	if len(r.roomSessions) == 0 {
		return true
	}
	return false
}

func (r *room) addRoomSession(rs *roomSession) {
	r.roomSessions[rs.id] = rs
}

func (r *room) removeRoomSession(rs *roomSession) {
	delete(r.roomSessions, rs.id)
}

func (r *room) broadcastMessage(m []byte) {
	wg := sync.WaitGroup{}
	wg.Add(len(r.roomSessions))
	for _, rs := range r.roomSessions {
		go func(rs *roomSession) {
			defer wg.Done()
			rs.writeMessage(m)
		}(rs)
	}
	wg.Wait()
}

type ChatServer struct {
	Address   string
	Pattern   string
	clients   map[int]*client
	chatRooms map[string]*room
	onConnect chan *websocket.Conn
	onClose   chan *clientMessage
	onMessage chan *clientMessage
	upgrader  *websocket.Upgrader
}

func NewChatServer(address string, pattern string) *ChatServer {
	cs := &ChatServer{Address: address, Pattern: pattern}
	cs.chatRooms = make(map[string]*room)
	cs.clients = make(map[int]*client)
	cs.onConnect = make(chan *websocket.Conn)
	cs.onClose = make(chan *clientMessage)
	cs.onMessage = make(chan *clientMessage)
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

func (cs *ChatServer) readFromConnection(c *client) {
	for {
		_, p, err := c.conn.ReadMessage()
		msg := &clientMessage{clientId: c.id, rawMessage: p}
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
	log.Printf("New connection established - {clientId: %d}", c.id)
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
			log.Fatalf("Start of web socket server failed - {error: %s}", err)
		}
	}()
	log.Printf("Chat server is listening - {address: %s}", cs.Address)
	clientIdGenerator := Generator{}
	for {
		select {
		case conn := <-cs.onConnect:
			cs.openClient(clientIdGenerator.generateId(), conn)
		case clientMsg := <-cs.onClose:
			cs.closeClient(clientMsg.clientId)
		case clientMsg := <-cs.onMessage:
			log.Printf(
				"New message received from client - {clientId: %d, message: %s}",
				clientMsg.clientId,
				string(clientMsg.rawMessage),
			)
			msg, err := ParseClientMessages(clientMsg.rawMessage)
			if err != nil {
				log.Printf("Failed to parse message - {nessage: %s}", clientMsg.rawMessage)
			}
			c := cs.clients[clientMsg.clientId]
			switch msg.Type {
			case JoinType:
				joinData, _ := msg.Data.(Join)
				cs.handleJoinMessage(joinData, c)
			case PartType:
				partData, _ := msg.Data.(Part)
				cs.handlePartMessage(partData, c)
			case TextType:
				textData, _ := msg.Data.(Text)
				cs.handleTextMessage(textData, c)
			}
		}
	}
}

func (cs *ChatServer) handleJoinMessage(msg Join, c *client) {
	r := cs.chatRooms[msg.RoomHandle]
	if r == nil {
		r = NewChatRoom(msg.RoomHandle)
		cs.chatRooms[msg.RoomHandle] = r
		log.Printf(
			"Room was created - {roomhandle: %s}",
			r.handle,
		)
	}
	rs := c.roomSessions[msg.RoomHandle]
	if rs == nil {
		rs := &roomSession{
			id:     r.roomSessionIdGenerator.generateId(),
			handle: msg.SessionHandle,
			room:   r,
			client: c,
		}
		r.addRoomSession(rs)
		c.roomSessions[r.handle] = rs
		rs.writeMessage(NewSuccessJoin(r.handle, rs.id))
		log.Printf("Client joined room - {clientId: %d, roomSessionId: %d, roomHandle: %s}",
			c.id,
			rs.id,
			r.handle,
		)
	} else {
		log.Printf("Client already in room - {clientId; %d, roomSessionId: %d, roomHandle: %s}",
			c.id,
			rs.id,
			r.handle,
		)
	}
}

func (cs *ChatServer) handlePartMessage(msg Part, c *client) {
	rs := c.roomSessions[msg.RoomHandle]
	if rs != nil {
		r := rs.room
		r.removeRoomSession(rs)
		log.Printf(
			"Client left room - {clientId; %d, roomSessionId: %d, roomHandle: %s}",
			c.id,
			rs.id,
			r.handle,
		)
		rs.writeMessage(NewSuccessPart(r.handle))
		if r.isEmpty() {
			delete(cs.chatRooms, r.handle)
			log.Printf(
				"Room was destroyed - {roomHandle: %s}",
				r.handle,
			)
		}
	}
}

func (cs *ChatServer) handleTextMessage(msg Text, c *client) {
	rs := c.roomSessions[msg.RoomHandle]
	if rs != nil {
		rs.room.broadcastMessage(NewReceiveTextMessage(msg.Content, msg.RoomHandle, rs.id))
	}
}
