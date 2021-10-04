package chat

import (
	"sync"

	"github.com/chiwon99881/gone-chat/entity"
	"github.com/gorilla/websocket"
)

type rooms struct {
	IDs map[string]*room
	m   sync.Mutex
}

type room struct {
	ID      string
	members map[string]*participant
}

type participants struct {
	participants map[string]*participant
	m            sync.Mutex
}

type participant struct {
	userID string
	conn   *websocket.Conn
	hub    chan []byte
}

type payload struct {
	RoomID  uint         `json:"roomId"`
	From    *entity.User `json:"from"`
	Message string       `json:"message"`
	Created int          `json:"created"`
}
