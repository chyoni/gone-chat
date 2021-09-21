package chat

import (
	"sync"

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
	RoomID  string `json:"roomId"`
	FromID  string `json:"fromId"`
	Message string `json:"message"`
}
