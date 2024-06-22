package logic

import (
	"fmt"
)

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

var MaxLengthOldMsg = 50

func NewHub() *Hub {
	rooms := make(map[string]*Room)
	rooms["room123"] = &Room{
		ID:      "room123",
		Name:    "Tám chuyện xuyên biên giới",
		Clients: make(map[string]*Client),
		LastMsg: "Room was created at 2024-06-20 10:42:22.6076085 +0700 +07 m=+45.546873001",
		Public:  1,
	}
	rooms["room456"] = &Room{
		ID:      "room456",
		Name:    "Coding interview",
		Clients: make(map[string]*Client),
		LastMsg: "Room was created at 2024-06-20 10:42:22.6076085 +0700 +07 m=+45.546873001",
		Public:  1,
	}
	return &Hub{
		Rooms:      rooms,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			fmt.Println("Have a new member joined room: ", cl.RoomID, " with name: ", cl.Username)
			if _, ok := h.Rooms[cl.RoomID]; ok {
				r := h.Rooms[cl.RoomID]
				memebers := r.Members
				var exist bool = false
				for _, mem := range memebers {
					if mem == cl.Username {
						exist = true
					}
				}
				if !exist {
					memebers = append(memebers, cl.Username)
					r.Members = memebers
				}

				if client, ok := r.Clients[cl.ID]; !ok || client.Conn == nil {
					r.Clients[cl.ID] = cl
				}
				for _, msg := range h.Rooms[cl.RoomID].OldMsg {
					cl.Message <- msg
				}

			}
		case cl := <-h.Unregister:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok {
					// if len(h.Rooms[cl.RoomID].Clients) != 0 {
					// h.Broadcast <- &Message{
					// 	Content:  cl.Username + " left the chat! at " + time.Now().String(),
					// 	RoomID:   cl.RoomID,
					// 	Username: cl.Username,
					// }
					// }

					delete(h.Rooms[cl.RoomID].Clients, cl.ID)
					close(cl.Message)
				}
			}

		case m := <-h.Broadcast:
			fmt.Println("Message broadcast: ", m.Content)
			if room, ok := h.Rooms[m.RoomID]; ok {

				for _, cl := range h.Rooms[m.RoomID].Clients {
					cl.Message <- m
				}
				if len(room.OldMsg) < MaxLengthOldMsg {
					room.OldMsg = append(room.OldMsg, m)
				} else if len(room.OldMsg) == MaxLengthOldMsg {
					room.OldMsg = append(room.OldMsg[1:], m)
				}
				h.Rooms[m.RoomID] = room
			}
		}
	}
}
