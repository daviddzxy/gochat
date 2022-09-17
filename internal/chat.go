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

func (c *client) addRoomSession(roomHandle string, rs *roomSession) {
	c.roomSessions[roomHandle] = rs
}

func (c *client) removeRoomSession(roomHandle string) {
	delete(c.roomSessions, roomHandle)
}

func NewClient(id int, conn *websocket.Conn) *client {
	c := &client{id: id, conn: conn}
	c.roomSessions = make(map[string]*roomSession)
	return c
}

type roomSession struct {
	Id     int     `json:"id"`
	Handle string  `json:"handle"`
	Room   *Room   `json:"-"`
	Client *client `json:"-"`
}

func (rs *roomSession) writeMessage(m []byte) {
	err := rs.Client.conn.WriteMessage(websocket.TextMessage, m)
	if err != nil {
		log.Printf("Write message to room session failed - {sessionId: %d, clientId %d, message: %s, error: %s}",
			rs.Id,
			rs.Client.id,
			string(m),
			err,
		)
		return
	}
	log.Printf("Write message to room session - {sessionId: %d, clientId %d, message: %s}",
		rs.Id,
		rs.Client.id,
		string(m),
	)
}

type Room struct {
	handle                 string
	roomSessions           map[int]*roomSession
	roomSessionIdGenerator Generator
}

func NewChatRoom(handle string) *Room {
	r := &Room{handle: handle}
	r.roomSessions = make(map[int]*roomSession)
	return r
}

func (r *Room) isEmpty() bool {
	if len(r.roomSessions) == 0 {
		return true
	}
	return false
}

func (r *Room) addRoomSession(rs *roomSession) {
	r.roomSessions[rs.Id] = rs
}

func (r *Room) removeRoomSession(rs *roomSession) {
	delete(r.roomSessions, rs.Id)
}

func (r *Room) broadcastMessage(m []byte) {
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
	Host      string
	Port      string
	clients   map[int]*client
	chatRooms map[string]*Room
	onConnect chan *websocket.Conn
	onClose   chan *clientMessage
	onMessage chan *clientMessage
	upgrader  *websocket.Upgrader
}

func NewChatServer(host string, port string) *ChatServer {
	cs := &ChatServer{Host: host, Port: port}
	cs.chatRooms = make(map[string]*Room)
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

func (cs *ChatServer) terminateAllClientSessions(id int) {
	c := cs.clients[id]
	for _, rs := range c.roomSessions {
		r := rs.Room
		rs.Room.removeRoomSession(rs)
		c.removeRoomSession(rs.Room.handle)
		r.broadcastMessage(NewRoomSessionPart(r.handle, rs.Id))
		log.Printf(
			"Client left room - {clientId; %d, roomSessionId: %d, roomHandle: %s}",
			c.id,
			rs.Id,
			r.handle,
		)
		if r.isEmpty() {
			delete(cs.chatRooms, r.handle)
			log.Printf(
				"Room was destroyed - {roomHandle: %s}",
				r.handle,
			)
		}
	}
}

func (cs *ChatServer) Run() {
	url := cs.Host + ":" + cs.Port
	go func() {
		http.HandleFunc("/", cs.connectionRequestHandler)
		if err := http.ListenAndServe(url, nil); err != http.ErrServerClosed {
			log.Fatalf("Start of web socket server failed - {error: %s}", err)
		}
	}()
	log.Printf("Chat server is listening - {address: %s}", url)
	clientIdGenerator := Generator{}
	for {
		select {
		case conn := <-cs.onConnect:
			cs.openClient(clientIdGenerator.generateId(), conn)
		case clientMsg := <-cs.onClose:
			cs.terminateAllClientSessions(clientMsg.clientId)
			cs.closeClient(clientMsg.clientId)
		case clientMsg := <-cs.onMessage:
			log.Printf(
				"New message received from Client - {clientId: %d, message: %s}",
				clientMsg.clientId,
				string(clientMsg.rawMessage),
			)
			msg, err := ParseClientMessages(clientMsg.rawMessage)
			if err != nil {
				log.Printf("Failed to parse message - {nessage: %s, error: %s}", clientMsg.rawMessage, err)
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
			Id:     r.roomSessionIdGenerator.generateId(),
			Handle: msg.SessionHandle,
			Room:   r,
			Client: c,
		}
		r.addRoomSession(rs)
		c.addRoomSession(r.handle, rs)
		rs.writeMessage(NewSuccessJoin(r.handle, rs.Id, r.roomSessions))
		r.broadcastMessage(NewRoomSessionJoin(r.handle, rs))
		log.Printf("Client joined room - {clientId: %d, roomSessionId: %d, roomHandle: %s}",
			c.id,
			rs.Id,
			r.handle,
		)
	} else {
		log.Printf("Client already in room - {clientId; %d, roomSessionId: %d, roomHandle: %s}",
			c.id,
			rs.Id,
			r.handle,
		)
	}
}

func (cs *ChatServer) handlePartMessage(msg Part, c *client) {
	rs := c.roomSessions[msg.RoomHandle]
	if rs != nil {
		r := rs.Room
		r.removeRoomSession(rs)
		c.removeRoomSession(r.handle)
		log.Printf(
			"Client left room - {clientId; %d, roomSessionId: %d, roomHandle: %s}",
			c.id,
			rs.Id,
			r.handle,
		)
		rs.writeMessage(NewSuccessPart(r.handle))
		r.broadcastMessage(NewRoomSessionPart(r.handle, rs.Id))
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
		rs.Room.broadcastMessage(NewReceiveTextMessage(msg.Content, msg.RoomHandle, rs.Id))
	}
}
