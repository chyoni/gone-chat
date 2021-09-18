package chat

import (
	"sync"

	"github.com/chiwon99881/gone-chat/utils"
	"github.com/gorilla/websocket"
)

type Rooms struct {
	IDs map[string]*room
	m   sync.Mutex
}

type room struct {
	ID      string
	members map[string]*participant
}

type Participants struct {
	participants map[string]*participant
	m            sync.Mutex
}

type participant struct {
	conn *websocket.Conn
	hub  chan []byte
}

var Rs *Rooms = &Rooms{
	IDs: make(map[string]*room),
}

var Ps *Participants = &Participants{
	participants: make(map[string]*participant),
}

func (p *participant) write() {
	defer Ps.m.Unlock()
	for {
		m, ok := <-p.hub
		if !ok {
			Ps.m.Lock()
			p.conn.Close()
			delete(Ps.participants, p.conn.RemoteAddr().String())
			break
		}

		message := &Message{}
		utils.FromBytes(message, m)
		err := p.conn.WriteJSON(message.Message)
		if err != nil {
			Ps.m.Lock()
			p.conn.Close()
			delete(Ps.participants, p.conn.RemoteAddr().String())
			break
		}
	}
}

func initMember(conn *websocket.Conn, roomID string) {
	Ps.m.Lock()
	Rs.m.Lock()
	defer Ps.m.Unlock()
	defer Rs.m.Unlock()

	p := &participant{
		conn: conn,
		hub:  make(chan []byte),
	}
	r := &room{
		ID:      roomID,
		members: make(map[string]*participant),
	}

	room, exist := Rs.IDs[roomID]
	if exist {
		room.members[p.conn.RemoteAddr().String()] = p
	} else {
		r.members[p.conn.RemoteAddr().String()] = p
		Rs.IDs[roomID] = r
	}

	Ps.participants[p.conn.RemoteAddr().String()] = p
	go p.write()
}
