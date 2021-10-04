package chat

import (
	"fmt"

	"github.com/chiwon99881/gone-chat/entity"
	"github.com/chiwon99881/gone-chat/utils"
)

func HandleMessage(message, roomID string, user *entity.User, created int) {
	rs.m.Lock()
	defer rs.m.Unlock()

	room, ok := rs.IDs[roomID]
	if !ok {
		fmt.Println("error")
	}
	m := &payload{
		RoomID:  utils.ToUintFromString(roomID),
		From:    user,
		Message: message,
		Created: created,
	}
	mBytes := utils.ToBytes(m)
	for _, member := range room.members {
		member.hub <- mBytes
	}
}
