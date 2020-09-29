package main

import (
	"TCPClient/model"
	"log"
	"net"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	server := model.StartServer()

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}
	defer listener.Close()
	log.Println("Listening on localhost:8081.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		server.Join(model.NewUser(conn))
	}
}
