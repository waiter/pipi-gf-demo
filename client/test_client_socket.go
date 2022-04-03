package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/net/gtcp"
)

func main() {
	if conn, err := gtcp.NewPoolConn("127.0.0.1:8187"); err == nil {
		if b, err := conn.SendRecv([]byte("ping"), -1); err == nil {
			fmt.Println(string(b), conn.LocalAddr(), conn.RemoteAddr())
		} else {
			fmt.Println(err)
		}

		time.Sleep(time.Second * 5)

		if b, err := conn.SendRecv([]byte("{\"cmd\":\"testAdd\",\"data\":[1,83,3]}"), -1); err == nil {
			fmt.Println(string(b), conn.LocalAddr(), conn.RemoteAddr())
		} else {
			fmt.Println(err)
		}

		conn.Close()
	} else {
		fmt.Println(err)
	}
}
