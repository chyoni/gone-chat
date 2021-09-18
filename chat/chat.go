package chat

import (
	"net/http"

	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func UpgradeWithRoom(rw http.ResponseWriter, r *http.Request) {
	ids := mux.Vars(r)
	roomID, ok := ids["roomID"]
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return ok
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleError(err)

	initMember(conn, roomID)
}
