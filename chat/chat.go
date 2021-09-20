package chat

import (
	"fmt"
	"net/http"

	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func UpgradeWithRoom(rw http.ResponseWriter, r *http.Request) {
	ids := mux.Vars(r)
	roomID, ok := ids["roomID"]
	fmt.Println(roomID, ok)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return ok
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleError(err)

	initRoom(conn, roomID)
}

// room을 만들어 -> 나는 그 room에 커넥션이 생기지?, db에도 저장이 돼 -> 근데 내가 만약에 앱을 종료해 그럼 지금상태에서는 conn이 끊겨 그럼, 로그인을 하면 db에서 본인의 room정보들을 받아서
// 모든 connection을 다시 연결해 upgradewithroom을 http.get으로 때리면 되고 header에 내 정보를 담아서 upgradewithroom에서 그거를 받아와서 initRoom할때 던져줘야겠다 그래야 p에 현재 conn의 유저를 저장할 수 있으니까
