package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/gookit/color"
)

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	connection, err := net.Dial("tcp", "localhost:8080")
	logFatal(err)

	defer connection.Close()
	color.Cyan.Print("Enter your name: ")
	// color.Cyan.Print("Enter you name: ")
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	logFatal(err)
	username = strings.Trim(username, " \r\n")
	color.Green.Printf("%s, Welcome  to the party\n==============================\n\n", username)
	connection.Write([]byte(fmt.Sprintf("%s is joins to the party\n", username)))

	go read(connection)
	write(connection, username, reader)
}

func read(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		incomingMsg, err := reader.ReadString('\n')
		if err == io.EOF {
			conn.Close()
			fmt.Println("Connection closed")
			os.Exit(0)
		}
		color.Red.Printf("%s", incomingMsg) //------------------------------\n
	}
}

func write(conn net.Conn, username string, reader *bufio.Reader) {
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		msg = fmt.Sprintf("[%s]: %s\n", username, strings.Trim(msg, " \r\n"))
		conn.Write([]byte(msg))
	}
}
