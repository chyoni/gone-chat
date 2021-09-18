package chat

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Client struct {
	participants map[string]*websocket.Conn
	m            sync.Mutex
}

var clients *Client = &Client{
	participants: make(map[string]*websocket.Conn),
}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleError(err)

	clients.m.Lock()
	clients.participants[conn.RemoteAddr().String()] = conn
	clients.m.Unlock()

	for {
		mType, payload, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			conn.Close()
		}
		for _, client := range clients.participants {
			client.WriteMessage(mType, payload)
		}
	}
}
