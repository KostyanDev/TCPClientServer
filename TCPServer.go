package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {

	fmt.Println("Launching server...")

	ln, _ := net.Listen("tcp", ":8088")

	conn, _ := ln.Accept()

	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message Received:", string(message))
		newmessage := strings.ToUpper(message)
		conn.Write([]byte(newmessage + "\n"))
	}
}
