package chat

import (
	"fmt"

	"github.com/chiwon99881/gone-chat/utils"
)

func HandleMessage(message, roomID string) {
	rs.m.Lock()
	defer rs.m.Unlock()

	room, ok := rs.IDs[roomID]
	if !ok {
		fmt.Println("error")
	}
	m := &payload{
		RoomID:  roomID,
		Message: message,
	}
	mBytes := utils.ToBytes(m)
	for _, member := range room.members {
		member.hub <- mBytes
	}
}
