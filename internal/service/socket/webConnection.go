package socket

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gmlock"
	"github.com/gogf/gf/v2/os/gtimer"
	"pipi.com/gogf/pipi-gf-demo/internal/consts"
)

type WebConnection struct {
	Socket    *ghttp.WebSocket
	closed    bool
	Unique    string
	Ctx       context.Context
	CtxClose  func()
	HeartBeat *gtimer.Entry
}

func CreateWebConnection(unique string, socket *ghttp.WebSocket, ctx context.Context, ctxClose func()) *WebConnection {
	client := &WebConnection{
		Socket:    socket,
		closed:    false,
		Unique:    unique,
		Ctx:       ctx,
		CtxClose:  ctxClose,
		HeartBeat: nil,
	}
	client.HeartBeat = webheartbeat(client)
	SocketManager.AddWebConnection(client)
	go client.read()
	client.log("connected")
	return client
}

func webheartbeat(client *WebConnection) *gtimer.Entry {
	return gtimer.AddSingleton(client.Ctx, time.Second*consts.WebSocketHeartBeatTime, func(ctx context.Context) {
		client.log("heart beat close")
		client.Close()
	})
}

func (c *WebConnection) Close() {
	if gmlock.TryLock(c.Unique) {
		if !c.closed {
			_ = c.Socket.Close()
			c.closed = true
			//必须优先关闭读写协程
			c.CtxClose()
			c.HeartBeat.Close()
			SocketManager.RemoveWebConnection(c)
			c.log("closed")
		}
		gmlock.Unlock(c.Unique)
	}
}

func (c *WebConnection) read() {
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
		result, err := CallWebLogic(cmd, data, c)
		if err != nil {
			c.writeError(err.Error())
			continue
		}
		if result != nil {
			c.WritePack(result)
		}
	}
}

func (c *WebConnection) write(msg []byte) {
	err := c.Socket.WriteMessage(1, msg)
	if err != nil {
		c.log("write error")
		c.Close()
	}
}

func (c *WebConnection) log(msg string) {
	glog.Print(c.Ctx, "【WebSocket Connection】"+c.Unique+": "+msg)
}

func (c *WebConnection) WritePack(data g.Map) {
	data["pack"] = time.Now().Unix()
	data["unique"] = c.Unique
	encode, err := gjson.Encode(data)
	if err != nil {
		fmt.Println(err)
	}
	c.write(encode)
}

func (c *WebConnection) writeError(msg string) {
	data := make(map[string]interface{})
	data["cmd"] = "error"
	data["msg"] = msg
	c.WritePack(data)
}
