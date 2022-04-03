package controller

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtcp"
	"pipi.com/gogf/pipi-gf-demo/internal/service/socket"
	"pipi.com/gogf/pipi-gf-demo/internal/utils"
)

var ClientSocket = cClientSocket{}

type cClientSocket struct{}

func (c *cClientSocket) Start() {
	address, _ := g.Cfg().Get(context.Background(), "socket.address", "")
	fmt.Println("start socket server: " + address.String())
	go gtcp.NewServer(address.String(), func(conn *gtcp.Conn) {
		ctx, cancel := context.WithCancel(context.Background())
		unique := "cs-" + utils.RandString(7)
		socket.CreateClientConnection(unique, conn, ctx, cancel, "unknown")
	}).Run()
}
