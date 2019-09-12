package handlers

import (
	"encoding/json"
	"final-project-crew/videoCharade/servers/gateway/sessions"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/streadway/amqp"

	"github.com/gorilla/websocket"
)

//MSG is a struct that stores a list of WSClients
type MSG struct {
	Type     string      `json:"type"`
	Data     interface{} `json:"data"`
	UserList []int64     `json:"userList"`
}

//Notifier is a struct that stores a list of WSClients
type Notifier struct {
	connectionList map[int64]*websocket.Conn
	lock           sync.Mutex
}

//store new struct for user list
//- userObject
//- connection

//NewNotifier Returns a new notifer struct
func NewNotifier() *Notifier {
	return &Notifier{
		connectionList: make(map[int64]*websocket.Conn),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		if r.Header.Get("Origin") != "api.videocharade.me" {
			return true
		}
		return false
	},
}

//WebSocketConnectionHandler handles incoming websocket requests
func (ctx *Context) WebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	//handle the websocket handshake
	var state SessionState
	_, err := sessions.GetState(r, ctx.SigningKey, ctx.SessionStore, &state)
	if err != nil {
		http.Error(w, "Status Unauthorized", http.StatusUnauthorized)
		return
	}
	user := state.User

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket connection", 401)
		return
	}
	fmt.Println("SOCKET ADDED")
	log.Printf("SOCKET ADDED")
	ctx.InsertConnection(user.ID, conn)

	// Invoke a goroutine for handling control messages from this connection
	go (func(conn *websocket.Conn) {
		defer conn.Close()
		for {

			if _, _, err := conn.NextReader(); err != nil {
				break
			}
		}
	})(conn)
}

// InsertConnection Adds a new connection to the notifier struct
func (ctx *Context) InsertConnection(id int64, conn *websocket.Conn) {
	ctx.Notifier.lock.Lock()
	fmt.Println("connected to UserID: ", id)
	ctx.Notifier.connectionList[id] = conn
	ctx.Notifier.lock.Unlock()
}

// RemoveConnection removes a connection from the notifier stuct
func (ctx *Context) RemoveConnection(id int64) {
	ctx.Notifier.lock.Lock()
	ctx.Notifier.connectionList[id].Close()
	delete(ctx.Notifier.connectionList, id)
	ctx.Notifier.lock.Unlock()
}

// WriteToAllConnections is a
// Simple method for writing a message to all live connections.
// In your homework, you will be writing a message to a subset of connections
// (if the message is intended for a private channel), or to all of them (if the message
// is posted on a public channel
func (ctx *Context) WriteToAllConnections(message interface{}, idList []int64) error {
	var writeError error

	if len(idList) < 1 {
		for _, conn := range ctx.Notifier.connectionList {
			writeError = conn.WriteJSON(message)
			if writeError != nil {
				return writeError
			}
		}
		return nil
	} else {
		for _, id := range idList {
			if _, ok := ctx.Notifier.connectionList[id]; ok {
				writeError = ctx.Notifier.connectionList[id].WriteJSON(message)
				if writeError != nil {
					return writeError
				}
			}
		}
		return nil
	}
}

//Process converts channels into variables and iterates through it
func (ctx *Context) Process(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		message := &MSG{}
		err := json.Unmarshal([]byte(msg.Body), message)
		if err != nil {
			log.Printf("had a problem processing message queue. Due to %s", err)
		}
		fmt.Println(message)
		err = ctx.WriteToAllConnections(message, message.UserList)
		if err != nil {
			log.Printf("had a problem broadcasting. Due to %s", err)
		}
	}
}
