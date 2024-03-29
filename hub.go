package gExcel

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mayur-tolexo/gExcel/message"
	"github.com/tidwall/gjson"
)

var messages map[string]message.EData

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func init() {
	messages = make(map[string]message.EData)
}

type Hub struct {
	clients    []*Client
	nextID     int
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		nextID:     0,
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

func (hub *Hub) broadcast(message interface{}, ignore *Client) {
	data, _ := json.Marshal(message)
	for _, c := range hub.clients {
		if c != ignore {
			c.outbound <- data
		}
	}
}

func (hub *Hub) send(message interface{}, client *Client) {
	data, _ := json.Marshal(message)
	client.outbound <- data
}

func (hub *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "could not upgrade", http.StatusInternalServerError)
		return
	}
	client := newClient(hub, socket)
	hub.register <- client

	go client.read()
	go client.write()
}

func (hub *Hub) onConnect(client *Client) {
	log.Println("client connected: ", client.socket.RemoteAddr())

	// Make new client
	client.id = hub.nextID
	hub.nextID++
	client.color = generateColor()
	hub.clients = append(hub.clients, client)

	// Make list of all users
	users := []message.User{}
	users = append(users, message.User{ID: client.id, Color: client.color})
	for _, c := range hub.clients {
		users = append(users, message.User{ID: c.id, Color: c.color})
	}

	// Notify that a user joined
	hub.send(message.NewConnected(client.id, client.color, users, messages), client)
	hub.broadcast(message.NewUserJoined(client.id, client.color), client)
}

func (hub *Hub) onDisconnect(client *Client) {
	log.Println("client disconnected: ", client.socket.RemoteAddr())

	client.close()

	// Find index of client
	i := -1
	for j, c := range hub.clients {
		if c.id == client.id {
			i = j
			break
		}
	}
	// Delete client from list
	copy(hub.clients[i:], hub.clients[i+1:])
	hub.clients[len(hub.clients)-1] = nil
	hub.clients = hub.clients[:len(hub.clients)-1]

	hub.broadcast(message.NewUserLeft(client.id), nil)
}

func (hub *Hub) onMessage(data []byte, client *Client) {
	log.Println("onMessage: ", string(data))

	kind := gjson.GetBytes(data, "kind").Int()
	if kind == message.KindStroke {
		var msg message.Stroke
		if err := json.Unmarshal(data, &msg); err != nil {
			fmt.Println(err)
			return
		}
		msg.UserID = client.id
		for _, curMsg := range msg.EDatas {
			messages[curMsg.ID] = curMsg
		}
		hub.broadcast(msg, client)
	} else {
		var msg message.Clear
		if json.Unmarshal(data, &msg) != nil {
			return
		}
		msg.UserID = client.id
		hub.broadcast(msg, client)
	}
}
