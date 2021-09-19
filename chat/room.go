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
			delete(ps.participants, p.conn.RemoteAddr().String())
			break
		}

		message := &payload{}
		utils.FromBytes(message, m)
		err := p.conn.WriteJSON(message.Message)
		if err != nil {
			ps.m.Lock()
			p.conn.Close()
			delete(ps.participants, p.conn.RemoteAddr().String())
			break
		}
	}
}

func initMember(conn *websocket.Conn, roomID string) {
	ps.m.Lock()
	rs.m.Lock()
	defer ps.m.Unlock()
	defer rs.m.Unlock()

	p := &participant{
		conn: conn,
		hub:  make(chan []byte),
	}
	r := &room{
		ID:      roomID,
		members: make(map[string]*participant),
	}

	room, exist := rs.IDs[roomID]
	if exist {
		room.members[p.conn.RemoteAddr().String()] = p
	} else {
		r.members[p.conn.RemoteAddr().String()] = p
		rs.IDs[roomID] = r
	}

	ps.participants[p.conn.RemoteAddr().String()] = p
	go p.write()
}
