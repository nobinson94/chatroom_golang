package chat

import (
	"bufio"
	"fmt" // package for formatted IO
	"net" // package for network interface
	"os"
)

func Client(service string) {
	// if len(os.Args) != 2 {
	// 	fmt.Println("Error: Arguments Not Enough")
	// 	os.Exit(1)
	// }
	// service := os.Args[1]

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErrorNum(err, 1)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkErrorNum(err, 2)

	fmt.Println("Welcome to SENTBE CHAT Channel")
	fmt.Print("YOU: ")
	go func() {
		serverReader := bufio.NewReader(conn)
		//buf := make([]byte, 1024)
		//fmt.Print("\r\x1b[2K")
		for {
			//buf := make([]byte, 1024)
			serverLine, err := serverReader.ReadBytes('\n')
			checkErrorNum(err, 3)
			fmt.Print("\r\x1b[2K")
			fmt.Print(string(serverLine))
			fmt.Print("YOU:")
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		msg, err := reader.ReadBytes('\n')
		checkErrorNum(err, 4)
		if string(msg) != "" {
			//msg = msg[:len(msg)-1]
			_, err = conn.Write(msg)
			checkErrorNum(err, 5)
			fmt.Print("YOU: ")
		} else {
			fmt.Print("\r\x1b[2K")
			fmt.Print("YOU: ")
		}
	}
}

func checkErrorNum(err error, num int) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "(%d)Fatal error : %s", num, err.Error())
		os.Exit(1)
	}
}
