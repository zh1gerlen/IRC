package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	openConnections = make(map[net.Conn]bool)
	newConnection   = make(chan net.Conn)
	deadConnection  = make(chan net.Conn)
)

func main() {
	lnr, err := net.Listen("tcp", ":8080")
	logFatal(err)
	fmt.Println("listen tcp:8080")
	defer lnr.Close()

	go func() {
		for {
			conn, err := lnr.Accept()
			logFatal(err)

			openConnections[conn] = true
			newConnection <- conn
		}
	}()

	for {
		select {
		case conn := <-newConnection:
			go broadcastMessage(conn)
		case conn := <-deadConnection:
			delete(openConnections, conn)
		}
	}

}

func broadcastMessage(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		// write message to all other connections
		for item := range openConnections {
			if item == conn {
				continue
			}
			item.Write([]byte(msg))
		}
	}
	deadConnection <- conn
}
