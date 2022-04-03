package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/net/gtcp"
)

func main() {
	if conn, err := gtcp.NewPoolConn("127.0.0.1:8187"); err == nil {
		go func() {
			for {
				b, err := conn.RecvLine()
				if err != nil {
					break
				}
				fmt.Println(string(b), conn.LocalAddr(), conn.RemoteAddr())
			}
		}()
		for {
			if err := conn.Send([]byte("ping\n")); err != nil {
				fmt.Println(err)
				break
			}

			time.Sleep(time.Second * 5)

			if err := conn.Send([]byte("{\"cmd\":\"testAdd\",\"data\":[1,83,3]}\n")); err != nil {
				fmt.Println(err)
				break
			}

			time.Sleep(time.Second * 5)

			if err := conn.Send([]byte("{\"cmd\":\"testPush\",\"data\":[1,83,3]}\n")); err != nil {
				fmt.Println(err)
				break
			}
		}
		conn.Close()
	} else {
		fmt.Println(err)
	}
}
