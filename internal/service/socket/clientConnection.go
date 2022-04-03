package socket

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gmlock"
	"github.com/gogf/gf/v2/os/gtimer"
	"pipi.com/gogf/pipi-gf-demo/internal/consts"
)

type ClientConnection struct {
	Socket    *gtcp.Conn
	closed    bool
	Unique    string
	Name      string
	Ctx       context.Context
	CtxClose  func()
	HeartBeat *gtimer.Entry
}

func CreateClientConnection(unique string, socket *gtcp.Conn, ctx context.Context, ctxClose func(), name string) *ClientConnection {
	client := &ClientConnection{
		Socket:    socket,
		closed:    false,
		Unique:    unique,
		Name:      name,
		Ctx:       ctx,
		CtxClose:  ctxClose,
		HeartBeat: nil,
	}
	client.HeartBeat = clientheartbeat(client)
	SocketManager.AddClientConnection(client)
	go client.read()
	client.log("connected")
	return client
}

func clientheartbeat(client *ClientConnection) *gtimer.Entry {
	return gtimer.AddSingleton(client.Ctx, time.Second*consts.ClientSocketHeartBeatTime, func(ctx context.Context) {
		client.log("heart beat close")
		client.Close()
	})
}

func (c *ClientConnection) Close() {
	if gmlock.TryLock(c.Unique) {
		if !c.closed {
			_ = c.Socket.Close()
			c.closed = true
			//必须优先关闭读写协程
			c.CtxClose()
			c.HeartBeat.Close()
			SocketManager.RemoveClientConnection(c)
			c.log("closed")
		}
		gmlock.Unlock(c.Unique)
	}
}

func (c *ClientConnection) read() {
	for {
		msg, err := c.Socket.RecvLine()
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
		data := gjson.New(msg)
		if data.IsNil() {
			c.writeError("empty:" + string(msg))
			continue
		}
		cmd := data.Get("cmd", "").String()
		if len(cmd) == 0 {
			c.writeError("no cmd:" + string(msg))
			continue
		}
		result, err := CallClientLogic(cmd, data, c)
		if err != nil {
			c.writeError(err.Error())
			continue
		}
		if result != nil {
			c.WritePack(result)
		}
	}
}

func (c *ClientConnection) write(msg []byte) {
	msg = append(msg, []byte("\n")...)
	err := c.Socket.Send(msg)
	if err != nil {
		c.log("write error")
		c.Close()
	}
}

func (c *ClientConnection) log(msg string) {
	glog.Print(c.Ctx, "【ClientSocket Connection】"+c.Unique+": "+msg)
}

func (c *ClientConnection) WritePack(data g.Map) {
	data["pack"] = time.Now().Unix()
	data["unique"] = c.Unique
	encode, err := gjson.Encode(data)
	if err != nil {
		fmt.Println(err)
	}
	c.write(encode)
}

func (c *ClientConnection) writeError(msg string) {
	data := make(map[string]interface{})
	data["cmd"] = "error"
	data["msg"] = msg
	c.WritePack(data)
}
