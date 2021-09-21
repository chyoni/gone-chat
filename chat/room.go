package chat

import (
	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/websocket"
)

var rs *rooms = &rooms{
	IDs: make(map[string]*room),
}

var ps *participants = &participants{
	participants: make(map[string]*participant),
}

func (p *participant) write() {
	defer ps.m.Unlock()
	for {
		m, ok := <-p.hub
		if !ok {
			ps.m.Lock()
			p.conn.Close()
			delete(ps.participants, p.userID)
			break
		}

		message := &payload{}
		utils.FromBytes(message, m)
		err := p.conn.WriteJSON(message)
		if err != nil {
			ps.m.Lock()
			p.conn.Close()
			delete(ps.participants, p.userID)
			break
		}
	}
}

func initRoom(conn *websocket.Conn, roomID, userID string) {
	ps.m.Lock()
	rs.m.Lock()
	defer ps.m.Unlock()
	defer rs.m.Unlock()

	p := &participant{
		userID: userID,
		conn:   conn,
		hub:    make(chan []byte),
	}
	r := &room{
		ID:      roomID,
		members: make(map[string]*participant),
	}

	room, exist := rs.IDs[roomID]
	if exist {
		room.members[p.userID] = p
	} else {
		r.members[p.userID] = p
		rs.IDs[roomID] = r
	}

	ps.participants[p.userID] = p
	go p.write()
}
