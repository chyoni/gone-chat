package chat

import (
	"fmt"
	"net/http"

	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleError(err)

	mType, payload, err := conn.ReadMessage()
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
	}
	conn.WriteMessage(mType, payload)
}
