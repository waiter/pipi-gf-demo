package socket

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var webLogic = map[string]func(data *gjson.Json, c *WebConnection) (g.Map, error){
	"testAdd": testAdd,
}

func CallWebLogic(cmd string, data *gjson.Json, c *WebConnection) (g.Map, error) {
	if webLogic[cmd] != nil {
		return webLogic[cmd](data, c)
	}
	return nil, gerror.New("can't do cmd: " + cmd)
}

func testAdd(data *gjson.Json, c *WebConnection) (g.Map, error) {
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