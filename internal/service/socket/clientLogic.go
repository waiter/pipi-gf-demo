package socket

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var clientLogic = map[string]func(data *gjson.Json, c *ClientConnection) (g.Map, error){
	"testAdd":  testAddClient,
	"testPush": testClientPush,
}

func CallClientLogic(cmd string, data *gjson.Json, c *ClientConnection) (g.Map, error) {
	if clientLogic[cmd] != nil {
		return clientLogic[cmd](data, c)
	}
	return nil, gerror.New("can't do cmd: " + cmd)
}

func testAddClient(data *gjson.Json, c *ClientConnection) (g.Map, error) {
	info := data.Get("data").Ints()
	back := g.Map{}
	back["cmd"] = "add"
	sum := 0
	for _, v := range info {
		sum += v
	}
	back["data"] = sum
	return back, nil
}

func testClientPush(data *gjson.Json, c *ClientConnection) (g.Map, error) {
	back := g.Map{}
	back["cmd"] = "testPush"
	back["data"] = "yyyyyyyyy"
	SocketManager.Push2Web(c, back)
	return nil, nil
}
