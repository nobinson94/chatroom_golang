package main

import (
	"app/chat"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	fmt.Println("START TEST")
	go chat.Server()

	time.Sleep(3 * time.Second)

	service := "0.0.0.0:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error : %s", err.Error())
	}
	for i := 0; i < 20; i++ {
		go func(idx int) {
			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Fatal error : %s", err.Error())
			}
			msg := "(" + strconv.Itoa(idx) + ") This is My Msg\n" + time.Now().String()
			_, err = conn.Write([]byte(msg))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Fatal error : %s", err.Error())
			}
		}(i)
	}
	time.Sleep(3 * time.Second)
}
