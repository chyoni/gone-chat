package chat

import (
	"fmt"

	"github.com/chiwon99881/gone-chat/utils"
)

type Message struct {
	RoomID  string `json:"roomId"`
	Message string `json:"message"`
}

func HandleMessage(message, roomID string) {
	Rs.m.Lock()
	defer Rs.m.Unlock()

	room, ok := Rs.IDs[roomID]
	if !ok {
		fmt.Println("error")
	}
	m := &Message{
		RoomID:  roomID,
		Message: message,
	}
	mBytes := utils.ToBytes(m)
	for _, member := range room.members {
		member.hub <- mBytes
	}
}
