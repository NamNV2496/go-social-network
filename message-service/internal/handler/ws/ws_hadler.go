package wsHandler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/message-service/internal/logic"
)

type Handler struct {
	hub *logic.Hub
}

func NewHandler(h *logic.Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

func (h *Handler) CreateRoom(c *gin.Context) {
	roomId := c.Query("roomId")
	roomName := c.Query("roomName")
	h.hub.Rooms[roomId] = &logic.Room{
		ID:      roomId,
		Name:    roomName,
		Clients: make(map[string]*logic.Client),
		LastMsg: fmt.Sprintf("Room was created at %v", time.Now()),
	}
}

func (h *Handler) JoinRoom(c *gin.Context) {
	newClient, err := logic.CreateNewClient(c)
	if err != nil {
		fmt.Println("Error when create new client")
		return
	}
	notifyMsg := &logic.Message{
		Content:  "Hello, I'm new member. Nice to join room!",
		RoomID:   newClient.RoomID,
		Username: newClient.Username,
	}
	// register new client
	h.hub.Register <- newClient
	// send a broadcast message
	h.hub.Broadcast <- notifyMsg

	// write new message
	go newClient.WriteMessage()
	// read message and display on screen
	newClient.ReadMessage(h.hub)

	c.JSON(http.StatusOK, "Thành công")
}

type RoomRes struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	LastMsg string `json:"last_msg"`
}

func (h *Handler) GetRooms(c *gin.Context) {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:      r.ID,
			Name:    r.Name,
			LastMsg: r.LastMsg,
		})
	}

	c.JSON(http.StatusOK, rooms)
}
