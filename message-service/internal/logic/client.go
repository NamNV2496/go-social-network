package logic

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func CreateNewClient(req *gin.Context) (*Client, error) {

	conn, err := upgrader.Upgrade(req.Writer, req.Request, nil)
	if err != nil {
		req.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, nil
	}

	roomId := req.Param("roomId")
	username := req.Param("username")
	id := uuid.New()
	return &Client{
		Conn:     conn,
		Message:  make(chan *Message),
		ID:       id.String(),
		RoomID:   roomId,
		Username: username,
	}, nil
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}
		c.Conn.WriteJSON(message)
	}
}

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}

		room := hub.Rooms[msg.RoomID]
		room.LastMsg = msg.Username + ": " + msg.Content
		hub.Rooms[msg.RoomID] = room
		hub.Broadcast <- msg
	}
}
