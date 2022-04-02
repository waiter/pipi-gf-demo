package socket

import "github.com/gogf/gf/v2/frame/g"

type SocketManagerStruct struct {
	WebSocket map[string]*WebConnection
}

var (
	SocketManager = &SocketManagerStruct{
		WebSocket: map[string]*WebConnection{},
	}
)

func (s *SocketManagerStruct) BroadcastMsg(data g.Map) {
	for _, c := range s.WebSocket {
		c.WritePack(data)
	}
}

func (s *SocketManagerStruct) AddWebConnection(conn *WebConnection) {
	s.WebSocket[conn.Unique] = conn
}

func (s *SocketManagerStruct) RemoveWebConnection(conn *WebConnection) {
	delete(s.WebSocket, conn.Unique)
}