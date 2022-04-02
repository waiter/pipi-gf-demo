package ws

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gmlock"
	"github.com/gogf/gf/v2/os/gtimer"
)

type Connection struct {
	Socket    *ghttp.WebSocket
	closed    bool
	Unique    string
	Ctx       context.Context
	CtxClose  func()
	HeartBeat *gtimer.Entry
}

func CreateConnection(unique string, socket *ghttp.WebSocket, ctx context.Context, ctxClose func()) *Connection {
	client := &Connection{
		Socket:    socket,
		closed:    false,
		Unique:    unique,
		Ctx:       ctx,
		CtxClose:  ctxClose,
		HeartBeat: nil,
	}
	client.HeartBeat = HeartBeat(client)
	go client.read()
	client.log("connected")
	return client
}

func HeartBeat(client *Connection) *gtimer.Entry {
	return gtimer.AddSingleton(client.Ctx, time.Second*60, func(ctx context.Context) {
		client.log("heart beat close")
		client.Close()
	})
}

func (c *Connection) Close() {
	if gmlock.TryLock(c.Unique) {
		if !c.closed {
			_ = c.Socket.Close()
			c.closed = true
			//必须优先关闭读写协程
			c.CtxClose()
			c.HeartBeat.Close()
			c.log("closed")
		}
		gmlock.Unlock(c.Unique)
	}
}

func (c *Connection) read() {
	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			c.log("read error")
			c.Close()
			return
		}
		c.HeartBeat.Reset()
		if strings.ToLower(string(msg)) == "ping" {
			c.write([]byte("pong"))
			continue
		}
		c.write(msg)
	}
}

func (c *Connection) write(msg []byte) {
	err := c.Socket.WriteMessage(1, msg)
	if err != nil {
		c.log("write error")
		c.Close()
	}
}

func (c *Connection) log(msg string) {
	glog.Print(c.Ctx, "【WebSocket Connection】" + c.Unique + ": " + msg)
}