package client

import (
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/sadeepa24/netshoot/server"
)

func TestSpeedTest(t *testing.T) {
	// ls, _ := server.NewMixedLs(config.LsConfig{
	// 	ListenAddr: "127.0.0.1:3033",
	// })

	laddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:3033")
	ls,_ := net.ListenTCP("tcp", laddr)

	go func(){ for {
		conn, err := ls.Accept()
		if err != nil {
			fmt.Println("server err " +err.Error())
			continue
		}
		go server.Speedtest(conn)
	}}()

	for i := 0; i < 10; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1:3033")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("connection dialed")
		test, err := speedtest(conn, 2)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(test)
		fmt.Println()
	}

}