package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/aelpxy/xoniaapp/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 10000
)

var newline = []byte{'\n'}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	ID    string
	conn  *websocket.Conn
	hub   *Hub
	send  chan []byte
	rooms map[*Room]bool
}

func newClient(conn *websocket.Conn, hub *Hub, id string) *Client {
	return &Client{
		ID:    id,
		conn:  conn,
		hub:   hub,
		send:  make(chan []byte, 256),
		rooms: make(map[*Room]bool),
	}
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)

	_ = client.conn.SetReadDeadline(time.Now().Add(pongWait))

	client.conn.SetPongHandler(func(string) error {
		_ = client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
		client.handleNewMessage(jsonMessage)
	}

}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			n := len(client.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	client.hub.unregister <- client
	for room := range client.rooms {
		room.unregister <- client
	}
	close(client.send)
	_ = client.conn.Close()
}

func ServeWs(hub *Hub, ctx *gin.Context) {

	userId := ctx.MustGet("userId").(string)
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, hub, userId)

	go client.writePump()
	go client.readPump()

	hub.register <- client
}

func (client *Client) handleNewMessage(jsonMessage []byte) {

	var message model.ReceivedMessage
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
	}

	switch message.Action {
	case JoinChannelAction:
		client.handleJoinChannelMessage(message)
	case JoinGuildAction:
		client.handleJoinGuildMessage(message)
	case JoinUserAction:
		client.handleJoinRoomMessage(message)

	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)
	case LeaveGuildAction:
		client.handleLeaveGuildMessage(message)

	case StartTypingAction:
		client.handleTypingEvent(message, AddToTypingAction)
	case StopTypingAction:
		client.handleTypingEvent(message, RemoveFromTypingAction)

	case ToggleOnlineAction:
		client.toggleOnlineStatus(true)
	case ToggleOfflineAction:
		client.toggleOnlineStatus(false)

	case GetRequestCountAction:
		client.handleGetRequestCount()
	}
}

func (client *Client) handleJoinChannelMessage(message model.ReceivedMessage) {
	roomName := message.Room

	cs := client.hub.channelService
	channel, err := cs.Get(roomName)

	if err != nil {
		return
	}

	if err = cs.IsChannelMember(channel, client.ID); err != nil {
		return
	}

	client.handleJoinRoomMessage(message)
}

func (client *Client) handleJoinGuildMessage(message model.ReceivedMessage) {
	roomName := message.Room

	gs := client.hub.guildService
	guild, err := gs.GetGuild(roomName)

	if err != nil {
		return
	}

	if !isMember(guild, client.ID) {
		return
	}

	client.handleJoinRoomMessage(message)
}

func (client *Client) handleJoinRoomMessage(message model.ReceivedMessage) {
	roomName := message.Room

	room := client.hub.findRoomById(roomName)
	if room == nil {
		room = client.hub.createRoom(roomName)
	}

	client.rooms[room] = true

	room.register <- client
}

func (client *Client) handleLeaveGuildMessage(message model.ReceivedMessage) {
	_ = client.hub.guildService.UpdateMemberLastSeen(client.ID, message.Room)
	client.handleLeaveRoomMessage(message)
}

func (client *Client) handleLeaveRoomMessage(message model.ReceivedMessage) {
	room := client.hub.findRoomById(message.Room)
	delete(client.rooms, room)

	if room != nil {
		room.unregister <- client
	}
}

func (client *Client) handleGetRequestCount() {
	if room := client.hub.findRoomById(client.ID); room != nil {
		count, err := client.hub.userService.GetRequestCount(client.ID)

		if err != nil {
			return
		}

		msg := model.WebsocketMessage{
			Action: RequestCountEmission,
			Data:   count,
		}
		room.broadcast <- &msg
	}
}

func (client *Client) handleTypingEvent(message model.ReceivedMessage, action string) {
	roomID := message.Room
	if room := client.hub.findRoomById(roomID); room != nil {
		msg := model.WebsocketMessage{
			Action: action,
			Data:   message.Message,
		}
		room.broadcast <- &msg
	}
}

func (client *Client) toggleOnlineStatus(isOnline bool) {
	uid := client.ID
	us := client.hub.userService

	user, err := us.Get(uid)

	if err != nil {
		log.Printf("could not find user: %v", err)
		return
	}

	user.IsOnline = isOnline

	if err := us.UpdateAccount(user); err != nil {
		log.Printf("could not update user: %v", err)
		return
	}

	ids, err := us.GetFriendAndGuildIds(uid)

	if err != nil {
		log.Printf("could not find ids: %v", err)
		return
	}

	action := ToggleOfflineEmission
	if isOnline {
		action = ToggleOnlineEmission
	}

	for _, id := range *ids {
		if room := client.hub.findRoomById(id); room != nil {
			msg := model.WebsocketMessage{
				Action: action,
				Data:   uid,
			}
			room.broadcast <- &msg
		}
	}
}

func isMember(guild *model.Guild, userId string) bool {
	for _, v := range guild.Members {
		if v.ID == userId {
			return true
		}
	}
	return false
}
