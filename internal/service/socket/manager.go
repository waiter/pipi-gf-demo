package socket

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type SocketManagerStruct struct {
	WebSocket    map[string]*WebConnection
	ClientSocket map[string]*ClientConnection
	Web2Client   map[string]string
	Client2Web   map[string]map[string]int
}

var (
	SocketManager = &SocketManagerStruct{
		WebSocket:    map[string]*WebConnection{},
		ClientSocket: map[string]*ClientConnection{},
		Web2Client:   map[string]string{},
		Client2Web:   map[string]map[string]int{},
	}
)

func (s *SocketManagerStruct) BroadcastMsg2Web(data g.Map) {
	for _, c := range s.WebSocket {
		c.WritePack(data)
	}
}

func (s *SocketManagerStruct) AddWebConnection(conn *WebConnection) {
	s.WebSocket[conn.Unique] = conn
	conn.WritePack(s.allClientConnections())
}

func (s *SocketManagerStruct) RemoveWebConnection(conn *WebConnection) {
	if w, ok := s.Web2Client[conn.Unique]; ok {
		if _, ok := s.Client2Web[w]; ok {
			delete(s.Client2Web[w], conn.Unique)
		}
		delete(s.Web2Client, conn.Unique)
	}
	delete(s.WebSocket, conn.Unique)
}

func (s *SocketManagerStruct) AddClientConnection(conn *ClientConnection) {
	s.ClientSocket[conn.Unique] = conn
	s.BroadcastMsg2Web(s.allClientConnections())
}

func (s *SocketManagerStruct) RemoveClientConnection(conn *ClientConnection) {
	if c, ok := s.Client2Web[conn.Unique]; ok {
		for k := range c {
			delete(s.Web2Client, k)
		}
		delete(s.Client2Web, conn.Unique)
	}
	delete(s.ClientSocket, conn.Unique)
	s.BroadcastMsg2Web(s.allClientConnections())
}

func (s *SocketManagerStruct) allClientConnections() g.Map {
	back := g.Map{}
	back["cmd"] = "updateClients"
	data := make(map[string]string)
	for k, v := range s.ClientSocket {
		data[k] = v.Name
	}
	back["data"] = data
	return back
}

func (s *SocketManagerStruct) BindClient(webUnique string, clientUnique string) error {
	if _, ok := s.ClientSocket[clientUnique]; !ok {
		return gerror.New("Bind Failed, can't find: " + clientUnique)
	}
	s.Web2Client[webUnique] = clientUnique
	if v, ok := s.Client2Web[clientUnique]; ok {
		v[webUnique] = 1
	} else {
		data := make(map[string]int)
		data[webUnique] = 1
		s.Client2Web[clientUnique] = data
	}
	return nil
}

func (s *SocketManagerStruct) Push2Web(conn *ClientConnection, data g.Map) {
	if v, ok := s.Client2Web[conn.Unique]; ok {
		for k := range v {
			if w, ok := s.WebSocket[k]; ok {
				w.WritePack(data)
			}
		}
	}
}

func (s *SocketManagerStruct) Push2Client(conn *WebConnection, data g.Map) {
	if v, ok := s.Web2Client[conn.Unique]; ok {
		if s, ok := s.ClientSocket[v]; ok {
			s.WritePack(data)
		}
	}
}
