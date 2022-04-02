package controller

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"pipi.com/gogf/pipi-gf-demo/internal/service/ws"
	"pipi.com/gogf/pipi-gf-demo/internal/utils"
)

var WebSocket = cWebSocket{}

type cWebSocket struct{}

func (c *cWebSocket) Upgrade(r *ghttp.Request) {
	websocket, err := r.WebSocket()
	if err != nil {
		// glog.Error(err)
		r.Exit()
	}
	ctx, cancel := context.WithCancel(context.Background())
	unique := "ws-" + utils.RandString(7)
	ws.CreateConnection(unique, websocket, ctx, cancel)
}
