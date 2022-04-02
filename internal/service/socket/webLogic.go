package socket

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var logic = map[string]func(data *gjson.Json, c *WebConnection) (g.Map, error){
	"testAdd": TestAdd,
}

func Call(cmd string, data *gjson.Json, c *WebConnection) (g.Map, error) {
	if logic[cmd] != nil {
		return logic[cmd](data, c)
	}
	return nil, gerror.New("can't do cmd: " + cmd)
}

func TestAdd(data *gjson.Json, c *WebConnection) (g.Map, error) {
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
