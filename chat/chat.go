package chat

import (
	"net/http"

	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	_, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleError(err)
}
