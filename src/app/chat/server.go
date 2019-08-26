package chat

import (
	"fmt" // package for formatted IO
	"net" // package for network interface
	"os"
	"strconv"
	"sync"
	"time"
)

type Msg struct {
	currentConn net.Conn
	msg         []byte
}

func Server() {
	var fileMutex sync.Mutex
	newConns := make(chan net.Conn, 128)  // 새로 연결되는 Client 받는 채널
	deadConns := make(chan net.Conn, 128) // 연결이 끊긴 Client 받는 채널
	delivery := make(chan Msg, 128)       // 연결된 다른 clinet들에게 msg를 전달하는 채널
	connList := make(map[net.Conn]bool)   // 연결이 되어있는 Client 소켓을 저장하는 리스트

	userNum := 0
	service := "0.0.0.0:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	defer listener.Close()

	fmt.Println("Chat Server Open")
	t := time.Now()

	// ho, mi, sec := Clock(t)
	// year, month, day := Date(t)
	fileName := "chat_log/sentbeChat-" + t.String() + ".txt"
	txtFo, err := os.Create(fileName)
	checkError(err)
	txtFo.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			//connList = append(connList, conn)
			if err != nil {
				panic(err)
			}
			newConns <- conn
		}
	}()

	for {
		select {
		case newConn := <-newConns: // 새로운 연결
			remoteAddr := newConn.RemoteAddr().String()
			connList[newConn] = true
			userNum++
			newUserAlertMsg := "New User [" + remoteAddr + "] Join .. (" + strconv.Itoa(userNum) + ") Users\n"
			fmt.Print(newUserAlertMsg)
			deliverMsg(newConn, []byte(newUserAlertMsg), delivery)

			go func() { // 새로 연결된 Connection과 계속해서 통신
				cConn := newConn
				fo, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
				if err != nil {
					fmt.Print("FILE OPEN ERROR")
					return
				}
				for {
					buf := make([]byte, 1024)
					n, err := cConn.Read(buf)
					if err != nil {
						deadConns <- cConn
						break
					} else {
						recordStr := remoteAddr + ": " + string(buf[0:n])
						recordBuf := []byte(recordStr + "\n")
						fmt.Println(recordStr)
						fileMutex.Lock()
						_, err = fo.Write(recordBuf)
						fileMutex.Unlock()
						checkError(err)
						deliverMsg(cConn, recordBuf, delivery)
					}
				}
				fo.Close()
			}()

		case deadConn := <-deadConns: // 종료된 Connection
			remoteAddr := deadConn.RemoteAddr().String()
			_ = deadConn.Close()
			delete(connList, deadConn)
			userNum--
			deadUserAlertMsg := "User [" + remoteAddr + "] Left .. (" + strconv.Itoa(userNum) + ") Users\n"
			fmt.Println(deadUserAlertMsg)

			deliverMsg(deadConn, []byte(deadUserAlertMsg), delivery)

		case dMsg := <-delivery: // 다른 client에 msg전달
			deliveredBuf := dMsg.msg
			currentConn := dMsg.currentConn
			for otherConn, _ := range connList {
				if otherConn != currentConn {
					totalWritten := 0
					for totalWritten < len(deliveredBuf) {
						writtenThisCall, err := otherConn.Write(deliveredBuf[totalWritten:])
						if err != nil {
							deadConns <- otherConn
							break
						}
						totalWritten += writtenThisCall
					}
				}
			}
		}
	}
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error : %s", err.Error())
		os.Exit(1)
	}
}

func deliverMsg(c net.Conn, b []byte, deliverChan chan Msg) {
	dm := new(Msg)
	dm.currentConn = c
	dm.msg = b
	deliverChan <- *dm
}
