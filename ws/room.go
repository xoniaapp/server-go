package ws

import (
	"context"
	"log"
	"github.com/aelpxy/xoniaapp/model"
	"github.com/go-redis/redis/v8"
)
// Room is a struct
type Room struct {
	id         string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *model.WebsocketMessage
	redis      *redis.Client
}

var ctx = context.Background()

func NewRoom(id string, rds *redis.Client) *Room {
	return &Room{
		id:         id,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *model.WebsocketMessage),
		redis:      rds,
	}
}

func (room *Room) RunRoom() {

	go room.subscribeToRoomMessages()

	for {
		select {

		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.broadcast:
			room.publishRoomMessage(message.Encode())
		}
	}
}

func (room *Room) registerClientInRoom(client *Client) {
	room.clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *Client) {
	delete(room.clients, client)
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) GetId() string {
	return room.id
}

func (room *Room) publishRoomMessage(message []byte) {
	err := room.redis.Publish(ctx, room.GetId(), message).Err()

	if err != nil {
		log.Println(err)
	}
}

func (room *Room) subscribeToRoomMessages() {
	pubsub := room.redis.Subscribe(ctx, room.GetId())

	ch := pubsub.Channel()

	for msg := range ch {
		room.broadcastToClientsInRoom([]byte(msg.Payload))
	}
}
// LICENSE@GNUv3
