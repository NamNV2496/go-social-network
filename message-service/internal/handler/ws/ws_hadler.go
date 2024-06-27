package wsHandler

import (
	"fmt"
	"log"
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

type CreateRoomRequest struct {
	RoomId   string   `json:"roomId"`
	RoomName string   `json:"roomName"`
	Public   int      `json:"public"`
	Member   []string `json:"member"`
}

func (h *Handler) CreateRoom(c *gin.Context) {

	var roomRequest CreateRoomRequest
	if err := c.BindJSON(&roomRequest); err != nil {
		log.Fatalln("Invalid input")
		return
	}

	if _, ok := h.hub.Rooms[roomRequest.RoomId]; ok {
		return
	}
	h.hub.Rooms[roomRequest.RoomId] = &logic.Room{
		ID:      roomRequest.RoomId,
		Name:    roomRequest.RoomName,
		Clients: make(map[string]*logic.Client),
		LastMsg: fmt.Sprintf("Room was created at %v", time.Now()),
		Public:  roomRequest.Public,
		Members: roomRequest.Member,
	}
}

func (h *Handler) JoinRoom(c *gin.Context) {
	newClient, err := logic.CreateNewClient(c)
	if err != nil {
		log.Println("Error when create new client")
		return
	}
	// notifyMsg := &logic.Message{
	// 	Content:  "Hello, I'm new member. Nice to join room!",
	// 	RoomID:   newClient.RoomID,
	// 	Username: newClient.Username,
	// }
	// register new client
	h.hub.Register <- newClient
	// send a broadcast message
	// h.hub.Broadcast <- notifyMsg

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
	username := c.Query("username")
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		if r.Public == 1 {
			rooms = append(rooms, RoomRes{
				ID:      r.ID,
				Name:    r.Name,
				LastMsg: r.LastMsg,
			})
		} else {
			for _, member := range r.Members {
				if member == username {
					rooms = append(rooms, RoomRes{
						ID:      r.ID,
						Name:    r.Name,
						LastMsg: r.LastMsg,
					})
				}
			}
		}
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *Handler) GetMembers(c *gin.Context) {
	username := c.Query("username")
	roomId := c.Query("roomId")

	room := h.hub.Rooms[roomId]
	for _, member := range room.Members {
		if member == username {
			c.JSON(http.StatusOK, gin.H{"member": room.Members})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{})
}
