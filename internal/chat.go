package internal

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

type ChatRoom struct {
	id      int
	clients []*websocket.Conn
}

type ChatMessage struct {
	clientId int
	payload  []byte
}

type ChatServer struct {
	Address   string
	Pattern   string
	conns     map[int]*websocket.Conn
	onConnect chan *websocket.Conn
	onClose   chan *ChatMessage
	onMessage chan *ChatMessage
	chatRooms []*ChatRoom
}

var upgrader = websocket.Upgrader{}
var idConnGenerator = Generator{}
var idRoomGenerator = Generator{}

func (cs *ChatServer) Run() {
	cs.conns = make(map[int]*websocket.Conn)
	cs.onConnect = make(chan *websocket.Conn)
	cs.onClose = make(chan *ChatMessage)

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
			id := idConnGenerator.generateId()
			cs.addConnection(conn, id)
			log.Printf("New connection established: %d.\n", id)
			go cs.readFromConn(conn, id)
		case msg := <-cs.onClose:
			cs.closeAndRemoveConnection(msg.clientId)
			log.Printf("Connection %d closed.\n", msg.clientId)
		case msg := <-cs.onMessage:
			log.Printf("New message received from client %d.\n", msg.clientId)
		}
	}
}

func (cs *ChatServer) connectionRequestHandler(responseWriter http.ResponseWriter, request *http.Request) {
	upgrader.CheckOrigin = func(request *http.Request) bool { return true } // TODO: implement check origin function
	conn, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Println(err)
	}
	cs.onConnect <- conn
}

func (cs *ChatServer) readFromConn(conn *websocket.Conn, id int) {
	for {
		_, p, err := conn.ReadMessage()
		msg := &ChatMessage{clientId: id, payload: p}
		if err != nil {
			log.Println(err)
			cs.onClose <- msg
			return
		}
		cs.onMessage <- msg
	}
}

func (cs *ChatServer) addConnection(conn *websocket.Conn, id int) {
	cs.conns[id] = conn
}

func (cs *ChatServer) closeAndRemoveConnection(id int) {
	err := cs.conns[id].Close()
	if err != nil {
		log.Println(err)
	}
	delete(cs.conns, id)
}

// chatClient is used for testing purposes only
type chatClient struct {
	Address string
	Pattern string
	conn    *websocket.Conn
}

func (cc *chatClient) Connect() {
	u := url.URL{Scheme: "ws", Host: cc.Address, Path: cc.Pattern}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Connected to the chat server.")
	cc.conn = c
}

func (cc *chatClient) SendMessage(message string) {
	err := cc.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println(err)
	}
}

func (cc *chatClient) SendRawMessage(message []byte) {
	err := cc.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println(err)
	}
}

func (cc *chatClient) JoinChatRoom(roonId int) {
	msg := NewJoinMessage(roonId)
	cc.SendRawMessage(msg)
}

func (cc *chatClient) Close() {
	err := cc.conn.Close()
	if err != nil {
		log.Println(err)
	}
}
